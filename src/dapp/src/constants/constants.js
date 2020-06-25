export const ruleTypes = {
  MUTABLE: 'mutable',
  IMMUTABLE: 'immutable'
}

export const ruleSetTypes = {
  SAVED: 'saved',
  IMMUTABLE: 'immutable',
  MUTABLE: 'mutable',
  QUEUED: 'queued'
}

export const proposalTypes = {
  CREATE: 'create',
  UPDATE: 'update',
  TRANSMUTE: 'transmute',
  DELETE: 'delete',
}

export const voteTypes = {
  YES: 1,
  NO: 0,
  ABSTAIN: -1
}

export const TZ_WALLET_PATTERN = "(tz(?:1|2|3)[a-zA-Z0-9]{33})";

// XXX: removed regex end-of-string token ('$') on these for now while testing 
export const eventMessages = {
  RULE_PROPOSED: `^${TZ_WALLET_PATTERN} proposed a rule in round (\\d+)`,
  VOTE_CAST: `^${TZ_WALLET_PATTERN} successfully voted (YES|NO) in round (\\d+)`,
  ROUND_OVER: `^Round (\\d+) has concluded`
}

// export const CURRENT_RULES = require('../rules/currentRules.json');
