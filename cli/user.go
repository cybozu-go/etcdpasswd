package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/google/subcommands"
)

type user struct{}

func (c user) SetFlags(f *flag.FlagSet) {}

func (c user) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	newc := NewCommander(f, "user")
	newc.Register(userListCommand(), "")
	newc.Register(userGetCommand(), "")
	newc.Register(userAddCommand(), "")
	newc.Register(userUpdateCommand(), "")
	newc.Register(userRemoveCommand(), "")
	return newc.Execute(ctx)
}

// UserCommand implements "user" subcommand.
func UserCommand() subcommands.Command {
	return subcmd{
		user{},
		"user",
		"manage user database",
		"user ACTION ...",
	}
}

type userList struct{}

func (c userList) SetFlags(f *flag.FlagSet) {}

func (c userList) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	users, err := client.ListUser(ctx)
	if err != nil {
		return handleError(err)
	}

	for _, u := range users {
		fmt.Println(u)
	}
	return handleError(nil)
}

func userListCommand() subcommands.Command {
	return subcmd{
		userList{},
		"list",
		"list users",
		"list",
	}
}

type userGet struct{}

func (c userGet) SetFlags(f *flag.FlagSet) {}

func (c userGet) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}

	name := f.Arg(0)
	user, _, err := client.GetUser(ctx, name)
	if err != nil {
		return handleError(err)
	}

	fmt.Printf(`uid: %d
display-name: %s
group: %s
groups: %v
shell: %s
public-keys: %d
`, user.UID, user.DisplayName, user.Group, user.Groups, user.Shell, len(user.PubKeys))
	return handleError(nil)
}

func userGetCommand() subcommands.Command {
	return subcmd{
		userGet{},
		"get",
		"get user information",
		"get NAME",
	}
}

type userAdd struct {
	displayName string
	group       string
	groups      commaStrings
	shell       string
}

func (c *userAdd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.displayName, "display", "", "display name")
	f.StringVar(&c.group, "group", "", "primary group")
	f.Var(&c.groups, "groups", "comma-separated supplementary groups")
	f.StringVar(&c.shell, "shell", etcdpasswd.DefaultShell, "shell program")
}

func (c *userAdd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}

	name := f.Arg(0)
	user := &etcdpasswd.User{
		Name:        name,
		DisplayName: c.displayName,
		Group:       c.group,
		Groups:      []string(c.groups),
		Shell:       c.shell,
	}

	err := client.AddUser(ctx, user)
	return handleError(err)
}

func userAddCommand() subcommands.Command {
	return subcmd{
		&userAdd{},
		"add",
		"add a new user",
		"add [OPTIONS] NAME",
	}
}

type userUpdate struct {
	displayName string
	group       string
	groups      commaStrings
	shell       string
}

func (c *userUpdate) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.displayName, "display", "", "display name")
	f.StringVar(&c.group, "group", "", "primary group")
	f.Var(&c.groups, "groups", "comma-separated supplementary groups")
	f.StringVar(&c.shell, "shell", "", "shell program")
}

func (c *userUpdate) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}

	name := f.Arg(0)
	user, rev, err := client.GetUser(ctx, name)
	if err != nil {
		return handleError(err)
	}

	if len(c.displayName) > 0 {
		user.DisplayName = c.displayName
	}
	if len(c.group) > 0 {
		user.Group = c.group
	}
	if len(c.groups) > 0 {
		user.Groups = []string(c.groups)
	}
	if len(c.shell) > 0 {
		user.Shell = c.shell
	}

	err = client.UpdateUser(ctx, user, rev)
	return handleError(err)
}

func userUpdateCommand() subcommands.Command {
	return subcmd{
		&userUpdate{},
		"update",
		"update an existing user",
		"update [OPTIONS] NAME",
	}
}

type userRemove struct{}

func (c userRemove) SetFlags(f *flag.FlagSet) {}

func (c userRemove) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}

	name := f.Arg(0)
	err := client.RemoveUser(ctx, name)
	return handleError(err)
}

func userRemoveCommand() subcommands.Command {
	return subcmd{
		userRemove{},
		"remove",
		"remove an existing user",
		"remove NAME",
	}
}
