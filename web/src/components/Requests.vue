<template>

  <div class="container">
    <div>
      <button v-on:click="fetchRequests">fetch headers</button>
    </div>
    <div>
      <table class="table">
        <thead>
        <tr>
          <th>key</th>
          <th>value</th>
          <th>action</th>
        </tr>
        </thead>

        <tbody>
        <tr v-for="request in requests" v-bind="requests" v-bind:key="request.key">
          <td>{{ request.key }}</td>
          <td>{{ request.value }}</td>
          <td>
            <button>remove</button>
          </td>
        </tr>
        </tbody>
      </table>

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
  name: 'Requests',
  data () {
    return {
      proxyheaders: [],
      message: 'testing'
    }
  },
  methods: {
    forward: function () {
      console.log('hello there.')
    },
    fetchRequests () {
      axiosInstance.get('/requests').then((response) => {
        console.log(`requests: ${JSON.stringify(response.data)}`)
        this.requests = response.data
      }).catch((error) => {
        console.log(`Error fetching headers: ${error}`)
      })
    }
  },
  mounted () {
    this.fetchHeaders()
  }
}
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
</style>
