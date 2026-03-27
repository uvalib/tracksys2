<template>
   <DPGButton @click="show" severity="secondary" label="Create Digitization Project" :disabled="createDisabled"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Create Digitization Project" :style="{width: '400px'}" :closable="false">
      <form @submit="createProject">
         <FormField id="workdlow" label="Workflow" :error="errors.workflowID" :required="true">
            <Select id="workdlow" v-model="workflowID"  :options="workflows" optionLabel="label" optionValue="value" placeholder="Select a workflow" />   
         </FormField>
         <FormField  v-if="workflowID==6" id="containertype" label="Container Tyoe" :error="errors.containerTypeID" :required="true">
            <Select id="containertype" v-model="containerTypeID"  :options="containerTypes" optionLabel="label" optionValue="value" placeholder="Select a container type" />   
         </FormField>
         <FormField id="category" label="Category" :error="errors.categoryID" :required="true">
            <Select id="category" v-model="categoryID"  :options="categories" optionLabel="label" optionValue="value" placeholder="Select a category" />   
         </FormField>
         <FormField id="condition" label="Condition" :error="errors.condition" :required="true">
            <Select id="condition" v-model="condition"  :options="conditions" optionLabel="label" optionValue="value" placeholder="Select a condition" />   
         </FormField>
         <FormField id="notes" label="Notes">
            <Textarea id="notes" v-model="notes" rows="4"/>   
         </FormField>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" severity="secondary"/>
            <DPGButton label="Create Project" type="submit"/>
         </div>
      </form>
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import Dialog from 'primevue/dialog'
import { useUnitsStore } from '@/stores/units'
import { useSystemStore } from '@/stores/system'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import axios from 'axios'
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import FormField from '@/components/FormField.vue'

const schema = yup.object().shape({
      workflowID: yup.number().min(1, 'Workflow is required'),
      categoryID: yup.number().min(1, 'Category is required'),
      condition: yup.bool().required("Condition is required"),
      containerTypeID: yup.number().when('workflowID', {
         is: (value) => value == 6,
         then: (schema) => schema.min(1,"Container type is required"),
      }),
   })
const { errors, resetForm, handleSubmit, defineField } = useForm({validationSchema: schema})

const [workflowID] = defineField('workflowID')
const [containerTypeID] = defineField('containerTypeID')
const [categoryID] = defineField('categoryID')
const [condition] = defineField('condition')
const [notes] = defineField('notes')

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

const createProject = handleSubmit( async (values) => {
   axios.post(`${systemStore.projectsURL}/api/projects/create`, values).then(response => {
      unitsStore.detail.projectID = parseInt(response.data, 10)
      systemStore.toastMessage("Project Created", "A new project has been created for this unit")
      hide()
   }).catch( err => {
      systemStore.setError(err)
   })
})

const hide = (() => {
   isOpen.value=false
})

const show = (() => {
   let vals = {
      unitID: unitsStore.detail.id,
      orderID: unitsStore.detail.orderID,
      dateDue: unitsStore.detail.order.dateDue,
      title: unitsStore.detail.metadata.title,
      callNumber: unitsStore.detail.metadata.callNumber,
      agencyID: 0,
      customerID: unitsStore.detail.order.customer.id,
      workflowID: 1,
      categoryID: 0,
      containerTypeID: 0,
      condition: null,
      notes: "",
   }
   if ( unitsStore.detail.order.agency ) {
       vals.agencyID = unitsStore.detail.order.agency.id
   }
   resetForm({ values: vals })
   
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
form {
   display: flex;
   flex-direction: column;
   gap: 15px;
   .acts {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      gap: 10px;
   }
}
</style>
