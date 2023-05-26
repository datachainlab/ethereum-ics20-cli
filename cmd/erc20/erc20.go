package erc20

import (
	"github.com/spf13/cobra"
)

func Erc20Cmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "erc20",
		Short: "erc20 command",
	}

	cmd.AddCommand(balanceCmd())
	cmd.AddCommand(transferCmd())

	return cmd
}
