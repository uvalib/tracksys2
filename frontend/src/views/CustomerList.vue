<template>
   <h2>Customers</h2>
   <DPGButton label="Add" class="p-button-secondary create" @click="addCustomer()"/>
   <div class="customers">
      <div class="filter-controls">
         <span>
            <span class="p-input-icon-right">
               <i class="pi pi-search" />
               <InputText v-model="filter['global'].value" placeholder="Customer Search"/>
            </span>
            <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()"/>
         </span>
      </div>
      <DataTable :value="customersStore.customers" ref="customerTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="false" :paginator="true"  v-model:filters="filter"
         :globalFilterFields="['lastName','email']"
         sortField="lastName" :sortOrder="1"
         :rows="10" :totalRecords="customersStore.total"
         v-model:expandedRows="expandedRows"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[10,30,100]"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <Column :expander="true" headerStyle="width: 3rem" />
         <Column field="lastName" header="Last Name" :sortable="true"/>
         <Column field="firstName" header="First Name"/>
         <Column field="email" header="Email" :sortable="true"/>
         <Column field="academicStatus.name" header="Acedemic Status"/>
         <Column header="" class="row-acts">
            <template #body="slotProps">
               <DPGButton label="Edit" class="p-button-text"  @click="edit(slotProps.data)" />
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
       <Dialog v-model:visible="showEdit" :style="{width: '500px'}" header="Customer Details" :modal="true" position="top">
         <FormKit type="form" id="customer-detail" :actions="false" @submit="submitChanges">
            <TabView>
               <TabPanel header="Customer">
                  <FormKit label="Last Name" type="text" v-model="customerDetails.lastName" validation="required" autofocus />
                  <FormKit label="First Name" type="text" v-model="customerDetails.firstName" validation="required" />
                  <FormKit label="Email" type="email" v-model="customerDetails.email" validation="required" />
                  <FormKit label="Acedemic Status" type="select" v-model="customerDetails.academicStatus" :options="academicStatuses" required/>
               </TabPanel>
               <TabPanel header="Primary Address">
                  <div v-if="customerDetails.addresses.length == 0">
                     <p>No primary address is defined for this customer.</p>
                     <DPGButton label="Add Primary Address" @click="addAddress('primary')"/>
                  </div>
                  <template v-else>
                     <FormKit label="Addresss 1" type="text" v-model="customerDetails.addresses[0].address1" validation="required" />
                     <FormKit label="Addresss 2" type="text" v-model="customerDetails.addresses[0].address2" />
                     <FormKit label="City" type="text" v-model="customerDetails.addresses[0].city" validation="required" />
                     <div class="two-col">
                        <FormKit label="State" type="text" v-model="customerDetails.addresses[0].state" outer-class="state"/>
                        <FormKit label="Zip" type="text" v-model="customerDetails.addresses[0].zip" />
                     </div>
                     <div class="two-col">
                        <FormKit label="Country" type="text" v-model="customerDetails.addresses[0].country" validation="required" outer-class="state"/>
                        <FormKit label="Phone" type="text" v-model="customerDetails.addresses[0].phone" />
                     </div>
                  </template>
               </TabPanel>
               <TabPanel header="Billing Address">
                  <div v-if="customerDetails.addresses.length == 0">
                     <p>No primary nor billing address is defined for this customer. Please add a primary address.</p>
                  </div>
                   <div v-else-if="customerDetails.addresses.length == 1">
                     <p>No billing address is defined for this customer.</p>
                     <DPGButton label="Add Billing Address" @click="addAddress('billable_address')"/>
                  </div>
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
            </TabView>
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
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import { FilterMatchMode } from 'primevue/api'

const customersStore = useCustomersStore()
const systemStore = useSystemStore()

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
   document.title = `Customers`
})

function submitChanges() {
   customersStore.addOrUpdateCustomer(customerDetails.value)
   showEdit.value = false
}

function addAddress( addrType ) {
   let newAddr = {address1: "", address2: "", city: "", state: "", zip: "", country: "", phone: "", addressType: addrType}
   customerDetails.value.addresses.push(newAddr)
}

function formattedAddress(data) {
   let out = data.address1
   if (data.address2) {
      out += ` ${data.address2}`
   }
   out += `, ${data.city} ${data.state} ${data.zip}, ${data.country}`
   return  out
}

function addCustomer() {
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
}

function edit(data) {
   customerDetails.value = {...data} // clone the data so edits dont change the store
   showEdit.value = true
}

function clearSearch() {
   filter.value = ""
}
</script>

<style scoped lang="scss">
button.p-button-secondary.create {
   position: absolute;
   right:15px;
   top: 15px;
}
.customers {
   min-height: 600px;
   text-align: left;
   padding: 0 25px 25px 25px;
   .filter-controls {
      padding: 10px 0;
      display: flex;
      flex-flow: row wrap;
      justify-content: flex-end;
      button.p-button-secondary.p-button {
         margin-left: 5px;
      }
   }
   :deep(.row-acts) {
      text-align: center;
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
      margin-top: 5px;
      text-align: right;
      padding: 10px 0;
      :deep(.cancel-button button) {
         @include base-button();
         width: auto;
         margin-right: 10px;
      }
      :deep(.submit-button button) {
         @include primary-button();
         width: auto;
      }
   }
}
</style>