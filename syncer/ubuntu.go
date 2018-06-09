package syncer

import (
	"context"

	"github.com/cybozu-go/etcdpasswd"
)

// UbuntuSyncer synchronizes with local users/groups of Debian/Ubuntu OS.
type UbuntuSyncer struct{}

// LookupUser implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) LookupUser(ctx context.Context, name string) (*etcdpasswd.User, error) {
	return nil, nil
}

// LookupGroup implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) LookupGroup(ctx context.Context, name string) (*etcdpasswd.Group, error) {
	return nil, nil
}

// AddUser implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) AddUser(ctx context.Context, user *etcdpasswd.User) error {
	return nil
}

// RemoveUser implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) RemoveUser(ctx context.Context, name string) error {
	return nil
}

// SetDisplayName implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetDisplayName(ctx context.Context, name, displayName string) error {
	return nil
}

// SetPrimaryGroup implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetPrimaryGroup(ctx context.Context, name, group string) error {
	return nil
}

// SetSupplementalGroups implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetSupplementalGroups(ctx context.Context, name string, groups []string) error {
	return nil
}

// SetShell implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetShell(ctx context.Context, name, shell string) error {
	return nil
}

// SetPubKeys implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetPubKeys(ctx context.Context, name string, pubkeys []string) error {
	return nil
}

// LockPassword implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) LockPassword(ctx context.Context, name string) error {
	return nil
}

// AddGroup implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) AddGroup(ctx context.Context, group etcdpasswd.Group) error {
	return nil
}

// RemoveGroup implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) RemoveGroup(ctx context.Context, name string) error {
	return nil
}
