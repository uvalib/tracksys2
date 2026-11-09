<template>
   <UModal v-model:open="open" :modal="true" :dismissible="false" :close="false" title="Create Collection Facet">
      <UButton @click="show()" label="Create Collection Facet" />
      <template #body>
         <div class="row">
            <label>Enter the name of the new collection facet</label>
            <UInput v-model="newCollectionFacet" autofocus/>
         </div>
      </template>
      <template #footer="{ close }">
         <UButton label="Cancel" color="secondary" @click="close" />
         <UButton label="Create" @click="createCollection()" />
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useSystemStore } from '../stores/system'

const system = useSystemStore()

const open = ref(false)
const newCollectionFacet = ref("")

const show = ( () => {
   newCollectionFacet.value = ""
   open.value = true
})

const createCollection = ( async () => {
   await system.createCollectionFacet(newCollectionFacet.value)
   open.value = false
})
</script>

<style scoped lang="scss">
.row {
   display: flex;
   flex-direction: column;
   gap: 5px;  
}
</style>