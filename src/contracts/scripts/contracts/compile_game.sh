#!/bin/bash

echo "Compiling Game contract..."
ligo compile-contract src/game.ligo main > build/game.tz

echo "Compiling Game contract's initial storage..."
ligo compile-storage src/game.ligo main 'record [
  gameOwner = ("tz1codeYURj5z49HKX9zmLHms2vJN2qDjrtt" : address);
  gameOracle = ("tz1codeYURj5z49HKX9zmLHms2vJN2qDjrtt" : address);
  gameNft = ("tz1codeYURj5z49HKX9zmLHms2vJN2qDjrtt" : address);
  gameVars = record [
    startPoints = 0;
    winPoints = 100;
    rulePassPoints = 10;
    voteAgainstPoints = 10;
    ruleFailedPenalty = 10;
    gameWindowDuration = 7200;
    gameWindowHourUTC = 4;
    turnWindowDuration = 300;
    voteRatioRequired = 100;
    transmutationVoteRatioRequired = 100
  ];
  gameScore = (map [] : map(address, int));
  gameWinner = ("tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb" : address)
]' > build/game.init.tz
