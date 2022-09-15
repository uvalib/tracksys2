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
      masterFiles: []
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

      getMasterFiles( unitID ) {
         const system = useSystemStore()
         this.masterFiles = []
         axios.get( `/api/units/${unitID}/masterfiles` ).then(response => {
            this.masterFiles = response.data
         }).catch( e => {
            system.setError(e)
         })

      }

   }
})