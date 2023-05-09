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

const (
	TransferUsage = "Usage: transfer <configDir> <ics20BankAddress> <ics20TransferBankAddress> <fromIndex> <toIndex> <amount> <tokenAddress>"
)

func transferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: TransferUsage,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 7 {
				log.Println(TransferUsage)
			}
			configDir := args[0]
			ics20BankAddress := args[1]
			ics20TransferBankAddress := args[2]
			fromIndex, err := strconv.ParseInt(args[3], 10, 32)
			if err != nil {
				log.Fatalln(err)
			}
			toIndex, err := strconv.ParseInt(args[4], 10, 32)
			if err != nil {
				log.Fatalln(err)
			}
			amount, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			tokenAddress := args[6]
			err = Transfer(configDir, ics20BankAddress, ics20TransferBankAddress, uint32(fromIndex), uint32(toIndex), amount, tokenAddress)
			if err != nil {
				log.Fatalln("transfer Error: ", err)
			}
		},
	}

	return cmd
}

func Transfer(configDir, ics20BankAddress, ics20TransferBankAddress string, fromIndex, toIndex uint32, amount int64, tokenAddress string) error {
	chainA, chainB, err := geth.InitializeChains(configDir, tokenAddress, ics20TransferBankAddress, ics20BankAddress)
	if err != nil {
		return err
	}
	ctx := context.Background()
	const (
		relayer  = 0
		deployer = 0
	)
	chanA := chainA.PathEnd
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

	var heightB uint64
	if header, err := chainB.LastHeader(ctx); err != nil {
		log.Fatalf("failed to the latest header of chain B: %v", err)
	} else {
		heightB = header.Number.Uint64()
	}
	baseDenom := strings.ToLower(tokenAddress)
	_, err = chainA.ICS20Transfer.SendTransfer(
		chainA.TxOpts(ctx, fromIndex),
		baseDenom,
		uint64(amount),
		chainB.CallOpts(ctx, toIndex).From,
		chanA.PortID, chanA.ChannelID,
		heightB+1000,
	)
	if err != nil {
		return err
	}
	log.Println("3. sendTransfer success")

	return nil
}
