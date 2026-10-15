<template>
   <UModal :dismissible="false" :close="false" :title="props.title" @enter="modalOpened">
      <template #body>
         <UForm :state="state"  :schema="schema" @submit="submitChanges" class="column w-full">
            <UFormField name="lastName" label="Last Name" :required="true" class="grow">  
               <UInput v-model="state.lastName" class="w-full"/>   
            </UFormField>
            <UFormField name="firstName" label="First Name" :required="true" class="grow">  
               <UInput v-model="state.firstName" class="w-full"/>   
            </UFormField>
            <UFormField name="computingID" label="UVA Computing ID" :required="true" class="grow">  
               <UInput v-model="state.computingID" class="w-full"/>   
            </UFormField>
            <UFormField name="email" label="Email" :required="true" class="grow">  
               <UInput v-model="state.email" class="w-full"/>   
            </UFormField>
            <UFormField name="roleID" label="Role" :required="true" class="grow">  
               <USelect v-model="state.roleID" :items="roles" class="w-full" placeholder="Select a role..."/>   
            </UFormField>
            <UFormField name="active" label="Active" class="grow">  
               <USelect v-model="state.active" :items="yesNo" class="w-full"/>   
            </UFormField>
            <div class="notes">
               <b>IMPORTANT:</b>
               <span>All new staff must be added to a group named lb-digiserv. This can be done here: </span>
               <a href="https://mygroups.virginia.edu/groups/" target="_blank">MyGroups</a>
            </div>
            <div class="row-right gap-2!">
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
import { useStaffStore } from '@/stores/staff'

const staffStore = useStaffStore()

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
   staff: {
      type: Object,
      default: null
   }
})

const schema = yup.object({
   lastName: yup.string().required('Last name is required'),
   firstName: yup.string().required('First name is required'),
   computingID: yup.string().required('ComputingID is required'),
   email: yup.string().email("Email is invalid").required("Email is required"),
   roleID: yup.number().min(1, "Role is required"),
})

const state = ref({
   id: 0,
   lastName: "",
   firstName: "",
   email: "",
   computingID: "",
   roleID: null,
   role: "",
   active: false
})

const roles = computed( () => {
   return [
      {label:"Admin", value: "admin"}, {label:"Supervisor", value: 1}, 
      {label:"Student", value: 2}, {label:"Viewer", value: 3}
   ]
})
const yesNo = computed(() => {
   return [ {label: "No", value: false}, {label: "Yes", value: true} ]
})

const modalOpened = (() => {
   if (props.mode == "edit" ) {
      state.value = {
         id: props.staff.id,
         lastName: props.staff.lastName,
         firstName: props.staff.firstName,
         email: props.staff.email,
         computingID: props.staff.computingID,
         roleID: props.staff.roleID,
         role: "",
         active: props.staff.active,
      }
   }
})

const submitChanges = (async () => {
   const roleMatch = roles.value.find( r => r.value == state.value.roleID)
   if (roleMatch) {
      state.value.role = roleMatch.label
   }
   await staffStore.addOrUpdateStaff(state.value)
   emit('close', true)
}) 
</script>

<style scoped lang="scss">
.notes {
   font-size: 0.9em;
   margin: 15px 0 10px 0;
   border:  1px solid var(--uvalib-teal);
   padding: 10px 20px;
   border-radius: 5px;
   background-color: var(--uvalib-teal-lightest);
   b {
      font-weight: bold;
      padding-right: 10px;
   }
   a {
      color: var(--uvalib-blue-alt-dark);
   }
}
</style>