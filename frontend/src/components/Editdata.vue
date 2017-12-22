<template>
  <main role="main" class="col-sm-9 ml-sm-auto col-md-10 pt-3">
    <h1>edit {{this.$route.params.name}}</h1>
    <section class="row placeholders">
      <div class="table-responsive p-4">

        <button class="btn btn-success" type="submit" @click="goBack">Return</button>

        <br/>
        <br/>

        <span class="text-primary">{{saveMessage}}</span>

        <br/>
        <br/>

        <table class="table table-sm table-bordered">
          <tbody>
            <tr v-for="(item, index) in columnlist" :key="index">
                <template v-if="item.Ai">
                    <td>{{item.Name}} (auto increment)</td>
                    <td>                            
                        <input type="text" name="fields" disabled v-model="item.Value"/>
                    </td>
                </template>
                <template v-else>
                    <td>{{item.Name}}</td>
                    <td>                            
                        <input type="text" name="fields" v-model="item.Value"/>
                    </td>
                </template>                              
            </tr>
            <tr>
                <td><button class="btn btn-default" type="submit" @click="saveData">Save</button>
                </td><td></td>
            </tr>
          </tbody>
        </table>
        <input type="hidden" name="name" :value="tablename" />
        <input type="hidden" name="id" :value="id" />
        <input type="hidden" name="ids" :value="ids" />        
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
      id: '',
      ids: '',
      fields: [],
      saveMessage: ''
    }
  },
  mounted () {
    this.readEditData(this.$route.params.name, this.$route.params.primcols, this.$route.params.ids, this.$route.params.id)
  },
  methods: {
    readEditData (tablename, primcols, ids, id) {
      axios.get('http://localhost:8081/reditdata?name=' + tablename + '&primcols=' + primcols + '&ids=' + ids + '&id=' + id)
      .then(response => {
        // JSON responses are automatically parsed.
        this.columnlist = response.data.cols
        this.rows = response.data.datas
        this.id = response.data.id
        this.ids = response.data.ids
        this.tablename = this.$route.params.name

        for (var i = 0; i < this.columnlist.length; i++) {
          this.fields.push(this.columnlist[i].Value)
        }
      })
      .catch(e => {
        console.log(e)
        this.errors.push(e)
      })
    },
    saveData () {
      console.log('this.tablename:' + this.tablename)
      console.log('fields:', this.fields)
      var paramFields = []
      for (var i = 0; i < this.columnlist.length; i++) {
        paramFields.push(this.columnlist[i].Value)
      }
      axios.post('http://localhost:8081/reditdatam', {
        'name': this.tablename,
        'primcols': this.primcols,
        'ids': this.ids,
        'id': this.id,
        'fields': paramFields
      })
      .then(response => {
        this.saveMessage = 'Saved Successfully'
      })
      .catch(e => {
        console.log(e)
        this.errors.push(e)
      })
    },
    goBack () {
      this.$router.go(-1)
    }
  },
  watch: {
    '$route' (to, from) {
      // react to route changes...
      this.readEditData(this.$route.params.name, this.$route.params.primcols, this.$route.params.ids, this.$route.params.id)
    }
  }
}
</script>

<style>

</style>
