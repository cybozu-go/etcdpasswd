package cmd

import (
	"github.com/spf13/cobra"
)

var lockerRemoveCmd = &cobra.Command{
	Use:   "remove NAME",
	Short: "remove NAME from the list of password-locked users",
	Long:  "remove NAME from the list of password-locked users.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return client.Unlock(cmd.Context(), name)
	},
}

func init() {
	lockerCmd.AddCommand(lockerRemoveCmd)
}
