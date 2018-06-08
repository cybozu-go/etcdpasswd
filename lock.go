package etcdpasswd

import (
	"context"
	"errors"
	"path"
)

// Lock adds name to locked user database on etcd.
func (c Client) Lock(ctx context.Context, name string) error {
	if !IsValidUserName(name) {
		return errors.New("invalid user name: " + name)
	}
	key := path.Join(KeyLocked, name)

	_, err := c.Put(ctx, key, "")
	return err
}

// Unlock removes name from locked user database on etcd.
func (c Client) Unlock(ctx context.Context, name string) error {
	key := path.Join(KeyLocked, name)

	_, err := c.Delete(ctx, key)
	return err
}
