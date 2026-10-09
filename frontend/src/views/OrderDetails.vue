<template>
   <h2>
      <span>Order {{route.params.id}}</span>
      <div class="row-right gap-2" v-if="(user.isAdmin || user.isSupervisor)" >
         <UButton label="HathiTrust Metadata Accepted" @click="hathiTrustMetadataAccepted" v-if="canAcceptHathiTrustMetadata"/>
         <UButton label="Submit HathiTrust Packages" @click="submitHathiTrustPackage" v-if="canSubmitHathiTrustPackage"/>
         <UButton label="Package for HathiTrust" @click="packageForHathiTrust" v-if="canPackageHathiTrust"/>
         <HathiTrustMetadataDialog @submit="submitHathiTrustMetadata" :order="detail.id" v-if="canSubmitHathiTrustMetadata"/>
         <UButton label="Flag for HathiTrust" @click="flagForHathiTrust" v-if="canFlagForHathiTrust"/>
         <UButton label="Delete" @click="deleteOrder()" v-if="canDelete"/>
         <UButton label="Edit" @click="editOrder()"/>
      </div>
   </h2>
   <div class="row-left gap-4 pt-8 px-8">
      <div class="column flex-1">
         <UCard title="General Information" class="w-full">
            <dl>
               <DataDisplay label="Status" :value="detail.status">
                  <div class="status">
                     <span :class="`status ${detail.status}`">{{displayStatus(detail.status)}}</span>
                  </div>
               </DataDisplay>
               <DataDisplay v-if="detail.status=='completed'" label="Date Completed" :value="$formatDate(detail.dateCompleted)"/>
               <DataDisplay v-if="detail.customer" label="Customer" value="placeholder">
                  <CustomerPopover :customer="detail.customer" />
               </DataDisplay>
               <DataDisplay v-else label="Customer" value=""/>
               <DataDisplay v-if="detail.agency" label="Agency" :value="detail.agency.name"/>
               <DataDisplay v-else label="Agency" value=""/>
               <DataDisplay label="Title" :value="detail.title"/>
               <DataDisplay label="Special Instructions" :value="detail.specialInstructions"/>
               <DataDisplay label="Staff Notes" :value="detail.staffNotes"/>
               <DataDisplay v-if="detail.processor" label="Order Processor" :value="detail.processor.lastName">
                  <span>{{ detail.processor.firstName }} {{ detail.processor.lastName }}</span>
               </DataDisplay>
            </dl>
         </UCard>
         <UCard class="messages w-full" title="Messages" v-if="hasMessages" >
            <div class="msg" v-if="detail.status== 'requested'">Order is not yet approved. Units must be added and approved before order can be approved.</div>
            <div class="msg" v-if="detail.status== 'deferred'">Order has been deferred.</div>
            <div class="msg" v-if="detail.customer.academicStatusID==1 && !detail.fee && !detail.feeWaived">Either enter a fee, defer or cancel this order.</div>
            <template v-if="detail.status== 'await_fee'">
               <div class="msg">Order is awaiting customer fee payment.</div>
               <div class="msg" v-if="ordersStore.isFeePaid == false">Fee payment information must be added to the invoice.</div>
               <div class="msg" v-if="ordersStore.hasUnitsBeingPrepared">You must approve or cancel units in this order.</div>
            </template>
         </UCard>
      </div>
      <UCard title="Workflow"  class="flex-1">
         <dl>
            <DataDisplay label="Date Submitted" :value="$formatDate(detail.dateSubmitted)"/>
            <DataDisplay label="Date Due" :value="$formatDate(detail.dateDue)"/>
            <template v-if="isExternalCustomer">
               <DataDisplay v-if="detail.feeWaived" label="Date Fee Waived" :value="$formatDate(detail.dateFeeWaived)"/>
               <template v-else>
                  <DataDisplay label="Fee" :value="formatFee(detail.fee)"/>
                  <DataDisplay label="Date Fee Sent to Customer" :value="$formatDate(detail.dateFeeEstimateSent)"/>
               </template>
            </template>
            <DataDisplay v-if="detail.dateCanceled" label="Date Canceled" :value="$formatDate(detail.dateCanceled)"/>
            <DataDisplay v-if="detail.dateDeferred" label="Date Deferred" :value="$formatDate(detail.dateDeferred)"/>
            <DataDisplay label="Date Finalization Started" :value="$formatDateTime(detail.dateFinalizationBegun)"/>
            <DataDisplay label="Date Archiving Complete" :value="$formatDateTime(detail.dateArchivingComplete)"/>
            <DataDisplay label="Date Patron Deliverables Complete" :value="$formatDateTime(detail.datePatronDeliverablesComplete)"/>
            <DataDisplay label="Date Customer Notified" :value="$formatDateTime(detail.dateCustomerNotified)"/>
         </dl>
         <div class="acts-wrap" v-if="user.isAdmin || user.isSupervisor">
            <div class="actions" v-if="detail.status != 'completed'">
               <UButton label="Claim for Processing" color="secondary" @click="claimOrder()" :disabled="isProcessor"/>
               <AssignModal />
            </div>
            <div class="actions" v-if="detail.status == 'await_fee'">
               <SendEmailDialog mode="fee" />
               <UButton label="Customer Declines Fee" color="secondary" @click="declineFeeClicked()"/>
               <UButton label="Customer Paid Fee" color="secondary" :disabled="isPaidDisabled"  @click="payFeeClicked()"/>
            </div>
            <template v-else>
               <div class="actions" v-if="detail.status != 'completed' && detail.status != 'canceled'">
                  <UButton v-if="canWaiveFee" label="Waive Fee" color="secondary" @click="waiveFeeClicked()"/>
                  <UButton v-if="isExternalCustomer && detail.feeWaived == false" label="Send Fee Estimate" color="secondary"
                     :disabled="isSendFeeDisabled" @click="sendFeeEstimateCllicked()"/>
                  <UButton v-if="detail.status == 'deferred'" label="Resume Order" color="secondary" @click="resumeOrderClicked()"/>
                  <UButton v-else label="Defer Order" color="secondary" @click="deferOrderClicked()"/>
                  <UButton label="Approve Order" color="secondary" :disabled="isApproveDisabled" @click="approveOrderClicked()"/>
                  <UButton label="Cancel Order" color="secondary" @click="cancelOrderClicked()"/>
                  <UButton label="Complete Order" color="secondary" :disabled="isCompleteOrderDisabled" @click="completeOrderClicked()"/>
               </div>
            </template>
            <div class="actions" v-if="(detail.status == 'approved' || detail.status == 'completed') && ordersStore.hasPatronDeliverables && detail.email" >
               <UButton label="View Customer Email" color="secondary" @click="viewEmailClicked()" :style="{marginLeft:0}"/>
               <UButton label="Recreate Email" color="secondary" @click="recreateEmailClicked()" />
               <SendEmailDialog mode="order" />
            </div>
            <div class="actions" v-if="ordersStore.hasPatronDeliverables && (detail.status == 'approved' || detail.status == 'completed')">
               <UButton v-if="!detail.email" label="Check Order Completeness" color="secondary" @click="checkOrderComplete()" />
               <UButton v-if="detail.email" label="View Order Summary" color="secondary" @click="viewSummaryClicked()" />
               <UButton v-if="detail.email" label="Recreate Order Summary" color="secondary" @click="recreateSummaryClicked()" />
            </div>
            <div class="actions" v-if="(detail.invoice || detail.fee) && !detail.feeWaived">
               <InvoiceDialog  v-if="detail.invoice" />
            </div>
         </div>
      </UCard>
   </div>
   <div class="details" v-if="ordersStore.items.length> 0">
      <UCard title="Order Details" class="w-full">
         <p>The following is all of the raw data submitted by the patron. Use it to create units or discard it. Once all units have been created and the order approved, this data will be deleted.</p>
         <dl class="item-intended-use">
            <DataDisplay label="Intended Use" :value="ordersStore.items[0].intendedUse.name"/>
            <DataDisplay label="Format" :value="ordersStore.items[0].intendedUse.deliverableFormat"/>
            <DataDisplay label="Resolution" :value="ordersStore.items[0].intendedUse.deliverableResolution"/>
         </dl>
         <Divider />
         <div class="items">
            <div class="item" v-for="item in ordersStore.items" :key="item.id">
               <i v-if="item.converted" class="used fas fa-check-circle"></i>
               <dl>
                  <DataDisplay label="Title" :value="item.title"/>
                  <DataDisplay label="Pages" :value="item.pages"/>
                  <DataDisplay v-if="item.author" label="Author" :value="item.author"/>
                  <DataDisplay v-if="item.callNumber" label="Call Number" :value="item.callNumber"/>
                  <DataDisplay v-if="item.year" label="Year Published" :value="item.year"/>
                  <DataDisplay v-if="item.location" label="Location" :value="item.location"/>
                  <DataDisplay v-if="item.sourceURL" label="Web Link" :value="item.sourceURL"/>
                  <DataDisplay v-if="item.description" label="Description" :value="item.description"/>
               </dl>
               <div class="item-acts">
                  <UButton label="Discard" autofocus color="secondary" @click="discardItem(item)"/>
                  <AddUnitDialog label="Create Unit" :item="item" />
               </div>
            </div>
         </div>
      </UCard>
   </div>
   <div class="details" v-if="systemStore.working==false" >
      <UCard title="Units" class="w-full">
         <RelatedUnits :orderID="detail.id" :units="ordersStore.units" :hathiTrust="canUpdateHathiTrust" :canAdd="canAddUnit"/>
      </UCard>
   </div>
   <Dialog v-model:visible="showEmail" :modal="true" header="Customer Email" @hide="emailClosed()" :style="{width: '650px'}">
      <div v-html="detail.email" class="email"></div>
      <template #footer>
         <UButton label="OK" autofocus color="secondary" @click="emailClosed()"/>
      </template>
   </Dialog>
</template>

<script setup>
import Dialog from 'primevue/dialog'
import { onBeforeMount, ref, computed } from 'vue'
import { useRoute, onBeforeRouteUpdate, useRouter } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useOrdersStore } from '@/stores/orders'
import { useUserStore } from '@/stores/user'
import { useCustomersStore } from '@/stores/customers'
import DataDisplay from '../components/DataDisplay.vue'
import { storeToRefs } from 'pinia'
import InvoiceDialog from '@/components/order/InvoiceDialog.vue'
import RelatedUnits from '../components/related/RelatedUnits.vue'
import Divider from 'primevue/divider'
import SendEmailDialog from '../components/order/SendEmailDialog.vue'
import AddUnitDialog from '../components/order/AddUnitDialog.vue'
import HathiTrustMetadataDialog from '../components/order/HathiTrustMetadataDialog.vue'
import { useConfirm } from "@/composables/useConfirm"
import AssignModal from '../components/order/AssignModal.vue'
import CustomerPopover from '../components/order/CustomerPopover.vue'

const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()
const ordersStore = useOrdersStore()
const user = useUserStore()
const customerStore = useCustomersStore()

const { detail } = storeToRefs(ordersStore)

const showEmail = ref(false)
const customer = ref(null)

const canUpdateHathiTrust = computed( () => {
   return user.isAdmin &&  ordersStore.hathiTrustMetadataCount > 0
})

const canFlagForHathiTrust = computed( () => {
   return user.isAdmin && ordersStore.hasHathiTrustCandidateUnits
})
const canSubmitHathiTrustMetadata = computed( () => {
   return user.isAdmin && ordersStore.hasHathiTrustMetadataCandidate
})
const canAcceptHathiTrustMetadata = computed( () => {
   return user.isAdmin && ordersStore.hasSubmittedHathiTrustMetadata
})
const canPackageHathiTrust = computed( () => {
   return user.isAdmin && ordersStore.hathiTrustPackageCandidate
})
const canSubmitHathiTrustPackage = computed( () => {
   return user.isAdmin && ordersStore.hathiTrustPackageSubmitCandidate
})

const canDelete = computed(() => {
   return (user.isAdmin || user.isSupervisor) && ordersStore.detail.status=='requested' && ordersStore.units.length == 0
})

const canAddUnit = computed(() =>{
   return detail.status != 'completed' && detail.status != 'canceled'
})

const hasMessages = computed(() => {
   if ( ordersStore.detail.id != 0 ) {
      if ( ordersStore.detail.status== 'requested' || ordersStore.detail.status == 'deferred' || ordersStore.detail.status== 'await_fee') return true
      if ( ordersStore.detail.customer.academicStatusID==1 && !ordersStore.detail.fee && !ordersStore.detail.feeWaived) return true
   }
   return false
})

const isProcessor = computed(() => {
   if (!ordersStore.detail.processor ) return false
   return (ordersStore.detail.processor.id == user.ID)
})

const isPaidDisabled = computed(() =>{
   return ordersStore.isFeePaid == false
})

const isCompleteOrderDisabled = computed(() =>{
   let disabled = true
   if ( ordersStore.detail.status == 'approved' ) {
      disabled = false
      ordersStore.units.forEach( u => {
         if (!u.datePatronDeliverablesReady && !u.dateDLDeliverablesReady) {
            disabled = true
         }
      })
   }
   return disabled
})

const isApproveDisabled = computed(() =>{
   if (  ordersStore.detail.status == 'approved' ) return true // already approved; disable
   if (  ordersStore.hasApprovedUnits == false ) return true // no approved untis; disable

   // external unpaid, not waived; disable
   if ( isExternalCustomer.value && (ordersStore.detail.fee == null || ordersStore.isFeePaid == false) && !ordersStore.detail.feeWaived) return true
   return false
})

const canWaiveFee = computed(() => {
   if ( ordersStore.detail.customer == null ) return false
   if ( customerStore.isExternal(ordersStore.detail.customer.id) == false ) return false
   if ( ordersStore.detail.feeWaived ) return false
   return ( user.isAdmin || user.isSupervisor)
})

const isSendFeeDisabled = computed(() => {
   // Only enable send estimate when estimate is populated, fee has not been
   // sent or paid and status is not deferred/canceled/approved
   let feeDisabled = ordersStore.detail.fee == null || ordersStore.detail.dateFeeEstimateSent != null  ||
      ordersStore.detail.status == 'deferred' || ordersStore.detail.status == 'canceled' || ordersStore.detail.status == 'approved'
   return feeDisabled
})

const isExternalCustomer = computed( () => {
   if (ordersStore.detail.customer == null) return false
   return customerStore.isExternal(ordersStore.detail.customer.id)
})

onBeforeRouteUpdate(async (to) => {
   let orderID = to.params.id
   ordersStore.getOrderDetails(orderID)
})

onBeforeMount( async () => {
   let orderID = route.params.id
   document.title = `Order #${orderID}`
   await ordersStore.getOrderDetails(orderID)
   await customerStore.getCustomers()
})

const flagForHathiTrust = (async () => {
   const msg = `Are you sure you want flag all digital collection building units in this order for inclusion in HathiTrust?`
   const resp = await useConfirm("Confirm HathiTrust Inclusion", msg, "Include")
   if (resp) {
      ordersStore.flagForHathiTrust( user.computeID )
   } 
})
const packageForHathiTrust = (async () => {
   const msg = `Are you sure you want package all units in this order for submission to HathiTrust?`
   const resp = await useConfirm("Confirm HathiTrust Package", msg, "Package")
   if (resp) {
      ordersStore.packageForHathiTrust( user.computeID )
   } 
})
const submitHathiTrustPackage = (() => {
   ordersStore.submitHathiTrustPackage( user.computeID )  
})
const hathiTrustMetadataAccepted = (() => {
   ordersStore.hathiTrustMetadataAccepted()
})
const submitHathiTrustMetadata = (( info ) => {
   ordersStore.submitHathiTrustMetadata( user.computeID, info.mode, info.name )
})

const deleteOrder = (async () => {
   const msg = `Are you sure you want delete this order? All data will be lost. This cannot be reversed.`
   const resp = await useConfirm("Confirm Delete Order", msg, "Delete")
   if (resp) {
      await ordersStore.deleteOrder()
      router.push("/orders")
   } 
})

const editOrder = (() => {
   router.push(`/orders/${route.params.id}/edit`)
})

const recreateEmailClicked = (() => {
   ordersStore.recreateEmail()
})

const recreateSummaryClicked = (() => {
   ordersStore.recreateSummary()
})

const viewSummaryClicked = (() => {
   let url = `${systemStore.jobsURL}/orders/${ordersStore.detail.id}/summary`
   window.open(url)
})

const formatFee = (( fee ) => {
   if (fee) {
      let floatFee = parseFloat(fee).toFixed(2)
      return `$${floatFee}`
   }
   return ""
})

const displayStatus = ((id) => {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
})

const viewEmailClicked = (() => {
   showEmail.value = true
})

const emailClosed = (() => {
   showEmail.value = false
})

const discardItem = (async (item) => {
   const msg = `Are you sure you want discard this item? All data will be lost. This cannot be reversed.`
   const resp = await useConfirm("Confirm Discard", msg, "Discard")
   if (resp) {
        await ordersStore.discardItem(item.id)
   } 
})

const waiveFeeClicked = ( async () => {
   const msg = `Waive the fee for this order? This cannot be reversed.`
   const resp = await useConfirm("Confirm Fee Waive", msg, "Waive Fee")
   if (resp) {
       await ordersStore.waiveFee( user.computeID )
   } 
})

const sendFeeEstimateCllicked = (() => {
   ordersStore.sendFeeEstimate( user.computeID )
})

const deferOrderClicked = (() => {
   ordersStore.deferOrder( user.computeID )
})

const resumeOrderClicked = (() => {
   ordersStore.resumeOrder( user.computeID )
})

const approveOrderClicked = (() => {
   ordersStore.approveOrder( user.computeID )
})

const cancelOrderClicked = (async () => {
   const msg = `Are you sure you want cancel this order? All related units and projects will be canceled. This cannot be reversed.`
   const resp = await useConfirm("Confirm Cancel Order", msg, "Cancel Order")
   if (resp) {
       await ordersStore.cancelOrder( user.computeID )
   } 
})

const completeOrderClicked = (() => {
   ordersStore.completeOrder( user.computeID )
})

const payFeeClicked = (() => {
   ordersStore.feeAccepted( user.computeID )
})

const declineFeeClicked = (() => {
   ordersStore.feeDeclined( user.computeID )
})

const checkOrderComplete = (() => {
   ordersStore.checkOrderComplete()
})

const claimOrder = (async () => {
   const msg = `Are you sure you want claim this order for processing?`
   const resp = await useConfirm("Confirm Claim Order", msg, "Claim")
   if (resp) {
      ordersStore.setProcessor( user.ID )
   } 
})

</script>

<style scoped lang="scss">
dl.item-intended-use {
   dd {
      margin: 0 0 5px 0 !important;
   }
}
div.item {
   margin: 15px;
   padding: 10px;
   border: 1px solid var(--uvalib-grey-light);
   border-radius: 5px;
   position: relative;
   i.used {
      font-size: 1.25em;
      position: absolute;
      color: var(--uvalib-green-dark);
      top: 10px;
      left: 10px;
   }
   .item-acts {
      font-size: 0.8em;
      margin: 10px;
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      align-items: flex-start;
      gap: 10px;
   }
}
.details {
   padding: 20px 20px 0 20px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   .no-units {
      padding: 20px 0 0 0;
      h3 {
         text-align: center;
      }
   }
   p {
      margin: 5px;
   }

   dl {
      margin-bottom: 10px !important;
   }

   div.customer, div.status {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: center;
      .name {
         color: var(--uvalib-blue-alt-dark);
         font-weight: 500;
         text-decoration: none;
         display: inline-block;
         cursor: pointer;
         &:hover {
            text-decoration: underline;
         }
      }
   }
}

.acts-wrap {
   border-top: 1px solid var(--uvalib-grey-light);
   padding-top: 15px;
   .actions {
      padding: 5px 0;
      font-size: 0.8em;
      display: flex;
      flex-flow: row wrap;
      justify-content: flex-start;
      gap: 5px;
   }
}

dl {
   margin: 0;
   display: inline-grid;
   grid-template-columns: max-content 1fr;
   grid-column-gap: 10px;
   font-size: 0.9em;
   text-align: left;
   box-sizing: border-box;

   dt {
      font-weight: bold;
      text-align: right;
   }

   dd {
      margin: 0 0 5px 0;
      word-break: break-word;
      -webkit-hyphens: auto;
      -moz-hyphens: auto;
      hyphens: auto;
      white-space: break-spaces;
      margin-inline-start: 5px;
   }
}
div.email {
   padding: 10px;
   font-size: 0.85em;
}
</style>