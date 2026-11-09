<template>
   <div class="row-right bg-white border-b border-brand-grey-100 py-3 px-2 items-center">
      <span>{{ table?.tableApi?.getFilteredRowModel().rows.length }} Images<span v-if="filters.length > 0"> (filtered)</span></span>
   </div>
   <UTable ref="table" :columns="columns" v-model:column-filters="filters"  :data="props.masterFiles" sticky virtualize class="h-150">
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
            <img :src="row.original.thumbnailURL"/>
         </a>
      </template>
   </UTable>
</template>

<script setup>
import { ref } from 'vue'
import FilterPopover from '@/components/FilterPopover.vue'

const props = defineProps({
   masterFiles: {
      type: Array,
      required: true
   }
})

const columns = [
   { accessorKey: 'id', header: "ID" },
   { accessorKey: 'pid', header: "PID" },
   { accessorKey: 'title', header: "Title" },
   { accessorKey: 'description', header: "Description" },
   { accessorKey: 'thumbnailURL', header: "Thumb" }
]

const table = ref()
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

</script>

<style scoped lang="scss">
</style>