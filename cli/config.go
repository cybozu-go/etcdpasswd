package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/google/subcommands"
)

type configSet struct{}

func (c configSet) SetFlags(f *flag.FlagSet) {}

func (c configSet) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 2 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	configName := f.Arg(0)
	configValue := f.Arg(1)

RETRY:
	config, rev, err := client.GetConfig(ctx)
	if err != nil {
		return handleError(err)
	}

	switch configName {
	case "start-uid":
		v, err := strconv.Atoi(configValue)
		if err != nil {
			return handleError(err)
		}
		config.StartUID = v
	case "start-gid":
		v, err := strconv.Atoi(configValue)
		if err != nil {
			return handleError(err)
		}
		config.StartGID = v
	case "default-group":
		config.DefaultGroup = configValue
	case "default-groups":
		config.DefaultGroups = strings.Split(configValue, ",")
	case "default-shell":
		config.DefaultShell = configValue
	default:
		return handleError(errors.New("unknown config: " + configName))
	}

	err = client.SetConfig(ctx, config, rev)
	if err == etcdpasswd.ErrCASFailure {
		goto RETRY
	}
	return handleError(err)
}

// SetCommand implements "set" subcommand.
func SetCommand() subcommands.Command {
	return subcmd{
		configSet{},
		"set",
		"set configuration",
		fmt.Sprintf(`Usage: %s set CONFIG VALUE

CONFIG is one of:
    start-uid:      the beginning UID assigned to managed users.
    start-gid:      the beginning GID assigned to managed users.
    default-group:  default primary group.
    default-groups: comma-separated list of supplementary groups.
    default-shell:  default shell program.
`, os.Args[0]),
	}
}

type configGet struct{}

func (c configGet) SetFlags(f *flag.FlagSet) {}

func (c configGet) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}
	configName := f.Arg(0)

	config, _, err := client.GetConfig(ctx)
	if err != nil {
		return handleError(err)
	}

	switch configName {
	case "start-uid":
		fmt.Println(config.StartUID)
	case "start-gid":
		fmt.Println(config.StartGID)
	case "default-group":
		fmt.Println(config.DefaultGroup)
	case "default-groups":
		for _, g := range config.DefaultGroups {
			fmt.Println(g)
		}
	case "default-shell":
		fmt.Println(config.DefaultShell)
	default:
		return handleError(errors.New("unknown config: " + configName))
	}

	return handleError(nil)
}

// GetCommand implements "get" subcommand.
func GetCommand() subcommands.Command {
	return subcmd{
		configGet{},
		"get",
		"get configurations",
		fmt.Sprintf(`Usage: %s get CONFIG

CONFIG is one of:
    start-uid:      the beginning UID assigned to managed users.
    start-gid:      the beginning GID assigned to managed users.
    default-group:  default primary group.
    default-groups: comma-separated list of supplementary groups.
    default-shell:  default shell program.
`, os.Args[0]),
	}
}
