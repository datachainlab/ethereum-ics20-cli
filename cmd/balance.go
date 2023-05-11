package cmd

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20bank"
	"github.com/spf13/cobra"
)

func balanceCmd() *cobra.Command {
	var rpcAddress, ics20BankAddress, walletAddress, tokenAddress string
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Query the account balance of the address",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			balance, err := balanceOf(rpcAddress, ics20BankAddress, walletAddress, tokenAddress)
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
	cmd.Flags().StringVar(&tokenAddress, "token-address", "", "Token address")

	return cmd
}

func balanceOf(rpcAddress, ics20BankAddress, walletAddress, tokenAddress string) (*big.Int, error) {
	baseDenom := strings.ToLower(tokenAddress)
	ethClient, err := client.NewETHClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	ics20bank, err := ics20bank.NewIcs20bank(common.HexToAddress(ics20BankAddress), ethClient)
	if err != nil {
		return nil, err
	}
	balance, err := ics20bank.BalanceOf(nil, common.HexToAddress(walletAddress), baseDenom)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
