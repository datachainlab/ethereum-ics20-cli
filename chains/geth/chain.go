package geth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/erc20"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20bank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20transferbank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
)

type Chain struct {
	chainID        int64
	mnemonicPhrase string
	keys           map[uint32]*ecdsa.PrivateKey

	Client        *client.ETHClient
	Erc20Token    erc20.Erc20
	ICS20Transfer ics20transferbank.Ics20transferbank
	ICS20Bank     ics20bank.Ics20bank
}

func NewChain(client *client.ETHClient, ethChainId int64, mnemonic string) *Chain {
	return &Chain{
		chainID:        ethChainId,
		mnemonicPhrase: mnemonic,
		keys:           make(map[uint32]*ecdsa.PrivateKey),

		Client: client,
	}
}

func (chain *Chain) AddERC20Token(erc20TokenAddress string) error {
	erc20Token, err := erc20.NewErc20(common.HexToAddress(erc20TokenAddress), chain.Client)
	if err != nil {
		return err
	}
	chain.Erc20Token = *erc20Token
	return nil
}

func (chain *Chain) AddICS20TransferBank(ics20TransferBankAddress string) error {
	ics20transfer, err := ics20transferbank.NewIcs20transferbank(common.HexToAddress(ics20TransferBankAddress), chain.Client)
	if err != nil {
		return err
	}
	chain.ICS20Transfer = *ics20transfer
	return nil
}

func (chain *Chain) AddICS20Bank(ics20BankAddress string) error {
	ics20bank, err := ics20bank.NewIcs20bank(common.HexToAddress(ics20BankAddress), chain.Client)
	if err != nil {
		return err
	}
	chain.ICS20Bank = *ics20bank
	return nil
}

func (chain *Chain) TxOpts(ctx context.Context, index uint32) *bind.TransactOpts {
	return makeGenTxOpts(big.NewInt(chain.chainID), chain.prvKey(index))(ctx)
}

func (chain *Chain) prvKey(index uint32) *ecdsa.PrivateKey {
	key, ok := chain.keys[index]
	if ok {
		return key
	}
	key, err := wallet.GetPrvKeyFromMnemonicAndHDWPath(chain.mnemonicPhrase, fmt.Sprintf("m/44'/60'/0'/0/%v", index))
	if err != nil {
		panic(err)
	}
	chain.keys[index] = key
	return key
}

func (chain *Chain) CallOpts(ctx context.Context, index uint32) *bind.CallOpts {
	opts := chain.TxOpts(ctx, index)
	return &bind.CallOpts{
		From:    opts.From,
		Context: opts.Context,
	}
}

func (chain *Chain) WaitAndCheckStatus(ctx context.Context, tx *types.Transaction) error {
	receipt, err := chain.Client.WaitForReceiptAndGet(ctx, tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return errors.New("tx status error")
	}
	return nil
}

func makeGenTxOpts(chainID *big.Int, prv *ecdsa.PrivateKey) func(ctx context.Context) *bind.TransactOpts {
	signer := types.LatestSignerForChainID(chainID)
	addr := gethcrypto.PubkeyToAddress(prv.PublicKey)
	return func(ctx context.Context) *bind.TransactOpts {
		return &bind.TransactOpts{
			From:     addr,
			GasLimit: 6382056,
			Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				if address != addr {
					return nil, errors.New("not authorized to sign this account")
				}
				signature, err := gethcrypto.Sign(signer.Hash(tx).Bytes(), prv)
				if err != nil {
					return nil, err
				}
				return tx.WithSignature(signer, signature)
			},
		}
	}
}

func InitializeChain(ctx context.Context, rpcAddress string, mnemonic string) (*Chain, error) {
	ethClient, err := client.NewETHClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	ethChainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	chain := NewChain(ethClient, ethChainID.Int64(), mnemonic)

	return chain, nil
}
