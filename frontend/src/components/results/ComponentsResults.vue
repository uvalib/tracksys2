<template>
   <ResultsToolbar :count="searchStore.components.hits.length" :total="searchStore.components.total" :hasFilter="hasFilter"
      @more="loadMore" @csv="downloadCSV" @clear-filter="clearFilters"
   />
   <DataTable :value="searchStore.components.hits" ref="componentHitsTable" dataKey="id"
      stripedRows showGridlines size="small" :lazy="true" 
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
   >
      <template #empty><h3>No matching components found</h3></template>
      <Column field="id" header="ID">
         <template #body="slotProps">
            <router-link :to="`/components/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="pid" header="PID" class="nowrap"/>
      <Column field="title" header="Title" filterField="title" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
         </template>
      </Column>
      <Column field="label" header="Label" filterField="label" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Label"/>
         </template>
      </Column>
      <Column field="description" header="Content Description" filterField="description" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Description"/>
         </template>
      </Column>
      <Column field="date" header="Date" class="nowrap" filterField="date" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Date"/>
         </template>
      </Column>
      <Column field="finding_aid" header="EAD ID" />
      <Column field="mf_cnt" header="Master Files"/>
   </DataTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { FilterMatchMode } from '@primevue/core/api'
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import { useRoute, useRouter } from 'vue-router'
import { usePinnable } from '@/composables/pin'
import ResultsToolbar from './ResultsToolbar.vue'

usePinnable("results-toolbar")

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const componentHitsTable = ref()

const filters = ref( {
   'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'label': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'description': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'date': {value: null, matchMode: FilterMatchMode.CONTAINS},
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

function downloadCSV() {
   componentHitsTable.value.exportCSV()
}

function clearFilters() {
   Object.values(filters.value).forEach( fv => fv.value = null )
   searchStore.components.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.resetSearch("orders")
   searchStore.executeSearch("components")
}

function onFilter(event) {
   searchStore.components.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         searchStore.components.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   let query = Object.assign({}, route.query)
   query.filters = searchStore.filtersAsQueryParam("components")
   delete query.filters
   if ( searchStore.components.filters.length > 0) {
      query.filters = searchStore.filtersAsQueryParam("components")
   }
   router.push({query})
   searchStore.resetSearch("orders")
   searchStore.executeSearch("components")
}

function loadMore() {
   searchStore.executeSearch("components")
}
</script>

<style scoped lang="scss">
</style>