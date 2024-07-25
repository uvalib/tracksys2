<template>
   <DPGButton @click="show" label="Add Equipment"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Add Equipment" :style="{width: '400px'}" :closable="false">
      <FormKit type="form" id="add-equipment" :actions="false" @submit="addEquipment">
         <FormKit label="Type" type="select" v-model="equipType" outer-class="first" :options="equipmentTypes" placeholder="Select equipment type" required autofocus/>
         <FormKit label="Name" type="text" v-model="name" required/>
         <FormKit label="Serial Number" type="text" v-model="serialNumber" required/>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" severity="secondary"/>
            <FormKit type="submit" label="Add" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import Dialog from 'primevue/dialog'
import { useEquipmentStore } from '@/stores/equipment'

const equipmentStore = useEquipmentStore()

const isOpen = ref(false)
const name = ref("")
const serialNumber = ref("")
const equipType = ref()

const equipmentTypes = computed( () => {
   return [
      {label: "Camera Body", value: "CameraBody"},
      {label: "Digital Back", value: "DigitalBack"},
      {label: "Lens", value: "Lens"},
      {label: "Scanner", value: "Scanner"},
   ]
})

const addEquipment= ( async () => {
   await equipmentStore.addEquipment( equipType.value, name.value, serialNumber.value)
   hide()
})

const hide = (() => {
   isOpen.value=false
})

const show = (() => {
   name.value = ""
   serialNumber.value = ""
   equipType.value = null
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
