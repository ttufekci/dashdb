<template>
  <div class="container">
      <h1>Dashboard</h1>
  </div>
</template>

<script>

export default {
  name: 'Dashboard',
  data () {
    return {
      msg: 'Ready-made Dashboard for your Database',
      errors: [],
      testsuccess: '',
      host: '',
      user: '',
      password: '',
      schema: 'dashdb'
    }
  }
}

</script>