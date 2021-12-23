<script setup lang="ts">
import { useUserStore } from '~/stores/user'
const user = useUserStore()

</script>

<template>
    <div>{{ feedback }}</div>
</template>
