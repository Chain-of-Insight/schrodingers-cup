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

      <div id="ide" class="ide input">
        <label>
          <strong>Input New Rule:</strong>
        </label>
        <codemirror 
          v-model="ide.input"
          :options="ide.options"
        ></codemirror>
        <div class="execute">
          <button class="btn btn-primary" @click="testRuleSet()">Compile</button>
          <!--
          <button class="btn btn-danger" @click="clearEditor()">Erase</button>
          -->
        </div>
      </div>

      <div class="ide output">
        <textarea 
          v-model="ide.output" 
          v-if="ide.output"
          readonly
        ></textarea>
      </div>
    </div>

  </div>
</template>


<script>
import { 
  Tezos,
  mountProvider,
  getBalance
} from '../../services/tezProvider';

// API
import { testNomic } from '../../services/apiProvider';

// IDE
//import 'codemirror/mode/lua/lua';
import 'codemirror/mode/shell/shell';
import 'codemirror/theme/dracula.css';

const DEMO_CODE = `### In Nomsu, variables have a "$" prefix, and you can just assign to them
### without declaring them first:
$x = 1
test that ($x == 1)

### Variables which have not yet been set have the value (nil)
test that ($not_yet_set == (nil))

### Variables can be nameless:
$ = 99

### Or have spaces, if surrounded with parentheses:
$(my favorite number) = 23

### Figure out what value $my_var should have:
$my_var = 100
$my_favourite_number = 1
$x = 0
$my_var = ($my_var + $x + $my_favourite_number)
test that ($my_var == 101)
say("OK!")`;

export default {
  data: () => ({
    title: "Practice Zone",
    subtitle: "Familiarize yourself with creating and editing Nomic rules",
    network: (process.env.hasOwnProperty('CURRENT_NETWORK')) ? process.env.CURRENT_NETWORK : 'carthagenet',
    address: null,
    getBalance: getBalance,
    connected: false,
    Tezos: Tezos,
    mountProvider: mountProvider,
    ide: {
      input: DEMO_CODE,
      output: null,
      options: {
        // IDE Options
        // mode: 'text/x-lua',
        tabSize: 2,
        theme: 'dracula',
        lineNumbers: true,
        line: true
      },
      state: {
        loading: false
      },
      execute: testNomic
    }
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
    },
    testRuleSet: async function () {
      console.log('Preparing to test rule', this.ide);
      let result = await this.ide.execute(this.ide.input);
      console.log('Result', result);
    },
    saveRuleSet: function () {
      // TODO: this
    },
    loadRuleSet: function () {
      // TODO: this
    }//,
    // clearEditor: function () {
    //   console.log('Clearing editor...', this.ide);
    //   this.ide.input = null;
    // }
  }
};
</script>

<style scoped>
  #ide {
    margin-top: 2rem;
  }
</style>