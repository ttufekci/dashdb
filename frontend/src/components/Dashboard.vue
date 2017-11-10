<template>
  <div>
    <header>
      <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
        <a class="navbar-brand" href="#">Dashboard</a>
        <button class="navbar-toggler d-lg-none" type="button" data-toggle="collapse" data-target="#navbarsExampleDefault" aria-controls="navbarsExampleDefault" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbarsExampleDefault">
          <ul class="navbar-nav mr-auto">
            <li class="nav-item active">
              <a class="nav-link" href="#">Home <span class="sr-only">(current)</span></a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="#">Settings</a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="#">Profile</a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="#">Help</a>
            </li>
          </ul>
          <form class="form-inline mt-2 mt-md-0">
            <input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search">
            <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
          </form>
        </div>
      </nav>
    </header>

    <div class="container-fluid">
      <div class="row">
        <nav class="col-sm-3 col-md-2 d-none d-sm-block bg-light sidebar">
          <ul class="nav nav-pills flex-column">
            <li class="nav-item">
              <router-link :to="{name: 'overview'}" class="nav-link" active-class="active">Overview <span class="sr-only">(current)</span></router-link>                            
            </li>
            <li class="nav-item">
              <router-link :to="{name: 'reports'}" active-class="active" class="nav-link">Reports</router-link>              
            </li>
          </ul>

          <ul class="nav nav-pills flex-column">
            <li v-for="item in tablelist" :key="item" class="nav-item">              
              <router-link :to="{name:'tabledata', params: {name:  item}}" active-class="active" class="nav-link">{{ item }}</router-link>              
            </li>
          </ul>

        </nav>

        <router-view/>
      </div>
    </div>    
  </div>
</template>

<script>
import axios from 'axios'

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
      schema: 'dashdb',
      tablelist: [

      ]
    }
  },
  mounted () {
    // this.$router.push('overview')
    this.readTableList()
  },
  methods: {
    readTableList () {
      axios.get('http://localhost:8081/tablelist')
      .then(response => {
        // JSON responses are automatically parsed.
        this.tablelist = response.data
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