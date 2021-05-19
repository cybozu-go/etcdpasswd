package cmd

import (
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "user subcommand",
	Long:  `user subcommand`,
}

func init() {
	rootCmd.AddCommand(userCmd)
}
