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
      async deleteAttachment(item) {
         const system = useSystemStore()
         let url = `${system.jobsURL}/units/${this.detail.id}/attachments/${item.filename}`
         return axios.delete(url).then( async () => {
            let aIdx = this.detail.attachments.findIndex( a => a.id == item.id)
            if (aIdx > -1) {
               this.detail.attachments.splice(aIdx, 1)
            }
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
      },

      regenerateIIIF() {
         const system = useSystemStore()
         let url = `${system.iiifManifestURL}/pid/${this.detail.metadata.pid}?refresh=true`
         axios.get( url ).then( () => {
            system.toastMessage("IIIF Manifest Regenerated", "The IIIF manifest for this unit has been regenerated.")
         }).catch( e => {
            system.setError(e)
         })
      },

      generateDeliverables() {
         const system = useSystemStore()
         let url = `${system.jobsURL}/units/${this.detail.id}/deliverables`
         axios.post(url).then( () => {
            system.toastMessage("Generate Deliverables", "Generating unit deliverables.")
         }).catch( e => {
            system.setError(e)
         })
      },

      downloadFromArchive( computeID, filename='all' ) {
         const system = useSystemStore()
         let payload = {computeID: computeID, filename: filename}
         let url = `${system.jobsURL}/units/${this.detail.id}/copy`
         axios.post(url, payload).then( () => {
            let tgt = "Master File"
            if (filename == 'all') tgt = "Unit"
            system.toastMessage("Archive Download", `${tgt} is being downloaded from the archive.`)
         }).catch( e => {
            system.setError(e)
         })
      },

      startUnitOCR() {
         const system = useSystemStore()
         let payload = {type: "unit", id: this.detail.id}
         axios.post(`${system.jobsURL}/ocr`, payload).then( () => {
            system.toastMessage("OCR Started", 'OCR has begun. Check the job status page for updates.')
         }).catch( e => {
            system.setError(e)
         })
      },

      setExemplar( mfID ) {
         const system = useSystemStore()
         axios.post( `/api/units/${this.detail.id}/exemplar/${mfID}` ).then( () => {
            this.masterFiles.forEach( mf => {
               mf.exemplar = false
               if (mf.id == mfID) {
                  mf.exemplar = true
               }
            })
            system.toastMessage("Exemplar", "New exemplar has been set.")
         }).catch( e => {
            system.setError(e)
         })
      }
   }
})