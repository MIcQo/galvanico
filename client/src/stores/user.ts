import {defineStore} from "pinia";
import {getCurrentUser} from "@/query/user.ts";
import {type UserCitiesResponse} from "@/query/userCity.ts";
import {ref} from "vue";

export const useCurrentUser = defineStore('currentUser', () => {
  const {data: currentUser, isLoading, error} = getCurrentUser()

  return {
    currentUser,
    isLoading,
    error
  }
})

export const useCurrentCity = defineStore('currentCity', () => {
  const currentCity = ref<UserCitiesResponse>()

  const setCurrentCity = (city: UserCitiesResponse) => {
    currentCity.value = city;
  }

  return {
    setCurrentCity,
    currentCity,
  }
})
