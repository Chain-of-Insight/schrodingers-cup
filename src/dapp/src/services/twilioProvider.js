const Twilio = require('twilio-chat');
const axios = require('axios').default;

const TOKEN_URL = process.env.TWILIO_SERVERLESS_ENDPOINT;

async function getToken(username) {
  const tokenUrl = `${TOKEN_URL}?username=${username}`;
  const res = await axios.get(tokenUrl);
  
  console.log('Token Response =>', res);
  
  if (res.data) {
    return res.data;
  } else {
    return false;
  }
};

module.exports = {
  getToken: getToken,
  Twilio: Twilio
}


/*
// The Twilio Lambda service contains this code:

exports.handler = function(context, event, callback) {  
  // make sure you enable ACCOUNT_SID and AUTH_TOKEN in Functions/Configuration
  const ACCOUNT_SID = context.ACCOUNT_SID;

  // you can set these values in Functions/Configuration or set them here
  const SERVICE_SID = 'SERVICE_SID';
  const API_KEY = 'API_KEY';
  const API_SECRET = 'API_SECRET';

  // REMINDER: This identity is only for prototyping purposes
  const IDENTITY = event.username || 'Default_Username';

  const AccessToken = Twilio.jwt.AccessToken;
  const IpMessagingGrant = AccessToken.IpMessagingGrant;

  const chatGrant = new IpMessagingGrant({
        serviceSid: SERVICE_SID
  });

  const accessToken = new AccessToken(
        ACCOUNT_SID,
        API_KEY,
        API_SECRET
  );

  accessToken.addGrant(chatGrant);
  accessToken.identity = IDENTITY;
 
  const response = new Twilio.Response();
  response.appendHeader('Access-Control-Allow-Origin', '*');
  response.appendHeader('Access-Control-Allow-Methods', 'GET');
  response.appendHeader('Access-Control-Allow-Headers', 'Content-Type');
  response.appendHeader('Content-Type', 'application/json');
  response.setBody({ token: accessToken.toJwt() });
  callback(null, response);
}
*/