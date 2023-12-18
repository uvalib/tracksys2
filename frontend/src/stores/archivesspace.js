import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useArchivesSpaceStore = defineStore('archivesspace', {
   state: () => ({
      working: false,
      reviews: [],
      total: 0,
      viewerBaseURL: "",
      searchOpts: {
         sortField: "submittedAt",
         sortOrder: "asc",
         query: "",
      },
   }),
   getters: {
   },
   actions: {
      getReviews() {
         const system = useSystemStore()
         this.working = true
         let so = this.searchOpts
         let url = `/api/archivesspace?by=${so.sortField}&order=${so.sortOrder}`
         if ( so.query != "") {
            url += `&q=${encodeURIComponent(so.query)}`
         }
         axios.get( url ).then(response => {
            this.reviews = response.data.submissions
            this.total = response.data.total
            this.viewerBaseURL = response.data.viewerBaseURL
         }).catch( e => {
            system.setError(e)
            this.reviews = []
            this.total = 0
         }).finally( () => {
           this.working = false
         })
      }
   }
})