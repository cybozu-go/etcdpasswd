package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var lockerListCmd = &cobra.Command{
	Use:   "list",
	Short: "list password-locked users",
	Long:  "list password-locked users.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := client.ListLocked(context.Background())
		if err != nil {
			return err
		}

		for _, u := range users {
			fmt.Println(u)
		}
		return nil
	},
}

func init() {
	lockerCmd.AddCommand(lockerListCmd)
}
