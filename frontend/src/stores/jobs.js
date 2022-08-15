import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'
import dayjs from 'dayjs'

export const useJobsStore = defineStore('jobs', {
	state: () => ({
      jobs: [],
      totalJobs: 0,
      details: {status: "",  error: "", events: []},
      searchOpts: {
         start: 0,
         limit: 30,
         sortBy: 'startedAt',
         sortType: 'desc',
      }
	}),
	getters: {
	},
	actions: {
      getJobs() {
         const system = useSystemStore()
         system.working = true
         let url = `/api/jobs?start=${this.searchOpts.start}&limit=${this.searchOpts.limit}`
         axios.get( url ).then(response => {
            this.jobs = []
            response.data.jobs.forEach( js => {
               let obj = `${js.originatorType} ${js.originatorID}`
               if (!js.originatorType) {
                  obj = "None"
               }
               let finished = "N/A"
               if (js.finishedAt ) {
                  finished = dayjs(js.finishedAt).format("YYYY-MM-DD hh:mm A")
               }
               this.jobs.push({
                  id: js.id,
                  name: js.name,
                  associatedObject: obj,
                  status: js.status,
                  warnings: js.failures,
                  error: js.error,
                  startedAt: dayjs(js.startedAt).format("YYYY-MM-DD hh:mm A"),
                  finishedAt: finished,
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
                  text: evt.text, timeStamp: dayjs(evt.finishedAt).format("YYYY-MM-DD HH:mm:ss")
               })
            })
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})