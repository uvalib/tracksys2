<template>
   <UModal v-model:open="open" :modal="true" :dismissible="false" :close="false" title="Create Metadata">
      <UButton @click="open = true" label="Create Metadata" />
      <template #body>
         <NewMetadataPanel @canceled="open = false" @created="metadataCreated" class="w-full"/>
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useSystemStore } from '../stores/system'
import { useMetadataStore } from '../stores/metadata'
import NewMetadataPanel from '@/components/metadata/NewMetadataPanel.vue'

const system = useSystemStore()
const metadataStore = useMetadataStore()

const open = ref(false)

const metadataCreated = (() => {
   system.toastMessage("Metadata Created", `Metadata ${metadataStore.detail.pid}: ${metadataStore.detail.title} has been created.`)
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