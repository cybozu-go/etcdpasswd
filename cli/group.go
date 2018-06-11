package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type group struct{}

func (c group) SetFlags(f *flag.FlagSet) {}

func (c group) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	newc := NewCommander(f, "group")
	newc.Register(groupListCommand(), "")
	newc.Register(groupAddCommand(), "")
	newc.Register(groupRemoveCommand(), "")
	return newc.Execute(ctx)
}

// GroupCommand implements "group" subcommand.
func GroupCommand() subcommands.Command {
	return subcmd{
		group{},
		"group",
		"manage group database",
		"group ACTION ...",
	}
}

type groupList struct{}

func (c groupList) SetFlags(f *flag.FlagSet) {}

func (c groupList) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	groups, err := client.ListGroups(ctx)
	if err != nil {
		return handleError(err)
	}

	for _, g := range groups {
		fmt.Printf("%s (%d)\n", g.Name, g.GID)
	}
	return handleError(nil)
}

func groupListCommand() subcommands.Command {
	return subcmd{
		groupList{},
		"list",
		"list groups",
		"list",
	}
}

type groupAdd struct{}

func (c groupAdd) SetFlags(f *flag.FlagSet) {}

func (c groupAdd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}

	name := f.Arg(0)
	err := client.AddGroup(ctx, name)
	return handleError(err)
}

func groupAddCommand() subcommands.Command {
	return subcmd{
		groupAdd{},
		"add",
		"add a group",
		"add NAME",
	}
}

type groupRemove struct{}

func (c groupRemove) SetFlags(f *flag.FlagSet) {}

func (c groupRemove) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}

	name := f.Arg(0)
	err := client.RemoveGroup(ctx, name)
	return handleError(err)
}

func groupRemoveCommand() subcommands.Command {
	return subcmd{
		groupRemove{},
		"remove",
		"remove a group",
		"remove NAME",
	}
}
