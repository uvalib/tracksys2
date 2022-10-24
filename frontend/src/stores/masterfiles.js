import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useMasterFilesStore = defineStore('masterfiles', {
	state: () => ({
      details: {},
      thumbURL: "",
      viewerURL: "",
      orderID: 0
   }),
	getters: {
	},
	actions: {
      async getDetails( masterFileID ) {
         const system = useSystemStore()
         system.working = true
         return axios.get( `/api/masterfiles/${masterFileID}` ).then(response => {
            this.details = response.data.masterFile
            this.thumbURL = response.data.thumbURL
            this.viewerURL = response.data.viewerURL
            this.orderID = response.data.orderID
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
   }
})