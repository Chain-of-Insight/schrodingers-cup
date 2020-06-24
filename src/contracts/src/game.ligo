type trusted is address
type player is address
type points is int
type playerPoints is map (player, points)

type nomicVars is
  record [
    startPoints                     : points;
    winPoints                       : points;
    rulePassPoints                  : points;
    voteAgainstPoints               : points;
    ruleFailedPenalty               : points;
    gameWindowDuration              : int;
    gameWindowHourUTC               : int;
    turnWindowDuration              : int;
    voteRatioRequired               : int;
    transmutationVoteRatioRequired  : int
  ]

(* contract storage *)
type storage is
  record [
    gameOwner     : trusted;
    gameOracle    : trusted;
    gameNft       : trusted;
    gameVars      : nomicVars;
    gameScore     : playerPoints;
    gameWinner    : player
  ]

(* mint nft params *)
type nft is
  record [
    owner         : player;
    data          : bytes
  ]

(* define return for readability *)
type return is list (operation) * storage

(* define noop for readability *)
const noOperations : list (operation) = nil;

(* Valid entry points *)
type entry_action is
  | SetOracle of trusted
  | SetNft of trusted
  | SetVars of nomicVars
  | UpdateScore of playerPoints

(* Set the game vars *)
function set_vars(const input : nomicVars; var s : storage) : return is
  block {
    (* Can only be called by game oracle *)
    if Tezos.source =/= s.gameOracle then
      failwith("NOPERM")
    else skip;

    (* Update game vars *)
    s.gameVars := input;

  } with (noOperations, s)

(* Set the game oracle *)
function set_oracle(const addr : trusted; var s : storage) : return is
  block {
    (* Can only be called by game owner *)
    if Tezos.source =/= s.gameOwner then
      failwith("NOPERM")
    else skip;

    (* Update storage *)
    s.gameOracle := addr;

  } with (noOperations, s)

(* Set the game nft contract *)
function set_nft(const addr : trusted; var s : storage) : return is
  block {
    (* Can only be called by game owner *)
    if Tezos.source =/= s.gameOwner then
      failwith("NOPERM")
    else skip;

    (* Update storage *)
    s.gameNft := addr;

  } with (noOperations, s)

(* Update game scores *)
function update_score(const input : playerPoints; var s : storage) : return is
  block {
    (* Can only be called by game oracle *)
    if Tezos.source =/= s.gameOracle then
      failwith("NOPERM")
    else skip;

    (* Update game score in storage *)
    s.gameScore := input;

    (* initialize operations *)
    var operations : list (operation) := nil;

    (* check if anyone won *)
    var winner : bool := False;
    for p -> score in map s.gameScore block {
      if score >= s.gameVars.winPoints and winner = False then {
        winner := True;
        s.gameWinner := p;

        (* @todo send any tez in contract to winner? *)

        (* Send NFT to winner *)
        const nft_params : nft =
          record [
            owner = p;
            data = Bytes.pack("NOMSU Winner")
          ];
        const nft_entrypoint : contract (nft) = get_entrypoint("%mint", s.gameNft);
        const nftOperation : operation = transaction (nft_params, 0tz, nft_entrypoint);
        operations := nftOperation # operations;
      } else skip;
    }

  } with (operations, s)

(* Main entrypoint *)
function main (const action : entry_action; var s : storage) : return is
  block {
    skip
  } with case action of
    | SetOracle(param) -> set_oracle(param, s)
    | SetNft(param) -> set_nft(param, s)
    | SetVars(param) -> set_vars(param, s)
    | UpdateScore(param) -> update_score(param, s)
  end;
