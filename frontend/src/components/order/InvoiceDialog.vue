<template>
   <Dialog v-model:visible="ordersStore.showInvoice" :modal="true" header="Invoice" @hide="invoiceClosed()" :style="{width: '650px'}">
      <div v-if=" ordersStore.editInvoice == false">
         <Panel header="Date Information" :style="{marginBottom: '20px'}">
            <dl>
               <DataDisplay label="Date Approved" :value="formatDate(detail.dateApproved)"/>
               <DataDisplay label="Date Customer Notified" :value="formatDate(detail.invoice.dateNoticeSent)"/>
               <DataDisplay label="Date Invoice" :value="formatDate(detail.invoice.invoiceDate)"/>
               <DataDisplay label="Date Fee Paid" :value="formatDate(detail.invoice.dateFeePaid)"/>
               <DataDisplay label="Date Fee Declined" :value="formatDate(detail.invoice.dateFeeDeclined)"/>
            </dl>
         </Panel>
         <Panel header="Billing Information">
            <dl>
               <DataDisplay label="Permanent Non-Payment" :value="formatDate(detail.invoice.permanentNonpayment)"/>
               <DataDisplay label="Fee Amount Paid" :value="formatFee(detail.invoice.feeAmountPaid)"/>
               <DataDisplay label="Transmittal Number" :value="detail.invoice.transmittalNumber"/>
               <DataDisplay label="Notes" :value="detail.invoice.notes"/>
            </dl>
         </Panel>
      </div>
      <template #footer>
         <DPGButton label="OK" autofocus class="p-button-secondary" @click="invoiceClosed()"/>
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

const ordersStore = useOrdersStore()
const { detail } = storeToRefs(ordersStore)


function formatFee( fee ) {
   if (fee) {
      return `$${fee}.00`
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

function invoiceClosed() {
   ordersStore.showInvoice = false
}

</script>

<style scoped lang="scss">
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