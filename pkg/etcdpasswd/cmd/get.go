package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get CONFIG",
	Short: "get configurations",
	Long: `Usage: etcdpasswd get CONFIG

CONFIG is one of:
    start-uid:      the beginning UID assigned to managed users.
    start-gid:      the beginning GID assigned to managed users.
    default-group:  default primary group.
    default-groups: comma-separated list of supplementary groups.
    default-shell:  default shell program.`,

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configName := args[0]

		config, _, err := client.GetConfig(context.Background())
		if err != nil {
			return err
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
			return errors.New("unknown config: " + configName)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
