<script setup lang="ts">
import { useUserStore } from '~/stores/user'

const user = useUserStore()
// const name = ref(user.name)

const { t } = useI18n()
</script>

<template>
<div>

</div>
</template>

<route lang="yaml">
meta:
  layout: home
</route>
