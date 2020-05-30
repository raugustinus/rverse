<template>

  <div :id="id" class="modal" tabindex="-1" role="dialog">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Headers</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>

        <div class="modal-body">

          <div class="card">
            <div class="card-body">

              <!-- === modal body === -->
              <div v-if="hasError()" class="alert alert-danger" role="alert">
                Error fetching headers: {{ listHeadersError }}
              </div>

              <h5 class="card-title">Current headers <span class="float-right"><button type="button" class="btn btn-secondary" v-on:click="fetchHeaders"><i class="ion-refresh"></i> reload</button></span></h5>
              <h6 class="card-subtitle mb-2 text-muted">The list of headers available, see the action column for their status.</h6>

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

        <div class="modal-footer">
          <button type="button" class="btn btn-primary" v-on:click="save">Save changes</button>
          <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>

</template>

<script>
import axios from 'axios'
import HeaderForm from '@/components/HeaderForm.vue'

let axiosInstance = axios.create({
  baseURL: '/rservice',
  timeout: 5000
})

export default {
  components: {
    HeaderForm
  },
  name: 'headers-modal',
  props: ['id'],
  data () {
    return {
      newHeader: { headerKey: '', headerValue: '' },
      proxyheaders: [],
      listHeadersError: []
    }
  },
  methods: {
    hasError: function () {
      return this.listHeadersError.length > 0
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
    },
    saveHeader(header) {
      console.log(`Saving header: ${JSON.stringify(header)}`)
      axiosInstance.post('/saveHeader', header).then(() => {
        console.log(`Saved header!`);
      }).catch( error => {
        console.log(`Error saving header: ${error.messages}`)
      })
    }
  },
  mounted () {
    console.log(`Mounted HeadersModal...`)
    this.fetchHeaders()
  }
}
</script>
