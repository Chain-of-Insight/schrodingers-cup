<template>
  <div>
    <!-- Notifications -->
    <Notification 
      :type="alert.type" 
      :msg="alert.msg" 
      v-on:reset="alert = {type: null, msg: null}"
    ></Notification>

    <div class="container main">
      <h1>{{ title }}</h1>
      <h5 class="mb-4">{{ subtitle }}</h5>

      <!-- Not Connected -->
      <template v-if="!connected">
        <ul class="list-unstyled">
          <li @click="connectUser()">
            <button class="btn btn-primary btn-connect">Login With Tezos</button>
          </li>
        </ul>
        <p>Connect your Tezos wallet to get started</p>
      </template>

      <!-- Connected -->
      <template v-else>
        <!-- Round totals -->
        <section>
          <Totals v-bind:round="currentRound" v-bind="currentTotals"></Totals>
        </section>

        <!-- Player Chat -->
        <section>
          <div id="messages" class="message-container">
            <!-- Chat Messages -->
            <div v-for="(message, index) in chatMessages" v-bind:key="index">
              <p v-bind:class="[message.type, (message.author) ? message.author : 'system', 'chat-msg']">
                <span v-if="message.author" class="chat-author">{{ message.author }}:</span>
                <span v-bind:class="['chat-msg-body', message.type]" v-if="message.msg">{{ message.msg }}</span>
              </p>
            </div>
          </div>

          <!-- Chat Form Input -->
          <div class="input-group chat-controls">
            <input 
              id="chat-input" 
              type="text" 
              class="form-control"
              aria-label="Send a chat message..."
              aria-describedby="Nomic player chat"
              placeholder="Send a chat message..." 
              v-model="chatInput.value"
              @focus="chatInput.focused = true"
              @blur="chatInput.focused = false"
              v-on:keyup="chatKeyListener($event)"
            />
            <!-- Send Message -->
            <div class="input-group-append">
              <button class="btn btn-outline-secondary" type="button" @click="submitChatMessage()">Send</button>
            </div>
          </div>

          <!-- Send test chat messages with fake API wallet address -->
          <div class="btn-group">
            <button class="btn btn-outline-secondary" type="button" @click="testSystemMessage(0)">Your Turn</button>
            <button class="btn btn-outline-secondary" type="button" @click="testSystemMessage(1)">Another Player's Turn</button>
            <button class="btn btn-outline-secondary" type="button" @click="testSystemMessage(2)">Create Rule</button>
            <button class="btn btn-outline-secondary" type="button" @click="testSystemMessage(3)">Update Rule</button>
            <button class="btn btn-outline-secondary" type="button" @click="testSystemMessage(4)">Transmute Rule</button>
            <button class="btn btn-outline-secondary" type="button" @click="testSystemMessage(5)">Delete Rule</button>
          </div>
        </section>

        <!-- IDE -->
        <div class="editor-toggle">
          <!-- IDE Shown -->
          <button class="btn btn-inverse toggle-rules-editor" @click="toggleEditor()">{{ showEditor ? "Hide Rules Editor" : "Show Rules Editor" }}</button>
          
          <!-- For testing: -->
          <div class="btn-group" role="group" aria-label="">
            <!-- Test Voting -->
            <!-- <button type="button" class="btn btn-outline-primary" @click="votingHandler()">
              Test Voting
            </button> -->
            <!-- Test Rule Proposal -->
            <button type="button" class="btn btn-outline-primary" @click="ruleProposalHandler()">
              Test Rule Proposal
            </button>
          </div>
          <Voting
            v-bind:voting-duration="votingDuration"
            v-on:vote-cast="castVote"
            ref="voting"
            v-bind:voting-candidate="votingCandidate"
          ></Voting>
          <RuleProposal ref="proposal"></RuleProposal>
          <Practice :activeGame="true" v-if="showEditor"></Practice>
        </div>
      </template>
    </div>
  </div>
</template>


<script>
import { 
  Tezos,
  mountProvider,
  getBalance,
  signMessage,
  validateAddress
} from '../../services/tezProvider';

import { 
  getToken,
  Twilio 
} from '../../services/twilioProvider';

// API
import { PerformAuth } from '../../services/apiProvider';

// Child components
import Notification from '../common/Notifications.vue';
import RuleProposal from '../common/RuleProposal.vue';
import Voting from '../common/Voting.vue';
import Totals from '../common/Totals.vue';
// IDE Component
import Practice from '../practice/Practice.vue';

const voteType = {
  YES: 0,
  NO: 1,
  ABSTAIN: -1
}
if (Object.freeze) Object.freeze(voteType);

const TZ_WALLET_PATTERN = "tz(1|2|3)[a-zA-Z0-9]{33}";

export default {
  components: {
    Notification,
    Practice,
    RuleProposal,
    Voting,
    Totals
  },
  data: () => ({
    title: "Nomic Battlegrounds",
    subtitle: "Pwned or be pwned, the choice is yours",
    alert: {
      type: null,
      msg: null
    },
    network: (process.env.hasOwnProperty('CURRENT_NETWORK')) ? process.env.CURRENT_NETWORK : 'carthagenet',
    address: null,
    getBalance: getBalance,
    connected: false,
    Tezos: Tezos,
    getToken: getToken,
    TwilioToken: null,
    TwilioIdentity: null,
    TwilioChat: Twilio,
    chatClient: null,
    chatMessages: [],
    chatChannelJoined: false,
    chatInput: {
      focused: false,
      value: null
    }, 
    chatChannel: null,
    mountProvider: mountProvider,
    signMessage: signMessage,
    loginSigned: null,
    jwtToken: null,
    showEditor: false,
    votingDuration: 300, // from rule 'external: $bl_turnWindowDuration = 300'
    votingCandidate: null,
    currentRound: 1,
    currentTotals: {
      yes: 0,
      no: 0,
      abstain: 0
    },
    apiWallet: 'BECjcQFZnoVYu94ns5HMLW7yaDsoJZYbU1zt',
  }),
  computed: {
    msgPatterns: function () {
      return {
        NEW_TURN_PATTERN: `^It's ${TZ_WALLET_PATTERN}'s turn to propose a rule change$`,
        CREATE_RULE_PATTERN: `^${TZ_WALLET_PATTERN} has proposed a new rule$`,
        UPDATE_RULE_PATTERN: `^${TZ_WALLET_PATTERN} has proposed an update to rule \\d+$`,
        TRANSMUTE_RULE_PATTERN: `^${TZ_WALLET_PATTERN} has proposed to transmute rule \\d+$`,
        DELETE_RULE_PATTERN: `^${TZ_WALLET_PATTERN} has proposed to delete rule \\d+$`,
        YES_VOTE_PATTERN: `${TZ_WALLET_PATTERN} voted YES in round ${this.currentRound}$`,
        NO_VOTE_PATTERN: `${TZ_WALLET_PATTERN} voted NO in round ${this.currentRound}$`,
        ABSTAIN_VOTE_PATTERN: `${TZ_WALLET_PATTERN} abstained in round ${this.currentRound}$`
      }
    }
  },
  mounted: async function () {
    await this.mountProvider();
    console.log('Nomic mounted', this.network);
    let returningUser = sessionStorage.getItem('tzAddress');
    if (returningUser) {
      this.connected = true;
      this.address = returningUser;
      // Request Twilio token
      try {
        // Request auth
        let tokenResponse = await this.getToken(this.address);
        if (tokenResponse) {
          // Valid auth
          this.TwilioToken = tokenResponse.token
          this.TwilioIdentity = tokenResponse.identity;
          console.log('Twilio Chat Params =>', [this.TwilioToken, this.TwilioIdentity]);
          console.log('TwilioChat', this.TwilioChat);
          // Connect to chat room
          this.connectChat();
          // await this.doLoginMessageSigning();
        }
      } catch (e) {
        // Auth failed
        console.log('Error requesting Twilio Auth Token', e);
      }
    }
  },
  methods: {
    connectUser: async function () {
      // Connect as required
      if (this.connected)
        return;

      // Signer instance
      this.address = await tezbridge.request({method: 'get_source'});
      
      // Fetch balance / Connection callbacks
      if (typeof this.address == 'string') {
        if (this.address.length === 36) {
          console.log('User XTZ Address =>', this.address);
          sessionStorage.setItem('tzAddress', this.address);
          this.connected = true;
          let balance = await this.getBalance(this.address);
          console.log("User balance =>", balance);
          // Request Twilio token
          try {
            // Request authvotingHandler()
            let tokenResponse = await this.getToken(this.address);
            if (tokenResponse) {
              // Valid auth
              this.TwilioToken = tokenResponse.token
              this.TwilioIdentity = tokenResponse.identity;
              console.log('Twilio Chat Params =>', [this.TwilioToken, this.TwilioIdentity]);
              console.log('TwilioChat', this.TwilioChat);
              // Connect to chat room
              this.connectChat();

              // Sign login auth message for API
              // await this.doLoginMessageSigning();
            }
          } catch (e) {
            // Auth failed
            console.log('Error requesting Twilio Auth Token', e);
          }
        }
      }
    },
    doLoginMessageSigning: async function () {
      let timestamp = new Date().getTime();
      let signedMsg = await this.signMessage(timestamp);
      this.loginSigned = signedMsg
      console.log('Signed message result', this.loginSigned);

      // get pubkey
      pubKey = await Tezos.rpc.getManagerKey(this.address);
      console.log('Pubkey', pubKey);

      // address not revealed
      if (pubKey == this.address) {
          //@todo what do here?
          console.warn('Warning: Address not revealed. Message signing will fail.');
      }

      // auth and get JWT token
      let result;
      try {
        result = await PerformAuth(this.loginSigned.bytes, this.loginSigned.prefixSig, pubKey, this.address);
      } catch(e) {
        // unauthorized if we are here
        result = e.response;
      }

      if (result.data && result.data.token) {
        this.jwtToken = result.data.token;
      }
      console.log("JWT token", this.jwtToken);
    },
    // Player chat connector
    connectChat: function () {
      if (!this.TwilioChat || !this.TwilioToken) {
        return;
      }
      
      const Chat = this.TwilioChat;
      
      Chat.Client.create(this.TwilioToken).then(client => {
        console.log('Chat client created...');
        this.chatClient = client;
        this.chatClient.getSubscribedChannels().then(this.createOrJoinGeneralChat);
      });
    },
    // Joins user in chat room
    createOrJoinGeneralChat: function() {
      if (!this.chatClient) {
        return;
      }

      const chatRoom = 'nomic-chat';

      // Joining chat
      let joinMessage = {
          type: 'info',
          msg: 'Attempting to join chat...'
      };

      let joinedMessage = {
          type: 'info',
          msg: 'Joined Nomic Player Chat'
      };

      this.chatMessages.push(joinMessage);

      this.chatClient.getChannelByUniqueName(chatRoom)
      .then(async (channel) => {
        console.log('Existing chat channel found', channel);
        // A pre-existing channel was found
        let chatChannel = channel;
        this.chatChannel = channel;

        // Join user in channel as required
        if (chatChannel.status !== "joined") {
          this.chatMessages.push(joinedMessage);
          if (!this.chatChannelJoined) {
            await this.chatChannel.join();
            console.log('Channel joined', channel);
          }
          this.chatChannelJoined = true;
          this.setupChatChannel();
        } else {
          this.chatMessages.push(joinedMessage);
          this.setupChatChannel();
          this.chatChannelJoined = true;
        }
      }).catch((error) => {
        console.log('Error joining chat', error);
        // No existing channel found, creating one instead
        this.chatClient.createChannel({
          uniqueName: chatRoom,
          friendlyName: 'Nomic Player Chat'
        }).then((channel) => {
          this.chatChannel = channel;
          console.log("Chat Channel =>", this.chatChannel);
          this.chatChannel.join().then((channel) => {
            // Release joined message
            this.chatMessages.push(joinMessage);
            this.setupChatChannel();
          });
        }).catch((channel) => {
          console.log('Chat channel could not be created:');
          console.log(channel);
        });
      });
    },
    // Subscribe to chat message websockets
    setupChatChannel: function () {
      console.log("Setting up chat channel", this.chatChannel);

      /**
       * Subscribe to incoming messages sent to the channel
       * @param {Object} message : A Twilio Message object container the {String} properties: `author` and `body`
       */
      this.chatChannel.on('messageAdded', (message) => {
        apiWalletPattern = new RegExp('^' + this.apiWallet + ':\\s'); // XXX: FOR TESTING ONLY
        isSystemMessage = apiWalletPattern.test(message.body); // XXX: FOR TESTING ONLY

        if (isSystemMessage) {
          // TODO: replace above condition with `message.author == this.apiWallet` when api wallet set up
          this.onSystemMessage(message);
        } else {
          this.onUserMessage(message);
        }
      });

      /**
       * Subcsribe to incoming chat participants
       * @param {Object} participant : A Twilio `participant` object
       */
      this.chatChannel.on('participantConnected', participant => {
        let messageOutput = {
          type: 'info',
          msg: participant.identity + ' has joined the chat'
        };
        // Release message to UI
        this.chatMessages.push(messageOutput);
      });

      /**
       * Subcsribe to chat participant disconnections
       * @param {Object} participant : A Twilio `participant` object
       */
      this.chatChannel.once('participantDisconnected', participant => {
        console.log(`Remote participant "${participant.identity}" has disconnected from the audio channel`);
        let messageOutput = {
          type: 'warn',
          msg: participant.identity + ' has disconnected from the chat'
        };
        // Release message to UI
        this.chatMessages.push(messageOutput);
      });
    },
    onSystemMessage: function (message) {
      const apiWalletPattern = new RegExp('^' + this.apiWallet + ':\\s'); // XXX: FOR TESTING ONLY
      let messageBody = message.body.replace(apiWalletPattern, ''); // XXX: FOR TESTING ONLY

      switch (messageBody) {
        case (messageBody.match(RegExp(this.msgPatterns.CREATE_RULE_PATTERN)) || {}).input:
        case (messageBody.match(RegExp(this.msgPatterns.UPDATE_RULE_PATTERN)) || {}).input:
        case (messageBody.match(RegExp(this.msgPatterns.TRANSMUTE_RULE_PATTERN)) || {}).input:
        case (messageBody.match(RegExp(this.msgPatterns.DELETE_RULE_PATTERN)) || {}).input:
          // On another player proposing a rule
          this.votingCandidate = {
            name: 'testRule',
            code: '$test_string = "this is test code"\nsay($test_string)'
          }
          this.$refs.voting.promptForVote(this.votingCandidate);
          break;
        case (messageBody.match(RegExp(this.msgPatterns.NEW_TURN_PATTERN)) || {}).input:
          // On new turn
          // Extract player's address from string
          const playerAddress = RegExp(TZ_WALLET_PATTERN).exec(messageBody)[0];
          // Check if it's logged in player
          if (playerAddress === this.TwilioIdentity) {
            // Format message
            messageBody = "It's your turn to propose a rule!";
            // Trigger proposal modal
            this.$refs.proposal.promptForProposal();
          }
          break;
        case (messageBody.match(RegExp(this.msgPatterns.YES_VOTE_PATTERN)) || {}).input:
          this.currentTotals.yes += 1;
          messageBody = messageBody.replace(this.TwilioIdentity, 'You');
          break;
        case (messageBody.match(RegExp(this.msgPatterns.NO_VOTE_PATTERN)) || {}).input:
          this.currentTotals.no += 1;
          messageBody = messageBody.replace(this.TwilioIdentity, 'You');
          break;
        case (messageBody.match(RegExp(this.msgPatterns.ABSTAIN_VOTE_PATTERN)) || {}).input:
          this.currentTotals.abstain += 1;
          messageBody = messageBody.replace(this.TwilioIdentity, 'You');
          break;
        default:
          return false;
      }

      // Parse message args
      let messageOutput = {
        type: 'system',
        msg: messageBody // XXX: replace this with message.body later
      };
      // Release message to UI
      this.chatMessages.push(messageOutput);
    },
    onUserMessage: function (message) {
        let messageType = message.author == this.TwilioIdentity ? 'local' : 'remote';
        // Parse message args
        let messageOutput = {
          type: messageType,
          author: message.author,
          msg: message.body
        };

        // Release message to UI
        this.chatMessages.push(messageOutput);
    },
    // Submit chat message
    submitChatMessage: function () {
      console.log("chatChannel from submitChatMessage =>", this.chatChannel);
      if (this.chatInput.value) {
        // Read from input
        let msgText = this.chatInput.value;
        // Push message to channel
        this.chatChannel.sendMessage(msgText);
        // Clear UI input
        this.chatInput.value = null;
        // console.log('Chat Messages', [this.chatMessages, this.chatChannel]);
      }
    },
    // Send messages when pressing 'enter' key
    chatKeyListener: function (e) {
      // console.log('keydown', e);
      if (e.keyCode == 13) {
        if (!this.chatChannel) {
          console.warn('Error submitting chat message, channel not setup', this.chatChannel);
          return;
        }
        // Send chat message to all recipients
        this.submitChatMessage();
      }
    },
    _retireNotification: function () {
      this.alert = {
        type: null,
        msg: null
      };
    },
    toggleEditor: function () {
      this.showEditor = this.showEditor ? false : true;
    },
    castVote: function (type) {
      if (!this.chatChannel)
        return false;

      let msgBody = null;

      switch (type) {
        case voteType.YES:
          msgBody = `${this.TwilioIdentity} voted YES in round ${this.currentRound}`;
          break;
        case voteType.NO:
          msgBody = `${this.TwilioIdentity} voted NO in round ${this.currentRound}`;
          break;
        case voteType.ABSTAIN:
          msgBody = `${this.TwilioIdentity} abstained in round ${this.currentRound}`;
          break;
        default:
          return;
      }

      const msgText = `${this.apiWallet}: ` + msgBody;
      this.chatChannel.sendMessage(msgText);
      this.votingCandidate = null;
    },
    ruleProposalHandler: function () {
      this.$refs.proposal.promptForProposal();
    },
    testSystemMessage: function (type) {
      if (!this.chatChannel)
        return false;

      let msgBody = null;

      switch (type) {
        case 0:
          // your turn
          msgBody = `It's ${this.TwilioIdentity}'s turn to propose a rule change`;
          break;
        case 1:
          // another player's turn
          msgBody = `It's tz1UbYZJosDay7WLMH5sn49uYVonZFQcjCEC's turn to propose a rule change`;
          break;
        case 2:
          // another player's turn
          msgBody = `tz1UbYZJosDay7WLMH5sn49uYVonZFQcjCEC has proposed a new rule`;
          break;
        case 3:
          // another player's turn
          msgBody = `tz1UbYZJosDay7WLMH5sn49uYVonZFQcjCEC has proposed an update to rule 12`;
          break;
        case 4:
          // another player's turn
          msgBody = `tz1UbYZJosDay7WLMH5sn49uYVonZFQcjCEC has proposed to transmute rule 12`;
          break;
        case 5:
          // another player's turn
          msgBody = `tz1UbYZJosDay7WLMH5sn49uYVonZFQcjCEC has proposed to delete rule 12`;
          break;
        default:
          return false;
      }

      const msgText = `${this.apiWallet}: ` + msgBody;
      this.chatChannel.sendMessage(msgText);
    }
  }
};
</script>

<style scoped>
  /* Chat */
  .input-group.chat-controls {
    width: 100%;
  }
  button.toggle-rules-editor {
    margin-top: 1rem;
    margin-bottom: 1rem;
    border: 1px solid #333333;
  }
</style>