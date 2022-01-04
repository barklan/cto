<script setup lang="ts">
// import { isDark, toggleDark } from '~/composables'
import { useUserStore } from '~/stores/user'

const { t, availableLocales, locale } = useI18n()

const user = useUserStore()

const toggleLocales = () => {
  // change to some real logic
  const locales = availableLocales
  locale.value = locales[(locales.indexOf(locale.value) + 1) % locales.length]
}
</script>

<template>
  <div id="menu-container" class="w-screen fixed p-2 z-49 text-light-50">
  <nav id="menubar" w="full" class="text-lg z-50 mt-1">
    <router-link class="icon-btn mx-2" to="/" :title="t('button.home')">
       <codicon:home /><div style="top:-4px;" class="mx-4 relative inline-block">Home</div>
    </router-link>

    <router-link class="icon-btn mx-2" to="/log" :title="t('button.logs')">
       <carbon:data-view-alt /><div style="top:-4px;" class="mx-4 relative inline-block">Logs</div>
    </router-link>

    <router-link class="icon-btn mx-2" to="/status" :title="t('button.status')">
       <carbon:ai-status /><div style="top:-4px;" class="mx-4 relative inline-block">Status</div>
    </router-link>

    <router-link class="icon-btn mx-2" to="/signin" :title="t('button.signin')">
       <bytesize:sign-in /><div style="top:-4px;" class="mx-4 relative inline-block">Sign in</div>
    </router-link>

    <div style="top:-4px;" class="border-1 rounded-md px-2 inline-block icon-btn relative">{{user.name}}</div>

    <!-- <button class="icon-btn mx-2 !outline-none" :title="t('button.toggle_dark')" @click="toggleDark()">
      <carbon-moon v-if="isDark" />
      <carbon-sun v-else />
    </button> -->

    <a class="icon-btn mx-2" :title="t('button.toggle_langs')" @click="toggleLocales">
      <carbon-language />
    </a>

    <a class="icon-btn mx-2" rel="noreferrer" href="https://github.com/barklan/cto" target="_blank" title="GitHub">
      <carbon-logo-github />
    </a>
  </nav>
  </div>
</template>

<style>
#menu-container {
  top: 0px;
  background-color: #161b22;
  left: 0px;
  opacity: 1 !important;
}

</style>
