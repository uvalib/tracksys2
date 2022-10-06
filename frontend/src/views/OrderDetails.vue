<template>
   <h2>
      <span>Order {{route.params.id}}</span>
   </h2>
   <DPGButton label="Edit" class="p-button-secondary edit" @click="editOrder()"/>
   <div class="details">
      <div class="left">
         <Panel header="General Information">
            <dl>
               <DataDisplay label="Status" :value="detail.status">
                  <div class="status">
                     <span :class="`status ${detail.status}`">{{displayStatus(detail.status)}}</span>
                     <DPGButton icon="pi pi-info-circle" class="p-button-rounded p-button-text"
                        @click="toggleEvents" aria-haspopup="true" aria-controls="events-panel" />
                  </div>
               </DataDisplay>
               <OverlayPanel ref="events" id="events-panel" :showCloseIcon="true">
                  <DataTable :value="ordersStore.events" ref="eventsTable" dataKey="id" :lazy="false"
                     stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
                  >
                     <Column header="Date">
                        <template #body="slotProps">
                           {{formatDateTime(slotProps.data.createdAt)}}
                        </template>
                     </Column>
                     <Column header="User">
                        <template #body="slotProps">
                           {{slotProps.data.staffMember.firstName}} {{slotProps.data.staffMember.lastName}}
                        </template>
                     </Column>
                     <Column field="details" header="Details" />
                  </DataTable>
               </OverlayPanel>
               <DataDisplay v-if="detail.status=='completed'" label="Date Completed" :value="formatDateTime(detail.dateCompleted)"/>
               <DataDisplay v-if="detail.customer" label="Customer" :value="`${detail.customer.lastName}, ${detail.customer.firstName}`"/>
               <DataDisplay v-else label="Customer" value=""/>
               <DataDisplay v-if="detail.agency" label="Agency" :value="detail.agency.name"/>
               <DataDisplay v-else label="Agency" value=""/>
               <DataDisplay label="Title" :value="detail.title"/>
               <DataDisplay label="Special Instructions" :value="detail.specialInstructions"/>
               <DataDisplay label="Staff Notes" :value="detail.staffNotes"/>
            </dl>
         </Panel>
         <Panel header="Messages" v-if="hasMessages">
            <div class="msg" v-if="detail.status== 'requested'">Order is not yet approved.</div>
            <div class="msg" v-if="detail.status== 'deferred'">Order has been deferred.</div>
            <div class="msg" v-if="detail.customer.academicStatusID==1 && !detail.fee">Either enter a fee, defer or cancel this order.</div>
            <template v-if="detail.status== 'await_fee'">
               <div class="msg">Order is awaiting customer fee payment.</div>
               <div class="msg" v-if="ordersStore.isFeePaid == false">Fee payment information must be added to the invoice.</div>
               <div class="msg" v-if="ordersStore.hasUnitsBeingPrepared">You must approve or cancel units in this order.</div>
            </template>
         </Panel>
      </div>
      <Panel header="Workflow">
         <dl>
            <DataDisplay label="Date Submitted" :value="formatDate(detail.dateSubmitted)"/>
            <DataDisplay label="Date Due" :value="formatDate(detail.dateDue)"/>
            <DataDisplay label="Fee" :value="formatFee(detail.fee)"/>
            <DataDisplay label="Date Deferred" :value="formatDate(detail.dateDeferred)"/>
            <DataDisplay label="Date Fee Sent to Customer" :value="formatDate(detail.dateFeeEstimateSent)"/>
            <DataDisplay label="Date Finalization Started" :value="formatDateTime(detail.dateFinalizationBegun)"/>
            <DataDisplay label="Date Archiving Complete" :value="formatDateTime(detail.dateArchivingComplete)"/>
            <DataDisplay label="Date Patron Deliverables Complete" :value="formatDateTime(detail.datePatronDeliverablesComplete)"/>
            <DataDisplay label="Date Customer Notified" :value="formatDateTime(detail.dateCustomerNotified)"/>
         </dl>
         <div class="acts-wrap" v-if="user.isAdmin || user.isSupervisor">
            <div class="actions" v-if="detail.status == 'await_fee'">
               <SendEmailDialog mode="fee" />
               <DPGButton label="Customer Declines Fee" class="p-button-secondary" @click="declineFeeClicked()"/>
               <DPGButton label="Customer Paid Fee" class="p-button-secondary" :disabled="isPaidDisabled"  @click="payFeeClicked()"/>
            </div>
            <template v-else>
               <div class="actions" v-if="detail.status != 'completed' && detail.status != 'canceled'">
                  <DPGButton v-if="isExternalCustomer" label="Send Fee Estimate" class="p-button-secondary right-pad"
                     :disabled="isSendFeeDisabled" @click="sendFeeEstimateCllicked()"/>
                  <DPGButton v-if="detail.status == 'deferred'" label="Resume Order" class="p-button-secondary" @click="resumeOrderClicked()"/>
                  <DPGButton v-else label="Defer Order" class="p-button-secondary" @click="deferOrderClicked()"/>
                  <DPGButton label="Approve Order" class="p-button-secondary" :disabled="isApproveDisabled" @click="approveOrderClicked()"/>
                  <DPGButton label="Cancel Order" class="p-button-secondary" @click="cancelOrderClicked()"/>
                  <DPGButton label="Complete Order" class="p-button-secondary" :disabled="isCompleteOrderDisabled" @click="completeOrderClicked()"/>
               </div>
            </template>
            <div class="actions" v-if="(detail.status == 'approved' || detail.status == 'completed') && ordersStore.hasPatronDeliverables && detail.email">
               <DPGButton label="View Customer Email" class="p-button-secondary" @click="viewEmailClicked()" :style="{marginLeft:0}"/>
               <DPGButton label="Recreate Email" class="p-button-secondary" @click="recreateEmailClicked()" />
               <SendEmailDialog mode="order" />
            </div>
            <div class="actions" v-if="ordersStore.hasPatronDeliverables">
               <DPGButton label="View Customer PDF" class="p-button-secondary" @click="viewPDFClicked()" />
               <DPGButton v-if="detail.email" label="Recreate Customer PDF" class="p-button-secondary" @click="recreatePDFClicked()" />
            </div>
            <div class="actions" v-if="detail.invoice || detail.fee">
               <DPGButton v-if="detail.invoice" label="View Invoice" class="p-button-secondary" @click="viewInvoiceClicked()"/>
               <DPGButton v-else-if="detail.fee" label="Create Invoice" class="p-button-secondary" @click="createInvoiceClicked()"/>
            </div>
         </div>
      </Panel>
   </div>
   <div class="details" v-if="ordersStore.items.length> 0">
      <Panel header="Order Details">
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
                  <DPGButton label="Discard" autofocus class="p-button-secondary" @click="discardItem(item)"/>
                  <DPGButton label="Create Unit" autofocus class="p-button-secondary" @click="createUnitFromItem(item)"/>
               </div>
            </div>
         </div>
      </Panel>
   </div>
   <div class="details" v-if="ordersStore.units.length> 0">
      <Panel header="Units">
         <RelatedUnits :units="ordersStore.units" />
      </Panel>
   </div>
   <Dialog v-model:visible="showEmail" :modal="true" header="Customer Email" @hide="emailClosed()" :style="{width: '650px'}">
      <div v-html="detail.email" class="email"></div>
      <template #footer>
         <DPGButton label="OK" autofocus class="p-button-secondary" @click="emailClosed()"/>
      </template>
   </Dialog>
   <InvoiceDialog />
</template>


<script setup>
import Dialog from 'primevue/dialog'
import OverlayPanel from 'primevue/overlaypanel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { onBeforeMount, ref, computed } from 'vue'
import { useRoute, onBeforeRouteUpdate, useRouter } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useOrdersStore } from '@/stores/orders'
import { useUserStore } from '@/stores/user'
import { useCustomersStore } from '@/stores/customers'
import Panel from 'primevue/panel'
import DataDisplay from '../components/DataDisplay.vue'
import { storeToRefs } from 'pinia'
import dayjs from 'dayjs'
import InvoiceDialog from '@/components/order/InvoiceDialog.vue'
import RelatedUnits from '../components/related/RelatedUnits.vue'
import Divider from 'primevue/divider'
import SendEmailDialog from '../components/order/SendEmailDialog.vue'

const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()
const ordersStore = useOrdersStore()
const user = useUserStore()
const customerStore = useCustomersStore()

const { detail } = storeToRefs(ordersStore)

const showEmail = ref(false)
const events = ref(null)

const hasMessages = computed(() => {
   if ( ordersStore.detail.id != 0 ) {
      if ( ordersStore.detail.status== 'requested' || ordersStore.detail.status == 'deferred' || ordersStore.detail.status== 'await_fee') return true
      if ( ordersStore.detail.customer.academicStatusID==1 && !ordersStore.detail.fee) return true
   }
   return false
})

const isPaidDisabled = computed(() =>{
   return ordersStore.isFeePaid == false || ordersStore.hasUnitsBeingPrepared
})

const isCompleteOrderDisabled = computed(() =>{
   return ordersStore.detail.status != 'approved'
})

const isApproveDisabled = computed(() =>{
   if (  ordersStore.detail.status == 'approved' ) return true // already approved; disable
   if (  ordersStore.hasApprovedUnits == false ) return true // no approved untis; disable
   if ( isExternalCustomer() && (ordersStore.detail.fee == null || ordersStore.isFeePaid == false)) return true // external unpaid; disable
   return false
})

const isSendFeeDisabled = computed(() => {
   // Only enable send estimate when estimate is populated, fee has not been
   // sent or paid and status is not deferred/approved
   let feeDisabled = ordersStore.detail.fee == null || ordersStore.detail.dateFeeEstimateSent != null  ||
      ordersStore.detail.status == 'deferred' || ordersStore.detail.status == 'approved'
                       // order.fee_paid?
   return feeDisabled
})

const isExternalCustomer = computed( () => {
   if (ordersStore.detail.customer == null) return false
   return customerStore.isExternal(ordersStore.detail.customer.id)
})

function editOrder() {
   router.push(`/orders/${route.params.id}/edit`)
}

function toggleEvents(e) {
   events.value.toggle(e)
}

onBeforeRouteUpdate(async (to) => {
   let orderID = to.params.id
   ordersStore.getOrderDetails(orderID)
})

onBeforeMount(() => {
   let orderID = route.params.id
   ordersStore.getOrderDetails(orderID)
})
function recreateEmailClicked() {
   ordersStore.recreateEmail()
}
function recreatePDFClicked() {
   ordersStore.recreatePDF()
}
function viewPDFClicked() {
   let url = `${systemStore.jobsURL}/orders/${ordersStore.detail.id}/pdf`
   window.open(url)
}

function formatFee( fee ) {
   if (fee) {
      return `$${fee}.00`
   }
   return ""
}

function displayStatus( id) {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
}

function formatDateTime( dateStr ) {
   if (dateStr) {
      let d = dayjs(dateStr)
      return d.format("YYYY-MM-DD HH:mm")
   }
   return ""
}

function formatDate( dateStr ) {
   if (dateStr) {
      return dateStr.split("T")[0]
   }
   return ""
}

function viewEmailClicked() {
   showEmail.value = true
}
function emailClosed() {
   showEmail.value = false
}
function viewInvoiceClicked() {
    ordersStore.editInvoice = false
    ordersStore.showInvoice = true
}
function createInvoiceClicked() {
    ordersStore.editInvoice = true
    ordersStore.showInvoice = true
}

function createUnitFromItem(item) {
   alert("create from "+item.id)
}
function discardItem(item) {
   alert("discard "+item.id)
}

function sendFeeEstimateCllicked() {
   ordersStore.sendFeeEstimate( user.ID )
}
function deferOrderClicked() {
   ordersStore.deferOrder( user.ID )
}
function resumeOrderClicked() {
   ordersStore.resumeOrder( user.ID )
}
function approveOrderClicked() {
   //
}
function cancelOrderClicked() {
   //
}
function completeOrderClicked() {
   //
}
function payFeeClicked() {
   ordersStore.feeAccepted( user.ID )
}
function declineFeeClicked() {
   ordersStore.feeDeclined( user.ID )
}

</script>

<style scoped lang="scss">
:deep(dl.item-intended-use) {
   dd {
      margin: 0 0 5px 0 !important;
   }
}
button.p-button-secondary.edit {
   position: absolute;
   right:15px;
   top: 15px;
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
      text-align: right;
      button.p-button {
         margin-left: 10px;
      }
   }
}
.details {
   padding: 0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   p {
      margin: 5px;
   }

   dl {
      margin-bottom: 10px !important;
   }
   div.status {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: center;
      margin-top: -5px;
   }

   div.left {
      margin: 10px;
      flex: 45%;
      text-align: left;
      .msg {
         margin: 5px 0;
         font-size: 0.9em;
      }
   }
   :deep(div.p-panel) {
      margin: 10px;
      flex: 45%;
      text-align: left;
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
         :deep(button.p-button) {
            margin-right: 10px;
         }
      }
   }
}
div.email {
   padding: 10px;
   font-size: 0.85em;
}
</style>