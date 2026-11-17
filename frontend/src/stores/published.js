import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const usePublishedStore = defineStore('published', {
   state: () => ({
      records: [],
      total: 0,
      searchOpts: {
         currPage: 1,
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
      },
      exportCSV( pubType ) {
         axios.get(`/api/published/${pubType}/csv`, null, {responseType: "blob"}).then((response) => {
            const fileURL = window.URL.createObjectURL(new Blob([response.data], { type: 'application/vnd.ms-excel' }))
            const fileLink = document.createElement('a')
            fileLink.href =  fileURL
            fileLink.setAttribute('download', `published-${pubType}.csv`)
            document.body.appendChild(fileLink)
            fileLink.click()
            window.URL.revokeObjectURL(fileURL)
         }).catch((error) => {
            const system = useSystemStore()
            system.setError(error)
         })
      },
      resetSearch() {
         this.searchOpts.currPage = 1
         this.searchOpts.start = 0
         this.orders = []
         this.total = 0
      },
   }
})