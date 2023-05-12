package cmd

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/spf13/cobra"
)

func chainCmd() *cobra.Command {
	var rpcAddress string
	cmd := &cobra.Command{
		Use:   "chain",
		Short: "chain info",
	}
	heightCmd := &cobra.Command{
		Use:   "height",
		Short: "Get the current height of the blockchain",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			headerNumber, err := headerNumber(ctx, rpcAddress)
			if err != nil {
				return err
			}
			fmt.Printf("%d\n", headerNumber)
			return nil
		},
	}
	heightCmd.Flags().StringVar(&rpcAddress, "rpc-address", "", "Ethereum RPC Address")
	heightCmd.MarkFlagRequired("rpc-address")

	cmd.AddCommand(heightCmd)
	return cmd
}

func headerNumber(ctx context.Context, rpcAddress string) (uint64, error) {
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
