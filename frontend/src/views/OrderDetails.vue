<template>
   <h2>Order {{route.params.id}}</h2>
   <div class="details" v-if="systemStore.working==false">
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
               <DataDisplay label="Customer" :value="`${detail.customer.lastName}, ${detail.customer.firstName}`"/>
               <DataDisplay label="Agency" :value="detail.agency.name"/>
               <DataDisplay label="Title" :value="detail.title"/>
               <DataDisplay label="Special Instructions" :value="detail.specialInstructions"/>
               <DataDisplay label="Staff Notes" :value="detail.staffNotes"/>
            </dl>
         </Panel>
         <Panel header="Messages">
            <div class="msg" v-if="detail.status== 'requested'">Order is not yet approved.</div>
            <div class="msg" v-if="detail.status== 'deferred'">Order has been deferred.</div>
            <div class="msg" v-if="detail.status== 'await_fee'">Order is awaiting customer fee payment.</div>
            <div class="msg" v-if="detail.customer.academicStatusID==1 && !detail.fee">Either enter a fee, defer or cancel this order.</div>
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
         <div class="actions">
            <DPGButton v-if="detail.email" label="View Customer Email" class="p-button-secondary" @click="viewEmailClicked()" :style="{marginLeft:0}"/>
            <DPGButton v-if="detail.invoice" label="View Invoice" class="p-button-secondary" @click="viewInvoiceClicked()"/>
            <DPGButton v-else-if="detail.fee" label="Create Invoice" class="p-button-secondary" @click="createInvoiceClicked()"/>
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
import { onBeforeMount, ref } from 'vue'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useOrdersStore } from '@/stores/orders'
import Panel from 'primevue/panel'
import DataDisplay from '../components/DataDisplay.vue'
import { storeToRefs } from 'pinia'
import dayjs from 'dayjs'
import InvoiceDialog from '@/components/order/InvoiceDialog.vue'
import RelatedUnits from '../components/related/RelatedUnits.vue'
import Divider from 'primevue/divider'

const route = useRoute()
const systemStore = useSystemStore()
const ordersStore = useOrdersStore()

const { detail } = storeToRefs(ordersStore)

const showEmail = ref(false)
const events = ref(null)

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
      let d = dayjs(dateStr)
      return d.format("YYYY-MM-DD")
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

</script>

<style scoped lang="scss">
:deep(dl.item-intended-use) {
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
   .actions {
      padding: 15px 0 0 0;
      font-size: 0.8em;
      display: flex;
      flex-flow: row wrap;
      justify-content: flex-start;
      button.p-button {
         margin-left: 10px;
      }
   }
}
div.email {
   padding: 10px;
   font-size: 0.85em;
}
</style>