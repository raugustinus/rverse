<template>

  <div>

    <h3>Headers</h3>

    <p>
      This is the headers management page for the rproxy reverse proxy. You can manage the headers that will augment the
      http requests sent to webapplications and record scenarios which can be played back. This can be usefull in testing
      backend application services.
    </p>

    <div v-if="hasError()" class="alert alert-danger" role="alert">
      Error fetching headers: {{ listHeadersError }}
    </div>

    <div class="row">&nbsp;</div> <!-- spacer -->

    <div class="card">
      <div class="card-body">

        <h5 class="card-title">Profiles</h5>
        <h6 class="card-subtitle mb-2 text-muted">List of profiles provided in login template</h6>

        <profile-list />

      </div>
    </div>

    <div class="row">&nbsp;</div> <!-- spacer -->

    <div class="card">
      <div class="card-body">

        <div v-if="hasError()" class="alert alert-danger" role="alert">
          Error fetching headers: {{ listHeadersError }}
        </div>

        <h5 class="card-title">Current headers <span class="float-right"><button type="button" class="btn btn-secondary" v-on:click="fetchHeaders"><i class="ion-refresh"></i> reload</button></span></h5>
        <h6 class="card-subtitle mb-2 text-muted">This is the list of headers currently added to each request</h6>

        <table class="table table-bordered table-condensed table-striped table-hover table-hover-cursor table-sm">
          <thead>
          <tr>
            <th>
              <a href="#">
                <div>key<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
              </a>
            </th>
            <th>
              <a href="#">
                <div>value<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
              </a>
            </th>
            <th><a href="#">
              <div>action<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
            </a>
            </th>
          </tr>
          </thead>

          <tbody>
          <tr v-for="header in proxyheaders" v-bind="proxyheaders" v-bind:key="header.id">
            <td>{{ header.key }}</td>
            <td>
              <textarea class="form-control" v-model="header.value"></textarea>
            </td>
            <td>
              <a v-on:click="rmHeader(header)" data-toggle="tooltip" data-placement="top" title="Tooltip on top"><i class="ion-backspace-outline icon"></i></a>
              <a v-on:click="toggleHeader(header)"><i v-bind:class="[ header.active ? 'ion-android-checkbox-outline' : 'ion-android-checkbox-outline-blank' ]" class="icon"></i></a>
            </td>
          </tr>
          </tbody>

        </table>

      </div>
    </div>

    <div class="row">&nbsp;</div> <!-- spacer -->

    <div class="card">
      <div class="card-body">
        <h5 class="card-title">Add header</h5>
        <h6 class="card-subtitle mb-2 text-muted">Use this form to augment requests with the headers you add here</h6>
        <form>
          <div class="form-group">
            <label>Header (key:value)</label>
            <header-form v-on:headers_updated="fetchHeaders"/>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>

<script>
import $ from 'jquery'
import axios from 'axios'
import HeaderForm from './HeaderForm.vue'
import '../../node_modules/bootstrap/js/dist/tooltip'
import ProfileList from '@/components/ProfileList'

let axiosInstance = axios.create({
  baseURL: '/rservice',
  timeout: 5000
})

export default {
  components: {
    ProfileList,
    HeaderForm
  },
  data () {
    return {
      newHeader: { headerKey: '', headerValue: '' },
      proxyheaders: [],
      message: 'testing',
      listHeadersError: []
    }
  },
  methods: {
    hasError: function () {
      return this.listHeadersError.length > 0
    },
    forward: function () {
      console.log('hello there.')
    },
    fetchHeaders () {
      this.listHeadersError = []
      axiosInstance.get('/listHeaders').then((response) => {
        this.proxyheaders = response.data
      }).catch((error) => {
        console.log(`Error fetching headers: ${error}`)
        this.listHeadersError.push(error.message)
      })
    },
    rmHeader (header) {
      axiosInstance.delete('/removeheader', { data: header }).then(() => {
        console.log(`removed header for key: ${header.key}`)
        this.fetchHeaders()
      }).catch((error) => {
        console.log(`Error removing header with key ${header.key} : ${error}`)
      })
    },
    toggleHeader (header) {
      console.log(`Toggle header: ${JSON.stringify(header)}`)
      axiosInstance.post('/toggleHeader', header).then(() => {
        console.log(`Toggled header for key: ${header.key}`)
        this.fetchHeaders()
      }).catch((error) => {
        console.log(`Error toggling header with key ${header.key} : ${error}`)
      })
    }
  },
  mounted () {
    this.fetchHeaders()
  }
}

$(function () {
  $('[data-toggle="tooltip"]').tooltip()
})
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1, h2 {
    font-weight: normal;
  }
  ul {
    list-style-type: none;
    padding: 0;
  }
  li {
    display: inline-block;
    margin: 0 10px;
  }
  a {
    color: #42b983;
  }
  .table-hover-cursor tbody {
    cursor: pointer;
  }
  .icon {
    padding: 0 8px;
    font-size: 140%;
    color: #8ac9d4;
  }
  .icon:hover {
    color: #45a450;
  }
</style>
