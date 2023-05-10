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
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "balance of the address",
		Long:  "Usage: balance <rpcAddress> <ics20BankAddress> <walletAddress> <tokenAddress>",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rpcAddress := args[0]
			ics20BankAddress := args[1]
			walletAddress := args[2]
			tokenAddress := args[3]
			balance, err := balanceOf(rpcAddress, ics20BankAddress, walletAddress, tokenAddress)
			if err != nil {
				return err
			}
			fmt.Printf("%d", balance)
			return nil
		},
	}
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
