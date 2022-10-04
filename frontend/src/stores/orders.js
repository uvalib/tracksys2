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
      },
      detail: {
         id: 0,
         status: "",
         title: "",
         dateDue: "",
         customer: null,
         agency: null,
         fee: null,
         invoice: null,
         email: "",
         staffNotes: "",
         specialInstructions: "",
         dateSubitted: "",
         dateApproved: "",
         dateDeferred: "",
         dateCanceled: "",
         dateCustomerNotified: "",
         datePatronDeliverablesComplete: "",
         dateArchivingComplete: "",
         dateFinalizationBegun: "",
         dateFeeEstimateSent: "",
         dateCompleted: "",
      },
      events: [],
      items: [],
      units: [],
      showInvoice: false,
      editInvoice: false
	}),
	getters: {
      hasPatronDeliverables: state => {
         let hasDeliverables = false
         state.units.some( u=>{
            /// only units that are NOT for digital collection building can have patron deliverables
            if (u.intendedUse && u.intendedUse.description != "Digital Collection Building" ) {
               hasDeliverables = true
            }
            return hasDeliverables==true
         })
         return hasDeliverables
      },
   },
	actions: {
     async  getOrderDetails(orderID) {
         if ( this.detail.id == orderID ) return
         const system = useSystemStore()
         system.working = true
         return axios.get( `/api/orders/${orderID}` ).then(response => {
            this.detail = response.data.order
            this.events = response.data.events
            this.items = response.data.items
            this.units = response.data.units
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      async submitEdit( edit ) {
         const system = useSystemStore()
         system.working = true
         return axios.post( `/api/orders/${this.detail.id}/update`, edit ).then( (resp) => {
            this.detail.status = resp.data.status
            this.detail.dateDue = resp.data.dateDue
            this.detail.title = resp.data.title
            this.detail.specialInstructions = resp.data.specialInstructions
            this.detail.staffNotes = resp.data.staffNotes
            this.detail.fee = resp.data.fee
            this.detail.agency = resp.data.agency
            this.detail.customer = resp.data.customer

            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      async sendEmail( toCustomer, toAlt = false, altAddress = "") {
         const system = useSystemStore()
         system.working = true
         let url = `${system.jobsURL}/orders/${this.detail.id}/email/send`
         if (toCustomer) {
            await axios.post( url )
         }
         if ( toAlt ) {
            url = `${url}?alt=${altAddress}`
            await axios.post( url  )
         }
         system.toastMessage("Email Sent", "Order email has been sent to the selected recipients")
         system.working = false
      },
      recreateEmail() {
         const system = useSystemStore()
         system.working = true
         let url = `${system.jobsURL}/orders/${this.detail.id}/email`
         axios.post( url ).then( () => {
            system.toastMessage("Email Recreated", "New email generated, but not sent.")
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
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