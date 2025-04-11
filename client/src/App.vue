<script setup lang="ts">
import {RouterView} from 'vue-router'
import Alert from "@/components/Alert.vue";
import {onMounted} from "vue";
import {useI18n} from "vue-i18n";

const i18n = useI18n({useScope: "global"})

onMounted(() => {

  const lang = localStorage.getItem("lang");
  if (lang) {
    i18n.locale.value = lang;
  }
})
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
