<script setup lang="ts">

import { useUserStore } from '~/stores/user'

// https://github.com/vueuse/head
// you can use this to manipulate the document head in any components,
// they will be rendered correctly in the html results with vite-ssg
useHead({
  title: 'CTO',
  meta: [
    { name: 'description', content: 'CTO - panel' },
  ],
})

const user = useUserStore()
const name = localStorage.getItem("name");
const token = localStorage.getItem("token")
const project = localStorage.getItem("project")

if (name) {
  user.setName(name)

  if (name == "guest" && token && project) {
    user.setToken(token)
    user.setProject(project)
  }
}

</script>

<template>
  <router-view />
</template>
