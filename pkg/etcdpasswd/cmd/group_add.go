package cmd

import (
	"github.com/spf13/cobra"
)

var groupAddCmd = &cobra.Command{
	Use:   "add GROUP",
	Short: "add a group",
	Long:  "add a group.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return client.AddGroup(cmd.Context(), name)
	},
}

func init() {
	groupCmd.AddCommand(groupAddCmd)
}
