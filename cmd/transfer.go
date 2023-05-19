package cmd

import (
	"context"
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/datachainlab/ethereum-ics20-cli/chains/geth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
)

func transferCmd() *cobra.Command {
	var rpcAddress, mnemonic, ics20BankAddress, ics20TransferBankAddress string
	var fromIndex uint32
	var toAddress string
	var amount int64
	var tokenAddress string
	var portID string
	var channelID string
	var timeoutHeight uint64
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer token from one account to another chain's wallet",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			if err := Transfer(ctx, rpcAddress, mnemonic, ics20BankAddress, ics20TransferBankAddress, uint32(fromIndex), toAddress, amount, tokenAddress, portID, channelID, timeoutHeight); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&rpcAddress, "rpc-address", "", "config file path")
	cmd.Flags().StringVar(&mnemonic, "mnemonic", "", "mnemonic phrase")
	cmd.Flags().StringVar(&ics20BankAddress, "ics20-bank-address", "", "address of ics20 bank contract")
	cmd.Flags().StringVar(&ics20TransferBankAddress, "ics20-transfer-bank-address", "", "address of ics20 transfer bank contract")
	cmd.Flags().Uint32Var(&fromIndex, "from-index", 0, "index of the from wallet")
	cmd.Flags().StringVar(&toAddress, "to-address", "", "address of the recipient")
	cmd.Flags().Int64Var(&amount, "amount", 0, "amount of the token")
	cmd.Flags().StringVar(&tokenAddress, "token-address", "", "address of the token contract")
	cmd.Flags().StringVar(&portID, "port-id", "", "port id")
	cmd.Flags().StringVar(&channelID, "channel-id", "", "channel id")
	cmd.Flags().Uint64Var(&timeoutHeight, "timeout-height", 0, "timeout height")

	cmd.MarkFlagRequired("rpc-address")
	cmd.MarkFlagRequired("chain-id")
	cmd.MarkFlagRequired("mnemonic")
	cmd.MarkFlagRequired("ics20-bank-address")
	cmd.MarkFlagRequired("ics20-transfer-bank-address")
	cmd.MarkFlagRequired("from-index")
	cmd.MarkFlagRequired("to-address")
	cmd.MarkFlagRequired("amount")
	cmd.MarkFlagRequired("token-address")
	cmd.MarkFlagRequired("port-id")
	cmd.MarkFlagRequired("channel-id")
	cmd.MarkFlagRequired("timeout-height")

	return cmd
}

func Transfer(ctx context.Context, rpcAddress string, mnemonic, ics20BankAddress, ics20TransferBankAddress string, fromIndex uint32, toAddress string, amount int64, tokenAddress, portID, channelID string, timeoutHeight uint64) error {
	chain, err := geth.InitializeChain(ctx, rpcAddress, mnemonic, tokenAddress, ics20TransferBankAddress, ics20BankAddress)
	if err != nil {
		return err
	}
	const (
		relayer  = 0
		deployer = 0
	)
	tx, err := chain.SimpleToken.Approve(chain.TxOpts(ctx, deployer), common.HexToAddress(ics20BankAddress), big.NewInt(amount))
	if err != nil {
		return err
	}
	if err := waitAndCheckStatus(ctx, chain, tx); err != nil {
		return err
	}
	log.Printf("1. token approve success (TxHash: %s)\n", tx.Hash().Hex())

	tx, err = chain.ICS20Bank.Deposit(
		chain.TxOpts(ctx, deployer),
		common.HexToAddress(tokenAddress),
		big.NewInt(amount),
		chain.CallOpts(ctx, fromIndex).From,
	)
	if err != nil {
		return err
	}
	if err := waitAndCheckStatus(ctx, chain, tx); err != nil {
		return err
	}
	log.Printf("2. deposit success (TxHash: %s)\n", tx.Hash().Hex())

	baseDenom := strings.ToLower(tokenAddress)
	tx, err = chain.ICS20Transfer.SendTransfer(
		chain.TxOpts(ctx, fromIndex),
		baseDenom,
		uint64(amount),
		common.HexToAddress(toAddress),
		portID, channelID,
		timeoutHeight,
	)
	if err != nil {
		return err
	}
	if err := waitAndCheckStatus(ctx, chain, tx); err != nil {
		return err
	}
	log.Printf("3. sendTransfer success (TxHash: %s)\n", tx.Hash().Hex())

	return nil
}

func waitAndCheckStatus(ctx context.Context, chain *geth.Chain, tx *types.Transaction) error {
	receipt, err := chain.Client.WaitForReceiptAndGet(ctx, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return errors.New("tx status error")
	}
	return nil
}
