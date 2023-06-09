package erc20

import (
	"context"
	"log"
	"math/big"

	"github.com/datachainlab/ethereum-ics20-cli/chains/geth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/erc20"
	"github.com/spf13/cobra"
)

func transferCmd() *cobra.Command {
	var rpcAddress, mnemonic string
	var fromIndex uint32
	var toAddress string
	var amount int64
	var denom string

	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer erc20 token from fromIndex account to toAddress account",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()
			if err := transfer(ctx, rpcAddress, mnemonic, uint32(fromIndex), toAddress, amount, denom); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&rpcAddress, "rpc-address", "", "Ethereum RPC Address")
	cmd.Flags().StringVar(&mnemonic, "mnemonic", "", "mnemonic phrase")
	cmd.Flags().Uint32Var(&fromIndex, "from-index", 0, "index of the from wallet")
	cmd.Flags().StringVar(&toAddress, "to-address", "", "address of the recipient")
	cmd.Flags().Int64Var(&amount, "amount", 0, "amount of the token")
	cmd.Flags().StringVar(&denom, "denom", "", "denom of the token")

	cmd.MarkFlagRequired("rpc-address")
	cmd.MarkFlagRequired("mnemonic")
	cmd.MarkFlagRequired("from-index")
	cmd.MarkFlagRequired("to-address")
	cmd.MarkFlagRequired("amount")
	cmd.MarkFlagRequired("denom")

	return cmd
}

func transfer(ctx context.Context, rpcAddress string, mnemonic string, fromIndex uint32, toAddress string, amount int64, denom string) error {
	chain, err := geth.InitializeChain(ctx, rpcAddress, mnemonic)
	if err != nil {
		return err
	}
	erc20Token, err := erc20.NewErc20(common.HexToAddress(denom), chain.Client)
	if err != nil {
		return err
	}

	tx, err := erc20Token.Transfer(chain.TxOpts(ctx, fromIndex), common.HexToAddress(toAddress), big.NewInt(amount))
	if err != nil {
		return err
	}
	if err := chain.WaitAndCheckStatus(ctx, tx); err != nil {
		return err
	}
	log.Printf("ERC20 token Transfer success (TxHash: %s)\n", tx.Hash().Hex())

	return nil
}
