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
      <div
        id="rules-current"
        class="list-group tab-pane fade show active"
        role="tab-panel"
      >
        <p
          class="list-group-item border-0 text-muted"
          v-if="gameRuleSets.length === 0"
        >No current rule sets loaded...</p>
        <!-- Current Rule Sets -->
        <a
          class="ruleset btn btn-outline-success list-group-item list-group-item-action"
          role="button"
          tabindex="0"
          v-bind:class="{
            'active':
              selectedRuleSet.type === ruleSetTypes.CURRENT &&
              selectedRuleSet.index === index
          }"
          v-bind:key="index"
          v-for="(ruleSet, index) in gameRuleSets"
          v-on:click="loadRuleSet(index, ruleSetTypes.CURRENT)"
          style="cursor:pointer;"
        >
          <span>{{index + 1}}. {{ ruleSet.name }}</span>
        </a>
      </div>
      <div
        id="rules-saved"
        class="list-group tab-pane fade"
        role="tab-panel"
      >
        <p
          class="list-group-item border-0 text-muted"
          v-if="savedRuleSets.length === 0"
        >No saved rule sets...</p>
        <p
          class="list-group-item border-0"
          v-if="queuedRuleSets.length > 0 && unQueuedRuleSets.length === 0"
        >All of your saved rule sets are queued!</p>
        <!-- Saved Rule Sets -->
        <div
          class="ruleset btn-group"
          role="group"
          v-bind:key="index"
          v-for="(savedIndex, index) in unQueuedRuleSets"
          v-on:click="loadRuleSet(savedIndex, ruleSetTypes.SAVED)"
        >
          <a
            class="btn btn-outline-success list-group-item list-group-item-action"
            role="button"
            tabindex="0"
            v-bind:class="{
              'active':
                selectedRuleSet.type === ruleSetTypes.SAVED &&
                selectedRuleSet.index === savedIndex
            }"
            style="cursor:pointer;"
          >
            <span>{{index + 1}}. {{ savedRuleSets[savedIndex].name }}</span>
          </a>
          <button
            class="btn btn-outline-warning"
            @click="queueRuleSet(savedIndex)"
          ><small>Queue</small></button>
        </div>
      </div>
      <div
        id="rules-queued"
        class="list-group tab-pane fade"
        role="tab-panel"
      >
        <!-- Queued Rule Sets -->
        <p
          class="list-group-item border-0 text-muted"
          v-if="queuedRuleSets.length === 0"
        >No queued rule sets...</p>
        <div
          class="ruleset btn-group"
          role="group"
          v-bind:key="index"
          v-for="(savedIndex, index) in queuedRuleSets"
          v-on:click="loadRuleSet(savedIndex, ruleSetTypes.QUEUED)"
        >
          <a
            class="btn btn-outline-success list-group-item list-group-item-action"
            role="button"
            tabindex="0"
            v-bind:class="{
              'active':
                selectedRuleSet.type === ruleSetTypes.QUEUED &&
                selectedRuleSet.index === savedIndex
            }"
            style="cursor:pointer;"
          >
            <span>{{index + 1}}. {{ savedRuleSets[savedIndex].name }}</span>
          </a>
          <button
            class="btn btn-outline-warning"
            @click="unQueueRuleSet(savedIndex, index)"
            style="word-break: keep-all;"
          ><small>Un&#8209;queue</small></button>
        </div>
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
    data: function () {
      return {
        currentPane: ruleSetTypes.CURRENT,
        ruleSetTypes: ruleSetTypes,
        savedRuleSets: [],
        gameRuleSets: [],
        selectedRuleSet: {
          type: null,
          index: null
        }
      }
    },
    computed: {
      unQueuedRuleSets: function () {
        return this.savedRuleSets;
      },
      queuedRuleSets: function () {
        return this.savedRuleSets;
      }
    },
    methods: {
      loadRuleSet: function (index, type) {
        // ...
      },
      queueRuleSet: async function (index) {
        // ...
      },
      unQueueRuleSet: async function (savedIndex, queuedIndex) {
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