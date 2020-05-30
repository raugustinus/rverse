<template>
  <div>
    <h5>Login</h5>
    <p>
      Select a profile to login with
    </p>

    <div class="card">
      <div class="card-body">

        <h5 class="card-title">Profiles</h5>
        <h6 class="card-subtitle mb-2 text-muted">List of profiles provided in login template</h6>

        <table class="table table-bordered table-condensed table-striped table-hover table-hover-cursor table-sm">
          <thead>
          <tr>
            <th>
              <a href="#">
                <div>name<span class="float-sm-right"><i class="ion-arrow-down-b"></i></span></div>
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
          <tr v-for="profile in profiles" v-bind="profiles" v-bind:key="profile.id">
            <td>{{ profile.name }}</td>
            <td>
              <a v-on:click="rmHeader(header)" data-toggle="tooltip" data-placement="top" title="Tooltip on top"><i class="ion-backspace-outline icon"></i></a>
              <a v-on:click="toggleHeader(header)"><i v-bind:class="[ header.active ? 'ion-android-checkbox-outline' : 'ion-android-checkbox-outline-blank' ]" class="icon"></i></a>
            </td>
          </tr>
          </tbody>

        </table>

      </div>
    </div>

  </div>
</template>

<script>
export default {
  name: 'login'
}
</script>
