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

      async claimForReview( item, reviewerID ) {
         const system = useSystemStore()
         this.working = true
         await axios.post(`/api/metadata/${item.metadata.id}/archivesspace/review?reviewer=${reviewerID}`).then( (resp) => {
            system.toastMessage('Review Started', 'You have successfully claimed this item for review')
            let tgtIdx = this.reviews.findIndex( r => r.id == item.id)
            this.reviews[tgtIdx].reviewer = resp.data.reviewer
            this.reviews[tgtIdx].reviewStartedAt = resp.data.reviewStartedAt
            this.reviews[tgtIdx].status = resp.data.status
         }).catch( e => {
            system.toastError('Review Failed', e)
         }).finally( () => {
           this.working = false
         })
      },

      async resubmit( item ) {
         const system = useSystemStore()
         this.working = true
         await axios.post(`/api/metadata/${item.metadata.id}/archivesspace/resubmit`).then( (resp) => {
            system.toastMessage('Resubmitted', 'You have successfully resubmitted this item for review')
            let tgtIdx = this.reviews.findIndex( r => r.id == item.id)
            this.reviews[tgtIdx].status = "requested"
         }).catch( e => {
            system.toastError('Resubmit Failed', e)
         }).finally( () => {
           this.working = false
         })
      }
   }
})