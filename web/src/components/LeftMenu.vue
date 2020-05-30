<template>
  <div>
    <h5>Menu</h5>
    <ul class="list-group-flush no-left-padding">
      <a v-for="item in menuItems" :href="item.href"
         class="list-group-item list-group-item-action"
         v-bind="menuItems" v-bind:key="item.href"
         v-on:click="selectItem(item)"
         v-bind:class="{ active : item.active }"
      >{{ item.name }}</a>
    </ul>
  </div>
</template>

<script>
import Headers from '@/components/Headers'
import Config from '@/components/Config'
import About from '@/components/About.vue'
import Overrides from '@/components/Overrides'

export default {
  name: 'left-menu',
  data () {
    return {
      menuItems: [
        { name: 'About', About, active: true },
        { name: 'Config', Config, active: false },
        { name: 'Headers', Headers, active: false },
        { name: 'Overrides', Overrides, active: false }
      ]
    }
  },
  methods: {
    selectItem: function (item) {
      this.$emit('select-component', item.name)
      this.menuItems.filter((v) => v.active).forEach((element) => { element.active = false })
      this.menuItems.filter((v) => (v.name === item.name)).forEach((element) => { element.active = true })
    }
  }
}
</script>

<style>
.no-left-padding {
  padding-left: 0;
}
</style>
