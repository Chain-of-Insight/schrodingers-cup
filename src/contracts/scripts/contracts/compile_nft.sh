#!/bin/bash

echo "Compiling NFT contract..."
ligo compile-contract src/nft.ligo main > build/nft.tz

echo "Compiling NFT contract's initial storage..."
ligo compile-storage src/nft.ligo main 'record [
  contractOwner = ("tz1codeYURj5z49HKX9zmLHms2vJN2qDjrtt" : address);
  contractOracle = ("tz1codeYURj5z49HKX9zmLHms2vJN2qDjrtt" : address);
  nfts = (map [] : map(nat, record [ owner : address; data : bytes ]))
]' > build/nft.init.tz
