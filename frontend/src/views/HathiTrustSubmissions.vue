<template>
   <h2>HathiTrust Submissions</h2>
   <div class="row-between bg-white border-b border-t border-brand-grey-100 py-3 px-2 items-center">
      <div>{{ hathiTrust.total }} Submissions</div>
      <div class="row-left">
         <UPagination color="neutral" variant="ghost"
            v-model:page="hathiTrust.searchOpts.currPage" :items-per-page="hathiTrust.searchOpts.limit" 
            :total="hathiTrust.total" @update:page="pageChanged"
         />
         <USelect v-model="hathiTrust.searchOpts.limit" :items="[30,50,100]" @update:modelValue="perPageChanged" />
      </div>
      <UInput v-model="hathiTrust.searchOpts.query" placeholder="Search submissions..." @update:modelValue="hathiTrust.getSubmissions(false)"/>
   </div>
   <UTable :data="hathiTrust.submissions" :columns="columns" v-model:sorting="sorting" sticky virtualize class="h-[calc(100vh-185px)]">
      <template #pid-cell="{ row }">
         <router-link :to="`/metadata/${row.original.metadataID}`">{{row.original.metadata.pid}}</router-link>
      </template>
       <template #requestedAt-header="{ column }">
         <UButton color="neutral" variant="ghost" :icon="sortIcon(column)" @click="sortClicked(column)"> 
            <span class="text-start">{{ column.columnDef.header }}</span>
         </UButton>
      </template>
      <template #requestedAt-cell="{ row }">
         <span class="text-nowrap">{{ toCalendarDate(parseAbsolute(row.original.requestedAt)).toString() }}</span>
      </template>
      <template #metadataSubmittedAt-header="{ column }">
         <UButton color="neutral" variant="ghost" :icon="sortIcon(column)" @click="sortClicked(column)"> 
            <span class="text-start">{{ column.columnDef.header }}</span>
         </UButton>
      </template>
      <template #metadataSubmittedAt-cell="{ row }">
         <span v-if="row.original.metadataSubmittedAt" class="text-nowrap">{{ toCalendarDate(parseAbsolute(row.original.metadataSubmittedAt)).toString() }}</span>
      </template>
      <template #metadataStatus-header>
         <FilterPopover label="Metadata Status"  :applied="isFilterApplied('metadataStatus')" @clear="clearFilter('metadataStatus')" @apply="getSubmissions">
            <USelect v-model="columnFilters.metadataStatus.value" :items="statuses" placeholder="Status..." />
         </FilterPopover>
      </template>
      <template #packageCreatedAt-header="{ column }">
         <UButton color="neutral" variant="ghost" :icon="sortIcon(column)" @click="sortClicked(column)"> 
            <span class="text-start">{{ column.columnDef.header }}</span>
         </UButton>
      </template>
      <template #packageCreatedAt-cell="{ row }">
         <span v-if="row.original.packageCreatedAt" class="text-nowrap">{{ toCalendarDate(parseAbsolute(row.original.packageCreatedAt)).toString() }}</span>
      </template>
      <template #packageSubmittedAt-header="{ column }">
         <UButton color="neutral" variant="ghost" :icon="sortIcon(column)" @click="sortClicked(column)"> 
            <span class="text-start">{{ column.columnDef.header }}</span>
         </UButton>
      </template>
      <template #packageSubmittedAt-cell="{ row }">
         <span v-if="row.original.packageSubmittedAt" class="text-nowrap">{{ toCalendarDate(parseAbsolute(row.original.packageSubmittedAt)).toString() }}</span>
      </template>
      <template #packageStatus-header>
         <FilterPopover label="Package Status"  :applied="isFilterApplied('packageStatus')" @clear="clearFilter('packageStatus')" @apply="getSubmissions">
            <USelect v-model="columnFilters.packageStatus.value" :items="statuses" placeholder="Status..." />
         </FilterPopover>
      </template>
      <template #finishedAt-header="{ column }">
         <UButton color="neutral" variant="ghost" :icon="sortIcon(column)" @click="sortClicked(column)"> 
            <span class="text-start">{{ column.columnDef.header }}</span>
         </UButton>
      </template>
      <template #finishedAt-cell="{ row }">
         <span v-if="row.original.finishedAt" class="text-nowrap">{{ toCalendarDate(parseAbsolute(row.original.finishedAt)).toString() }}</span>
      </template>
   </UTable>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useHathiTrustStore } from '@/stores/hathitrust'
import { toCalendarDate, parseAbsolute } from '@internationalized/date'

const hathiTrust = useHathiTrustStore()

const columns = [
   { accessorKey: 'metadata.pid', id: 'pid', header: 'PID'},
   { accessorKey: 'metadata.title', id: 'title', header: 'Title'},
   { accessorKey: 'metadata.barcode', id: 'barcode', header: 'Barcode'},
   { accessorKey: 'requestedAt', header: 'Requested'},
   { accessorKey: 'metadataSubmittedAt', header: 'Metadata Submitted'},
   { accessorKey: 'metadataStatus', header: 'Metadata Status'},
   { accessorKey: 'packageCreatedAt', header: 'Package Created'},
   { accessorKey: 'packageSubmittedAt', header: 'Package Submitted'},
   { accessorKey: 'packageStatus', header: 'Package Status Status'},
   { accessorKey: 'finishedAt', header: 'Finished'}
]
const sorting = ref([{id: 'requestedAt', desc: true}])
const columnFilters = ref( {
   metadataStatus: {value: null, mode: "equals"},
   packageStatus: {value: null, mode: "equals"},
})

const statuses = ref([
   {label: "Pending", value: "pending"},
   {label: "Submitted", value: "submitted"},
   {label: "Accepted", value: "accepted"},
   {label: "Failed", value: "failed"},
])

onMounted(() => {
   getSubmissions()
   document.title = `HathiTrust Submissions`
})

const sortIcon = ((column) => {
   const isSorted = column.getIsSorted()
   if (isSorted) {
      if (isSorted === 'asc') return 'i-lucide-arrow-up-narrow-wide'
      return 'i-lucide-arrow-down-wide-narrow'
   }
   return 'i-lucide-arrow-up-down'
})
const sortClicked = (async (column) => {
   const isSorted = column.getIsSorted()
   hathiTrust.searchOpts.sortField = column.id
   hathiTrust.searchOpts.sortOrder = "asc"
   if (isSorted === "asc") {
      hathiTrust.searchOpts.sortOrder = "desc"   
   }
   hathiTrust.resetSearch()
   await hathiTrust.getSubmissions( )
   column.toggleSorting(column.getIsSorted() === 'asc')
})

const isFilterApplied = ((name) => {
   return columnFilters.value[name].value != null 
})
const clearFilter = ((name) => {
   columnFilters.value[name].value = null 
   getSubmissions()
}) 
const getSubmissions = (() => {
   hathiTrust.searchOpts.filters = []
   Object.entries(columnFilters.value).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         hathiTrust.searchOpts.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   hathiTrust.getSubmissions()
})

const perPageChanged = (() => {
   hathiTrust.searchOpts.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   hathiTrust.searchOpts.start = (hathiTrust.searchOpts.currPage-1) * hathiTrust.searchOpts.limit
   hathiTrust.getSubmissions()
})

// const onPage = ((event) => {
//    hathiTrust.searchOpts.start = event.first
//    hathiTrust.searchOpts.limit = event.rows
//    getSubmissions()
// })

// const onSort = ((event) => {
//    hathiTrust.searchOpts.sortField = event.sortField
//    hathiTrust.searchOpts.sortOrder = "asc"
//    if (event.sortOrder == -1) {
//       hathiTrust.searchOpts.sortOrder = "desc"
//    }
//    getSubmissions()
// })

</script>

<style scoped lang="scss">
:deep(th.p-sortable-column)  {
   white-space: break-spaces !important;
}
.submissions {
   min-height: 600px;
   text-align: left;
   padding: 0 25px;

   button.pad {
      margin-left: 10px;
   }
   button.notes {
      font-size: 0.75em;
      padding: 5px 10px;
   }
   .none {
      color: var(--uvalib-grey-light);
      font-style: italic;
   }
}
div.note-text {
   height: 150px;
   overflow-y: scroll;
   margin: 15px 0 10px 0;
   padding: 5px 10px;
   border: 1px solid var(--uvalib-grey-light);
   border-radius: 5px;
   background-color: white;
}
</style>