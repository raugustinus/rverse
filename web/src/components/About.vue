<template>
  <div>
    <h3>About</h3>
    <p>
      RProxy is a reverse proxy for developers. It can augment requests with headers for use cases like authentication, authorization and such.
      Numerous organisations are using loadbalancers or http servers to handle authentication and then augment the requests to backends
      with security headers. Not all organisations have these loadbalancers setup on their development and or test environments.
      This reverse proxy can be used in these environments.
    </p>
  </div>
</template>

<script>
export default {
  name: 'About'
}
</script>
