import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useHathiTrustStore = defineStore('hathitrust', {
   state: () => ({
      working: false,
      submissions: [],
      total: 0,
      searchOpts: {
         start: 0,
         limit: 30,
         sortField: "pid",
         sortOrder: "desc",
         query: "",
      },
   }),
   getters: {
   },
   actions: {
      getSubmissions() {
         const system = useSystemStore()
         this.working = true
         let so = this.searchOpts
         let url = `/api/hathitrust?start=${so.start}&limit=${so.limit}&by=${so.sortField}&order=${so.sortOrder}`
         if ( so.query != "") {
            url += `&q=${encodeURIComponent(so.query)}`
         }
         axios.get( url ).then(response => {
            this.submissions = response.data.submissions
            this.total = response.data.total
         }).catch( e => {
            system.setError(e)
         }).finally( () => {
           this.working= false
         })
      }
   }
})