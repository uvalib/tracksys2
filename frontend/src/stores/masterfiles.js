import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useMasterFilesStore = defineStore('masterfiles', {
	state: () => ({
      details: {},
      thumbURL: "",
      viewerURL: "",
      orderID: 0,
      nextID: 0,
      prevID: 0,
      replaceInProgress: false,
   }),
	getters: {
      hasText: state => {
         return state.details.transcription
      },
      isOCRCandidate: state => {
         if ( state.details.metadata == null) return false
         if ( state.details.metadata.ocrHint == null ) return false
         return state.details.metadata.ocrHint.ocrCandidate
      }
	},
	actions: {
      async getDetails( masterFileID ) {
         const system = useSystemStore()
         system.working = true
         this.nextID = 0
         this.prevID = 0
         return axios.get( `/api/masterfiles/${masterFileID}` ).then(response => {
            this.details = response.data.masterFile
            this.thumbURL = response.data.thumbURL
            this.viewerURL = response.data.viewerURL
            this.orderID = response.data.orderID
            if ( response.data.nextID) {
               this.nextID = response.data.nextID
            }
            if ( response.data.prevID) {
               this.prevID = response.data.prevID
            }
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      audit() {
         const system = useSystemStore()
         let payload = {type: "id", data: `${this.details.id}`}
         axios.post(`${system.jobsURL}/audit`, payload).then( ( resp ) => {
            this.details.audit = resp.data
         }).catch( e => {
            system.setError(e)
         })
      },
      async submitEdit( edit ) {
         const system = useSystemStore()
         system.working = true
         return axios.post( `/api/masterfiles/${this.details.id}/update`, edit ).then( (resp) => {
            this.details = resp.data
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      downloadFromArchive( computeID ) {
         const system = useSystemStore()
         let payload = {computeID: computeID, filename: this.details.filename}
         let url = `${system.jobsURL}/units/${this.details.unitID}/copy`
         axios.post(url, payload).then( () => {
            system.toastMessage("Archive Download", `${this.details.filename} is being downloaded from the archive.`)
         }).catch( e => {
            system.setError(e)
         })
      },
      addTag( tag ) {
         axios.post(`/api/masterfiles/${this.details.id}/tags?tag=${tag.id}`).then( () => {
            this.details.tags.push(tag)
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
         })
      },
      removeTag( tgt ) {
         axios.delete(`/api/masterfiles/${this.details.id}/tags?tag=${tgt.id}`).then( () => {
            let idx = this.details.tags.findIndex( t => t.tag == tgt.tag)
            if (idx > -1) {
               this.details.tags.splice(idx,1)
            }
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
         })
      },
      ocr() {
         const system = useSystemStore()
         let payload = {type: "masterfile", id: this.details.id}
         axios.post(`${system.jobsURL}/ocr`, payload).then( () => {
            system.toastMessage("OCR Started", `OCR nas been started form ${this.details.filename}. Check the Job Statuses page for updates.`)
         }).catch( e => {
            system.setError(e)
         })
      },
      replace() {
         const system = useSystemStore()
         this.replaceInProgress = true
         axios.post(`${system.jobsURL}/units/${this.details.unitID}/masterfiles/replace`).then( resp => {
            system.toastMessage("Please Wait", 'Replacing unit master file...')
            let jobID = resp.data
            var tid = setInterval( ()=> {
               axios.get(`${system.jobsURL}/jobs/${jobID}`).then( resp => {
                  let status = resp.data.status
                  if (status == 'failure') {
                     clearInterval(tid)
                     this.replaceInProgress = false
                     system.setError(`Replace failed: ${resp.data.error}. Check the job status logs for more information.`)
                  } else if (status == 'finished') {
                     clearInterval(tid)
                     this.replaceInProgress = false
                     window.location.reload()
                  }
               }).catch( e => {
                  system.setError(e)
                  this.replaceInProgress = false
                  clearInterval(tid)
               })
            }, 1000)
         }).catch( e => {
            system.setError(e)
            this.replaceInProgress = false
         })
      },
   },
})