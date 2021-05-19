package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var userRemoveCmd = &cobra.Command{
	Use:   "remove NAME",
	Short: "remove an existing user",
	Long:  "remove an existing user.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return client.RemoveUser(context.Background(), name)
	},
}

func init() {
	userCmd.AddCommand(userRemoveCmd)
}
