<template>
  <div class="container-fluid p-0 h-100">
    <div class="row h-100">
      <div class="col">
        <!-- <h2>{{ typeHeadings[changeType] }}</h2> -->
        <component :is="ideView"></component>
      </div>
    </div>
  </div>
</template>

<script>
import RuleSetList from '../ide/RuleSetList.vue';

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
      ideView: 'RuleSetList'
    }
  }
}
</script>

<style scoped>

</style>