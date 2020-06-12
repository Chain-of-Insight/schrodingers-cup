<template>
  <div class="rule">
    <!-- Rule Update (Change, Remove, Add) -->
    <div 
      v-if="rule && changeSet && type == 'update' && !queued" 
      :class="'rule-proposal type-' + type" 
      role="rule"
    >
      <!-- New Rule -->
      <div id="ide_new" class="ide input">
        <label>
          <strong>You are proposing this rule change:</strong>
        </label>
        <codemirror 
          v-model="changeSet.new"
          :options="ide.options"
          readonly
        ></codemirror>
        <!-- E.g. Delete, Edit, Add -->
        <p>Change Type: {{ type }}</p>
        <!-- Mutable / Immutable -->
        <p v-if="rule.type">Rule Type: {{ rule.type }}</p>
        <p v-if="changeSet.old">
          <span v-if="changeSet.old.id">Updating existing rule: {{ changeSet.old.id }}</span>
        </p>
      </div>

      <!-- Old Rule -->
      <div id="ide_old" class="ide input" v-if="changeSet.old">
        <label>
          <strong>Replacing this rule:</strong>
        </label>
        <codemirror 
          v-model="changeSet.old"
          :options="ide.options"
          readonly
        ></codemirror>
      </div>

      <button class="btn btn-primary">Queue Rule</button>
    </div>
  
    <!-- Else, Queued -->
    <div v-if="queued" class="alert-success">
      <p>Your rule has been queued, good luck!</p>
      <span class="close close-x" @click="$emit('close-rule-proposal')">&times;</span>
    </div>
  </div>
</template>

<script>
// IDE
import 'codemirror/mode/shell/shell';
import 'codemirror/theme/dracula.css';

export default {
  props: {
    type: String,
    rule: Object, // [rule ID on tezos, rule type (immutable/mutable)]
    changeSet: Object,
    queued: Boolean
  },
  data: () => ({
    ide: {
      options: {
        tabSize: 4,
        theme: 'dracula',
        lineNumbers: true,
        line: true
      }
    }
  })
};
</script>