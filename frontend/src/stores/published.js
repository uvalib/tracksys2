import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const usePublishedStore = defineStore('published', {
   state: () => ({
      records: [],
      total: 0,
      searchOpts: {
         start: 0,
         limit: 25,
         filters: [],
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
      getRecords( pubType ) {
         const system = useSystemStore()
         this.total = 0
         this.records = []
         system.working = true
         let so = this.searchOpts
         let url = `/api/published/${pubType}?filters=${this.filtersAsQueryParam}&start=${so.start}&limit=${so.limit}`
         axios.get( url ).then(response => {
            this.total = response.data.total
            this.records = response.data.records
            system.working = false
         }).catch( e => {
            system.setError(e)
            system.working = false
         })
      }
   }
})