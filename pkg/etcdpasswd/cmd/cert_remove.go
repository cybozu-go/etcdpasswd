package cmd

import (
	"context"
	"errors"
	"strconv"

	"github.com/spf13/cobra"
)

var certRemoveCmd = &cobra.Command{
	Use:   "remove NAME INDEX",
	Short: "remove a SSH public key of a user",
	Long:  "remove a SSH public key of a user.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		name := args[0]
		index, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		user, rev, err := client.GetUser(ctx, name)
		if err != nil {
			return err
		}

		if index >= len(user.PubKeys) || index < 0 {
			return errors.New("invalid index")
		}

		user.PubKeys = append(user.PubKeys[:index], user.PubKeys[(index+1):]...)
		return client.UpdateUser(ctx, user, rev)
	},
}

func init() {
	certCmd.AddCommand(certRemoveCmd)
}
