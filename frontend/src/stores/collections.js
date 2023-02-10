import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useCollectionsStore = defineStore('collections', {
   state: () => ({
      working: false,
      records: [],
      total: 0,
      colectionID: -1,
      searchOpts: {
         start: 0,
         limit: 30,
         query: "",
      },
   }),
   getters: {
	},
	actions: {
      setCollection( collectionID ) {
         this.collectionID = collectionID
         this.records = []
         this.total = 0
         this.searchOpts.start = 0
         this.searchOpts.limit = 30
         this.searchOpts.query = ""
      },
      getRecords( ) {
         this.working = true
         let so = this.searchOpts
         let url = `/api/collections/${this.collectionID}?start=${so.start}&limit=${so.limit}`
         if ( so.query != "") {
            url += `&q=${encodeURIComponent(so.query)}`
         }
         axios.get( url ).then(response => {
            this.records = response.data.records
            this.total = response.data.total
            this.working = false
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
            this.working = false
         })
      },
   }
})