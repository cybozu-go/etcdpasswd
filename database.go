package etcdpasswd

import (
	"context"
	"encoding/json"
	"strconv"

	clientv3 "go.etcd.io/etcd/client/v3"
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

	resp, err := etcd.Get(ctx, KeyUsers,
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

	resp, err = etcd.Get(ctx, KeyGroups,
		clientv3.WithPrefix(),
		clientv3.WithRev(rev),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		groups := make([]Group, resp.Count)
		for i, kv := range resp.Kvs {
			name := string(kv.Key[len(KeyGroups):])
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

	delus, err := getKeys(KeyDeletedUsers)
	if err != nil {
		return nil, err
	}
	db.DeletedUsers = delus

	delgs, err := getKeys(KeyDeletedGroups)
	if err != nil {
		return nil, err
	}
	db.DeletedGroups = delgs

	locked, err := getKeys(KeyLocked)
	if err != nil {
		return nil, err
	}
	db.LockedUsers = locked

	return db, nil
}
