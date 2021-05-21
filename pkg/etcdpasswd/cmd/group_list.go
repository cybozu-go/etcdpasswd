package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "list groups",
	Long:  "list groups",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		groups, err := client.ListGroups(cmd.Context())
		if err != nil {
			return err
		}

		for _, g := range groups {
			fmt.Printf("%s (%d)\n", g.Name, g.GID)
		}
		return nil
	},
}

func init() {
	groupCmd.AddCommand(groupListCmd)
}
