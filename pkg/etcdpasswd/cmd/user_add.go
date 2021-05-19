package cmd

import (
	"context"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/spf13/cobra"
)

var userAddConfig struct {
	displayName string
	group       string
	groups      commaStrings
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
			Groups:      []string(userAddConfig.groups),
			Shell:       userAddConfig.shell,
		}

		return client.AddUser(context.Background(), user)
	},
}

func init() {
	userCmd.AddCommand(userAddCmd)

	f := userAddCmd.Flags()
	f.StringVar(&userAddConfig.displayName, "display", "", "display name")
	f.StringVar(&userAddConfig.group, "group", "", "primary group")
	f.Var(&userAddConfig.groups, "groups", "comma-separated supplementary groups")
	f.StringVar(&userAddConfig.shell, "shell", "", "shell program")
}
