#!/bin/bash

mkdir -p ./build
./scripts/contracts/compile_nft.sh
./scripts/contracts/compile_game.sh
