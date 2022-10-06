<template>
   <h2>
      <span v-if="route.params.id">Order {{route.params.id}}</span>
      <span v-else>New Order [{{newOrder}}]</span>
   </h2>
   <div class="edit-form">
      <FormKit type="form" id="customer-detail" :actions="false" @submit="submitChanges">
         <div class="split">
            <FormKit label="Status" type="select" v-model="edited.status" :options="orderStatuses" required/>
            <div class="sep"></div>
            <FormKit label="Date Due" type="date" v-model="edited.dateDue" required/>
         </div>
         <FormKit label="Order Title" type="text" v-model="edited.title"/>
         <FormKit label="Special Instructions" type="textarea" rows="5" v-model="edited.specialInstructions"/>
         <FormKit label="Staff Notes" type="textarea" rows="5" v-model="edited.staffNotes"/>
         <FormKit label="Fee" type="number" v-model="edited.fee"/>
         <div class="split">
            <FormKit label="Agency" type="select" v-model="edited.agencyID" :options="agencies" placeholder="Select an agency"/>
            <div class="sep"></div>
            <FormKit label="Customer" type="select" v-model="edited.customerID" :options="customers" placeholder="Select a customer" required/>
         </div>

         <div class="acts">
            <DPGButton label="Cancel" class="p-button-secondary" @click="cancelEdit()"/>
            <FormKit type="submit" label="Save" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'
import { useSystemStore } from '@/stores/system'
import { useCustomersStore } from '@/stores/customers'
import { onMounted, ref, computed } from 'vue'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const ordersStore = useOrdersStore()
const systemStore = useSystemStore()
const customersStore = useCustomersStore()
const edited = ref({
   status: "",
   dateDue: 0,
   title: "",
   specialInstructions: "",
   staffNotes: "",
   fee: null,
   agencyID: "",
   customerID: "",
})
const newOrder = ref(false)

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
   console.log(orderID)
   if (orderID) {
      newOrder.value = false
      await ordersStore.getOrderDetails(orderID)
   } else {
      ordersStore.clearDetails()
      newOrder.value = true
   }

   edited.value.status = ordersStore.detail.status
   edited.value.dateDue = dayjs(ordersStore.detail.dateDue).format("YYYY-MM-DD")
   edited.value.title = ordersStore.detail.title
   edited.value.specialInstructions = ordersStore.detail.specialInstructions
   edited.value.staffNotes = ordersStore.detail.staffNotes
   edited.value.fee = ordersStore.detail.fee
   if (ordersStore.detail.agency) {
      edited.value.agencyID = ordersStore.detail.agency.id
   } else {
      edited.value.agencyID = ""
   }
   if (ordersStore.detail.customer) {
      edited.value.customerID = ordersStore.detail.customer.id
   } else {
      edited.value.customerID = ""
   }
})

function cancelEdit() {
   router.push(`/orders/${route.params.id}`)
}
async function submitChanges() {
   if ( newOrder.value == true) {
      console.log("CREATE ORDER")
      await ordersStore.createOrder( edited.value )
   } else {
      console.log("UPDATE ORDER")
      await ordersStore.submitEdit( edited.value )
   }
   if (systemStore.showError == false) {
      router.push(`/orders/${ordersStore.detail.id}`)
   }
}

</script>

<style lang="scss" scoped>
.edit-form {
   width: 50%;
   margin: 0 auto;

   .split {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      :deep(.formkit-outer) {
         flex-grow: 0.6;
      }
      .sep {
         display: inline-block;
         width: 20px;
      }
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 25px 0;
   button {
      margin-right: 10px;
   }
}
</style>