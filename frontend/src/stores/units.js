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
      clearDetails() {
         this.detail.id = 0
         this.detail.status = ""
         this.detail.intendedUse = null
         this.detail.includeInDL = false
         this.detail.removeWaterMark = false
         this.detail.reorder = false
         this.detail.completeScan = false
         this.detail.throwAway = false
         this.detail.ocrMasterFiles = false
         this.detail.attachments = []
         this.detail.staffNotes = ""
         this.detail.dateArchived = ""
         this.detail.datePatronDeliverablesReady = ""
         this.detail.dateDLDeliverablesReady = ""
      },
      async getDetails( unitID ) {
         if ( this.detail.id == unitID ) return

         const system = useSystemStore()
         system.working = true
         return axios.get( `/api/units/${unitID}` ).then(response => {
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