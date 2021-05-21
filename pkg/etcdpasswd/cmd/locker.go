package cmd

import (
	"github.com/spf13/cobra"
)

var lockerCmd = &cobra.Command{
	Use:   "locker",
	Short: "locker subcommand",
	Long:  "locker subcommand",
}

func init() {
	rootCmd.AddCommand(lockerCmd)
}
