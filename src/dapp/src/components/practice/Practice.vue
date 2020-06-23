<template>

  <div :class="{ 'h-100': !activeGame || !ruleProposal }">
    <!-- Notifications -->
    <Notification 
      :type="alert.type" 
      :msg="alert.msg" 
      v-on:reset="alert = {type: null, msg: null}"
    ></Notification>

    <div :class="{
      'p-0': activeGame || ruleProposal,
      'py-4': !activeGame && !ruleProposal
    }" class="container h-100 d-flex flex-column">

      <!-- Heading -->
      <div class="row" v-if="!activeGame && !ruleProposal">
        <div class="col">
          <h1>{{ title }}</h1>
          <h5 class="mb-4">{{ subtitle }}</h5>
        </div>
      </div>

      <!-- Not Connected -->
      <div class="row" v-if="!connected">
        <div class="col">
          <ul class="list-unstyled">
            <li @click="connectUser()">
              <button class="btn btn-primary btn-connect">Login With Tezos</button>
            </li>
          </ul>
          <p>Connect your Tezos wallet to get started</p>
        </div>
      </div>

      <!-- Connected -->
      <template v-else>
        <div class="row flex-shrink-1 flex-grow-1 overflow-hidden">

          <!-- Rule Lists -->
          <div
            class="d-flex flex-column mh-100"
            :class="{
              'col-4': !noEditor,
              'col': noEditor
            }"
          >
            <div class="row">
              <div class="col">
                <label>
                  <strong>Rules:</strong>
                </label>
              </div>
            </div>
            <div class="row flex-grow-1 flex-shrink-1 overflow-hidden h-100">
              <div class="col mh-100">
                <RuleSetList
                  :loaded-rule="selectedRule"
                  :current-rules="ruleSetLists.current"
                  :saved-rules="ruleSetLists.saved"
                  :un-queued-rules="unQueuedRules"
                  :queued-rules="queuedRules"
                  v-on:select-rule="loadRule"
                  v-on:queue-rule="queueRule"
                  v-on:unqueue-rule="unQueueRule"
                  :queued-only="queuedOnly"
                  :current-only="currentOnly"
                ></RuleSetList>
              </div>
            </div>
          </div>

          <!-- IDE Editor -->
          <div
            class="d-flex flex-column mh-100"
            :class="{
              'col-8': !noEditor,
              'col': noEditor
            }"
            v-if="!noEditor"
          >
            <div class="row">
              <div class="col">
                <label>
                  <strong>Rule Editor:</strong>
                </label>
              </div>
            </div>
            <div class="row flex-shrink-1 overflow-hidden">
              <div id="ide-input" class="ide-pane col h-100">
                <div class="row">
                  <div class="col">
                    <codemirror
                      v-model="ide.input"
                      :options="ide.options"
                      @input="clearEditorOutput()"
                      @change="inputHandler()"
                    ></codemirror>
                  </div>
                </div>
              </div>
            </div>

            <div class="row">
              <div class="col pt-2">
                <!-- Compile Nomic -->
                <button class="btn btn-primary" @click="testRuleSet()" v-if="!ruleProposal">Compile</button>
                
                <!-- Save Rule Set -->
                <button 
                  class="btn btn-success" 
                  @click="saveRuleSetHandler()"
                  :disabled="!ide.output || compilerError"
                  v-if="!ruleProposal"
                >Save</button>

                <!-- Clear Editor -->
                <button class="btn btn-danger" @click="clearEditor()">
                  <span class="oi oi-trash" title="Clear editor" aria-hidden="true"></span>
                </button>
              </div>
            </div>
          </div>

        </div>

        <div
          class="row"
          style="max-height: 40%"
          v-if="!noEditor"
        >
          <div class="col pt-3 mh-100 overflow-hidden d-flex flex-column">
            <div class="row">
              <div class="col">
                <label>
                  <strong>Output:</strong>
                </label>
              </div>
            </div>
            <div class="row flex-shrink-1 flex-grow-1 overflow-hidden h-100" v-if="ide.output">
              <!-- IDE Output -->
              <div id="ide-output" class="ide-pane col mh-100 overflow-hidden d-flex flex-column h-100">
                <div class="row">
                  <div class="col">
                    <div class="executed clear-output">
                      <span class="clear" @click="clearEditorOutput()">Clear Output</span>
                    </div>
                  </div>
                </div>
                <div class="row pt-2 flex-grow-1 flex-shrink-1 overflow-hidden">
                  <div class="col mh-100">
                    <div class="term-container overflow-auto mh-100" v-html="ide.output"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
          
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
            <p><strong>Are you sure you want to delete ruleset '{{ this.ide.selectedRule }}'?</strong></p>
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
import RuleSetList from '../ide/RuleSetList.vue'

// IDE
//import 'codemirror/mode/lua/lua';
import 'codemirror/mode/shell/shell';
import 'codemirror/theme/dracula.css';

const $ = window.jQuery;

const ruleSetTypes = {
  SAVED: 'SAVED',
  CURRENT: 'CURRENT',
  QUEUED: 'QUEUED'
}

const CURRENT_RULES = require('./rules/currentRules.json');

export default {
  components: {
    Notification,
    RuleSetList
  },
  props: {
    activeGame: {
      default: false,
      type: Boolean
    },
    ruleProposal: {
      default: false,
      type: Boolean
    },
    queuedOnly: {
      type: Boolean,
      default: false
    },
    currentOnly: {
      type: Boolean,
      default: false
    },
    noEditor: {
      type: Boolean,
      default: false
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
    ruleSetLists: {
      current: [],
      saved: [],
    },
    ide: {
      input: '',
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
        lineWrapping: true,
        line: true
      },
      state: {
        loading: false
      },
      execute: testNomic
    },
    compilerError: false,
    selectedRule: {
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
    unQueuedRules: function () {
      const ruleSetEntries = Array.from(this.ruleSetLists.saved.entries());
      return ruleSetEntries
        .filter(([index, ruleSet]) => !ruleSet.hasOwnProperty('queued'))
        .sort(([indexA, ruleSetA], [indexB, ruleSetB]) => ruleSetA.queued - ruleSetB.queued)
        .map(([index, ruleSet]) => index);
    },
    queuedRules: function () {
      const ruleSetEntries = Array.from(this.ruleSetLists.saved.entries());
      return ruleSetEntries
        .filter(([index, ruleSet]) => ruleSet.hasOwnProperty('queued'))
        .sort(([indexA, ruleSetA], [indexB, ruleSetB]) => ruleSetA.queued - ruleSetB.queued)
        .map(([index, ruleSet]) => index);
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
          // Emit unsuccessful flag for rule proposal
          this.$emit('compiled', false);
          setTimeout(() => {
            this._retireNotification();
          }, 5000);
        } else if (result.status == 200) {
          // console.log('Compiled successfully,', result.status);
          this.alert.type = 'success';
          this.alert.msg = 'Compiled successfully';
          this.compilerError = false;
          // Emit success flag and data to submit for proposal
          this.$emit(
            'compiled',
            true,
            this.ide.input,
            this.selectedRule.index
          );
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
      this.ruleSetLists.current = CURRENT_RULES;
    },
    getSavedRuleSets: async function (index = false) {
      // Load rule sets if exists
      let storedRuleSets = await localStorage.getItem('ruleSets');
      if (storedRuleSets) {
        storedRuleSets = await JSON.parse(storedRuleSets);
        this.ruleSetLists.saved = storedRuleSets;
      }
    },
    saveRuleSetHandler: function () {
      switch (this.selectedRule.type) {
        case ruleSetTypes.QUEUED:
        case ruleSetTypes.SAVED:
          this.saveRuleSet(this.selectedRule.index);
          break;
        case ruleSetTypes.CURRENT:
        case typeof(this.selectedRule.index) !== 'number':
        default:
          $('#save-modal').modal('show');
          break;
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
        ruleSets[index].code = this.ide.input;
        ruleSet = ruleSets[index];
      } else {
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
        if (this.ruleSetLists.saved) {
          if (this.ruleSetLists.saved.length) {
            index = this.ruleSetLists.saved.length;
          } else {
            index = 0;
          }
        } else {
          index = 0;
        }
      }

      await this.getSavedRuleSets();
      this.loadRule(index, ruleSetTypes.SAVED);
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
    loadRule: function (savedIndex, ruleType) {
      let ruleSetList = null;

      switch (ruleType) {
        case ruleSetTypes.CURRENT:
          ruleSetList = this.ruleSetLists.current;
          break;
        case ruleSetTypes.SAVED || ruleSetTypes.QUEUED:
          ruleSetList = this.ruleSetLists.saved;
          break;
        default:
          ruleSetList = this.ruleSetLists.current;
          break;
      }

      let ruleSet = ruleSetList[savedIndex];

      if (!ruleSet) {
        return;
      }

      console.log('Loading rule set =>', [ruleSet, savedIndex]);

      if (ruleSet.hasOwnProperty('code')) {
        try {
          // Set IDE state
          this.ide.input = ruleSet.code;
          // Set UI state
          this.selectedRule.type = ruleType;
          this.selectedRule.index = Number(savedIndex);
          // Force update cycle
          this.$forceUpdate();
        } catch(e) {
          console.log('Error loading ruleset =>', [e, ruleSet])
          this.alert.type = 'danger';
          this.alert.msg = `Error loading rule set '${JSON.stringify(ruleSet)}'`;
          setTimeout(() => {
            this._retireNotification();
          }, 5000);
        }
      }
    },
    updateRuleSets: async function (index, data) {
      // Update rule sets in localStorage and refresh
      await localStorage.setItem('ruleSets', JSON.stringify(this.ruleSetLists.saved));
      await this.getSavedRuleSets();
    },
    queueRule: async function (savedIndex) {
      // console.log('Queueing rule!', this.selectedRule);
      // const savedIndex = this.selectedRule.index
      if (typeof(savedIndex) !== 'number')
        return false;
      const rule = this.ruleSetLists.saved[savedIndex];

      rule.queued = this.queuedRules.length + 1;
      await this.updateRuleSets();

      this.alert.type = 'info';
      this.alert.msg = `Queued ruleset '${rule.name}'`;
      setTimeout(() => {
        this._retireNotification();
      }, 5000);
    },
    unQueueRule: async function (savedIndex, queuedIndex) {
      // console.log('Un-queueing rule!', this.selectedRule);
      // const savedIndex = this.queuedRules[queuedIndex];
      if (typeof(savedIndex) !== 'number')
        return false;
      const rule = this.ruleSetLists.saved[savedIndex];
      
      // Remove queued order/flag
      delete rule.queued;
      // Shift other queued rule sets' orders
      const toBeShifted = this.queuedRules.slice(queuedIndex + 1);
      toBeShifted.forEach(savedIndex => {
        this.ruleSetLists.saved[savedIndex].queued -= 1;
      });
      await this.updateRuleSets();

      this.alert.type = 'info';
      this.alert.msg = `Un-queued rule set '${rule.name}'`;
      setTimeout(() => {
        this._retireNotification();
      }, 5000);
    },
    clearEditor: function () {
      console.log('Clearing editor...', this.ide);
      this.ide.input = '';
      this.ide.output = null;
      this.selectedRule.type = null;
      this.selectedRule.index = null;
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

  .container {
    min-height: 100%;
  }
</style>