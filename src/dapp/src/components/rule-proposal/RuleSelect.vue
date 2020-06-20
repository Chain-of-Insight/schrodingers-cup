<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col">
        <h2>{{ typeHeadings[changeType] }}</h2>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    props: {
      changeType: String,
      typeHeadings: Object
    }
  }
</script>

<style scoped>

</style>