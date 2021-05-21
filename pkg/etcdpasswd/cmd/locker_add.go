package cmd

import (
	"github.com/spf13/cobra"
)

var lockerAddCmd = &cobra.Command{
	Use:   "add NAME",
	Short: "add NAME to the list of password-locked users",
	Long:  "add NAME to the list of password-locked users.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return client.Lock(cmd.Context(), name)
	},
}

func init() {
	lockerCmd.AddCommand(lockerAddCmd)
}
