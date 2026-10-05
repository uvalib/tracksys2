<template>
   <div class="row-between p-2 sticky z-50 bg-white border-b-1 border-brand-grey-100" :style="{top: headerHeight}">
      <div class="row-left gap-2!">
         <UButton color="secondary" label="Download Results CSV" @click="downloadCSV"/>
         <UButton v-if="hasFilter" color="secondary" label="Clear All Filters" @click="clearFilters"/>
      </div>
      <div class="row-right">
         <UPagination color="neutral" variant="ghost"
            v-model:page="searchStore.units.currPage" :items-per-page="searchStore.units.limit" 
            :total="searchStore.units.total" @update:page="pageChanged"
         />
         <USelect v-model="searchStore.units.limit" :items="[15,30,100]" @change="perPageChanged" />
      </div>
   </div>
   <UTable :data="searchStore.units.hits" :columns="columns">
      <template #id-cell="{ row }">
         <router-link :to="`/units/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
      <template #status-header>
         <FilterPopover label="Status" :applied="filterApplied('status')" @clear="clearFilter('status')" @apply="applyFilter()">
            <USelect v-model="filters.status.value" :items="unitStatuses" placeholder="Select a status" />
         </FilterPopover>
      </template>
      <template #status-cell="{ row }">
         <span :class="`status ${row.original.status}`">{{displayStatus(row.original.status)}}</span>
      </template>
      <template #staff_notes-header>
         <FilterPopover label="Staff Notes"  :applied="filterApplied('staff_notes')" @clear="clearFilter('staff_notes')" @apply="applyFilter()">
            <UInput v-model="filters.staff_notes.value" placeholder="Staff notes..." />
         </FilterPopover>
      </template>
      <template #special_instructions-header>
         <FilterPopover label="Special Instructions"  :applied="filterApplied('special_instructions')" @clear="clearFilter('special_instructions')" @apply="applyFilter()">
            <UInput v-model="filters.special_instructions.value" placeholder="Instructions..." />
         </FilterPopover>
      </template>
      <template #date_dl_deliverables_ready-cell="{ row }">
         <span>{{ $formatDate(row.original.date_dl_deliverables_ready) }}</span>  
      </template>
      <template #date_patron_deliverables_ready-cell="{ row }">
         <span>{{ $formatDate(row.original.date_patron_deliverables_ready) }}</span>  
      </template>
   </UTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSearchStore } from '../../stores/search'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

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
      header: 'Status'
   },
   {
      accessorKey: 'staff_notes',
      header: 'Staff Notes',
      meta: {
         class: {
            td: 'w-1/3 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'special_instructions',
      header: 'Special Instructions',
      meta: {
         class: {
            td: 'w-1/3 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'date_dl_deliverables_ready',
      header: 'DL Deliverable Date'
   },
   {
      accessorKey: 'date_patron_deliverables_ready',
      header: 'Patron Deliverable Date'
   },
]

const filters = ref({
   status: {value: null, mode: "equals"},
   staff_notes: {value: null, mode: "contains"},
   special_instructions: {value: null, mode: "contains"},
})

const unitStatuses = ref([
   {label: "Approved", value: "approved"},
   {label: "Unapproved", value: "unapproved"},
   {label: "Canceled", value: "canceled"},
   {label: "Done", value: "done"},
   {label: "Error", value: "error"},
])

const hasFilter = computed(() => {
   let idx = Object.values(filters.value).findIndex( fv => fv.value && fv.value != "")
   return idx >= 0
})

onMounted(() => {
   searchStore.units.filters.forEach( fv => {
      filters.value[fv.field].value = fv.value
   })
})

const displayStatus = ((id) => {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
})

const downloadCSV = (() => {
   searchStore.downloadCSV('units', columns )
})

const clearFilters = (() => {
   filters.value = {
      system_name: {value: null, mode: "equals"},
      title: {value: null, mode: "contains"},
      creator_name: {value: null, mode: "contains"},
      barcode: {value: null, mode: "startsWith"},
      call_number: {value: null, mode: "startsWith"},
      catalog_key: {value: null, mode: "startsWith"},
      virgo: {value: null, mode: "equals"},
      dpla: {value: null, mode: "equals"},
      hathitrust: {value: null, mode: "equals"}
   }
   searchStore.units.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.executeSearch("units")
})

const filterApplied = ((name) => {
   return filters.value[name].value != null 
})
const clearFilter = ((name) => {
   filters.value[name].value = null 
   applyFilter()
}) 

const applyFilter =(() => {
   searchStore.units.filters = []
   Object.entries(filters.value).forEach(([filterName, data]) => {
      if (data.value && data.value != "") {
         searchStore.units.filters.push({field: filterName, match: data.mode, value: data.value})
      }
   })
   let query = Object.assign({}, route.query)
   query.filters = searchStore.filtersAsQueryParam("units")
   delete query.filters
   if ( searchStore.units.filters.length > 0) {
      query.filters = searchStore.filtersAsQueryParam("units")
   }
   router.push({query})
   searchStore.executeSearch("units")
})

const perPageChanged = (() => {
   searchStore.units.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   searchStore.units.start = (searchStore.units.currPage-1) * searchStore.units.limit
   searchStore.executeSearch("units")
})

</script>

<style scoped lang="scss">
</style>