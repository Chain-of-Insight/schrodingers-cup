<template>
  <div>
    <!-- Button trigger modal -->
    <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#voting-modal">
      Voting Test
    </button>
  
    <!-- Modal -->
    <div class="modal fade" id="voting-modal" data-backdrop="static" data-keyboard="false" tabindex="-1" role="dialog" aria-labelledby="voting-modal-label" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="voting-modal-label"><strong>Time to vote!</strong></h5>
            <div>
              <h5 class="d-inline">Time remaining:</h5>
              <h5 class="d-inline h4 border border-dark rounded ml-2 p-2 text-monospace">{{ hours }}:{{ minutes }}:{{ seconds }}</h5>
            </div>
          </div>
          <div class="modal-body">
            <pre class="term-container">
$test_string = "this is test code"
say($test_string)
            </pre>
          </div>
          <div class="modal-footer container-fluid">
            <div class="container-fluid p-0">
              <div class="row">
                <div class="col-3">
                  <button type="button" class="btn btn-block btn-lg btn-secondary" data-dismiss="modal">Abstain</button>
                </div>
                <div class="col-9">
                  <div class="btn-group btn-block">
                    <button type="button" class="btn btn-lg btn-danger" data-dismiss="modal">
                      <span class="oi mirrored oi-thumb-down" title="Vote down"></span>
                    </button>
                    <button type="button" class="btn btn-lg btn-success" data-dismiss="modal">
                      <span class="oi oi-thumb-up" title="Vote up"></span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  const $ = window.jQuery;

  export default {
    components: {},
    props: {},
    data: function () {
      return {
        votingDuration: 300,
        secondsLeft: 300,
        timer: null
      }
    },
    mounted: function () {
      // Start timer
      // TODO: move this to modal open
      console.log('timer', this.timer);
      $('#voting-modal').on('shown.bs.modal', this.startTimer.bind(this));
      $('#voting-modal').on('hidden.bs.modal', this.resetTimer.bind(this));
    },
    destroyed: function () {
      // clearInterval(this.timer);
    },
    computed: {
      hours: function () {
        const hours = Math.floor(this.secondsLeft / 3600);
        return hours.toString().padStart(2, '0');
      },
      minutes: function () {
        const minutes = Math.floor((this.secondsLeft % 3600) / 60);
        return minutes.toString().padStart(2, '0');
      },
      seconds: function () {
        const seconds = this.secondsLeft % 60;
        return seconds.toString().padStart(2, '0');
      }
    },
    methods: {
      startTimer: function () {
        this.timer = setInterval(this.timerDecrement, 1000);
      },
      resetTimer: function () {
        clearInterval(this.timer);
        this.secondsLeft = this.votingDuration;
      },
      timerDecrement: function () {
        if (this.secondsLeft > 0) {
          this.secondsLeft -= 1;
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
</style>