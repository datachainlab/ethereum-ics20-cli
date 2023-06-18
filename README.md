# ethereum-ics20-cli

## Install
```bash
$ git clone https://github.com/datachainlab/ethereum-ics20-cli.git
$ cd ethereum-ics20-cli
$ go install
CLI="ethereum-ics20-cli"
```

## erc20
```bash
${CLI} erc20 balance --rpc-address=${rpc_address} --wallet-address=${wallet_address} --denom=${denom}
${CLI} erc20 transfer --rpc-address=${RPC_ADDRESS} --mnemonic="${MNEMONIC}" --from-index=${from_index} --to-address=${to_address} --amount=${amount} --denom=${DENOM}
```

## ics20
```bash
${CLI} ics20 balance --rpc-address=${rpc_address} --ics20-bank-address=${ICS20_BANK_ADDRESS} --wallet-address=${wallet_address} --denom=${denom}
${CLI} ics20 transfer --rpc-address=${from_rpc_address} --mnemonic="${MNEMONIC}" --ics20-bank-address=${ICS20_BANK_ADDRESS} --ics20-transfer-bank-address=${ICS20_TRANSFER_BANK_ADDRESS} --from-index=${from_index} --to-address=${to_address} --amount=${amount} --denom=${denom} --port-id=${PORT_ID} --channel-id=${CHANNEL_ID} --timeout-height=${timeout_height}
```
