<template>
   <UModal :dismissible="false" :close="false" :title="props.title" @enter="modalOpened">
      <template #body>
         <UForm :state="state"  :schema="schema" @submit="submitChanges" class="column w-full" @error="onError">
            <UTabs :items="tabs" variant="link" valueKey="slot">
               <template #customer>
                  <div class="column gap-4">
                     <UFormField name="lastName" label="Last Name" :required="true" class="grow">  
                        <UInput v-model="state.lastName" class="w-full"/>   
                     </UFormField>
                     <UFormField name="firstName" label="First Name" :required="true" class="grow">  
                        <UInput v-model="state.firstName" class="w-full"/>   
                     </UFormField>
                     <UFormField name="email" label="Email" :required="true" class="grow">  
                        <UInput v-model="state.email" class="w-full"/>   
                     </UFormField>
                     <UFormField name="academicStatusID" label="Academic Status" :required="true" class="grow">  
                        <USelect v-model="state.academicStatusID" :items="academicStatuses" class="w-full" placeholder="Select a role..."/>   
                     </UFormField>
                  </div>
               </template>
               <template v-for="info in addresses" #[info.name]>
                  <div v-if="info.name == 'primary' && state.addresses.length == 0" class="pt-2 column  items-center">
                     <div>No primary address is defined for this customer.</div>
                     <UButton label="Add Primary Address" @click="addAddress('primary')"/>
                  </div>
                  <div v-else-if="info.name == 'billable_address' && state.addresses.length < 2" class="pt-2 column  items-center">
                     <template v-if="state.addresses.length == 0">
                        <div>No primary nor billing address is defined for this customer.<br/>Please add a primary address.</div>
                     </template>
                     <template v-else>
                        <div>No billing address is defined for this customer.</div>
                        <UButton label="Add Billing Address" @click="addAddress('billable_address')"/>
                     </template>
                  </div>
                  <div v-else class="column gap-4">
                     <UFormField :name="`addresses.${info.index}.address1`" label="Address 1" :required="true" class="grow">  
                        <UInput v-model="state.addresses[info.index].address1" class="w-full"/>   
                     </UFormField>
                     <UFormField :name="`addresses.${info.index}.addresss2`" label="Address 2" class="grow">  
                        <UInput v-model="state.addresses[info.index].addres2" class="w-full"/>   
                     </UFormField>
                     <UFormField :name="`addresses.${info.index}.city`" label="City" class="grow" :required="true">  
                        <UInput v-model="state.addresses[info.index].city" class="w-full"/>   
                     </UFormField>
                     <div class="row-left">
                        <UFormField :name="`addresses.${info.index}.state`" label="State" class="grow">  
                           <UInput v-model="state.addresses[info.index].state" class="w-full"/>   
                        </UFormField>    
                        <UFormField :name="`addresses.${info.index}.zip`" label="Zip" class="grow">  
                           <UInput v-model="state.addresses[info.index].zip" class="w-full"/>   
                        </UFormField>   
                     </div>
                     <div class="row-left">
                        <UFormField :name="`addresses.${info.index}.country`" label="Country" class="grow" :required="true">  
                           <UInput v-model="state.addresses[info.index].country" class="w-full"/>   
                        </UFormField>    
                        <UFormField :name="`addresses.${info.index}.phone`" label="Phone" class="grow">  
                           <UInput v-model="state.addresses[info.index].phone" class="w-full"/>   
                        </UFormField>   
                     </div>
                  </div>
               </template>
            </UTabs>
            <div class="row-right gap-2 pt-4">
               <UButton label="Cancel" color="secondary" @click="emit('close', false)"/>
               <UButton label="Save" type="submit" /> 
            </div>
         </UForm>
      </template>
   </UModal>
</template>

<script setup>
import { ref, computed } from 'vue'
import * as yup from 'yup'
import { useCustomersStore } from '@/stores/customers'
import { useSystemStore } from '@/stores/system'

const customerStore = useCustomersStore()
const systemStore = useSystemStore()

const emit = defineEmits( ['close' ])
const props = defineProps({
   mode: {
      type: String,
      required: true,
      validator: (value) => {
         return ['add', 'edit'].includes(value)    
      }
   },
   title: {
      type: String,
      required: true,
   },
   customer: {
      type: Object,
      default: null
   }
})

const onError = ((err) => {
   console.log(err)
})

const schema = yup.object({
   lastName: yup.string().required('Last name is required'),
   firstName: yup.string().required('First name is required'),
   email: yup.string().email("Email is invalid").required("Email is required"),
   academicStatusID: yup.string().required("Academic status is required"),
   addresses: yup.array().of(
      yup.object({
         address1: yup.string().required('Address 1 is required'),
         city: yup.string().required('City is required'),
         country: yup.string().required('Country is required'),
      })   
   ).min(1),
})

const state = ref({
   id: 0,
   lastName: "",
   firstName: "",
   email: "",
   academicStatusID: "",
   addresses: []
})

const addresses = computed(() => {
   return [ {name: "primary", index: 0}, {name: "billable_address", index: 1} ]
})

const tabs = computed(() => {
   let out = [
      { label: `Customer`, slot: 'customer'},
      { label: `Primary Address`, slot: 'primary'},
      { label: `Billing Address`, slot: 'billable_address'}
   ]
   return out
})

const academicStatuses = computed(() => {
   let out = []
   systemStore.academicStatuses.forEach( i => {
      out.push( {label: i.name, value: i.id} )
   })
   return out
})

const modalOpened = (() => {
   state.value = {
      id: 0,
      lastName: "",
      firstName: "",
      email: "",
      academicStatusID: "",
      addresses: []
   }
   if (props.mode == "edit" ) {
      state.value = {
         id: props.customer.id,
         lastName: props.customer.lastName,
         firstName: props.customer.firstName,
         email: props.customer.email,
         academicStatusID: props.customer.academicStatusID,
         addresses: [],
      }
      props.customer.addresses.forEach( a => state.value.addresses.push(a))
   }   
})

const submitChanges = ( async () => {
   await customerStore.addOrUpdateCustomer( state.value )
   emit('close', true)
})

const addAddress = (( addrType ) => {
   let newAddr = {address1: "", address2: "", city: "", state: "", zip: "", country: "", phone: "", addressType: addrType}
   state.value.addresses.push( newAddr )
})

</script>