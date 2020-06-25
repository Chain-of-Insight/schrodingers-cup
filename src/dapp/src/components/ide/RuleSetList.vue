<template>
  <!-- IDE Saved Rule Sets -->
  <div id="ide-saved" class="ide-pane h-100 d-flex flex-column">
    <nav v-if="!currentOnly" class="nav nav-tabs">
      <a
        class="nav-link active"
        id="nav-current-tab"
        data-toggle="tab"
        href="#rules-current"
        role="tab"
        v-if="!queuedOnly"
      >Current</a>
      <a
        class="nav-link"
        id="nav-saved-tab"
        data-toggle="tab"
        href="#rules-saved"
        role="tab"
        v-if="!queuedOnly && !currentOnly"
      >Saved</a>
      <a
        class="nav-link"
        :class="{ 'active': queuedOnly }"
        id="nav-queued-tab"
        data-toggle="tab"
        href="#rules-queued"
        role="tab"
        v-if="!currentOnly"
      >Queued</a>
    </nav>
    <div
      class="tab-content border rounded-bottom flex-grow-1 flex-shrink-1 overflow-auto"
      :class="{ 'border-top-0': !currentOnly }"
    >
      <!-- CURRENT RULES -->
      <div
        id="rules-current"
        class="list-group list-group-flush tab-pane fade active show"
        role="tab-panel"
        v-if="!queuedOnly"
      >
        <nav class="nav nav-pills p-2">
          <a
            class="nav-link active"
            id="nav-immutable-tab"
            data-toggle="tab"
            href="#rules-immutable"
            role="tab"
          >Immutable</a>
          <a
            class="nav-link"
            id="nav-mutable-tab"
            data-toggle="tab"
            href="#rules-mutable"
            role="tab"
          >Mutable</a>
        </nav>
        <div class="tab-content flex-grow-1 flex-shrink-1 overflow-auto">
          <!-- IMMUTABLE RULES -->
          <div
            id="rules-immutable"
            class="list-group list-group-flush tab-pane fade active show"
            role="tab-panel"
            v-if="!queuedOnly"
          >
            <p
              class="list-group-item border-0 text-muted"
              v-if="immutableRules.length === 0"
            >No immutable rules loaded...</p>
            <a
              class="ruleset btn btn-outline-primary list-group-item list-group-item-action"
              role="button"
              tabindex="0"
              v-bind:class="{
                'active':
                  loadedRule.type === ruleSetTypes.IMMUTABLE &&
                  loadedRule.index === index
              }"
              v-bind:key="index"
              v-for="(ruleSet, index) in immutableRules"
              v-on:click="selectRule(index, ruleSetTypes.IMMUTABLE)"
              style="cursor:pointer;"
            >
              <span class="font-weight-bold">Rule {{ index }}</span>
            </a>
          </div>
          <!-- MUTABLE RULES -->
          <div
            id="rules-mutable"
            class="list-group list-group-flush tab-pane fade"
            role="tab-panel"
            v-if="!queuedOnly"
          >
            <p
              class="list-group-item border-0 text-muted"
              v-if="mutableRules.length === 0"
            >No mutable rules loaded...</p>
            <a
              class="ruleset btn btn-outline-primary list-group-item list-group-item-action"
              role="button"
              tabindex="0"
              v-bind:class="{
                'active':
                  loadedRule.type === ruleSetTypes.MUTABLE &&
                  loadedRule.index === index
              }"
              v-bind:key="index"
              v-for="(ruleSet, index) in mutableRules"
              v-on:click="selectRule(index, ruleSetTypes.MUTABLE)"
              style="cursor:pointer;"
            >
              <span class="font-weight-bold">Rule {{ index }}</span>
            </a>
          </div>
        </div>
      </div>
      <!-- PLAYER-SAVED RULES -->
      <div
        id="rules-saved"
        class="list-group list-group-flush tab-pane fade"
        role="tab-panel"
        v-if="!currentOnly && !queuedOnly"
      >
        <p
          class="list-group-item border-0 text-muted"
          v-if="savedRules.length === 0"
        >No saved rules...</p>
        <p
          class="list-group-item border-0 text-muted"
          v-if="queuedRules.length > 0 && unQueuedRules.length === 0"
        >All of your saved rules are queued!</p>
        <a
          class="ruleset btn btn-outline-primary list-group-item list-group-item-action d-flex justify-content-between align-items-center"
          role="button"
          tabindex="0"
          v-bind:class="{
            'active':
              loadedRule.type === ruleSetTypes.SAVED &&
              loadedRule.index === savedIndex
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
          'active': queuedOnly,
          'show': queuedOnly
        }"
        v-if="!currentOnly"
      >
        <p
          class="list-group-item border-0 text-muted"
          v-if="queuedRules.length === 0"
        >No queued rules...</p>
        <a
          class="ruleset btn btn-outline-primary list-group-item list-group-item-action d-flex justify-content-between align-items-center"
          role="button"
          tabindex="0"
          v-bind:class="{
            'active':
              loadedRule.type === ruleSetTypes.QUEUED &&
              loadedRule.index === savedIndex
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
            v-if="!queuedOnly"
          ><small>Un&#8209;queue</small></button>
        </a>
      </div>
    </div>
  </div>

  <!-- Immutable and Mutable rules only -->
</template>

<script>
import { ruleSetTypes, ruleTypes } from '../../constants/constants.js';

export default {
  props: {
    loadedRule: {
      type: Object,
      default: () => ({
        index: 0,
        type: ruleSetTypes.IMMUTABLE
      })
    },
    immutableRules: {
      type: Array,
      default: () => []
    },
    mutableRules: {
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
    queuedOnly: {
      type: Boolean,
      default: false
    },
    currentOnly: {
      type: Boolean,
      default: false
    }
  },
  data: function () {
    return {
      ruleSetTypes: ruleSetTypes
    }
  },
  mounted: function () {
    // Auto-select Rule 0
    if (!this.queuedOnly) {
      this.$nextTick(() => {
        this.selectRule(0, ruleSetTypes.IMMUTABLE);
      })
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

      this.$emit('rule-selected', savedIndex, ruleType);
    },
    queueRule: async function (savedIndex) {
      this.$emit('queue-rule', savedIndex);
    },
    unQueueRule: async function (savedIndex, queuedIndex) {
      this.$emit('unqueue-rule', savedIndex, queuedIndex);
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