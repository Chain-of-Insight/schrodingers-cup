<template>
  <div class="h-100">
    <!-- Notifications -->
    <Notification 
      :type="alert.type" 
      :msg="alert.msg" 
      v-on:reset="alert = {type: null, msg: null}"
    ></Notification>

    <div class="container py-4">
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
          <Totals
            v-if="currentRound !== -1"
            v-bind:round="currentRound"
            v-bind="currentTotals"
            :turn-window="turnWindow"
            ref="totals"
          ></Totals>
        </section>

        <section>
          <Voting
            v-if="chatChannelJoined && votingCandidate"
            v-bind:turn-window="turnWindow"
            v-on:vote-cast="onVoteCast"
            ref="voting"
            :voting-candidate="votingCandidate"
          ></Voting>
        </section>

        <!-- Player Chat -->
        <section>
          <div ref="chatWindow" id="messages" class="message-container">
            <!-- Chat Messages -->
            <p
              v-for="(message, index) in chatMessages"
              v-bind:key="index"
              v-bind:class="[message.type, (message.author) ? message.author : 'system', 'chat-msg']"
            >
              <span v-if="message.author" class="chat-author">{{ message.author }}:</span>
              <span v-bind:class="['chat-msg-body', message.type]" v-if="message.msg">{{ message.msg }}</span>
            </p>
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
          <p class="h5 mt-3">Testing:</p>
          <button class="btn btn-primary" type="button" @click="ruleProposalHandler()">Proposal</button>
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
            <!-- <button type="button" class="btn btn-outline-primary" @click="ruleProposalHandler()">
              Test Rule Proposal
            </button> -->
          </div>
        </div>
        <RuleProposal
          v-if="chatChannelJoined"
          ref="proposal"
          v-on:rule-proposed="onRuleProposed"
          :turn-window="turnWindow"
        ></RuleProposal>
        <Practice
          :active-game="true"
          v-if="showEditor"
          style="height: 500px; max-height: 700px;"
        ></Practice>
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
import {
  PerformAuth,
  proposeRule,
  castVote,
  getRoundNumber,
  getPlayers,
  getProposedRule
} from '../../services/apiProvider';

// Child components
import Notification from '../common/Notifications.vue';
import RuleProposal from '../rule-proposal/RuleProposal.vue';
import Voting from '../common/Voting.vue';
import Totals from '../common/Totals.vue';
// IDE Component
import Practice from '../practice/Practice.vue';

const voteTypes = {
  YES: 1,
  NO: 0,
  ABSTAIN: -1
}

const ruleChangeTypes = {
  CREATE: 'create',
  UPDATE: 'update',
  TRANSMUTE: 'transmute',
  DELETE: 'delete',
}

const CURRENT_RULES = require('../practice/rules/currentRules.json');

const TZ_WALLET_PATTERN = "(tz(?:1|2|3)[a-zA-Z0-9]{33})";

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
    turnWindow: 300, // from rule 'external: $bl_turnWindowDuration = 300'
    votingCandidate: null,
    currentRound: -1,
    currentTotals: {
      yes: 0,
      no: 0,
      abstain: 0
    },
    players: [],
    currentTurn: null,
    nextTurn: null,
    ruleSets: {
      current: [],
      saved: []
    }
  }),
  computed: {
    msgPatterns: function () {
      return {
        NEW_TURN_PATTERN: `^It's ${TZ_WALLET_PATTERN}'s turn to propose a rule change$`,
        PROPOSAL_PATTERN: `^${TZ_WALLET_PATTERN} proposed a rule in round (${this.currentRound})$`,
        VOTE_PATTERN: `^${TZ_WALLET_PATTERN} successfully voted (YES|NO) in round (${this.currentRound})$`,
        // YES_VOTE_PATTERN: `^${TZ_WALLET_PATTERN} voted YES in round ${this.currentRound}$`,
        // NO_VOTE_PATTERN: `^${TZ_WALLET_PATTERN} voted NO in round ${this.currentRound}$`,
        // ABSTAIN_VOTE_PATTERN: `^${TZ_WALLET_PATTERN} abstained in round ${this.currentRound}$`
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

      // Get current round number
      await this.getCurrentRound();

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
          await this.doLoginMessageSigning();
          // Get players
          await this.getCurrentPlayers();
          // Get current rule set
          this.getCurrentRules();
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
              await this.doLoginMessageSigning();
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
        console.log('chat message!', message);

        if (message.author === 'system') {
          // TODO: replace above condition with `message.author == this.apiWallet` when api wallet set up
          this.onSystemMessage(message);
        } else {
          this.onUserMessage(message);
        }

        // auto-scroll to bottom to show new message
        this.$nextTick(function () {
          const messageWindow = this.$refs.chatWindow
          if (messageWindow && messageWindow.children.length) {
            messageWindow.scrollTop = messageWindow.scrollHeight - messageWindow.clientHeight;
          }
        })
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
    onSystemMessage: async function (message) {
      let messageBody = message.body;
      let playerAddress = null;

      switch (messageBody) {
        case (messageBody.match(RegExp(this.msgPatterns.PROPOSAL_PATTERN)) || {}).input:
          // TODO: GET (/game/proposals?) to get latest proposed rule
          playerAddress = RegExp(TZ_WALLET_PATTERN).exec(messageBody)[0];

          this.votingCandidate = null;
          await this.getLastProposed();

          // if (playerAddress !== this.TwilioIdentity) {
          //   // On another player proposing a rule
          //   this.votingCandidate = {
          //     name: 'testRule',
          //     code: '$test_string = "this is test code"\nsay($test_string)'
          //   }
          //   console.log('Time to vote!');
          //   this.$refs.voting.promptForVote(this.votingCandidate);
          // }
          break;
        case (messageBody.match(RegExp(this.msgPatterns.NEW_TURN_PATTERN)) || {}).input:
          // On new turn
          // Extract player's address from string
          playerAddress = RegExp(TZ_WALLET_PATTERN).exec(messageBody)[0];
          // Check if it's logged in player
          if (playerAddress === this.TwilioIdentity) {
            // Format message
            messageBody = "It's your turn to propose a rule!";
            // Trigger proposal modal
            this.$refs.proposal.promptForProposal();
          }
          break;
        case (messageBody.match(RegExp(this.msgPatterns.VOTE_PATTERN)) || {}).input:
          const matches = messageBody.match(RegExp(this.msgPatterns.VOTE_PATTERN));
          const vote = matches[2];
          const round = parseInt(matches[3]);

          // If vote msg applies to this round, increment vote counter
          if (round === this.currentRound) {
            switch (vote) {
              case "YES":
                this.currentTotals.yes ++;
                break;
              case "NO":
                this.currentTotals.no ++;
                break;
            }
          }
          
          // If msg applies to your vote:
          messageBody = messageBody.replace(this.TwilioIdentity, 'You');
          break;
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
    onVoteCast: async function (vote) {
      if (vote === voteTypes.ABSTAIN || typeof this.votingCandidate !== 'object') {
        // Don't send anything on abstain
        return false;
      }

      let result = null;
      try {
        result = await castVote(this.jwtToken, Boolean(vote), this.currentRound);
      } catch (error) {
        result = error.response;
      }

      if (result.status == 200) {
        if (!result.data) {
          console.error('Response successful but no data present:', result);
          return false;
        }

        if (result.data.success) {
          this.alert.type = 'success';
          this.alert.msg = 'Your vote was cast successfully';
          setTimeout(() => {
            this._retireNotification();
          }, 5000);

          // Update round number
          this.currentRound = result.data.round;
          // clear voting candidate and voting ribbon
          this.votingCandidate = null;
        } else {
          // Response OK but vote cast failed
          this.alert.type = 'danger';
          this.alert.msg = 'Vote cast unsuccessful: "' + result.data.message + '"... Please try again.';
          setTimeout(() => {
            this._retireNotification();
          }, 5000);
        }
      } else if (result.status == 500) {
        console.error('Error while trying to propose rule: ', result);
        this.alert.type = 'danger';
        this.alert.msg = 'There was an error while trying to cast your vote... Please try again.';
        setTimeout(() => {
          this._retireNotification();
        }, 5000);
      }
    },
    ruleProposalHandler: async function () {
      if (!this.jwtToken) {
        this.alert.type = 'danger';
        this.alert.msg = "It's your turn to propose a rule, but you havent' been authenticated yet! Have you signed a message in TezBridge yet?";
        setTimeout(async () => {
          this._retireNotification();
        }, 3000);
      } else {
        // Check if you already have a rule up for vote this round

        if (!this.votingCandidate)
          return false;

        if (
          this.votingCandidate.round === this.currentRound &&
          this.votingCandidate.author === this.TwilioIdentity
        ) {
          return false;
        }

        // Otherwise, prompt to propose a rule
        this.$refs.proposal.promptForProposal();
      }
    },
    onRuleProposed: async function (code, index, kind, type) {
      let result = null;
      try {
        result = await proposeRule(this.jwtToken, code, index, kind, type);
      } catch (error) {
        result = error.response;
      }

      if (result.status == 200) {
        if (!result.data) {
          console.error('Response successful but no data present:', result);
          return false;
        }

        if (result.data.success) {
          this.$refs.proposal.closeModal();
          this.alert.type = 'success';
          this.alert.msg = 'Your rule was proposed successfully';
          setTimeout(() => {
            this._retireNotification();
          }, 5000);

          // Update round number and reset vote totals
          this.currentRound = result.data.round;
          this.currentTotals.yes = 0;
          this.currentTotals.no = 0;
          this.currentTotals.abstain = 0;
        } else {
          // Response OK but rule proposal failed
          this.$refs.proposal.alert.type = 'danger';
          this.$refs.proposal.alert.msg = 'Rule proposal unsuccessful: "' + result.data.message + '"... Please try again.';
          setTimeout(() => {
            this.$refs.proposal._retireNotification();
          }, 5000);
        }
      } else if (result.status == 500) {
        console.error('Error while trying to propose rule: ', result);
        this.$refs.proposal.alert.type = 'danger';
        this.$refs.proposal.alert.msg = 'There was an error while trying to propose your rule... Please try again.';
        setTimeout(() => {
          this.$refs.proposal._retireNotification();
        }, 5000);
      }
    },
    getCurrentRound: async function () {
      let result = null;
      try {
        result = await getRoundNumber();
      } catch (error) {
        result = error.response;
      }

      if (result.status == 200) {
        if (!result.data) {
          console.error('Response successful but no data present:', result);
          return false;
        }

        this.currentRound = result.data.round;
      } else if (result.status == 500) {
        console.error('Error while trying to get current round number: ', result);
      }
    },
    getCurrentPlayers: async function () {
      let result = null;
      try {
        result = await getPlayers();
      } catch (error) {
        result = error.response;
      }

      if (result.status == 200) {
        if (!result.data) {
          console.error('Response successful but no data present:', result);
          return false;
        }

        // console.log('Players =====>', result);
        this.players = result.data.players;
        this.currentTurn = result.data.currentTurn;
        this.nextTurn = result.data.nextTurn;

        // If user's turn, prompt for rule proposal immediately?
        if (this.currentTurn === this.TwilioIdentity) {
          await this.getLastProposed();
          this.ruleProposalHandler();
        } else {
          // Otherwise just prompt to vote
          await this.getLastProposed();
        }

        // Start round timer
        // this.$refs.totals.startTimer();
      } else if (result.status == 500) {
        console.error('Error while trying to get players: ', result);
      }
    },
    getLastProposed: async function () {
      let result = null;
      try {
        result = await getProposedRule(this.jwtToken, this.currentRound);
      } catch (error) {
        result = error.response;
      }

      if (result.status == 200) {
        if (!result.data) {
          console.error('Response successful but no data present:', result);
          return false;
        }

        // console.log('Proposed rule =====>', result);
        const proposedRule = result.data;
        
        if (typeof this.ruleSets.current[proposedRule.index] !== 'undefined') {
          if (proposedRule.proposal !== ruleChangeTypes.CREATE) {
            proposedRule.original = this.ruleSets.current[proposedRule.index].code;
          }

          if (this.votingCandidate.author !== this.TwilioIdentity) {
            this.votingCandidate = proposedRule;
          }
        }
      } else if (result.status == 500) {
        console.error('Error while trying to get proposed rule: ', result);
      }
    },
    getCurrentRules: async function () {
      this.ruleSets.current = CURRENT_RULES;
    },
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