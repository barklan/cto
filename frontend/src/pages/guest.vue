<script setup lang="ts">
import { useUserStore } from '~/stores/user'

const route = useRoute();
const user = useUserStore()
const feedback = ref("")

feedback.value = "Redirecting..."

const token = route.query.token?.toString()
const project = route.query.project?.toString()
const name = "guest"

const router = useRouter()
const go = () => {
    if (token === undefined) {
        feedback.value = "No token provided."
        return
    }
    if (project === undefined) {
        feedback.value = "No project provided."
        return
    }
    user.setName(name)
    user.setToken(token)
    user.setProject(project)
    router.push(`/log`)
}

setTimeout(go, 50)

</script>

<template>
    <div>{{ feedback }}</div>
</template>
