import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useAPTrustStore = defineStore('aptrust', {
   state: () => ({
      working: false,
      loadingReport: false,
      collectionStatus: {
         totalSubmitted: 0,
         successCount: 0,
         failures: [],
         errorMessage: ""
      },
      itemStatus: {
         id: 0,
         bag: "",
         etag: "",
         objectIdentifier: "",
         storage: "",
         note: "",
         status: "",
         requestedAt: "",
         submittedAt: "",
         finishedAt: "",
         errorMessage: ""
      },
      submissions: [],
      total: 0,
      searchOpts: {
         start: 0,
         limit: 30,
         sortField: "id",
         sortOrder: "desc",
         query: "",
      },
   }),
   getters: {
      hasItemStatus: state => {
         return state.itemStatus.requestedAt != ""
      },
   },
   actions: {
      getSubmissions(showWorking = true) {
         const system = useSystemStore()
         if ( showWorking == true ) this.working = true
         let so = this.searchOpts
         let url = `/api/aptrust?start=${so.start}&limit=${so.limit}&by=${so.sortField}&order=${so.sortOrder}`
         if ( so.query != "") {
            url += `&q=${encodeURIComponent(so.query)}`
         }
         axios.get( url ).then(response => {
            this.submissions = response.data.submissions
            this.total = response.data.total
         }).catch( e => {
            system.setError(e)
         }).finally( () => {
           this.working= false
         })
      },
      clearItemStatus() {
         this.itemStatus.id = 0
         this.itemStatus.bag = ""
         this.itemStatus.etag = ""
         this.itemStatus.objectIdentifier = ""
         this.itemStatus.storage = ""
         this.itemStatus.note = ""
         this.itemStatus.status = ""
         this.itemStatus.requestedAt = ""
         this.itemStatus.submittedAt = ""
         this.itemStatus.finishedAt = ""
         this.itemStatus.errorMessage = ""
      },
      getItemStatus( metadaID ) {
         this.working = true
         this.clearItemStatus()
         axios.get(`/api/metadata/${metadaID}/aptrust`).then((response) => {
            this.itemStatus = response.data
            if ( this.itemStatus.finishedAt ) {
               // update processedAt and success in the submissions list for this item
               let idx = this.submissions.findIndex( s => s.id == response.data.id)
               if (idx > -1 ) {
                  this.submissions[idx].processedAt = this.response.data.finishedAt
                  this.submissions[idx].status = (this.response.data.status == "Success")
               }
            }
         }).catch((error) => {
            this.itemStatus.errorMessage = error
         }).finally( () =>{
            this.working = false
         })
      },
      async getCollectionStatusReport( collectionID ) {
         this.loadingReport = true
         this.collectionStatus.totalSubmitted = 0
         this.collectionStatus.successCount = 0
         this.collectionStatus.failures = []
         this.collectionStatus.errorMessage = ""
         return axios.get(`/api/collections/${collectionID}/aptrust`).then((response) => {
            this.collectionStatus.totalSubmitted = response.data.length
            this.collectionStatus.errorMessage = ""
            this.collectionStatus.successCount = 0
            response.data.forEach( (s) => {
               if ( s.status == "Success" ) {
                  this.collectionStatus.successCount++
               } else {
                  this.collectionStatus.failures.push( {id: s.metadata_id, pid: s.metadata_pid,  title: s.metadata_title, error: s.note} )
               }
            })
         }).catch((error) => {
            this.collectionStatus.errorMessage = error
         }).finally( () =>{
            this.loadingReport = false
         })
      },
      resubmitCollectionItems( collectionID, metadataIDs ) {
         this.working = true
         let req = {collectionID: collectionID, metadataRecords: metadataIDs}
         const system = useSystemStore()
         axios.post(`${system.jobsURL}/aptrust`, req).then(() => {
            system.toastMessage('Submitted', 'The selected items have begun the APTrust resubmission process; check the job status page for updates')
            this.working = false
         }).catch((error) => {
            system.toastError('Submit Failed', `APTrust resubmission failed: ${error}`)
            this.working = false
         })
      },
      async submitItem( metadataID, resubmit ) {
         const system = useSystemStore()
         this.working = true
         let url = `${system.jobsURL}/metadata/${metadataID}/aptrust`
         if (resubmit) {
            url += "?resubmit=1"
         }
         return axios.post( url ).catch( e => {
            const system = useSystemStore()
            system.setError(e)
         }).finally (
            this.working = false
         )
      },
   }
})