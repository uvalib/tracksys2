<template>
   <DataTable :value="searchStore.metadata.hits" ref="metadataTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll"
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
      :lazy="true" :paginator="true" @page="onMetadataPage($event)"
      :rows="searchStore.metadata.limit" :totalRecords="searchStore.metadata.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[15,30,100]" :first="searchStore.metadata.start"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"   paginatorPosition="top"
   >
      <template #empty><h3>No matching metadata records found</h3></template>
      <template #paginatorstart>
         <DPGButton label="Download Results CSV" class="p-button-secondary download" @click="downloadCSV" v-if="searchStore.metadata.total>0" />
         <DPGButton v-if="hasFilter" label="Clear All Filters" class="p-button-secondary" @click="clearFilters"/>
      </template>
      <Column field="id" header="ID">
         <template #body="slotProps">
            <router-link :to="`/metadata/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="pid" header="PID" class="nowrap">
         <template #body="slotProps">
            <div>{{slotProps.data.pid}}</div>
            <div v-if="slotProps.data.virgoURL && slotProps.data.type=='XmlMetadata'"><a :href="slotProps.data.virgoURL" target="_blank">VIRGO</a></div>
         </template>
      </Column>
      <Column field="type" header="Type" filterField="type" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="mdTypes" optionLabel="name" optionValue="code" placeholder="Select a type" />
         </template>
         <template #body="slotProps">
            <div v-if="slotProps.data.type != 'ExternalMetadata'">{{slotProps.data.type}}</div>
            <div v-else>{{slotProps.data.externalSystem.name}}</div>
         </template>
      </Column>
      <Column field="title" header="Title" filterField="title" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
         </template>
      </Column>
      <Column field="creatorName" header="Creator Name" filterField="creator_name" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Creator name"/>
         </template>
      </Column>
      <Column field="barcode" header="Barcode" class="nowrap" filterField="barcode" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Barcode"/>
         </template>
      </Column>
      <Column field="callNumber" header="Call Number" class="nowrap" filterField="call_number" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Call number"/>
         </template>
      </Column>
      <Column field="catalogKey" header="Catalog Key" class="nowrap" filterField="catalog_key" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Catalog key"/>
         </template>
         <template #body="slotProps">
            <div>{{slotProps.data.catalogKey}}</div>
            <div v-if="slotProps.data.virgoURL && slotProps.data.catalogKey"><a :href="slotProps.data.virgoURL" target="_blank">VIRGO</a></div>
         </template>
      </Column>
      <Column field="virgo" header="Virgo" class="nowrap" filterField="virgo" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="yesNo" optionLabel="label" optionValue="value" placeholder="Select a value" />
         </template>
         <template #body="slotProps">
            <span v-if="slotProps.data.virgo">Yes</span>
            <span v-else>No</span>
         </template>
      </Column>
      <Column field="dpla" header="DPLA" class="nowrap" filterField="dpla" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="yesNo" optionLabel="label" optionValue="value" placeholder="Select a value" />
         </template>
         <template #body="slotProps">
            <span v-if="slotProps.data.dpla">Yes</span>
            <span v-else>No</span>
         </template>
      </Column>
      <Column field="hathitrust" header="HathiTrust" class="nowrap" filterField="hathitrust" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="yesNo" optionLabel="label" optionValue="value" placeholder="Select a value" />
         </template>
         <template #body="slotProps">
            <span v-if="slotProps.data.hathitrust">Yes</span>
            <span v-else>No</span>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSearchStore } from '../../stores/search'
import { useSystemStore } from '../../stores/system'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import { FilterMatchMode } from '@primevue/core/api'
import { useRoute, useRouter } from 'vue-router'
import { usePinnable } from '@/composables/pin'

usePinnable("p-paginator-top")

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()
const system = useSystemStore()

const metadataTable = ref()

const filters = ref( {
    'type': {value: null, matchMode: FilterMatchMode.EQUALS},
    'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
    'creator_name': {value: null, matchMode: FilterMatchMode.CONTAINS},
    'barcode': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
    'call_number': {value: null, matchMode: FilterMatchMode.CONTAINS},
    'catalog_key': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
    'virgo': {value: null, matchMode: FilterMatchMode.EQUALS},
    'dpla': {value: null, matchMode: FilterMatchMode.EQUALS},
    'hathitrust': {value: null, matchMode: FilterMatchMode.EQUALS}
})

const yesNo = computed(() => {
   let out = []
   out.push( {label: "No", value: "false"} )
   out.push( {label: "Yes", value: "true"} )
   return out
})
const mdTypes = computed(() => {
   let out = []
   out.push({name: "Sirsi", code: "SirsiMetadata"})
   out.push({name: "XML", code: "XmlMetadata"})
   system.externalSystems.forEach( es => {
      out.push({name: es.name, code: `ExternalMetadata:${es.id}`})
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

function downloadCSV() {
   metadataTable.value.exportCSV()
}

function clearFilters() {
   Object.values(filters.value).forEach( fv => fv.value = null )
   searchStore.metadata.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.executeSearch("metadata")
}

function onMetadataPage(event) {
   searchStore.metadata.start = event.first
   searchStore.metadata.limit = event.rows
   searchStore.executeSearch("metadata")
}

function onFilter(event) {
   searchStore.metadata.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         searchStore.metadata.filters.push({field: key, match: data.matchMode, value: data.value})
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
}

</script>

<stype scoped lang="scss">
</stype>