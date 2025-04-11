<script setup lang="ts">
import {useRouter} from "vue-router";
import {removeToken} from "@/services/auth.ts";
import {useAlert} from "@/stores/alert.ts";
import {useI18n} from "vue-i18n";

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
  <div class="hero bg-base-200 min-h-screen">
    <div class="hero-content text-center">
      <div class="max-w-md">
        <h1 class="text-2xl font-bold max-w-18 mx-auto">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path
              d="M4 7m0 1a1 1 0 0 1 1 -1h14a1 1 0 0 1 1 1v7a1 1 0 0 1 -1 1h-14a1 1 0 0 1 -1 -1z"/>
            <path d="M7 16v4"/>
            <path d="M7.5 16l9 -9"/>
            <path d="M13.5 16l6.5 -6.5"/>
            <path d="M4 13.5l6.5 -6.5"/>
            <path d="M17 16v4"/>
            <path d="M5 20h4"/>
            <path d="M15 20h4"/>
            <path d="M17 7v-2"/>
            <path d="M7 7v-2"/>
          </svg>
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
