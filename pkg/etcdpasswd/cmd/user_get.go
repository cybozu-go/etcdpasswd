package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userGetCmd = &cobra.Command{
	Use:   "get NAME",
	Short: "get user information",
	Long:  "get user information.",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		user, _, err := client.GetUser(cmd.Context(), name)
		if err != nil {
			return err
		}

		fmt.Printf(`uid: %d
display-name: %s
group: %s
groups: %v
shell: %s
public-keys: %d
`, user.UID, user.DisplayName, user.Group, user.Groups, user.Shell, len(user.PubKeys))
		return nil
	},
}

func init() {
	userCmd.AddCommand(userGetCmd)
}
