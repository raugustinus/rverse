<template>
  <div class="input-group">
    <label>Parameter name</label>
    <input type="text" class="form-control" v-model="param.name">
  </div>
</template>

<script>
export default {
  name: 'config-param',
  data () {
    return {
      param: { name: '' }
    }
  },
  methods: {
    addParam () {
      console.log(`Adding param: ${JSON.stringify(this.param)}`)
    }
  }
}
</script>
