<template>
  <main role="main" class="col-sm-9 ml-sm-auto col-md-10 pt-3">
    <h1>{{this.$route.params.name}}</h1>
    <section class="row placeholders">
      <div class="table-responsive p-4">

        <router-link :to="{name:'adddata', params: {name:  tablename}}" class="btn btn-primary">Add</router-link>        

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
              <td v-for="(itemdetail, index) in item.S" :key="index">{{itemdetail}}</td>
              <td>
                  <router-link :to="{name:'editdata', params: {name:  tablename, primcols: primcols, id: item.Id, ids: item.Ids}}" class="btn btn-secondary btn-sm">Edit</router-link>
              </td>
              <td>
                  <button type="button" class="btn btn-danger btn-sm" @click="deleteRow(primcols, item.Id, item.Ids)">Delete</button>                 
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
      tablename: 'test',
      primcols: '',
      ids: '',
      id: ''
    }
  },
  mounted () {
    this.readColumnList(this.$route.params.name)
  },
  methods: {
    deleteRow (primcols, id, ids) {
      console.log(primcols, 'id:', id, ',ids:', ids)
      var txt
      var r = confirm('record will be deleted!')
      if (r === true) {
        txt = 'Record will be deleted!'
        axios.post('http://localhost:8081/deleterowdata', {
          'name': this.tablename,
          'primcols': this.primcols,
          'ids': ids,
          'id': String(id),
          'fields': 'test'
        })
        .then(response => {
          this.saveMessage = 'Deleted Successfully'
          this.readColumnList(this.$route.params.name)
        })
        .catch(e => {
          console.log(e)
          this.errors.push(e)
        })
      } else {
        txt = 'Cancel!'
      }
      console.log(txt)
    },
    readColumnList (tablename) {
      this.tablename = tablename
      axios.get('http://localhost:8081/columnlist?name=' + tablename)
      .then(response => {
        console.log('columnlist loaded')
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
