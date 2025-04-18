<script setup lang="ts">
import UnderConstruction from "@/components/UnderConstruction.vue";
import {getUserCities, type UserCitiesResponse} from "@/query/userCity.ts";
import {useCurrentCity} from "@/stores/user.ts";

const {isPending, data: userCities, promise} = getUserCities();
const {currentCity, setCurrentCity} = useCurrentCity();

promise.value.then((r: UserCitiesResponse[]) => {
  setCurrentCity(r[0])
})
</script>

<template>
  <div class="leftbox w-3/6 h-32 bg-base-200 absolute top-[40px] left-0 z-400">
    <select name="" id="">
      <option v-for="uc of userCities" :value="uc.ID" :selected="uc.ID===currentCity?.ID"
              @click="setCurrentCity(uc)">{{ uc.Name }}
      </option>
    </select>
    <UnderConstruction/>
  </div>
</template>

<style scoped>

</style>
