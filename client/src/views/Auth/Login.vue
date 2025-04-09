<script setup lang="ts">
import {ref} from "vue";
import {defaultInstance} from "@/services/api.ts";
import {AlertType, useAlert} from "@/stores/alert.ts";
import {useI18n} from "vue-i18n";
import {useRouter} from "vue-router";

const alert = useAlert();
const {t} = useI18n();
const router = useRouter();
const email = ref("");
const password = ref("");

interface tokenResponse {
  token: string
}

const login = async () => {
  const user = await defaultInstance.post('auth/login', {
    json: {username: email.value, password: password.value},
    noAuthHeader: true,
  }).json().catch(async (r) => {
    if (r.response.status >= 500) {
      alert.open(t('global.errors.errorOccurred'), AlertType.danger)
    } else {
      const errorData = await r.response.json();
      console.log(errorData);
      alert.open(t(`auth.responses.${errorData.message}`), AlertType.warning)
    }
  }) as tokenResponse;

  alert.open(t('auth.alert.successLogin'), AlertType.success)
  localStorage.setItem("token", JSON.stringify(user.token));
  router.push("/");
}

</script>

<template>
  <h2 class="card-title">{{ $t('auth.login') }}</h2>
  <div class="items-center mt-2">
    <label class="input input-bordered flex items-center gap-2 mb-2">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
           class="w-4 h-4 opacity-70">
        <path
          d="M2.5 3A1.5 1.5 0 0 0 1 4.5v.793c.026.009.051.02.076.032L7.674 8.51c.206.1.446.1.652 0l6.598-3.185A.755.755 0 0 1 15 5.293V4.5A1.5 1.5 0 0 0 13.5 3h-11Z"/>
        <path
          d="M15 6.954 8.978 9.86a2.25 2.25 0 0 1-1.956 0L1 6.954V11.5A1.5 1.5 0 0 0 2.5 13h11a1.5 1.5 0 0 0 1.5-1.5V6.954Z"/>
      </svg>
      <input v-model="email" type="text" class="grow" :placeholder="$t('auth.fields.email')"/>
    </label>
    <label class="input input-bordered flex items-center gap-2 mb-2">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
           class="w-4 h-4 opacity-70">
        <path fill-rule="evenodd"
              d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
              clip-rule="evenodd"/>
      </svg>
      <input v-model="password" type="password" class="grow"
             :placeholder="$t('auth.fields.password')" value=""/>
    </label>
  </div>
  <div class="card-actions justify-center">
    <!-- TODO: in the future   <a href="#" class="text-center">Forgot password?</a>-->
    <button @click="login" class="btn btn-primary w-full">{{ $t('auth.login') }}</button>
    <RouterLink :to="{name: 'auth.register'}" class="text-center">{{
        $t('auth.doesNotHaveAccount')
      }}
    </RouterLink>
  </div>
</template>

<style scoped>

</style>
