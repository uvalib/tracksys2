import { defineStore } from 'pinia'
import axios from 'axios'
import { useSystemStore } from './system'

export const useSearchStore = defineStore('search', {
	state: () => ({
      query: "",
      scope: "all",
      field: "all",
      searched: false,
      components: {
         start: 0,
         limit: 30,
         total: 0,
         hits: [],
      },
      masterFiles: {
         start: 0,
         limit: 30,
         total: 0,
         hits: [],
         filters: []
      },
      metadata: {
         start: 0,
         limit: 30,
         total: 0,
         hits: [],
         filters: []
      },
      orders: {
         start: 0,
         limit: 30,
         total: 0,
         hits: [],
      },
      searchFields: {},
	}),
	getters: {
      hasResults: state => {
         return state.masterFiles.total > 0 || state.metadata.total > 0 ||
            state.orders.total > 0 || state.components.total > 0
      }
	},
	actions: {
      setGlobalSearchFields( data ) {
         this.searchFields = data
      },
      resetResults() {
         this.components.start = 0
         this.components.limit = 30
         this.components.total = 0
         this.components.hits = []

         this.masterFiles.start = 0
         this.masterFiles.limit = 30
         this.masterFiles.total = 0
         this.masterFiles.hits = []

         this.metadata.start = 0
         this.metadata.limit = 30
         this.metadata.total = 0
         this.metadata.hits = []

         this.orders.start = 0
         this.orders.limit = 30
         this.orders.total = 0
         this.orders.hits = []
      },
      executeSearch(searchOrigin) {
         const system = useSystemStore()
         system.working = true
         let tgtScope = this.scope
         if (searchOrigin != "all") {
            tgtScope = searchOrigin
         }
         let url = `/api/search?scope=${tgtScope}&q=${encodeURIComponent(this.query)}`
         if (this.field != "all" ) {
            url += `&field=${this.field}`
         }

         if (searchOrigin == "components") {
            url += `&start=${this.components.start}&limit=${this.components.limit}`
         } else if (searchOrigin == "masterfiles") {
            url += `&start=${this.masterFiles.start}&limit=${this.masterFiles.limit}`
         } else if (searchOrigin == "metadata") {
            url += `&start=${this.metadata.start}&limit=${this.metadata.limit}`
            if (this.metadata.filters.length > 0) {
               let params = []
               this.metadata.filters.forEach( fv => params.push(`{"filter":"${fv.field}|${fv.match}|${encodeURIComponent(fv.value)}"}`) )
               url += `&filters=[${params.join(",")}]`
            }
         } else if (searchOrigin == "orders") {
            url += `&start=${this.orders.start}&limit=${this.orders.limit}`
         } else {
            this.resetResults()
         }

         axios.get(url).then(response => {
            if (searchOrigin == "components" || searchOrigin == "all") {
               this.components.hits = response.data.components.hits
               this.components.total = response.data.components.total
            }
            if (searchOrigin == "masterfiles" || searchOrigin == "all") {
               this.masterFiles.hits = response.data.masterFiles.hits
               this.masterFiles.total = response.data.masterFiles.total
            }
            if (searchOrigin == "metadata" || searchOrigin == "all") {
               this.metadata.hits = response.data.metadata.hits
               this.metadata.total = response.data.metadata.total
            }
            if (searchOrigin == "orders" || searchOrigin == "all") {
               this.orders.hits = response.data.orders.hits
               this.orders.total = response.data.orders.total
            }
            system.working = false
            this.searched = true
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})