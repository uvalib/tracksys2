<template>
   <UModal v-model:open="isOpen" :modal="true" size="sm" :dismissible="false" :close="false" title="Add to Collection" >
      <UButton @click="showClicked" label="Add to Collection" color="secondary" />
      <template #body>
         <UListbox v-model="collectionID" :items="collectionStore.collections" valueKey="id" labelKey="title" virtualize />
      </template>
      <template #footer={close}>
         <UButton @click="close" color="secondary" label="Cancel"/>
         <UButton @click="addToCollection" label="Add" :disabled="collectionID==0"/>
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useCollectionsStore } from '@/stores/collections'

const props = defineProps({
   metadataID: {
      type: Number,
      required: true,
   },
})

const collectionStore = useCollectionsStore()

const isOpen = ref(false)
const collectionID = ref(0)

const addToCollection = (async () => {
   await collectionStore.addToCollection( collectionID.value, props.metadataID )
   isOpen.value=false
})

const showClicked = (() => {
   isOpen.value = true
   collectionID.value = 0
   collectionStore.getCollections()
   
})
</script>

<style lang="scss" scoped>
</style>
