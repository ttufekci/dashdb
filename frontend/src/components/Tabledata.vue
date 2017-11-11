<template>
  <main role="main" class="col-sm-9 ml-sm-auto col-md-10 pt-3">
    <h1>{{this.$route.params.name}}</h1>
    <section class="row placeholders">
      <div class="table-responsive p-4">

        <a href="#" class="btn btn-primary">Add</a>

        <br/>
        <br/>

        <table class="table table-sm table-bordered">
          <thead class="thead-light">
              <tr>
                <th v-for="(item, index) in columnlist" :key="item.Name">
                  {{item.Name}}
                </th>
                <th></th>
                <th></th>                
              </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in rows" :key="index">
              <td v-for="(item, index) in item.S" :key="index">{{item}}</td>
              <td>
                  <router-link :to="{name:'editdata', params: {name:  tablename, primcols: primcols, id: item.Id, ids: ids}}" class="btn btn-secondary btn-sm">Edit</router-link>
              </td>
              <td>
                  <button type="button" class="btn btn-danger btn-sm">Delete</button>                 
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </main>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Tabledata',
  data () {
    return {
      msg: 'Ready-made Dashboard for your Database',
      errors: [],
      testsuccess: '',
      host: '',
      user: '',
      password: '',
      schema: 'dashdb',
      tablelist: [],
      columnlist: [],
      rows: [],
      tablename: '',
      primcols: '',
      ids: '',
      id: ''
    }
  },
  mounted () {
    // this.$router.push('overview')
    this.tablename = this.$route.params.name
    this.readColumnList(this.$route.params.name)
  },
  methods: {
    readColumnList (tablename) {
      axios.get('http://localhost:8081/columnlist?name=' + tablename)
      .then(response => {
        // JSON responses are automatically parsed.
        this.columnlist = response.data.cols
        this.rows = response.data.datas
        this.primcols = response.data.primcols
        this.ids = response.data.ids
        this.id = response.data.id
      })
      .catch(e => {
        console.log(e)
        this.errors.push(e)
      })
    }
  },
  watch: {
    '$route' (to, from) {
      // react to route changes...
      this.readColumnList(this.$route.params.name)
    }
  }
}
</script>

<style>

</style>
