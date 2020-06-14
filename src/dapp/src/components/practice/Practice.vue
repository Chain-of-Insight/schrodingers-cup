<template>

  <div>
    <!-- Notifications -->
    <Notification 
      :type="alert.type" 
      :msg="alert.msg" 
      v-on:reset="alert = {type: null, msg: null}"
    ></Notification>

    <div class="container-xl main">
      <div class="row" v-if="!activeGame">
        <div class="col">
          <h1>{{ title }}</h1>
          <h5>{{ subtitle }}</h5>
        </div>
      </div>
      <div class="row">
        <div class="col">
      
          <!-- Not Connected -->
          <div class="container-fluid p-0" v-if="!connected">
            <ul class="list-unstyled">
              <li @click="connectUser()">
                <button class="btn btn-primary btn-connect">Login With Tezos</button>
              </li>
            </ul>
            <p>Connect your Tezos wallet to get started</p>
          </div>
      
          <!-- Connected -->
          <div class="ide container-fluid p-0" v-else>
            <div class="row no-gutters">
              <!-- IDE Saved Rule Sets -->
              <div id="ide-saved" class="ide-pane col">
                <label>
                  <strong>Saved Rules:</strong>
                </label>
                <div class="list-group">
                  
                  <!-- Saved Rule Sets -->
                  <a
                    class="ruleset"
                    v-bind:key="index"
                    v-for="(ruleSet, index) in ide.savedRuleSets.player"
                    v-on:click="loadRuleSet(index)"
                    style="cursor:pointer;"
                  >
                    <span class="ruleset success bg-success active list-group-item" v-if="ruleSet.active">{{index + 1}}. {{ ruleSet.name }}</span>
                    <span v-else class="ruleset list-group-item">{{index + 1}}. {{ ruleSet.name }}</span>
                  </a>

                </div>
              </div>
              <!-- IDE Input -->
              <div id="ide-input" class="ide-pane col">
                <label>
                  <strong>Input New Rule:</strong>
                </label>
                <codemirror 
                  v-model="ide.input"
                  :options="ide.options"
                  @input="clearEditorOutput()"
                ></codemirror>
                <div class="execute">
                  <!-- Compile Nomic -->
                  <button class="btn btn-primary" @click="testRuleSet()">Compile</button>
                  
                  <!-- Save Rule Set -->
                  <button 
                    class="btn btn-success" 
                    data-toggle="modal" 
                    data-target="#save-modal"
                    v-if="typeof selectedRuleSet !== 'number'"
                    :disabled="!ide.output || compilerError"
                  >Save</button>

                  <button 
                    class="btn btn-success" 
                    v-if="typeof selectedRuleSet == 'number'"
                    @click="saveRuleSet(selectedRuleSet)"
                    :disabled="!ide.output || compilerError"
                  >Save</button>

                  <!-- Clear Editor -->
                  <button class="btn btn-danger" @click="clearEditor()">
                    <span class="oi oi-trash" title="Clear editor" aria-hidden="true"></span>
                  </button>
                </div>
              </div>
            </div>
            <div class="row no-gutters">
              <!-- IDE Output -->
              <div id="ide-output" class="ide-pane col">
                <label>
                  <strong>Output:</strong>
                </label>
                <div v-if="ide.output">
                  <div class="term-container" v-html="ide.output"></div>
                  <div class="executed">
                    <span class="clear" @click="clearEditorOutput()">clear output</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Save Ruleset Modal -->
    <div class="modal fade" id="save-modal" tabindex="-1" role="dialog" aria-labelledby="save-modal-label" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="save-modal-label">Save Ruleset</h5>
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label for="ruleset-name">Give your ruleset a name:</label>
              <div class="input-group" v-bind:class="{ 'is-invalid': ide.nameError }">
                <input
                  v-model="ide.ruleSetName"
                  type="text"
                  id="ruleset-name"
                  class="form-control"
                  v-bind:class="{ 'is-invalid': ide.nameError }"
                  placeholder="my-awesome-ruleset"
                  required
                >
              </div>
              <div class="invalid-feedback">
                You already have a saved ruleset with that name!
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="cancelSave()" data-dismiss="modal">Cancel</button>
            <button v-if="ide.state.loading" type="button" class="btn btn-success" disabled>
              <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
              <span class="sr-only">Saving...</span>
            </button>
            <button v-else class="btn btn-success" @click="saveRuleSetHandler()">Save</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Ruleset Modal -->
    <div class="modal fade" id="delete-modal" tabindex="-1" role="dialog" aria-labelledby="delete-modal-label" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="delete-modal-label">Delete Ruleset</h5>
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <p><strong>Are you sure you want to delete ruleset '{{ this.ide.selectedRuleSet }}'?</strong></p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="cancelDelete()" data-dismiss="modal">Cancel</button>
            <button v-if="ide.state.loading" type="button" class="btn btn-danger" disabled>
              <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
              <span class="sr-only">Deleting...</span>
            </button>
            <button v-else type="button" class="btn btn-danger" @click="deleteRuleSet()">Delete</button>
          </div>
        </div>
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

// Child components
import Notification from '../common/Notifications.vue';

// IDE
//import 'codemirror/mode/lua/lua';
import 'codemirror/mode/shell/shell';
import 'codemirror/theme/dracula.css';

const $ = window.jQuery;

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
  components: { Notification },
  props: {
    activeGame: {
      default: false,
      type: Boolean
    }
  },
  data: () => ({
    title: "Practice Zone",
    subtitle: "Familiarize yourself with creating and editing Nomic rules",
    network: (process.env.hasOwnProperty('CURRENT_NETWORK')) ? process.env.CURRENT_NETWORK : 'carthagenet',
    address: null,
    getBalance: getBalance,
    connected: false,
    Tezos: Tezos,
    mountProvider: mountProvider,
    alert: {
      type: null,
      msg: null
    },
    ide: {
      input: DEMO_CODE,
      output: null,
      savedRuleSets: [],
      ruleSetName: '',
      nameError: false,
      options: {
        // IDE Options
        // mode: 'text/x-lua',
        tabSize: 4,
        theme: 'dracula',
        lineNumbers: true,
        line: true
      },
      state: {
        loading: false
      },
      execute: testNomic
    },
    compilerError: false,
    selectedRuleSet: null
  }),
  mounted: async function () {
    await this.mountProvider();
    console.log('Practice mounted', this.network);
    let returningUser = sessionStorage.getItem('tzAddress');
    if (returningUser) {
      this.connected = true;
      this.address = returningUser;
    }
    // Get list of saved ruleset names from localStorage
    this.getSavedRuleSets();
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
      console.log('Preparing to test rule', this.ide.input);
      
      // Fetch compiler result
      let result;
      try {
        result = await this.ide.execute(this.ide.input);
      } catch(e) {
        result = e.response;
      }

      console.log('Result', result);

      if (result.status) {
        if (result.status == 500) {
          // console.log('Compile failed', result.status);
          this.alert.type = 'danger';
          this.alert.msg = 'Compile failed';
          this.compilerError = true;
          setTimeout(() => {
            this._retireNotification();
          }, 5000);
        } else if (result.status == 200) {
          // console.log('Compiled successfully,', result.status);
          this.alert.type = 'success';
          this.alert.msg = 'Compiled successfully';
          this.compilerError = false;
          setTimeout(() => {
            this._retireNotification();
          }, 5000);
        }
      }

      // Show compiler feedback
      let output;
      if (result && result.data) {
        if (result.data.resultHtml) {
          output = result.data.resultHtml;
        } 
      }
      if (output) {
        this.ide.output = output;
      } else {
        this.ide.output = "Warning: compilation passed but compiler output was empty 😅";
      }
    },
    getSavedRuleSets: async function (index = false) {
      // Load rule sets if exists
      let storedRuleSets = await localStorage.getItem('ruleSets');
      let playerSavedRuleSets = null;
      if (storedRuleSets) {
        storedRuleSets = await JSON.parse(storedRuleSets);
        playerSavedRuleSets = {
          player: storedRuleSets,
          nomic: []
        };
      }
      this.ide.savedRuleSets = await (playerSavedRuleSets) ? playerSavedRuleSets : {
        player: [],
        nomic: []
      };
      console.log('Stored Rule Sets =>', this.ide.savedRuleSets);
      if (typeof index == 'number') {
        this.loadRuleSet(index);
      }
    },
    saveRuleSetHandler: function () {
      let index = false;
      if (typeof this.selectedRuleSet == 'number') {
        this.saveRuleSet(this.selectedRuleSet);
      } else {
        this.saveRuleSet();
      }

    },
    saveRuleSet: async function (index = false) {
      if (this.compilerError) {
        this.alert.type = 'danger';
        this.alert.msg = 'Error saving rule, compilation failed';
        setTimeout(() => {
          this._retireNotification();
        }, 5000);
        return;
      }

      let ruleSets = localStorage.getItem('ruleSets');
      let hasRuleSets;

      // Create new collection / Add to existing collection
      if (ruleSets) {
        hasRuleSets = true;
        ruleSets = JSON.parse(ruleSets);
        console.log(ruleSets)
      } else {
        hasRuleSets = false;
        ruleSets = [];
      }

      // Loading
      this.ide.state.loading = true;
      
      const ruleSet = {
        code: this.ide.input,
        name: this.ide.ruleSetName
      };

      console.log('Rule set to be saved:', ruleSet);

      // Prepare rule set
      if (hasRuleSets && index) {
        let rName = ruleSets[index].name;
        ruleSet.name = rName;
        ruleSets[index] = ruleSet;
      } else {
        ruleSets.push(ruleSet);
      }
      ruleSets = JSON.stringify(ruleSets)

      // Save rule set
      await localStorage.setItem('ruleSets', ruleSets);

      // Update rule sets
      if (!index) {
        if (this.ide.savedRuleSets.player) {
          if (this.ide.savedRuleSets.player.length) {
            index = this.ide.savedRuleSets.player.length;
          } else {
            index = 0;
          }
        } else {
          index = 0;
        }
      }
      this.getSavedRuleSets(index);

      // Reset app state
      this.alert.type = 'success';
      this.alert.msg = `Ruleset '${this.ide.ruleSetName}' saved!`;
      setTimeout(() => {
        this._retireNotification();
      }, 5000);
      $('#save-modal').modal('hide');
      this.ide.state.loading = false;
      this.ide.ruleSetName = '';
    },
    cancelSave: function () {
      this.ide.ruleSetName = '';
      this.ide.nameError = false;
    },
    loadRuleSet: function (index) {
      let ruleSet = this.ide.savedRuleSets.player[index];
      console.log('Loading rule set =>', [ruleSet, index]);

      if (!ruleSet) {
        // setTimeout(() => {
        //   this.loadRuleSet(index);
        // }, 500);
        return;
      }

      if (ruleSet.hasOwnProperty('code')); {
        
        try {
          // Set IDE state
          this.ide.input = ruleSet.code;
          // Set UI state
          this.ide.savedRuleSets.player[index].active = true;
          this.selectedRuleSet = Number(index);
          for (let i = 0; i < this.ide.savedRuleSets.player.length; i++) {
            if (i !== index) {
              if (this.ide.savedRuleSets.player[i].hasOwnProperty('active')) {
                this.ide.savedRuleSets.player[i].active = false;
              }
              
            }
          }
          // Force update cycle
          this.$forceUpdate();
        } catch(e) {
          console.log('Error loading ruleset =>', [e, ruleSet])
          this.alert.type = 'danger';
          this.alert.msg = `Error loading ruleset '${JSON.stringify(ruleSet)}'`;
          setTimeout(() => {
            this._retireNotification();
          }, 5000);
        }
      }
    },
    clearEditor: function () {
      console.log('Clearing editor...', this.ide);
      this.ide.input = '';
      this.ide.output = null;
      this.selectedRuleSet = null;
      this.compilerError = false;
      
      if (this.ide.savedRuleSets.player) {
        if (this.ide.savedRuleSets.player.length) {
          // Remove active file edit
          for (let i = 0; i < this.ide.savedRuleSets.player.length; i++) {
            if (this.ide.savedRuleSets.player[i].hasOwnProperty('active')) {
              if (this.ide.savedRuleSets.player[i].active) {
                this.ide.savedRuleSets.player[i].active = false;
              }
            }
          }
        }
      }
    },
    clearEditorOutput: function () {
      if (!this.ide.output) {
        return;
      }
      console.log('Clearing output...', this.ide);
      this.ide.output = null;
      this.compilerError = false;
    },
    _retireNotification: function () {
      this.alert = {
        type: null,
        msg: null
      };
    }
  }
};
</script>

<style scoped>
  .ide {
    margin-top: 2rem;
  }
  span.clear {
    text-decoration: underline;
    cursor: pointer;
    clear: both;
  }
  a.ruleset {
    text-decoration: none;
  }
  span.ruleset,
  span.ruleset:hover {
    color: black;
    border: none;
    margin-right: 1em;
  }

  span.ruleset.active,
  span.ruleset.active:hover {
    color: white;
    border: none;
  }
</style>