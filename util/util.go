package util

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// ref https://github.com/hyperledger-labs/yui-ibc-solidity/blob/v0.3.21/contracts/apps/20-transfer/ICS20Lib.sol#L210
func ToCanonicalICS20Denom(denom string) (string, error) {
	if common.IsHexAddress(denom) {
		return strings.ToLower(denom), nil
	}

	if strings.Contains(denom, "/") {
		identifiers := strings.Split(denom, "/")
		// If a prefix exists, at least one port-id and one channel-id must be specified
		// ex) ${PortID}/${ChannelID}/${Denom}
		if len(identifiers) < 3 {
			return "", fmt.Errorf("invalid denom format: %s", denom)
		}

		return strings.Join(
			append(identifiers[0:len(identifiers)-1], strings.ToLower(identifiers[len(identifiers)-1])),
			"/"), nil
	}

	return "", fmt.Errorf("invalid denom format: %s", denom)
}
