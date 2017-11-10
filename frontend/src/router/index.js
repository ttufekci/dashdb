import Vue from 'vue'
import Router from 'vue-router'
import Entrance from '@/components/Entrance'
import Dashboard from '@/components/Dashboard'
import Overview from '@/components/Overview'
import Reports from '@/components/Reports'

Vue.use(Router)

export default new Router({
  routes: [{
    path: '/',
    name: 'Entrance',
    component: Entrance
  }, {
    path: '/dashboard',
    component: Dashboard,
    children: [
      {
        // UserProfile will be rendered inside User's <router-view>
        // when /user/:id/profile is matched
        path: '/reports',
        name: 'reports',
        component: Reports
      },
      {
        // UserProfile will be rendered inside User's <router-view>
        // when /user/:id/profile is matched
        path: '/overview',
        name: 'overview',
        component: Overview
      }
    ]
  }]
})
