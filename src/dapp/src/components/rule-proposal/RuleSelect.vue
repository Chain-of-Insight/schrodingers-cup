<template>
  <div class="container-fluid p-0 h-100">
    <div class="row h-100">
      <div class="col" v-if="changeType !== ruleChangeTypes.TRANSMUTE">
        <component
          :is="ideView"
          :active-lists="{
            [ruleSetTypes.QUEUED]: true
          }"
          :default-pane="ruleSetTypes.QUEUED"
        ></component>
      </div>
      <div class="col" v-else>
        <h2>{{ typeHeadings[changeType] }}</h2>
      </div>
    </div>
  </div>
</template>

<script>
import RuleSetList from '../ide/RuleSetList.vue';

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
    RuleSetList
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
  }
}
</script>

<style scoped>

</style>