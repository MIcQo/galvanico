<script setup lang="ts">
import {IconBrandGithub} from '@tabler/icons-vue';
import {useCurrentUser} from "@/stores/user.ts";
import {storeToRefs} from "pinia";
import {useLogout} from "@/composables/logout.ts";

const {currentUser} = storeToRefs(useCurrentUser());

const {logout} = useLogout();

</script>

<template>
  <div class="container mx-auto">
    <div class="absolute top-[10px] right-2"><a href="https://github.com/MIcQo/galvanico"
                                                target="_blank">
      <IconBrandGithub stroke="2"/>
    </a></div>
    <div class="topbox w-full bg-base-100 h-[40px] flex justify-center items-center px-12 text-xs">
      <ul class="grid gap-3 grid-cols-6 w-full content-center text-center">
        <li class="font-bold"><a href="#">{{
            currentUser?.user.Username
          }}</a></li>
        <li><a target="_blank" href="https://discord.gg/EcYcJA3f">Discord</a></li>
        <li><a target="_blank" href="https://github.com/MIcQo/galvanico"> Github</a></li>
        <li><a target="_blank" href="https://github.com/MIcQo/galvanico/issues/new">Report
          problem</a></li>
        <li><a href="#">Item</a></li>
        <li class="font-bold">
          <button class="cursor-pointer" @click="logout">Logout</button>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>

</style>
