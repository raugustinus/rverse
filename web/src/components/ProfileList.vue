<template>
  <div>

    <div class="row">&nbsp;</div> <!-- spacer -->

    <div class="card">
      <div class="card-body">

        <h5 class="card-title">Current headers <span class="float-right"><button type="button" class="btn btn-secondary" v-on:click="listProfiles"><i class="ion-refresh"></i> reload</button></span></h5>
        <h6 class="card-subtitle mb-2 text-muted">This is the list of headers currently added to each request</h6>

        <div class="list-group">

          <a v-for="profile in profiles"
             v-bind="profiles"
             v-bind:key="profile.id"
             href="#" class="list-group-item list-group-item-action flex-column align-items-start active">
            <div class="d-flex w-100 justify-content-between">
              <h5 class="mb-1">{{ profile.name }}</h5>
              <small>3 days ago</small>
            </div>
            <p class="mb-1">Donec id elit non mi porta gravida at eget metus. Maecenas sed diam eget risus varius blandit.</p>
            <small>Donec id elit non mi porta.</small>
          </a>

          <a href="#" class="list-group-item list-group-item-action flex-column align-items-start">
            <div class="d-flex w-100 justify-content-between">
              <h5 class="mb-1">List group item heading</h5>
              <small class="text-muted">3 days ago</small>
            </div>
            <p class="mb-1">Donec id elit non mi porta gravida at eget metus. Maecenas sed diam eget risus varius blandit.</p>
            <small class="text-muted">Donec id elit non mi porta.</small>
          </a>
          <a href="#" class="list-group-item list-group-item-action flex-column align-items-start">
            <div class="d-flex w-100 justify-content-between">
              <h5 class="mb-1">List group item heading</h5>
              <small class="text-muted">3 days ago</small>
            </div>
            <p class="mb-1">Donec id elit non mi porta gravida at eget metus. Maecenas sed diam eget risus varius blandit.</p>
            <small class="text-muted">Donec id elit non mi porta.</small>
          </a>
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
  name: 'profile-list',
  data () {
    return {
      profiles: []
    }
  },
  methods: {
    listProfiles: function () {
      axiosInstance.get('/listProfiles').then((response) => {
        this.profiles = response.data
      }).catch((error) => {
        console.log(`Error fetching profiles: ${error}`)
        this.errors.push(error.message)
      })
    }
  },
  mounted () {
    this.listProfiles()
  }
}
</script>
