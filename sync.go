package etcdpasswd

import "context"

// Syncer is an interface for user and group synchronization.
type Syncer interface {
	// LookupUser looks up the named user in the system.
	// If the user is not found, this should return (nil, nil).
	LookupUser(ctx context.Context, name string) (*User, error)

	// LookupGroup looks up the named group in the system.
	// If the group is not found, this should return (nil, nil).
	LookupGroup(ctx context.Context, name string) (*Group, error)

	// AddUser adds a user to the system.
	AddUser(ctx context.Context, user *User) error

	// RemoveUser removes a user from the system.
	RemoveUser(ctx context.Context, name string) error

	// SetDisplayName sets the display name of the user.
	SetDisplayName(ctx context.Context, name, displayName string) error

	// SetPrimaryGroup sets the primary group of the user.
	SetPrimaryGroup(ctx context.Context, name, group string) error

	// SetSupplementalGroups sets the supplemental groups of the user.
	SetSupplementalGroups(ctx context.Context, name string, groups []string) error

	// SetShell sets the login shell of the user.
	SetShell(ctx context.Context, name, shell string) error

	// SetPubKeys sets SSH authorized keys of the user.
	SetPubKeys(ctx context.Context, name string, pubkeys []string) error

	// LockPassword locks the password of the user to prohibit login attempts using password.
	LockPassword(ctx context.Context, name string) error

	// AddGroup adds a group to the system.
	AddGroup(ctx context.Context, group Group) error

	// RemoveGroup removes a group from the system.
	RemoveGroup(ctx context.Context, name string) error
}
