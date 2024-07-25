<template>
   <DPGButton @click="show" label="Add Workstation"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Add Workstation" :style="{width: '400px'}" :closable="false">
      <FormKit type="form" id="add-workstation" :actions="false" @submit="addWorkstation">
         <FormKit label="Name" type="text" v-model="workstationName" outer-class="first" required autofocus/>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" severity="secondary"/>
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

const addWorkstation= (async () => {
   await equipmentStore.addWorkstation( workstationName.value)
   hide()
})

const hide = (() => {
   isOpen.value=false
})

const show = (() => {
   workstationName.value = ""
   isOpen.value = true
})
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
   gap: 10px;
}
</style>
