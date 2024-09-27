import { useDateFormat } from '@vueuse/core'

export default {
   install: (app) => {
      app.config.globalProperties.$formatDateTime = (date) => {
         if (date) {
            return useDateFormat(date, "YYYY-MM-DD HH:mm A").value
         }
         return ""
      }
   }
}