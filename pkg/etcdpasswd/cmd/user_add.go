package cmd

import (
	"github.com/cybozu-go/etcdpasswd"
	"github.com/spf13/cobra"
)

var userAddConfig struct {
	displayName string
	group       string
	groups      []string
	shell       string
}

var userAddCmd = &cobra.Command{
	Use:   "add [OPTIONS] NAME",
	Short: "add a new user",
	Long:  "add a new user.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		user := &etcdpasswd.User{
			Name:        name,
			DisplayName: userAddConfig.displayName,
			Group:       userAddConfig.group,
			Groups:      userAddConfig.groups,
			Shell:       userAddConfig.shell,
		}

		return client.AddUser(cmd.Context(), user)
	},
}

func init() {
	userCmd.AddCommand(userAddCmd)

	f := userAddCmd.Flags()
	f.StringVar(&userAddConfig.displayName, "display", "", "display name")
	f.StringVar(&userAddConfig.group, "group", "", "primary group")
	f.StringSliceVar(&userAddConfig.groups, "groups", []string{}, "comma-separated supplementary groups")
	f.StringVar(&userAddConfig.shell, "shell", "", "shell program")
}
