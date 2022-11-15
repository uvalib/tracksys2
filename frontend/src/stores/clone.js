import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import { useUnitsStore } from './units'
import axios from 'axios'

export const useCloneStore = defineStore('clone', {
   state: () => ({
      sourceUnits: [],
      masterFiles: [],
      initialized: false,
      status: "pending"
	}),
	getters: {
      inProgress: state => {
         return state.status == "cloning"
      },
	},
	actions: {
      getSourceUnits( destUnitID ) {
         const system = useSystemStore()
         this.sourceUnits = []
         this.initialized = false
         this.status = "pending"
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
      async  cloneMasterFiles( destUnitID, masterFiles ) {
         const system = useSystemStore()
         var data = []
         masterFiles.forEach( mf => {
            let dataIdx = data.findIndex( d => d.unitID == mf.unitID)
            if (dataIdx > -1) {
               let rec =  data[dataIdx]
               rec.masterfiles.push( {id: mf.id, title: mf.title} )
            } else {
               let rec = {unitID: mf.unitID, masterfiles:[]};
               rec.masterfiles.push( {id: mf.id, title: mf.title} )
               data.push(rec);
            }
         })
         this.status = "cloning"
         axios.post(`${system.jobsURL}/units/${destUnitID}/masterfiles/clone`, data).then( resp => {
            system.toastMessage("Please Wait", 'Master files are being cloned...')
            this.awaitCloneCompletion( resp.data )
         }).catch( e => {
            system.setError(e)
            this.status = "failed"
         })
      },
      awaitCloneCompletion( jobID ) {
         const system = useSystemStore()
         var tid = setInterval( ()=> {
            axios.get(`${system.jobsURL}/jobs/${jobID}`).then( resp => {
               let status = resp.data.status
               if (status == 'failure') {
                  clearInterval(tid)
                  this.cloning = false
                  system.setError(`Clone failed: ${resp.data.error}. Check the job status logs for more information.`)
                  this.status = "failed"
               } else if (status == 'finished') {
                  clearInterval(tid)
                  this.status = "success"
                  const unitsStore = useUnitsStore()
                  unitsStore.flagAsReorder()
               }
            }).catch( e => {
               system.setError(e)
               this.status = "failed"
               clearInterval(tid)
            })
         }, 1000)
      },
   }
})