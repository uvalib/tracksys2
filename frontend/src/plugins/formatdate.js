import dayjs from 'dayjs'

export default {
   install: (app) => {
      // inject a globally available $translate() method
      app.config.globalProperties.$formatDate = (dateStr, includeTime = true) => {
         if (dateStr) {
            let d = dayjs(dateStr)
            if ( includeTime ) {
               return d.format("YYYY-MM-DD HH:mm")
            }
            return d.format("YYYY-MM-DD")
         }
         return ""
      }
   }
}