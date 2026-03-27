import { createApp, markRaw } from 'vue'
import { createPinia } from 'pinia'
import formatDatePlugin from './plugins/formatdate'
import formatDateTimePlugin from './plugins/formatdatetime'
import formatBoolPlugin from './plugins/formatbool'

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

// Styles
import './assets/styles/main.scss'

// Primevue setup
import PrimeVue from 'primevue/config'
import UVA from './assets/theme/uva'
import 'primeicons/primeicons.css'
import ConfirmationService from 'primevue/confirmationservice'
import ToastService from 'primevue/toastservice'
import Tooltip from 'primevue/tooltip'
import Ripple from 'primevue/ripple'

app.directive('ripple', Ripple)
app.directive('tooltip', Tooltip)

app.use(PrimeVue, {
   ripple: true,
   Tooltip: true,
   theme: {
      preset: UVA,
      options: {
         prefix: 'p',
         darkModeSelector: '.dpg-dark'
      }
   }
})

app.use(ConfirmationService)
app.use(ToastService)

import Button from 'primevue/button'
import ConfirmDialog from 'primevue/confirmdialog'
app.component("DPGButton", Button)
app.component("ConfirmDialog", ConfirmDialog)

app.mount('#app')