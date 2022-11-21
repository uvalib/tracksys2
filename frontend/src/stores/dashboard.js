import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useDashboardStore = defineStore('dsahboard', {
	state: () => ({
      dueInOneWeek: 0,
      overdue: 0,
      readyForDelivery: 0
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
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})