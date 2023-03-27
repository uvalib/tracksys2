<template>
   <div v-if="hasFilter" class="filters">
      <div class="filter-head">Filters</div>
      <div class="content">
            <ul>
               <li v-for="(vf,idx) in selectedFilters" :key="`mf-filter=${idx}`">
                  <label>{{vf.filter}}:</label>
                  <span>{{vf.value}}</span>
               </li>
            </ul>
         <div class="filter-acts">
            <DPGButton label="Clear all" class="p-button-secondary" @click="clearFilters"/>
         </div>
      </div>
   </div>
   <div v-if="searchStore.masterFiles.total == 0">
      <h3>No matching master files found</h3>
   </div>
   <DataTable v-else :value="searchStore.masterFiles.hits" ref="masterFileHitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm" :rowStyle="rowStyle"
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
      :lazy="true" :paginator="searchStore.masterFiles.total > 15" @page="onPage($event)"
      :rows="searchStore.masterFiles.limit" :totalRecords="searchStore.masterFiles.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[15,30,100]" :first="searchStore.masterFiles.start"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <template #header>
         <div class="results-toolbar">
            <div class="matches">{{searchStore.masterFiles.total}} matches found</div>
            <DPGButton label="Download Results CSV" class="p-button-secondary" @click="downloadCSV"/>
         </div>
      </template>
      <Column field="id" header="ID">
         <template #body="slotProps">
            <router-link :to="`/masterfiles/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="pid" header="PID" class="nowrap" />
      <Column field="unitID" header="Unit ID" class="nowrap" filterField="unit_id" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Unit ID"/>
         </template>
         <template #body="slotProps">
            <router-link :to="`/units/${slotProps.data.unitID}`">{{slotProps.data.unitID}}</router-link>
         </template>
      </Column>
      <Column field="originalID" header="Clone" class="nowrap" filterField="clone" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="yesNo" optionLabel="label" optionValue="value" placeholder="Select a value" />
         </template>
         <template #body="slotProps">
            <span v-if="slotProps.data.originalID > 0">Yes</span>
            <span v-else>No</span>
         </template>
      </Column>
      <Column field="metadata.callNumber" header="Call Number" class="nowrap" filterField="call_number" :showFilterMatchModes="false">
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Call Number"/>
         </template>
         <template #body="slotProps">
            <router-link :to="`/metadata/${slotProps.data.metadata.id}`">{{slotProps.data.metadata.callNumber}}</router-link>
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
      <Column field="thumbnailURL" header="Thumb">
         <template #body="slotProps">
            <a :href="slotProps.data.imageURL" target="_blank">
               <img :src="slotProps.data.thumbnailURL" />
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
import {FilterMatchMode} from 'primevue/api'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const masterFileHitsTable = ref()

const filters = ref( {
   'clone': {value: null, matchMode: FilterMatchMode.EQUALS},
   'unit_id': {value: null, matchMode: FilterMatchMode.EQUALS},
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

const selectedFilters = computed(() => {
   let out = []
   Object.entries(filters.value).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         out.push( {filter: key, value: data.value})
      }
   })
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
   console.log(filters.value)
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

<stype scoped lang="scss">
.results {
   margin: 20px;
   font-size: 0.9em;
   td.nowrap, th {
      white-space: nowrap;
   }
   th, td {
      font-size: 0.85em;
   }
   .results-toolbar {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      .matches {
         padding: 5px 0;
         text-align: left;
      }
      button {
         font-size: 0.8em;
      }
   }
}
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
   :deep(td.nowrap) {
      white-space: nowrap;
   }
   :deep(.row-acts) {
      text-align: center;
      padding: 0;
      a {
         display: inline-block;
         margin: 0;
         padding: 5px 10px;
      };
   }
}
</stype>