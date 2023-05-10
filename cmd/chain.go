package cmd

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/spf13/cobra"
)

func chainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain",
		Short: "chain info",
	}
	heightCmd := &cobra.Command{
		Use:   "height",
		Short: "get the current height of the blockchain",
		Long:  "Usage: height <rpcAddress>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rpcAddress := args[0]
			headerNumber, err := headerNumber(rpcAddress)
			if err != nil {
				return err
			}
			fmt.Printf("%d", headerNumber)
			return nil
		},
	}
	cmd.AddCommand(heightCmd)
	return cmd
}

func headerNumber(rpcAddress string) (uint64, error) {
	ctx := context.Background()
	ethClient, err := client.NewETHClient(rpcAddress)
	if err != nil {
		return 0, err
	}
	header, err := ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}
