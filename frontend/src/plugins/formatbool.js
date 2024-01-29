export default {
   install: (app) => {
      app.config.globalProperties.$formatBool = (boolFlag) => {
         if (boolFlag) return "Yes"
         return "No"
      }
   }
}