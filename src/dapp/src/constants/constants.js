export const ruleTypes = {
  MUTABLE: 'mutable',
  IMMUTABLE: 'immutable'
}

export const ruleSetTypes = {
  SAVED: 'SAVED',
  CURRENT: 'CURRENT',
  QUEUED: 'QUEUED'
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

export const CURRENT_RULES = require('../rules/currentRules.json');
