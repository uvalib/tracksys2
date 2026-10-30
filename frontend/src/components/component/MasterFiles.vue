<template>
   <UCard :title="props.title" :ui="{body: 'pt-0!'}">
      <div class="row-between bg-white border-b border-brand-grey-100 py-3 px-2 items-center">
         <div class="row-left gap-2">
            <UButton label="PDF" @click="pdfClicked()" color="secondary" :disabled="!filesSelected" />
            <UButton label="Download" @click="downloadClicked()" color="secondary" :disabled="!filesSelected" />
            <template  v-if="userStore.isAdmin || userStore.isSupervisor">
               <LookupDialog :disabled="!filesSelected" label="Assign Metadata" @selected="assignMetadata" target="metadata" :create="true"/>
            </template>
         </div>
         <span>{{ table?.tableApi?.getFilteredRowModel().rows.length }} Master Files<span v-if="filters.length > 0"> (filtered)</span></span>
      </div>
      <UTable ref="table" :data="component.relatedMasterFiles" :columns="columns" v-model:column-filters="filters"  sticky virtualize class="h-150">
         <template #select-header="{ table }">
            <UCheckbox size="lg" :modelValue="table.getIsSomePageRowsSelected() ? 'indeterminate' : table.getIsAllPageRowsSelected()" @update:modelValue="table.toggleAllRowsSelected()"/>
         </template>
         <template #select-cell="{ row }">
            <UCheckbox size="lg" :modelValue="row.getIsSelected()" @update:modelValue="row.toggleSelected($event)"/>
         </template>
         <template #id-cell="{ row }">
            <router-link :to="`/masterfiles/${row.original.id}`">{{row.original.id}}</router-link>
         </template>  
         <template #filename-header="{ column }">
            <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
               @clear="clearFilter(column)" @apply="applyFilter(column)">
               <UInput v-model="pendingFilter" placeholder="File name..." />
            </FilterPopover>
         </template> 
         <template #title-header="{ column }">
            <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
               @clear="clearFilter(column)" @apply="applyFilter(column)">
               <UInput v-model="pendingFilter" placeholder="Title..." />
            </FilterPopover>
         </template> 
         <template #description-header="{ column }">
            <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
               @clear="clearFilter(column)" @apply="applyFilter(column)">
               <UInput v-model="pendingFilter" placeholder="Description..." />
            </FilterPopover>
         </template> 
         <template #metadata-header="{ column }">
            <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
               @clear="clearFilter(column)" @apply="applyFilter(column)">
               <UInput v-model="pendingFilter" placeholder="Metadata..." />
            </FilterPopover>
         </template> 
         <template #metadata-cell="{ row }">
            <router-link :to="`/metadata/${row.original.metadata?.id}`">{{row.original.metadata?.title}}</router-link>
         </template>
         <template #thumbnailURL-cell="{ row }">
            <a :href="row.original.viewerURL" target="_blank">
               <img :src="row.original.thumbnailURL" :class="{exemplar: row.original.exemplar}"/>
            </a>
         </template>
      </UTable>
   </UCard>

   <UModal v-model:open="pdfStore.downloading" :modal="true" :dismissible="false" title="Generating PDF">
      <template #body>
         <div>PDF generation in progress...</div>
         <UProgress v-model="pdfStore.percent" status />
      </template>
   </UModal>
</template>

<script setup>
import { useComponentsStore } from '@/stores/components'
import { useUserStore } from '@/stores/user'
import { usePDFStore } from '@/stores/pdf'
import { ref, computed, watch } from 'vue'
import LookupDialog from '@/components/LookupDialog.vue'
import FilterPopover from '@/components/FilterPopover.vue'

const props = defineProps({
   title: {
      type: String,
      required: true
   },
})

const columns = [
   { accessorKey: "select", header: ""},
   { accessorKey: 'id', header: 'ID'},
   { accessorKey: 'pid', header: 'PID'},
   { accessorKey: 'filename', header: 'File Name'},
   { accessorKey: 'title', header: 'Title'},
   { accessorKey: 'description', header: 'Description'},
   { 
      id: 'metadata', 
      header: 'Metadata',
      accessorFn: (row) => row.metadata?.title
   },
   { accessorKey: 'thumbnailURL', header: 'Thumb'},
]

const component = useComponentsStore()
const userStore = useUserStore()
const pdfStore = usePDFStore()

const table = ref()
const filters = ref([])
const pendingFilter = ref()

const filesSelected = computed(() => {
   return table.value?.tableApi?.getIsSomePageRowsSelected()
})
const selectedFileNames = computed(() => {
   let out = table.value?.tableApi?.getSelectedRowModel().rows.map( r => r.original.filename )
   if (!out) return [] 
   return out
})
const selectedIDs = computed(() => {
   let out = table.value.tableApi.getSelectedRowModel().rows.map( r => r.original.id )
   if (!out) return [] 
   return out
})

watch(() => component.loadingMasterFiles, (newVal) => {
   if ( newVal == false) {
      table.value?.tableApi.resetRowSelection()
   }
})

const downloadClicked = (() => {
   const unitID = table.value.tableApi.getSelectedRowModel().rows[0].original.unitID
   component.downloadFromArchive(userStore.computeID, unitID, selectedFileNames.value)
   clearSelections()
})

const pdfClicked = (() => {
   const unitID = table.value.tableApi.getSelectedRowModel().rows[0].original.unitID
   pdfStore.requestPDF( unitID, selectedIDs.value )
})

const assignMetadata = ( async ( metadataID ) => {
   const unitID = table.value.tableApi.getSelectedRowModel().rows[0].original.unitID
   await component.assignMetadata(metadataID, unitID, selectedIDs.value)
   clearSelections()
})

const clearSelections = (() => {
   table.value?.tableApi.resetRowSelection()
})

const isFilterApplied = ((name) => {
   return filters.value.findIndex( f => f.id == name) > -1
})
const filterOpened = ((column) => {
   pendingFilter.value = column.getFilterValue()  
})
const clearFilter = ((column) => {
   column.setFilterValue(null)
})
const applyFilter = ((column) => {
   column.setFilterValue(pendingFilter.value)
})
</script>

<style scoped lang="scss">
img.exemplar {
   border: 5px solid var(--uvalib-teal);
}
</style>