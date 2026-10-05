<template>
   <div class="row-between p-2 sticky z-50 bg-white border-b-1 border-brand-grey-100" :style="{top: headerHeight}">
      <div class="row-left gap-2!">
         <UButton color="secondary" label="Download Results CSV" @click="downloadCSV"/>
         <UButton v-if="hasFilter" color="secondary" label="Clear All Filters" @click="clearFilters"/>
      </div>
      <div class="row-right">
         <UPagination color="neutral" variant="ghost"
            v-model:page="searchStore.metadata.currPage" :items-per-page="searchStore.metadata.limit" 
            :total="searchStore.metadata.total" @update:page="pageChanged"
         />
         <USelect v-model="searchStore.metadata.limit" :items="[15,30,100]" @change="perPageChanged" />
      </div>
   </div>
   <UTable :data="searchStore.metadata.hits" :columns="columns">
      <template #id-cell="{ row }">
         <router-link :to="`/metadata/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
      <template #system_name-header>
         <FilterPopover label="Type" :applied="filterApplied('system_name')" @clear="clearFilter('system_name')" @apply="applyFilter()">
            <USelect v-model="filters.system_name.value" :items="mdTypes" placeholder="Select a type" />
         </FilterPopover>
      </template>
      <template #title-header>
         <FilterPopover label="Title"  :applied="filterApplied('title')" @clear="clearFilter('title')" @apply="applyFilter()">
            <UInput v-model="filters.title.value" placeholder="Title..." />
         </FilterPopover>
      </template>
      <template #creator_name-header>
         <FilterPopover label="Creator Name"  :applied="filterApplied('creator_name')" @clear="clearFilter('creator_name')" @apply="applyFilter()">
            <UInput v-model="filters.creator_name.value" placeholder="Creator name..." />
         </FilterPopover>
      </template>
      <template #barcode-header>
         <FilterPopover label="Barcode"  :applied="filterApplied('barcode')" @clear="clearFilter('barcode')" @apply="applyFilter()">
            <UInput v-model="filters.barcode.value" placeholder="Barcode..." />
         </FilterPopover>
      </template>
      <template #call_number-header>
         <FilterPopover label="Call Number"  :applied="filterApplied('call_number')" @clear="clearFilter('call_number')" @apply="applyFilter()">
            <UInput v-model="filters.call_number.value" placeholder="Call Number..." />
         </FilterPopover>
      </template>
      <template #catalog_key-header>
         <FilterPopover label="Catalog Key"  :applied="filterApplied('catalog_key')" @clear="clearFilter('catalog_key')" @apply="applyFilter()">
            <UInput v-model="filters.catalog_key.value" placeholder="Catalog Key..." />
         </FilterPopover>
      </template>
      <template #catalog_key-cell="{ row }">
         <div>{{row.original.catalog_key}}</div>
         <div v-if="row.original.virgo_url && row.original.catalog_key"><a :href="row.original.virgo_url" target="_blank">VIRGO</a></div>
      </template>
      <template #virgo-cell="{ row }">
         <span v-if="row.original.virgo">Yes</span>
         <span v-else>No</span>
      </template>
      <template #virgo-header>
         <FilterPopover label="Virgo" :applied="filterApplied('virgo')" @clear="clearFilter('virgo')" @apply="applyFilter()">
            <USelect v-model="filters.virgo.value" :items="yesNo" placeholder="Select a value" />
         </FilterPopover>
      </template>
      <template #dpla-cell="{ row }">
         <span v-if="row.original.dpla">Yes</span>
         <span v-else>No</span>
      </template>
      <template #dpla-header>
         <FilterPopover label="DPLA" :applied="filterApplied('dpla')" @clear="clearFilter('dpla')" @apply="applyFilter()">
            <USelect v-model="filters.dpla.value" :items="yesNo" placeholder="Select a value" />
         </FilterPopover>
      </template>
      <template #hathitrust-cell="{ row }">
         <span v-if="row.original.hathitrust">Yes</span>
         <span v-else>No</span>
      </template>
      <template #hathitrust-header>
         <FilterPopover label="HathiTrust" :applied="filterApplied('hathitrust')" @clear="clearFilter('hathitrust')" @apply="applyFilter()">
            <USelect v-model="filters.hathitrust.value" :items="yesNo" placeholder="Select a value" />
         </FilterPopover>
      </template>
   </UTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSearchStore } from '../../stores/search'
import { useSystemStore } from '../../stores/system'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()
const system = useSystemStore()

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
      accessorKey: 'pid',
      header: 'PID'
   },
   {
      accessorKey: 'system_name',
      header: 'Type',
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
            td: 'w-1/6 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'creator_name',
      header: 'Creator Name',
       meta: {
         class: {
            td: 'w-1/8 whitespace-break-spaces'
         }
      }
   },
   {
      accessorKey: 'barcode',
      header: 'Barcode'
   },
   {
      accessorKey: 'call_number',
      header: 'Call Number'
   },
   {
      accessorKey: 'catalog_key',
      header: 'Catalog Key'
   },
   {
      accessorKey: 'virgo',
      header: 'Virgo'
   },
   {
      accessorKey: 'dpla',
      header: 'DPLA'
   },
   {
      accessorKey: 'hathitrust',
      header: 'HathiTrust'
   },
]

const filters = ref({
   system_name: {value: null, mode: "equals"},
   title: {value: null, mode: "contains"},
   creator_name: {value: null, mode: "contains"},
   barcode: {value: null, mode: "startsWith"},
   call_number: {value: null, mode: "startsWith"},
   catalog_key: {value: null, mode: "startsWith"},
   virgo: {value: null, mode: "equals"},
   dpla: {value: null, mode: "equals"},
   hathitrust: {value: null, mode: "equals"}
})

const yesNo = computed(() => {
   let out = []
   out.push( {label: "No", value: "false"} )
   out.push( {label: "Yes", value: "true"} )
   return out
})
const mdTypes = computed(() => {
   let out = []
   out.push({label: "Sirsi", value: "SirsiMetadata"})
   out.push({label: "XML", value: "XmlMetadata"})
   system.externalSystems.forEach( es => {
      out.push({label: es.name, value: es.name})
   })
   return out
})
const hasFilter = computed(() => {
   let idx = Object.values(filters.value).findIndex( fv => fv.value && fv.value != "")
   return idx >= 0
})

onMounted(() =>{
   searchStore.metadata.filters.forEach( fv => {
      filters.value[fv.field].value = fv.value
   })
})

const downloadCSV = (() => {
   searchStore.downloadCSV('metadata', columns )
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
   searchStore.metadata.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.executeSearch("metadata")
})

const filterApplied = ((name) => {
   return filters.value[name].value != null 
})
const clearFilter = ((name) => {
   filters.value[name].value = null 
   applyFilter()
}) 

const applyFilter =(() => {
   searchStore.metadata.filters = []
   Object.entries(filters.value).forEach(([filterName, data]) => {
      if (data.value && data.value != "") {
         searchStore.metadata.filters.push({field: filterName, match: data.mode, value: data.value})
      }
   })
   let query = Object.assign({}, route.query)
   query.filters = searchStore.filtersAsQueryParam("metadata")
   delete query.filters
   if ( searchStore.metadata.filters.length > 0) {
      query.filters = searchStore.filtersAsQueryParam("metadata")
   }
   router.push({query})
   searchStore.executeSearch("metadata")
})

const perPageChanged = (() => {
   searchStore.metadata.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   searchStore.metadata.start = (searchStore.metadata.currPage-1) * searchStore.metadata.limit
   searchStore.executeSearch("metadata")
})

</script>

<style scoped lang="scss">
</style>