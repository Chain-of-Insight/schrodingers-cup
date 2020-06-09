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

module.exports = {
    testNomic: testNomic
};