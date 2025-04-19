<script setup lang="ts">

import {onMounted} from "vue";
import {isAuthenticated} from "@/services/auth.ts";
import {useRouter} from "vue-router";
import Modal from "@/components/Modal.vue";
import Menu from "@/components/Game/Menu.vue";
import Advisors from "@/components/Game/Advisors.vue";
import TopBar from "@/components/Game/TopBar.vue";

const router = useRouter();
onMounted(() => {
  if (!isAuthenticated()) {
    router.push({name: "auth.login"});
  }
})

</script>

<template>
  <div id="container" class="w-screen h-screen">
    <TopBar/>
    <Menu/>
    <Advisors/>

    <RouterView v-slot="{Component}">
      <transition>
        <component :is="Component"/>
      </transition>
    </RouterView>
    <Modal/>
  </div>
</template>

<style scoped>

</style>
