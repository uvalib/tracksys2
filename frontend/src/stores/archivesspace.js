import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useArchivesSpaceStore = defineStore('archivesspace', {
   state: () => ({
      working: false,
      reviews: [],
      total: 0,
      viewerBaseURL: "",
      searchOpts: {
         sortField: "submittedAt",
         sortOrder: "asc",
         query: "",
      },
   }),
   getters: {
   },
   actions: {
      getReviews() {
         const system = useSystemStore()
         this.working = true
         let so = this.searchOpts
         let url = `/api/archivesspace?by=${so.sortField}&order=${so.sortOrder}`
         if ( so.query != "") {
            url += `&q=${encodeURIComponent(so.query)}`
         }
         axios.get( url ).then(response => {
            this.reviews = response.data.submissions
            this.total = response.data.total
            this.viewerBaseURL = response.data.viewerBaseURL
         }).catch( e => {
            system.setError(e)
            this.reviews = []
            this.total = 0
         }).finally( () => {
           this.working = false
         })
      },

      async publish( userID, metadata ) {
         const system = useSystemStore()
         this.working = true
         await axios.post(`/api/metadata/${metadata.id}/archivesspace/publish?user=${userID}`).then( () => {
            let idx = this.reviews.findIndex( r => r.metadataID == metadata.id )
            if ( idx > -1 )  {
               this.reviews.splice(idx,1)
            }
            system.toastMessage('Published', `You have successfully published metadata ${metadata.pid} to ArchivesSpace. Images should appear within a few minutes.`)
         }).catch( e => {
            system.setError(e)
         }).finally( () => {
            this.working = false
          })
      },

      async reject( userID, item, notes ) {
         const system = useSystemStore()
         this.working = true
         var data = { userID: userID, notes: notes}
         await axios.post(`/api/metadata/${item.id}/archivesspace/reject`, data).then( () => {
            let tgtIdx = this.reviews.findIndex( r => r.metadataID == item.id)
            if ( tgtIdx > -1 ) {
               this.reviews[tgtIdx].notes = notes
               this.reviews[tgtIdx].status = "rejected"
               system.toastMessage('Submission Rejected', `You have successfully rejected submission ${item.pid}`)
            } else {
               console.error("couldnt find review to reject")
               console.error(item)
            }
         }).catch( e => {
            system.toastError('Reject Failed', e)
         }).finally( () => {
           this.working = false
         })
      },

      async claimForReview( item, reviewerID ) {
         const system = useSystemStore()
         this.working = true
         await axios.post(`/api/metadata/${item.id}/archivesspace/review?user=${reviewerID}`).then( (resp) => {
            let tgtIdx = this.reviews.findIndex( r => r.metadataID == item.id)
            if ( tgtIdx > -1 ) {
               this.reviews[tgtIdx].reviewer = resp.data.reviewer
               this.reviews[tgtIdx].reviewStartedAt = resp.data.reviewStartedAt
               this.reviews[tgtIdx].status = resp.data.status
               system.toastMessage('Review Started', 'You have successfully claimed this item for review')
            } else {
               console.error("couldnt find review to claim")
               console.error(item)
            }
         }).catch( e => {
            system.toastError('Review Failed', e)
         }).finally( () => {
           this.working = false
         })
      },

      async resubmit( item ) {
         const system = useSystemStore()
         this.working = true
         await axios.post(`/api/metadata/${item.id}/archivesspace/resubmit`).then( () => {
            let tgtIdx = this.reviews.findIndex( r => r.metadataID == item.id)
            this.reviews[tgtIdx].status = "requested"
            system.toastMessage('Resubmitted', 'You have successfully resubmitted this item for review')
         }).catch( e => {
            system.toastError('Resubmit Failed', e)
         }).finally( () => {
           this.working = false
         })
      }
   }
})