<template>
  <div>

    <!-- Not Connected -->
    <div class="container" v-if="!connected">
      <h1>{{ title }}</h1>
      <h5>{{ subtitle }}</h5>
      <ul class="list-unstyled" v-if="!connected">
        <li @click="connectUser()">
          <button class="btn btn-primary btn-connect">Login With Tezos</button>
        </li>
      </ul>
      <p>Connect your Tezos wallet to get started</p>
    </div>
    
    <!-- Connected -->
    <div class="container" v-else>
      <h1>{{ title }}</h1>
      <h5>{{ subtitle }}</h5>

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
      </section>


    </div>

  </div>
</template>


<script>
import { 
  Tezos,
  mountProvider,
  getBalance,
  signMessage
} from '../../services/tezProvider';

import { 
  getToken,
  Twilio 
} from '../../services/twilioProvider';

export default {
  data: () => ({
    title: "Nomic Battlegrounds",
    subtitle: "Pwned or be pwned: the choice is yours",
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
    loginSigned: null
  }),
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
        }
      } catch (e) {
        // Auth failed
        console.log('Error requesting Twilio Auth Token', e);
      }

      // Sign login auth message for API
      await this.doLoginMessageSigning();
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
            }
          } catch (e) {
            // Auth failed
            console.log('Error requesting Twilio Auth Token', e);
          }

          // Sign login auth message for API
          await this.doLoginMessageSigning();
        }
      }
    },
    doLoginMessageSigning: async function () {
      let timestamp = new Date().getTime();
      let signedMsg = await this.signMessage(timestamp);
      this.loginSigned = signedMsg
      console.log('Signed message result', this.loginSigned);
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
    setupChatChannel: function () {//here
      console.log("Setting up chat channel", this.chatChannel);

      /**
       * Subscribe to incoming messages sent to the channel
       * @param {Object} message : A Twilio Message object container the {String} properties: `author` and `body`
       */
      this.chatChannel.on('messageAdded', (message) => {
        //console.log("Received message from " + message.author, message.body);
        // Determine message origins
        let messageType;
        if (message.author == this.TwilioIdentity) {
          messageType = "local";
        } else {
          messageType = "remote";
        }

        // Parse message args
        let messageOutput = {
          type: messageType,
          author: message.author,
          msg: message.body
        };

        // Release message to UI
        this.chatMessages.push(messageOutput);
        // console.log(this.chatMessages);
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
    }
  }
};
</script>

<style scoped>
  /* Chat */
  .input-group.chat-controls {
    width: 80vw;
  }
</style>