package cmd

import (
	"fmt"

	"github.com/datachainlab/ethereum-ics20-cli/chains/geth"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
	"github.com/spf13/cobra"
)

func walletCmd() *cobra.Command {
	var configFile string
	var walletIndex uint64
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "wallet commands",
	}
	addressCmd := &cobra.Command{
		Use:   "address",
		Short: "Get address by specifying an index of the wallet",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			address, err := address(configFile, walletIndex)
			if err != nil {
				return err
			}
			fmt.Printf("%s", address)
			return nil
		},
	}
	addressCmd.Flags().StringVar(&configFile, "config", "", "config file path")
	addressCmd.Flags().Uint64Var(&walletIndex, "wallet-index", 0, "index of the wallet")

	cmd.AddCommand(addressCmd)
	return cmd
}

func address(configFile string, walletIndex uint64) (string, error) {
	chainConfig, err := geth.ParseChainConfig(configFile)
	if err != nil {
		return "", err
	}
	key, err := wallet.GetPrvKeyFromMnemonicAndHDWPath(chainConfig.Chain.HdwMnemonic, fmt.Sprintf("m/44'/60'/0'/0/%d", walletIndex))
	if err != nil {
		return "", err
	}
	addr := gethcrypto.PubkeyToAddress(key.PublicKey)
	return addr.Hex(), nil
}
