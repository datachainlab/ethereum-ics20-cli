package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	var rootCmd = &cobra.Command{
		Use:   "ethereum-ics20-cli",
		Short: "ethereum-ics20-cli",
	}

	rootCmd.AddCommand(
		balanceCmd(),
		transferCmd(),
		walletCmd(),
	)

	return rootCmd.Execute()
}
