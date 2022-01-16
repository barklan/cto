<script setup lang="ts">
import { useUserStore } from '~/stores/user'
const user = useUserStore()
import axios from "axios"
const router = useRouter()

const issues = ref([])
const envs = ref([])
const envToServices = ref(new Map())

const getList = async (path: string, extraQuery: string) => {
  let url = import.meta.env.VITE_PROTOCOL + "://" + import.meta.env.VITE_HOSTNAME +
    "/api/porter/project/" + user.project + path + "?token=" + user.token +
    extraQuery
  let resp = await axios.get(url)
    .then(response => {
      return response
    })
    .catch((error) => console.log(error));
  return resp
}

const getStatus = async () => {
  let response = await getList("/issues", "")
  issues.value = response?.data;

  response = await getList("/environments", "")
  envs.value = response?.data;

  for (var env in envs.value) {
    response = await getList("/services", "&env=" + env)
    envToServices.value.set(env, response?.data);
  }
}

setTimeout(getStatus, 50)

const goToExact = (key) => {
  router.push(`/log/exact?key=` + key)
}
</script>

<template>
  <div class="m-auto text-left" style="width:500px; max-width: 100%;">
    <h2 class="m-auto mb-2 mt-16 text-lg text-left">&nbsp;Recent Issues:</h2>
    <ul id="items" class="text-left m-auto">
      <li v-for="(item: any, index) in issues" :key="index">
        <div
          class="m-auto btn m-2 text-left text-sm break-all"
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

    <h2 class="m-auto mb-2 mt-16 text-lg text-left">&nbsp;Envs and Services:</h2>
    <ul id="items" class="text-left m-auto">
      <li v-for="[env, services] in envToServices">
        <div
          class="m-auto m-2 text-left text-sm break-all"
        >
          {{ env }} =>
            <ul class="ml-20">
              <li v-for="(service, key) in services" :key="key">
                {{key}}
              </li>
            </ul>
        </div>
      </li>
    </ul>
  </div>
</template>
