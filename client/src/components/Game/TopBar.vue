<script setup lang="ts">
import {useQuery} from "@tanstack/vue-query";
import {defaultInstance} from "@/services/api.ts";
import {useRouter} from "vue-router";
import {useAlert} from "@/stores/alert.ts";
import {useI18n} from "vue-i18n";
import {removeToken} from "@/services/auth.ts";
import {IconBrandGithub} from '@tabler/icons-vue';

interface UserResponse {
  user: User
}

interface User {
  Username: string;
  Email: string;
}

const getUser = (): Promise<UserResponse> => {
  return defaultInstance.get('api/user').json();
}
const {data, error} = useQuery({
  queryKey: ['user'],
  queryFn: getUser,
})

const router = useRouter();
const alert = useAlert();
const {t} = useI18n();

const logout = () => {
  removeToken()
  alert.open(t("auth.alert.logout"));
  router.push({name: "auth.login"})
}

</script>

<template>
  <div class="container mx-auto">
    <div class="absolute top-[10px] right-2"><a href="#">
      <IconBrandGithub stroke="2"/>
    </a></div>
    <div class="topbox w-full bg-base-100 h-[40px] flex justify-center items-center px-12 text-xs">
      <ul class="grid gap-3 grid-cols-6 w-full content-center text-center">
        <li class="font-bold"><a href="https://github.com/MIcQo/galvanico">{{
            data?.user.Username
          }}</a></li>
        <li><a href="https://discord.gg/EcYcJA3f">Discord</a></li>
        <li><a href="https://github.com/MIcQo/galvanico"> Github</a></li>
        <li><a href="https://github.com/MIcQo/galvanico/issues/new">Report problem</a></li>
        <li><a href="#">Item</a></li>
        <li class="font-bold">
          <button class="cursor-pointer" @click="logout">Logout</button>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>

</style>
