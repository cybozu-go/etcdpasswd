package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

func pprintPubKey(pubkey string) string {
	pk, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(pubkey))
	if err != nil {
		return pubkey
	}
	return fmt.Sprintf("%s (%s)", comment, pk.Type())
}

var certListCmd = &cobra.Command{
	Use:   "list NAME",
	Short: "list SSH public keys of a user",
	Long:  "list SSH public keys of a user.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		user, _, err := client.GetUser(context.Background(), name)
		if err != nil {
			return err
		}

		for i, pubkey := range user.PubKeys {
			fmt.Printf("%d: %s\n", i, pprintPubKey(pubkey))
		}
		return nil
	},
}

func init() {
	certCmd.AddCommand(certListCmd)
}
