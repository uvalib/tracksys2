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
         sortField: "lastName",
         sortOrder: "asc",
      }
	}),
	getters: {
	},
	actions: {
      getStaff( queryStr ) {
         const system = useSystemStore()
         system.working = true
         let so = this.searchOpts
         let url = `/api/staff?start=${so.start}&limit=${so.limit}&by=${so.sortField}&order=${so.sortOrder}`
         if (queryStr != "") {
            url += `&q=${encodeURIComponent(queryStr)}`
         }
         axios.get( url ).then(response => {
            this.staff = response.data.staff
            this.total = response.data.total
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      addOrUpdateStaff( staffData ) {
         const system = useSystemStore()
         system.working = true
         let add = false
         if (staffData.id == 0) {
            add = true
         }
         axios.post( "/api/staff", staffData ).then(response => {
            if ( add == false) {
               let idx  = this.staff.findIndex( s => s.id == staffData.id)
               if (idx >= 0) {
                  this.staff.splice(idx, 1, response.data)
               }
            } else {
               // insert new staff at head of current page. the order will get fixed with a later load/sort/filter
               this.staff.unshift(response.data)
            }
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      }
	}
})