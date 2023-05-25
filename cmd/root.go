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
		walletCmd(),
		chainCmd(),
		erc20Cmd(),
		ics20Cmd(),
	)

	return rootCmd.Execute()
}
