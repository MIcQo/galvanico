<script setup lang="ts">
import {useModal} from "@/stores/modal.ts";
import {ref} from "vue";

const props = defineProps({
  position: {
    required: true,
    type: Number,
  },
  building: {
    required: true,
    type: String,
  },
  level: Number,
  upgrading: Boolean,
})

const modal = useModal();
const linkOpacity = ref(0);

const openModal = (building: string) => {
  modal.open(building)
}

const linkOnMouseOver = () => {
  linkOpacity.value = 1;

}
const linkOnMouseLeave = () => {
  linkOpacity.value = 0;
}

</script>

<template>

  <div class="building"
       :class="`position${props.position} ${props.building} ${props.upgrading ? 'constructionSite' : ''}`"
       :data-id="props.position">
    <div class="buildingItem">
      <div class="buildingimg img_pos animated"/>
      <div :style="{opacity: linkOpacity}" class="hover img_pos"/>
      <a @mouseover="linkOnMouseOver" @mouseleave="linkOnMouseLeave" href="#"
         @click.prevent="openModal(props.building)" class="hoverable"></a>
    </div>
  </div>

  <div v-if="props.upgrading" style="cursor:pointer;" class="buildingSpeedup timetofinish"
       :class="`position${props.position}`">
    <div class="before"/>
    <div class="buildingUpgradeIcon">47m 28s</div>
    <div class="buildingSpeedupButton" title="Skrátiť dobu výstavby"/>
    <div class="after"/>
  </div>
  <div v-else-if="!props.upgrading && props.building != 'land'" style="cursor:pointer;"
       class="cityScroll timetofinish"
       :class="`position${props.position}`">
    <div class="before"></div>
    <div class="green">{{ props.building }} (<span>{{ props.level || 1 }}</span>)</div>
    <div class="after"></div>
  </div>
</template>

<style>

</style>
