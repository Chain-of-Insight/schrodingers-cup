<template>
  <div>
    <h5 class="d-inline">Time remaining:</h5>
    <h5
      class="countdown d-inline h4 border rounded ml-2 p-2 text-monospace"
      :class="[
        secondsLeft <= 10 ? [
          'critical',
          'text-danger',
          'border-danger'
        ] : 'border-dark'
      ]"
    >{{ hours }}:{{ minutes }}:{{ seconds }}</h5>
  </div>
</template>

<script>
  export default {
    props: {
      duration: {
        required: true,
        type: Number
      }
    },
    data: function () {
      return {
        secondsLeft: null,
        timer: null
      }
    },
    destroyed: function () {
      clearInterval(this.timer);
    },
    computed: {
      hours: function () {
        let hours = null;
        hours = Math.floor(this.secondsLeft / 3600);
        return hours.toString().padStart(2, '0');
      },
      minutes: function () {
        let minutes = null;
        minutes = Math.floor((this.secondsLeft % 3600) / 60);
        return minutes.toString().padStart(2, '0');
      },
      seconds: function () {
        let seconds = null;
        seconds = this.secondsLeft % 60;
        return seconds.toString().padStart(2, '0');
      }
    },
    methods: {
      start: function () {
        this.secondsLeft = this.duration;
        this.timer = setInterval(this.timerDecrement, 1000);
      },
      reset: function () {
        clearInterval(this.timer);
        this.secondsLeft = this.duration;
      },
      timerDecrement: function () {
        if (this.secondsLeft > 0) {
          this.secondsLeft -= 1;
        } else {
          this.secondsLeft = 0;
          this.$emit('ended');
        }
      }
    }
  }
</script>

<style scoped>
.oi.mirrored {
  -webkit-transform: scale(-1, 1);
  -moz-transform: scale(-1, 1);
  -ms-transform: scale(-1, 1);
  -o-transform: scale(-1, 1);
  transform: scale(-1, 1);
}

.countdown.critical {
  -webkit-animation: blink 1s steps(1) infinite;
  -moz-animation: blink 1s steps(1) infinite;
  -o-animation: blink 1s steps(1) infinite;
  animation: blink 1s steps(1) infinite;
}

@keyframes blink {
  75% {
    opacity: 0;
  }
}
</style>