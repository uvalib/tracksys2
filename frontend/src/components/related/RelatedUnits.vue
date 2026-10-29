<template>
   <div class="row-between bg-white border-b border-brand-grey-100 py-3 px-2 items-center">
      <div class="row-between gap-2">
         <HathiTrustUpdateDialog :orderID="props.orderID" v-if="props.hathiTrust" />
         <AddUnitDialog v-if="props.canAdd"/>
         <UButton label="Download Units CSV" color="secondary" @click="downloadCSV" />
      </div>
      <span>{{ table?.tableApi?.getFilteredRowModel().rows.length }} Units<span v-if="filters.length > 0"> (filtered)</span></span>
   </div>
   <UTable ref="table" :data="props.units" :columns="columns" v-model:column-filters="filters" v-model:sorting="sorting" sticky virtualize class="h-150">
      <template #id-cell="{ row }">
         <router-link :to="`/units/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
      <template #title-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('title')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Title..." />
         </FilterPopover>
      </template>
      <template #title-cell="{ row }">
         <router-link v-if="row.original.metadata" :to="`/metadata/${row.original.metadataID}`">{{row.original.metadata.title}}</router-link>
         <span v-else class="empty">N/A</span>
      </template>
      <template #status-header="{column}">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('status')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <USelect v-model="pendingFilter" :items="statusValues" placeholder="Status..." />
         </FilterPopover>
      </template>
      <template #callNumber-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('callNumber')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Call Number..." />
         </FilterPopover>
      </template>
      <template #intendedUse-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('intendedUse')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <USelect v-model="pendingFilter" :items="system.intendedUses" labelKey="name" valueKey="name" placeholder="Intended use..." />
         </FilterPopover>
      </template>
      <template #datePatronDeliverablesReady-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('datePatronDeliverablesReady')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Date..." />
         </FilterPopover>
      </template>
      <template #datePatronDeliverablesReady-cell="{ row }">
         <span v-if="row.original.datePatronDeliverablesReady">{{ toCalendarDate(parseAbsolute(row.original.datePatronDeliverablesReady)).toString() }}</span>
      </template>
      <template #dateDLDeliverablesReady-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('dateDLDeliverablesReady')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Date..." />
         </FilterPopover>
      </template>
      <template #dateDLDeliverablesReady-cell="{ row }">
         <span v-if="row.original.dateDLDeliverablesReady">{{ toCalendarDate(parseAbsolute(row.original.dateDLDeliverablesReady)).toString() }}</span>
      </template>
      <template #reorder-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('reorder')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <USelect v-model="pendingFilter" :items="yesNo" placeholder="Select a value..." />
         </FilterPopover>
      </template>
   </UTable>
</template>

<script setup>
import { ref, computed } from 'vue'
import AddUnitDialog from '@/components/order/AddUnitDialog.vue'
import HathiTrustUpdateDialog from '@/components/HathiTrustUpdateDialog.vue'
import FilterPopover from '@/components/FilterPopover.vue'
import { useSystemStore } from '@/stores/system'
import { toCalendarDate, parseAbsolute } from '@internationalized/date'

const props = defineProps({
   orderID: {
      type: Number,
      default: -1
   },
   units: {
      type: Array,
      required: true
   },
   canAdd: {
      type: Boolean,
      default: false
   },
   showMetadata: {
      type: Boolean,
      default: true
   },
   hathiTrust: {
      type: Boolean,
      default: false
   },
})

const columns = [
   {
      accessorKey: 'id',
      header: 'ID',
   },
   {
      accessorKey: 'metadata.title',
      id: 'title',
      header: 'Title',
      enableSorting: true,
   },
   {
      accessorKey: 'status',
      header: 'Status',
   },
   {
      accessorKey: 'metadata.callNumber',
      id: 'callNumber',
      header: 'Call Number',
   },
   {
      id: 'intendedUse',
      accessorFn: (row) => row.intendedUse ? row.intendedUse.name : '',
      header: 'Intended Use',
   },
   {
      accessorKey: 'datePatronDeliverablesReady',
      header: 'Date Patron Deliverables Ready',
   },
   {
      accessorKey: 'dateDLDeliverablesReady',
      header: 'Date DL Deliverables Ready',
   },
   {
      accessorKey: 'reorder',
      header: 'Reorder?',
   },
   {
      accessorKey: 'masterFilesCount',
      header: 'Master Files Count',
   },
]

const system = useSystemStore()

const table = ref()
const sorting = ref([])
const filters = ref([])
const pendingFilter = ref()

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

const statusValues = computed(() => {
   let out = []
   out.push( {label: "Approved", value: 'approved'} )
   out.push( {label: "Canceled", value: 'canceled'} )
   out.push( {label: "Done", value: 'done'} )
   out.push( {label: "Error", value: 'error'} )
   out.push( {label: "Unapproved", value: 'unapproved'} )
   return out
})

const yesNo = computed(() => {
   let out = []
   out.push( {label: "No", value: false} )
   out.push( {label: "Yes", value: true} )
   return out
})

const downloadCSV = (() => {
   const header = columns.map( c => `"${c.header}"`)
   const fields = columns.map( c => c.accessorKey )

   let csv = "data:text/csv;charset=utf-8,"
   csv += (header.join(",")+"\n")
   props.units.forEach( r => {
      fields.forEach( fieldName => {
         let val = ""
         if (fieldName.includes(".")) {
            let bits = fieldName.split(".")
            let child = r[bits[0]]
            val = child[bits[1]]
         } else {
            val = `${r[fieldName]}`
         }
         val = val.replaceAll("\"", "'")
         csv += `"${val}",`
      })
      csv += "\n"
   })

   // create a link for the data and download it
   var encodedUri = encodeURI(csv)
   var link = document.createElement("a")
   link.setAttribute("href", encodedUri)
   link.setAttribute("download", "units.csv")
   document.body.appendChild(link)
   link.click()
   document.body.removeChild(link)
})

</script>

<style scoped lang="scss">
</style>