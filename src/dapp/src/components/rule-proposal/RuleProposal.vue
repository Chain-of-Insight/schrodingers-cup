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
                :rule-proposal="true"
                :queued-only="currentView === 'Practice' ? true : false"
                v-on:rule-selected="selectQueuedRule"
                v-on:current-selected="selectCurrentRule"
              ></component>
            </transition>
          </div>
          <div class="modal-footer d-block">
            <button
              v-if="currentView !== 'ChangeType'"
              type="button"
              class="btn btn-info float-left"
              @click="goBack()"
            >Back</button>
            <button
              v-if="submitEnabled || nextEnabled"
              type="button"
              class="btn float-right"
              :class="{
                'btn-success': submitEnabled,
                'btn-primary': nextEnabled
              }"
              @click="submitRule()"
            >{{ submitEnabled ? 'Submit' : nextEnabled ? 'Next' : null }}</button>
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

const ruleSetTypes = {
  SAVED: 'SAVED',
  CURRENT: 'CURRENT',
  QUEUED: 'QUEUED'
}

import Practice from '../practice/Practice.vue';
import ChangeType from '../rule-proposal/ChangeType.vue';
import RuleSelect from '../rule-proposal/RuleSelect.vue';
import Notification from '../common/Notifications.vue';
import Countdown from '../common/Countdown.vue';

export default {
  components: {
    ChangeType,
    RuleSelect,
    Notification,
    Countdown,
    Practice
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
    ruleCandidate: {
      code: null,
      index: null,
      kind: 'mutable'
    },
    queuedIndex: 0,
  }),
  computed: {
    submitEnabled: function () {
      if (this.changeType === ruleChangeTypes.CREATE) {
        return this.currentView === 'Practice' ? true : false;
      } else {
        return this.currentView === 'RuleSelect' ? true : false
      }

      return false;
    },
    nextEnabled: function () {
      if (
        (
          this.changeType === ruleChangeTypes.UPDATE ||
          this.changeType === ruleChangeTypes.TRANSMUTE
        ) && this.currentView === 'Practice'
      ) {
        return true;
      }

      return false;
    }
  },
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
      
      if (this.changeType === ruleChangeTypes.DELETE) {
        this.currentView = 'RuleSelect';
      } else {
        this.currentView = 'Practice';
      }
    },
    selectCurrentRule: function (currentIndex) {
      this.ruleCandidate.index = currentIndex;
    },
    selectQueuedRule: function (queuedIndex, ruleSetType) {
      if (ruleSetType !== ruleSetTypes.QUEUED)
        return false;

      this.queuedIndex = queuedIndex;
    },
    submitRule: async function () {
      if (!this.ruleCandidate.code && this.changeType !== ruleChangeTypes.DELETE) {
        await this.$refs.proposal.testRuleSet();
      } else {
        this.proposeRule();
      }
    },
    onCompiled: async function (successful, code) {
      if (!successful) {
        this.alert.type = 'danger';
        this.alert.msg = 'Fix your rule! Your rule must compile successfully before trying to propose it.';
        setTimeout(() => {
          this._retireNotification();
        }, 5000);
        return;
      }

      this.ruleCandidate.code = code;

      if (this.nextEnabled) {
        this.currentView = 'RuleSelect';
        this._retireNotification();
      } else if (this.submitEnabled) {
        this.proposeRule();
      }
    },
    proposeRule () {
      if (this.changeType === ruleChangeTypes.CREATE) {
        this.ruleCandidate.index = -1;
      } else if (typeof(this.ruleCandidate.index) !== 'number') {
        this.alert.type = 'danger';
        this.alert.msg = 'You must select a current rule to change.';
        setTimeout(() => {
          this._retireNotification();
        }, 5000);
        return;
      }

      if (this.changeType === ruleChangeTypes.DELETE) {
        this.ruleCandidate.code = "";
      }

      this.$emit(
        'rule-proposed',
        this.ruleCandidate.code,
        this.ruleCandidate.index,
        this.ruleCandidate.kind,
        this.changeType
      );
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
      this.ruleCandidate.code = null;
      this.ruleCandidate.index = null;
      this.ruleCandidate.kind = ruleTypes.MUTABLE;
      this.$refs.timer.reset();
      this._retireNotification();
    },
    goBack: function () {
      // For now:
      this.ruleCandidate.code = null;
      this.ruleCandidate.index = null;

      if (this.changeType !== ruleChangeTypes.DELETE && this.currentView === 'RuleSelect') {
        this.currentView = 'Practice';
        return;
      }

      this.currentView = 'ChangeType';
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