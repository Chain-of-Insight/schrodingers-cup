<template>

  <div>
    <!-- Notifications -->
    <Notification 
      :type="alert.type" 
      :msg="alert.msg" 
      v-on:reset="alert = {type: null, msg: null}"
    ></Notification>

    <div class="container-xl main">
      <div class="row">
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
              <!-- IDE Saved Rulesets -->
              <div id="ide-saved" class="ide-pane col">
                <label>
                  <strong>Saved Rulesets:</strong>
                </label>
                <ul class="list-group">
                  <li
                    href="#"
                    class="list-group-item list-group-item-action d-flex justify-content-between align-items-center"
                    v-for="ruleSet in ide.savedRuleSets"
                    :key="ruleSet.name"
                  >{{ ruleSet.name }}.nom
                    <div class="btn-group">
                      <button class="btn btn-primary" @click="loadRuleSet(ruleSet.name)">Load</button>
                      <button
                        class="btn btn-danger"
                        data-toggle="modal"
                        data-target="#delete-modal"
                        @click="ide.selectedRuleSet = ruleSet.name;"
                      >
                        <span class="oi oi-trash" title="Delete ruleset" aria-hidden="true"></span>
                      </button>
                    </div>
                  </li>
                </ul>
              </div>
              <!-- IDE Input -->
              <div id="ide-input" class="ide-pane col">
                <label>
                  <strong>Input New Rule:</strong>
                </label>
                <codemirror 
                  v-model="ide.input"
                  :options="ide.options"
                ></codemirror>
                <div class="execute">
                  <button class="btn btn-primary" @click="testRuleSet()">Compile</button>
                  <button class="btn btn-success" data-toggle="modal" data-target="#save-modal">Save</button>
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
          <form v-on:submit.prevent="saveRuleSet()">
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
                  <div class="input-group-append">
                    <span class="input-group-text">.nom</span>
                  </div>
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
              <button v-else type="submit" class="btn btn-success">Save</button>
            </div>
          </form>
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
      ruleSetName: null,
      selectedRuleSet: null,
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
      
      // Fetch compiler result\
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
        } else if (result.status == 200) {
          // console.log('Compiled successfully,', result.status);
          this.alert.type = 'success';
          this.alert.msg = 'Compiled successfully';
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
        this.ide.output = "Error: compiler failed for unknown reasons 😅";
      }
    },
    getSavedRuleSets: function () {
      this.ide.savedRuleSets = Object.entries(localStorage)
        .map(([name, data]) => {
          const ruleSet = JSON.parse(data);
          ruleSet.name = name;
          return ruleSet;
        })
        .sort((a, b) => a.sortOrder - b.sortOrder);
    },
    saveRuleSet: async function () {
      // Check if ruleset name already exists
      if (localStorage.hasOwnProperty(this.ide.ruleSetName)) {
        this.ide.nameError = true;
      } else {
        this.ide.state.loading = true;
        const ruleSet = {
          code: this.ide.input,
          sortOrder: localStorage.length
        }

        console.log('Ruleset to be saved:', ruleSet);
        localStorage.setItem(this.ide.ruleSetName, JSON.stringify(ruleSet));
        this.getSavedRuleSets();

        this.alert.type = 'success';
        this.alert.msg = `Ruleset '${this.ide.ruleSetName}' saved!`;

        $('#save-modal').modal('hide');
        this.ide.state.loading = false;
        this.ide.ruleSetName = '';
      }
    },
    loadRuleSet: function (name) {
      console.log('Loading ruleset:', name);
      const loadedRuleSet = JSON.parse(localStorage.getItem(name));
      this.ide.input = loadedRuleSet.code;

      this.alert.type = 'success';
      this.alert.msg = `Loaded ruleset '${name}'.`;
    },
    deleteRuleSet: function () {
      console.log('Ruleset to be deleted:', this.ide.selectedRuleSet);
      localStorage.removeItem(this.ide.selectedRuleSet);
      this.getSavedRuleSets();

      this.alert.type = 'info';
      this.alert.msg = `Ruleset '${this.ide.selectedRuleSet}' deleted.`;

      $('#delete-modal').modal('hide');
      this.ide.selectedRuleSet = ''
    },
    cancelSave: function () {
      // Quick and dirty solution to handle modal close (jQuery doesn't play nice with vue).
      this.ide.ruleSetName = '';
      this.ide.nameError = false;
    },
    cancelDelete: function () {
      // Ditto.
      this.ide.selectedRuleSet = ''
    },
    clearEditor: function () {
      console.log('Clearing editor...', this.ide);
      this.ide.input = '';
    },
    clearEditorOutput: function () {
      console.log('Clearing output...', this.ide);
      this.ide.output = null;
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
  }
</style>