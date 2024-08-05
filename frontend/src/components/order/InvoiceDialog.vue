<template>
   <Dialog v-model:visible="ordersStore.showInvoice" :modal="true" header="Invoice" @show="invoiceOpened()"
      @hide="invoiceClosed()" :style="{width: '650px'}" :closable="false">
      <div v-if="ordersStore.editInvoice == false">
         <Panel header="Date Information" :style="{marginBottom: '20px'}">
            <dl>
               <DataDisplay label="Date Invoice" :value="$formatDate(detail.invoice.invoiceDate)"/>
               <DataDisplay label="Date Fee Paid" :value="$formatDate(detail.invoice.dateFeePaid)"/>
               <DataDisplay label="Date Fee Declined" :value="$formatDate(detail.invoice.dateFeeDeclined)"/>
            </dl>
         </Panel>
         <Panel header="Billing Information">
            <dl>
               <DataDisplay label="Fee Amount Paid" :value="formatFee(detail.invoice.feeAmountPaid)"/>
               <DataDisplay label="Transmittal/Confirmation Number" :value="detail.invoice.transmittalNumber"/>
               <DataDisplay label="Notes" :value="detail.invoice.notes"/>
            </dl>
         </Panel>
      </div>
      <FormKit v-else type="form" id="customer-detail" :actions="false" @submit="submitChanges">
         <div class="split">
            <FormKit label="Date Fee Paid" type="date" v-model="edit.dateFeePaid"/>
            <div class="sep"></div>
            <FormKit label="Fee Amount Paid" type="text" v-model="edit.feeAmountPaid"/>
         </div>
         <div class="split">
            <FormKit label="Date Fee Declined" type="date" v-model="edit.dateFeeDeclined"/>
            <div class="sep"></div>
            <FormKit label="Transmittal/Confirmation Number" type="text" v-model="edit.transmittalNumber"/>
         </div>
         <FormKit label="Notes" type="textarea" rows="5" v-model="edit.notes"/>
         <div class="acts">
            <DPGButton label="Cancel" severity="secondary" @click="invoiceClosed()"/>
            <FormKit type="submit" label="Save" wrapper-class="submit-button" />
         </div>
      </FormKit>
      <template #footer v-if="ordersStore.editInvoice == false">
         <DPGButton label="Edit" autofocus severity="secondary" @click="editInvoice()"/>
         <DPGButton label="OK" autofocus severity="secondary" @click="invoiceClosed()"/>
      </template>
   </Dialog>
</template>

<script setup>
import Dialog from 'primevue/dialog'
import { useOrdersStore } from '@/stores/orders'
import DataDisplay from '@/components/DataDisplay.vue'
import Panel from 'primevue/panel'
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { ref } from 'vue'

const ordersStore = useOrdersStore()
const { detail } = storeToRefs(ordersStore)

const edit = ref({
   dateFeePaid: "",
   dateFeeDeclined: "",
   feeAmountPaid: null,
   transmittalNumber: "",
   notes: ""
})

function submitChanges() {
   ordersStore.updateInvoice( edit.value )
   ordersStore.showInvoice = false
}

function formatFee( fee ) {
   if (fee) {
      return `$${fee}`
   }
   return ""
}

function invoiceOpened() {
   if (ordersStore.editInvoice) {
      updateEditData()
   }
}

function editInvoice() {
   ordersStore.editInvoice = true
   updateEditData()
}

function updateEditData() {
   if ( ordersStore.detail.invoice ) {
      if (ordersStore.detail.invoice.dateFeePaid) {
         edit.value.dateFeePaid = dayjs(ordersStore.detail.invoice.dateFeePaid).format("YYYY-MM-DD")
      }
      if (ordersStore.detail.invoice.dateFeeDeclined) {
         edit.value.dateFeeDeclined = dayjs(ordersStore.detail.invoice.dateFeeDeclined).format("YYYY-MM-DD")
      }
      edit.value.feeAmountPaid = ordersStore.detail.invoice.feeAmountPaid
      edit.value.transmittalNumber = ordersStore.detail.invoice.transmittalNumber
      edit.value.notes = ordersStore.detail.invoice.notes
   }
}

function invoiceClosed() {
   ordersStore.showInvoice = false
}

</script>

<style scoped lang="scss">
.split {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: baseline;
   :deep(.formkit-outer) {
      flex-grow: 0.6;
   }
   .sep {
      display: inline-block;
      width: 20px;
   }
   .checkbox {
      margin-top: 20px;
      label {
         color: var(--uvalib-text-dark);
         font-weight: normal;
         display: block;
      }
      input[type=checkbox] {
         display: block;
         width: 18px;
         height: 18px;
         margin: 10px 0 0 0;
      }
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 15px 0 5px 5px;
   gap: 10px;
}

:deep(dl) {
   margin: 10px 30px 0 30px;
   display: inline-grid;
   grid-template-columns: max-content 2fr;
   grid-column-gap: 10px;
   font-size: 0.9em;
   text-align: left;
   box-sizing: border-box;

   dt {
      font-weight: bold;
      text-align: right;
   }

   dd {
      margin: 0 0 10px 0;
      word-break: break-word;
      -webkit-hyphens: auto;
      -moz-hyphens: auto;
      hyphens: auto;
      white-space: break-spaces;
   }
}
</style>