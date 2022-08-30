import { defineStore } from 'pinia'
import axios from 'axios'
import { useSystemStore } from './system'

export const useSearchStore = defineStore('search', {
	state: () => ({
      query: "",
      scope: "all",
      field: "all",
      components: {
         start: 0,
         limit: 30,
         total: 0,
         sortField: "id",
         sortOrder: "desc",
         hits: [],
      },
      masterFiles: {
         start: 0,
         limit: 30,
         total: 0,
         sortField: "id",
         sortOrder: "desc",
         hits: [],
      },
      metadata: {
         start: 0,
         limit: 30,
         total: 0,
         sortField: "id",
         sortOrder: "desc",
         hits: [],
      },
      orders: {
         start: 0,
         limit: 30,
         total: 0,
         sortField: "id",
         sortOrder: "desc",
         hits: [],
      },
      searchFields: {}
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
         this.components.sortField = "id"
         this.components.sortOrder = "desc"
         this.components.hits = []

         this.masterFiles.start = 0
         this.masterFiles.limit = 30
         this.masterFiles.total = 0
         this.masterFiles.sortField = "id"
         this.masterFiles.sortOrder = "desc"
         this.masterFiles.hits = []

         this.metadata.start = 0
         this.metadata.limit = 30
         this.metadata.total = 0
         this.metadata.sortField = "id"
         this.metadata.sortOrder = "desc"
         this.metadata.hits = []

         this.orders.start = 0
         this.orderslimit = 30
         this.orderstotal = 0
         this.orderssortField = "id"
         this.orderssortOrder = "desc"
         this.ordershits = []
      },
      globalSearch() {
         const system = useSystemStore()
         system.working = true
         this.resetResults()
         let url = `/api/search?scope=${this.scope}&q=${encodeURIComponent(this.query)}`
         if (this.field != "all" ) {
            url += `&field=${this.field}`
         }
         axios.get(url).then(response => {
            this.components.hits = response.data.components.hits
            this.components.total = response.data.components.total
            this.masterFiles.hits = response.data.masterFiles.hits
            this.masterFiles.total = response.data.masterFiles.total
            this.metadata.hits = response.data.metadata.hits
            this.metadata.total = response.data.metadata.total
            this.orders.hits = response.data.orders.hits
            this.orders.total = response.data.orders.total
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})