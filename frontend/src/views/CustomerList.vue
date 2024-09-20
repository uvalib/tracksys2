<template>
   <h2>Customers</h2>
   <div class="customers">
      <DataTable :value="customersStore.customers" ref="customerTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="false" :paginator="true"  v-model:filters="filter"
         :globalFilterFields="['lastName','email']"
         sortField="lastName" :sortOrder="1"
         :rows="30" :totalRecords="customersStore.total"
         v-model:expandedRows="expandedRows"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[30,50,100]" paginatorPosition="top"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <template #paginatorstart  v-if="(userStore.isAdmin || userStore.isSupervisor)" >
            <DPGButton label="Add Customer" severity="secondary" @click="addCustomer()"/>
         </template>
         <template #paginatorend>
            <IconField iconPosition="left">
               <InputIcon class="pi pi-search" />
               <InputText v-model="filter['global'].value" placeholder="Search Customers" />
            </IconField>
         </template>
         <Column :expander="true" headerStyle="width: 3rem" />
         <Column field="lastName" header="Last Name" :sortable="true"/>
         <Column field="firstName" header="First Name"/>
         <Column field="email" header="Email" :sortable="true"/>
         <Column field="academicStatus.name" header="Academic Status"/>
         <Column header="" class="row-acts">
            <template #body="slotProps">
               <DPGButton label="Edit" class="edit-btn" severity="secondary"  @click="edit(slotProps.data)" />
            </template>
         </Column>
         <template #expansion="slotProps">
            <div class="address-data">
               <div v-for="(address,idx) in slotProps.data.addresses" class="address" :key="`customer${slotProps.data}-addr${idx}`">
                  <div v-if="address.addressType=='primary'" class="addr-type">Primary Address</div>
                  <div v-else class="addr-type">Billing Address</div>
                  <div>{{slotProps.data.firstName}} {{slotProps.data.lastName}}</div>
                  <div>{{formattedAddress(address)}}</div>
                   <div>{{address.phone}}</div>
               </div>
            </div>
         </template>
      </DataTable>
       <Dialog v-model:visible="showEdit" :style="{width: '500px'}" header="Customer Details" :modal="true" position="top" :closable="false">
         <FormKit type="form" id="customer-detail" :actions="false" @submit="submitChanges">
            <Tabs value="customer" :lazy="true">
               <TabList>
                  <Tab value="customer">Customer</Tab>
                  <Tab value="address1">Primary Address</Tab>
                  <Tab value="address2">Billing Address</Tab>
               </TabList>
               <TabPanels>
                  <TabPanel value="customer">
                     <FormKit label="Last Name" type="text" name="lastName" v-model="customerDetails.lastName" validation="required" autofocus />
                     <FormKit label="First Name" type="text"name="firstName"  v-model="customerDetails.firstName" validation="required" />
                     <FormKit label="Email" type="email" name="email" v-model="customerDetails.email" validation="required" />
                     <FormKit label="Academic Status" type="select" name="academicStatus" v-model="customerDetails.academicStatus" :options="academicStatuses" required/>
                  </TabPanel>
                  <TabPanel value="address1">
                     <template v-if="customerDetails.addresses.length == 0">
                        <p>No primary address is defined for this customer.</p>
                        <DPGButton label="Add Primary Address" @click="addAddress('primary')"/>
                     </template>
                     <template v-else>
                        <FormKit label="Addresss 1" type="text" name="address1" v-model="customerDetails.addresses[0].address1" validation="required" />
                        <FormKit label="Addresss 2" type="text" name="address2" v-model="customerDetails.addresses[0].address2" />
                        <FormKit label="City" type="text" name="city" v-model="customerDetails.addresses[0].city" validation="required" />
                        <div class="two-col">
                           <FormKit label="State" type="text" name="state" v-model="customerDetails.addresses[0].state" outer-class="state"/>
                           <FormKit label="Zip" type="text" name="zip" v-model="customerDetails.addresses[0].zip" />
                        </div>
                        <div class="two-col">
                           <FormKit label="Country" type="text" name="country" v-model="customerDetails.addresses[0].country" validation="required" outer-class="state"/>
                           <FormKit label="Phone" type="text" name="phone" v-model="customerDetails.addresses[0].phone" />
                        </div>
                     </template>
                  </TabPanel>
                  <TabPanel value="address2">
                     <template v-if="customerDetails.addresses.length == 0">
                        <p>No primary nor billing address is defined for this customer. Please add a primary address.</p>
                     </template>
                     <template v-else-if="customerDetails.addresses.length == 1">
                        <p>No billing address is defined for this customer.</p>
                        <DPGButton label="Add Billing Address" @click="addAddress('billable_address')"/>
                     </template>
                     <template v-else>
                        <FormKit label="Addresss 1" type="text" v-model="customerDetails.addresses[1].address1" validation="required" />
                        <FormKit label="Addresss 2" type="text" v-model="customerDetails.addresses[1].address2" />
                        <FormKit label="City" type="text" v-model="customerDetails.addresses[1].city" validation="required" />
                        <div class="two-col">
                           <FormKit label="State" type="text" v-model="customerDetails.addresses[1].state" outer-class="state"/>
                           <FormKit label="Zip" type="text" v-model="customerDetails.addresses[1].zip" />
                        </div>
                        <div class="two-col">
                           <FormKit label="Country" type="text" v-model="customerDetails.addresses[1].country" validation="required" outer-class="state"/>
                           <FormKit label="Phone" type="text" v-model="customerDetails.addresses[1].phone" />
                        </div>
                     </template>
                  </TabPanel>
               </TabPanels>
            </Tabs>
            <div class="form-controls">
               <FormKit type="button" label="Cancel" wrapper-class="cancel-button" @click="showEdit = false" />
               <FormKit type="submit" label="Save" wrapper-class="submit-button" />
            </div>
         </FormKit>
      </Dialog>
   </div>
</template>

<script setup>
import { ref, computed, onBeforeMount } from 'vue'
import { useCustomersStore } from '@/stores/customers'
import { useSystemStore } from '@/stores/system'
import { useUserStore } from '../stores/user'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import { FilterMatchMode } from '@primevue/core/api'
import { usePinnable } from '@/composables/pin'

usePinnable("p-datatable-paginator-top")

const customersStore = useCustomersStore()
const systemStore = useSystemStore()
const userStore = useUserStore()

const filter = ref( {'global': {value: null, matchMode: FilterMatchMode.STARTS_WITH}})
const expandedRows = ref([])
const showEdit = ref(false)
const customerDetails = ref({
   lastName: "",
   firstName: "",
   academicStatusID: 0,
   academicStatus: {id: 0},
   email: "",
   addresses: [],
   id: 0
})

const academicStatuses = computed(() => {
   let out = []
   systemStore.academicStatuses.forEach( i => {
      out.push( {label: i.name, value: {id: i.id, name: i.name}} )
   })
   return out
})

onBeforeMount(() => {
   customersStore.getCustomers()
   document.title = `Customers`
})

const submitChanges = (() => {
   customersStore.addOrUpdateCustomer(customerDetails.value)
   showEdit.value = false
})

const addAddress = (( addrType ) => {
   let newAddr = {address1: "", address2: "", city: "", state: "", zip: "", country: "", phone: "", addressType: addrType}
   customerDetails.value.addresses.push(newAddr)
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
   customerDetails.value = {
      lastName: "",
      firstName: "",
      academicStatusID: 0,
      academicStatus: {id: 0},
      email: "",
      addresses: [],
      id: 0
   }
   showEdit.value = true
})

const edit = ((data) => {
   customerDetails.value = {...data} // clone the data so edits dont change the store
   showEdit.value = true
})
</script>

<style scoped lang="scss">
.customers {
   min-height: 600px;
   text-align: left;
   padding: 0 25px 25px 25px;
   :deep(td.row-acts) {
      text-align: center;
      .edit-btn {
         font-size: 0.85em;
         padding: 2px 12px;
      }
   }
   .address-data {
      padding: 10px 0 0px 3rem;
      .address {
         margin-bottom: 20px;
         div {
            margin: 5px 0;
         }
         .addr-type {
            font-weight: bold;
            margin: 5px 0;
         }
      }
   }
}

#customer-detail {
   .two-col {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      :deep(.formkit-outer) {
         flex-grow: 1;
      }
      :deep(.formkit-outer.state) {
         margin-right: 15px;
      }
   }
   .form-controls {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      gap: 10px;
      margin-top: 5px;
      text-align: right;
      padding: 10px 0;
   }
}
</style>