<template>
   <div class="row-right p-2 sticky z-50 bg-white border-b-1 border-brand-grey-100" :style="{top: headerHeight}">
      <UPagination color="neutral" variant="ghost"
         v-model:page="searchStore.orders.currPage" :items-per-page="searchStore.orders.limit" 
         :total="searchStore.orders.total" @update:page="pageChanged"
      />
      <USelect v-model="searchStore.orders.limit" :items="[15,30,100]" @change="perPageChanged" />
   </div>
   <UTable :data="searchStore.orders.hits" :columns="columns">
      <template #id-cell="{ row }">
         <router-link :to="`/orders/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
      <template #status-header>
         <FilterPopover label="Status" :applied="filterApplied('status')" @clear="clearFilter('status')" @apply="applyFilter()">
            <USelect v-model="filters.status.value" :items="orderStatuses" placeholder="Select a status" />
         </FilterPopover>
      </template>
      <template #status-cell="{ row }">
         <span :class="`status ${row.original.status}`">{{displayStatus(row.original.status)}}</span>
      </template>
      <template #customer-header>
         <FilterPopover label="Customer"  :applied="filterApplied('customer')" @clear="clearFilter('customer')" @apply="applyFilter()">
            <UInput v-model="filters.customer.value" placeholder="Customer name..." />
         </FilterPopover>
      </template>
      <template #agency-header>
         <FilterPopover label="Agency"  :applied="filterApplied('agency')" @clear="clearFilter('agency')" @apply="applyFilter()">
            <UInput v-model="filters.agency.value" placeholder="Agency name..." />
         </FilterPopover>
      </template>
      <template #title-header>
         <FilterPopover label="Title"  :applied="filterApplied('title')" @clear="clearFilter('title')" @apply="applyFilter()">
            <UInput v-model="filters.title.value" placeholder="Title..." />
         </FilterPopover>
      </template>
      <template #title-cell="{ row }">
         <div class="max-width: 150px; white-space: break-spaces;">{{ row.original.title   }}</div>
      </template>
      <template #staff_notes-header>
         <FilterPopover label="Staff Notes"  :applied="filterApplied('staff_notes')" @clear="clearFilter('staff_notes')" @apply="applyFilter()">
            <UInput v-model="filters.staff_notes.value" placeholder="Staff notes..." />
         </FilterPopover>
      </template>
      <template #special_instructions-header>
         <FilterPopover label="Special Instructions" :applied="filterApplied('special_instructions')" @clear="clearFilter('special_instructions')" @apply="applyFilter()">
            <UInput v-model="filters.special_instructions.value" placeholder="Special instructions..." />
         </FilterPopover>
      </template>
   </UTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSearchStore } from '../../stores/search'
import { useRoute, useRouter } from 'vue-router'
import FilterPopover from './FilterPopover.vue'

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const orderHitsTable = ref()

// you cannot custruct tailwind class values dynamically, so you  cant do `top-${hdr.clientHeight}`. 
// Instead use this to bind an inline style 'top' param to stick the controls below the header
const headerHeight = computed(() => {
   let hdr = document.getElementById('uva-header')
   return `${hdr.clientHeight}px`
})

const columns = [
   {
      accessorKey: 'id',
      header: 'ID'
   },
   {
      accessorKey: 'status',
      header: 'Status',
      meta: {
         class: {
            td: 'w-fit'
         }
      }
   },  
   {
      accessorKey: 'customer',
      header: 'Customer',
      meta: {
         class: {
            td: 'w-fit'
         }
      }
   }, 
   {
      accessorKey: 'agency',
      header: 'Agency',
      meta: {
         class: {
            td: 'w-1/8 whitespace-break-spaces'
         }
      }
   }, 
   {
      accessorKey: 'title',
      header: 'Order Title',
      meta: {
         class: {
            td: 'w-1/4 whitespace-break-spaces'
         }
      }
   }, 
   {
      accessorKey: 'staff_notes',
      header: 'Staff Notes',
      meta: {
         class: {
            td: 'w-1/4 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'special_instructions',
      header: 'Special Instructions',
      meta: {
         class: {
            td: 'w-1/4 whitespace-break-spaces'
         }
      }
   },
]

const filters = ref({
   status: {value: null, mode: "equals"},
   customer: {value: null, mode: "contains"},
   agency: {value: null, mode: "contains"},
   title: {value: null, mode: "contains"},
   staff_notes: {value: null, mode: "contains"},
   special_instructions: {value: null, mode: "contains"},
})

const orderStatuses = ref([
   {label: "Requested", value: "requested"},
   {label: "Approved", value: "approved"},
   {label: "Await Fee", value: "await_fee"},
   {label: "Completed", value: "completed"},
   {label: "Canceled", value: "canceled"},
   {label: "Deferred", value: "deferred"},
])

const hasFilter = computed(() => {
   let idx = Object.values(filters.value).findIndex( fv => fv.value && fv.value != "")
   return idx >= 0
})

onMounted(() =>{
   searchStore.orders.filters.forEach( fv => {
      filters.value[fv.field].value = fv.value
   })
})

function downloadCSV() {
   orderHitsTable.value.exportCSV()
}

function displayStatus( id) {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
}

function clearFilters() {
   Object.values(filters.value).forEach( fv => fv.value = null )
   searchStore.orders.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.executeSearch("orders")
}

const filterApplied = ((name) => {
   return filters.value[name].value != null 
})
const clearFilter = ((name) => {
   filters.value[name].value = null 
   applyFilter()
}) 

const applyFilter =(() => {
   searchStore.orders.filters = []
   Object.entries(filters.value).forEach(([filterName, data]) => {
      if (data.value && data.value != "") {
         searchStore.orders.filters.push({field: filterName, match: data.mode, value: data.value})
      }
   })
   let query = Object.assign({}, route.query)
   query.filters = searchStore.filtersAsQueryParam("orders")
   delete query.filters
   if ( searchStore.orders.filters.length > 0) {
      query.filters = searchStore.filtersAsQueryParam("orders")
   }
   router.push({query})
   searchStore.executeSearch("orders")
})

const perPageChanged = (() => {
   console.log("new page size "+searchStore.orders.limit)
   searchStore.orders.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   searchStore.orders.start = (searchStore.orders.currPage-1) * searchStore.orders.limit
   console.log("new start: "+searchStore.orders.start)
   searchStore.executeSearch("orders")
})

</script>

<style scoped lang="scss">
.acts{
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: center;
   gap: 10px;
}
</style>