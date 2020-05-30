<template>
    <nav class="navbar navbar-expand-md">
      <a class="navbar-brand" href="#">
        <img src="../assets/logo7.png" height="100" alt="rverse.io logo">
      </a>
    </nav>
</template>

<script>
import $ from 'jquery'

export default {
  name: 'navbar',
  components: {
  },
  data () {
    return {
      url: '/',
      delay: 0,
      navlinks: []
    }
  },
  mounted () {
    this.readProxyConfig()
  },
  methods: {}
}

// might be able to move this to vuejs methods and use v-on:click
$(document).ready(function () {
  $('#menu-toggle').on('click', function (event) {
    event.preventDefault()
    let currentMarginLeft = $('#leftMenu').css('margin-left')
    if (currentMarginLeft.endsWith('px')) {
      currentMarginLeft = currentMarginLeft.substr(0, currentMarginLeft.length - 2)
    }
    let leftMenuWidth = $('#leftMenu').width() + 2 * 15 // 2 * 15 gutter
    if (currentMarginLeft < 0) {
      $('#leftMenu').css('margin-left', '0px')
    } else {
      $('#leftMenu').css('margin-left', `-${leftMenuWidth}px`)
    }
  })
})
</script>

<style>
.ionicon {
  font-size: 180%;
}
.hover-glow-orange:hover {
  color: orange !important;
}
.wider-input-nav {
  width: 600px;
}
.lean-navbar {
  padding: 0 15px !important;
  line-height: 1.0em;
}
.nav-item {
  padding: 0 10px;
}
.dropdown-item {
  color: white;
}
.dropdown-item:hover {
  color: #000000;
}
</style>
