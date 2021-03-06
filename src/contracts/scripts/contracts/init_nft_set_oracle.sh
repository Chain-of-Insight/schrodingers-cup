#!/bin/bash

export network="$1"
export contract="$2"
export owner="$3"
export oracle="$4"

if [ "$network" == "" ] || [ "$contract" == "" ] || [ "$owner" == "" ] || [ "$oracle" == "" ]; then
  echo "Usage: $0 <network> <nft-contract> <contract-owner> <oracle-contract>"
  exit 0
fi

if [[ ! " testnet mainnet " =~ " $network " ]]; then
  echo "Must specify network (mainnet|testnet)"
  exit 1
fi

export tezos_client_connection=""
if [ "$network" == "mainnet" ]; then
  tezos_client_connection="-A mainnet-tezos.giganode.io -S -P 443"
fi

if [ "$network" == "testnet" ]; then
  tezos_client_connection="-A testnet-tezos.giganode.io -S -P 443"
fi

params=$(ligo compile-parameter src/nft.ligo main "SetOracle((\"$oracle\": address))")
cmd="tezos-client $tezos_client_connection transfer 0 from $owner to $contract --arg '$params' --dry-run"

echo "Run the dry-run command to get the burn-cap; then re-run without --dry-run:"
echo ""
echo "  $cmd"
echo ""
