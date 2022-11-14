<template>
   <div v-if="hasFilter" class="filters">
      <div class="filter-head">Filters</div>
      <div class="content">
            <ul>
               <li v-for="(vf,idx) in selectedFilters" :key="`mdfilter=${idx}`">
                  <label>{{vf.filter}}:</label>
                  <span>{{vf.value}}</span>
               </li>
            </ul>
         <div class="filter-acts">
            <DPGButton label="Clear all" class="p-button-secondary" @click="clearFilters()"/>
         </div>
      </div>
   </div>
   <div v-if="searchStore.metadata.total == 0">
      <h3>No matching metadata records found</h3>
   </div>
   <DataTable v-else :value="searchStore.metadata.hits" ref="metadataTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll"
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
      :lazy="true" :paginator="searchStore.metadata.hits.length > 15" @page="onMetadataPage($event)"
      :rows="searchStore.metadata.limit" :totalRecords="searchStore.metadata.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[15,30,100]"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <template #header>
         <div class="matches">{{searchStore.metadata.total}} matches found</div>
      </template>
      <Column field="id" header="ID">
         <template #body="slotProps">
            <router-link :to="`/metadata/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="pid" header="PID" class="nowrap"/>
      <Column field="type" header="Type" filterField="type" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="mdTypes" optionLabel="name" optionValue="code" placeholder="Select a type" />
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
            <div v-if="slotProps.data.virgoURL"><a :href="slotProps.data.virgoURL" target="_blank">VIRGO</a></div>
         </template>
      </Column>
      <Column header="" class="row-acts nowrap">
         <template #body="slotProps">
            <router-link :to="`/metadata/${slotProps.data.id}`">View</router-link>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import {FilterMatchMode} from 'primevue/api'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const filters = ref( {
    'type': {value: null, matchMode: FilterMatchMode.EQUALS},
    'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
    'creator_name': {value: null, matchMode: FilterMatchMode.CONTAINS},
    'barcode': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
    'call_number': {value: null, matchMode: FilterMatchMode.CONTAINS},
    'catalog_key': {value: null, matchMode: FilterMatchMode.STARTS_WITH}
})
const mdTypes = ref([
   {name: "SirsiMetadata", code: "SirsiMetadata"},
   {name: "XmlMetadata", code: "XmlMetadata"},
   {name: "ExternalMetadata", code: "ExternalMetadata"},
])

const selectedFilters = computed(() => {
   let out = []
   Object.entries(filters.value).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         out.push( {filter: key, value: data.value})
      }
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
   query.scope = "metadata"
   router.push({query})
   searchStore.executeSearch("metadata")
}

</script>

<stype scoped lang="scss">
div.filters {
   text-align: left;
   border: 1px solid #e9ecef;
   margin-bottom: 15px;
   div.filter-head {
      padding: 5px 10px;
      font-size: 1em;
      background: var(--uvalib-grey-lightest);
      border-bottom: 1px solid #e9ecef;
   }
   ul {
      list-style: none;
      margin: 10px;
      padding: 5px 10px;
      label {
         font-weight: bold;
         display: inline-block;
         margin-right: 10px;
      }
   }
   .content {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
   }
   .filter-acts {
      padding: 10px;
      font-size: 0.85em;
   }
}
.results {
   margin: 20px;
   font-size: 0.9em;
   td.nowrap, th {
      white-space: nowrap;
   }
   th, td {
      font-size: 0.85em;
   }
   .matches {
      padding: 5px 10px;
      text-align: center;
   }
}
</stype>