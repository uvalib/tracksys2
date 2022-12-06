import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useCustomersStore = defineStore('customers', {
	state: () => ({
      customers: [],
      total: 0,
	}),
	getters: {
      isExternal: state => {
         return (customerID) => {
            let c = state.customers.find( c => c.id == customerID)
            if ( c ) {
               return (c.academicStatusID == 1)
            }
            return false
         }
      }
	},
	actions: {
      getCustomers( ) {
         const system = useSystemStore()
         system.working = true
         axios.get( `/api/customers` ).then(response => {
            this.customers = response.data
            this.total = response.data.length
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      addOrUpdateCustomer( data ) {
         const system = useSystemStore()
         system.working = true
         let add = false
         if (data.id == 0) {
            add = true
         }
         data.academicStatusID = data.academicStatus.id
         axios.post( "/api/customers", data ).then(response => {
            if ( add == false) {
               let idx  = this.customers.findIndex( s => s.id == data.id)
               if (idx >= 0) {
                  this.customers.splice(idx, 1, response.data)
               }
            } else {
               // insert new customer at head of current page. the order will get fixed with a later load/sort/filter
               this.customers.unshift(response.data)
            }
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      }
	}
})