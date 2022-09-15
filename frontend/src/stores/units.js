import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useUnitsStore = defineStore('units', {
	state: () => ({
      detail: {
         id: 0,
         status: "",
         intendedUse: null,
         includeInDL: false,
         removeWaterMark: false,
         reorder: false,
         completeScan: false,
         throwAway: false,
         ocrMasterFiles: false,
         attachments: [],
         staffNotes: "",
         dateArchived: "",
         datePatronDeliverablesReady: "",
         dateDLDeliverablesReady: ""
      },
   }),
	getters: {
	},
	actions: {
      getDetails( unitID ) {
         const system = useSystemStore()
         system.working = true
         axios.get( `/api/units/${unitID}` ).then(response => {
            this.detail = response.data
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },

      getRelatedItems( unitID ) {
         const system = useSystemStore()
         this.related.units = []
         this.related.orders = []
         axios.get( `/api/units/${unitID}/related` ).then(response => {
               console.log(response.data)
         }).catch( e => {
            system.setError(e)
         })

      }

   }
})