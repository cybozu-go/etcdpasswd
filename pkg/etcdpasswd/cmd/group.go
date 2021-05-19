package cmd

import (
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "group subcommand",
	Long:  "group subcommand",
}

func init() {
	rootCmd.AddCommand(groupCmd)
}
