package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/datachainlab/ethereum-ics20-cli/chains/geth"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
	"github.com/spf13/cobra"
)

const (
	WalletUsage = "Usage: address <configDir> <chainIndex> <walletIndex>"
)

func walletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: WalletUsage,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				log.Fatalln(WalletUsage)
			}
			configDir := args[0]
			chainIndex, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			walletIndex, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			address, err := address(configDir, chainIndex, walletIndex)
			if err != nil {
				log.Fatalln("address Error: ", err)
			}
			fmt.Printf("%s", address)
		},
	}
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
