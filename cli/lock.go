package cli

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type lockCommand struct{}

func (c lockCommand) SetFlags(f *flag.FlagSet) {}

func (c lockCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	name := f.Arg(0)

	err := client.Lock(ctx, name)
	return handleError(err)
}

// LockCommand implements "lock" subcommand.
func LockCommand() subcommands.Command {
	return subcmd{
		lockCommand{},
		"lock",
		"add NAME to locked user database",
		"lock NAME",
	}
}

type unlockCommand struct{}

func (c unlockCommand) SetFlags(f *flag.FlagSet) {}

func (c unlockCommand) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	name := f.Arg(0)

	err := client.Unlock(ctx, name)
	return handleError(err)
}

// UnlockCommand implements "lock" subcommand.
func UnlockCommand() subcommands.Command {
	return subcmd{
		unlockCommand{},
		"unlock",
		"remove NAME from locked user database",
		"unlock NAME",
	}
}
