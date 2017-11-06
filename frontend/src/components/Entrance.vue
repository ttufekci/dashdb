<template>
  <div class="container">
    <br>
    <h1 class="text-center dotted-border">{{ msg }}</h1>
    <br>
    <div class="row">
      <div class="col-6">
        <!-- Nav tabs -->
        <ul class="nav nav-tabs" id="myTab" role="tablist">
          <li class="nav-item">
            <a class="nav-link active" id="home-tab" data-toggle="tab" href="#home" role="tab" aria-controls="home" aria-selected="true">MySQL Settings</a>
          </li>
        </ul>

        <!-- Tab panes -->
        <div class="tab-content">
          <div class="tab-pane active" id="home" role="tabpanel" aria-labelledby="home-tab">
            <form>
              <br>
              <div class="form-group">
                <label for="inputHost">Host</label>
                <input type="host" class="form-control" id="inputHost" placeholder="Host" v-model="host">
              </div>
              <div class="form-group">
                <label for="inputSchema">Schema</label>
                <input type="schema" class="form-control" id="inputSchema" placeholder="Schema" v-model="schema">
              </div>          
              <div class="form-group">
                <label for="inputUser">User</label>
                <input type="user" class="form-control" id="inputUser" placeholder="User" v-model="user">
              </div>          
              <div class="form-group">
                <label for="exampleInputPassword1">Password</label>
                <input type="password" class="form-control" id="exampleInputPassword1" placeholder="Password" v-model="password">
              </div>
              <br>
              <button type="submit" class="btn btn-primary">Save</button>
              <button type="submit" class="btn btn-primary" v-on:click="testConnection">Test Connection</button>
              <label class="text-success ml-3">{{testsuccess}}</label>
            </form>        
          </div>
        </div>
      </div>
      <div class="col-1">        
      </div>
    </div>        
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Entrance',
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
  },
  methods: {
    testConnection () {
      var setGetParams = 'host=' + this.host + '&user=' + this.user + '&password=' + this.password + '&schema=' + this.schema
      axios.get('http://localhost:8081/testconnection?' + setGetParams)
      .then(response => {
        // JSON responses are automatically parsed.
        var success = response.data.success
        if (success === 'true') {
          this.testsuccess = 'Success'
        } else {
          this.testsuccess = 'Fails'
        }
        console.log(response)
        console.log(response.data)
      })
      .catch(e => {
        console.log(e)
        this.errors.push(e)
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.dotted-border {
  border-bottom: thick dotted #ff0000;
  color: green;
}
</style>
