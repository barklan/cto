<script setup lang="ts">
import { useUserStore } from '~/stores/user'

const route = useRoute();
const user = useUserStore()
const feedback = ref("")

feedback.value = "Redirecting..."

const token = route.query.token?.toString()
const name = route.query.name?.toString()
const project = route.query.project?.toString()

const router = useRouter()
const go = () => {
    if (token === undefined) {
        feedback.value = "No token provided."
        return
    }

    if (name === undefined) {
        feedback.value = "No name provided."
        return
    }

    if (project === undefined) {
        feedback.value = "No project provided."
        return
    }
    user.setName(name)
    user.setToken(token)
    user.setProject(project)
    localStorage.setItem("name", name)
    localStorage.setItem("token", token)
    localStorage.setItem("project", project)
    if (name == "guest") {
        router.push(`/log`)
    } else if (name != "") {
        router.push(`/`)
    }
}

setTimeout(go, 50)

</script>

<template>
    <div>{{ feedback }}</div>
</template>
