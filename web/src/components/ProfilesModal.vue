<template>
  <div class="modal" tabindex="-1" role="dialog">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Modal title</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body">
          <p>Modal body text goes here.</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-primary">Save changes</button>
          <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'profiles-modal',
  data () {
    return {
      url: '/publiek/',
      delay: 0,
      delayRange: 0,
      navlinks: []
    }
  },
  methods: {
    frameUrl () {
      console.log(`Framing URL: ${this.url}`)
      this.$router.push('/frame', this.frameGoTo())
      // $('#proxyFrame').attr('src', '/Config')
    },
    frameGoTo () {
      console.log(`hello..`)
      let iframe = document.getElementById('proxyFrame')
      iframe.contentWindow.location.replace(this.url)
    }
  }
}

// might be able to move this to vuejs methods and use v-on:click
$(document).ready(function () {
  $('#menu-toggle').on('click', function (event) {
    event.preventDefault()
    let currentMarginLeft = $('#leftMenu').css('margin-left')
    if (currentMarginLeft.endsWith('px')) {
      currentMarginLeft = currentMarginLeft.substr(0, currentMarginLeft.length - 2)
    }
    let leftMenuWidth = $('#leftMenu').width() + 2 * 15 // 2 * 15 gutter
    if (currentMarginLeft < 0) {
      $('#leftMenu').css('margin-left', '0px')
    } else {
      $('#leftMenu').css('margin-left', `-${leftMenuWidth}px`)
    }
  })
})
</script>
