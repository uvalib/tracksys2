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
         filters: [],
         sortField: "id",
         sortOrder: "desc",
         query: "",
         ownerID: -1
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
      editInvoice: false,
      lookupHits: [],
      totalLookupHits: 0,
	}),
	getters: {
      isFeePaid: state => {
         if ( state.detail.invoice == null) return false
         return state.detail.invoice.dateFeePaid != null
      },
      hasUnitsBeingPrepared: state => {
         // Returns units belonging to current order that are not ready to proceed with digitization and would prevent an order from being approved.
         // Only units whose unit_status = 'approved' or 'canceled' are removed from consideration by this method.
         let beingPrepared = false
         state.units.some( u=>{
            if (u.status == 'unapproved' || u.status == 'condition' || u.status == 'copyright') {
               beingPrepared = true
            }
            return beingPrepared==true
         })
         return beingPrepared
      },
      hasApprovedUnits: state => {
         let hasApproved = false
         state.units.some( u=>{
            if (u.status == 'approved') {
               hasApproved = true
            }
            return hasApproved==true
         })
         return hasApproved
      },
      hasPatronDeliverables: state => {
         let hasDeliverables = false
         state.units.some( u=>{
            /// only units that are NOT for digital collection building can have patron deliverables
            if (u.intendedUse && u.intendedUse.id != 110 ) {
               hasDeliverables = true
            }
            return hasDeliverables==true
         })
         return hasDeliverables
      },
      filtersAsQueryParam: state => {
         let out = []
         state.searchOpts.filters.forEach( fv => out.push(`${fv.field}|${fv.match}|${fv.value}`) )
         return JSON.stringify(out)
      }
   },
	actions: {
      getOrders() {
         const system = useSystemStore()
         system.working = true
         let so = this.searchOpts
         let url = `/api/orders?filters=${this.filtersAsQueryParam}&start=${so.start}&limit=${so.limit}&by=${so.sortField}&order=${so.sortOrder}`
         if ( so.query != "") {
            url += `&q=${encodeURIComponent(so.query)}`
         }
         if ( this.ownerID > -1) {
            url += `&owner=${this.ownerID}`
         }
         axios.get( url ).then(response => {
            this.orders = []
            response.data.orders.forEach( js => {
               js.dateDue = js.dateDue.split("T")[0]
               js.dateSubmitted =  js.dateSubmitted.split("T")[0]
               if (js.dateCustomerNotified) {
                  js.dateCustomerNotified =  js.dateCustomerNotified.split("T")[0]
               }
               if (js.dateArchivingComplete) {
                  js.dateArchivingComplete =  js.dateArchivingComplete.split("T")[0]
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
      async lookup( query ) {
         const system = useSystemStore()
         let url = `/api/search?scope=orders&q=${encodeURIComponent(query)}&start=0&limit=30`
         return axios.get(url).then(response => {
            this.lookupHits = response.data.orders.hits
            this.totalLookupHits = response.data.orders.total
         }).catch( e => {
            system.setError(e)
         })
      },
      clearDetails()  {
         this.detail.id = 0
         this.detail.status = ""
         this.detail.title = ""
         this.detail.dateDue = ""
         this.detail.customer = null
         this.detail.agency = null
         this.detail.fee = null
         this.detail.invoice = null
         this.detail.email = ""
         this.detail.staffNotes = ""
         this.detail.specialInstructions = ""
         this.detail.dateSubitted = ""
         this.detail.dateApproved = ""
         this.detail.dateDeferred = ""
         this.detail.dateCanceled = ""
         this.detail.dateCustomerNotified = ""
         this.detail.datePatronDeliverablesComplete = ""
         this.detail.dateArchivingComplete = ""
         this.detail.dateFinalizationBegun = ""
         this.detail.dateFeeEstimateSent = ""
         this.detail.dateCompleted = ""
      },
      async addUnit( metadataID, unitInfo, itemID=0) {
         const system = useSystemStore()
         system.working = true
         let req = {metadataID: metadataID, intendedUseID: unitInfo.intendedUseID, sourceURL: unitInfo.sourceURL,
            specialInstructions: unitInfo.specialInstructions, staffNotes: unitInfo.staffNotes,
            completeScan: unitInfo.completeScan, throwAway: unitInfo.throwAway, includeInDL: unitInfo.includeInDL}
         if (itemID != 0) {
            req.ItemID = itemID
         }
         return axios.post( `/api/orders/${this.detail.id}/units`, req ).then(response => {
            this.units.push( response.data )
            if (itemID != 0) {
               let item = this.items.find( i => i.id == itemID)
               if ( item ) {
                  item.converted = true
               }
            }
            system.working = false
         }).catch( e => {
            system.setError(e)
         })

      },
      async discardItem(itemID) {
         const system = useSystemStore()
         system.working = true
         return axios.delete( `/api/orders/${this.detail.id}/items/${itemID}` ).then( () => {
            let itemIdx = this.items.findIndex( i => i.id == itemID)
            if (itemIdx > -1) {
               this.items.splice(itemIdx, 1)
            }
            system.toastMessage("Item Discarded", 'Order item has been discarded')
            system.working = false
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },

      async getOrderDetails(orderID) {
         if ( this.detail.id == orderID ) return
         const system = useSystemStore()
         system.working = true
         return axios.get( `/api/orders/${orderID}` ).then(response => {
            this.detail = response.data.order
            if (this.detail.fee != null) {
               this.detail.fee = `${this.detail.fee}`
            }
            if (this.detail.invoice && this.detail.invoice.feeAmountPaid) {
               this.detail.invoice.feeAmountPaid = `${this.detail.invoice.feeAmountPaid}`
            }
            this.events = response.data.events
            this.items = response.data.items
            this.units = response.data.units
            system.working = false
         }).catch( e => {
            this.clearDetails()
            if (e.response && e.response.status == 404) {
               system.working = false
               this.router.push("/not_found")
            } else {
               system.setError(e)
            }
         })
      },
      async deleteOrder() {
         const system = useSystemStore()
         system.working = true
         return axios.delete( `/api/orders/${this.detail.id}` ).then( () => {
            this.clearDetails()
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },

      updateInvoice( edit ) {
         const system = useSystemStore()
         system.working = true
         axios.post( `/api/invoices/${this.detail.invoice.id}/update`, edit ).then( () => {
            this.detail.invoice.dateFeePaid = edit.dateFeePaid
            this.detail.invoice.feeAmountPaid = edit.feeAmountPaid
            this.detail.invoice.dateFeeDeclined = edit.dateFeeDeclined
            this.detail.invoice.transmittalNumber = edit.transmittalNumber
            this.detail.invoice.notes = edit.notes
            system.toastMessage("Invoice Updated", "Order invoice has been updated")
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },

      async createOrder( data ) {
         const system = useSystemStore()
         system.working = true
         if (data.fee) {
            data.fee = `${data.fee}`
         }
         return axios.post( `/api/orders`, data ).then( (response ) => {
            this.detail = response.data
            if (this.detail.fee) {
               this.detail.fee = `${this.detail.fee}`
            }
            if (this.detail.invoice && this.detail.invoice.feeAmountPaid) {
               this.detail.invoice.feeAmountPaid = `${this.detail.invoice.feeAmountPaid}`
            }
            system.toastMessage("Order Created", "Order has been created")
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },

      async submitEdit( edit ) {
         const system = useSystemStore()
         system.working = true
         if (edit.fee) {
            edit.fee = `${edit.fee}`
         }
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

      async sendEmail( staffID, toCustomer, toAlt = false, altAddress = "") {
         const system = useSystemStore()
         system.working = true
         let url = `${system.jobsURL}/orders/${this.detail.id}/email/send`
         if (toCustomer) {
            await axios.post( url )
            await axios.post( `/api/orders/${this.detail.id}/complete?staff=${staffID}` ).then( (resp) => {
               this.detail.dateCompleted = resp.data.dateCompleted
               this.detail.dateArchivingComplete = resp.data.dateArchivingComplete
               this.detail.status = resp.data.status
               system.toastMessage("Email Sent", `Email has been sent to ${this.detail.customer.email} and order marked as complete.`)
            })
         }
         if ( toAlt ) {
            url = `${url}?alt=${altAddress}`
            await axios.post( url  )
            system.toastMessage("Email Sent", `Order email has been sent to ${altAddress}`)
         }
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

      async recreatePDF() {
         const system = useSystemStore()
         system.working = true
         let url = `${system.jobsURL}/orders/${this.detail.id}/pdf`
         await axios.post( url ).then( () => {
            system.toastMessage("PDF Recreated", "Customer PDF has been recreated")
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },

      async resendFeeEstimate( staffID,  toCustomer, toAlt = false, altAddress = "") {
         const system = useSystemStore()
         system.working = true
         let url = `${system.jobsURL}/orders/${this.detail.id}/fees?staff=${staffID}&resend=1`
         if (toCustomer) {
            await axios.post( url )
         }
         if ( toAlt ) {
            url = `${url}&alt=${altAddress}`
            await axios.post( url  )
         }
         system.toastMessage("Email Sent", "Fee estimate email has been resent to the selected recipients")
         system.working = false
      },

      sendFeeEstimate( staffID) {
         const system = useSystemStore()
         system.working = true
         let url = `${system.jobsURL}/orders/${this.detail.id}/fees?staff=${staffID}`
         axios.post( url ).then( () => {
            system.toastMessage("A fee estimate email has been sent to the customer")
            this.detail.status = "await_fee"
            this.detail.dateFeeEstimateSent = dayjs(new Date()).format("YYYY-MM-DD")
            let tgtO = this.orders.find( o => o.id == this.detail.id)
            if (tgtO ) {
               tgtO.status = "await_fee"
               tgtO.dateFeeEstimateSent = this.detail.dateFeeEstimateSent
            }
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      feeAccepted( staffID ) {
         const system = useSystemStore()
         system.working = true
         let url = `/api/orders/${this.detail.id}/fee/accept?staff=${staffID}`
         axios.post( url ).then( (resp) => {
            this.detail.dateApproved = resp.data.dateApproved
            this.detail.status = resp.data.status
            let tgtO = this.orders.find( o => o.id == this.detail.id)
            if (tgtO ) {
               tgtO.status = this.detail.status
               tgtO.dateApproved = this.detail.dateApproved
            }
            this.items = []
            system.toastMessage("Fee Accepted", "Fee accepted and order approved")
            system.working = false
         }).catch( e => {
            if (e.response) {
               system.setError("Unable to accept fee: "+e.response.data)
            } else {
               system.setError(e)
            }
         })
      },
      feeDeclined( staffID ) {
         const system = useSystemStore()
         system.working = true
         let url = `/api/orders/${this.detail.id}/fee/decline?staff=${staffID}`
         axios.post( url ).then( (resp) => {
            this.detail.dateCanceled = resp.data.dateCanceled
            this.detail.status = resp.data.status
            this.detail.invoice.dateFeeDeclined = resp.data.dateCanceled
            let tgtO = this.orders.find( o => o.id == this.detail.id)
            if (tgtO ) {
               tgtO.status = this.detail.status
               tgtO.dateCanceled = this.detail.dateCanceled
            }
            system.toastMessage("Fee Declined", "Fee declined and order canceled")
            system.working = false
         }).catch( e => {
            if (e.response) {
               system.setError("Unable to decline feev: "+e.response.data)
            } else {
               system.setError(e)
            }
         })
      },
      approveOrder( staffID ) {
         const system = useSystemStore()
         system.working = true
         let url = `/api/orders/${this.detail.id}/approve?staff=${staffID}`
         axios.post( url ).then( (resp) => {
            this.detail.dateApproved = resp.data.dateApproved
            this.detail.status = resp.data.status
            let tgtO = this.orders.find( o => o.id == this.detail.id)
            if (tgtO ) {
               tgtO.status = this.detail.status
               tgtO.dateApproved = this.detail.dateApproved
            }
            this.items = []
            system.toastMessage("Order Approved", "Order has been approved")
            system.working = false
         }).catch( e => {
            if (e.response) {
               system.setError("Unable to approve order: "+e.response.data)
            } else {
               system.setError(e)
            }
         })
      },
      cancelOrder( staffID ) {
         const system = useSystemStore()
         system.working = true
         let url = `/api/orders/${this.detail.id}/cancel?staff=${staffID}`
         axios.post( url ).then( (resp) => {
            this.detail.dateCanceled = resp.data.dateCanceled
            this.detail.status = resp.data.status
            let tgtO = this.orders.find( o => o.id == this.detail.id)
            if (tgtO ) {
               tgtO.status = this.detail.status
               tgtO.dateCanceled = this.detail.dateCanceled
            }
            system.toastMessage("Order Canceled", "Order has been canceled")
            system.working = false
         }).catch( e => {
            if (e.response) {
               system.setError("Unable to cancel order: "+e.response.data)
            } else {
               system.setError(e)
            }
         })
      },
      deferOrder( staffID ) {
         const system = useSystemStore()
         system.working = true
         let url = `/api/orders/${this.detail.id}/defer?staff=${staffID}`
         axios.post( url ).then( (resp) => {
            this.detail.dateDeferred = resp.data.dateDeferred
            this.detail.status = resp.data.status
            let tgtO = this.orders.find( o => o.id == this.detail.id)
            if (tgtO ) {
               tgtO.status = this.detail.status
               tgtO.dateDeferred = this.detail.dateDeferred
            }
            system.toastMessage("Order Deferred", "Order has been deferred")
            system.working = false
         }).catch( e => {
            if (e.response) {
               system.setError("Unable to defer order: "+e.response.data)
            } else {
               system.setError(e)
            }
         })
      },
      completeOrder( staffID ) {
         const system = useSystemStore()
         system.working = true
         let url = `/api/orders/${this.detail.id}/complete?staff=${staffID}`
         axios.post( url ).then( (resp) => {
            this.detail.dateCompleted = resp.data.dateCompleted
            this.detail.dateArchivingComplete = resp.data.dateArchivingComplete
            this.detail.status = resp.data.status
            system.toastMessage("Order Complete", "Order has been completed.")
            system.working = false
         }).catch( e => {
            if (e.response) {
               system.setError("Unable to mark order complete: "+e.response.data)
            } else {
               system.setError(e)
            }
         })
      },
      resumeOrder( staffID ) {
         const system = useSystemStore()
         system.working = true
         let url = `/api/orders/${this.detail.id}/resume?staff=${staffID}`
         axios.post( url ).then( (resp) => {
            this.detail.status = resp.data.status
            let tgtO = this.orders.find( o => o.id == this.detail.id)
            if (tgtO ) {
               tgtO.status = this.detail.status
            }
            system.toastMessage("Order Resumed", "Order has been reactivated")
            system.working = false
         }).catch( e => {
            if (e.response) {
               system.setError("Unable to resume order: "+e.response.data)
            } else {
               system.setError(e)
            }
         })
      },
      setProcessor( staffID ) {
         const system = useSystemStore()
         system.working = true
         let url = `/api/orders/${this.detail.id}/processor?staff=${staffID}`
         axios.post( url ).then( (response) => {
            this.detail = response.data
            if (this.detail.fee != null) {
               this.detail.fee = `${this.detail.fee}`
            }
            if (this.detail.invoice && this.detail.invoice.feeAmountPaid) {
               this.detail.invoice.feeAmountPaid = `${this.detail.invoice.feeAmountPaid}`
            }
            system.working = false
         }).catch( e => {
            if (e.response) {
               system.setError("Unable to set order processor: "+e.response.data)
            } else {
               system.setError(e)
            }
         })
      },
      checkOrderComplete() {
         const system = useSystemStore()
         system.working = true
         let url = `${system.jobsURL}/orders/${this.detail.id}/check`
         axios.post( url ).then( () => {
            system.toastMessage("Order Checked", `Order ${this.detail.id} has been checked. See details in the job status log.`)
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      setTargetOwner( ownerID) {
         this.ownerID = ownerID
      },
      clearTargetOwner() {
         this.ownerID = -1
      }
	}
})