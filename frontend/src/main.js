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
import DPGButton from "@/components/DPGButton.vue"
import ConfirmModal from "@/components/ConfirmModal.vue"
app.component("WaitSpinner", WaitSpinner)
app.component("ErrorMessage", ErrorMessage)
app.component("DPGButton", DPGButton)
app.component("ConfirmModal", ConfirmModal)

import Vue3EasyDataTable from 'vue3-easy-data-table'
import 'vue3-easy-data-table/dist/style.css'
app.component('EasyDataTable', Vue3EasyDataTable)

// Styles
import '@fortawesome/fontawesome-free/css/all.css'
import './assets/styles/forms.scss'
import './assets/styles/main.scss'
import './assets/styles/uva-colors.css'

app.mount('#app')
