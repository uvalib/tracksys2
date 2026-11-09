<template>
   <UModal v-model:open="showDialog" fullscreen :modal="true" :dismissible="false" title="Add to Collection" >
      <UButton @click="showClicked" label="Add Item(s)" size="sm" color="secondary" />
      <template #body>
         <WaitSpinner v-if="collectionStore.searching" :overlay="true" message="Searching..." />
         <div class="row-between bg-white border-b border-brand-grey-100 py-3 px-2 items-center">
            <div class="row-left gap-2">
               <UInput v-model="query" @keyup.enter="queryMetadata" class="grow" placeholder="Search..." />
               <UButton label="Search" @click="queryMetadata()" :disabled="query.length == 0"/>
               <UButton label="Clear" color="secondary" @click="clearSearch()" :disabled="query.length == 0"/>
            </div>
            <span>{{ collectionStore.metadataHits.length }} Hits</span>
         </div>
         <UTable ref="table" :data="collectionStore.metadataHits" :columns="columns" sticky virtualize class="h-150">
            <template #select-header="{ table }">
               <UCheckbox size="lg" :modelValue="table.getIsSomePageRowsSelected() ? 'indeterminate' : table.getIsAllPageRowsSelected()" @update:modelValue="table.toggleAllRowsSelected()"/>
            </template>
            <template #select-cell="{ row }">
               <UCheckbox size="lg" :modelValue="row.getIsSelected()" @update:modelValue="row.toggleSelected($event)"/>
            </template>
         </UTable>
      </template>
      <template #footer>
         <UButton label="Cancel" color="secondary" @click="cancelClicked()" />
         <UButton label="Add Selected" @click="submitClicked()" :disabled="!(table?.tableApi?.getIsSomePageRowsSelected() || table?.tableApi?.getIsAllPageRowsSelected())" />
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useCollectionsStore } from '@/stores/collections'
import { useSystemStore } from '@/stores/system'
import WaitSpinner from "@/components/WaitSpinner.vue"
import { useConfirm } from "@/composables/useConfirm"

const columns = [ 
   { accessorKey: "select", header: ""},
   { accessorKey: "id", header: "ID"},
   { accessorKey: "pid", header: "PID"},
   { accessorKey: "title", header: "Title"},
   { accessorKey: "barcode", header: "Barcode"},
   { accessorKey: "catalogKey", header: "Catalog Key"},
   { accessorKey: "callNumber", header: "Call Number"}
]

const collectionStore = useCollectionsStore()
const systemStore = useSystemStore()

const showDialog = ref(false)
const table = ref()
const query = ref("")

const showClicked = (() => {
   query.value = ""
   collectionStore.metadataHits = []
   showDialog.value = true
})

const queryMetadata = (() => {
   collectionStore.metadataSearch( query.value )
})
const clearSearch = (() => {
   query.value = ""
   collectionStore.metadataHits = []
})
const cancelClicked = (() => {
   showDialog.value = false
   collectionStore.metadataHits = []
})
const submitClicked = ( async () => {
   let metadataIDs = table.value.tableApi.getSelectedRowModel().rows.map( r => r.original.id )
   const msg = `Add the selected metadata records to the collection?`
   const resp = await useConfirm("Confirm Add", msg, "Add")
   if (resp) {
      await collectionStore.addRecords( metadataIDs )
      showDialog.value = false
      systemStore.toastMessage("Add Started", "The selected items are being added to the collection. Check the job status page for updates.")
   }
})

</script>

<style scoped lang="scss">
</style>