<template>
  <!-- IDE Saved Rule Sets -->
  <div id="ide-saved" class="ide-pane h-100 d-flex flex-column">
    <!-- <nav class="nav nav-tabs">
      <a
        class="nav-link"
        @click="currentPane = ruleSetTypes.CURRENT"
        v-bind:class="{ active: currentPane === ruleSetTypes.CURRENT }"
      >Current</a>
      <a
        class="nav-link"
        @click="currentPane = ruleSetTypes.SAVED"
        v-bind:class="{ active: currentPane === ruleSetTypes.SAVED }"
      >Saved</a>
      <a
        class="nav-link"
        @click="currentPane = ruleSetTypes.QUEUED"
        v-bind:class="{ active: currentPane === ruleSetTypes.QUEUED }"
      >Queued</a>
    </nav> -->
    <nav class="nav nav-tabs">
      <a
        class="nav-link active"
        id="nav-current-tab" data-toggle="tab" href="#rules-current" role="tab"
      >Current</a>
      <a
        class="nav-link"
        id="nav-saved-tab" data-toggle="tab" href="#rules-saved" role="tab"
      >Saved</a>
      <a
        class="nav-link"
        id="nav-queued-tab" data-toggle="tab" href="#rules-queued" role="tab"
      >Queued</a>
    </nav>
    <div class="tab-content border rounded-bottom border-top-0 flex-grow-1">
      <!-- CURRENT RULES -->
      <div
        id="rules-current"
        class="list-group list-group-flush tab-pane fade show active"
        role="tab-panel"
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
    SAVED: 'ACTIVE',
    CURRENT: 'CURRENT',
    QUEUED: 'QUEUED'
  }

  export default {
    props: {
      currentRules: {
        // required: true,
        type: Array,
        default: () => []
      },
      savedRules: {
        // required: true,
        type: Array,
        default: () => []
      }
    },
    data: function () {
      return {
        currentPane: ruleSetTypes.CURRENT,
        ruleSetTypes: ruleSetTypes,
        selectedRule: {
          type: null,
          index: null
        }
      }
    },
    computed: {
      unQueuedRules: function () {
      const ruleSetEntries = Array.from(this.savedRules.entries());
      return ruleSetEntries
        .filter(([index, ruleSet]) => !ruleSet.hasOwnProperty('queued'))
        .sort(([indexA, ruleSetA], [indexB, ruleSetB]) => ruleSetA.queued - ruleSetB.queued)
        .map(([index, ruleSet]) => index);
    },
    queuedRules: function () {
      const ruleSetEntries = Array.from(this.savedRules.entries());
      return ruleSetEntries
        .filter(([index, ruleSet]) => ruleSet.hasOwnProperty('queued'))
        .sort(([indexA, ruleSetA], [indexB, ruleSetB]) => ruleSetA.queued - ruleSetB.queued)
        .map(([index, ruleSet]) => index);
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

        this.selectedRule = {
          type: ruleType,
          index: savedIndex
        }

        this.$emit('select-rule', this.selectedRule);
      },
      loadRule: function (index, type) {
        // ...
      },
      queueRule: async function (index) {
        // ...
      },
      unQueueRule: async function (savedIndex, queuedIndex) {
        // ...
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