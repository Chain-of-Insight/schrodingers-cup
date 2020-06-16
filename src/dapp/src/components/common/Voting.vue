<template>
  <div>
    <div class="container-fluid voting">
      <div class="row">
        <div class="col">
          <!-- Button trigger modal -->
          <button type="button" class="btn btn-primary" data-toggle="modal" data-target="#voting-modal">
            Voting Test
          </button>
        </div>
      </div>
    </div>
  
    <!-- Modal -->
    <div class="modal fade" id="voting-modal" data-backdrop="static" data-keyboard="false" tabindex="-1" role="dialog" aria-labelledby="voting-modal-label" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="voting-modal-label">Time to vote!</h5>
          </div>
          <div class="modal-body">
            ...
          </div>
          <div class="modal-footer">
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
</template>

<script>
  export default {
    components: {},
    props: {},
    data: () => ({}),
    mounted: () => {},
    methods: {}
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