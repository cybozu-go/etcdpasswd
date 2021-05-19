package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "list users",
	Long:  "list users.",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := client.ListUsers(context.Background())
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
	userCmd.AddCommand(userListCmd)
}
