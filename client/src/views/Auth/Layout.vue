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
      <div class="card-body">
        <h1>Galvanico</h1>

        <RouterView v-slot="{Component}">
          <transition name="fade">
            <component :is="Component"/>
          </transition>
        </RouterView>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
