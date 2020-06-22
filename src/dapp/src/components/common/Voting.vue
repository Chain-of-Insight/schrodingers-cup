<template>
  <div>  
    <!-- Modal -->
    <div class="modal fade" id="voting-modal" data-backdrop="static" data-keyboard="false" tabindex="-1" role="dialog" aria-labelledby="voting-modal-label" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="voting-modal-label"><strong>Time to vote!</strong></h5>
            <Countdown
              :duration="votingDuration"
              v-on:ended="closeVotingWindow()"
              ref="timer"
            ></Countdown>
            <!-- <div>
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
            </div> -->
          </div>
          <div class="modal-body">
            <pre v-if="votingCandidate" class="term-container">{{ votingCandidate.code }}</pre>
          </div>
          <div class="modal-footer container-fluid">
            <div class="container-fluid p-0">
              <div class="row">
                <div class="col-3">
                  <button
                    type="button"
                    class="btn btn-block btn-lg btn-secondary"
                    @click="vote(voteType.ABSTAIN)"
                    :class="{ 'disabled': votingWindowClosed } "
                    :disabled="votingWindowClosed"
                  >Abstain</button>
                </div>
                <div class="col-9">
                  <div class="btn-group btn-block">
                    <button
                      type="button"
                      class="btn btn-lg btn-danger"
                      @click="vote(voteType.NO)"
                      :class="{ 'disabled': votingWindowClosed } "
                      :disabled="votingWindowClosed"
                    >
                      <span class="oi mirrored oi-thumb-down" title="Vote down"></span>
                    </button>
                    <button
                      type="button"
                      class="btn btn-lg btn-success"
                      @click="vote(voteType.YES)"
                      :class="{ 'disabled': votingWindowClosed } "
                      :disabled="votingWindowClosed"
                    >
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

  import Countdown from '../common/Countdown.vue';

  export default {
    components: { Countdown },
    props: {
      votingDuration: {
        required: true,
        type: Number
      }
    },
    data: function () {
      return {
        votingCandidate: null,
        votingWindowClosed: false,
        voteType: {
          YES: 0,
          NO: 1,
          ABSTAIN: -1
        }
      }
    },
    mounted: function () {
      // Start timer
      $('#voting-modal').on('shown.bs.modal', this.startTimer.bind(this));
      $('#voting-modal').on('hidden.bs.modal', this.resetTimer.bind(this));
    },
    methods: {
      promptForVote: function (votingCandidate) {
        if (!votingCandidate.code) {
          return;
        }

        this.votingCandidate = votingCandidate
        $('#voting-modal').modal('show');
      },
      // startTimer: function () {
      //   this.secondsLeft = this.votingDuration;
      //   this.timer = setInterval(this.timerDecrement, 1000);
      // },
      // resetTimer: function () {
      //   clearInterval(this.timer);
      //   this.secondsLeft = this.votingDuration;
      //   this.votingCandidate = null;
      // },
      // timerDecrement: function () {
      //   if (this.secondsLeft > -2) {
      //     this.secondsLeft -= 1;
      //   } else {
      //     $('#voting-modal').modal('hide');
      //     // Register vote as abstained if no action when time runs out
      //     this.$emit('vote-cast', 'abstain');
      //   }
      // },
      vote: function (type) {
        if (
          typeof(type) !== 'number' || (
            type !== this.voteType.YES &&
            type !== this.voteType.NO &&
            type !== this.voteType.ABSTAIN
          )
        ) {
          return false;
        }

        this.$emit('vote-cast', type);
        $('#voting-modal').modal('hide');
      },
      closeVotingWindow: function () {
        this.votingWindowClosed = true;
        $emit('vote-cast', voteType.ABSTAIN);
      },
      startTimer: function () {
        this.$refs.timer.start();
      },
      resetTimer: function () {
        this.$refs.timer.reset();
      }
    }
  }
</script>