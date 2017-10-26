import Vue from 'vue'
import Router from 'vue-router'
import Entrance from '@/components/Entrance'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Entrance',
      component: Entrance
    }
  ]
})
