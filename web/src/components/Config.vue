<template>
  <div>
    <h3>Configuration</h3>
    <p>
      The configuration contains the url's that should be intercepted for authentication.
    </p>
    <div class="col-sm-8"> <!-- because of this col we have indent -->

      <!-- === Proxy server host & port configuration === -->
      <div class="card">
        <div class="card-body">
          <h5 class="card-title">Proxy host & port</h5>
          <h6 class="card-subtitle mb-2 text-muted">The hostname and port number the proxy will be running on.</h6>

          <form class="form">
            <div class="form-group row">
              <label for="proxyhost" class="col-sm-2 col-form-label">Proxy host</label>
              <div class="col-sm-6">
                <div class="input-group">
                  <input id="proxyhost" type="text" class="form-control" v-model="config.host" placeholder="http://localhost:3000/">
                  <input id="proxyport" type="text" class="form-control" v-model="config.port">
                  <div class="input-group-append">
                    <button v-on:click="setHost" class="btn btn-outline-secondary" type="button">save</button>
                  </div>
                </div>
              </div>
              <label class="col-form-label text-info"><small>* Requires a server restart.</small></label>
            </div>
          </form>
        </div>
      </div>

      <div class="row">&nbsp;</div> <!-- spacer -->

      <!-- === Proxypass configuration === -->
      <div class="card">
        <div class="card-body">
          <h5 class="card-title">Proxypass</h5>
          <h6 class="card-subtitle mb-2">To pass a request to an HTTP proxied server, the proxypass directive is specified.</h6>
          <form class="form">
            <div class="form-group row">
              <label for="proxypassHost" class="col-sm-2 col-form-label">Proxypass</label>
              <div class="col-sm-6">
                <div class="input-group">
                  <input id="proxypassHost" type="text" class="form-control" v-model="config.proxypass.host" placeholder="...">
                  <input type="text" class="form-control" v-model="config.proxypass.port" placeholder="...">
                  <div class="input-group-append">
                    <button v-on:click="setProxypass" class="btn btn-outline-secondary" type="button">save</button>
                  </div>
                </div>
              </div>
            </div>
          </form>
        </div>
      </div>

      <div class="row">&nbsp;</div> <!-- spacer -->

      <!-- === Delay configuration === -->
      <div class="card">
        <div class="card-body">
          <h5 class="card-title">Delay</h5>
          <h6 class="card-subtitle mb-2">For testing purposes you can set a delay on responses.</h6>
          <form class="form">
            <div class="form-group row">
              <label for="delay" class="col-sm-2 col-form-label">Delay</label>
              <div class="col-sm-6">
                <div class="input-group">
                  <input id="delay" type="number" class="form-control" v-model.number="config.delay" placeholder="...">
                  <div class="input-group-append">
                    <button v-on:click="setDelay" class="btn btn-outline-secondary" type="button">save</button>
                  </div>
                </div>
              </div>
            </div>
          </form>
        </div>
      </div>

      <div class="row">&nbsp;</div> <!-- spacer -->

      <!-- === Intercept Routes === -->
      <div class="card">
        <div class="card-body">
          <h5 class="card-title">Intercept routes</h5>
          <h6 class="card-subtitle mb-2">An application often has some authentication url's in order to use this reverse proxy to intercept
          those requests you provide those URLs here. RProxy will intercept these requests and provide the sample logins you have configured.</h6>

          <form class="form">
            <div class="form-group row">
              <label for="proxypassHost" class="col-sm-2 col-form-label">Proxypass</label>
              <div class="col-sm-6">
                <div class="input-group">
                  <input id="routeUrl" type="text" class="form-control" v-model="newRoute.url" placeholder="...">
                  <input id="routeParam" type="text" class="form-control" v-model="newRoute.param" placeholder="...">
                  <div class="input-group-append">
                    <button v-on:click="addRoute" class="btn btn-outline-secondary" type="button">save</button>
                  </div>
                </div>
              </div>
            </div>
          </form>

          <!-- === table of intercept routes === -->
          <table class="table table-bordered table-condensed table-striped table-hover table-hover-cursor table-sm">
            <thead>
            <tr>
              <th>
                <a href="#">
                  <div>Route<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
                </a>
              </th>
              <th>
                <a href="#">
                  <div>Redirect parameter<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
                </a>
              </th>
              <th><a href="#">
                <div>Actions<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
              </a>
              </th>
            </tr>
            </thead>

            <tbody>
            <tr v-for="route in routes" v-bind="routes" v-bind:key="route.url">
              <td>{{ route.path }}</td>
              <td>{{ route.parameters }}</td>
              <td>
                <a v-on:click="rmRoute(route)" data-toggle="tooltip" data-placement="top" title="Tooltip on top"><i class="ion-backspace-outline icon"></i></a>
              </td>
            </tr>
            </tbody>

          </table>

        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

let axiosInstance = axios.create({
  baseURL: '/rservice',
  timeout: 5000
})

export default {
  name: 'config',
  data () {
    return {
      newRoute: { url: '', param: '' },
      config: { delay: 0, host: 'localhost', port: 3000, routes: [], proxypass: { host: '', port: '' } },
      routes: []
    }
  },
  mounted () {
    this.readConfig()
  },
  methods: {
    readConfig () {
      this.fetchRoutes()
      this.readProxypass()
    },
    readProxypass () {
      axiosInstance.get('/readProxypass').then((response) => {
        console.log(`Read proxypass: ${JSON.stringify(response.data)}`)
        this.config.proxypass = response.data
      }).catch((error) => {
        console.log(`Error reading proxypass: ${error.message}`)
      })
    },
    fetchRoutes () {
      this.routes = []
      axiosInstance.get('/listRoutes').then((response) => {
        this.routes = response.data
      }).catch((error) => {
        console.log(`Error fetching headers: ${error}`)
        // this.errors.push(error.message)
      })
    },
    saveConfig () {
      console.log(`Saving rproxy configuration: ${JSON.stringify(this.config)}`)
      this.readConfig()
    },
    setProxypass () {
      console.log(`Setting proxypass: ${JSON.stringify(this.config.proxypass)}`)
      axiosInstance.post('/setProxypass', { host: this.config.proxypass.host, port: this.config.proxypass.port }).then(() => {
        console.log(`Set proxypass: ${JSON.stringify(this.config.proxypass)}`)
      }).catch((error) => {
        console.log(`Error removing route: ${error.message}`)
      })
    },
    setDelay () {
      axiosInstance.post('/setDelay', {delay: this.config.delay}).then(() => {
        console.log(`Setting delay: ${JSON.stringify(this.config.delay)}`)
      }).catch((error) => {
        console.log(`Error removing route: ${error.message}`)
      })
    },
    addRoute () {
      console.log(`Adding route: ${this.newRoute}`)
      axiosInstance.post('/addRoute', { path: this.newRoute.url, parameters: [this.newRoute.param] }).then(() => {
        console.log(`Added route: ${JSON.stringify(this.newRoute)}`)
      }).catch((error) => {
        console.log(`Error removing route: ${error.message}`)
      })
    },
    rmRoute (route) {
      console.log(`Removing route: ${route}`)
    },
    setHost () {
      console.log(`Setting host to: ${this.config.host} and port: ${this.config.port}`)
      axiosInstance.post('/setHost', { host: this.config.host, port: this.config.port }).then(() => {
        console.log(`Success setting host & port.`)
      }).catch((error) => {
        console.log(`Error setting host: ${error.message}`)
      })
    }
  }
}
</script>
