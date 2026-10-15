<template>
   <h2>Customers</h2>
   <div class="row-between p-2 sticky  bg-white border-b border-brand-grey-100" :style="{top: system.headerHeight}">
      <UButton v-if="(userStore.isAdmin || userStore.isSupervisor)" label="Add Customer" color="secondary" @click="addCustomer" />
      <UInput v-model="filter" placeholder="Search customers..." />
   </div>
   <UTable :data="customersStore.customers" :columns="columns" v-model:global-filter="filter" v-model:expanded="expanded"
      sticky virtualize class="h-[calc(100vh-185px)]">
      <template #expand-cell  ="{ row }">
         <UButton color='neutral' variant='ghost' icon='i-lucide-chevron-right' aria-label='Expand' 
            :ui="{ leadingIcon: ['transition-transform',row.getIsExpanded() ? 'rotate-90' : '']}"
            @click="row.toggleExpanded()"/>
      </template>
      <template #expanded="{ row }">
         <div class="address-data">
            <div v-for="(address,idx) in row.original.addresses" class="address" :key="`customer${row.original}-addr${idx}`">
               <div v-if="address.addressType=='primary'" class="addr-type">Primary Address</div>
               <div v-else class="addr-type">Billing Address</div>
               <div>{{row.original.firstName}} {{row.original.lastName}}</div>
               <div>{{formattedAddress(address)}}</div>
                  <div>{{address.phone}}</div>
            </div>
         </div>   
      </template>
      <template #acts-cell="{ row }">
         <UButton label="Edit" color="secondary" size="xs" @click="edit(row.original)" />
      </template>  
   </UTable>
</template>

<script setup>
import { ref, onBeforeMount } from 'vue'
import { useCustomersStore } from '@/stores/customers'
import { useUserStore } from '../stores/user'
import { useSystemStore } from '../stores/system'
import CustomerEditModal from '@/components/CustomerEditModal.vue'

const customersStore = useCustomersStore()
const userStore = useUserStore()
const system = useSystemStore()
const customerModal = useOverlay().create(CustomerEditModal)

const filter = ref("")
const expanded = ref({})

const columns = [
   { accessorKey: 'expand', header: ""},
   { accessorKey: 'lastName', header: "Last Name"},
   { accessorKey: 'firstName', header: "First Name"},
   { accessorKey: 'email', header: "Email"},
   { accessorKey: 'academicStatus.name', id: "academicStatus", header: "Academic Status"},
   { accessorKey: 'acts', header: ""}
]

onBeforeMount(() => {
   customersStore.getCustomers()
   document.title = `Customers`
})

const formattedAddress = ((data) => {
   let out = data.address1
   if (data.address2) {
      out += ` ${data.address2}`
   }
   out += `, ${data.city} ${data.state} ${data.zip}, ${data.country}`
   return  out
})

const addCustomer = (() => {
   customerModal.open({mode: "add", title: "Add Customer"})
})

const edit = ((data) => {
   customerModal.open({mode: "edit", title: "Edit Customer", customer: {...data}})
})
</script>

<style scoped lang="scss">
.address-data {
   padding: 10px 0 0px 3rem;
   .address {
      margin-bottom: 20px;
      .addr-type {
         font-weight: bold;
      }
   }
}
</style>