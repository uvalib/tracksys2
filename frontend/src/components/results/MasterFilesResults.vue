<template>
   <DataTable :value="searchStore.masterFiles.hits" ref="masterFileHitsTable" dataKey="id"
      stripedRows showGridlines size="small" :rowStyle="rowStyle"
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
      :lazy="true" :paginator="true" @page="onPage($event)" paginatorPosition="top"
      :rows="searchStore.masterFiles.limit" :totalRecords="searchStore.masterFiles.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[15,30,100]" :first="searchStore.masterFiles.start"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <template #empty><h3>No matching master files found</h3></template>
      <template #paginatorstart>
         <div class="acts">
            <DPGButton label="Download Results CSV" severity="secondary" @click="downloadCSV" v-if="searchStore.masterFiles.total>0" />
            <DPGButton v-if="hasFilter" label="Clear All Filters" severity="secondary" @click="clearFilters"/>
         </div>
      </template>
      <Column field="id" header="ID">
         <template #body="slotProps">
            <router-link :to="`/masterfiles/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="pid" header="PID" class="nowrap" />
      <Column field="unit_id" header="Unit" class="nowrap" >
         <template #body="slotProps">
            <router-link :to="`/units/${slotProps.data.unit_id}`">{{slotProps.data.unit_id}}</router-link>
         </template>
      </Column>
      <Column field="is_clone" header="Clone" class="nowrap" filterField="is_clone" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Select v-model="filterModel.value" :options="yesNo" optionLabel="label" optionValue="value" placeholder="Select a value" />
         </template>
         <template #body="slotProps">
            <span v-if="slotProps.data.is_clone > 0">Yes</span>
            <span v-else>No</span>
         </template>
      </Column>
     <Column field="call_number" header="Call Number" class="nowrap" filterField="call_number" :showFilterMatchModes="false">
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Call Number"/>
         </template>
         <template #body="slotProps">
            <router-link :to="`/metadata/${slotProps.data.metadata_id}`">{{slotProps.data.call_number}}</router-link>
         </template>
      </Column>
      <Column field="filename" header="Filename"/>
      <Column field="title" header="Title" filterField="title" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
         </template>
      </Column>
      <Column field="description" header="Description" filterField="description" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Description"/>
         </template>
      </Column>
      <Column field="thumbnail_url" header="Thumb">
         <template #body="slotProps">
            <a :href="slotProps.data.image_url" target="_blank">
               <img :src="slotProps.data.thumbnail_url" />
            </a>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { FilterMatchMode } from '@primevue/core/api'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import { useRoute, useRouter } from 'vue-router'
import { usePinnable } from '@/composables/pin'

usePinnable("p-datatable-paginator-top")

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const masterFileHitsTable = ref()

const filters = ref( {
   'is_clone': {value: null, matchMode: FilterMatchMode.EQUALS},
   'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'description': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'call_number': {value: null, matchMode: FilterMatchMode.CONTAINS},
})

const yesNo = computed(() => {
   let out = []
   out.push( {label: "No", value: "false"} )
   out.push( {label: "Yes", value: "true"} )
   return out
})

const rowStyle = (data) => {
    if (data.originalID) {
        return { background: '#f5f5f5' };
    }
}

const hasFilter = computed(() => {
   let idx = Object.values(filters.value).findIndex( fv => fv.value && fv.value != "")
   return idx >= 0
})

onMounted(() =>{
   searchStore.masterFiles.filters.forEach( fv => {
      filters.value[fv.field].value = fv.value
   })
})

function downloadCSV() {
   masterFileHitsTable.value.exportCSV()
}

function clearFilters() {
   Object.values(filters.value).forEach( fv => fv.value = null )
   searchStore.masterFiles.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.executeSearch("masterfiles")
}

function onFilter(event) {
   searchStore.masterFiles.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         searchStore.masterFiles.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   let query = Object.assign({}, route.query)
   delete query.filters
   if ( searchStore.masterFiles.filters.length > 0) {
      query.filters = searchStore.filtersAsQueryParam("masterfiles")
   }
   router.push({query})
   searchStore.executeSearch("masterfiles")
}

function onPage(event) {
   searchStore.masterFiles.start = event.first
   searchStore.masterFiles.limit = event.rows
   searchStore.executeSearch("masterfiles")
}

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