<template>
  <div class="container-fluid p-0 h-100">
    <div class="row h-100">
      <div class="col">
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

const ruleChangeTypes = {
  CREATE: 'create',
  UPDATE: 'update',
  TRANSMUTE: 'transmute',
  DELETE: 'delete',
}

const ruleSetTypes = {
  SAVED: 'SAVED',
  CURRENT: 'CURRENT',
  QUEUED: 'QUEUED'
}

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
      ruleChangeTypes: ruleChangeTypes,
      ruleSetTypes: ruleSetTypes
    }
  },
  methods: {
    tryCompile: async function () {
      await this.$refs.ide.testRuleSet();
    },
    selectCurrentRule: function (index, ruleSetType) {
      if (ruleSetType !== ruleSetTypes.CURRENT) {
        return false;
      }

      this.$emit('current-selected', index);
    }
  }
}
</script>

<style scoped>

</style>