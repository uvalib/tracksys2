import { defineStore } from 'pinia'
import axios from 'axios'
import { useSystemStore } from './system'

export const useSearchStore = defineStore('search', {
	state: () => ({
      query: "",
      scope: "all",
      searched: false,
      unitValid: false,
      view: "",               // name of the view for below. used in client query params
      components: {
         start: 0,
         limit: 15,
         scroll: "",
         total: 0,
         hits: [],
         filters: []
      },
      masterFiles: {
         start: 0,
         limit: 15,
         scroll: "",
         total: 0,
         hits: [],
         filters: []
      },
      metadata: {
         start: 0,
         limit: 15,
         scroll: "",
         total: 0,
         hits: [],
         filters: []
      },
      orders: {
         start: 0,
         limit: 15,
         scroll: "",
         total: 0,
         hits: [],
         filters: []
      },
      units: {
         start: 0,
         limit: 15,
         scroll: "",
         total: 0,
         hits: [],
         filters: []
      },
      searchPHash: 0,
      similarSearch: false,
      distance: 5,
      similarImages: {
         total: 0,
         hits: [],
      },
	}),
	getters: {
      hasResults: state => {
         return state.masterFiles.total > 0 || state.metadata.total > 0 ||
            state.orders.total > 0 || state.components.total > 0 || state.units.total > 0
      },
      filtersAsQueryParam: state => {
         return (filterTarget) => {
            let tgtFilters = null
            if (filterTarget == "components") {
               tgtFilters = state.components.filters
            } else if (filterTarget == "masterfiles") {
               tgtFilters = state.masterFiles.filters
            } else if (filterTarget == "metadata") {
               tgtFilters = state.metadata.filters
            } else if (filterTarget == "orders") {
               tgtFilters = state.orders.filters
            } else if (filterTarget == "units") {
               tgtFilters = state.units.filters
            }else {
               return ""
            }
            if (tgtFilters != null && tgtFilters.length > 0) {
               let out = {type: filterTarget, params: []}
               tgtFilters.forEach( fv => out.params.push(`${fv.field}|${fv.match}|${fv.value}`) )
               return JSON.stringify(out)
            }
            return ""
         }
      }
	},
	actions: {
      resetSearch() {
         this.$reset()
      },

      setFilter( filterQueryParm) {
         let parsedFilters = []
         let filterObj = JSON.parse(filterQueryParm)
         filterObj.params.forEach( f => {
            let bits = f.split("|") // ex: title|contains|charlottesville
            parsedFilters.push({field: bits[0].trim(), match: bits[1].trim(), value: bits[2].trim()})
         })
         if (filterObj.type == "components") {
            this.components.filters = parsedFilters
         } else if (filterObj.type == "masterfiles") {
            this.masterFiles.filters = parsedFilters
         } else if (filterObj.type == "metadata") {
            this.metadata.filters = parsedFilters
         } else if (filterObj.type == "orders") {
            this.orders.filters = parsedFilters
         } else if (filterObj.type == "units") {
            this.units.filters = parsedFilters
         }
      },

      imageSearch( pHash ) {
         const system = useSystemStore()
         system.working = true
         this.searchPHash = pHash
         this.similarSearch = true
         this.similarImages = {
            total: 0,
            hits: [],
         }
         axios.get(`/api/search/images?phash=${pHash}&distance=${this.distance}`).then(response => {
            this.similarImages.hits = response.data.hits
            this,this.similarImages.total = response.data.total
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },

      executeSearch( scopeOverride ) {
         const system = useSystemStore()
         system.working = true
         // console.log("exec search. scopeOverride: ["+scopeOverride+", scope: "+this.scope+", view: "+this.view)

         // this lets secondary queries on specific item types with different filter and paginiation settings
         // Ex; initial scope is all, but user is viewing masterfiles and goes to next page. Override scope to masterfiles
         // and apply the pagination changes
         let tgtScope = scopeOverride
         if ( !tgtScope ) {
            tgtScope = this.scope
         }

         let url = `/api/search?scope=${tgtScope}&q=${encodeURIComponent(this.query)}`

         if (tgtScope == "components") {
            url += `&start=${this.components.start}&limit=${this.components.limit}`
            if (this.components.scroll != "") {
               url += `&scroll=${this.components.scroll}`
            }
         } else if (tgtScope == "masterfiles") {
            url += `&start=${this.masterFiles.start}&limit=${this.masterFiles.limit}`
            if (this.masterFiles.scroll != "") {
               url += `&scroll=${this.masterFiles.scroll}`
            }
         } else if (tgtScope == "metadata") {
            url += `&start=${this.metadata.start}&limit=${this.metadata.limit}`
            if (this.metadata.scroll != "") {
               url += `&scroll=${this.metadata.scroll}`
            }
         } else if (tgtScope == "orders") {
            url += `&start=${this.orders.start}&limit=${this.orders.limit}`
            if (this.orders.scroll != "") {
               url += `&scroll=${this.orders.scroll}`
            }
         } else if (tgtScope == "units") {
            url += `&start=${this.units.start}&limit=${this.units.limit}`
            if (this.units.scroll != "") {
               url += `&scroll=${this.units.scroll}`
            }
         }

         // filter is always based on active view
         let filterParam = this.filtersAsQueryParam(this.view)
         if ( filterParam != "") {
            url += `&filters=${filterParam}`
         }

         console.log("SEARCH URL "+url)
         axios.get(url).then(response => {
            if (tgtScope == "components" || tgtScope == "all") {
               this.components.hits = response.data.components.hits
               this.components.scroll = response.data.components.scroll
               this.components.total = response.data.components.total
            }
            if (tgtScope == "masterfiles" || tgtScope == "all") {
               this.masterFiles.hits = response.data.masterFiles.hits
               this.masterFiles.scroll = response.data.masterFiles.scroll
               this.masterFiles.total = response.data.masterFiles.total
            }
            if (tgtScope == "metadata" || tgtScope == "all") {
               this.metadata.hits = response.data.metadata.hits
               this.metadata.scroll = response.data.metadata.scroll
               this.metadata.total = response.data.metadata.total
            }
            if (tgtScope == "orders" || tgtScope == "all") {
               this.orders.hits = response.data.orders.hits
               this.orders.scroll = response.data.orders.scroll
               this.orders.total = response.data.orders.total
            }
            if (tgtScope == "units" || tgtScope == "all") {
               this.units.hits = response.data.units.hits
               this.units.scroll = response.data.units.scroll
               this.units.total = response.data.units.total
            }
            if ( this.scope == "all" ) {
               if ( this.orders.total > 0) {
                  this.view = "orders"
               } else if  ( this.metadata.total > 0) {
                  this.view = "metadata"
               } else if  ( this.masterFiles.total > 0) {
                  this.view = "masterfiles"
               } else if  ( this.components.total > 0) {
                  this.view = "components"
               } else if  ( this.units.total > 0) {
                  this.view = "units"
               }
            } else {
               this.view = tgtScope
            }
            system.working = false
            this.searched = true
         }).catch( e => {
            system.setError(e)
         })
      },

      setActiveView( viewName ) {
         this.view = viewName
      }
	},
})