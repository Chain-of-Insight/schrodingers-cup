<template>
  <!-- IDE Saved Rule Sets -->
  <div id="ide-saved" class="ide-pane h-100 d-flex flex-column">
    <nav class="nav nav-tabs">
      <a
        class="nav-link"
        :class="{ 'active': defaultPane === ruleSetTypes.CURRENT }"
        id="nav-current-tab"
        data-toggle="tab"
        href="#rules-current"
        role="tab"
        v-if="activeLists[ruleSetTypes.CURRENT] === true"
      >Current</a>
      <a
        class="nav-link"
        :class="{ 'active': defaultPane === ruleSetTypes.SAVED }"
        id="nav-saved-tab"
        data-toggle="tab"
        href="#rules-saved"
        role="tab"
        v-if="activeLists[ruleSetTypes.SAVED] === true"
      >Saved</a>
      <a
        class="nav-link"
        :class="{ 'active': defaultPane === ruleSetTypes.QUEUED }"
        id="nav-queued-tab"
        data-toggle="tab"
        href="#rules-queued"
        role="tab"
        v-if="activeLists[ruleSetTypes.QUEUED] === true"
      >Queued</a>
    </nav>
    <div class="tab-content border rounded-bottom border-top-0 flex-grow-1">
      <!-- CURRENT RULES -->
      <div
        id="rules-current"
        class="list-group list-group-flush tab-pane fade"
        :class="{
          'active': defaultPane === ruleSetTypes.CURRENT,
          'show': defaultPane === ruleSetTypes.CURRENT
        }"
        role="tab-panel"
        v-if="activeLists[ruleSetTypes.CURRENT] === true"
      >
        <p
          class="list-group-item border-0 text-muted"
          v-if="currentRules.length === 0"
        >No current rule sets loaded...</p>
        <a
          class="ruleset btn btn-outline-success list-group-item list-group-item-action"
          role="button"
          tabindex="0"
          v-bind:class="{
            'active':
              selectedRule.type === ruleSetTypes.CURRENT &&
              selectedRule.index === index
          }"
          v-bind:key="index"
          v-for="(ruleSet, index) in currentRules"
          v-on:click="selectRule(index, ruleSetTypes.CURRENT)"
          style="cursor:pointer;"
        >
          <span>{{index + 1}}. {{ ruleSet.name }}</span>
        </a>
      </div>
      <!-- PLAYER-SAVED RULES -->
      <div
        id="rules-saved"
        class="list-group list-group-flush tab-pane fade"
        role="tab-panel"
        :class="{
          'active': defaultPane === ruleSetTypes.SAVED,
          'show': defaultPane === ruleSetTypes.SAVED
        }"
        v-if="activeLists[ruleSetTypes.SAVED] === true"
      >
        <p
          class="list-group-item border-0 text-muted"
          v-if="savedRules.length === 0"
        >No saved rule sets...</p>
        <p
          class="list-group-item border-0 text-muted"
          v-if="queuedRules.length > 0 && unQueuedRules.length === 0"
        >All of your saved rule sets are queued!</p>
        <a
          class="ruleset btn btn-outline-success list-group-item list-group-item-action d-flex justify-content-between align-items-center"
          role="button"
          tabindex="0"
          v-bind:class="{
            'active':
              selectedRule.type === ruleSetTypes.SAVED &&
              selectedRule.index === savedIndex
          }"
          style="cursor:pointer;"
          v-bind:key="index"
          v-for="(savedIndex, index) in unQueuedRules"
          v-on:click="selectRule(savedIndex, ruleSetTypes.SAVED)"
        >
          <span>{{index + 1}}. {{ savedRules[savedIndex].name }}</span>
          <button
            class="btn btn-outline-warning"
            @click="queueRule(savedIndex)"
          ><small>Queue</small></button>
        </a>
      </div>
      <!-- QUEUED RULES -->
      <div
        id="rules-queued"
        class="list-group list-group-flush tab-pane fade"
        role="tab-panel"
        :class="{
          'active': defaultPane === ruleSetTypes.QUEUED,
          'show': defaultPane === ruleSetTypes.QUEUED
        }"
        v-if="activeLists[ruleSetTypes.QUEUED] === true"
      >
        <p
          class="list-group-item border-0 text-muted"
          v-if="queuedRules.length === 0"
        >No queued rule sets...</p>
        <a
          class="ruleset btn btn-outline-success list-group-item list-group-item-action d-flex justify-content-between align-items-center"
          role="button"
          tabindex="0"
          v-bind:class="{
            'active':
              selectedRule.type === ruleSetTypes.QUEUED &&
              selectedRule.index === savedIndex
          }"
          style="cursor:pointer;"
          v-bind:key="index"
          v-for="(savedIndex, index) in queuedRules"
          v-on:click="selectRule(savedIndex, ruleSetTypes.QUEUED)"
        >
          <span>{{index + 1}}. {{ savedRules[savedIndex].name }}</span>
          <button
            class="btn btn-outline-warning"
            @click="unQueueRule(savedIndex, index)"
            style="word-break: keep-all;"
          ><small>Un&#8209;queue</small></button>
        </a>
      </div>
    </div>
  </div>
</template>

<script>
  const ruleSetTypes = {
    SAVED: 'SAVED',
    CURRENT: 'CURRENT',
    QUEUED: 'QUEUED'
  }

  export default {
    props: {
      defaultPane: {
        type: String,
        default: ruleSetTypes.SAVED
      },
      currentRules: {
        type: Array,
        default: () => []
      },
      savedRules: {
        type: Array,
        default: () => []
      },
      unQueuedRules: {
        type: Array,
        default: () => []
      },
      queuedRules: {
        type: Array,
        default: () => []
      },
      activeLists: {
        type: Object,
        default: () => ({
          [ruleSetTypes.CURRENT]: true,
          [ruleSetTypes.SAVED]: true,
          [ruleSetTypes.QUEUED]: true,
        })
      }
    },
    data: function () {
      return {
        // defaultPane: ruleSetTypes.CURRENT,
        ruleSetTypes: ruleSetTypes,
        selectedRule: {
          type: null,
          index: null
        }
      }
    },
    methods: {
      selectRule: function (savedIndex, ruleType) {
        if (
          typeof(savedIndex) !== 'number' ||
          !Object.values(ruleSetTypes).includes(ruleType)
        ) {
          return false;
        }

        this.selectedRule.type = ruleType;
        this.selectedRule.index = savedIndex;

        this.$emit('select-rule', this.selectedRule);
      },
      queueRule: async function (savedIndex) {
        if (
          savedIndex !== this.selectedRule.index ||
          this.selectedRule.type !== ruleSetTypes.SAVED
        ) {
          this.selectRule(savedIndex, ruleSetTypes.SAVED);
        }

        this.$emit('queue-rule');
      },
      unQueueRule: async function (savedIndex, queuedIndex) {
        if (
          savedIndex !== this.selectedRule.index ||
          this.selectedRule.type !== ruleSetTypes.QUEUED
        ) {
          this.selectRule(savedIndex, ruleSetTypes.QUEUED);
        }

        this.$emit('unqueue-rule', queuedIndex);
      }
    }
  }
</script>

<style scoped>
.nav-link {
  margin-bottom: -1px;
  border-bottom: 0;
}

</style>