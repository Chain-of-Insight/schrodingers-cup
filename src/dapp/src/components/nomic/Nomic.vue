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
            <button type="button" class="btn btn-outline-primary" @click="getLastProposed()">
              Test Voting
            </button>
          </div>
        </div>
        <RuleProposal
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
import * as api from '../../services/apiProvider';

// Child components
import Notification from '../common/Notifications.vue';
import RuleProposal from '../rule-proposal/RuleProposal.vue';
import Voting from '../common/Voting.vue';
import Totals from '../common/Totals.vue';
// IDE Component
import Practice from '../practice/Practice.vue';

// Constants
import {
  proposalTypes,
  voteTypes,
  TZ_WALLET_PATTERN,
  CURRENT_RULES
} from '../../constants/constants.js';


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
    },
    currentVotes: [],
    votedThisRound: false,
    proposedThisRound: false,
  }),
  computed: {
    msgPatterns: function () {
      return {
        PROPOSAL_PATTERN: `^${TZ_WALLET_PATTERN} proposed a rule in round (\\d+)$`,
        VOTE_PATTERN: `^${TZ_WALLET_PATTERN} successfully voted (YES|NO) in round (\\d+)$`,
      }
    }
  },
  watch: {
    currentVotes: function (votes) {
      votes.forEach(castVote => {
        if (typeof castVote.player !== 'string' || typeof castVote.vote !== 'boolean')
          return;

        if (castVote.player === this.address) {
          this.votedThisRound = true;

          if (this.currentTurn === this.address) {
            this.proposedThisRound = true;
          }
        }

        switch (castVote.vote) {
          case true:
            this.currentTotals.yes ++;
            break;
          case false:
            this.currentTotals.no ++;
            break;
        }
      });
    }
  },
  mounted: async function () {
    await this.mountProvider();
    console.log('Nomic mounted', this.network);

    // Check if user already logged in with TezBridge
    this.address = sessionStorage.getItem('tzAddress');
    if (this.address != null) {
      this.connected = true;
      await this.gameSetup();
    }
  },
  methods: {
    gameSetup: async function () {
      // Get current rule set
      this.getCurrentRules();
      // Get current round number
      await this.getCurrentRound();
      // Get players
      await this.getCurrentPlayers();
      // Get votes
      await this.getCurrentVotes();
      // Get Twilio token
      await this.twilioAuth();
      // Connect to chat room
      this.connectChat();
      // Sign login auth message for API
      await this.doLoginMessageSigning();
      // Get last proposed rule if applicable
      await this.getLastProposed();
      // If user's turn, prompt for rule proposal immediately
      if (this.currentTurn === this.address) {
        this.ruleProposalHandler();
      }
    },
    connectUser: async function () {
      // Connect as required
      if (this.connected)
        return;

      // Signer instance
      this.address = await tezbridge.request({method: 'get_source'});
      
      // Fetch balance / Connection callbacks
      if (typeof this.address == 'string' && this.address.length === 36) {
        console.log('User XTZ Address =>', this.address);
        sessionStorage.setItem('tzAddress', this.address);
        this.connected = true;
        await this.gameSetup();
        let balance = await this.getBalance(this.address);
        console.log("User balance =>", balance);
      }
    },
    twilioAuth: async function () {
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
        }
      } catch (e) {
        // Auth failed
        console.error('Error requesting Twilio Auth Token', e);
      }
    },
    doLoginMessageSigning: async function () {
      const jwt = sessionStorage.getItem('jwt');
      // don't try to sign again if jwt already present
      // TODO: handle jwt expiry?
      if (jwt != null) {
        this.jwtToken = jwt;
        return;
      }

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
        result = await api.PerformAuth(this.loginSigned.bytes, this.loginSigned.prefixSig, pubKey, this.address);
      } catch(e) {
        // unauthorized if we are here
        result = e.response;
      }

      if (result.data && result.data.token) {
        this.jwtToken = result.data.token;
        sessionStorage.setItem('jwt', this.jwtToken);
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
        // console.log('chat message!', message);

        if (message.author === 'system') {
          this.systemMsgHandler(message);
        } else {
          this.userMsgHandler(message);
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
    systemMsgHandler: async function (message) {
      let matches = null;
      let round = null;
      let vote = null;
      let player = null;

      // GAME EVENTS
      switch (message.body) {

        // NEW RULE PROPOSED
        case (message.body.match(RegExp(this.msgPatterns.PROPOSAL_PATTERN)) || {}).input:
          matches = message.body.match(RegExp(this.msgPatterns.PROPOSAL_PATTERN));
          // Parse player wallet address
          player = matches[1];
          // Parse round number for incrementing
          round = parseInt(matches[2]);

          // Update round number
          await this.getCurrentRound();
          // Increment YES vote once for the player who proposed
          if (player !== this.address && round === this.currentRound) {
            this.currentTotals.yes ++;
            // Get proposed rule for voting
            await this.getLastProposed();
          }
          break;

        // VOTE CAST
        case (message.body.match(RegExp(this.msgPatterns.VOTE_PATTERN)) || {}).input:
          matches = message.body.match(RegExp(this.msgPatterns.VOTE_PATTERN));
          // Parse vote (YES/NO)
          vote = matches[2];
          // Parse round number for incrementing
          round = parseInt(matches[3]);

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
          break;
      }

      // Release message to UI
      let messageOutput = {
        type: 'system',
        msg: message.body
      };
      this.chatMessages.push(messageOutput);
    },
    userMsgHandler: function (message) {
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
      if (vote === voteTypes.ABSTAIN) {
        // Don't send anything on abstain, just 'close' voting section
        this.votingCandidate = null;
        return false;
      }

      let result = null;
      try {
        result = await api.castVote(this.jwtToken, Boolean(vote), this.currentRound);
      } catch (error) {
        console.error('Error while trying to cast vote:', error);
        if (error.response) {
          result = error.response;
        } else {
          result = error;
        }
      }

      if (result.status && result.status == 200) {
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

          // clear voting candidate and voting ribbon
          this.votingCandidate = null;
        } else {
          // Response OK but vote cast failed
          this.alert.type = 'danger';
          this.alert.msg = 'Vote cast unsuccessful: "' + result.data.message + '"';
          setTimeout(() => {
            this._retireNotification();
          }, 5000);
        }
      } else {
        console.error('Error while trying to propose rule: ', result);
        this.alert.type = 'danger';
        this.alert.msg = 'There was an error while trying to cast your vote';
        setTimeout(() => {
          this._retireNotification();
        }, 5000);
      }
    },
    ruleProposalHandler: function () {
      if (this.proposedThisRound) {
        console.log('You already proposed this round! Skipping proposal prompt...');
        return;
      }

      if (!this.jwtToken) {
        this.alert.type = 'danger';
        this.alert.msg = "It's your turn to propose a rule, but you havent' been authenticated yet! Have you signed a message in TezBridge?";
        setTimeout(() => {
          this._retireNotification();
        }, 3000);
      } else {
        this.votingCandidate = null;

        // Otherwise, prompt to propose a rule
        this.$refs.proposal.promptForProposal();
      }
    },
    onRuleProposed: async function (code, index, kind, type) {
      let result = null;
      try {
        result = await api.proposeRule(this.jwtToken, code, index, kind, type);
      } catch (error) {
        console.error('Error while trying to propose rule:', error);
        if (error.response) {
          result = error.response;
        } else {
          result = error;
        }
      }

      if (result.status && result.status == 200) {
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
          this.currentTotals.yes = 1; // You vote YES on your own rule by default
          this.currentTotals.no = 0;
          this.currentTotals.abstain = 0;
          this.votedThisRound = true;
          this.proposedThisRound = true;
        } else {
          // Response OK but rule proposal failed
          this.$refs.proposal.alert.type = 'danger';
          this.$refs.proposal.alert.msg = 'Rule proposal unsuccessful: "' + result.data.message + '"';
          setTimeout(() => {
            this.$refs.proposal._retireNotification();
          }, 5000);
        }
      } else {
        console.error('Error while trying to propose rule: ', result);
        this.$refs.proposal.alert.type = 'danger';
        this.$refs.proposal.alert.msg = 'There was an error while trying to propose your rule';
        setTimeout(() => {
          this.$refs.proposal._retireNotification();
        }, 5000);
      }
    },
    getCurrentRound: async function () {
      let result = null;
      try {
        result = await api.getRoundNumber();
      } catch (error) {
        console.error('Error while trying to get current round:', error);
        if (error.response) {
          result = error.response;
        } else {
          result = error;
        }
      }

      if (result.status && result.status == 200) {
        if (!result.data) {
          console.error('Response successful but no data present:', result);
          return false;
        }

        this.currentRound = result.data.round;
      } else {
        console.error('Error while trying to get current round number: ', result);
      }
    },
    getCurrentPlayers: async function () {
      let result = null;
      try {
        result = await api.getPlayers();
      } catch (error) {
        console.error('Error while trying to get players:', error);
        if (error.response) {
          result = error.response;
        } else {
          result = error;
        }
      }

      if (result.status && result.status == 200) {
        if (!result.data) {
          console.error('Response successful but no data present:', result);
          return false;
        }

        // console.log('Players =====>', result);
        this.players = result.data.players;
        this.currentTurn = result.data.currentTurn;
        this.nextTurn = result.data.nextTurn;
      } else {
        console.error('Error while trying to get players: ', result);
      }
    },
    getLastProposed: async function () {
      if (this.votedThisRound) {
        console.log('You already voted this round! Skipping fetch for proposed rule...');
        return;
      }

      let result = null;
      try {
        result = await api.getProposedRule(this.currentRound);
      } catch (error) {
        if (error.response) {
          result = error.response;
        } else {
          result = error;
        }
      }

      if (result.status && result.status == 200) {
        if (!result.data) {
          console.error('Response successful but no data present:', result);
          return false;
        }


        // console.log('Proposed rule =====>', result);
        const proposedRule = result.data;
        
        // Make sure the response data has the important stuff...
        if (typeof proposedRule.code === 'string' && typeof proposedRule.index === 'number') {
          if (proposedRule.proposal !== proposalTypes.CREATE) {
            // Store current rule code for voting display
            proposedRule.original = this.ruleSets.current[proposedRule.index].code;
          }
          // Store proposed rule for voting (triggers showing of voting 'ribbon')
          this.votingCandidate = proposedRule;
        }
      } else {
        console.error('Error while trying to get proposed rule:', result);
      }
    },
    getCurrentRules: function () {
      this.ruleSets.current = CURRENT_RULES;
    },
    getCurrentVotes: async function () {
      let result = null;
      try {
        result = await api.getVotes(this.currentRound);
      } catch (error) {
        console.error('Error while trying to get votes:', error);
        if (error.response) {
          result = error.response;
        } else {
          result = error;
        }
      }

      if (result.status && result.status == 200) {
        if (!result.data) {
          console.error('Response successful but no data present:', result);
          return false;
        }

        console.log('Votes =====>', result);
        if (result.data.Votes instanceof Array) {
          this.currentVotes = result.data.Votes;
        }
      } else {
        console.error('Error while trying to get votes: ', result);
      }
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