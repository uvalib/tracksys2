import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// Global component registration. All components can use these without import
import WaitSpinner from "@/components/WaitSpinner.vue"
import ErrorMessage from "@/components/ErrorMessage.vue"
app.component("WaitSpinner", WaitSpinner)
app.component("ErrorMessage", ErrorMessage)

// Styles
import '@fortawesome/fontawesome-free/css/all.css'
import './assets/styles/forms.scss'
import './assets/styles/main.scss'
import './assets/styles/uva-colors.css'
import './assets/styles/styleoverrides.scss'


// Primevue setup
import PrimeVue from 'primevue/config'
import ConfirmationService from 'primevue/confirmationservice'

app.use(PrimeVue)
app.use(ConfirmationService)

import 'primevue/resources/themes/saga-blue/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'

import Button from 'primevue/button'
import ConfirmDialog from 'primevue/confirmdialog'
app.component("DPGButton", Button)
app.component("ConfirmDialog", ConfirmDialog)


app.mount('#app')
