import Vue from 'vue'
import Router from 'vue-router'
import Entrance from '@/components/Entrance'
import Dashboard from '@/components/Dashboard'

Vue.use(Router)

export default new Router({
  routes: [{
    path: '/',
    name: 'Entrance',
    component: Entrance
  }, {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard
  }]
})