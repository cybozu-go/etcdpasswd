package cmd

import (
	"bytes"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var certAddCmd = &cobra.Command{
	Use:   "add NAME [FILE]",
	Short: "add a SSH public key of a user",
	Long: `add a SSH public key of a user.

If FILE is not specified, public key is read from stdin.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		var file string
		switch len(args) {
		case 2:
			file = args[1]
			fallthrough
		case 1:
			name = args[0]
		}

		input := os.Stdin
		if len(file) > 0 {
			g, err := os.Open(file)
			if err != nil {
				return err
			}
			defer g.Close()
			input = g
		}
		pubkey, err := io.ReadAll(input)
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		user, rev, err := client.GetUser(ctx, name)
		if err != nil {
			return err
		}

		user.PubKeys = append(user.PubKeys, string(bytes.TrimSpace(pubkey)))
		return client.UpdateUser(ctx, user, rev)
	},
}

func init() {
	certCmd.AddCommand(certAddCmd)
}
