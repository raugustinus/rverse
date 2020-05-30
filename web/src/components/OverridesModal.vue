<template>

  <div :id="id" class="modal" tabindex="-1" role="dialog">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Overrides - error interceptors</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>

        <div class="modal-body">

          <div class="card">
            <div class="card-body">

              <!-- === modal body === -->
              <div v-if="hasError()" class="alert alert-danger" role="alert">
                Error fetching overrides: {{ errors }}
              </div>

              <h5 class="card-title">Current overrides <span class="float-right"><button type="button" class="btn btn-secondary" v-on:click="fetchOverrides"><i class="ion-refresh"></i> reload</button></span></h5>
              <h6 class="card-subtitle mb-2 text-muted">The list of headers available, see the action column for their status.</h6>

              <table class="table table-bordered table-condensed table-striped table-hover table-hover-cursor table-sm">
                <thead>
                <tr>
                  <th>
                    <a href="#">
                      <div>Pattern<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
                    </a>
                  </th>
                  <th>
                    <a href="#">
                      <div>Message<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
                    </a>
                  </th>
                  <th>
                    <a href="#">
                      <div>Content-Type<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
                    </a>
                  </th>
                  <th>
                    <a href="#">
                    <div>Response code<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
                  </a>
                  </th>
                  <th>
                    actions
                  </th>
                </tr>
                </thead>

                <tbody>
                <tr v-for="override in overrides" v-bind="overrides" v-bind:key="override.pattern">
                  <td>{{ override.pattern }}</td>
                  <td>{{ override.message }}</td>
                  <td>{{ override.contentType }}</td>
                  <td>{{ override.responsecode }}</td>
                  <td>
                    <button type="button" v-on:click="rmOverride(override)" class="btn btn-secondary"><i class="ion-backspace-outline icon"></i> remove</button>
                    <button type="button" v-on:click="toggleOverride(override)" class="btn btn-secondary"><i v-bind:class="[ override.active ? 'ion-android-checkbox-outline' : 'ion-android-checkbox-outline-blank' ]" class="icon"></i></button>
                  </td>
                </tr>
                </tbody>

              </table>

            </div>
          </div>

          <div class="row">&nbsp;</div> <!-- spacer -->

          <div class="card">
            <div class="card-body">
              <h5 class="card-title">Add override</h5>
              <h6 class="card-subtitle mb-2 text-muted">Use this form to intercept URL's and provide a static response (like an error for instance)</h6>
              <override-form v-on:overrides_updated="fetchOverrides"/>
            </div>
          </div>

        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>

</template>

<script>
import axios from 'axios'
import OverrideForm from '@/components/OverrideForm.vue'

let axiosInstance = axios.create({
  baseURL: '/rservice',
  timeout: 5000
})

export default {
  components: {
    OverrideForm
  },
  name: 'overrides-modal',
  props: ['id'],
  data () {
    return {
      overrides: [],
      errors: []
    }
  },
  methods: {
    hasError: function () {
      return this.errors.length > 0
    },
    fetchOverrides () {
      axiosInstance.get('/listOverrides').then((response) => {
        this.overrides = response.data
        console.log(`listed overrides from server: ${response.data}`)
      }).catch((error) => {
        console.log(`Error fetching overrides: ${error}`)
        this.errors = error.message
      })
    },
    rmOverride (override) {
      axiosInstance.delete('/rmOverride', { data: override }).then(() => {
        console.log(`Removed override for url: ${override.url}`)
        this.fetchOverrides()
      }).catch((error) => {
        console.log(`Error removing override with url: ${override.url} : ${error}`)
      })
    },
    toggleOverride (override) {
      console.log(`Toggle override: ${JSON.stringify(override)}`)
      axiosInstance.post('/toggleOverride', override).then(() => {
        console.log(`Toggled override for url: ${override.url}`)
        this.fetchOverrides()
      }).catch((error) => {
        console.log(`Error toggling override with url: ${override.url} : ${error}`)
      })
    }
  },
  mounted () {
    console.log(`Mounted OverrideModal...`)
    this.fetchOverrides()
  }
}
</script>
