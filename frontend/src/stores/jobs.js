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
         page: 1,
         rowsPerPage: 30,
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
         let url = `/api/jobs?page=${this.searchOpts.page}&limit=${this.searchOpts.rowsPerPage}`
         axios.get( url ).then(response => {
            this.jobs = []
            response.data.jobs.forEach( js => {
               this.jobs.push({
                  id: js.id,
                  name: js.name,
                  associatedObject: `${js.originatorType} ${js.originatorID}`,
                  status: js.status,
                  warnings: js.failures,
                  error: js.error,
                  startedAt: dayjs(js.startedAt).format("YYYY-MM-DD hh:mm A"),
                  finishedAt: dayjs(js.finishedAt).format("YYYY-MM-DD hh:mm A"),
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