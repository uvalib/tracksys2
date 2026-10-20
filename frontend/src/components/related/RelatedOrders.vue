<template>
   <div class="row-right bg-white border-b border-brand-grey-100 py-3 px-2 items-center">
      <span>{{ table?.tableApi?.getFilteredRowModel().rows.length }} Units<span v-if="filters.length > 0"> (filtered)</span></span>
   </div>
   <UTable ref="table" :columns="columns" v-model:column-filters="filters"  :data="props.orders" sticky virtualize class="h-150">
      <template #id-cell="{ row }">
         <router-link :to="`/orders/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
      <template #customer-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Last name..." />
         </FilterPopover>
      </template>
      <template #agency-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <USelectMenu v-model="pendingFilter" :items="systemStore.agencies" class="w-full" 
               descriptionKey="" labelKey="name" valueKey="name" placeholder="Select an agency..." virtualize/>   
         </FilterPopover>
      </template>
      <template #title-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Title..." />
         </FilterPopover>
      </template>
      <template #staffNotes-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Staff notes..." />
         </FilterPopover>
      </template>
      <template #specialInstructions-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Instructions..." />
         </FilterPopover>
      </template>
   </UTable>
</template>

<script setup>
import { ref } from 'vue'
import FilterPopover from '@/components/FilterPopover.vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()

const props = defineProps({
   orders: {
      type: Array,
      required: true
   }
})

const columns = [
   {
      accessorKey: 'id',
      header: 'ID',
   },
   { 
      id: 'customer', 
      header: 'Customer',
      accessorFn: (row) => `${row.customer.lastName}, ${row.customer.firstName}`
   },
   { 
      id: 'agency',
      header: 'Agency',
      accessorKey: 'agency.name'
   },
   { accessorKey: 'title', header: 'Order Title' },
   { accessorKey: 'staffNotes', header: 'Staff Notes' },
   { accessorKey: 'specialInstructions', header: 'Special Instructions' },
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


// const filters = ref( {
//    'agency.id': {value: null, matchMode: FilterMatchMode.EQUALS},
//    'customer.lastName': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
//    'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
//    'specialInstructions': {value: null, matchMode: FilterMatchMode.CONTAINS},
//    'staffNotes': {value: null, matchMode: FilterMatchMode.CONTAINS},
// })
</script>

<style scoped lang="scss">
</style>