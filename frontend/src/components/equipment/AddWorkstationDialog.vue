<template>
   <DPGButton @click="show" label="Add Workstation"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Add Workstation" :style="{width: '400px'}">
      <FormKit type="form" id="add-workstation" :actions="false" @submit="addWorkstation">
         <FormKit label="Name" type="text" v-model="workstationName" outer-class="first" required autofocus/>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" class="p-button-secondary"/>
            <span class="spacer"></span>
            <FormKit type="submit" label="Add" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import Dialog from 'primevue/dialog'
import { useEquipmentStore } from '@/stores/equipment'

const equipmentStore = useEquipmentStore()

const isOpen = ref(false)
const workstationName = ref("")

async function addWorkstation() {
   await equipmentStore.addWorkstation( workstationName.value)
   hide()
}
function hide() {
   isOpen.value=false
}
function show() {
   workstationName.value = ""
   isOpen.value = true

}
</script>

<style lang="scss" scoped>
button.p-button {
   margin-left: 10px;
   font-size: 0.9em;
}

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
