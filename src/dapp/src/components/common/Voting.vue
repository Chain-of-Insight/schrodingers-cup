<template>
  <div>
    <div class="border-dark border-top border-bottom py-3 my-3">
      <div class="row">
        <!-- Rule info -->
        <div class="col-7 border-right border-dark">
          <h3 class="mb-3"><strong>Time to vote!</strong></h3>
          <h5 class="text-primary"><strong>{{ votingCandidate.proposal.toUpperCase() }} RULE</strong></h5>
          <p class="mb-1">Proposed by: <strong class="text-info">{{ votingCandidate.author }}</strong></p>
          <p v-if="votingCandidate.proposal !== 'create'">Rule number: <strong>{{ votingCandidate.index }}</strong></p>
          <button class="btn btn-sm btn-warning" type="button" @click="codeVisible = codeVisible ? false : true" data-toggle="collapse" data-target="#rule-code" aria-expanded="false" aria-controls="rule-code">
            {{ codeVisible ? 'Hide rule code' : 'Show rule code' }}
          </button>
        </div>

        <div class="col-5 d-flex flex-column">
          <button
            type="button"
            class="btn btn-block btn-success flex-grow-1"
            @click="vote(voteTypes.YES)"
          >
            <span class="oi mirrored oi-thumb-up" title="Vote up"></span>
          </button>
          <button
            type="button"
            class="btn btn-block btn-danger flex-grow-1"
            @click="vote(voteTypes.NO)"
          >
            <span class="oi oi-thumb-down" title="Vote down"></span>
          </button>
          <button
            type="button"
            class="btn btn-block btn-secondary flex-grow-1"
            @click="vote(voteTypes.ABSTAIN)"
          >Abstain</button>
        </div>
      </div>
      <div
        class="row collapse"
        id="rule-code"
        v-if="votingCandidate.proposal !== 'delete' && votingCandidate.original"
      >
        <div class="col-6 mt-4">
          <h6 class="font-weight-bold">Current:</h6>
          <pre class="term-container rounded mb-0">{{ votingCandidate.original }}</pre>
        </div>
        <div class="col-6 mt-4">
          <h6 class="font-weight-bold">Proposed:</h6>
          <pre class="term-container rounded mb-0">{{ votingCandidate.code }}</pre>
        </div>
      </div>
      <div
        class="row collapse"
        id="rule-code"
        v-else
      >
        <div class="col mt-4">
          <h6 class="font-weight-bold" v-if="votingCandidate.proposal === 'delete'">To be deleted:</h6>
          <pre class="term-container rounded mb-0">{{ votingCandidate.proposal === 'delete' ? votingCandidate.original : votingCandidate.code }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  const $ = window.jQuery;

  import Notification from '../common/Notifications.vue';

  export default {
    components: {
      Notification
    },
    props: {
      votingCandidate: {
        type: Object,
        default: () => ({})
      }
    },
    data: function () {
      return {
        codeVisible: false,
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
    methods: {
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
      toggleRuleCode: function () {
        
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

<style scoped>
#rule-code .term-container {
  max-height: 330px;
  overflow-y: auto;
}

.oi.mirrored {
  -webkit-transform: scale(-1, 1);
  -moz-transform: scale(-1, 1);
  -ms-transform: scale(-1, 1);
  -o-transform: scale(-1, 1);
  transform: scale(-1, 1);
}
</style>