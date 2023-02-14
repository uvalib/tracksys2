import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useCollectionsStore = defineStore('collections', {
   state: () => ({
      working: false,
      records: [],
      totalRecords: 0,
      colectionID: -1,
      searchOpts: {
         start: 0,
         limit: 30,
         query: "",
      },
      collections: [],
      totalCollections: 0,
   }),
   getters: {
	},
	actions: {
      setCollection( collectionID ) {
         this.collectionID = collectionID
         this.records = []
         this.totalRecords = 0
         this.searchOpts.start = 0
         this.searchOpts.limit = 30
         this.searchOpts.query = ""
      },
      addCollection( data ) {
         let rec = {
            barcode: data.barcode,
            callNumber: data.callNumber,
            catalogKey: data.catalogKey,
            creatorName: data.creatorName,
            id: data.id,
            pid: data.pid,
            recordCount: 0,
            title: data.title,
            type: data.type,
         }
         this.collections.push(rec)
      },
      getCollections() {
         const system = useSystemStore()
         system.working = true
         this.collections = []
         this.totalCollections = 0
         axios.get( "/api/collections" ).then(response => {
            this.totalCollections = response.data.total
            this.collections = response.data.collections
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
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
            this.totalRecords = response.data.total
            this.working = false
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
            this.working = false
         })
      },
   }
})