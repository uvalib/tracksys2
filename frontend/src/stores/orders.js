import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'
import dayjs from 'dayjs'

export const useOrdersStore = defineStore('orders', {
	state: () => ({
      orders: [],
      total: 0,
      searchOpts: {
         start: 0,
         limit: 30,
         filter: "active",
         sortField: "id",
         sortOrder: "desc",
         query: "",
      }
	}),
	getters: {
	},
	actions: {
      getOrders() {
         const system = useSystemStore()
         system.working = true
         let so = this.searchOpts
         let url = `/api/orders?filter=${so.filter}&start=${so.start}&limit=${so.limit}&by=${so.sortField}&order=${so.sortOrder}`
         if ( so.query != "") {
            url += `&q=${encodeURIComponent(so.query)}`
         }
         axios.get( url ).then(response => {
            this.orders = []
            response.data.orders.forEach( js => {
               js.dateDue =  dayjs(js.dateDue).format("YYYY-MM-DD")
               js.dateSubmitted =  dayjs(js.dateSubmitted).format("YYYY-MM-DD")
               if (js.dateCustomerNotified) {
                  js.dateCustomerNotified =  dayjs(js.dateCustomerNotified).format("YYYY-MM-DD")
               }
               if (js.dateArchivingComplete) {
                  js.dateArchivingComplete =  dayjs(js.dateArchivingComplete).format("YYYY-MM-DD")
               }
               let fee = js.fee
               if (fee.Valid) {
                  js.fee = `$${fee.Float64}`
               } else {
                  js.fee = null
               }
               js.customerName = `${js.customer.lastName}, ${js.customer.firstName}`
               this.orders.push(js)
            })
            this.total = response.data.total
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})