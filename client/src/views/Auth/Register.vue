<script setup lang="ts">
import {AlertType, useAlert} from "@/stores/alert.ts";
import {defaultInstance, type HttpRequestOptions} from "@/services/api.ts";
import {ref} from "vue";
import {useI18n} from "vue-i18n";
import {useRouter} from "vue-router";

const {t} = useI18n()
const router = useRouter();
const alert = useAlert();

const email = ref<string>();
const password = ref<string>();
const confirmPassword = ref<string>();
const loading = ref(false);

const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

const register = async (e: Event) => {
  if (email.value === '' ||
    email.value === undefined ||
    password.value === '' ||
    password.value === undefined ||
    confirmPassword.value === '' ||
    confirmPassword.value === undefined
  ) {
    alert.open(t('auth.errors.invalidInput'), AlertType.warning)
    return
  }

  if (email.value && !isValidEmail(email.value)) {
    alert.open(t('auth.errors.invalidEmail'), AlertType.warning)
    return
  }

  if (password.value !== confirmPassword.value) {
    alert.open(t("auth.errors.passwordMismatch"), AlertType.warning);
    return
  }

  const opts: HttpRequestOptions = {
    json: {Email: email.value, Password: password.value},
    noAuthHeader: true,
  }

  loading.value = true;

  await defaultInstance.post('auth/register', opts).json().catch(async (r) => {
    loading.value = false;
    if (r.response.status >= 500) {
      alert.open(t('global.errors.errorOccurred'), AlertType.danger)
    } else {
      const errorData = await r.response.json();
      console.log(errorData);
      alert.open(t(`auth.responses.${errorData.message}`), AlertType.warning)
    }
  });

  alert.open(t('auth.alert.successRegister'), AlertType.success)
  loading.value = false;
  await router.push({name: 'auth.login'})
}
</script>

<template>
  <form>
  <h2 class="card-title">{{ $t("auth.register") }}</h2>
  <div class="items-center mt-2">
    <label class="w-full input input-bordered flex items-center gap-2 mb-2">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
           class="w-4 h-4 opacity-70">
        <path
          d="M2.5 3A1.5 1.5 0 0 0 1 4.5v.793c.026.009.051.02.076.032L7.674 8.51c.206.1.446.1.652 0l6.598-3.185A.755.755 0 0 1 15 5.293V4.5A1.5 1.5 0 0 0 13.5 3h-11Z"/>
        <path
          d="M15 6.954 8.978 9.86a2.25 2.25 0 0 1-1.956 0L1 6.954V11.5A1.5 1.5 0 0 0 2.5 13h11a1.5 1.5 0 0 0 1.5-1.5V6.954Z"/>
      </svg>
      <input type="text" class="grow" v-model="email" :placeholder="$t('auth.fields.email')"/>
    </label>
    <label class="w-full input input-bordered flex items-center gap-2 mb-2">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
           class="w-4 h-4 opacity-70">
        <path fill-rule="evenodd"
              d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
              clip-rule="evenodd"/>
      </svg>
      <input type="password" class="grow" v-model="password"
             :placeholder="$t('auth.fields.password')" value=""/>
    </label>
    <label class="w-full input input-bordered flex items-center gap-2 mb-2">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor"
           class="w-4 h-4 opacity-70">
        <path fill-rule="evenodd"
              d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
              clip-rule="evenodd"/>
      </svg>
      <input type="password" class="grow" v-model="confirmPassword"
             :placeholder="$t('auth.fields.confirmPassword')"
             value=""/>
    </label>
  </div>
  <div class="card-actions justify-center">
    <button type="submit" :disabled="!email || !password || !confirmPassword"
            @click.prevent="register"
            class="btn btn-primary w-full">{{ $t("auth.register") }}
    </button>
    <RouterLink :to="{name: 'auth.login'}" class="text-center">{{ $t("auth.alreadyHaveAccount") }}
    </RouterLink>
  </div>
  </form>
</template>

<style scoped>

</style>
