import Vue from 'vue'
import Router from 'vue-router'
import HeaderList from '@/components/HeaderList'
import Config from '@/components/Config'
import About from '@/components/About'
import Testing from '@/components/Testing'
import Login from '@/components/Login'
import ProxyFrame from '@/components/ProxyFrame'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'HeaderList',
      component: HeaderList
    }, {
      path: '/about',
      name: 'About',
      component: About
    }, {
      path: '/config',
      name: 'Config',
      component: Config
    }, {
      path: '/testing',
      name: 'Testing',
      component: Testing
    }, {
      path: '/login',
      name: 'Login',
      component: Login
    }, {
      path: '/frame',
      name: 'Frame',
      component: ProxyFrame
    }
  ]
})
