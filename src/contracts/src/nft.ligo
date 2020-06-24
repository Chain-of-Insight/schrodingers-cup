type trusted is address
type nftId is nat

type nft is
  record [
    owner : address;
    data  : bytes
  ]

type nfts is map( nftId, nft )

(* Contract storage *)
type storage is
  record [
    contractOwner  : trusted;
    contractOracle : trusted;
    nfts           : nfts
  ]

type transfer_params is
  record [
    nftId        : nftId;
    destination  : address
  ]

(* define return for readability *)
type return is list (operation) * storage

(* define noop for readability *)
const noOperations : list (operation) = nil;

(* Valid entry points *)
type entry_action is
  | SetOracle of trusted
  | Mint of nft
  | Transfer of transfer_params
  | Burn of nftId

(* Set the contract oracle *)
function set_oracle(const addr : trusted; var s : storage) : return is
  block {
    (* Can only be called by contract owner *)
    if Tezos.source =/= s.contractOwner then
      failwith("NOPERM")
    else skip;

    (* Update storage *)
    s.contractOracle := addr;

  } with (noOperations, s)

(* Mint NFT *)
function mint(const input : nft; var s : storage) : return is
  block { 
    (* check permission *)
    if Tezos.sender =/= s.contractOwner and Tezos.sender =/= s.contractOracle then
      failwith("NOPERM")
    else skip;

    (* Get next ID *)
    const nextId : nat = Map.size(s.nfts) + 1n;

    (* Create NFT for next ID *)
    s.nfts[nextId] := input;

  } with (noOperations, s)

(* Transfer NFT *)
function transfer(const input : transfer_params; var s : storage) : return is
  block {
    (* Retrieve NFT from storage or fail *)
    const nft : nft =
      case s.nfts[input.nftId] of
        Some (instance) -> instance
      | None -> (failwith ("Unknown NFT") : nft)
      end;

    (* Check permissions *)
    if Tezos.source =/= nft.owner then
      failwith("NOPERM")
    else skip;

    (* change owner's address *)
    nft.owner := input.destination;

    (* Update storage *)
    s.nfts[input.nftId] := nft;

  } with (noOperations, s)

(* Burn an NFT *)
function burn(const input : nftId; var s : storage) : return is
  block {
    (* Retrieve NFT from storage or fail *)
    const nft : nft =
      case s.nfts[input] of
        Some (instance) -> instance
      | None -> (failwith ("Unknown NFT") : nft)
      end;

    (* Check permissions *)
    if Tezos.source =/= nft.owner then
      failwith("NOPERM")
    else skip;

    (* remove NFT *)
    remove input from map s.nfts;

  } with (noOperations, s)

(* Main entrypoint *)
function main (const action : entry_action; var s : storage) : return is
  block {
    skip
  } with case action of
    | SetOracle(param) -> set_oracle(param, s)
    | Mint(param) -> mint(param, s)
    | Transfer(param) -> transfer(param, s)
    | Burn(param) -> burn(param, s)
  end;
