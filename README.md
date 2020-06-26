# Schrodinger's Cup
A game of Peter Suber's Nomic running on the Tezos network

<br/>
<p align="center">
  <img width="450px" height="auto" src="https://cdn.discordapp.com/attachments/709953915550171206/717048880000270431/schrodingers-cup.png">
</p>
<br/>

## The legend of "Tezos-Nomic"

Tezos supports meta upgrades: the protocols can evolve by amending their own code. This is not unlike philosopher Peter Suber’s Nomic, a game built around a fully introspective set of rules that are changed and develop as the game progresses. Actually, the idea for Tezos was based on Nomic.

## How it works

Schrodinger's Cup consists of a website, API, and Tezos smart contracts. Play sessions are conducted daily, the duration and timing of which is coded into the orignal rule set stored on Tezos. However, just like any other rule of Nomic they can be updated during game play. 

## NFTs and Esports

Players who successfully perform an in-game action (changing the rules, voting, etc.) receive a Tezos NFT for playing. Additionally, players with the top highest scores are ranked by the NFT distribution system. These ranking NFTs contain a zero knowledge proof attribute which can be used to calculate and prove their final ranking on the overall Nomic ELO ladder. Proofs for ELO tokens can be updated later as the player participates in future games of Nomic, providing the basis for the first long term ELO / Glicko-2 style ladder ranking built on blockchain technology.

The overall winner of the Nomic game receives an XTZ reward transferred to their wallet in the same operation as mints their ranking NFT.

## Why
We hope to produce a blockchain game that's entertaining but psychologically intense, and that educates players about how Tezos actually works. Additionally, our NFT and rewards distribution system presents a reasonable implementation to show how esports tournament ladders can be run on blockchain technology.

## Links
- Dapp (Player interface): https://nomic.schrodingerscup.com/
- API: https://api.schrodingerscup.com/
- Game Contract (Successfully completed game of Nomic): https://better-call.dev/carthagenet/KT1UBKq1mMym1vHwgk3JbcLqmRhu3ActcLYB/operations
- NFT Rewards Contract: https://better-call.dev/carthagenet/KT1QkCAYZT4Eo94RML6yXzXFZMXx6gSFgAg5/operations


---

# Gameplay Gallery

## Pre-game Warmup

### After connecting with your wallet you'll see a countdown to the next game. Play sessions are held daily at 16:00 UTC until a winner is declared in the contract.

<br/>
<p align="center">
  <img src="https://i.imgur.com/rfCZ3gh.png">
</p>
<br/>

### In between game play sessions, you'll want to brush up on your skills. The game's meta-programming language, for changing the rules, is called [Nomsu](https://nomsu.org/). Head to https://nomic.schrodingerscup.com/practice where you can test compiling rules and save or queue them for use in the next game session.

<br/>
<p align="center">
  <img src="https://i.imgur.com/hHpEHMu.png">
</p>
<br/>

## The Game Knows Who You Are 😱

### Before joining live play, users sign a message proving ownership of their Tezos address. The API server sends some of the game's transactions keeping it a "free-to-play" game, but this also ensures nobody else can claim your wallet address as a display name in the game chat.

<br/>
<p align="center">
  <img src="https://i.imgur.com/WBS6dE4.png">
</p>
<br/>

<br/>
<p align="center">
  <img src="https://i.imgur.com/Y65TKb3.png">
</p>
<br/>

## Live Play

### In an active game, the order of turns is decided by the order players joined the game session. Once a rule is proposed by the "first joiner", a voting interface is presented to each other player appearing above the game chat.

<br/>
<p align="center">
  <img src="https://i.imgur.com/yaSSkPe.png">
</p>
<br/>

### Clicking the "Show rule code" button displays information you'll want to check out. Here we see a side-by-side if a previous rule is being changed, or just the new code if it's a rule being added. If the proposal is to delete or transmute a previous rule, only that rule's current code is shown.

<br/>
<p align="center">
  <img src="https://i.imgur.com/XDGXCun.png">
</p>
<br/>

### As defined in [Mutable Rule 0](https://github.com/Chain-of-Insight/schrodingers-cup/blob/master/src/api/nomsu/rules/mutable/rule0.nom) (at the outset of the game at least) everyone will need to vote before triggering a decision. While you're waiting on those slow-pokes you can watch voting results in real time 😎

<br/>
<p align="center">
  <img src="https://i.imgur.com/C4R7t3e.png">
</p>
<br/>

### Once there is a quorum of votes, the API server updates the game chat with results of the round, and updates Tezos storage on the game contract. 

<br/>
<p align="center">
  <img src="https://i.imgur.com/00w5oXC.png">
</p>
<br/>

### Once there is a rule change occurs that triggers a game over event, the API initiates an operation on the game contract which transfers an NFT reward to the winning player. [Example operation](https://better-call.dev/search?text=oocc7Sz4nyaCcn9ucz5NEru9dop7Ugsyc6LN4PEhghptfoFEhiH)

<br/>
<p align="center">
  <img src="https://i.imgur.com/XmNiZRA.png">
</p>
<br/>
