<template>
  <div class="container-fluid">
    <navbar></navbar>
    <div class="row">
      <div class="col-12 col-md-3 col-lg-2">
        <left-menu v-on:select-component="selectComponent"></left-menu>
      </div>
      <div class="col-12 col-md-9 col-lg-10">
        <component v-bind:is="currentTabComponent"></component>
      </div>
    </div>
  </div>
</template>

<script>
import LeftMenu from './components/LeftMenu'
import About from './components/About'
import Headers from './components/Headers'
import Config from './components/Config'
import Overrides from './components/Overrides'
import Navbar from './components/Navbar'

export default {
  name: 'App',
  components: {
    Navbar, About, LeftMenu, Headers, Config, Overrides
  },
  data () {
    return {
      currentComponent: 'About'
    }
  },
  computed: {
    currentTabComponent: function () {
      return this.currentComponent !== undefined || null ? this.currentComponent : 'About'
    }
  },
  methods: {
    selectComponent: function (componentName) {
      console.log(`hello component: ${componentName}`)
      this.currentComponent = componentName
    }
  }
}
// +++ the material design icons
// https://material.io/tools/icons/?style=baseline
</script>

<style lang="css">
@import '../node_modules/bootstrap/dist/css/bootstrap.css';
@import '/assets/css/base.css';
@import '../static/css/ionicons.min.css';
@import '../static/css/material-icons.css';
</style>
