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
         limit: 15,
         total: 0,
         hits: [],
         filters: []
      },
      masterFiles: {
         start: 0,
         limit: 15,
         total: 0,
         hits: [],
         filters: []
      },
      metadata: {
         start: 0,
         limit: 15,
         total: 0,
         hits: [],
         filters: []
      },
      orders: {
         start: 0,
         limit: 15,
         total: 0,
         hits: [],
         filters: []
      },
      searchFields: {},
	}),
	getters: {
      hasResults: state => {
         return state.masterFiles.total > 0 || state.metadata.total > 0 ||
            state.orders.total > 0 || state.components.total > 0
      },
      filtersAsQueryParam: state => {
         return (searchOrigin) => {
            let tgtFilters = null
            if (searchOrigin == "components") {
               tgtFilters = state.components.filters
            } else if (searchOrigin == "masterfiles") {
               tgtFilters = state.masterFiles.filters
            } else if (searchOrigin == "metadata") {
               tgtFilters = state.metadata.filters
            } else if (searchOrigin == "orders") {
               tgtFilters = state.orders.filters
            } else {
               return ""
            }
            if (tgtFilters != null && tgtFilters.length > 0) {
               let params = []
               tgtFilters.forEach( fv => params.push(`{"filter":"${fv.field}|${fv.match}|${encodeURIComponent(fv.value)}"}`) )
               return `[${params.join(",")}]`
            }
            return ""
         }
      }
	},
	actions: {
      setGlobalSearchFields( data ) {
         this.searchFields = data
      },
      resetResults() {
         this.components.start = 0
         this.components.limit = 15
         this.components.total = 0
         this.components.hits = []
         this.components.filters = []

         this.masterFiles.start = 0
         this.masterFiles.limit = 15
         this.masterFiles.total = 0
         this.masterFiles.hits = []
         this.masterFiles.filters = []

         this.metadata.start = 0
         this.metadata.limit = 15
         this.metadata.total = 0
         this.metadata.hits = []
         this.metadata.filters = []

         this.orders.start = 0
         this.orders.limit = 15
         this.orders.total = 0
         this.orders.hits = []
         this.orders.filters = []
      },
      resetSearch() {
         this.query = ""
         this.scope = "all"
         this.field = "all"
         this.searched = false
         this.resetResults()
      },
      setFilter( scope, filterQueryParm) {
         let parsedFilters = []
         let filterObj = JSON.parse(filterQueryParm)
         filterObj.forEach( f => {
            let bits = f.filter.split("|") // ex: title|contains|charlottesville
            parsedFilters.push({field: bits[0].trim(), match: bits[1].trim(), value: bits[2].trim()})
         })
         if (scope == "components") {
            this.components.filters = parsedFilters
         } else if (scope == "masterfiles") {
            this.masterFiles.filters = parsedFilters
         } else if (scope == "metadata") {
            this.metadata.filters = parsedFilters
         } else if (scope == "orders") {
            this.orders.filters = parsedFilters
         }
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
         } else if (searchOrigin == "orders") {
            url += `&start=${this.orders.start}&limit=${this.orders.limit}`
         } else {
            this.resetResults()
         }

         let filterParam = this.filtersAsQueryParam(searchOrigin)
         console.log("EXEC SEARCH in "+searchOrigin+ " FILTER "+filterParam)
         if ( filterParam != "") {
            url += `&filters=${filterParam}`
         }

         console.log("SEARCH URL "+url)
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