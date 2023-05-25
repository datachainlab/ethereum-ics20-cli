package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/datachainlab/ethereum-ics20-cli/chains/geth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20bank"
	"github.com/spf13/cobra"
)

func ics20Cmd() *cobra.Command {
	var rpcAddress, mnemonic, ics20BankAddress, ics20TransferBankAddress string
	var fromIndex uint32
	var toAddress string
	var amount int64
	var denom string
	var portID string
	var channelID string
	var timeoutHeight uint64
	var walletAddress string

	cmd := &cobra.Command{
		Use:   "ics20",
		Short: "ics20 command",
	}

	// balance command
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "balance of the account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			balance, err := ics20Balance(rpcAddress, ics20BankAddress, walletAddress, denom)
			if err != nil {
				return err
			}
			fmt.Printf("%d\n", balance)
			return nil
		},
	}
	balanceCmd.Flags().StringVar(&rpcAddress, "rpc-address", "", "Ethereum RPC Address")
	balanceCmd.Flags().StringVar(&ics20BankAddress, "ics20-bank-address", "", "Ics20Bank contract address")
	balanceCmd.Flags().StringVar(&walletAddress, "wallet-address", "", "Wallet address")
	balanceCmd.Flags().StringVar(&denom, "denom", "", "Token denom")

	balanceCmd.MarkFlagRequired("rpc-address")
	balanceCmd.MarkFlagRequired("ics20-bank-address")
	balanceCmd.MarkFlagRequired("wallet-address")
	balanceCmd.MarkFlagRequired("denom")

	cmd.AddCommand(balanceCmd)

	// transfer command
	transferCmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer token from one account to another chain's account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			if err := ics20Transfer(ctx, rpcAddress, mnemonic, ics20BankAddress, ics20TransferBankAddress, uint32(fromIndex), toAddress, amount, denom, portID, channelID, timeoutHeight); err != nil {
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
	transferCmd.Flags().StringVar(&portID, "port-id", "", "port id")
	transferCmd.Flags().StringVar(&channelID, "channel-id", "", "channel id")
	transferCmd.Flags().Uint64Var(&timeoutHeight, "timeout-height", 0, "timeout height")

	transferCmd.MarkFlagRequired("rpc-address")
	transferCmd.MarkFlagRequired("chain-id")
	transferCmd.MarkFlagRequired("mnemonic")
	transferCmd.MarkFlagRequired("ics20-bank-address")
	transferCmd.MarkFlagRequired("ics20-transfer-bank-address")
	transferCmd.MarkFlagRequired("from-index")
	transferCmd.MarkFlagRequired("to-address")
	transferCmd.MarkFlagRequired("amount")
	transferCmd.MarkFlagRequired("denom")
	transferCmd.MarkFlagRequired("port-id")
	transferCmd.MarkFlagRequired("channel-id")
	transferCmd.MarkFlagRequired("timeout-height")

	cmd.AddCommand(transferCmd)

	return cmd
}

func ics20Balance(rpcAddress, ics20BankAddress, walletAddress, denom string) (*big.Int, error) {
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

func ics20Transfer(ctx context.Context, rpcAddress string, mnemonic, ics20BankAddress, ics20TransferBankAddress string, fromIndex uint32, toAddress string, amount int64, denom, portID, channelID string, timeoutHeight uint64) error {
	chain, err := geth.InitializeChain(ctx, rpcAddress, mnemonic, denom, ics20TransferBankAddress, ics20BankAddress)
	if err != nil {
		return err
	}

	if common.IsHexAddress(denom) {
		tx, err := chain.Erc20Token.Approve(chain.TxOpts(ctx, fromIndex), common.HexToAddress(ics20BankAddress), big.NewInt(amount))
		if err != nil {
			return err
		}
		if err := geth.WaitAndCheckStatus(ctx, chain, tx); err != nil {
			return err
		}
		log.Printf("Token approve success (TxHash: %s)\n", tx.Hash().Hex())

		tx, err = chain.ICS20Bank.Deposit(
			chain.TxOpts(ctx, fromIndex),
			common.HexToAddress(denom),
			big.NewInt(amount),
			chain.CallOpts(ctx, fromIndex).From,
		)
		if err != nil {
			return err
		}
		if err := geth.WaitAndCheckStatus(ctx, chain, tx); err != nil {
			return err
		}
		log.Printf("Deposit success (TxHash: %s)\n", tx.Hash().Hex())
	}

	tx, err := chain.ICS20Transfer.SendTransfer(
		chain.TxOpts(ctx, fromIndex),
		denom,
		uint64(amount),
		common.HexToAddress(toAddress),
		portID, channelID,
		timeoutHeight,
	)
	if err != nil {
		return err
	}
	if err := geth.WaitAndCheckStatus(ctx, chain, tx); err != nil {
		return err
	}
	log.Printf("SendTransfer success (TxHash: %s)\n", tx.Hash().Hex())

	return nil
}
