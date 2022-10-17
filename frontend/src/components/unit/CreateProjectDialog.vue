<template>
   <DPGButton @click="show" class="p-button-secondary" label="Create Digitization Project"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Create Digitization Project" :style="{width: '400px'}">
      <FormKit type="form" id="create-project" :actions="false" @submit="createProject">
         <FormKit label="Workflow" type="select" v-model="project.workflowID" :options="workflows" required outer-class="first" />
         <FormKit label="Category" type="select" v-model="project.categoryID" :options="categories" required/>
         <FormKit label="Date Due" type="date" v-model="project.dueOn" required/>
         <FormKit label="Condition" type="select" v-model="project.condition" :options="conditions" required/>
         <FormKit label="" type="textarea" rows="4" v-model="project.notes"/>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" class="p-button-secondary"/>
            <span class="spacer"></span>
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

const unitsStore = useUnitsStore()
const systemStore = useSystemStore()

const isOpen = ref(false)
const project = ref({
   workflowID: 1,
   categoryID: 0,
   condition: 0,
   dueOn: null,
   notes: ""
})

const workflows = computed( () => {
   let out = []
   systemStore.workflows.forEach( w => {
      out.push({label: w.name, value: w.id})
   })
   return out
})
const categories = computed( () => {
   let out = []
   systemStore.categories.forEach( w => {
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
   await unitsStore.createProject(project.value)
   hide()
}
function hide() {
   isOpen.value=false
}
function show() {
   project.value.workflowID = 1
   project.value.categoryID = 0
   project.value.condition = 0
   project.value.dueOn = unitsStore.detail.order.dateDue.split("T")[0]
   project.value.notes = ""
   isOpen.value = true

}
</script>

<style lang="scss" scoped>
.spacer {
   display: inline-block;
   margin: 0 5px;
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 15px 0 10px 0;
   margin: 0;
   button {
      margin-right: 10px;
   }
}
</style>
