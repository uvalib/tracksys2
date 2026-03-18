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
         <Form v-slot="$form" :initialValues="customerDetails" :resolver @submit="submitChanges" id="customer-detail" :validateOnBlur="true">
            <Tabs value="customer">
               <TabList>
                  <Tab value="customer">Customer</Tab>
                  <Tab value="address1">Primary Address</Tab>
                  <Tab value="address2">Billing Address</Tab>
               </TabList>
               <TabPanels>
                  <TabPanel value="customer">
                     <FormField id="lname" label="Last Name" :error="$form.lastName?.invalid ? $form.lastName.error.message : ''" :required="true">
                        <InputText id="lname" name="lastName" type="text" autofocus/>   
                     </FormField>
                     <FormField id="fname" label="First Name" :error="$form.firstName?.invalid ? $form.firstName.error.message : ''" :required="true">
                        <InputText id="fname" name="firstName" type="text"/>   
                     </FormField>
                     <FormField id="email" label="Email" :error="$form.email?.invalid ? $form.email.error.message : ''" :required="true">
                        <InputText id="email" name="email" type="text"/>   
                     </FormField>
                     <FormField id="astatus" label="Academic Status" :error="$form.academicStatusID?.invalid ? $form.academicStatusID.error.message : ''" :required="true">
                        <Select id="astatus" name="academicStatusID"  :options="academicStatuses" optionLabel="label" optionValue="id" placeholder="Select a status" />   
                     </FormField>
                  </TabPanel>
                  <TabPanel value="address1">
                     <template v-if="customerDetails.addresses.length == 0">
                        <p>No primary address is defined for this customer.</p>
                        <DPGButton label="Add Primary Address" @click="addAddress('primary')"/>
                     </template>
                     <template v-else>
                        <FormField id="primary1" label="Addresss 1" :error="addrErr($form, 0, 'address1' )" :required="true">
                           <InputText id="primary1" name="addresses[0].address1"  v-model="customerDetails.addresses[0].address1" type="text"/>   
                        </FormField>
                        <FormField id="primary2" label="Addresss 2">
                           <InputText id="primary2" name="addresses[0].address2" v-model="customerDetails.addresses[0].address2" type="text"/>   
                        </FormField>
                        <FormField id="primarycity" label="City" :error="addrErr($form, 0, 'city' )" :required="true">
                           <InputText id="primarycity" name="addresses[0].city" v-model="customerDetails.addresses[0].city" type="text"/>   
                        </FormField>
                        <div class="two-col">
                           <FormField id="primarystate" label="State">
                              <InputText id="primarystate" name="addresses[0].state" v-model="customerDetails.addresses[0].state" type="text"/>   
                           </FormField>
                           <FormField id="primaryzip" label="Zip">
                              <InputText id="primaryzip" name="addresses[0].zip" v-model="customerDetails.addresses[0].zip" type="text"/>   
                           </FormField>
                        </div>
                        <div class="two-col">
                           <FormField id="primarycountry" label="Country" :error="addrErr($form, 0, 'country' )" :required="true">
                              <InputText id="primarycountry" name="addresses[0].country" v-model="customerDetails.addresses[0].country" type="text"/>   
                           </FormField>
                           <FormField id="primaryphone" label="Phone">
                              <InputText id="primaryphone" name="addresses[0].phone" v-model="customerDetails.addresses[0].phone" type="text"/>   
                           </FormField>
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
                        <FormField id="biz1" label="Addresss 1" :error="addrErr($form, 1, 'address1' )" :required="true">
                           <InputText id="biz1" name="addresses[1].address1"  v-model="customerDetails.addresses[1].address1" type="text"/>   
                        </FormField>
                        <FormField id="biz2" label="Addresss 2">
                           <InputText id="biz2" name="addresses[1].address2" v-model="customerDetails.addresses[1].address2" type="text"/>   
                        </FormField>
                        <FormField id="bizcity" label="City" :error="addrErr($form, 1, 'city' )" :required="true">
                           <InputText id="bizcity" name="addresses[1].city" v-model="customerDetails.addresses[1].city" type="text"/>   
                        </FormField>
                        <div class="two-col">
                           <FormField id="bizstate" label="State">
                              <InputText id="bizstate" name="addresses[1].state" v-model="customerDetails.addresses[1].state" type="text"/>   
                           </FormField>
                           <FormField id="bizzip" label="Zip">
                              <InputText id="bizzip" name="addresses[1].zip" v-model="customerDetails.addresses[1].zip" type="text"/>   
                           </FormField>
                        </div>
                        <div class="two-col">
                           <FormField id="bizcountry" label="Country" :error="addrErr($form, 1, 'country' )" :required="true">
                              <InputText id="bizcountry" name="addresses[1].country" v-model="customerDetails.addresses[1].country" type="text"/>   
                           </FormField>
                           <FormField id="bizphone" label="Phone">
                              <InputText id="bizphone" name="addresses[1].phone" v-model="customerDetails.addresses[1].phone" type="text"/>   
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
         </Form>
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
import { Form } from '@primevue/forms'
import { yupResolver } from '@primevue/forms/resolvers/yup'
import * as yup from 'yup'
import FormField from '@/components/FormField.vue'
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
   academicStatusID: null,
   email: "",
   addresses: [],
   id: 0
})

const resolver = yupResolver( 
   yup.object().shape({
      lastName: yup.string().required('Last name is required'),
      firstName: yup.string().required('First name is required'),
      email: yup.string().email("Email is invalid").required("Email is required"),
      academicStatusID: yup.string().required("Academic status is required"),
      addresses: yup.array(
         yup.object({
            address1: yup.string().required('Address 1 is required'),
            city: yup.string().required('City is required'),
            state: yup.string().required('State is required'),
            zip: yup.string().required('Zip is required'),
            country: yup.string().required('Country is required'),
         })   
      ),
   })
)

const academicStatuses = computed(() => {
   let out = []
   systemStore.academicStatuses.forEach( i => {
      out.push( {label: i.name, id: i.id} )
   })
   return out
})

const addrErr = ((form, addrIdx, field) => {
   if (form.addresses && form.addresses[addrIdx] && form.addresses[addrIdx][field].invalid) {
      return form.addresses[addrIdx][field].error.message
   }
   return ""
})

onBeforeMount(() => {
   customersStore.getCustomers()
   document.title = `Customers`
})

const submitChanges = async ({ valid, values }) => {
   if (valid ) {
      values.id = customerDetails.value.id
      customersStore.addOrUpdateCustomer(values)
      showEdit.value = false
   }
}

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
      academicStatusID: null,
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
   }
}
   
</style>