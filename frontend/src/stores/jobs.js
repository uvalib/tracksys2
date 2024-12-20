import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useJobsStore = defineStore('jobs', {
	state: () => ({
      jobs: [],
      totalJobs: 0,
      details: {status: "",  error: "", associatedObject:"", events: []},
      searchOpts: {
         start: 0,
         limit: 30,
         query: ""
      }
	}),
	getters: {
	},
	actions: {
      getJobs(showWorking=true ) {
         const system = useSystemStore()
         if ( showWorking ) system.working = true

         let url = `/api/jobs?start=${this.searchOpts.start}&limit=${this.searchOpts.limit}`
         if ( this.searchOpts.query != "") {
            url += `&q=${this.searchOpts.query}`
         }
         axios.get( url ).then(response => {
            this.jobs = []
            response.data.jobs.forEach( js => {
               let obj = `${js.originatorType} ${js.originatorID}`
               if (!js.originatorType) {
                  obj = "None"
               }
               this.jobs.push({
                  id: js.id,
                  name: js.name,
                  associatedObject: obj,
                  status: js.status,
                  warnings: js.failures,
                  error: js.error,
                  startedAt: js.startedAt,
                  finishedAt: js.finishedAt,
               })
            })
            this.totalJobs = response.data.total
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      getJobDetails( jobID ) {
         const system = useSystemStore()
         system.working = true
         axios.get(`/api/jobs/${jobID}`).then(response => {
            this.details.events = []
            this.details.error = response.data.error
            this.details.status = response.data.status
            let levels =["info", "warning", "error", "fatal"]
            response.data.events.forEach( evt => {
               this.details.events.push({
                  id: evt.id, jobID: evt.jobID, level: levels[evt.level],
                  text: evt.text, timeStamp: evt.createdAt
               })
            })
            let job = this.jobs.find( j => j.id == jobID)
            if (job) {
               this.details.associatedObject = job.associatedObject
            }
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
      async deleteJobs( delIDs ) {
         const system = useSystemStore()
         system.working = true
         await axios.delete(`/api/jobs/`, {data: {jobs: delIDs}}).then(response => {
            response.data.jobs.forEach( jobID => {
               let idx = this.jobs.findIndex( j => j.id == jobID)
               if (idx >= 0) {
                  this.jobs.splice(idx, 1);
               }
            })
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      }
	}
})