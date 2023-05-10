package cmd

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/datachainlab/ethereum-ics20-cli/chains/geth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func transferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer token",
		Long:  "Usage: transfer <configDir> <ics20BankAddress> <ics20TransferBankAddress> <fromIndex> <toAddress> <amount> <tokenAddress> <portID> <channelID> <timeout>",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			configDir := args[0]
			ics20BankAddress := args[1]
			ics20TransferBankAddress := args[2]
			fromIndex, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			toAddress := args[4]
			amount, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				return err
			}
			tokenAddress := args[6]
			portID := args[7]
			channelID := args[8]
			timeout, err := strconv.ParseUint(args[9], 10, 64)
			if err != nil {
				return err
			}
			if err := Transfer(configDir, ics20BankAddress, ics20TransferBankAddress, uint32(fromIndex), toAddress, amount, tokenAddress, portID, channelID, timeout); err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}

func Transfer(configDir, ics20BankAddress, ics20TransferBankAddress string, fromIndex uint32, toAddress string, amount int64, tokenAddress, portID, channelID string, timeout uint64) error {
	chainA, err := geth.InitializeChains(configDir, geth.PathSrc, tokenAddress, ics20TransferBankAddress, ics20BankAddress)
	if err != nil {
		return err
	}
	ctx := context.Background()
	const (
		relayer  = 0
		deployer = 0
	)
	_, err = chainA.SimpleToken.Approve(chainA.TxOpts(ctx, deployer), common.HexToAddress(ics20BankAddress), big.NewInt(amount))
	if err != nil {
		return err
	}
	log.Println("1. token approve success")

	_, err = chainA.ICS20Bank.Deposit(
		chainA.TxOpts(ctx, deployer),
		common.HexToAddress(tokenAddress),
		big.NewInt(amount),
		chainA.CallOpts(ctx, fromIndex).From,
	)
	if err != nil {
		return err
	}
	log.Println("2. deposit success")

	baseDenom := strings.ToLower(tokenAddress)
	_, err = chainA.ICS20Transfer.SendTransfer(
		chainA.TxOpts(ctx, fromIndex),
		baseDenom,
		uint64(amount),
		common.HexToAddress(toAddress),
		portID, channelID,
		timeout+1000,
	)
	if err != nil {
		return err
	}
	log.Println("3. sendTransfer success")

	return nil
}
