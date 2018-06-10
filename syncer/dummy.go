package syncer

import (
	"context"
	"errors"
	"os/user"
	"strconv"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/log"
)

// DummySyncer emulates synchronization with OS-local users/groups.
type DummySyncer struct {
	users         map[string]*etcdpasswd.User
	groups        map[string]*etcdpasswd.Group
	deletedUsers  map[string]bool
	deletedGroups map[string]bool
}

// NewDummySyncer creates a DummySyncer.
func NewDummySyncer() *DummySyncer {
	return &DummySyncer{
		users:         make(map[string]*etcdpasswd.User),
		groups:        make(map[string]*etcdpasswd.Group),
		deletedUsers:  make(map[string]bool),
		deletedGroups: make(map[string]bool),
	}
}

// LookupUser implements etcdpasswd.Syncer interface.
func (s *DummySyncer) LookupUser(ctx context.Context, name string) (*etcdpasswd.User, error) {
	if s.deletedUsers[name] {
		// the user was logically deleted
		return nil, nil
	}

	if u, ok := s.users[name]; ok {
		return u, nil
	}

	uu, err := user.Lookup(name)
	if err != nil {
		if _, ok := err.(user.UnknownUserError); ok {
			s.deletedUsers[name] = true
			return nil, nil
		}
		return nil, err
	}

	u, err := makeUser(uu)
	if err != nil {
		return nil, err
	}
	s.users[name] = u
	return u, nil
}

// LookupGroup implements etcdpasswd.Syncer interface.
func (s *DummySyncer) LookupGroup(ctx context.Context, name string) (*etcdpasswd.Group, error) {
	if s.deletedGroups[name] {
		// the group was logically deleted
		return nil, nil
	}

	if g, ok := s.groups[name]; ok {
		return g, nil
	}

	gg, err := user.LookupGroup(name)
	if err != nil {
		if _, ok := err.(user.UnknownGroupError); ok {
			s.deletedGroups[name] = true
			return nil, nil
		}
		return nil, err
	}
	gid, err := strconv.Atoi(gg.Gid)
	if err != nil {
		return nil, err
	}

	g := &etcdpasswd.Group{Name: gg.Name, GID: gid}
	s.groups[name] = g
	return g, nil
}

// AddUser implements etcdpasswd.Syncer interface.
func (s *DummySyncer) AddUser(ctx context.Context, u *etcdpasswd.User) error {
	if _, ok := s.users[u.Name]; ok {
		return errors.New("user already exists: " + u.Name)
	}
	s.users[u.Name] = u
	delete(s.deletedUsers, u.Name)
	return nil
}

// RemoveUser implements etcdpasswd.Syncer interface.
func (s *DummySyncer) RemoveUser(ctx context.Context, name string) error {
	if _, ok := s.users[name]; !ok {
		return errors.New("user does not exist: " + name)
	}
	delete(s.users, name)
	s.deletedUsers[name] = true
	return nil
}

// SetDisplayName implements etcdpasswd.Syncer interface.
func (s *DummySyncer) SetDisplayName(ctx context.Context, name, displayName string) error {
	if u, ok := s.users[name]; ok {
		u.DisplayName = displayName
		return nil
	}
	return errors.New("user does not exist: " + name)
}

// SetPrimaryGroup implements etcdpasswd.Syncer interface.
func (s *DummySyncer) SetPrimaryGroup(ctx context.Context, name, group string) error {
	if _, ok := s.groups[group]; !ok {
		return errors.New("group does not exist: " + group)
	}

	if u, ok := s.users[name]; ok {
		u.Group = group
		return nil
	}
	return errors.New("user does not exist: " + name)
}

// SetSupplementalGroups implements etcdpasswd.Syncer interface.
func (s *DummySyncer) SetSupplementalGroups(ctx context.Context, name string, groups []string) error {
	for _, gname := range groups {
		if _, ok := s.groups[gname]; !ok {
			return errors.New("group does not exist: " + gname)
		}
	}

	if u, ok := s.users[name]; ok {
		u.Groups = groups
		return nil
	}
	return errors.New("user does not exist: " + name)
}

// SetShell implements etcdpasswd.Syncer interface.
func (s *DummySyncer) SetShell(ctx context.Context, name, shell string) error {
	if u, ok := s.users[name]; ok {
		u.Shell = shell
		return nil
	}
	return errors.New("user does not exist: " + name)
}

// SetPubKeys implements etcdpasswd.Syncer interface.
func (s *DummySyncer) SetPubKeys(ctx context.Context, name string, pubkeys []string) error {
	if u, ok := s.users[name]; ok {
		u.PubKeys = pubkeys
		return nil
	}
	return errors.New("user does not exist: " + name)
}

// LockPassword implements etcdpasswd.Syncer interface.
func (s *DummySyncer) LockPassword(ctx context.Context, name string) error {
	if _, ok := s.users[name]; ok {
		log.Warn("dummy: lock password", map[string]interface{}{
			"user": name,
		})
	}
	return nil
}

// AddGroup implements etcdpasswd.Syncer interface.
func (s *DummySyncer) AddGroup(ctx context.Context, g etcdpasswd.Group) error {
	if _, ok := s.groups[g.Name]; ok {
		return errors.New("group already exists: " + g.Name)
	}
	s.groups[g.Name] = &g
	delete(s.deletedGroups, g.Name)
	return nil
}

// RemoveGroup implements etcdpasswd.Syncer interface.
func (s *DummySyncer) RemoveGroup(ctx context.Context, name string) error {
	if _, ok := s.groups[name]; !ok {
		return errors.New("group does not exist: " + name)
	}
	delete(s.groups, name)
	s.deletedGroups[name] = true
	return nil
}
