package ics20

import (
	"github.com/spf13/cobra"
)

func Ics20Cmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "ics20",
		Short: "ics20 command",
	}

	cmd.AddCommand(balanceCmd())
	cmd.AddCommand(transferCmd())

	return cmd
}
