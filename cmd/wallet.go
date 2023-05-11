package cmd

import (
	"fmt"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
	"github.com/spf13/cobra"
)

func walletCmd() *cobra.Command {
	var mnemonic string
	var walletIndex uint64
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "wallet commands",
	}
	addressCmd := &cobra.Command{
		Use:   "address",
		Short: "Get address by specifying an index of the wallet",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			address, err := address(mnemonic, walletIndex)
			if err != nil {
				return err
			}
			fmt.Printf("%s", address)
			return nil
		},
	}
	addressCmd.Flags().StringVar(&mnemonic, "mnemonic", "", "mnemonic phrase")
	addressCmd.Flags().Uint64Var(&walletIndex, "wallet-index", 0, "index of the wallet")

	cmd.AddCommand(addressCmd)
	return cmd
}

func address(mnemonic string, walletIndex uint64) (string, error) {
	key, err := wallet.GetPrvKeyFromMnemonicAndHDWPath(mnemonic, fmt.Sprintf("m/44'/60'/0'/0/%d", walletIndex))
	if err != nil {
		return "", err
	}
	addr := gethcrypto.PubkeyToAddress(key.PublicKey)
	return addr.Hex(), nil
}
