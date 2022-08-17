import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useStaffStore = defineStore('staff', {
	state: () => ({
      staff: [],
      total: 0,
      searchOpts: {
         start: 0,
         limit: 30,
      }
	}),
	getters: {
	},
	actions: {
      getStaff() {
         const system = useSystemStore()
         system.working = true
         let url = `/api/staff?start=${this.searchOpts.start}&limit=${this.searchOpts.limit}`
         axios.get( url ).then(response => {
            this.staff = response.data.staff
            this.total = response.data.total
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})