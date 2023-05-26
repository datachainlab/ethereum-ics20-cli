package cmd

import (
	"github.com/datachainlab/ethereum-ics20-cli/cmd/erc20"
	"github.com/datachainlab/ethereum-ics20-cli/cmd/ics20"
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
		erc20.Erc20Cmd(),
		ics20.Ics20Cmd(),
	)

	return rootCmd.Execute()
}
