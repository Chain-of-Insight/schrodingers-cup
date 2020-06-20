<template>
  <div class="container h-100">
    <div class="row h-100 justify-content-center align-items-center">
      <div class="col-6">
        <p class="lead text-center mb-4">Choose the type of rule change you'd like to propose:</p>
        <div class="btn-group-toggle">
          <label
            class="btn btn-info btn-block p-3"
            :class="{ 'active': selectedType === changeType }"
            v-for="changeType of Object.values(ruleChangeTypes)"
            :key="changeType"
          >
            <span class="oi oi-check float-left d-none"></span>
            <input
              type="radio"
              name="change-types"
              :id="changeType + '-rule'"
              :value="changeType"
              v-model="selectedType"
              @click="selectChangeType(changeType)"
            ><h4 class="m-0">{{ typeHeadings[changeType] }}</h4>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

const ruleChangeTypes = {
  CREATE: 'create',
  UPDATE: 'update',
  TRANSMUTE: 'transmute',
  DELETE: 'delete',
}

export default {
  props: {
    changeType: String,
    typeHeadings: Object
  },
  data: function () {
    return {
      ruleChangeTypes: ruleChangeTypes,
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