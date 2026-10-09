<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" title="Invoice">
      <UButton label="View Invoice" color="secondary" @click="viewInvoiceClicked()"/>
      <template #body>
         <div v-if="editInvoice == false" class="column">
            <UCard title="Date Information">
               <dl>
                  <DataDisplay label="Date Invoice" :value="$formatDate(detail.invoice.invoiceDate)"/>
                  <DataDisplay label="Date Fee Paid" :value="$formatDate(detail.invoice.dateFeePaid)"/>
                  <DataDisplay label="Date Fee Declined" :value="$formatDate(detail.invoice.dateFeeDeclined)"/>
               </dl>
            </UCard>
            <UCard title="Billing Information">
               <dl>
                  <DataDisplay label="Fee Amount Paid" :value="formatFee(detail.invoice.feeAmountPaid)"/>
                  <DataDisplay label="Transmittal/Confirmation Number" :value="detail.invoice.transmittalNumber"/>
                  <DataDisplay label="Notes" :value="detail.invoice.notes"/>
               </dl>
            </UCard>
         </div>
         <UForm v-else :state="state" class="column" @submit="submitChanges">
            <div class="row-left gap-4">
               <UFormField name="feeAmountPaid" label="Fee Amount Paid" class="grow">   
                  <UInputNumber v-model="state.feeAmountPaid" class="w-full" :step="0.01" :increment="false" :decrement="false"
                     :format-options="{style: 'currency',currency: 'USD'}" />
               </UFormField>
               <UFormField name="dateFeePaid" label="Date Fee Paid" class="grow">   
                  <UInputDate v-model="state.dateFeePaid" class="w-full"  />
               </UFormField>
            </div>
            <div class="row-left gap-4">
               <UFormField name="dateFeeDeclined" label="Date Fee Declined" class="grow">   
                  <UInputDate v-model="state.dateFeeDeclined" class="w-full"  />
               </UFormField>
               <UFormField name="transmittalNumber" label="Transmittal/Confirmation Number" class="grow">   
                  <UInput v-model="state.transmittalNumber" class="w-full"/>
               </UFormField>
            </div>
           <UFormField name="notes" label="Notes" class="grow">   
               <UTextarea v-model="state.notes" class="w-full"/>
            </UFormField>
            <div class="row-right gap-2">
               <UButton label="Cancel" color="secondary" @click="editInvoice = false"/>
               <UButton label="Save" type="submit" />
            </div>
         </UForm>
      </template>
      <template #footer v-if="editInvoice == false">
         <UButton label="Edit" color="secondary" @click="editInvoiceClicked"/>
         <UButton label="Close" color="secondary" @click="invoiceClosed"/>
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useOrdersStore } from '@/stores/orders'
import DataDisplay from '@/components/DataDisplay.vue'
import { storeToRefs } from 'pinia'
import { parseDate } from '@internationalized/date'

const ordersStore = useOrdersStore()
const { detail } = storeToRefs(ordersStore)

const isOpen = ref(false)
const editInvoice = ref(false)

const state = ref({
   dateFeePaid: null,
   dateFeeDeclined: null,
   feeAmountPaid: 0,
   transmittalNumber: "",
   notes: "",
})

const submitChanges = (async () => {
   const update = {
      dateFeePaid: "",
      dateFeeDeclined: "",
      feeAmountPaid: state.value.feeAmountPaid,
      transmittalNumber: state.value.transmittalNumber,
      notes: state.value.notes,
   }
   if ( state.value.dateFeePaid ) {
      update.dateFeePaid = state.value.dateFeePaid.toString()  
   }
   if ( state.value.dateFeeDeclined ) {
      update.dateFeeDeclined = state.value.dateFeeDeclined.toString()  
   }
   await ordersStore.updateInvoice( update )
   ordersStore.showInvoice = false
})

const formatFee = ( (fee) => {
   if (fee) {
      return `$${fee}`
   }
   return ""
})

const viewInvoiceClicked = (() => {
    editInvoice.value = false
    isOpen.value = true
})

const editInvoiceClicked = (() => {
   editInvoice.value = true
   updateEditData()
})

const updateEditData = (() => {
   state.value = {
      dateFeePaid: null,
      dateFeeDeclined: null,
      feeAmountPaid: 0,
      transmittalNumber: "",
      notes: "",
   }
   if ( ordersStore.detail.invoice ) {
      if (ordersStore.detail.invoice.dateFeePaid) {
         state.value.dateFeePaid = parseDate(ordersStore.detail.invoice.dateFeePaid.split("T")[0])
      }
      if (ordersStore.detail.invoice.dateFeeDeclined) {
         state.value.dateFeeDeclined = parseDate(ordersStore.detail.invoice.dateFeeDeclined.split("T")[0])
      }
      state.value.feeAmountPaid = parseFloat(ordersStore.detail.invoice.feeAmountPaid)
      state.value.transmittalNumber = ordersStore.detail.invoice.transmittalNumber
      state.value.notes = ordersStore.detail.invoice.notes
   }
})

const invoiceClosed = (() => {
   isOpen.value = false
})
</script>

<style scoped lang="scss">
dl {
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