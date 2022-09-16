<template>
   <h2>Order {{route.params.id}}</h2>
   <div class="details" v-if="systemStore.working==false">
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
            <DPGButton v-if="detail.email" label="View Customer Email" class="p-button-secondary" @click="viewEmailClicked()"/>
         </div>
      </Panel>
   </div>
   <Dialog v-model:visible="showEmail" :modal="true" header="Customer Email" @hide="emailClosed()" :style="{width: '650px'}">
      <div v-html="detail.email" class="email"></div>
      <template #footer>
         <DPGButton label="OK" autofocus class="p-button-secondary" @click="emailClosed()"/>
      </template>
   </Dialog>
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


</script>

<style scoped lang="scss">
.details {
   padding: 0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;

   div.status {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: center;
      margin-top: -5px;
   }

   :deep(div.p-panel) {
      margin: 10px;
      flex: 45%;
      text-align: left;
   }
   .actions {
      padding: 15px 0 0 0;
      font-size: 0.8em;
   }
}
div.email {
   padding: 10px;
   font-size: 0.85em;
}
</style>