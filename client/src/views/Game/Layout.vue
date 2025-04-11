<script setup lang="ts">

import {onMounted} from "vue";
import {isAuthenticated} from "@/services/auth.ts";
import {useRouter} from "vue-router";

const router = useRouter();
onMounted(() => {
  if (!isAuthenticated()) {
    router.push({name: "auth.login"});
  }
})

</script>

<template>
  <RouterView v-slot="{Component}">
    <transition>
      <component :is="Component"/>
    </transition>
  </RouterView>
</template>

<style scoped>

</style>
