package cmd

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set CONFIG VALUE",
	Short: "set configuration",
	Long: `Usage: etcdpasswd set CONFIG VALUE

CONFIG is one of:
    start-uid:      the beginning UID assigned to managed users.
    start-gid:      the beginning GID assigned to managed users.
    default-group:  default primary group.
    default-groups: comma-separated list of supplementary groups.
    default-shell:  default shell program.`,

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		configName := args[0]
		configValue := args[1]
		ctx := context.Background()

	RETRY:
		config, rev, err := client.GetConfig(ctx)
		if err != nil {
			return err
		}

		switch configName {
		case "start-uid":
			v, err := strconv.Atoi(configValue)
			if err != nil {
				return err
			}
			config.StartUID = v
		case "start-gid":
			v, err := strconv.Atoi(configValue)
			if err != nil {
				return err
			}
			config.StartGID = v
		case "default-group":
			config.DefaultGroup = configValue
		case "default-groups":
			config.DefaultGroups = strings.Split(configValue, ",")
		case "default-shell":
			config.DefaultShell = configValue
		default:
			return errors.New("unknown config: " + configName)
		}

		err = client.SetConfig(ctx, config, rev)
		if err == etcdpasswd.ErrCASFailure {
			goto RETRY
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
