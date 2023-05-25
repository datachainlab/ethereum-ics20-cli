package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/datachainlab/ethereum-ics20-cli/chains/geth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/erc20"
	"github.com/spf13/cobra"
)

func erc20Cmd() *cobra.Command {
	var rpcAddress, mnemonic, ics20BankAddress, ics20TransferBankAddress string
	var fromIndex uint32
	var toAddress string
	var amount int64
	var denom string
	var walletAddress string

	cmd := &cobra.Command{
		Use:   "erc20",
		Short: "erc20 command",
	}

	// balance command
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "balance of the account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			balance, err := erc20Balance(rpcAddress, walletAddress, denom)
			if err != nil {
				return err
			}
			fmt.Printf("%d\n", balance)
			return nil
		},
	}
	balanceCmd.Flags().StringVar(&rpcAddress, "rpc-address", "", "Ethereum RPC Address")
	balanceCmd.Flags().StringVar(&walletAddress, "wallet-address", "", "Wallet address")
	balanceCmd.Flags().StringVar(&denom, "denom", "", "Token denom")

	balanceCmd.MarkFlagRequired("rpc-address")
	balanceCmd.MarkFlagRequired("wallet-address")
	balanceCmd.MarkFlagRequired("denom")

	cmd.AddCommand(balanceCmd)

	// transfer command
	transferCmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer erc20 token from fromIndex account to toAddress account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			if err := erc20Transfer(ctx, rpcAddress, mnemonic, ics20BankAddress, ics20TransferBankAddress, uint32(fromIndex), toAddress, amount, denom); err != nil {
				return err
			}
			return nil
		},
	}
	transferCmd.Flags().StringVar(&rpcAddress, "rpc-address", "", "Ethereum RPC Address")
	transferCmd.Flags().StringVar(&mnemonic, "mnemonic", "", "mnemonic phrase")
	transferCmd.Flags().StringVar(&ics20BankAddress, "ics20-bank-address", "", "address of ics20 bank contract")
	transferCmd.Flags().StringVar(&ics20TransferBankAddress, "ics20-transfer-bank-address", "", "address of ics20 transfer bank contract")
	transferCmd.Flags().Uint32Var(&fromIndex, "from-index", 0, "index of the from wallet")
	transferCmd.Flags().StringVar(&toAddress, "to-address", "", "address of the recipient")
	transferCmd.Flags().Int64Var(&amount, "amount", 0, "amount of the token")
	transferCmd.Flags().StringVar(&denom, "denom", "", "denom of the token")

	transferCmd.MarkFlagRequired("rpc-address")
	transferCmd.MarkFlagRequired("mnemonic")
	transferCmd.MarkFlagRequired("ics20-bank-address")
	transferCmd.MarkFlagRequired("ics20-transfer-bank-address")
	transferCmd.MarkFlagRequired("from-index")
	transferCmd.MarkFlagRequired("to-address")
	transferCmd.MarkFlagRequired("amount")
	transferCmd.MarkFlagRequired("denom")

	cmd.AddCommand(transferCmd)

	return cmd
}

func erc20Balance(rpcAddress, walletAddress, denom string) (*big.Int, error) {
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

func erc20Transfer(ctx context.Context, rpcAddress string, mnemonic, ics20BankAddress, ics20TransferBankAddress string, fromIndex uint32, toAddress string, amount int64, denom string) error {
	chain, err := geth.InitializeChain(ctx, rpcAddress, mnemonic, denom, ics20TransferBankAddress, ics20BankAddress)
	if err != nil {
		return err
	}
	tx, err := chain.Erc20Token.Transfer(chain.TxOpts(ctx, fromIndex), common.HexToAddress(toAddress), big.NewInt(amount))
	if err != nil {
		return err
	}
	if err := geth.WaitAndCheckStatus(ctx, chain, tx); err != nil {
		return err
	}
	log.Printf("ERC20 token Transfer success (TxHash: %s)\n", tx.Hash().Hex())

	return nil
}
