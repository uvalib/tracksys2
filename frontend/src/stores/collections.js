import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useCollectionsStore = defineStore('collections', {
   state: () => ({
      working: false,
      searching: false,
      records: [],
      totalRecords: 0,
      collectionID: -1,
      inAPTrust: false,
      searchOpts: {
         start: 0,
         limit: 30,
         query: "",
         sortField: "id",
         sortOrder: "asc",
      },
      collections: [],
      totalCollections: 0,
      bulkAdd: false,
      metadataHits: [],
   }),
   getters: {
	},
	actions: {
      setCollection( collection) {
         this.collectionID = collection.id
         this.inAPTrust = collection.inAPTrust
         this.records = []
         this.totalRecords = 0
         this.searchOpts.start = 0
         this.searchOpts.limit = 30
         this.searchOpts.query = ""
         this.metadataHits = []
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
         this.metadataHits = []
         this.totalCollections = 0
         axios.get( "/api/collections" ).then(response => {
            this.totalCollections = response.data.total
            this.collections = response.data.collections
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      async addToCollection( collectionID, metadataID ) {
         const system = useSystemStore()
         return axios.post( `/api/collections/${collectionID}/item?rec=${metadataID}` ).then(() => {
            system.toastMessage( "Collection Updated", "This metadata record has been added to the selected collection" )
         }).catch( e => {
            system.setError(e)
         })
      },
      removeItem( item ) {
         let url = `/api/collections/${this.collectionID}/items/${item.id}`
         axios.delete( url ).then(() => {
            let idx = this.records.findIndex( r => r.id == item.id)
            if (idx > -1) {
               this.records.splice(idx,1)
            }
            this.totalRecords-=1
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
         })
      },
      getItems( showWorking=true ) {
         if ( showWorking ) this.working = true
         let so = this.searchOpts
         let url = `/api/collections/${this.collectionID}?start=${so.start}&limit=${so.limit}&by=${so.sortField}&order=${so.sortOrder}`
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
      exportCSV() {
         axios.get(`/api/collections/${this.collectionID}/export`, null, {responseType: "blob"}).then((response) => {
            const fileURL = window.URL.createObjectURL(new Blob([response.data], { type: 'application/vnd.ms-excel' }))
            const fileLink = document.createElement('a')
            fileLink.href =  fileURL
            fileLink.setAttribute('download', `collection-${this.collectionID}.csv`)
            document.body.appendChild(fileLink)
            fileLink.click()
            window.URL.revokeObjectURL(fileURL)
         }).catch((error) => {
            const system = useSystemStore()
            system.setError(error)
         })
      },
      toggleBulkAdd() {
         this.bulkAdd = !this.bulkAdd
      },
      metadataSearch( query ) {
         const system = useSystemStore()
         let url = `/api/collections/candidates?q=${encodeURIComponent(query)}`
         this.metadataHits = []
         this.searching = true
         axios.get(url).then(response => {
            response.data.hits.forEach( mh => {
               let hit = {
                  id: mh.id, pid: mh.pid, title: mh.title, catalogKey: mh.catalogKey, type: mh.type,
                  callNumber: mh.callNumber, barcode: mh.barcode }
               if ( mh.externalSystemID > 0) {
                  const es = system.externalSystems.find( s => s.id == mh.externalSystemID)
                  if (es) {
                     hit.type = es.name
                  }
               }
               this.metadataHits.push( hit )
            })
         }).catch( e => {
            system.setError(e)
         }).finally( () => {
            this.searching =  false
         })
      },
      addRecords( metadataIDs ) {
         const system = useSystemStore()
         let req = { items: metadataIDs}
         let url = `${system.jobsURL}/collections/${this.collectionID}/add`
         axios.post(url, req)
      }
   }
})