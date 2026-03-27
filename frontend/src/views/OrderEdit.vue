<template>
   <h2>
      <span v-if="route.params.id">Order {{route.params.id}}</span>
      <span v-else>New Order</span>
   </h2>
   <div class="edit-form">
      <form @submit="submitChanges">
         <div class="split">
            <FormField id="ostatus" label="Status">
               <Select id="ostatus" v-model="status"  :options="orderStatuses" optionLabel="label" optionValue="value" />   
            </FormField>
            <FormField id="datedue" label="Date Due" :error="errors.dateDue" :required="true">
               <DatePicker id="datedue" v-model="dateDue" showIcon dateFormat="yy-mm-dd" updateModelType="string"/>   
            </FormField>
         </div>
         <FormField id="otitle" label="Order Title">
            <InputText id="otitle" v-model="title" type="text"/>   
         </FormField>
         <FormField id="instruct" label="Special Instructions">
            <Textarea id="instruct" v-model="specialInstructions" rows="5"/>   
         </FormField>
         <FormField id="notes" label="Staff Notes">
            <Textarea id="notes" v-model="staffNotes" rows="5"/>   
         </FormField>
         <FormField  v-if="isFeeRequired" id="fee" label="Fee" :error="errors.fee" :required="true">
            <InputNumber id="fee" v-model="fee" mode="currency" currency="USD" locale="en-US"/>   
         </FormField>
         <div class="split">
             <FormField id="agency" label="Agency">
               <Select id="agency" v-model="agencyID" filter placeholder="Select an agency"
                  :options="agencies" optionLabel="label" optionValue="value" 
               />   
            </FormField>
            <FormField id="customer" label="Customer">
               <Select id="customer" v-model="customerID" filter :error="errors.customerID" :required="true"
                  :options="customers" optionLabel="label" optionValue="value" placeholder="Select a customer"
               />   
            </FormField>
         </div>
         <div class="acts">
            <DPGButton label="Cancel" severity="secondary" @click="cancelEdit()"/>
            <DPGButton label="Save" type="submit" />
         </div>
      </form>
   </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'
import { useSystemStore } from '@/stores/system'
import { useCustomersStore } from '@/stores/customers'
import { onMounted, ref, computed } from 'vue'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import DatePicker from 'primevue/datepicker'
import { useDateFormat } from '@vueuse/core'

import { useForm } from 'vee-validate'
import * as yup from 'yup'
import FormField from '@/components/FormField.vue'

const schema = yup.object().shape({
   dateDue: yup.string().required('Due date is required'),
   customerID: yup.number().min(1, "Customer is required"),
   fee: yup.number().when('feeRequired', {
      is: (value) => value == true,
      then: (schema) => schema.moreThan(0, "Non-zero fee is required").required("Fee is required"),
   }),
})
const { errors, resetForm, handleSubmit, defineField } = useForm({ validationSchema: schema })

const route = useRoute()
const router = useRouter()
const ordersStore = useOrdersStore()
const systemStore = useSystemStore()
const customersStore = useCustomersStore()

const [status] = defineField('status')
const [dateDue] = defineField('dateDue')
const [title] = defineField('title')
const [specialInstructions] = defineField('specialInstructions')
const [staffNotes] = defineField('staffNotes')
const [fee] = defineField('fee')
const [agencyID] = defineField('agencyID')
const [customerID] = defineField('customerID')

const newOrder = ref(false)

const isFeeRequired = computed( () => {
   if ( ordersStore.detail.customer == null ) return false
   if ( customersStore.isExternal(ordersStore.detail.customer.id) == false) return false 
   if ( ordersStore.detail.feeWaived ) return false 
   return true
})

const agencies = computed(() => {
   let out = []
   systemStore.agencies.forEach( a => {
      out.push( {label: a.name, value: a.id} )
   })
   return out
})
const customers = computed(() => {
   let out = []
   customersStore.customers.forEach( a => {
      out.push( {label: `${a.lastName}, ${a.firstName}`, value: a.id} )
   })
   return out
})

const orderStatuses = computed(() => {
   let out = []
   out.push( {label: "Requested", value: "requested"} )
   out.push( {label: "Approved", value: "approved"} )
   out.push( {label: "Completed", value: "completed"} )
   out.push( {label: "Deferred", value: "deferred"} )
   out.push( {label: "Canceled", value: "canceled"} )
   out.push( {label: "Await Fee", value: "await_fee"} )
   return out
})


onMounted( async () =>{
   let orderID = route.params.id
   if (orderID) {
      newOrder.value = false
      await ordersStore.getOrderDetails(orderID)
      document.title = `Edit | Order #${orderID}`
   } else {
      ordersStore.clearDetails()
      newOrder.value = true
      document.title = `New Order`
   }

   await customersStore.getCustomers()

   let val = {
      status: ordersStore.detail.status,
      dateDue: useDateFormat(ordersStore.detail.dateDue, "YYYY-MM-DD").value,
      title: ordersStore.detail.title,
      specialInstructions: ordersStore.detail.specialInstructions,
      staffNotes: ordersStore.detail.staffNotes,
      fee: 0.0,
      agencyID: 0,
      customerID: 0,
      feeRequired: isFeeRequired.value,
   }
   if (ordersStore.detail.fee && !ordersStore.detail.feeWaived ) {
      val.fee = parseFloat(ordersStore.detail.fee).toFixed(2)
   }
   if (ordersStore.detail.agency) {
      val.agencyID = ordersStore.detail.agency.id
   } 
   if (ordersStore.detail.customer) {
      val.customerID = ordersStore.detail.customer.id
   } 
   resetForm({values: val})
})

const cancelEdit = (() => {
   if ( newOrder.value == true) {
      router.push(`/orders`)
   } else {
      router.push(`/orders/${route.params.id}`)
   }
})

const submitChanges = handleSubmit( async (values) => {
   if ( newOrder.value == true) {
      await ordersStore.createOrder( values )
   } else {
      await ordersStore.submitEdit( values )
   }
   if (systemStore.showError == false) {
      router.push(`/orders/${ordersStore.detail.id}`)
   }
})

</script>

<style lang="scss" scoped>
.edit-form {
   width: 50%;
   margin: 20px auto;
   form {
      display: flex;
      flex-direction: column;
      gap: 15px;
      text-align: left;
   }
   .split {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      gap: 15px;
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   gap: 10px;
}
</style>