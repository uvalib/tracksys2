import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'
import dayjs from 'dayjs'

export const useJobsStore = defineStore('jobs', {
	state: () => ({
      jobs: [],
      totalJobs: 0,
	}),
	getters: {
	},
	actions: {
      getJobs() {
         const system = useSystemStore()
         system.working = true
         axios.get("/api/jobs").then(response => {
            this.jobs = []
            response.data.jobs.forEach( js => {
               this.jobs.push({
                  id: js.id,
                  name: js.name,
                  associatedObject: `${js.originatorType} ${js.originatorID}`,
                  status: js.status,
                  warnings: js.failures,
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
	}
})