<template>
   <UCard title="Master Files" :ui="{body: 'pt-0!'}">
      <div class="row-between bg-white border-b border-brand-grey-100 py-3 px-2 items-center">
         <div class="row-left gap-2 items-baseline">
            <template v-if="detail.reorder==false && userStore.isAdmin">
               <UButton size="sm" label="Add" @click="addClicked()" color="secondary" :loading="unitsStore.updateInProgress" />
               <UButton size="sm" label="Replace" @click="replaceClicked()" color="secondary" :loading="unitsStore.updateInProgress" />
            </template>
            <template v-if="userStore.isAdmin && (detail.dateArchived==null || detail.reorder || detail.dateDLDeliverablesReady == null)">
               <UButton size="sm" label="Delete" @click="deleteClicked()" color="secondary" :disabled="!filesSelected" />
            </template>
            <RenumberDialog v-if="userStore.isAdmin || userStore.isSupervisor" :disabled="!filesSelected" :filenames="selectedFileNames" />
            <template v-if="unitsStore.canDownload && detail.throwAway == false">
               <UButton size="sm" label="Download" @click="downloadClicked()" color="secondary" :disabled="!filesSelected" />
            </template>
            <template v-if="unitsStore.canPDF">
               <UButton size="sm" label="PDF" @click="pdfClicked()" color="secondary" :disabled="filesSelected == false" />
            </template>
            <template  v-if="(userStore.isAdmin || userStore.isSupervisor) && detail.throwAway == false">
               <LookupDialog :disabled="!filesSelected" label="Assign Metadata" @selectedObject="assignMetadata" target="metadata" :create="true"/>
               <LookupDialog :disabled="!filesSelected" label="Assign Componment" @selected="assignComponent" target="component" />
            </template>    
         </div>
         <span>{{ table?.tableApi?.getFilteredRowModel().rows.length }} Master Files<span v-if="filters.length > 0"> (filtered)</span></span>
      </div>
      <UTable ref="table" :data="unitsStore.masterFiles" :columns="columns" v-model:column-filters="filters" sticky virtualize class="h-200">
         <template #select-header="{ table }">
            <UCheckbox size="lg" :modelValue="table.getIsSomePageRowsSelected() ? 'indeterminate' : table.getIsAllPageRowsSelected()" @update:modelValue="table.toggleAllRowsSelected()"/>
         </template>
         <template #select-cell="{ row }">
            <UCheckbox size="lg" :modelValue="row.getIsSelected()" @update:modelValue="row.toggleSelected($event)"/>
         </template>
         <template #id-cell="{ row }">
            <router-link :to="`/masterfiles/${row.original.id}`">{{row.original.id}}</router-link>
         </template>
         <template #metadata-header="{ column }">
            <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
               @clear="clearFilter(column)" @apply="applyFilter(column)">
               <UInput v-model="pendingFilter" placeholder="Metadata..." />
            </FilterPopover>
         </template>
         <template #metadata-cell="{ row }">
            <router-link v-if="row.original.metadata" :to="`/metadata/${row.original.metadata.id}`">{{row.original.metadata.title}}</router-link>
            <span v-else class="empty">N/A</span>
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
         <template #thumbnailURL-cell="{ row }">
            <a :href="row.original.viewerURL" target="_blank">
               <img :src="row.original.thumbnailURL" :class="{exemplar: row.original.exemplar}"/>
            </a>
         </template>
         <template #actions-cell="{ row }">
            <UDropdownMenu :items="getDropdownActions(row.original)">
               <UButton icon="i-lucide-ellipsis-vertical" color="neutral" variant="ghost" aria-label="Actions"/>
            </UDropdownMenu>
         </template>
      </UTable>
   </UCard>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useUserStore } from '@/stores/user'
import { usePDFStore } from '@/stores/pdf'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import FilterPopover from '@/components/FilterPopover.vue'
import { useConfirm } from "@/composables/useConfirm"
import RenumberDialog from './RenumberDialog.vue'
import LookupDialog from '@/components/LookupDialog.vue'

const unitsStore = useUnitsStore()
const userStore = useUserStore()
const pdfStore = usePDFStore()
const router = useRouter()

const columns = [
   { accessorKey: "select", header: ""},
   { accessorKey: 'id', header: 'ID'},
   { accessorKey: 'metadata.title', id: 'metadata', header: 'Metadata'},
   { accessorKey: 'filename', header: 'File Name'},
   { accessorKey: 'title', header: 'Title'},
   { accessorKey: 'description', header: 'Description'},
   { accessorKey: 'thumbnailURL', header: 'Thumb'},
   { id: 'actions', header: "" }
]

const table = ref()
const filters = ref([])
const pendingFilter = ref()
const { detail } = storeToRefs(unitsStore)

const getDropdownActions = ( (item) =>{
   let out = [] 
   out.push( [ {label: "Actions",  type: "label"} ] )
   out.push( [ {label: "View",  onSelect: () => { viewClicked(item) }} ] )
   if (detail.value.throwAway == false && unitsStore.canDownload) {
      out.push( [ {label: "Download Image",  onSelect: () => { downloadFile(item) }} ] ) 
   }
   out.push( [ {label: "Download PDF",  onSelect: () => { downloadPDF(item) }} ] )
   if (item.exemplar==false && (detail.value.intendedUse && detail.value.intendedUse.id == 110 || detail.value.includeInDL)) {
      out.push( [ {label: "Set Exemplar",  onSelect: () => { exemplarClicked(item) }} ] ) 
   }
   if (detail.value.throwAway == false && detail.value.reorder==false && userStore.isAdmin) {
      out.push( [ {label: "Republish IIIF",  onSelect: () => { republishIIIF(item) }} ] ) 
   }
   return out
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

const downloadClicked = (() => {
   unitsStore.downloadFromArchive(userStore.computeID, selectedFileNames.value)
})
const pdfClicked = (() => {
   requestPDF( selectedIDs.value )
})
const requestPDF = ( async ( masterFileIDs ) => {
   if (unitsStore.hasText == false ) {
      pdfStore.requestPDF( unitsStore.detail.id, masterFileIDs, false, false )
   } else {
      const msg = 'This unit has transcription or OCR text. Include it with the PDF?'
      const resp = await useConfirm("Include Text", msg, "Include")
      if (resp) {
         pdfStore.requestPDF( unitsStore.detail.id, masterFileIDs, false, true )
      } else {
         pdfStore.requestPDF( unitsStore.detail.id, masterFileIDs, false, false )  
      }
   }
})
const downloadPDF = ((info) => {
   requestPDF( [info.id] )
})

const replaceClicked = ( async () => {
   let unitDir = `${unitsStore.detail.id}`.padStart(9, '0')
   const msg = `Replace master files with .tif files from ./finalization/unit_update/${unitDir}?`
   const resp = await useConfirm("Confirm Replace Master Files", msg, "Replace")
   if (resp) {
      unitsStore.replaceMasterFiles()
      clearSelections()
   } 
})

const addClicked = (async () => {
   let unitDir = `${unitsStore.detail.id}`.padStart(9, '0')
   const msg = `Add all .tif files from ./finalization/unit_update/${unitDir} to this unit?`
   const resp = await useConfirm("Confirm Add Master Files", msg, "Add")
   if (resp) {
      unitsStore.addMasterFiles()
   }  
})

const deleteClicked = (async () => {
   const msg = 'Are you sure you want delete the selected master files? All data will be lost. This cannot be reversed.'
   const resp = await useConfirm("Confirm Delete Master Files", msg, "Delete")
   if (resp) {
      unitsStore.deleteMasterFiles(selectedFileNames.value)
      clearSelections()
   }  
})
const downloadFile =(( info) => {
   unitsStore.downloadFromArchive(userStore.computeID, info.filename )
})
const viewClicked = ((data) => {
   router.push(`/masterfiles/${data.id}`)
})
const exemplarClicked = (( info) => {
   unitsStore.setExemplar( info.id )
})
const republishIIIF =((masterFileID) => {
   unitsStore.regenerateMasterFileIIIF( masterFileID )
})

const assignMetadata = ( async ( metadata ) => {
   await unitsStore.assignMetadata(metadata, selectedIDs.value)
   clearSelections()
})

const assignComponent= ( async ( componentID ) => {
   await unitsStore.assignComponent(componentID, selectedIDs.value)
   clearSelections()
})

const clearSelections = (() => {
   selectAll.value = false
   selectedMasterFiles.value = []
})
</script>

<style scoped lang="scss">
</style>