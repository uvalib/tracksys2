<template>
   <h2>
      <span v-if="route.params.id">Order {{route.params.id}}</span>
      <span v-else>New Order</span>
   </h2>
   <UForm :schema="schema" :state="state" class="column gap-4 w-full max-w-150 mx-auto! mt-10!" @submit="submitChanges">
      <div class="row-between gap-4">
         <UFormField name="status" label="Status" class="grow">   
            <USelect v-model="state.status" :items="orderStatuses" class="w-full"/>   
         </UFormField>
         <UFormField label="Date Due" :required="true" class="grow">   
            <DatePicker v-model="state.dateDue" class="w-full"/>   
         </UFormField>
      </div>
      <UFormField name="title" label="Title" class="grow">   
         <UInput  v-model="state.title" class="w-full"/>    
      </UFormField>
      <UFormField name="specialInstructions" label="Special Instructions" class="grow">   
         <UTextarea v-model="state.specialInstructions" class="w-full"/>    
      </UFormField>
      <UFormField name="staffNotes" label="Staff Notes" class="grow">   
         <UTextarea v-model="state.staffNotes" class="w-full"/>    
      </UFormField>
      <UFormField  v-if="isFeeRequired" name="fee" label="Fee" class="grow">   
         <UInputNumber v-model="state.fee" class="w-full" :step="0.01" :increment="false" :decrement="false"
            :format-options="{style: 'currency',currency: 'USD'}" />
      </UFormField>
      <div class="row-between gap-4">
         <UFormField name="agencyID" label="Agency" class="grow">   
            <USelectMenu v-model="state.agencyID" :items="agencies" class="w-full" value-key="value" placeholder="Select an agency..." virtualize/>   
         </UFormField>   
         <UFormField name="customerID" label="Customer" class="grow">   
            <USelectMenu v-model="state.customerID" :items="customers" class="w-full" value-key="value"  placeholder="Select a customer..." virtualize/>
         </UFormField>  
      </div>
      <div class="row-right gap-2 pt-4">
         <UButton label="Cancel" color="secondary" @click="cancelEdit()"/>
         <UButton label="Save" type="submit" />
      </div>
   </UForm>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'
import { useSystemStore } from '@/stores/system'
import { useCustomersStore } from '@/stores/customers'
import { onMounted, ref, computed } from 'vue'
import DatePicker from '@/components/DatePicker.vue'
import { toCalendarDate, parseAbsolute } from '@internationalized/date'

import * as yup from 'yup'

const route = useRoute()
const router = useRouter()
const ordersStore = useOrdersStore()
const systemStore = useSystemStore()
const customersStore = useCustomersStore()

const schema = yup.object({
   status: yup.string().required('Status is required'),
   dateDue: yup.string().required('Due date is required'),
   customerID: yup.number().min(1, "Customer is required"),
})

const state = ref({
   status: "",
   dateDue: null,
   title: "",
   specialInstructions: "",
   staffNotes: "",
   fee: 0.0,
   agencyID: null,
   customerID: null,
})

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

   state.value = {
      status: ordersStore.detail.status,
      dateDue: null,
      title: ordersStore.detail.title,
      specialInstructions: ordersStore.detail.specialInstructions,
      staffNotes: ordersStore.detail.staffNotes,
      fee: 0.0,
      agencyID: null,
      customerID: null,
      feeRequired: isFeeRequired.value,
   }
   if (ordersStore.detail.fee && !ordersStore.detail.feeWaived ) {
      state.value.fee = parseFloat(ordersStore.detail.fee)
   }
   if (ordersStore.detail.agency) {
      state.value.agencyID = ordersStore.detail.agency.id
   } 
   if (ordersStore.detail.customer) {
      state.value.customerID = ordersStore.detail.customer.id
   } 
   if (ordersStore.detail.dateDue != "") {
      state.value.dateDue = toCalendarDate( parseAbsolute(ordersStore.detail.dateDue) )
   } 
})

const cancelEdit = (() => {
   if ( newOrder.value == true) {
      router.push(`/orders`)
   } else {
      router.push(`/orders/${route.params.id}`)
   }
})

const submitChanges = (async () => {
   const data = {
      status: state.value.status,
      dateDue: state.value.dateDue.toString(),
      title: state.value.title,
      specialInstructions: state.value.specialInstructions,
      staffNotes: state.value.staffNotes,
      fee: state.value.fee,
      agencyID: state.value.agencyID,
      customerID: state.value.customerID,
   }

   if ( newOrder.value == true) {
      await ordersStore.createOrder( data )
   } else {
      await ordersStore.submitEdit( data )
   }
   if (systemStore.showError == false) {
      router.push(`/orders/${ordersStore.detail.id}`)
   }
   console.log(data)
})

</script>

<style lang="scss" scoped>
.edit-form {
   width: 50%;
   margin: 20px auto;
   .split {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      gap: 15px;
   }
}
</style>