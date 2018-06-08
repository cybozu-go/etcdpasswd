package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/google/subcommands"
)

type userCmd struct{}

func (c userCmd) SetFlags(f *flag.FlagSet) {}

func (c userCmd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
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
		userCmd{},
		"user",
		"manage user database",
		"user ACTION ...",
	}
}

type userListCmd struct{}

func (c userListCmd) SetFlags(f *flag.FlagSet) {}

func (c userListCmd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
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
		userListCmd{},
		"list",
		"list users",
		"list",
	}
}

type userGetCmd struct{}

func (c userGetCmd) SetFlags(f *flag.FlagSet) {}

func (c userGetCmd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
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
public-keys: %v
`, user.UID, user.DisplayName, user.Group, user.Groups, user.Shell, user.PubKeys)
	return handleError(nil)
}

func userGetCommand() subcommands.Command {
	return subcmd{
		userGetCmd{},
		"get",
		"get user information",
		"get NAME",
	}
}

type userAddCmd struct {
	displayName string
	group       string
	groups      commaStrings
	shell       string
}

func (c *userAddCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.displayName, "display", "", "display name")
	f.StringVar(&c.group, "group", "", "primary group")
	f.Var(&c.groups, "groups", "comma-separated supplementary groups")
	f.StringVar(&c.shell, "shell", etcdpasswd.DefaultShell, "shell program")
}

func (c *userAddCmd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
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
		&userAddCmd{},
		"add",
		"add a new user",
		"add [OPTIONS] NAME",
	}
}

type userUpdateCmd struct {
	displayName string
	group       string
	groups      commaStrings
	shell       string
}

func (c *userUpdateCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.displayName, "display", "", "display name")
	f.StringVar(&c.group, "group", "", "primary group")
	f.Var(&c.groups, "groups", "comma-separated supplementary groups")
	f.StringVar(&c.shell, "shell", "", "shell program")
}

func (c *userUpdateCmd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
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
		&userUpdateCmd{},
		"update",
		"update an existing user",
		"update [OPTIONS] NAME",
	}
}

type userRemoveCmd struct{}

func (c userRemoveCmd) SetFlags(f *flag.FlagSet) {}

func (c userRemoveCmd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
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
		userRemoveCmd{},
		"remove",
		"remove an existing user",
		"remove NAME",
	}
}
