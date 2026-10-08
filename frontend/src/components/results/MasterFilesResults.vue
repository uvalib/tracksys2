<template>
   <div class="row-between p-2 sticky z-50 bg-white border-b border-brand-grey-100" :style="{top: headerHeight}">
      <div class="row-left gap-2!">
         <UButton color="secondary" label="Download Results CSV" @click="downloadCSV"/>
         <UButton v-if="hasFilter" color="secondary" label="Clear All Filters" @click="clearFilters"/>
      </div>
      <div class="row-right">
         <UPagination color="neutral" variant="ghost"
            v-model:page="searchStore.masterFiles.currPage" :items-per-page="searchStore.masterFiles.limit" 
            :total="searchStore.masterFiles.total" @update:page="pageChanged"
         />
         <USelect v-model="searchStore.masterFiles.limit" :items="[15,30,100]" @change="perPageChanged" />
      </div>
   </div>
   <UTable :data="searchStore.masterFiles.hits" :columns="columns">
      <template #id-cell="{ row }">
         <router-link :to="`/masterfiles/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
      <template #unit_id-cell="{ row }">
         <router-link :to="`/units/${row.original.id}`">{{row.original.unit_id}}</router-link>
      </template>
      <template #is_clone-header>
         <FilterPopover label="Clone" :applied="filterApplied('is_clone')" @clear="clearFilter('is_clone')" @apply="applyFilter()">
            <USelect v-model="filters.is_clone.value" :items="yesNo" placeholder="Select a type" />
         </FilterPopover>
      </template>
      <template #is_clone-cell="{ row }">
         <span v-if="row.original.is_clone > 0">Yes</span>
         <span v-else>No</span>
      </template>
      <template #call_number-header>
         <FilterPopover label="Call Number"  :applied="filterApplied('call_number')" @clear="clearFilter('call_number')" @apply="applyFilter()">
            <UInput v-model="filters.call_number.value" placeholder="Call Number..." />
         </FilterPopover>
      </template>
      <template #call_number-cell="{ row }">
         <router-link :to="`/metadata/${row.original.metadata_id}`">{{row.original.call_number}}</router-link>
      </template>
      <template #title-header>
         <FilterPopover label="Title"  :applied="filterApplied('title')" @clear="clearFilter('title')" @apply="applyFilter()">
            <UInput v-model="filters.title.value" placeholder="Title..." />
         </FilterPopover>
      </template>
      <template #description-header>
         <FilterPopover label="Description"  :applied="filterApplied('description')" @clear="clearFilter('description')" @apply="applyFilter()">
            <UInput v-model="filters.description.value" placeholder="Description..." />
         </FilterPopover>
      </template>
      <template #thumbnail_url-cell="{ row }">
         <a :href="row.original.image_url" target="_blank">
            <img :src="row.original.thumbnail_url" />
         </a>
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
      header: 'ID',
      meta: {
         class: {
            td: 'w-fit'
         }
      }
   },
   {
      accessorKey: 'pid',
      header: 'PID',
      meta: {
         class: {
            td: 'w-fit'
         }
      }
   },
   {
      accessorKey: 'unit_id',
      header: 'Unit',
      meta: {
         class: {
            td: 'w-fit'
         }
      }
      
   },
   {
      accessorKey: 'is_clone',
      header: 'Clone',
      meta: {
         class: {
            td: 'w-fit'
         }
      }
      
   },
   {
      accessorKey: 'call_number',
      header: 'Call Number',
      meta: {
         class: {
            td: 'w-fit'
         }
      }
   },
   {
      accessorKey: 'filename',
      header: 'Filename',
      meta: {
         class: {
            td: 'w-fit'
         }
      }
   },
   {
      accessorKey: 'title',
      header: 'Title',
      meta: {
         class: {
            td: 'w-1/5 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'description',
      header: 'Description',
      meta: {
         class: {
            td: 'w-1/5 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'thumbnail_url',
      header: 'Thumb',
      meta: {
         class: {
            td: 'w-fit'
         }
      }
   },
]

const filters = ref({
   is_clone: {value: null, mode: "equals"},
   title: {value: null, mode: "contains"},
   description: {value: null, mode: "contains"},
   call_number: {value: null, mode: "startsWith"},
})

const yesNo = computed(() => {
   let out = []
   out.push( {label: "No", value: "false"} )
   out.push( {label: "Yes", value: "true"} )
   return out
})

const hasFilter = computed(() => {
   let idx = Object.values(filters.value).findIndex( fv => fv.value && fv.value != "")
   return idx >= 0
})

onMounted(() =>{
   searchStore.masterFiles.filters.forEach( fv => {
      filters.value[fv.field].value = fv.value
   })
})

const downloadCSV = (() => {
   searchStore.downloadCSV('mastefiles', columns )
})

const clearFilters = (() => {
   filters.value = {
      is_clone: {value: null, mode: "equals"},
      title: {value: null, mode: "contains"},
      description: {value: null, mode: "contains"},
      call_number: {value: null, mode: "startsWith"},
   }
   searchStore.masterFiles.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.executeSearch("masterfiles")
})

const filterApplied = ((name) => {
   return filters.value[name].value != null 
})
const clearFilter = ((name) => {
   filters.value[name].value = null 
   applyFilter()
}) 

const applyFilter =(() => {
   searchStore.masterFiles.filters = []
   Object.entries(filters.value).forEach(([filterName, data]) => {
      if (data.value && data.value != "") {
         searchStore.masterFiles.filters.push({field: filterName, match: data.mode, value: data.value})
      }
   })
   let query = Object.assign({}, route.query)
   query.filters = searchStore.filtersAsQueryParam("masterfiles")
   delete query.filters
   if ( searchStore.masterFiles.filters.length > 0) {
      query.filters = searchStore.filtersAsQueryParam("masterfiles")
   }
   router.push({query})
   searchStore.executeSearch("masterfiles")
})

const perPageChanged = (() => {
   searchStore.masterFiles.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   searchStore.masterFiles.start = (searchStore.masterFiles.currPage-1) * searchStore.masterFiles.limit
   searchStore.executeSearch("masterfiles")
})
</script>

<style scoped lang="scss">
</style>