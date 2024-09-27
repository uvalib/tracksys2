import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'
import { useDateFormat } from '@vueuse/core'

export const useHathiTrustStore = defineStore('hathitrust', {
   state: () => ({
      working: false,
      submissions: [],
      total: 0,
      searchOpts: {
         start: 0,
         limit: 30,
         filters: [],
         sortField: "pid",
         sortOrder: "desc",
         query: "",
      },
   }),
   getters: {
      filtersAsQueryParam: state => {
         let out = []
         state.searchOpts.filters.forEach( fv => out.push(`${fv.field}|${fv.match}|${fv.value}`) )
         return JSON.stringify(out)
      }
   },
   actions: {
      getSubmissions(showWorking = true) {
         const system = useSystemStore()
         if ( showWorking == true ) this.working = true
         let so = this.searchOpts
         let url = `/api/hathitrust?start=${so.start}&limit=${so.limit}&by=${so.sortField}&order=${so.sortOrder}`
         if ( so.query != "") {
            url += `&q=${encodeURIComponent(so.query)}`
         }
         if ( so.filters.length > 0) {
            url += `&filters=${this.filtersAsQueryParam}`
         }
         axios.get( url ).then(response => {
            this.submissions = response.data.submissions
            this.total = response.data.total
         }).catch( e => {
            system.setError(e)
         }).finally( () => {
           this.working= false
         })
      },
      batchUpdateOrder(orderID, field, value) {
         const system = useSystemStore()
         let req = {orderID: orderID, field: field, value: value}
         if ( field == "metadata_submitted_at" || field == "package_submitted_at" || field == "finished_at") {
            req.value = useDateFormat(value, "YYYY-MM-DD").value
         }
         axios.put( `/api/hathitrust`, req ).then( () => {
            system.toastMessage("Updated", `HathiTrust status records have been updated.`)
         }).catch( e => {
            system.setError(e)
         })
      },
      batchUpdate(IDs, field, value) {
         console.log("BATCH")
         const system = useSystemStore()
         let req = {statusIDs: IDs, field: field, value: value}
         console.log(req)
         if ( field == "metadata_submitted_at" || field == "package_submitted_at" || field == "finished_at") {
            req.value = useDateFormat(value, "YYYY-MM-DD").value
         }
         axios.put( `/api/hathitrust`, req ).then( () => {
            system.toastMessage("Updated", `HathiTrust status records have been updated.`)
         }).catch( e => {
            system.setError(e)
         })
      },
   }
})