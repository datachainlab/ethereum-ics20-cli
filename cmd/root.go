package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	var rootCmd = &cobra.Command{
		Use:   "ethereum-ics20-cli",
		Short: "command line tool for ethereum ics20 token",
	}

	rootCmd.AddCommand(
		balanceCmd(),
		transferCmd(),
		walletCmd(),
		chainCmd(),
	)

	return rootCmd.Execute()
}
