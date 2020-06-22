const axios = require('axios').default;

const API_URL = process.env.API_BASE_URL;

/**
 * Test compile a nomsu payload
 * @param {String} code : A nomsu payload to send to the API for test execution
 */
async function testNomic(code) {
  if (typeof code !== 'string') {
    return 'Payload rejected. String required, got ' + typeof code;
  }

  let apiEndpoint = API_URL + 'test';
  const res = await axios.post(apiEndpoint, {code: code});

  return res;
};

/**
 * Authorize user in nomsu backend
 */
async function auth(msg, sig, pubKey, address) {
  if (typeof msg !== 'string') {
    return 'msg rejected. String required, got ' + typeof msg;
  }

  if (typeof sig !== 'string') {
    return 'sig rejected. String required, got ' + typeof sig;
  }

  if (typeof pubKey !== 'string') {
    return 'pubKey rejected. String required, got ' + typeof pubKey;
  }

  let apiEndpoint = API_URL + 'auth';
  const res = await axios.post(apiEndpoint, {msg: msg, sig: sig, pubKey: pubKey, address: address});

  return res;
};

/**
 * Propose a nomsu rule
 */
async function proposeRule(jwt, code, index, kind, type) {
  if (typeof(jwt) !== 'string') {
    return 'jwt rejected. String required, got ' + typeof code;
  }
  if (typeof code !== 'string') {
    return 'code rejected. String required, got ' + typeof code;
  }

  if (typeof index !== 'number') {
    return 'index rejected. Number required, got ' + typeof index;
  }

  if (typeof code !== 'string') {
    return 'kind rejected. String required, got ' + typeof kind;
  }

  if (typeof code !== 'string') {
    return 'kind rejected. String required, got ' + typeof type;
  }

  const config = {
    headers: {
      'Authorization': 'Bearer ' + jwt
    }
  }

  let apiEndpoint = API_URL + 'game/propose';
  const res = await axios.post(apiEndpoint, {
    code: code,
    index: index,
    kind: kind,
    type: type
  }, config);

  return res;
};

module.exports = {
  testNomic: testNomic,
  PerformAuth: auth,
  proposeRule: proposeRule
};
