package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var groupRemoveCmd = &cobra.Command{
	Use:   "remove GROUP",
	Short: "remove a group",
	Long:  "remove a group.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return client.RemoveGroup(context.Background(), name)
	},
}

func init() {
	groupCmd.AddCommand(groupRemoveCmd)
}
