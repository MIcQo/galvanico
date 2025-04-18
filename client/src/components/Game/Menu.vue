<script setup lang="ts">
import UnderConstruction from "@/components/UnderConstruction.vue";
import {getUserCities} from "@/query/userCity.ts";
import {useCurrentCity} from "@/stores/user.ts";
import {computed, watch} from 'vue';


// Handle city change
function handleCityChange(event: Event) {
  const selectedId = (event.target as HTMLSelectElement).value;
  const city = userCities.value?.find(city => city.ID === selectedId);
  if (city) {
    setCurrentCity(city);
  }
}

const {data: userCities} = getUserCities();
const {currentCity, setCurrentCity} = useCurrentCity();
// Compute the selected city ID based on the current city
const selectedCityId = computed(() => currentCity?.ID || '');

watch(userCities, (newCities) => {
  if (newCities && newCities.length > 0 && !currentCity) {
    setCurrentCity(newCities[0]);
  }
}, {immediate: true});

</script>

<template>
  <div class="leftbox w-3/6 h-32 bg-base-200 absolute top-[40px] left-0 z-[400]">
    <select
      @change="handleCityChange"
      aria-label="Select city"
      class="select select-bordered w-full max-w-xs"
    >
      <option
        v-for="uc of userCities"
        :key="uc.ID"
        :value="uc.ID"
        :selected="selectedCityId === uc.ID"
      >
        {{ uc.Name }}
      </option>
    </select>
    <UnderConstruction/>
  </div>
</template>

<style scoped>

</style>
