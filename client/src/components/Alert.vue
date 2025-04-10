<script setup lang="ts">
import {AlertType, useAlert} from "@/stores/alert.ts";
import {computed} from "vue";

const alert = useAlert();

const alertClass = computed(() => {
  return {
    'alert-success': alert.alertType === AlertType.success,
    'alert-warning': alert.alertType === AlertType.warning,
    'alert-error': alert.alertType === AlertType.danger,
  }
})
</script>

<template>
  <div
    class="absolute top-0 w-auto max-w-md" style="left: 50%; transform: translateX(-50%)">
    <transition name="slide">
      <div v-if="alert.isOpen" role="alert" class="alert mt-5" :class="alertClass"
           aria-live="assertive">
        <svg v-if="alert.alertType === AlertType.success" xmlns="http://www.w3.org/2000/svg"
             class="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <svg v-if="alert.alertType === AlertType.warning" xmlns="http://www.w3.org/2000/svg"
             class="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
        </svg>
        <svg v-if="alert.alertType === AlertType.danger" xmlns="http://www.w3.org/2000/svg"
             class="h-6 w-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <span>{{ alert.message }}</span>
      </div>
    </transition>
  </div>
</template>

<style scoped>

</style>
