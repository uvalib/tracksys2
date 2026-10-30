<template>
   <h2>{{ title }}</h2>
   <div class="row-between gap-2! items-center bg-white border-b border-brand-grey-100">
      <span class="px-4">Total: {{ publishedStore.total }}</span>
      <div class="p-4 flex flex-row flex-nowrap justify-center items-center">
         <UPagination color="neutral" variant="ghost" show-edges
            v-model:page="publishedStore.searchOpts.currPage" :items-per-page="publishedStore.searchOpts.limit" 
            :total="publishedStore.total" @update:page="pageChanged"
         />
         <USelect v-model="publishedStore.searchOpts.limit" :items="[25,50,100]" @change="perPageChanged" />
      </div>
   </div>
   <UTable :data="publishedStore.records" :columns="columns" sticky class="h-[calc(100vh-140px)]" >
      <template #pid-cell="{ row }">
         <router-link :to="`/metadata/${row.original.id}`">{{row.original.pid}}</router-link>
      </template>
      <template #type-header="{ column }">
         <div  v-if="route.params.type == 'archivesspace'">{{ column.columnDef.header }}</div>
         <FilterPopover v-else :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column.id)"
            @clear="clearFilter(column.id)" @apply="setFilter(column.id)">
            <USelect v-model="pendingFilter" :items="mdTypes" placeholder="Type..." />
         </FilterPopover>
      </template>
      <template #type-cell="{ row }">
         <div v-if="row.original.type != 'ExternalMetadata'">{{row.original.type}}</div>
         <div v-else>ArchivesSpace</div>
      </template>
      <template #title-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column.id)"
            @clear="clearFilter(column.id)" @apply="setFilter(column.id)">
            <UInput v-model="pendingFilter" :placeholder="`${column.columnDef.header}...`" />
         </FilterPopover>
      </template>
      <template #creatorName-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column.id)"
            @clear="clearFilter(column.id)" @apply="setFilter(column.id)">
            <UInput v-model="pendingFilter" :placeholder="`${column.columnDef.header}...`" />
         </FilterPopover>
      </template>
      <template v-if="route.params.type == 'virgo'" #barcode-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column.id)"
            @clear="clearFilter(column.id)" @apply="setFilter(column.id)">
            <UInput v-model="pendingFilter" :placeholder="`${column.columnDef.header}...`" />
         </FilterPopover>
      </template>
      <template v-if="route.params.type == 'virgo'" #callNumber-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column.id)"
            @clear="clearFilter(column.id)" @apply="setFilter(column.id)">
            <UInput v-model="pendingFilter" :placeholder="`${column.columnDef.header}...`" />
         </FilterPopover>
      </template>
      <template v-if="route.params.type == 'virgo'" #catalogKey-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @opened="filterOpened(column.id)"
            @clear="clearFilter(column.id)" @apply="setFilter(column.id)">
            <UInput v-model="pendingFilter" :placeholder="`${column.columnDef.header}...`" />
         </FilterPopover>
      </template>
   </UTable>
</template>

<script setup>
import { onBeforeMount, onMounted, ref, computed } from 'vue'
import { usePublishedStore } from '@/stores/published'
import { useRoute, useRouter, onBeforeRouteUpdate } from 'vue-router'

const route = useRoute()
const router = useRouter()
const publishedStore = usePublishedStore()
const pendingFilter = ref()
const filters = ref([])

const columns = computed(() => {
   let out = [
      { accessorKey: 'pid', header: 'PID'}, { accessorKey: 'type', header: 'Type'},
      { accessorKey: 'title', header: 'Title'}, { accessorKey: 'creatorName', header: 'Creator Name'}
   ]
   if ( route.params.type == 'virgo' ) {
      out.push( { accessorKey: 'barcode', header: 'Barcode'} )  
      out.push( { accessorKey: 'callNumber', header: 'Call Number'} )   
      out.push( { accessorKey: 'catalogKey', header: 'Catalog Key'} )  
   }
   return out
})

const isFilterApplied = ((name) => {
   return filters.value.findIndex( f => f.id == name) > -1
})
const filterOpened = ((name) => {
   pendingFilter.value = null
   const existing = filters.value.find( f => f.id == name)
   if (existing ) {
      pendingFilter.value = existing.value
   }
})
const clearFilter = ((name) => {
   const idx = filters.value.findIndex( f => f.id == name) 
   if ( idx > -1 ) {
      filters.value.splice(idx,1)
   }
   applyFilter()
})
const setFilter = ((name) => {    
   filters.value.push({id: name, value: pendingFilter.value})
   applyFilter()
})
const applyFilter = (() => {    
   publishedStore.searchOpts.filters = []
   filters.value.forEach( (filter) => {
      if ( filter.id == "type") {
         publishedStore.searchOpts.filters.push({field: filter.id, match: "equals", value: filter.value})
      } else {
         publishedStore.searchOpts.filters.push({field: filter.id, match: "contains", value: filter.value})
      }
   })

   let query = Object.assign({}, route.query)
   if ( publishedStore.searchOpts.filters.length > 0 ) {
      query.filters = publishedStore.filtersAsQueryParam
   } else {
      delete query.filters
   }
   router.push({query})
   publishedStore.resetSearch()
   publishedStore.getRecords(route.params.type)
})

const mdTypes = computed(() => {
   return [ {label: "Sirsi", value: "SirsiMetadata"}, {label: "XML", value: "XmlMetadata"} ]
})

const title = computed(() => {
   if (route.params.type == "archivesspace") {
      return "Items Published to ArchivesSpace"
   }
   if (route.params.type == "dpla") {
      return "Items Published to DPLA"
   }
   return "Items Published to Virgo"
})

onBeforeRouteUpdate( (to) => {
   publishedStore.searchOpts.filters = []
   publishedStore.getRecords(to.params.type)
})

onBeforeMount( () => {
   publishedStore.searchOpts.filters = []
   if ( route.query.filters ) {
      let parsedFilters = JSON.parse(route.query.filters)
      filters.value = []
      parsedFilters.forEach( filter => {
         let bits = filter.split("|")
         publishedStore.searchOpts.filters.push( {field: bits[0], match: bits[1], value: bits[2]} )
         filters.value.push({id: bits[0], value: bits[2]})
      })
   }
})

onMounted(() => {
   publishedStore.getRecords(route.params.type)
   document.title = `Published`
})

const perPageChanged = (() => {
   publishedStore.searchOpts.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   publishedStore.searchOpts.start = (publishedStore.searchOpts.currPage-1) * publishedStore.searchOpts.limit
   publishedStore.getRecords(route.params.type)
})
</script>

<style scoped lang="scss">
</style>