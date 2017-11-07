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
              <button type="submit" class="btn btn-primary" v-on:click="saveConfig">Save Configuration</button>
              <button type="submit" class="btn btn-primary" v-on:click="testConnection">Test Connection</button>
              <label class="text-success ml-3">{{testsuccess}}</label>
            </form>        
          </div>
        </div>
      </div>
      <div class="col-3">
        <div style="height: 400px; padding-top: 200px; padding-left: 130px;">
          <a href="#" class="btn btn-success">Start Your Adventure =></a>
        </div>
      </div>
      <div class="col-3">
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
  created () {
    this.readConfig()
  },
  methods: {
    testConnection () {
      var setGetParams = 'host=' + this.host + '&user=' + this.user + '&password=' + this.password + '&schema=' + this.schema
      axios.get('http://localhost:8081/testconnection?' + setGetParams)
      .then(response => {
        // JSON responses are automatically parsed.
        var success = response.data.success
        if (success === 'true') {
          this.testsuccess = 'Connection is live :)'
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
    },
    readConfig () {
      axios.get('http://localhost:8081/readconfig')
      .then(response => {
        // JSON responses are automatically parsed.
        var success = response.data.success
        if (success === 'true') {
          this.testsuccess = 'Read successfully'
          this.user = response.data.user
          this.host = response.data.host
          this.password = response.data.password
          this.schema = response.data.schema
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
    },
    saveConfig () {
      var setGetParams = 'host=' + this.host + '&user=' + this.user + '&password=' + this.password + '&schema=' + this.schema
      axios.get('http://localhost:8081/saveconfig?' + setGetParams)
      .then(response => {
        // JSON responses are automatically parsed.
        var success = response.data.success
        if (success === 'true') {
          this.testsuccess = 'Saved Successfully'
        } else {
          this.testsuccess = 'Saved Fails'
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
