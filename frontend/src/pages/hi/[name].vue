<script setup lang="ts">
const props = defineProps<{ name: string }>()
const router = useRouter()
</script>

<template>
  <div>
    <p class="text-4xl">
      <carbon-pedestrian class="inline-block" />
    </p>
    <p>
      Your query: {{ props.name }}
    </p>

    <div>
      <button
        class="btn m-3 text-sm mt-8"
        @click="router.back()"
      >
        Back
      </button>
    </div>
  </div>
</template>
