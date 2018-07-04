package syncer

import (
	"context"
	"errors"
	"os/user"
	"strconv"
	"strings"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/etcdpasswd"
)

// UbuntuSyncer synchronizes with local users/groups of Debian/Ubuntu OS.
type UbuntuSyncer struct{}

// LookupUser implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) LookupUser(ctx context.Context, name string) (*etcdpasswd.User, error) {
	uu, err := user.Lookup(name)
	if err != nil {
		if _, ok := err.(user.UnknownUserError); ok {
			return nil, nil
		}
		return nil, err
	}

	return makeUser(uu)
}

// LookupGroup implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) LookupGroup(ctx context.Context, name string) (*etcdpasswd.Group, error) {
	gg, err := user.LookupGroup(name)
	if err != nil {
		if _, ok := err.(user.UnknownGroupError); ok {
			return nil, nil
		}
		return nil, err
	}
	gid, err := strconv.Atoi(gg.Gid)
	if err != nil {
		return nil, err
	}

	return &etcdpasswd.Group{Name: gg.Name, GID: gid}, nil
}

// AddUser implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) AddUser(ctx context.Context, u *etcdpasswd.User) error {
	_, err := user.Lookup(u.Name)
	if err == nil {
		return errors.New("user exists: " + u.Name)
	}

	args := []string{
		"-c", u.DisplayName, "-e", "", "-f", "-1", "-g", u.Group,
		"-m", "-s", u.Shell, "-u", strconv.Itoa(u.UID),
	}
	if len(u.Groups) > 0 {
		args = append(args, "-G")
		args = append(args, strings.Join(u.Groups, ","))
	}
	args = append(args, u.Name)

	// use background context to ignore cancellation.
	return cmd.CommandContext(context.Background(), "useradd", args...).Run()
}

// RemoveUser implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) RemoveUser(ctx context.Context, name string) error {
	_, err := user.Lookup(name)
	if err != nil {
		return err
	}

	// use background context to ignore cancellation.
	return cmd.CommandContext(context.Background(), "userdel", "-f", "-r", name).Run()
}

func userMod(args ...string) error {
	return cmd.CommandContext(context.Background(), "usermod", args...).Run()
}

// SetDisplayName implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetDisplayName(ctx context.Context, name, displayName string) error {
	return userMod("-c", displayName, name)
}

// SetPrimaryGroup implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetPrimaryGroup(ctx context.Context, name, group string) error {
	return userMod("-g", group, name)
}

// SetSupplementalGroups implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetSupplementalGroups(ctx context.Context, name string, groups []string) error {
	return userMod("-G", strings.Join(groups, ","), name)
}

// SetShell implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetShell(ctx context.Context, name, shell string) error {
	return userMod("-s", shell, name)
}

// SetPubKeys implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) SetPubKeys(ctx context.Context, name string, pubkeys []string) error {
	uu, err := user.Lookup(name)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(uu.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(uu.Gid)
	if err != nil {
		return err
	}

	return savePubKeys(uu.HomeDir, uid, gid, pubkeys)
}

// LockPassword implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) LockPassword(ctx context.Context, name string) error {
	return userMod("-L", name)
}

// AddGroup implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) AddGroup(ctx context.Context, g etcdpasswd.Group) error {
	_, err := user.LookupGroup(g.Name)
	if err == nil {
		return errors.New("group exists: " + g.Name)
	}

	// use background context to ignore cancellation.
	return cmd.CommandContext(context.Background(),
		"groupadd", "-g", strconv.Itoa(g.GID), g.Name).Run()
}

// RemoveGroup implements etcdpasswd.Syncer interface.
func (s UbuntuSyncer) RemoveGroup(ctx context.Context, name string) error {
	_, err := user.LookupGroup(name)
	if err != nil {
		return err
	}

	// use background context to ignore cancellation.
	return cmd.CommandContext(context.Background(), "groupdel", name).Run()
}
