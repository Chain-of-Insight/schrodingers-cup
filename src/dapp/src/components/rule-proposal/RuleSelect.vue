<template>
  <div class="container-fluid p-0 h-100">
    <div class="row h-100">
      <div class="col">
        <p class="lead font-weight-bold">Choose the current rule you're proposing to change:</p>
        <Practice
          :rule-proposal="true"
          :current-only="true"
          :disable-editing="true"
          ref="ide"
          v-on:rule-selected="selectCurrentRule"
        ></Practice>
      </div>
    </div>
  </div>
</template>

<script>
import RuleSetList from '../ide/RuleSetList.vue';
import Practice from '../practice/Practice.vue';

import {
  proposalTypes,
  ruleSetTypes
} from '../../constants/constants.js';

export default {
  components: {
    RuleSetList,
    Practice
  },
  props: {
    changeType: String,
    typeHeadings: Object
  },
  data: function () {
    return {
      ideView: 'RuleSetList',
      proposalTypes: proposalTypes,
      ruleSetTypes: ruleSetTypes
    }
  },
  methods: {
    tryCompile: async function () {
      await this.$refs.ide.testRuleSet();
    },
    selectCurrentRule: function (index, ruleSetType) {
      if (
        ruleSetType !== ruleSetTypes.IMMUTABLE &&
        ruleSetType !== ruleSetTypes.MUTABLE
      ) {
        return false;
      }

      this.$emit('current-selected', index, ruleSetType);
    }
  }
}
</script>

<style scoped>

</style>