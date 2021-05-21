package cmd

import (
	"github.com/spf13/cobra"
)

var userUpdateConfig struct {
	displayName string
	group       string
	groups      []string
	shell       string
}

var userUpdateCmd = &cobra.Command{
	Use:   "update [OPTIONS] NAME",
	Short: "update an existing user",
	Long:  "update an existing user.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]
		user, rev, err := client.GetUser(ctx, name)
		if err != nil {
			return err
		}

		if len(userUpdateConfig.displayName) > 0 {
			user.DisplayName = userUpdateConfig.displayName
		}
		if len(userUpdateConfig.group) > 0 {
			user.Group = userUpdateConfig.group
		}
		if len(userUpdateConfig.groups) > 0 {
			user.Groups = userUpdateConfig.groups
		}
		if len(userUpdateConfig.shell) > 0 {
			user.Shell = userUpdateConfig.shell
		}

		return client.UpdateUser(ctx, user, rev)
	},
}

func init() {
	userCmd.AddCommand(userUpdateCmd)

	f := userUpdateCmd.Flags()
	f.StringVar(&userUpdateConfig.displayName, "display", "", "display name")
	f.StringVar(&userUpdateConfig.group, "group", "", "primary group")
	f.StringSliceVar(&userUpdateConfig.groups, "groups", []string{}, "comma-separated supplementary groups")
	f.StringVar(&userUpdateConfig.shell, "shell", "", "shell program")
}
