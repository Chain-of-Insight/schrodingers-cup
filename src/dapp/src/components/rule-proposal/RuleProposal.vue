<template>
  <div>
    <!-- Modal -->
    <div class="modal fade" id="proposal-modal" data-backdrop="static" data-keyboard="false" tabindex="-1" role="dialog" aria-labelledby="proposal-modal-label" aria-hidden="true">
      <div class="modal-dialog modal-dialog-scrollable modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="proposal-modal-label">Time to propose a rule change!</h5>
          </div>
          <div class="modal-body">
            <transition name="slide" mode="out-in">
              <component
                :is="currentView"
                :change-type="changeType"
                :type-headings="typeHeadings"
                v-on:select-type="selectChangeType"
                v-on:go-back="currentView = 'ChangeType'"
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
              data-dismiss="modal"
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

import ChangeType from '../rule-proposal/ChangeType.vue';
import RuleSelect from '../rule-proposal/RuleSelect.vue';

export default {
  components: {
    ChangeType,
    RuleSelect
  },
  props: {
    rule: Object
  },
  data: () => ({
    currentView: 'ChangeType',
    changeType: null,
    typeHeadings: {
      [ruleChangeTypes.CREATE]: 'Create a New Rule',
      [ruleChangeTypes.UPDATE]: 'Update an Existing Rule',
      [ruleChangeTypes.TRANSMUTE]: 'Transmute a Rule',
      [ruleChangeTypes.DELETE]: 'Delete a Rule'
    }
  }),
  mounted: function () {
    $('#proposal-modal').on('hidden.bs.modal', this.resetModal.bind(this));
  },
  methods: {
    promptForProposal: function () {
      $('#proposal-modal').modal('show');
    },
    resetModal: function () {
      this.currentView = 'ChangeType';
      this.changeType = null;
    },
    selectChangeType: function (changeType) {
      const validType = Object.values(ruleChangeTypes).includes(changeType);
      if (!changeType || !validType) 
        return false;
      
      this.changeType = changeType
      this.currentView = 'RuleSelect';
    },
    selectRule: function  () {
      // this.currentView = 'RuleSelect';
    }
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