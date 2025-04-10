<script setup lang="ts">
import {RouterView} from 'vue-router'
import Alert from "@/components/Alert.vue";
</script>

<template>
  <Alert/>

  <RouterView v-slot="{Component}">
    <transition name="fade">
      <component :is="Component"/>
    </transition>
  </RouterView>
</template>

<style scoped>

</style>
