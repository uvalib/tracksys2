import { createApp, markRaw } from 'vue'
import { createPinia } from 'pinia'
import formatDatePlugin from './plugins/formatdate'
import formatDateTimePlugin from './plugins/formatdatetime'
import formatBoolPlugin from './plugins/formatbool'
import './assets/main.css'

import App from './App.vue'
import router from './router'

const app = createApp(App)

const pinia = createPinia()
pinia.use(({ store }) => {
   // all stores can access router with this.router
   store.router = markRaw(router)
})

app.use(formatBoolPlugin)
app.use(formatDatePlugin)
app.use(formatDateTimePlugin)
app.use(pinia)
app.use(router)

// NuxtUI defaulted to LIGHT mode
import ui from '@nuxt/ui/vue-plugin'
import { useColorMode } from '@vueuse/core'
useColorMode().value = 'light'
app.use(ui)

app.mount('#app')