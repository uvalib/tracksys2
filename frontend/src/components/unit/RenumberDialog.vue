<template>
   <DPGButton @click="show" severity="secondary" label="Renumber" :disabled="props.disabled"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Renumber Master Files" :style="{width: '400px'}" :closable="false">
      <label for="start">Starting page number:</label>
      <InputNumber id="start" v-model="startPage" :useGrouping="false" fluid/>  
      <template #footer>
         <DPGButton @click="hide" label="Cancel" severity="secondary"/>
         <DPGButton @click="renumberPages" label="Renumber"  />
      </template>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import Dialog from 'primevue/dialog'
import { useUnitsStore } from '@/stores/units'
import InputNumber from 'primevue/inputnumber'

const props = defineProps({
   disabled: {
      type: Boolean,
      default: false
   },
   filenames: {
      type: Array,
      required: true
   }
})

const unitsStore = useUnitsStore()

const isOpen = ref(false)
const startPage = ref(1)

async function renumberPages() {
   await unitsStore.renumberPages(startPage.value, props.filenames)
   hide()
}
function hide() {
   isOpen.value=false
}
function show() {
   startPage.value = 1
   isOpen.value = true

}
</script>

<style lang="scss" scoped>
label {
   display: block;
   margin-bottom: 5px;
}
</style>
