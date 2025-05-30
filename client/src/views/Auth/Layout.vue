<script setup lang="ts">
import {onMounted} from "vue";
import {useRouter} from "vue-router";
import {isAuthenticated} from "@/services/auth.ts";

let router = useRouter();

onMounted(() => {
  if (isAuthenticated()) {
    router.push("/");
  }
})
</script>

<template>
  <div class="flex justify-center">
    <div class="card w-96 bg-base-100 shadow-xl mt-20 mb-20">
      <div class="card-body relative">
        <h1 class="text-center text-xl">
          <RouterLink :to="{name: `home`}">Galvanico</RouterLink>
        </h1>

        <div class="transition-wrapper relative">
          <RouterView v-slot="{Component}">
            <transition name="fade">
              <component :is="Component"/>
            </transition>
          </RouterView>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
