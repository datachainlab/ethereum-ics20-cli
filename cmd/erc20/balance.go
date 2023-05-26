package erc20

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/erc20"
	"github.com/spf13/cobra"
)

func balanceCmd() *cobra.Command {
	var rpcAddress string
	var walletAddress string
	var denom string

	cmd := &cobra.Command{
		Use:   "balance",
		Short: "balance of the account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			balance, err := balance(rpcAddress, walletAddress, denom)
			if err != nil {
				return err
			}
			fmt.Printf("%d\n", balance)
			return nil
		},
	}
	cmd.Flags().StringVar(&rpcAddress, "rpc-address", "", "Ethereum RPC Address")
	cmd.Flags().StringVar(&walletAddress, "wallet-address", "", "Wallet address")
	cmd.Flags().StringVar(&denom, "denom", "", "Token denom")

	cmd.MarkFlagRequired("rpc-address")
	cmd.MarkFlagRequired("wallet-address")
	cmd.MarkFlagRequired("denom")

	return cmd
}

func balance(rpcAddress, walletAddress, denom string) (*big.Int, error) {
	ethClient, err := client.NewETHClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	erc20Token, err := erc20.NewErc20(common.HexToAddress(denom), ethClient)
	if err != nil {
		return nil, err
	}
	balance, err := erc20Token.BalanceOf(nil, common.HexToAddress(walletAddress))
	if err != nil {
		return nil, err
	}
	return balance, nil
}
