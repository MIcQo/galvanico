<script setup lang="ts">
import {useModal} from '@/stores/modal';
import UnderConstruction from "@/components/UnderConstruction.vue";

const modal = useModal();
</script>

<template>
  <div class="modal" role="dialog" :class="{'modal-open': modal.isOpen}">
    <div class="modal-box" :class="{'relative z-1000': modal.isOpen}">
      <h3 class="text-lg font-bold">Hello!</h3>
      <p class="py-4">
        This is {{ modal.buildingContext }} context
        <UnderConstruction/>
      </p>
      <div class="modal-action">
        <button @click="modal.close" class="btn">Yay!</button>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
