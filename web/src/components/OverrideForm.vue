<template>
  <form>
    <div class="form-group">
      <label for="overrideUrl">Pattern</label>
      <input id="overrideUrl" type="text" class="form-control" v-model="override.pattern">
    </div>
    <div class="form-group">
      <label for="responseContentType">Response Content-Type</label>
      <input id="responseContentType" type="text" class="form-control" v-model="override.contentType">
    </div>
    <div class="form-group">
      <label for="responseBody">Response body</label>
      <textarea id="responseBody" cols="200" rows="10" class="form-control" v-model="override.responseBody"></textarea>
    </div>
    <div class="form-group">
      <label for="overrideResponseCode">Response code</label>
      <input id="overrideResponseCode" type="number" class="form-control" v-model.number="override.responsecode">
    </div>
    <div class="form-group">
      <button v-on:click="addOverride" class="btn btn-outline-secondary" type="button"><i class="ion-plus-round"></i> Add</button>
    </div>
  </form>
</template>

<script>
import axios from 'axios'

let axiosInstance = axios.create({
  baseURL: '/rservice',
  timeout: 5000
})

export default {
  data () {
    return {
      errors: [],
      override: { pattern: '', message: 'default', responsecode: 0, contentType: '', responseBody: '', active: true }
    }
  },
  methods: {
    addOverride () {
      console.log(`override to add: ${JSON.stringify(this.override)}`)
      axiosInstance.post('/addOverride', this.override).then(() => {
        console.log(`Added override OK!`)
        this.$emit('overrides_updated') // notify parent that the list was changed
      }).catch((error) => {
        console.log(`Error adding override: ${error}`)
        this.errors = error.message
      })
    }
  },
  name: 'override-form'
}
</script>
