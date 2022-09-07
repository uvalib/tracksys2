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
      getDetails( masterFileID ) {
         const system = useSystemStore()
         system.working = true
         axios.get( `/api/masterfiles/${masterFileID}` ).then(response => {
            this.details = response.data.masterFile
            this.thumbURL = response.data.thumbURL
            this.viewerURL = response.data.viewerURL
            this.orderID = response.data.orderID
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      }
   }
})