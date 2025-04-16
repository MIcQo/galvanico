import './assets/main.css'

import {createApp} from 'vue'
import {createPinia} from 'pinia'
import {createI18n} from 'vue-i18n'
import {VueQueryPlugin} from '@tanstack/vue-query'

import App from './App.vue'
import router from './router'
import locale from "@/locale";

const i18n = createI18n({
  locale: 'en',
  fallbackLocale: 'en',
  messages: locale,
})

const app = createApp(App)

app.use(i18n)
app.use(createPinia())
app.use(router)
app.use(VueQueryPlugin)

app.mount('#app')
