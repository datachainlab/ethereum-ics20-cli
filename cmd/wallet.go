package cmd

import (
	"fmt"
	"strconv"

	"github.com/datachainlab/ethereum-ics20-cli/chains/geth"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
	"github.com/spf13/cobra"
)

func walletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "wallet commands",
	}
	addressCmd := &cobra.Command{
		Use:   "address",
		Short: "address of the wallet",
		Long:  "Usage: address <configDir> <chainIndex> <walletIndex>",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			configDir := args[0]
			chainIndex, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			walletIndex, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			address, err := address(configDir, chainIndex, walletIndex)
			if err != nil {
				return err
			}
			fmt.Printf("%s", address)
			return nil
		},
	}
	cmd.AddCommand(addressCmd)
	return cmd
}

func address(configDir string, chainIndex, walletIndex uint64) (string, error) {
	chainConfigs, err := geth.ParseChainConfigs(fmt.Sprintf("%s/%s", configDir, "chains"))
	if err != nil {
		return "", err
	}
	chain := chainConfigs[chainIndex]
	if chain == nil {
		return "", fmt.Errorf("chain not found")
	}
	key, err := wallet.GetPrvKeyFromMnemonicAndHDWPath(chain.Chain.HdwMnemonic, fmt.Sprintf("m/44'/60'/0'/0/%d", walletIndex))
	if err != nil {
		return "", err
	}
	addr := gethcrypto.PubkeyToAddress(key.PublicKey)
	return addr.Hex(), nil
}
