<template>
  <div class="w-100 d-flex justify-content-between align-items-center border-dark border-top border-bottom py-3 my-3">
    <p class="h2 m-0">Round {{ round }}</p>
    <p class="h2 m-0">Yes: <span class="badge badge-pill badge-success">{{ yesTotals }}</span></p>
    <p class="h2 m-0">No: <span class="badge badge-pill badge-danger">{{ noTotals }}</span></p>
    <p class="h2 m-0">Abstain: <span class="badge badge-pill badge-secondary">{{ abstainTotals }}</span></p>
  </div>
</template>

<script>
  export default {
    props: {
      round: {
        required: true,
        default: 1,
        type: Number
      },
      yesTotals: {
        required: true,
        default: 0,
        type: Number
      },
      noTotals: {
        required: true,
        default: 0,
        type: Number
      },
      abstainTotals: {
        required: true,
        default: 0,
        type: Number
      }
    }
  }
</script>