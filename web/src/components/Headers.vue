<template>
  <div>
    <div>
      <div v-if="hasError()" class="alert alert-danger" role="alert">
        Error fetching headers: {{ listHeadersError }}
      </div>

      <h5 class="card-title">Current headers <i class="material-icons btn-icon" v-on:click="fetchHeaders">refresh</i></h5>
      <h6 class="card-subtitle mb-2 text-muted">The list of headers available, see the action column for their status.</h6>

      <div class="table-responsive">
        <table class="table">
          <thead>
          <tr>
            <th>
              <a href="#">
                <div>Key<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
              </a>
            </th>
            <th>
              <a href="#"><div>Action</div></a>
            </th>
            <th >
              <a href="#">
                <div>Value<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
              </a>
            </th>
          </tr>
          </thead>

          <tbody>
          <tr v-for="header in proxyheaders" v-bind="proxyheaders" v-bind:key="header.id" class="rowcursor">
            <td class="actions">
              <i class="material-icons btn-icon" data-toggle="tooltip" data-placement="top" title="Toggle header on/off" v-bind:class="header.active ? 'btn-active' : ''" v-on:click="toggleHeader(header)" >{{ header.active ? 'check_box' : 'check_box_outline_blank' }}</i>
              <i class="material-icons btn-icon" data-toggle="tooltip" data-placement="top" title="Remove header" v-on:click="rmHeader(header)">delete</i>
              <i class="material-icons btn-icon" v-on:click="editHeader(header)" data-toggle="modal" data-target="#editModal">edit</i>
            </td>
            <td v-bind:class="[ header.active ? 'header-on' : '' ]">{{ header.key }}</td>
            <td>
              <span class="header-text" v-on:click="editHeader(header)" data-toggle="modal" data-target="#editModal">
                {{header.value}}
              </span>
            </td>
          </tr>
          </tbody>

        </table>
      </div>

      <br/>

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

      <button v-on:click="ttips()" class="btn btn-primary">activate tooltips!</button>

    </div>

    <div id="editModal" class="modal" tabindex="-1" role="dialog">
      <div class="modal-dialog" role="document">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{selectedHeader.key}}</h5>
            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
              <span aria-hidden="true">&times;</span>
            </button>
          </div>
          <div class="modal-body">
            <textarea class="form-control" rows="10" v-model="selectedHeader.value"></textarea>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary" data-dismiss="modal" v-on:click="updateHeader()">Save</button>
            <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import axios from 'axios'
import HeaderForm from '@/components/HeaderForm'
import $ from 'jquery'
import 'bootstrap/js/dist/tooltip'

let axiosInstance = axios.create({
  baseURL: '/rservice',
  timeout: 5000
})

export default {
  components: {
    HeaderForm
  },
  name: 'Headers',
  props: ['id'],
  data () {
    return {
      newHeader: {headerKey: '', headerValue: ''},
      selectedHeader: {},
      proxyheaders: [],
      listHeadersError: []
    }
  },
  methods: {
    hasError: function () {
      return this.listHeadersError.length > 0
    },
    editHeader (header) {
      this.selectedHeader = header
    },
    ttips () {
      console.log('hello there')
      $('[data-toggle="tooltip"]').tooltip()
    },
    fetchHeaders () {
      this.listHeadersError = []
      axiosInstance.get('/listHeaders').then((response) => {
        this.proxyheaders = response.data
        $(function () {
          $('[data-toggle="tooltip"]').tooltip()
        })
      }).catch((error) => {
        console.log(`Error fetching headers: ${error}`)
        this.listHeadersError.push(error.message)
      })
    },
    rmHeader (header) {
      axiosInstance.delete('/removeheader', {data: header}).then(() => {
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
    },
    updateHeader () {
      console.log(`Updating header: ${JSON.stringify(this.selectedHeader)}`)
      axiosInstance.post('/updateHeader', this.selectedHeader).then(() => {
        console.log(`Updated header!`)
      }).catch((error) => {
        console.log(`Error saving header: ${error.message}`)
      })
    }
  },
  mounted () {
    this.fetchHeaders()
    $(function () {
      $('[data-toggle="tooltip"]').tooltip()
    })
  }
}
</script>

<style lang="css">
.actions a {
  font-weight: 100 !important;
  font-size: 200%;
}
.actions a:hover {
  color: red !important;
}
.header-text {
  width: 500px;
  display: inline-block;
  word-wrap:break-word;
  display:inline-block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.btn-icon {
  padding: 5px;
  font-size: 1.4em;
}
.btn-active {
  color: #007bff;
}
.btn-icon:hover {
  cursor: pointer;
  background: #efefef;
  /*color: #007bff;*/
}
</style>
