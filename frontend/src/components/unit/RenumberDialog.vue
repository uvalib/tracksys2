<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" title="Renumber Master Files">
      <UButton @click="show()" label="Renumber" :disabled="props.disabled" size="sm" color="secondary"/>
      <template #body>
         <label for="start">Starting page number:</label>
         <UInputNumber v-model="startPage" class="w-full"/>
      </template>
      <template #footer>
         <UButton @click="hide" label="Cancel" color="secondary"/>
         <UButton @click="renumberPages" label="Renumber"  />
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useUnitsStore } from '@/stores/units'

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
