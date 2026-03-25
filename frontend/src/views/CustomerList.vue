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
         <form @submit="submitChanges" id="customer-detail">
            <Tabs value="customer" :lazy="false">
               <TabList>
                  <Tab value="customer">Customer</Tab>
                  <Tab value="address1">Primary Address</Tab>
                  <Tab value="address2">Billing Address</Tab>
               </TabList>
               <TabPanels>
                  <TabPanel value="customer">
                     <FormField id="lname" label="Last Name" :error="errors.lastName" :required="true">
                        <InputText id="lname" v-model="lastName" type="text" autofocus/>   
                     </FormField>
                     <FormField id="fname" label="First Name" :error="errors.firstName" :required="true">
                        <InputText id="fname" v-model="firstName" type="text"/>   
                     </FormField>
                     <FormField id="email" label="Email" :error="errors.email" :required="true">
                        <InputText id="email" v-model="email" type="text"/>   
                     </FormField>
                     <FormField id="astatus" label="Academic Status" :error="errors.academicStatusID" :required="true">
                        <Select id="astatus" v-model="academicStatusID"  :options="academicStatuses" optionLabel="label" optionValue="id" placeholder="Select a status" />   
                     </FormField>
                  </TabPanel>
                  <TabPanel value="address1">
                     <template v-if="addresses.length == 0">
                        <p>No primary address is defined for this customer.</p>
                        <DPGButton label="Add Primary Address" @click="addAddress('primary')"/>
                     </template>
                     <template v-else>
                        <FormField id="primary1" label="Addresss 1"  :error="addressError(0, 'address1')" :required="true"  >
                           <InputText id="primary1" v-model="pAddress1" type="text"/>
                        </FormField>
                        <FormField id="primary2" label="Addresss 2">
                           <InputText id="primary2" v-model="pAddress2" type="text"/>   
                        </FormField>
                        <FormField id="primarycity" label="City" :error="addressError(0, 'city')" :required="true"  >
                           <InputText id="primarycity" v-model="pCity" type="text"/>   
                        </FormField>
                        <div class="two-col">
                           <FormField id="primarystate" label="State">
                              <InputText id="primarystate" v-model="pState" type="text"/>   
                           </FormField>
                           <FormField id="primaryzip" label="Zip">
                              <InputText id="primaryzip" v-model="pZip" type="text"/>   
                           </FormField>
                        </div>
                        <div class="two-col">
                           <FormField id="primarycountry" label="Country" :error="addressError(0, 'country')" :required="true"  >
                              <InputText id="primarycountry" v-model="pCountry" type="text"/>   
                           </FormField>
                           <FormField id="primaryphone" label="Phone">
                              <InputText id="primaryphone" v-model="pPhone" type="text"/>   
                           </FormField>
                        </div> 
                     </template>
                  </TabPanel>
                  <TabPanel value="address2">
                     <template v-if="addresses.length == 0">
                        <p>No primary nor billing address is defined for this customer. Please add a primary address.</p>
                     </template>
                     <template v-else-if="addresses.length == 1">
                        <p>No billing address is defined for this customer.</p>
                        <DPGButton label="Add Billing Address" @click="addAddress('billable_address')"/>
                     </template>
                     <template v-else>
                        <FormField id="biz1" label="Addresss 1"  :error="addressError(1, 'address1')" :required="true"  >
                           <InputText id="biz1" v-model="bAddress1" type="text"/>
                        </FormField>
                        <FormField id="biz2" label="Addresss 2">
                           <InputText id="biz2" v-model="bAddress2" type="text"/>   
                        </FormField>
                        <FormField id="bizcity" label="City" :error="addressError(1, 'city')" :required="true"  >
                           <InputText id="bizcity" v-model="bCity" type="text"/>   
                        </FormField>
                        <div class="two-col">
                           <FormField id="bizstate" label="State">
                              <InputText id="bizstate" v-model="bState" type="text"/>   
                           </FormField>
                           <FormField id="bizzip" label="Zip">
                              <InputText id="bizzip" v-model="bZip" type="text"/>   
                           </FormField>
                        </div>
                        <div class="two-col">
                           <FormField id="bizcountry" label="Country" :error="addressError(1, 'country')" :required="true"  >
                              <InputText id="bizcountry" v-model="bCountry" type="text"/>   
                           </FormField>
                           <FormField id="bizphone" label="Phone">
                              <InputText id="bizphone" v-model="bPhone" type="text"/>   
                           </FormField>
                        </div> 
                     </template> 
                  </TabPanel>
               </TabPanels>
            </Tabs>
            <div class="form-controls">
               <DPGButton label="Cancel" severity="secondary" @click="showEdit=false"/>
               <DPGButton label="Save" type="submit" />
            </div>
         </form>
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
import Select from 'primevue/select'
import Dialog from 'primevue/dialog'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import { FilterMatchMode } from '@primevue/core/api'

import { usePinnable } from '@/composables/pin'
usePinnable("p-datatable-paginator-top")

import { useForm } from 'vee-validate'
import * as yup from 'yup'
import FormField from '@/components/FormField.vue'

const {  values, setFieldValue, errors, resetForm, handleSubmit, defineField } = useForm({
  validationSchema: yup.object().shape({
      lastName: yup.string().required('Last name is required'),
      firstName: yup.string().required('First name is required'),
      email: yup.string().email("Email is invalid").required("Email is required"),
      academicStatusID: yup.string().required("Academic status is required"),
      addresses: yup.array().of(
         yup.object({
            address1: yup.string().required('Address 1 is required'),
            city: yup.string().required('City is required'),
            state: yup.string().required('State is required'),
            zip: yup.string().required('Zip is required'),
            country: yup.string().required('Country is required'),
         })   
      ),
   })
})

const [lastName] = defineField('lastName')
const [firstName] = defineField('firstName')
const [email] = defineField('email')
const [academicStatusID] = defineField('academicStatusID')

const [addresses] = defineField('addresses')

const [pAddress1] = defineField('addresses[0].address1')
const [pAddress2] = defineField('addresses[0].address2')
const [pCity] = defineField('addresses[0].city')
const [pState] = defineField('addresses[0].state')
const [pZip] = defineField('addresses[0].zip')
const [pCountry] = defineField('addresses[0].country')
const [pPhone] = defineField('addresses[0].phone')

const [bAddress1] = defineField('addresses[1].address1')
const [bAddress2] = defineField('addresses[1].address2')
const [bCity] = defineField('addresses[1].city')
const [bState] = defineField('addresses[1].state')
const [bZip] = defineField('addresses[1].zip')
const [bCountry] = defineField('addresses[1].country')
const [bPhone] = defineField('addresses[1].phone')


const customersStore = useCustomersStore()
const systemStore = useSystemStore()
const userStore = useUserStore()

const filter = ref( {'global': {value: null, matchMode: FilterMatchMode.STARTS_WITH}})
const expandedRows = ref([])
const showEdit = ref(false)

const addressError = ((index, field) => {
   const addressField = `addresses[${index}].${field}`
   return errors.value[addressField]
})

const academicStatuses = computed(() => {
   let out = []
   systemStore.academicStatuses.forEach( i => {
      out.push( {label: i.name, id: i.id} )
   })
   return out
})

onBeforeMount(() => {
   customersStore.getCustomers()
   document.title = `Customers`
})

const submitChanges = handleSubmit(values => {
   customersStore.addOrUpdateCustomer(values)
   showEdit.value = false
})

const addAddress = (( addrType ) => {
   let newAddr = {address1: "", address2: "", city: "", state: "", zip: "", country: "", phone: "", addressType: addrType}
   console.log(addresses.value)
   let a = addresses.value.splice()
   a.push(newAddr)
   setFieldValue("addresses", a)
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
   resetForm({
      values: {
         lastName: "",
         firstName: "",
         academicStatusID: null,
         email: "",
         addresses: [],
         id: 0
      }
   })
   showEdit.value = true
})

const edit = ((data) => {
   resetForm({ values: data })
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

.form-controls {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   gap: 10px;
}

:deep(.p-tabpanel) {
   display: flex;
   flex-direction: column;
   gap: 15px;

   .two-col {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      gap: 15px;
   }
}
   
</style>