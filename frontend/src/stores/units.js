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

      async createProject( data ) {
         const system = useSystemStore()
         return axios.post( `/api/units/${this.detail.id}/project`, data ).then( (resp) => {
            this.detail.projectID = parseInt(resp.data, 10)
         }).catch( e => {
            system.setError(e)
         })
      },

      async attachFile( info ) {
         const system = useSystemStore()
         var formData = new FormData()
         let fileData = info.attachment[0]
         formData.append("description", info.description)
         formData.append("name", fileData.name)
         formData.append("file", fileData.file)
         let url = `${system.jobsURL}/units/${this.detail.id}/attach`
         return axios.post(url, formData, {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
         }).then( async () => {
            let uID = this.detail.id
            this.detail.id = 0
            await this.getDetails(uID)
         }).catch( e => {
            system.setError(e)
         })
      },

      async submitEdit( edit ) {
         const system = useSystemStore()
         system.working = true
         return axios.post( `/api/units/${this.detail.id}/update`, edit ).then( (resp) => {
            this.detail = resp.data
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