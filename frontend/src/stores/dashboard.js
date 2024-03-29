import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useDashboardStore = defineStore('dashboard', {
	state: () => ({
      dueInOneWeek: 0,
      overdue: 0,
      readyForDelivery: 0,
      asRequests: 0,
      asReviews: 0,
      asRejections: 0
	}),
	getters: {
	},
	actions: {
      getStatistics( ) {
         const system = useSystemStore()
         axios.get( `/api/dashboard` ).then(response => {
            this.dueInOneWeek = response.data.dueInOneWeek
            this.overdue = response.data.overdue
            this.readyForDelivery = response.data.readyForDelivery
            this.asRequests = response.data.asRequests
            this.asReviews = response.data.asReviews
            this.asRejections = response.data.asRejections
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})