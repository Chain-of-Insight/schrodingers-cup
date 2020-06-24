#!/bin/bash

export network="$1"
export cname="$2"
export owner="$3"

if [ "$network" == "" ] || [ "$cname" == "" ] || [ "$owner" == "" ] ; then
  echo "Usage: $0 <network> <contract-name> <contract-owner>"
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

cmd="tezos-client $tezos_client_connection originate contract $cname transferring 0 from $owner running build/game.tz --init '$(< build/game.init.tz)' --dry-run"

echo "Run the dry-run command to get the burn-cap to originate the contract; then re-run without --dry-run:"
echo ""
echo "  $cmd"
echo ""
