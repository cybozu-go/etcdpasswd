package cmd

import (
	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "cert subcommand",
	Long:  "cert subcommand",
}

func init() {
	rootCmd.AddCommand(certCmd)
}
