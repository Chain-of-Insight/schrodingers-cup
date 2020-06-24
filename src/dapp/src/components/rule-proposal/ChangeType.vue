<template>
  <div class="container h-100">
    <div class="row h-100 justify-content-center align-items-center">
      <div class="col-6">
        <p class="lead text-center mb-4">Choose the type of rule change you'd like to propose:</p>
        <div class="btn-group-toggle">
          <label
            class="btn btn-info btn-block p-3"
            :class="{ 'active': selectedType === type }"
            v-for="type of Object.values(proposalTypes)"
            :key="type"
          >
            <span class="oi oi-check float-left d-none"></span>
            <input
              type="radio"
              name="change-types"
              :id="type + '-rule'"
              :value="type"
              v-model="selectedType"
              @click="selectChangeType(type)"
            ><h4 class="m-0">{{ typeHeadings[type] }}</h4>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
// Constants
import proposalTypes from '../../constants/constants.js';

export default {
  props: {
    changeType: String,
    typeHeadings: Object
  },
  data: function () {
    return {
      proposalTypes: proposalTypes,
      selectedType: this.changeType,
    }
  },
  watch: {
    selectedType: function (value) {
      this.selectChangeType(value);
    }
  },
  methods: {
    selectChangeType: function (changeType) {
      this.$emit('select-type', changeType);
    }
  }
}
</script>

<style scoped>
  .active .oi-check {
    width: 0%;
    display: block !important;
    position: absolute;
    top: 35%;
  }

  .btn {
    position: relative;
  }
</style>