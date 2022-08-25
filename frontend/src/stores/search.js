import { defineStore } from 'pinia'
import axios from 'axios'
import { useSystemStore } from './system'

export const useSearchStore = defineStore('search', {
	state: () => ({
      query: "",
      scope: "all",
      field: "all",
      masterFiles: [],
      metadata: [],
      orders: [],
      searchFields: {}
	}),
	getters: {
      hasResults: state => {
         return (state.masterFiles && state.masterFiles.length > 0) ||
            (state.metadata && state.metadata.length > 0) ||
            (state.orders && state.orders.length > 0)
      }
	},
	actions: {
      setGlobalSearchFields( data ) {
         this.searchFields = data
      },
      globalSearch() {
         const system = useSystemStore()
         system.working = true
         this.masterFiles = []
         this.metadata = []
         this.orders = []
         axios.get(`/api/search?scope=${this.scope}&q=${encodeURIComponent(this.query)}`).then(response => {
            this.masterFiles = response.data.masterFiles
            this.metadata = response.data.metadata
            this.orders = response.data.orders
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})