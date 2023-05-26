package ics20

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20bank"
	"github.com/spf13/cobra"
)

func balanceCmd() *cobra.Command {
	var rpcAddress, ics20BankAddress string
	var walletAddress string
	var denom string

	cmd := &cobra.Command{
		Use:   "balance",
		Short: "balance of the account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			balance, err := balacne(rpcAddress, ics20BankAddress, walletAddress, denom)
			if err != nil {
				return err
			}
			fmt.Printf("%d\n", balance)
			return nil
		},
	}
	cmd.Flags().StringVar(&rpcAddress, "rpc-address", "", "Ethereum RPC Address")
	cmd.Flags().StringVar(&ics20BankAddress, "ics20-bank-address", "", "Ics20Bank contract address")
	cmd.Flags().StringVar(&walletAddress, "wallet-address", "", "Wallet address")
	cmd.Flags().StringVar(&denom, "denom", "", "Token denom")

	cmd.MarkFlagRequired("rpc-address")
	cmd.MarkFlagRequired("ics20-bank-address")
	cmd.MarkFlagRequired("wallet-address")
	cmd.MarkFlagRequired("denom")

	return cmd
}

func balacne(rpcAddress, ics20BankAddress, walletAddress, denom string) (*big.Int, error) {
	ethClient, err := client.NewETHClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	ics20bank, err := ics20bank.NewIcs20bank(common.HexToAddress(ics20BankAddress), ethClient)
	if err != nil {
		return nil, err
	}
	balance, err := ics20bank.BalanceOf(nil, common.HexToAddress(walletAddress), denom)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
