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
         dateDLDeliverablesReady: "",
         relatedUnits: [],
      },
      masterFiles: [],
      updateInProgress: false
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
      flagAsReorder() {
         this.detail.reorder = true
      },
      async getDetails( unitID ) {
         if ( this.detail.id == unitID ) return

         const system = useSystemStore()
         system.working = true
         return axios.get( `/api/units/${unitID}` ).then(response => {
            this.detail = response.data
            system.working = false
         }).catch( e => {
            if (e.response && e.response.status == 404) {
               this.router.push("/not_found")
               system.working = false
            } else {
               system.setError(e)
            }
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

      generateDeliverables() {
         const system = useSystemStore()
         let url = `${system.jobsURL}/units/${this.detail.id}/deliverables`
         axios.post(url).then( () => {
            system.toastMessage("Generate Deliverables", "Generating unit deliverables.")
         }).catch( e => {
            system.setError(e)
         })
      },

      downloadFromArchive( computeID, downloadTarget='all' ) {
         const system = useSystemStore()
         let msg = "Master file is"
         let payload = {computeID: computeID, filename: downloadTarget}
         if (Array.isArray(downloadTarget)) {
            payload = {computeID: computeID, files: downloadTarget}
            msg = "Master files are"
         } else {
            if (downloadTarget == 'all') {
               msg = "Unit is"
            }
         }
         let url = `${system.jobsURL}/units/${this.detail.id}/copy`
         axios.post(url, payload).then( () => {
            system.toastMessage("Archive Download", `${msg} being downloaded from the archive.`)
         }).catch( e => {
            system.setError(e)
         })
      },

      awaitUpdateCompletion( jobID, message ) {
         const system = useSystemStore()
         var tid = setInterval( ()=> {
            axios.get(`${system.jobsURL}/jobs/${jobID}`).then( resp => {
               let status = resp.data.status
               if (status == 'failure') {
                  clearInterval(tid)
                  this.updateInProgress = false
                  system.setError(`Update failed: ${resp.data.error}. Check the job status logs for more information.`)
               } else if (status == 'finished') {
                  this.getMasterFiles( this.detail.id )
                  clearInterval(tid)
                  system.toastMessage("Update Complete", message)
                  this.updateInProgress = false
               }
            }).catch( e => {
               system.setError(e)
               this.updateInProgress = false
               clearInterval(tid)
            })
         }, 1000)
      },

      addMasterFiles() {
         const system = useSystemStore()
         this.updateInProgress = true
         axios.post(`${system.jobsURL}/units/${this.detail.id}/masterfiles/add`).then( resp => {
            system.toastMessage("Please Wait", 'Master files are being added to the unit...')
            this.awaitUpdateCompletion( resp.data, 'Master files have successfully been added.' )
         }).catch( e => {
            system.setError(e)
            this.updateInProgress = false
         })
      },

      assignMetadata( metadataID, masterFileIDs) {
         const system = useSystemStore()
         let data = {ids: masterFileIDs, metadataID:  parseInt(metadataID,10) }
         axios.post(`${system.jobsURL}/units/${this.detail.id}/masterfiles/metadata`, data).then( () => {
            this.getMasterFiles( this.detail.id )
            system.toastMessage("Assign Metadata Success", 'The selected master files have been assigned new metadata.')
         }).catch( e => {
            system.setError(e)
         })
      },

      assignComponent( componentID, masterFileIDs) {
         const system = useSystemStore()
         let data = {ids: masterFileIDs, componentID:  parseInt(componentID,10) }
         axios.post(`${system.jobsURL}/units/${this.detail.id}/masterfiles/component`, data).then( () => {
            this.getMasterFiles( this.detail.id )
            system.toastMessage("Assign Component Success", `The selected master files have been assigned to component ${componentID}.`)
         }).catch( e => {
            system.setError(e)
         })
      },

      renumberPages( startPage, filenames) {
         const system = useSystemStore()
         let data = {filenames: filenames, startNum:  parseInt(startPage,10) }
         axios.post(`${system.jobsURL}/units/${this.detail.id}/masterfiles/renumber`, data).then( () => {
            this.getMasterFiles( this.detail.id )
            system.toastMessage("Renumber Success", 'The selected master files have been renumbered.')
         }).catch( e => {
            system.setError(e)
            this.updateInProgress = false
         })
      },

      replaceMasterFiles() {
         const system = useSystemStore()
         this.updateInProgress = true
         axios.post(`${system.jobsURL}/units/${this.detail.id}/masterfiles/replace`).then( resp => {
            system.toastMessage("Please Wait", 'Replacing unit master files...')
            this.awaitUpdateCompletion( resp.data, 'Master files have successfully been replaced.' )
         }).catch( e => {
            system.setError(e)
            this.updateInProgress = false
         })
      },

      deleteMasterFiles( filenames ) {
         const system = useSystemStore()
         let payload = {filenames: filenames}
         axios.post(`${system.jobsURL}/units/${this.detail.id}/masterfiles/delete`, payload).then( () => {
            this.getMasterFiles( this.detail.id )
            system.toastMessage("Delete Success", 'The selected master files have been deleted.')
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