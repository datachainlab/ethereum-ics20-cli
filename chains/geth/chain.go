package geth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20bank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20transferbank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/simpletoken"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
)

type Chain struct {
	chainID        int64
	mnemonicPhrase string
	keys           map[uint32]*ecdsa.PrivateKey

	Client        *client.ETHClient
	SimpleToken   simpletoken.Simpletoken
	ICS20Transfer ics20transferbank.Ics20transferbank
	ICS20Bank     ics20bank.Ics20bank
}

func NewChain(client *client.ETHClient, ethChainId int64, mnemonic, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress string) *Chain {
	simpletoken, err := simpletoken.NewSimpletoken(common.HexToAddress(simpleTokenAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	ics20transfer, err := ics20transferbank.NewIcs20transferbank(common.HexToAddress(ics20TransferBankAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	ics20bank, err := ics20bank.NewIcs20bank(common.HexToAddress(ics20BankAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &Chain{
		chainID:        ethChainId,
		mnemonicPhrase: mnemonic,
		keys:           make(map[uint32]*ecdsa.PrivateKey),

		Client:        client,
		SimpleToken:   *simpletoken,
		ICS20Transfer: *ics20transfer,
		ICS20Bank:     *ics20bank,
	}
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

func makeGenTxOpts(chainID *big.Int, prv *ecdsa.PrivateKey) func(ctx context.Context) *bind.TransactOpts {
	signer := gethtypes.LatestSignerForChainID(chainID)
	addr := gethcrypto.PubkeyToAddress(prv.PublicKey)
	return func(ctx context.Context) *bind.TransactOpts {
		return &bind.TransactOpts{
			From:     addr,
			GasLimit: 6382056,
			Signer: func(address common.Address, tx *gethtypes.Transaction) (*gethtypes.Transaction, error) {
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

func InitializeChain(rpcAddress string, ethChainID int64, mnemonic string, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress string) (*Chain, error) {
	ethClient, err := client.NewETHClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	chain := NewChain(ethClient, ethChainID, mnemonic, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress)

	return chain, nil
}
