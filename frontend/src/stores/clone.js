import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useCloneStore = defineStore('clone', {
   state: () => ({
      sourceUnits: [],
      masterFiles: [],
      initialized: false
	}),
	getters: {
	},
	actions: {
      getSourceUnits( destUnitID ) {
         const system = useSystemStore()
         this.sourceUnits = []
         this.initialized = false
         axios.get( `/api/units/${destUnitID}/clone-sources` ).then(response => {
            this.sourceUnits = response.data
         }).catch( e => {
            system.setError(e)
         }). finally( ()=> {
            this.initialized = true
         })
      },
     async  getMasterFiles( unitID ) {
         const system = useSystemStore()
         this.masterFiles = []
         return axios.get( `/api/units/${unitID}/masterfiles` ).then(response => {
            this.masterFiles = response.data
         }).catch( e => {
            system.setError(e)
         })
      },
   }
})