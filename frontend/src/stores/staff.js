import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useStaffStore = defineStore('staff', {
	state: () => ({
      staff: [],
	}),
	getters: {
	},
	actions: {
      getStaff() {
         const system = useSystemStore()
         system.working = true
         axios.get( '/api/staff' ).then(response => {
            this.staff = response.data
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      getAdmins() {
         const system = useSystemStore()
         this.$reset()
         let url = `/api/staff?admins=1`
         axios.get( url ).then(response => {
            this.staff = response.data
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