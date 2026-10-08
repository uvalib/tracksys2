<template>
   <ResultsToolbar :count="searchStore.components.hits.length" :total="searchStore.components.total" :hasFilter="hasFilter"
      @more="loadMore" @csv="downloadCSV" @clear-filter="clearFilters"
   />
   <UTable :data="searchStore.components.hits" :columns="columns">
      <template #id-cell="{ row }">
         <router-link :to="`/components/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
      <template #title-header>
         <FilterPopover label="Title"  :applied="filterApplied('title')" @clear="clearFilter('title')" @apply="applyFilter()">
            <UInput v-model="filters.title.value" placeholder="Title..." />
         </FilterPopover>
      </template>
      <template #label-header>
         <FilterPopover label="Label"  :applied="filterApplied('label')" @clear="clearFilter('label')" @apply="applyFilter()">
            <UInput v-model="filters.label.value" placeholder="Label..." />
         </FilterPopover>
      </template>
      <template #description-header>
         <FilterPopover label="Content Description"  :applied="filterApplied('description')" @clear="clearFilter('description')" @apply="applyFilter()">
            <UInput v-model="filters.description.value" placeholder="Description..." />
         </FilterPopover>
      </template>
      <template #date-header>
         <FilterPopover label="Date"  :applied="filterApplied('date')" @clear="clearFilter('date')" @apply="applyFilter()">
            <UInput v-model="filters.date.value" placeholder="Date..." />
         </FilterPopover>
      </template>
   </UTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSearchStore } from '../../stores/search'
import { useRoute, useRouter } from 'vue-router'
import ResultsToolbar from './ResultsToolbar.vue'

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const columns = [
   {
      accessorKey: 'id',
      header: 'ID'
   },
   {
      accessorKey: 'pid',
      header: 'PID'
   },
   {
      accessorKey: 'title',
      header: 'Title',
       meta: {
         class: {
            td: 'w-1/6 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'label',
      header: 'Label',
       meta: {
         class: {
            td: 'w-1/6 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'description',
      header: 'Content Description',
       meta: {
         class: {
            td: 'w-1/6 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'date',
      header: 'Date'
   },
   {
      accessorKey: 'ead_id',
      header: 'EAD ID'
   },
   {
      accessorKey: 'mf_cnt',
      header: 'Master Files'
   },
]

const filters = ref({
   title: {value: null, mode: "contains"},
   label: {value: null, mode: "contains"},
   description: {value: null, mode: "contains"},
   date: {value: null, mode: "contains"},
})

const hasFilter = computed(() => {
   let idx = Object.values(filters.value).findIndex( fv => fv.value && fv.value != "")
   return idx >= 0
})

onMounted(() =>{
   searchStore.components.filters.forEach( fv => {
      filters.value[fv.field].value = fv.value
   })
})

const downloadCSV = (() => {
   searchStore.downloadCSV('components', columns )
})

const clearFilters = (() => {
   filters.value = {
      title: {value: null, mode: "contains"},
      label: {value: null, mode: "contains"},
      description: {value: null, mode: "contains"},
      date: {value: null, mode: "contains"},
   }
   searchStore.components.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.resetSearch("components")
   searchStore.executeSearch("components")
})

const filterApplied = ((name) => {
   return filters.value[name].value != null 
})
const clearFilter = ((name) => {
   filters.value[name].value = null 
   applyFilter()
}) 

const applyFilter =(() => {
   searchStore.components.filters = []
   Object.entries(filters.value).forEach(([filterName, data]) => {
      if (data.value && data.value != "") {
         searchStore.components.filters.push({field: filterName, match: data.mode, value: data.value})
      }
   })
   let query = Object.assign({}, route.query)
   query.filters = searchStore.filtersAsQueryParam("components")
   delete query.filters
   if ( searchStore.components.filters.length > 0) {
      query.filters = searchStore.filtersAsQueryParam("components")
   }
   router.push({query})
   searchStore.resetSearch("components")
   searchStore.executeSearch("components")
})

const loadMore = (() => {
   searchStore.executeSearch("components")
})

</script>

<style scoped lang="scss">
</style>