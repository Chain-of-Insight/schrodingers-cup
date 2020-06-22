<template>
  <div>  
    <!-- Modal -->
    <div class="modal fade" id="voting-modal" data-backdrop="static" data-keyboard="false" tabindex="-1" role="dialog" aria-labelledby="voting-modal-label" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="voting-modal-label"><strong>Time to vote!</strong></h5>
            <Countdown
              :duration="turnWindow"
              v-on:ended="closeTurnWindow()"
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
            <!-- Notifications -->
            <Notification 
              :type="alert.type" 
              :msg="alert.msg" 
              v-on:reset="alert = {type: null, msg: null}"
              :local="true"
            ></Notification>
            <pre v-if="votingCandidate" class="term-container">{{ votingCandidate.code }}</pre>
          </div>
          <div class="modal-footer container-fluid">
            <div class="container-fluid p-0">
              <div class="row">
                <div class="col-3">
                  <button
                    type="button"
                    class="btn btn-block btn-lg btn-secondary"
                    @click="vote(voteTypes.ABSTAIN)"
                    :class="{ 'disabled': votingWindowClosed } "
                    :disabled="votingWindowClosed"
                  >Abstain</button>
                </div>
                <div class="col-9">
                  <div class="btn-group btn-block">
                    <button
                      type="button"
                      class="btn btn-lg btn-danger"
                      @click="vote(voteTypes.NO)"
                      :class="{ 'disabled': votingWindowClosed } "
                      :disabled="votingWindowClosed"
                    >
                      <span class="oi mirrored oi-thumb-down" title="Vote down"></span>
                    </button>
                    <button
                      type="button"
                      class="btn btn-lg btn-success"
                      @click="vote(voteTypes.YES)"
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
  import Notification from '../common/Notifications.vue';

  export default {
    components: {
      Countdown,
      Notification
    },
    props: {
      turnWindow: {
        required: true,
        type: Number
      }
    },
    data: function () {
      return {
        votingCandidate: null,
        votingWindowClosed: false,
        voteTypes: {
          YES: 1,
          NO: 0,
          ABSTAIN: -1
        },
        alert: {
          type: null,
          msg: null
        },
      }
    },
    mounted: function () {
      // Start timer
      $('#voting-modal').on('shown.bs.modal', this.startTimer.bind(this));
      $('#voting-modal').on('hidden.bs.modal', this.resetModal.bind(this));
    },
    methods: {
      promptForVote: function (votingCandidate) {
        if (!votingCandidate.code) {
          return;
        }

        this.votingCandidate = votingCandidate
        $('#voting-modal').modal('show');
      },
      closeModal: function () {
        $('#voting-modal').modal('hide');
      },
      vote: function (type) {
        if (
          typeof(type) !== 'number' || (
            type !== this.voteTypes.YES &&
            type !== this.voteTypes.NO &&
            type !== this.voteTypes.ABSTAIN
          )
        ) {
          return false;
        }

        this.$emit('vote-cast', type);
      },
      closeTurnWindow: function () {
        this.votingWindowClosed = true;
        $emit('vote-cast', voteTypes.ABSTAIN);
        $('#voting-modal').modal('hide');
      },
      startTimer: function () {
        this.$refs.timer.start();
      },
      resetModal: function () {
        this.votingCandidate = null;
        this.$refs.timer.reset();
        this._retireNotification();
      },
      _retireNotification: function () {
        this.alert = {
          type: null,
          msg: null
        };
      },
    }
  }
</script>