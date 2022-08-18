import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useCustomersStore = defineStore('customers', {
	state: () => ({
      customers: [],
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
      getCustomers( queryStr ) {
         const system = useSystemStore()
         system.working = true
         let so = this.searchOpts
         let url = `/api/customers?start=${so.start}&limit=${so.limit}&by=${so.sortField}&order=${so.sortOrder}`
         if (queryStr != "") {
            url += `&q=${encodeURIComponent(queryStr)}`
         }
         axios.get( url ).then(response => {
            this.customers = response.data.customers
            this.total = response.data.total
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
         axios.post( "/api/customers", data ).then(response => {
            if ( add == false) {
               let idx  = this.customers.findIndex( s => s.id == data.id)
               if (idx >= 0) {
                  this.customers.splice(idx, 1, response.data)
               }
            } else {
               // insert new customer at head of current page. the order will get fixed with a later load/sort/filter
               this.customers.unshift(response.data)
               this.customers.pop()
            }
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      }
	}
})