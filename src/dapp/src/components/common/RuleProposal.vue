<template>
  <div>
    <!-- Modal -->
    <div class="modal fade" id="proposal-modal" data-backdrop="static" data-keyboard="false" tabindex="-1" role="dialog" aria-labelledby="proposal-modal-label" aria-hidden="true">
      <div class="modal-dialog modal-xl">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="proposal-modal-label">Time to propose a rule change!</h5>
          </div>
          <div class="modal-body">
            ...
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
const $ = window.jQuery;

export default {
  props: {
    changeType: String,
    rule: Object
  },
  data: () => ({
    // ...
  }),
  methods: {
    promptForProposal: function () {
      $('#proposal-modal').modal('show');
    }
  }
};
</script>