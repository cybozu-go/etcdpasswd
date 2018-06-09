package syncer

import (
	"context"

	"github.com/cybozu-go/etcdpasswd"
)

// DummySyncer emulates synchronization with OS-local users/groups.
type DummySyncer struct{}

// LookupUser implements etcdpasswd.Syncer interface.
func (s DummySyncer) LookupUser(ctx context.Context, name string) (*etcdpasswd.User, error) {
	return nil, nil
}

// LookupGroup implements etcdpasswd.Syncer interface.
func (s DummySyncer) LookupGroup(ctx context.Context, name string) (*etcdpasswd.Group, error) {
	return nil, nil
}

// AddUser implements etcdpasswd.Syncer interface.
func (s DummySyncer) AddUser(ctx context.Context, user *etcdpasswd.User) error {
	return nil
}

// RemoveUser implements etcdpasswd.Syncer interface.
func (s DummySyncer) RemoveUser(ctx context.Context, name string) error {
	return nil
}

// SetDisplayName implements etcdpasswd.Syncer interface.
func (s DummySyncer) SetDisplayName(ctx context.Context, name, displayName string) error {
	return nil
}

// SetPrimaryGroup implements etcdpasswd.Syncer interface.
func (s DummySyncer) SetPrimaryGroup(ctx context.Context, name, group string) error {
	return nil
}

// SetSupplementalGroups implements etcdpasswd.Syncer interface.
func (s DummySyncer) SetSupplementalGroups(ctx context.Context, name string, groups []string) error {
	return nil
}

// SetShell implements etcdpasswd.Syncer interface.
func (s DummySyncer) SetShell(ctx context.Context, name, shell string) error {
	return nil
}

// SetPubKeys implements etcdpasswd.Syncer interface.
func (s DummySyncer) SetPubKeys(ctx context.Context, name string, pubkeys []string) error {
	return nil
}

// LockPassword implements etcdpasswd.Syncer interface.
func (s DummySyncer) LockPassword(ctx context.Context, name string) error {
	return nil
}

// AddGroup implements etcdpasswd.Syncer interface.
func (s DummySyncer) AddGroup(ctx context.Context, group etcdpasswd.Group) error {
	return nil
}

// RemoveGroup implements etcdpasswd.Syncer interface.
func (s DummySyncer) RemoveGroup(ctx context.Context, name string) error {
	return nil
}
