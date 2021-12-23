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

if (name) {
  user.setNewName(name)
}

</script>

<template>
  <router-view />
</template>
