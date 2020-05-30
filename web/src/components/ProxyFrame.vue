<template>
  <div>
    <div class="row frame-container">
      <div class="col-xs-12">
        <iframe id="proxyFrame"
                v-if="loaded"
                width="100%"
                height="100%"
                class="proxyframe"
                name="myProxyFrame"
                frameborder="0">
        </iframe>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  name: 'proxy-frame',
  data () {
    return {
      loaded: false,
      url: 'http://localhost:8081/#/About'
    }
  },
  methods: {
    go () {
      let iframe = document.getElementById('proxyFrame')
      iframe.contentWindow.location.replace(this.url)
    }
  },
  mounted () {
    this.loaded = true
  }
}
</script>

<style>
.frame-container {
  z-index: 1;
}
.proxyframe {
  position: absolute;
  height: 100%
}
</style>
