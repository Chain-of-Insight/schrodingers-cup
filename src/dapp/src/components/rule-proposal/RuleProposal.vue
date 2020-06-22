<template>
  <div>
    <!-- Modal -->
    <div class="modal fade" id="proposal-modal" data-backdrop="static" data-keyboard="false" tabindex="-1" role="dialog" aria-labelledby="proposal-modal-label" aria-hidden="true">
      <div class="modal-dialog modal-dialog-scrollable modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="proposal-modal-label">Time to propose a rule change!</h5>
            <Countdown
              :duration="turnWindow"
              v-on:ended="closeTurnWindow()"
              ref="timer"
            ></Countdown>
          </div>
          <div class="modal-body">
            <!-- Notifications -->
            <Notification 
              :type="alert.type" 
              :msg="alert.msg" 
              v-on:reset="alert = {type: null, msg: null}"
              :local="true"
            ></Notification>
            <transition name="slide" mode="out-in">
              <component
                :is="currentView"
                :change-type="changeType"
                :type-headings="typeHeadings"
                v-on:select-type="selectChangeType"
                v-on:go-back="currentView = 'ChangeType'"
                v-on:compiled="onCompiled"
                ref="proposal"
              ></component>
            </transition>
          </div>
          <div class="modal-footer d-block">
            <button
              v-if="currentView === 'RuleSelect'"
              type="button"
              class="btn btn-info float-left"
              @click="currentView = 'ChangeType'"
            >Back</button>
            <button
              v-if="currentView === 'RuleSelect'"
              type="button"
              class="btn btn-success float-right"
              @click="tryCompile()"
            >Submit</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
const $ = window.jQuery;

const ruleChangeTypes = {
  CREATE: 'create',
  UPDATE: 'update',
  TRANSMUTE: 'transmute',
  DELETE: 'delete',
}

const ruleTypes = {
  MUTABLE: 'mutable',
  IMMUTABLE: 'immutable'
}

import ChangeType from '../rule-proposal/ChangeType.vue';
import RuleSelect from '../rule-proposal/RuleSelect.vue';
import Notification from '../common/Notifications.vue';
import Countdown from '../common/Countdown.vue';

export default {
  components: {
    ChangeType,
    RuleSelect,
    Notification,
    Countdown
  },
  props: {
    rule: Object,
    turnWindow: {
      type: Number,
      required: true
    }
  },
  data: () => ({
    currentView: 'ChangeType',
    changeType: null,
    typeHeadings: {
      [ruleChangeTypes.CREATE]: 'Create a New Rule',
      [ruleChangeTypes.UPDATE]: 'Update an Existing Rule',
      [ruleChangeTypes.TRANSMUTE]: 'Transmute a Rule',
      [ruleChangeTypes.DELETE]: 'Delete a Rule'
    },
    proposalWindowClosed: false,
    alert: {
      type: null,
      msg: null
    },
  }),
  mounted: function () {
    // Start timer
    $('#proposal-modal').on('shown.bs.modal', this.startTimer.bind(this));
    $('#proposal-modal').on('hidden.bs.modal', this.resetModal.bind(this));
  },
  methods: {
    promptForProposal: function () {
      $('#proposal-modal').modal('show');
    },
    closeModal: function () {
      $('#proposal-modal').modal('hide');
    },
    selectChangeType: function (changeType) {
      const validType = Object.values(ruleChangeTypes).includes(changeType);
      if (!changeType || !validType) 
        return false;
      
      this.changeType = changeType
      this.currentView = 'RuleSelect';
    },
    tryCompile: async function () {
      if (this.currentView === 'RuleSelect') {
        await this.$refs.proposal.tryCompile();
      }
    },
    onCompiled: async function (successful, code, index) {
      if (!successful) {
        this.alert.type = 'danger';
        this.alert.msg = 'Fix your rule! Your rule must compile successfully before trying to propose it.';
        return;
      }

      let kind = ruleTypes.MUTABLE; // How to handle transmutation?
      this.$emit('rule-proposed', code, index, kind, this.changeType);
    },
    closeTurnWindow: function () {
        this.proposalWindowClosed = true;
        this.$emit('no-proposal');
        $('#proposal-modal').modal('hide');
      },
    startTimer: function () {
      this.$refs.timer.start();
    },
    resetModal: function () {
      this.currentView = 'ChangeType';
      this.changeType = null;
      this.$refs.timer.reset();
      this._retireNotification();
    },
    _retireNotification: function () {
      this.alert = {
        type: null,
        msg: null
      };
    },
  }
};
</script>

<style scoped>
  .slide-enter-active, .slide-leave-active {
    transition: all .3s ease;
  }
  .slide-enter, .slide-leave-to
  /* .slide-leave-active below version 2.1.8 */ {
    transform: translateX(10px);
    opacity: 0;
  }

  #proposal-modal .modal-dialog {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
  }

  #proposal-modal .modal-body {
    overflow-x: hidden;
  }
  
  #proposal-modal .modal-footer {
    min-height: 72px;
  }
</style>