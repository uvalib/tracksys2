<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" title="Create Digitization Project">
      <UButton @click="show()" color="secondary" label="Create Digitization Project" />
      <template #body>
         <UForm  :schema="schema" :state="state" class="flex flex-col gap-4" @submit="createProject">
            <UFormField name="workflowID" label="Workflow" :required="true">
               <USelect v-model="state.workflowID" :items="workflows" placeholder="Select a workflow..." class="w-full" />   
            </UFormField>
            <UFormField  v-if="state.workflowID==6"  name="containerTypeID" label="Container Type" :required="true">
               <USelect v-model="state.containerTypeID" :items="containerTypes" placeholder="Select a container type..." class="w-full" />   
            </UFormField>
            <UFormField name="categoryID" label="Category" :required="true">
               <USelect v-model="state.categoryID" :items="categories" placeholder="Select a category..." class="w-full" />   
            </UFormField>
            <UFormField name="condition" label="Condition" :required="true">
               <USelect v-model="state.condition" :items="conditions" placeholder="Select a condition..." class="w-full" />   
            </UFormField>
            <UFormField name="notes" label="Notes">
               <UTextarea v-model="state.notes" class="w-full" />   
            </UFormField>
            <div class="row-right gap-2">
               <UButton @click="isOpen = false" label="Cancel" color="secondary"/>
               <UButton label="Create Project" type="submit" :disabled="createDisabled"/>  
            </div>
         </UForm>
      </template>
   </UModal>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useSystemStore } from '@/stores/system'
import { toCalendarDate, parseAbsolute } from '@internationalized/date'
import axios from 'axios'
import * as yup from 'yup'

const schema = yup.object({
   workflowID: yup.number().min(1, 'Workflow is required').required("Workflow is required"),
   categoryID: yup.number().min(1, 'Category is required').required("Category is required"),
   condition: yup.bool().required("Condition is required"),
   containerTypeID: yup.number().when('workflowID', {
      is: (value) => value == 6,
      then: (schema) => schema.min(1,"Container type is required").required("Container type is required"),
   }),
})

const state = ref({
   unitID: null,
   orderID: null,
   dateDue: null,
   title: "",
   callNumber: "",
   agencyID: null,
   customerID: null,
   workflowID: null,
   categoryID: null,
   containerTypeID: null,
   condition: null,
   notes: "",
})

const unitsStore = useUnitsStore()
const systemStore = useSystemStore()

const isOpen = ref(false)
const workflows = ref([])
const categories = ref([])

const createDisabled = computed(() => {
   let approved = (unitsStore.detail.status == 'approved' && unitsStore.detail.order.status == 'approved')
   return !approved
})
const containerTypes = computed( () => {
   let out = []
   systemStore.containerTypes.forEach( w => {
      out.push({label: w.name, value: w.id})
   })
   return out
})
const conditions = computed( () => {
   let out = [
      {label: "Good", value: 0},
      {label: "Bad", value: 1},
   ]
   return out
})

const createProject = ( async () => {
   axios.post(`${systemStore.projectsURL}/api/projects/create`, values).then(response => {
      unitsStore.detail.projectID = parseInt(response.data, 10)
      systemStore.toastMessage("Project Created", "A new project has been created for this unit")
      hide()
   }).catch( err => {
      systemStore.setError(err)
   })
})

const show = (() => {
   state.value = {
      unitID: unitsStore.detail.id,
      orderID: unitsStore.detail.orderID,
      title: unitsStore.detail.metadata.title,
      callNumber: unitsStore.detail.metadata.callNumber,
      agencyID: null,
      customerID: unitsStore.detail.order.customer.id,
      workflowID: null,
      categoryID: null,
      containerTypeID: null,
      condition: null,
      notes: "",
   }
   if ( unitsStore.detail.order.dateDue ) {
      state.value.dateDue = toCalendarDate( parseAbsolute(unitsStore.detail.order.dateDue) )
   }
   if ( unitsStore.detail.order.agency ) {
       state.value.agencyID = unitsStore.detail.order.agency.id
   }
   
   // load constants from dpg-imaging
   axios.get(`${systemStore.projectsURL}/constants`).then(response => {
      workflows.value = response.data.workflows
      categories.value = response.data.categories
      console.log(response.data)
      isOpen.value = true
   }).catch( err => {
      systemStore.setError(err)
   })
})
</script>

<style lang="scss" scoped>
</style>
