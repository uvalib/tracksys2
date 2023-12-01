import { createApp, markRaw } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)

const pinia = createPinia()
pinia.use(({ store }) => {
   // all stores can access router with this.router
   store.router = markRaw(router)
})

app.use(pinia)
app.use(router)

// Styles
import '@fortawesome/fontawesome-free/css/all.css'
import './assets/styles/forms.scss'
import './assets/styles/main.scss'
import './assets/styles/uva-colors.css'
import './assets/styles/styleoverrides.scss'

// Primevue setup
import PrimeVue from 'primevue/config'
import ConfirmationService from 'primevue/confirmationservice'
import ToastService from 'primevue/toastservice'
import Tooltip from 'primevue/tooltip';

app.directive('tooltip', Tooltip)

app.use(PrimeVue, { ripple: true })
app.use(ConfirmationService)
app.use(ToastService)

import 'primevue/resources/themes/saga-blue/theme.css'
import 'primeicons/primeicons.css'

import Button from 'primevue/button'
import ConfirmDialog from 'primevue/confirmdialog'
app.component("DPGButton", Button)
app.component("ConfirmDialog", ConfirmDialog)

// FormKit
import { plugin, defaultConfig } from '@formkit/vue'
const dc = defaultConfig({
   plugins: [addErrorAlertIconPlugin, addRequiredNotePlugin],
   config: {
      classes: {
         input: '$reset dpg-form-input',
         label: '$reset dpg-form-label',
         messages: '$reset dpg-form-invalid',
         help: '$reset dpg-form-help',
      },
      incompleteMessage: false,
      validationVisibility: 'submit'
   }
})
app.use(plugin, dc)

app.mount('#app')

/* FORMKIT PLUGINS */
function addRequiredNotePlugin(node) {
   var showRequired = true
   node.on('created', () => {
      if (node.config.disableRequiredDecoration == true) {
         showRequired = false
      }
      const schemaFn = node.props.definition.schema
      node.props.definition.schema = (sectionsSchema = {}) => {
         const isRequired = node.props.parsedRules.some(rule => rule.name === 'required')

         if (isRequired && showRequired) {
            // this input has the required rule so we modify
            // the schema to add an astrics to the label.
            sectionsSchema.label = {
               attrs: {
                  innerHTML: `<i class="req fas fa-asterisk"></i><span class="req-label">${node.props.label}</span><span class="req">(required)</span>`
               },
               children: null//['$label', '*']
            }
         }
         return schemaFn(sectionsSchema)
      }
   })
}

function addErrorAlertIconPlugin(node) {
   node.on('created', () => {
      const schemaFn = node.props.definition.schema
      node.context.warningIcon = '<i class="fas fa-exclamation-triangle"></i>'
      node.props.definition.schema = (extensions) => {
         if (!extensions.message) {
            extensions.message = {
               attrs: {
                  innerHTML: '$warningIcon + " " + $message.value'
               },
               children: null
            }
         }
         return schemaFn(extensions)
      }
   })
}
