package syncer

import (
	"context"
	"errors"

	"github.com/cybozu-go/etcdpasswd"
)

// MockSyncer is a memory-only synchronizer for tests.
type MockSyncer struct {
	Users       map[string]*etcdpasswd.User
	Groups      map[string]*etcdpasswd.Group
	LockedUsers map[string]bool
}

// NewMockSyncer creates a MockSyncer instance.
func NewMockSyncer() *MockSyncer {
	return &MockSyncer{
		make(map[string]*etcdpasswd.User),
		make(map[string]*etcdpasswd.Group),
		make(map[string]bool),
	}
}

func (s *MockSyncer) checkGroup(u *etcdpasswd.User) error {
	if _, ok := s.Groups[u.Group]; !ok {
		return errors.New("no such group: " + u.Group)
	}
	for _, g := range u.Groups {
		if _, ok := s.Groups[g]; !ok {
			return errors.New("no such group: " + g)
		}
	}
	return nil
}

// LookupUser implements etcdpasswd.Syncer interface.
func (s *MockSyncer) LookupUser(ctx context.Context, name string) (*etcdpasswd.User, error) {
	return s.Users[name], nil
}

// LookupGroup implements etcdpasswd.Syncer interface.
func (s *MockSyncer) LookupGroup(ctx context.Context, name string) (*etcdpasswd.Group, error) {
	return s.Groups[name], nil
}

// AddUser implements etcdpasswd.Syncer interface.
func (s *MockSyncer) AddUser(ctx context.Context, user *etcdpasswd.User) error {
	if _, ok := s.Users[user.Name]; ok {
		return errors.New("exists")
	}
	err := s.checkGroup(user)
	if err != nil {
		return err
	}

	s.Users[user.Name] = user
	return nil
}

// RemoveUser implements etcdpasswd.Syncer interface.
func (s *MockSyncer) RemoveUser(ctx context.Context, name string) error {
	if _, ok := s.Users[name]; ok {
		delete(s.Users, name)
		return nil
	}
	return errors.New("not exists")
}

// SetDisplayName implements etcdpasswd.Syncer interface.
func (s *MockSyncer) SetDisplayName(ctx context.Context, name, displayName string) error {
	if u, ok := s.Users[name]; ok {
		u.DisplayName = displayName
		return nil
	}
	return errors.New("not exists")
}

// SetPrimaryGroup implements etcdpasswd.Syncer interface.
func (s *MockSyncer) SetPrimaryGroup(ctx context.Context, name, group string) error {
	if _, ok := s.Groups[group]; !ok {
		return errors.New("no such group: " + group)
	}

	if u, ok := s.Users[name]; ok {
		u.Group = group
		return nil
	}
	return errors.New("not exists")
}

// SetSupplementalGroups implements etcdpasswd.Syncer interface.
func (s *MockSyncer) SetSupplementalGroups(ctx context.Context, name string, groups []string) error {
	for _, g := range groups {
		if _, ok := s.Groups[g]; !ok {
			return errors.New("no such group: " + g)
		}
	}

	if u, ok := s.Users[name]; ok {
		sg := make([]string, len(groups))
		copy(sg, groups)
		u.Groups = sg
		return nil
	}
	return errors.New("not exists")
}

// SetShell implements etcdpasswd.Syncer interface.
func (s *MockSyncer) SetShell(ctx context.Context, name, shell string) error {
	if u, ok := s.Users[name]; ok {
		u.Shell = shell
		return nil
	}
	return errors.New("not exists")
}

// SetPubKeys implements etcdpasswd.Syncer interface.
func (s *MockSyncer) SetPubKeys(ctx context.Context, name string, pubkeys []string) error {
	if u, ok := s.Users[name]; ok {
		pks := make([]string, len(pubkeys))
		copy(pks, pubkeys)
		u.PubKeys = pks
		return nil
	}
	return errors.New("not exists")
}

// LockPassword implements etcdpasswd.Syncer interface.
func (s *MockSyncer) LockPassword(ctx context.Context, name string) error {
	s.LockedUsers[name] = true
	return nil
}

// AddGroup implements etcdpasswd.Syncer interface.
func (s *MockSyncer) AddGroup(ctx context.Context, group etcdpasswd.Group) error {
	if _, ok := s.Groups[group.Name]; ok {
		return errors.New("exists")
	}
	s.Groups[group.Name] = &group
	return nil
}

// RemoveGroup implements etcdpasswd.Syncer interface.
func (s *MockSyncer) RemoveGroup(ctx context.Context, name string) error {
	if _, ok := s.Groups[name]; ok {
		delete(s.Groups, name)
		return nil
	}
	return errors.New("not exists")
}
