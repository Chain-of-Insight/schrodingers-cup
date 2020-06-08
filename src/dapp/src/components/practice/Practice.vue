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
    title: "Practice Zone",
    subtitle: "Familiarize yourself with creating and editing Nomic rules",
    network: (process.env.hasOwnProperty('CURRENT_NETWORK')) ? process.env.CURRENT_NETWORK : 'carthagenet',
    address: null,
    getBalance: getBalance,
    connected: false,
    Tezos: Tezos,
    mountProvider: mountProvider
  }),
  mounted: async function () {
    await this.mountProvider();
    console.log('Home mounted', this.network);
    let returningUser = sessionStorage.getItem('tzAddress');
    if (returningUser) {
      this.connected = true;
      this.address = returningUser;
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
        }
      }
    }
  }
};
</script>

<style scoped>
  /* .live-game-countdown {
    text-align: center;
  } */
</style>