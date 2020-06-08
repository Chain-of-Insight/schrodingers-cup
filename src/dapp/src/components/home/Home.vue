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
      <!-- Live Gameplay Countdown -->
      <div class="live-game-countdown" v-if="!live && liveStateVerified">
        <div class="jumbotron">
          <p class="descr">Next gameplay session begins:<br/><h5 class="game-session-countdown counter-display">{{ nextGameCountdown }}<h5></p>
        </div>
      </div>
      <!-- Loading Gameplay Countdown -->
      <div class="loading" v-if="!live && !liveStateVerified">
        <img src="../../assets/img/loading.gif">
      </div>
      <!-- Gameplay is Live -->
      <div class="live-game" v-if="live && liveStateVerified">
        <p>A gameplay session is currently in progress</p>
        <!-- 
          TODO XXX (drew): 
          Add nav. button here to proceed to live game 
        -->
      </div>
    </div>


    <div class="footer">
      <img src="../../assets/img/schrodingers-cup.png" />
    </div>

  </div>
</template>


<script>
import { 
  Tezos,
  mountProvider,
  getBalance
} from '../../services/tezProvider';

export default {
  data: () => ({
    title: "Welcome to Schrodinger's Cup",
    subtitle: "A game of Nomic running on Tezos",
    network: (process.env.hasOwnProperty('CURRENT_NETWORK')) ? process.env.CURRENT_NETWORK : 'carthagenet',
    address: null,
    getBalance: getBalance,
    connected: false,
    Tezos: Tezos,
    mountProvider: mountProvider,
    live: null,
    liveStateVerified: false,
    liveDailyHourUTC: (process.env.hasOwnProperty('UTC_START')) ? process.env.UTC_START : 16,
    liveDailyHourClosedUTC: (process.env.hasOwnProperty('UTC_END')) ? process.env.UTC_END : 17,
    nextGamePlaySessionStart: null,
    nextGamePlayStartOffset: null
  }),
  mounted: async function () {
    await this.mountProvider();
    console.log('Home mounted', this.network);
    let returningUser = sessionStorage.getItem('tzAddress');
    if (returningUser) {
      this.connected = true;
      this.address = returningUser;
    }

    // Calculate next game time
    this.nextGame();
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
        }
      }
    },
    nextGame: function () {
      // Countdown bounding props.
      const MINUTES_TO_MILLISECONDS_MULTIPLIER = 60000;
      let baseDate = new Date();                        // <= Basic date obj. (user timezone)
      let userTime = baseDate.getTime();                // <= as epoch milliseconds
      let userOffsetUTC = baseDate.getTimezoneOffset(); // <= User UTC offset in minutes
      let mult = MINUTES_TO_MILLISECONDS_MULTIPLIER;
      userOffsetUTC = userOffsetUTC * mult;             // <= User UTC offset in milliseconds
      let utcTime = userTime + (userOffsetUTC);         // <= Current time as UTC date obj.

      // Determine occurence of next live gameplay session
      // 1) Has gameplay today occurred
      let utcHour = new Date(utcTime).getHours();
      
      // 2) Knowing the current UTC hour, determine if the next gameplay session is:
      // "TODAY", "TOMORROW" or "IN PROGRESS"
      
      // i) Live Game In Progress
      if (utcHour >= this.liveDailyHourUTC && utcHour < this.liveDailyHourClosedUTC) {
        this.live = true;
        this.liveStateVerified = true;
        
        // Set next gameplay start (now)
        this.nextGamePlayStartOffset = 0;
        this.nextGamePlaySessionTime = new Date();

      // ii) Next Gameplay Session is Tomorrow
      } else if (utcHour > this.liveDailyHourClosedUTC) {
        this.live = false;
        this.liveStateVerified = true;

        // Calculate date difference
        let today = new Date(utcTime);
        let tomorrow = new Date(today);
        tomorrow.setDate(tomorrow.getDate() + 1);
        let tomorrowDay = tomorrow.getDate();
        let tomorrowMonth = tomorrow.getMonth();
        let tomorrowYear = tomorrow.getFullYear();
        let tomorrowBegin = new Date(tomorrowYear, tomorrowMonth, tomorrowDay).getTime();
        let tomorrowGameStartOffset = (this.liveDailyHourUTC * 60) * MINUTES_TO_MILLISECONDS_MULTIPLIER;
        let nextGameplaySesssion = new Date(tomorrowBegin + tomorrowGameStartOffset);

        // Set next gameplay start
        this.nextGamePlayStartOffset = (nextGameplaySesssion.getTime()) - userTime;
        this.nextGamePlaySessionStart = nextGameplaySesssion;

        // Start timer
        setInterval(() => {
          this.nextGamePlayStartOffset = this.nextGamePlayStartOffset - 1000;
        }, 1000);

      // iii) Next Gameplay Session is Later Today
      } else {
        this.live = false;
        this.liveStateVerified = true;
        
        // Calculate date difference
        let today = new Date(utcTime);
        let todayHours = (today.getHours() * 60) * MINUTES_TO_MILLISECONDS_MULTIPLIER;
        let todayHoursRemaining = ((this.liveDailyHourUTC * 60) * MINUTES_TO_MILLISECONDS_MULTIPLIER) - todayHours;
        let todayMinutes = today.getMinutes() * MINUTES_TO_MILLISECONDS_MULTIPLIER;
        let todaySeconds = today.getSeconds() * (MINUTES_TO_MILLISECONDS_MULTIPLIER / 60);
        let gameStartOffset = todayHoursRemaining - todayMinutes - todaySeconds;

        // Do calculation for next session
        let nextGamePlaySessionTime = new Date(utcTime + gameStartOffset);

        // Set next gameplay start
        this.nextGamePlayStartOffset = gameStartOffset;
        this.nextGamePlaySessionStart = nextGamePlaySessionTime;

        // Start timer
        setInterval(() => {
          this.nextGamePlayStartOffset = this.nextGamePlayStartOffset - 1000;
        }, 1000);
      }
      console.log('Gameplay Start', [this.nextGamePlaySessionStart, this.nextGamePlayStartOffset])
    }
  },
  computed: {
    nextGameCountdown: function () {
      let countdown = '';
      // Test data integrity before display
      if (!this.liveDailyHourUTC || !this.nextGamePlaySessionStart || !this.nextGamePlayStartOffset) {
        return countdown;
      } else if (isNaN(this.liveDailyHourUTC) || isNaN(this.nextGamePlayStartOffset)) {
        return countdown;
      } else if (this.nextGamePlayStartOffset < 0) {
        // Gameplay already in progress
        this.live = true;
        return '00:00:00';
      } else {
        if (this.live) {
          this.live = false;
        }
      }

      // Build timer countdown
      let hours,
          minutes,
          seconds;
      
      // Hours
      hours = parseInt(this.nextGamePlayStartOffset / (3600 * 1000));
      
      // Minutes
      minutes = (this.nextGamePlayStartOffset / (3600 * 1000)) - hours;
      minutes = parseInt(minutes * 60);

      // Seconds
      seconds = (this.nextGamePlayStartOffset / 1000);
      seconds = seconds - (hours * 3600) - (minutes * 60);
      // console.log(seconds);
      
      // Zero prefixing and decimal safety
      if (hours < 10) {
        hours = String('0' + parseInt(hours));
      } else {
        hours = parseInt(hours);
      }
      if (minutes < 10) {
        minutes = String('0' + parseInt(minutes));
      } else {
        minutes = parseInt(minutes);
      }
      if (seconds < 10) {
        seconds = String('0' + parseInt(seconds));
      } else {
        seconds = parseInt(seconds);
      }

      // Countdown final format
      countdown = hours + ':' + minutes + ':' + seconds;
      
      // Update timer in UI
      return countdown;
    }
  }
};
</script>

<style scoped>
  .live-game-countdown {
    text-align: center;
  }
</style>