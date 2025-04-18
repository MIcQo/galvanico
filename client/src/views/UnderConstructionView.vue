<script setup lang="ts">
import {useLogout} from "@/composables/logout.ts";
import {IconBarrierBlock} from "@tabler/icons-vue";

const {logout} = useLogout();
</script>

<template>
  <div class="hero bg-base-200 min-h-screen">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <h1 class="text-2xl font-bold max-w-18 mx-auto">
          <IconBarrierBlock stroke={1} width="auto" height="auto"/>
        </h1>
        <p class="py-6">
          {{ $t('game.underConstruction') }}
        </p>
        <button @click="logout"
                class="btn btn-ghost">
          {{ $t('auth.logout') }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
