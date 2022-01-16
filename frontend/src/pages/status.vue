<script setup lang="ts">
import { useUserStore } from '~/stores/user'
const user = useUserStore()
import axios from "axios"
const router = useRouter()

const listDataString = ref('')
const listData = ref([])

const getStatus = () => {
    let url = import.meta.env.VITE_PROTOCOL + "://" + import.meta.env.VITE_HOSTNAME +
        "/api/porter/me/project/" + user.project + "/status?token=" + user.token
    axios.get(url)
        .then((response) => {
            listDataString.value = JSON.stringify(response.data, null, "\t");
            listData.value = response.data;
            return response;
        })
        .catch((error) => console.log(error));
}

setTimeout(getStatus, 50)
setInterval(getStatus, 5000)

const goToExact = (key) => {
    router.push(`/log/exact?key=` + key)
}
</script>

<template>
    <div class="m-auto text-left" style="width:500px; max-width: 100%;">
        <h2 class="m-auto mb-2 mt-16 text-lg text-left">&nbsp;Recent Issues:</h2>
        <ul id="items" class="text-left m-auto">
            <li v-for="(item: any, index) in listData" :key="index">
                <div
                    class="m-auto btn m-2 text-left text-sm"
                    @click="goToExact(item.origin_badger_key)"
                >
                    {{ `Environment: ${item.hostname}` }}
                    <br />
                    {{ `Service: ${item.service}` }}
                    <br />
                    {{ `Last seen: ${item.last_seen}` }}
                    <br />
                    {{ `Total count: ${item.counter}` }}
                </div>
            </li>
        </ul>
    </div>
</template>
