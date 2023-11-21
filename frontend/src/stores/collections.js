import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useCollectionsStore = defineStore('collections', {
   state: () => ({
      working: false,
      records: [],
      totalRecords: 0,
      collectionID: -1,
      inAPTrust: false,
      searchOpts: {
         start: 0,
         limit: 30,
         query: "",
         sortField: "id",
         sortOrder: "desc",
      },
      collections: [],
      totalCollections: 0,
      bulkAdd: false,
      metadataHits: [],
      apTrustStatus: {
         totalSubmitted: 0,
         successCount: 0,
         failures: [],
         errorMessage: ""
      },
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
      addItems( items ) {
         this.working = true
         axios.post( `/api/collections/${this.collectionID}/items`, {items: items} ).then(() => {
            this.working = false
            this.getItems()
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
            this.working = false
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
      getItems( ) {
         this.working = true
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
         let url = `/api/search?scope=metadata&q=${encodeURIComponent(query)}&collection=1`
         this.metadataHits = []
         axios.get(url).then(response => {
            response.data.metadata.hits.forEach( mh => {
               this.metadataHits.push(
                  { id: mh.id, pid: mh.pid, title: mh.title, catalogKey: mh.catalogKey,
                    callNumber: mh.callNumber, barcode: mh.barcode, masterFilesCount: mh.masterFilesCount }
               )
            })
         }).catch( e => {
            system.setError(e)
         })
      },
      apTrustResubmit( metadataIDs ) {
         if (this.collectionID == -1 || !this.inAPTrust) return

         let req = {metadataRecords: metadataIDs}
         const system = useSystemStore()
         axios.post(`${system.jobsURL}/aptrust`, req).then((response) => {
            system.toastMessage('Submitted', 'The selected items have begun the APTrust submission process; check the job status page for updates')
         }).catch((error) => {
            system.toastError('Submit Failed', `APTrust submission failed: ${error}`)
         })
      },
      async getAPTrustStatus() {
         if (this.collectionID == -1 || !this.inAPTrust) return

         this.working = true
         return axios.get(`/api/collections/${this.collectionID}/aptrust`).then((response) => {
            this.apTrustStatus.totalSubmitted = response.data.length
            this.apTrustStatus.errorMessage = ""
            this.apTrustStatus.successCount = 0
            response.data.forEach( (s) => {
               if ( s.status == "Success" ) {
                  this.apTrustStatus.successCount++
               } else {
                  this.apTrustStatus.failures.push( {id: s.metadata_id, pid: s.metadata_pid,  error: s.note} )
               }
            })
            for (let i=1400; i<1410; i++) {
               let title = "Fake Title "+i
               if ( i == 1400 ) {
                  title = "Declaration of Independence of the State of South Carolina : in convention, at the city of Charleston, December 20, 1860. : An ordinance to dissolve the Union between the state of South Carolina and other states united with her under the compact entitled \"The constitution of the United States of America.\""
               }
               this.apTrustStatus.failures.push( {id: i, pid: "tsb:"+i, title: title,  error: "This is fake error #"+i} )
            }
         }).catch((error) => {
            this.apTrustStatus.errorMessage = error
         }).finally( () =>{
            this.working = false
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