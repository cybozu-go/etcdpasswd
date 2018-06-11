package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type locker struct{}

func (c locker) SetFlags(f *flag.FlagSet) {}

func (c locker) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	newc := NewCommander(f, "locker")
	newc.Register(lockerListCommand(), "")
	newc.Register(lockerAddCommand(), "")
	newc.Register(lockerRemoveCommand(), "")
	return newc.Execute(ctx)
}

// LockerCommand implements "locker" subcommand.
func LockerCommand() subcommands.Command {
	return subcmd{
		locker{},
		"locker",
		"lock user passwords",
		"locker ACTION ...",
	}
}

type lockerList struct{}

func (c lockerList) SetFlags(f *flag.FlagSet) {}

func (c lockerList) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	users, err := client.ListLocked(ctx)
	if err != nil {
		return handleError(err)
	}

	for _, u := range users {
		fmt.Println(u)
	}
	return handleError(nil)
}

func lockerListCommand() subcommands.Command {
	return subcmd{
		lockerList{},
		"list",
		"list password-locked users",
		"locker list",
	}
}

type lockerAdd struct{}

func (c lockerAdd) SetFlags(f *flag.FlagSet) {}

func (c lockerAdd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	name := f.Arg(0)

	err := client.Lock(ctx, name)
	return handleError(err)
}

func lockerAddCommand() subcommands.Command {
	return subcmd{
		lockerAdd{},
		"add",
		"add NAME to the list of password-locked users",
		"locker add NAME",
	}
}

type lockerRemove struct{}

func (c lockerRemove) SetFlags(f *flag.FlagSet) {}

func (c lockerRemove) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	name := f.Arg(0)

	err := client.Unlock(ctx, name)
	return handleError(err)
}

func lockerRemoveCommand() subcommands.Command {
	return subcmd{
		lockerRemove{},
		"remove",
		"remove NAME from the list of password-locked users",
		"locker remove NAME",
	}
}
