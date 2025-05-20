<script setup lang="ts">
import {IconBarrierBlock} from '@tabler/icons-vue'
</script>

<template>
  <div class="max-w-md text-center w-full mx-auto">
    <h1 class="text-2xl font-bold max-w-18 mx-auto">
      <IconBarrierBlock stroke={1} width="auto" height="auto"/>
    </h1>
    <p class="py-6">
      {{ $t('game.underConstruction') }}
    </p>
  </div>
</template>

<style scoped>

</style>
