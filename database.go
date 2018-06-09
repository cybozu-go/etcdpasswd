package etcdpasswd

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"

	"github.com/coreos/etcd/clientv3"
	"github.com/cybozu-go/log"
)

// Database is a on-memory snapshot of users and groups in etcd database.
type Database struct {
	Users         []*User
	Groups        []Group
	DeletedUsers  []string
	DeletedGroups []string
	LockedUsers   []string
}

// GetDatabase takes a snapshot of etcd database at revision rev.
// If rev is 0, the snapshot will be the latest one.
func GetDatabase(ctx context.Context, etcd *clientv3.Client, rev int64) (*Database, error) {
	db := new(Database)

	resp, err := etcd.Get(ctx, KeyUsers+"/",
		clientv3.WithPrefix(),
		clientv3.WithRev(rev),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		users := make([]*User, resp.Count)
		for i, kv := range resp.Kvs {
			u := new(User)
			err := json.Unmarshal(kv.Value, u)
			if err != nil {
				return nil, err
			}
			users[i] = u
		}
		db.Users = users
	}

	gprefix := KeyGroups + "/"
	resp, err = etcd.Get(ctx, gprefix,
		clientv3.WithPrefix(),
		clientv3.WithRev(rev),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		groups := make([]Group, resp.Count)
		for i, kv := range resp.Kvs {
			name := string(kv.Key[len(gprefix):])
			gid, err := strconv.Atoi(string(kv.Value))
			if err != nil {
				return nil, err
			}
			groups[i] = Group{name, gid}
		}
		db.Groups = groups
	}

	getKeys := func(prefix string) ([]string, error) {
		resp, err := etcd.Get(ctx, prefix,
			clientv3.WithPrefix(),
			clientv3.WithRev(rev),
			clientv3.WithKeysOnly(),
		)
		if err != nil {
			return nil, err
		}
		if resp.Count == 0 {
			return nil, nil
		}
		ret := make([]string, resp.Count)
		for i, kv := range resp.Kvs {
			ret[i] = string(kv.Key[len(prefix):])
		}
		return ret, nil
	}

	delus, err := getKeys(KeyDeletedUsers + "/")
	if err != nil {
		return nil, err
	}
	db.DeletedUsers = delus

	delgs, err := getKeys(KeyDeletedGroups + "/")
	if err != nil {
		return nil, err
	}
	db.DeletedGroups = delgs

	locked, err := getKeys(KeyLocked + "/")
	if err != nil {
		return nil, err
	}
	db.LockedUsers = locked

	return db, nil
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	ca := make([]string, len(a))
	copy(ca, a)
	sort.Strings(ca)
	cb := make([]string, len(b))
	copy(cb, b)
	sort.Strings(cb)

	for i, s := range ca {
		if cb[i] != s {
			return false
		}
	}
	return true
}

// Sync synchronizes users and groups using sc.
func (db *Database) Sync(ctx context.Context, sc Syncer) error {
	// lock passwords
	for _, name := range db.LockedUsers {
		lu, err := sc.LookupUser(ctx, name)
		if err != nil {
			return err
		}
		if lu == nil {
			continue
		}
		err = sc.LockPassword(ctx, name)
		if err != nil {
			return err
		}
		log.Info("locked password", map[string]interface{}{
			"user": name,
		})
	}

	// remove groups
	for _, name := range db.DeletedGroups {
		lg, err := sc.LookupGroup(ctx, name)
		if err != nil {
			return err
		}
		if lg == nil {
			continue
		}
		err = sc.RemoveGroup(ctx, name)
		if err != nil {
			return err
		}
		log.Info("removed a group", map[string]interface{}{
			"group": name,
			"gid":   lg.GID,
		})
	}

	// add groups
	for _, g := range db.Groups {
		lg, err := sc.LookupGroup(ctx, g.Name)
		if err != nil {
			return err
		}

		if lg != nil {
			if g.GID == lg.GID {
				return nil
			}

			// the system has a group with the same name but different GID.
			// remove it first.
			err = sc.RemoveGroup(ctx, g.Name)
			if err != nil {
				return err
			}
			log.Info("removed a group", map[string]interface{}{
				"group": lg.Name,
				"gid":   lg.GID,
			})
		}

		err = sc.AddGroup(ctx, g)
		if err != nil {
			return err
		}
		log.Info("added a group", map[string]interface{}{
			"group": g.Name,
			"gid":   g.GID,
		})
	}

	// remove users
	for _, name := range db.DeletedUsers {
		lu, err := sc.LookupUser(ctx, name)
		if err != nil {
			return err
		}
		if lu == nil {
			continue
		}
		err = sc.RemoveUser(ctx, name)
		if err != nil {
			return err
		}
		log.Info("removed a user", map[string]interface{}{
			"user": name,
			"uid":  lu.UID,
		})
	}

	// add or update users
	for _, u := range db.Users {
		lu, err := sc.LookupUser(ctx, u.Name)
		if err != nil {
			return err
		}

		if lu != nil {
			if u.UID == lu.UID {
				goto UPDATE
			}

			// the system has a user with the same name but different UID.
			// remove it first.
			err = sc.RemoveUser(ctx, u.Name)
			if err != nil {
				return err
			}
			log.Info("removed a user", map[string]interface{}{
				"user": lu.Name,
				"uid":  lu.UID,
			})
		}

		err = sc.AddUser(ctx, u)
		if err != nil {
			return err
		}
		log.Info("added a user", map[string]interface{}{
			"user": u.Name,
			"uid":  u.UID,
		})
		continue

	UPDATE:
		if lu.DisplayName != u.DisplayName {
			err = sc.SetDisplayName(ctx, u.Name, u.DisplayName)
			if err != nil {
				return err
			}
			log.Info("updated display name", map[string]interface{}{
				"user": u.Name,
			})
		}

		if lu.Group != u.Group {
			err = sc.SetPrimaryGroup(ctx, u.Name, u.Group)
			if err != nil {
				return err
			}
			log.Info("updated primary group", map[string]interface{}{
				"user": u.Name,
			})
		}

		if !equalStringSlice(lu.Groups, u.Groups) {
			err = sc.SetSupplementalGroups(ctx, u.Name, u.Groups)
			if err != nil {
				return err
			}
			log.Info("updated supplementary groups", map[string]interface{}{
				"user": u.Name,
			})
		}

		if lu.Shell != u.Shell {
			err = sc.SetShell(ctx, u.Name, u.Shell)
			if err != nil {
				return err
			}
			log.Info("updated shell", map[string]interface{}{
				"user": u.Name,
			})
		}

		if !equalStringSlice(lu.PubKeys, u.PubKeys) {
			err = sc.SetPubKeys(ctx, u.Name, u.PubKeys)
			if err != nil {
				return err
			}
			log.Info("updated public keys", map[string]interface{}{
				"user": u.Name,
			})
		}
	}

	return nil
}
