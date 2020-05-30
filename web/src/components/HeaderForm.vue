<template>
  <div class="input-group">
    <input type="text" class="form-control" v-model="header.key">
    <div class="input-group-prepend input-group-append">
      <span class="input-group-text">:</span>
    </div>
    <input type="text" class="form-control" v-model="header.value">
    <div class="input-group-append">
      <button v-on:click="addHeader" class="btn btn-outline-secondary" type="button"><i class="ion-plus-round"></i> Add</button>
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
  data () {
    return {
      header: { key: '', value: '', active: true }
    }
  },
  methods: {
    addHeader () {
      axiosInstance.post('/addheader', this.header).then(() => {
        console.log(`Added header OK!`)
        this.$emit('headers_updated') // notify parent that the list was changed
      }).catch((error) => { console.log(`Error adding header: ${error}`) })
    }
  },
  name: 'header-form'
}
</script>
