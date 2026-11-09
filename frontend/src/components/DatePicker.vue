<template>
   <UInputDate ref="inputDate" v-model="model" >
      <template #trailing>
         <UPopover :reference="inputDate?.inputsRef[3]?.$el">
            <UButton color="neutral" variant="link" size="sm" icon="i-lucide-calendar" 
               aria-label="Select a date" class="px-0" />
            <template #content>
               <UCalendar v-model="model" class="p-2" />
            </template>
         </UPopover>
      </template>
   </UInputDate>
</template>

<script setup>
import { useTemplateRef } from 'vue'

const inputDate = useTemplateRef('inputDate')
const model = defineModel()
</script>


<style lang="scss" scoped>
</style>
