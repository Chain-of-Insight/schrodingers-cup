<template>
  <div class="w-100 d-flex justify-content-between align-items-center border-dark border-top border-bottom py-3 my-3">
    <p class="h2 m-0">Round {{ round }}</p>
    <p class="h2 m-0">Yes: <span class="badge badge-pill badge-success">{{ yes }}</span></p>
    <p class="h2 m-0">No: <span class="badge badge-pill badge-danger">{{ no }}</span></p>
    <p class="h2 m-0">Abstain: <span class="badge badge-pill badge-secondary">{{ abstain }}</span></p>
  </div>
</template>

<script>
  import Countdown from '../common/Countdown.vue';

  export default {
    components: { Countdown },
    props: {
      round: {
        required: true,
        default: 1,
        type: Number
      },
      yes: {
        required: true,
        default: 0,
        type: Number
      },
      no: {
        required: true,
        default: 0,
        type: Number
      },
      abstain: {
        required: true,
        default: 0,
        type: Number
      },
      turnWindow: {
        required: true,
        type: Number
      }
    },
    methods: {
      endRound: function () {
        this.$emit('round-over');
      },
      startTimer: function () {
        this.$refs.timer.start();
      },
      resetTimer: function () {
        this.$refs.timer.reset();
      },
    }
  }
</script>