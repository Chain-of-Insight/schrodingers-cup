<template>
  <div class="container-fluid p-0 h-100">
    <div class="row h-100">
      <div class="col" v-if="changeType === ruleChangeTypes.TRANSMUTE || changeType === ruleChangeTypes.DELETE">
        <h2>{{ typeHeadings[changeType] }}</h2>
      </div>
      <div class="col" v-else>
        <!-- <component
          :is="ideView"
          :active-lists="{
            [ruleSetTypes.QUEUED]: true
          }"
          :default-pane="ruleSetTypes.QUEUED"
        ></component> -->
        <Practice :rule-proposal="true" ref="ide" v-on:compiled="onCompiled"></Practice>
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
    onCompiled: function (successful, code, index) {
      this.$emit('compiled', successful, code, index);
    },
    tryCompile: async function () {
      await this.$refs.ide.testRuleSet();
    }
  }
}
</script>

<style scoped>

</style>