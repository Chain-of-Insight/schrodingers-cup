<template>

  <div>
    <!-- Notifications -->
    <Notification 
      :type="alert.type" 
      :msg="alert.msg" 
      v-on:reset="alert = {type: null, msg: null}"
    ></Notification>

    <div class="container-fluid main">
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
              <div id="ide-saved" class="ide-pane col-3">
                <label>
                  <strong>Rule Sets:</strong>
                </label>
                <nav class="nav nav-tabs">
                  <a
                    class="nav-link"
                    @click="ide.ruleSetPane = ruleSetTypes.CURRENT"
                    v-bind:class="{ active: ide.ruleSetPane === ruleSetTypes.CURRENT }"
                  >Current</a>
                  <a
                    class="nav-link"
                    @click="ide.ruleSetPane = ruleSetTypes.SAVED"
                    v-bind:class="{ active: ide.ruleSetPane === ruleSetTypes.SAVED }"
                  >Saved</a>
                  <a
                    class="nav-link"
                    @click="ide.ruleSetPane = ruleSetTypes.QUEUED"
                    v-bind:class="{ active: ide.ruleSetPane === ruleSetTypes.QUEUED }"
                  >Queued</a>
                </nav>
                <div
                  id="rules-current"
                  class="list-group"
                  v-if="ide.ruleSetPane === ruleSetTypes.CURRENT"
                >
                  <p
                    class="list-group-item border-0"
                    v-if="ide.savedRuleSets.nomic.length === 0"
                  >No current rule sets loaded...</p>
                  <!-- Current Rule Sets -->
                  <a
                    class="ruleset btn btn-outline-success list-group-item list-group-item-action"
                    role="button"
                    tabindex="0"
                    v-bind:class="{
                      'active':
                        selectedRuleSet.type === ruleSetTypes.CURRENT &&
                        selectedRuleSet.index === index
                    }"
                    v-bind:key="index"
                    v-for="(ruleSet, index) in ide.savedRuleSets.nomic"
                    v-on:click="loadRuleSet(index, ruleSetTypes.CURRENT)"
                    style="cursor:pointer;"
                  >
                    <span>{{index + 1}}. {{ ruleSet.name }}</span>
                  </a>
                </div>
                <div
                  id="rules-saved"
                  class="list-group"
                  v-if="ide.ruleSetPane === ruleSetTypes.SAVED"
                >
                  <p
                    class="list-group-item border-0"
                    v-if="ide.savedRuleSets.player.length === 0"
                  >No saved rule sets...</p>
                  <!-- Saved Rule Sets -->
                  <div
                    class="ruleset btn-group"
                    role="group"
                    v-bind:key="index"
                    v-for="(ruleSet, index) in ide.savedRuleSets.player"
                    v-on:click="loadRuleSet(index, ruleSetTypes.SAVED)"
                  >
                    <a
                      class="btn btn-outline-success list-group-item list-group-item-action"
                      role="button"
                      tabindex="0"
                      v-bind:class="{
                        'active':
                          selectedRuleSet.type === ruleSetTypes.SAVED &&
                          selectedRuleSet.index === index
                      }"
                      style="cursor:pointer;"
                    >
                      <span>{{index + 1}}. {{ ruleSet.name }}</span>
                    </a>
                    <button
                      class="btn btn-outline-warning"
                      @click="queueRuleSet(index)"
                    >Queue</button>
                  </div>
                </div>
                <div
                  id="rules-current"
                  class="list-group"
                  v-if="ide.ruleSetPane === ruleSetTypes.QUEUED"
                >
                  <!-- Queued Rule Sets -->
                  <p
                    class="list-group-item border-0"
                    v-if="queuedRuleSets.length === 0"
                  >No queued rule sets...</p>
                  <a
                    class="ruleset btn btn-outline-success list-group-item list-group-item-action"
                    v-bind:class="{
                      'active':
                        selectedRuleSet.type === ruleSetTypes.QUEUED &&
                        selectedRuleSet.index === index
                    }"
                    v-bind:key="index"
                    v-for="index in queuedRuleSets"
                    v-on:click="loadRuleSet(index, ruleSetTypes.QUEUED)"
                    style="cursor:pointer;"
                  >
                    <span>{{index + 1}}. {{ ide.savedRuleSets.player[index].name }}</span>
                  </a>
                </div>
              </div>
              <!-- IDE Input -->
              <div id="ide-input" class="ide-pane col-9">
                <label>
                  <strong>Rule Editor:</strong>
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
                    @click="saveRuleSetHandler()"
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
              <div id="ide-output" class="ide-pane-b col">
                <label>
                  <strong>Output:</strong>
                </label>
                <div class="ide-output-wrapper row no-gutters" v-if="ide.output">
                  <div class="executed clear-output">
                    <span class="clear" @click="clearEditorOutput()">clear output</span>
                  </div>
                  <div class="term-container-wrapper">
                    <div class="term-container" v-html="ide.output"></div>
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
            <button v-else class="btn btn-success" @click="saveRuleSet()">Save</button>
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

const ruleSetTypes = {
  SAVED: 'ACTIVE',
  CURRENT: 'CURRENT',
  QUEUED: 'QUEUED'
}

const CURRENT_RULES = require('./rules/currentRules.json');

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
    ruleSetTypes: ruleSetTypes,
    ide: {
      input: DEMO_CODE,
      output: null,
      ruleSetPane: ruleSetTypes.CURRENT,
      savedRuleSets: {
        player: [],
        nomic: []
      },
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
    selectedRuleSet: {
      type: null,
      index: null
    }
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
    this.getCurrentRuleSets();
  },
  computed: {
    queuedRuleSets: function () {
      // Store indexes of queued rules
      const queuedIndexes = [];
      this.ide.savedRuleSets.player.forEach((ruleSet, index) => {
        if (ruleSet.queued) {
          queuedIndexes.push(index);
        }
      });
      return queuedIndexes;
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
    getCurrentRuleSets: async function () {
      this.ide.savedRuleSets.nomic = CURRENT_RULES;
    },
    getSavedRuleSets: async function (index = false) {
      // Load rule sets if exists
      let storedRuleSets = await localStorage.getItem('ruleSets');
      if (storedRuleSets) {
        storedRuleSets = await JSON.parse(storedRuleSets);
        this.ide.savedRuleSets.player = storedRuleSets;
      }
      console.log('Stored Rule Sets =>', this.ide.savedRuleSets);
      // if (typeof index == 'number') {
      //   this.loadRuleSet(index);
      // }
    },
    saveRuleSetHandler: function () {
      if (
        this.selectedRuleSet.type === ruleSetTypes.CURRENT ||
        typeof(this.selectedRuleSet.index) !== 'number'
      ) {
        $('#save-modal').modal('show');
      } else {
        this.saveRuleSet(this.selectedRuleSet.index);
      }
    },
    saveRuleSet: async function (index = null) {
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

      let ruleSet = null;

      // Prepare rule set
      if (hasRuleSets && index !== null) {
        console.log(42, index);
        ruleSets[index].code = this.ide.input;
        ruleSet = ruleSets[index];
        console.log(1);
      } else {
        console.log(2, index);
        ruleSet = {
          code: this.ide.input,
          name: this.ide.ruleSetName
        };
        ruleSets.push(ruleSet);
      }
      ruleSets = JSON.stringify(ruleSets);

      console.log('Rule set to be saved:', ruleSet);

      // Save rule set
      await localStorage.setItem('ruleSets', ruleSets);

      // Update rule sets
      if (index === null) {
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

      await this.getSavedRuleSets();
      this.loadRuleSet(index, ruleSetTypes.SAVED);
      this.ide.ruleSetPane = ruleSetTypes.SAVED;

      // Reset app state
      this.alert.type = 'success';
      this.alert.msg = `Ruleset '${ruleSet.name}' saved!`;
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
    loadRuleSet: function (index, type) {
      let ruleSetList = null;

      switch (type) {
        case ruleSetTypes.CURRENT:
          ruleSetList = this.ide.savedRuleSets.nomic;
          break;
        case ruleSetTypes.SAVED || ruleSetTypes.QUEUED:
          ruleSetList = this.ide.savedRuleSets.player;
          break;
        // case ruleSetTypes.QUEUED:
        //   ruleSetList = this.queuedRuleSets;
        //   break;
        default:
          ruleSetList = this.ide.savedRuleSets.nomic;
          break;
      }

      let ruleSet = ruleSetList[index];
      console.log('Loading rule set =>', [ruleSet, index]);

      if (!ruleSet) {
        // setTimeout(() => {
        //   this.loadRuleSet(index);
        // }, 500);
        return;
      }

      if (ruleSet.hasOwnProperty('code')) {
        try {
          // Set IDE state
          this.ide.input = ruleSet.code;
          // Set UI state
          // ruleSetList[index].active = true;
          this.selectedRuleSet.type = type;
          this.selectedRuleSet.index = Number(index);
          // for (let i = 0; i < ruleSetList.length; i++) {
          //   if (i !== index) {
          //     if (ruleSetList[i].hasOwnProperty('active')) {
          //       ruleSetList[i].active = false;
          //     }
              
          //   }
          // }
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
    queueRuleSet: function (index) {
      console.log('Queuing rule set:', index);
      this.ide.savedRuleSets.player[index].queued = true;
      // Update rule set in localStorage and refresh
      localStorage.setItem('ruleSets', JSON.stringify(this.ide.savedRuleSets.player));
      this.getSavedRuleSets();
    },
    clearEditor: function () {
      console.log('Clearing editor...', this.ide);
      this.ide.input = '';
      this.ide.output = null;
      this.selectedRuleSet.type = null;
      this.selectedRuleSet.index = null;
      this.compilerError = false;
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
  a.ruleset,
  a.nav-link {
    text-decoration: none;
  }
  a.nav-link {
    cursor: pointer;
  }
  a.nav-link.active {
    color: #17a2b8;
    background-color: aliceblue;
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
  .container-fluid.main {
    margin: auto;
    margin-top: 2rem;
    margin-bottom: 2rem;
    max-width: 1175px;
  }
</style>