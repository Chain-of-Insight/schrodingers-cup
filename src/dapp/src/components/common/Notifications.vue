<template>
  <div
    class="notification"
    :class="{
      'fixed-top': !local,
      'w-100': !local
    }">
    <!-- Notification Display -->
    <div v-if="type && msg" :class="'alert alert-' + type" role="alert">{{ msg }}</div>
    <span v-if="type && msg" class="close close-x" @click="$emit('reset')">&times;</span>
  </div>
</template>

<script>
export default {
  props: {
    type: String,
    msg: String,
    local: {
      type: Boolean,
      default: false
    }
  }
};
</script>

<style scoped>
  .alert {
    text-align: center;
  }
  .close-x {
    position: relative;
    right: 1.25rem;
    top: -55px;
    text-align: right;
    cursor: pointer;
    color: #155724;
  }
</style>