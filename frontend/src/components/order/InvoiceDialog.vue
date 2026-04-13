<template>
   <Dialog v-model:visible="ordersStore.showInvoice" :modal="true" header="Invoice" @show="invoiceOpened"
      @hide="invoiceClosed" :style="{width: '650px'}" :closable="false">
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
      <form v-else @submit="submitChanges">
         <div class="split">
            
            <FormField id="feepaid" label="Fee Amount Paid" >
               <InputNumber id="feepaid" v-model="feeAmountPaid" mode="currency" currency="USD" locale="en-US"/>   
            </FormField>
            <FormField id="paiddate" label="Date Fee Paid">
               <DatePicker id="paiddate" v-model="dateFeePaid" showIcon dateFormat="yy-mm-dd" updateModelType="string"/>   
            </FormField>
         </div>
         <div class="split">
            <FormField id="declinedate" label="Date Fee Declined">
               <DatePicker id="declinedate" v-model="dateFeeDeclined" showIcon dateFormat="yy-mm-dd" updateModelType="string"/>   
            </FormField>
            <FormField id="confnum" label="Transmittal/Confirmation Number" >
               <InputText id="confnum" v-model="transmittalNumber" />   
            </FormField>
         </div>
         <FormField id="notes" label="Notes" >
            <Textarea id="notes" v-model="notes" rows="5" />   
         </FormField>
         <div class="acts">
            <DPGButton label="Cancel" severity="secondary" @click="invoiceClosed"/>
            <DPGButton label="Save" type="submit" />
         </div>
      </form>
      <template #footer v-if="ordersStore.editInvoice == false">
         <DPGButton label="Edit" autofocus severity="secondary" @click="editInvoice()"/>
         <DPGButton label="OK" autofocus severity="secondary" @click="invoiceClosed"/>
      </template>
   </Dialog>
</template>

<script setup>
import Dialog from 'primevue/dialog'
import { useOrdersStore } from '@/stores/orders'
import DataDisplay from '@/components/DataDisplay.vue'
import Panel from 'primevue/panel'
import { useDateFormat } from '@vueuse/core'
import { storeToRefs } from 'pinia'
import { useForm } from 'vee-validate'
import FormField from '@/components/FormField.vue'
import Textarea from 'primevue/textarea'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import DatePicker from 'primevue/datepicker'

const ordersStore = useOrdersStore()
const { detail } = storeToRefs(ordersStore)

const { resetForm, handleSubmit, defineField } = useForm({})

const [dateFeePaid] = defineField('dateFeePaid')
const [dateFeeDeclined] = defineField('dateFeeDeclined')
const [feeAmountPaid] = defineField('feeAmountPaid')
const [transmittalNumber] = defineField('transmittalNumber')
const [notes] = defineField('notes')

const submitChanges = handleSubmit(values => {
   ordersStore.updateInvoice( values )
   ordersStore.showInvoice = false
})

const formatFee = ( (fee) => {
   if (fee) {
      return `$${fee}`
   }
   return ""
})

const invoiceOpened = (() => {
   if (ordersStore.editInvoice) {
      updateEditData()
   }
})

const editInvoice = (() => {
   ordersStore.editInvoice = true
   updateEditData()
})

const updateEditData = (() => {
   let val = {
      dateFeePaid: "",
      dateFeeDeclined: "",
      feeAmountPaid: 0,
      transmittalNumber: "",
      notes: "",
   }
   if ( ordersStore.detail.invoice ) {
      if (ordersStore.detail.invoice.dateFeePaid) {
         val.dateFeePaid = useDateFormat(ordersStore.detail.invoice.dateFeePaid, "YYYY-MM-DD").value
      }
      if (ordersStore.detail.invoice.dateFeeDeclined) {
         val.dateFeeDeclined = useDateFormat(ordersStore.detail.invoice.dateFeeDeclined, "YYYY-MM-DD").value
      }
      val.feeAmountPaid = ordersStore.detail.invoice.feeAmountPaid
      val.transmittalNumber = ordersStore.detail.invoice.transmittalNumber
      val.notes = ordersStore.detail.invoice.notes
   }
   resetForm({values: val})
})

const invoiceClosed = (() => {
   ordersStore.showInvoice = false
})
</script>

<style scoped lang="scss">
.split {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: baseline;
   gap: 20px;
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