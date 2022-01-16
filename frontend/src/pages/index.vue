<script setup lang="ts">

import { useUserStore } from '~/stores/user'
import axios from "axios"

const user = useUserStore()

const listDataString = ref('')
const listData = ref([])

const getProjects = () => {
  let url = import.meta.env.VITE_PROTOCOL + "://" + import.meta.env.VITE_HOSTNAME +
    "/api/porter/me/project?token=" + user.token
  axios.get(url)
    .then((response) => {
      listDataString.value = JSON.stringify(response.data, null, "\t");
      listData.value = response.data;
      return response;
    })
    .catch((error) => console.log(error));
}

const setActiveProject = (id) => {
  user.setProject(id)
  localStorage.setItem("project", id)
}

setTimeout(getProjects, 50)

const tgInitLink = import.meta.env.VITE_PROTOCOL + "://" + import.meta.env.VITE_HOSTNAME +
  "/api/porter/me/project/new?token=" + user.token
const signInMsg = ref('')

if (user.name == "" || user.name == "guest") {
  signInMsg.value = 'SIGN IN FIRST!'
}

</script>

<template>
  <div class="m-auto text-left" style="width:500px; max-width: 100%;">
    <div class="text-left m-auto">
      <div class="m-2">
        {{ signInMsg }} To start a new project
        <a
          class="btn"
          :href="tgInitLink"
          target="_blank"
        >click here</a>. To
        remove any project call
        <code>`/remove`</code> in TG group.
      </div>
      <h2 class="m-auto mb-2 mt-8 text-lg text-left">&nbsp;Your Projects:</h2>
      <ul id="items" class="text-left m-auto">
        <li v-for="(item:any, index) in listData" :key="index">
          <div class="btn m-1 text-left break-all text-sm" @click="setActiveProject(item.ID)">
            {{ `Title: ${item.PrettyTitle.String}` }}
            <br />
            {{ `ID: ${item.ID}` }}
            <br />
            {{ `Secret: ${item.SecretKey}` }}
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<route lang="yaml">
meta:
  layout: default
</route>
