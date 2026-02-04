<template>
   <DPGButton @click="show" severity="secondary" label="Create Digitization Project" :disabled="createDisabled"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Create Digitization Project" :style="{width: '400px'}" :closable="false">
      <FormKit type="form" id="create-project" :actions="false" @submit="createProject">
         <FormKit label="Workflow" type="select" v-model="project.workflowID" :options="workflows" required outer-class="first" />
         <FormKit v-if="project.workflowID==6" label="Container Type" type="select" v-model="project.containerTypeID" :options="containerTypes" required placeholder="Select a container type"/>
         <FormKit label="Category" type="select" v-model="project.categoryID" :options="categories" required placeholder="Select a category"/>
         <FormKit label="Condition" type="select" v-model="project.condition" :options="conditions" required/>
         <FormKit label="Notes" type="textarea" rows="4" v-model="project.notes"/>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" severity="secondary"/>
            <FormKit type="submit" label="Create Project" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import Dialog from 'primevue/dialog'
import { useUnitsStore } from '@/stores/units'
import { useSystemStore } from '@/stores/system'
import axios from 'axios'

const unitsStore = useUnitsStore()
const systemStore = useSystemStore()

const isOpen = ref(false)
const workflows = ref([])
const categories = ref([])
const project = ref({
   unitID: 0,
   orderID: 0,
   title: "",
   callNumber: "",
   customerID: 0,
   agencyID: 0,
   workflowID: 1,
   containerTypeID: 0,
   categoryID: null,
   condition: 0,
   notes: ""
})

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

async function createProject() {
   // FIXME cache some TS specific data in the project; order_id, customer_id, agency_id, call_number (needed for search)
   // possibly cache some for display too: title, special instructions, intended use
   axios.post(`${systemStore.projectsURL}/api/projects/create`, project.value).then(response => {
      unitsStore.detail.projectID = parseInt(response.data, 10)
      systemStore.toastMessage("Project Created", "A new project has been created for this unit")
      hide()
   }).catch( err => {
      systemStore.setError(err)
   })
}
function hide() {
   isOpen.value=false
}
function show() {
   project.value.unitID = unitsStore.detail.id
   project.value.orderID = unitsStore.detail.orderID
   project.value.title = unitsStore.detail.metadata.title
   project.value.callNumber = unitsStore.detail.metadata.callNumber
   project.value.agencyID = 0
   if ( unitsStore.detail.order.agency ) {
       project.value.agencyID = unitsStore.detail.order.agency.id
   }
   project.value.customerID = unitsStore.detail.order.customer.id
   project.value.workflowID = 1
   project.value.categoryID = 0
   project.value.containerTypeID = 0
   project.value.condition = 0
   project.value.notes = ""

   // load constants from dpg-imaging
   axios.get(`${systemStore.projectsURL}/constants`).then(response => {
      workflows.value = response.data.workflows
      categories.value = response.data.categories
      isOpen.value = true
   }).catch( err => {
      systemStore.setError(err)
   })
}
</script>

<style lang="scss" scoped>
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 15px 0 10px 0;
   margin: 0;
   gap: 10px;
}
</style>
