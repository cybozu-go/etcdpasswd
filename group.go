package etcdpasswd

import (
	"context"
	"errors"
	"strconv"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/clientv3util"
)

// Group represents attributes of a group.
type Group struct {
	Name string
	GID  int
}

// ListGroups lists all groups registered in the database.
// The result is sorted alphabetically.
func (c Client) ListGroups(ctx context.Context) ([]Group, error) {
	prefix := KeyGroups
	resp, err := c.Get(ctx, prefix, clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return nil, err
	}

	ret := make([]Group, resp.Count)
	for i, kv := range resp.Kvs {
		gid, err := strconv.Atoi(string(kv.Value))
		if err != nil {
			return nil, err
		}
		ret[i] = Group{string(kv.Key[len(prefix):]), gid}
	}
	return ret, nil
}

// AddGroup adds a new managed group to the database.
// If a group having the same name already exists, ErrExists will be returned.
func (c Client) AddGroup(ctx context.Context, name string) error {
	cfg, _, err := c.GetConfig(ctx)
	if err != nil {
		return err
	}
	if cfg.StartGID == 0 {
		// not yet configured
		return errors.New("start-gid has not been set")
	}

	if !IsValidGroupName(name) {
		return errors.New("invalid group name: " + name)
	}

	key := KeyGroups + name
	delKey := KeyDeletedGroups + name

RETRY:
	gid, rev, err := c.getLastID(ctx, KeyLastGID, cfg.StartGID)
	if err != nil {
		return err
	}
	gidStr := strconv.Itoa(gid)
	nextGID := strconv.Itoa(gid + 1)

	resp, err := c.Txn(ctx).
		If(
			clientv3.Compare(clientv3.ModRevision(KeyLastGID), "=", rev),
		).
		Then(
			clientv3.OpTxn(
				[]clientv3.Cmp{clientv3util.KeyMissing(key)},
				[]clientv3.Op{
					clientv3.OpPut(KeyLastGID, nextGID),
					clientv3.OpPut(key, gidStr),
					clientv3.OpDelete(delKey),
				},
				nil,
			),
		).
		Commit()
	if err != nil {
		return err
	}

	if !resp.Succeeded {
		goto RETRY
	}
	if !resp.Responses[0].GetResponseTxn().Succeeded {
		return ErrExists
	}
	return nil
}

// RemoveGroup removes an existing managed group.
// If the group does not exist, ErrNotFound will be returned.
func (c Client) RemoveGroup(ctx context.Context, name string) error {
	key := KeyGroups + name
	delKey := KeyDeletedGroups + name

	resp, err := c.Txn(ctx).
		If(clientv3util.KeyExists(key)).
		Then(
			clientv3.OpDelete(key),
			clientv3.OpPut(delKey, ""),
		).
		Commit()
	if err != nil {
		return err
	}

	if !resp.Succeeded {
		return ErrNotFound
	}
	return nil
}
