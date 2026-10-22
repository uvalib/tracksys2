<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" title="Related Units">
      <UButton @click="isOpen=true" variant="link" label="...view more" />
      <template #body>
         <div class="flex flex-row flex-wrap justify-start gap-1">
            <template v-for="uid in props.units" :key="`related-${uid}`">
               <UButton v-if="uid != currentUnitID" :to="`/units/${uid}`" variant="link" :label="uid"/>
            </template>
         </div>
      </template>
      <template #footer="{close}">
         <UButton @click="close" color="secondary" label="Close" />   
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
   units: {
      type: Array,
      required: true
   },
   orderID: {
      type: Number,
      reqtured: true
   },
   currentUnitID: {
      type: Number,
      reqtured: true
   }
})

const isOpen = ref(false)
</script>

<style lang="scss" scoped>
</style>
