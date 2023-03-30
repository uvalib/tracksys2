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
      prevID: 0
   }),
	getters: {
      hasText: state => {
         return state.details.transcription
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
      }
   },
})