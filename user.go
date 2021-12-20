package etcdpasswd

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/clientv3util"
)

// User represents a user managed by etcdpasswd.
type User struct {
	Name        string   `json:"name"`
	UID         int      `json:"uid"`
	DisplayName string   `json:"display-name"`
	Group       string   `json:"group"`
	Groups      []string `json:"groups"`
	Shell       string   `json:"shell"`
	PubKeys     []string `json:"public-keys"`
}

// Validate validates User attributes.
// If some attribute is not valid, a non-nil error is returned.
func (u *User) Validate() error {
	if !IsValidUserName(u.Name) {
		return errors.New("invalid user name: " + u.Name)
	}
	if !IsValidGroupName(u.Group) {
		return errors.New("invalid group name: " + u.Group)
	}
	if len(u.Shell) == 0 {
		return errors.New("shell must be specified")
	}
	return nil
}

// ListUsers lists all user names registered in the database.
// The result is sorted alphabetically.
func (c Client) ListUsers(ctx context.Context) ([]string, error) {
	return c.list(ctx, KeyUsers)
}

// GetUser looks up named user from the database.
// If the user is not found, this returns ErrNotFound.
func (c Client) GetUser(ctx context.Context, name string) (*User, int64, error) {
	key := KeyUsers + name

	resp, err := c.Get(ctx, key)
	if err != nil {
		return nil, 0, err
	}

	if resp.Count == 0 {
		return nil, 0, ErrNotFound
	}

	u := new(User)
	err = json.Unmarshal(resp.Kvs[0].Value, u)
	if err != nil {
		return nil, 0, err
	}

	return u, resp.Kvs[0].ModRevision, nil
}

// AddUser adds a new managed user to the database.
// If a user having the same name already exists, ErrExists will be returned.
func (c Client) AddUser(ctx context.Context, user *User) error {
	cfg, _, err := c.GetConfig(ctx)
	if err != nil {
		return err
	}
	if cfg.StartUID == 0 {
		// not yet configured
		return errors.New("start-uid has not been set")
	}

	if len(user.Group) == 0 {
		user.Group = cfg.DefaultGroup
	}
	if len(user.Groups) == 0 {
		user.Groups = cfg.DefaultGroups
	}
	if len(user.Shell) == 0 {
		user.Shell = cfg.DefaultShell
	}
	err = user.Validate()
	if err != nil {
		return err
	}

	key := KeyUsers + user.Name
	delKey := KeyDeletedUsers + user.Name

RETRY:
	uid, rev, err := c.getLastID(ctx, KeyLastUID, cfg.StartUID)
	if err != nil {
		return err
	}
	user.UID = uid
	j, err := json.Marshal(user)
	if err != nil {
		return err
	}
	uidData := strconv.Itoa(uid + 1)

	resp, err := c.Txn(ctx).
		If(
			clientv3.Compare(clientv3.ModRevision(KeyLastUID), "=", rev),
		).
		Then(
			clientv3.OpTxn(
				[]clientv3.Cmp{clientv3util.KeyMissing(key)},
				[]clientv3.Op{
					clientv3.OpPut(KeyLastUID, uidData),
					clientv3.OpPut(key, string(j)),
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

// UpdateUser updates an existing managed user in the database.
// This operation does compare-and-swap with rev.  If CAS failed,
// ErrCASFailure will be returned.
func (c Client) UpdateUser(ctx context.Context, user *User, rev int64) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	key := KeyUsers + user.Name
	j, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := c.Txn(ctx).
		If(
			clientv3.Compare(clientv3.ModRevision(key), "=", rev),
		).
		Then(clientv3.OpPut(key, string(j))).
		Commit()
	if err != nil {
		return err
	}

	if !resp.Succeeded {
		return ErrCASFailure
	}
	return nil
}

// RemoveUser removes an existing managed user.
// If the user does not exist, ErrNotFound will be returned.
func (c Client) RemoveUser(ctx context.Context, name string) error {
	key := KeyUsers + name
	delKey := KeyDeletedUsers + name

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
