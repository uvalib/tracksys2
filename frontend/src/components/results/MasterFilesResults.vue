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
            <DPGButton label="Clear all" class="p-button-secondary" @click="clearFilters()"/>
         </div>
      </div>
   </div>
   <div v-if="searchStore.masterFiles.total == 0">
      <h3>No matching master files found</h3>
   </div>
   <DataTable v-else :value="searchStore.masterFiles.hits" ref="masterFileHitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
      :lazy="true" :paginator="true" @page="onPage($event)"
      :rows="searchStore.masterFiles.limit" :totalRecords="searchStore.masterFiles.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10,30,100]"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <Column field="id" header="ID"/>
      <Column field="pid" header="PID" class="nowrap" />
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
   </DataTable>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import {FilterMatchMode} from 'primevue/api'
import InputText from 'primevue/inputtext'

const searchStore = useSearchStore()

const filters = ref( {
    'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
    'description': {value: null, matchMode: FilterMatchMode.CONTAINS},
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

const hasFilter = computed(() => {
   let idx = Object.values(filters.value).findIndex( fv => fv.value && fv.value != "")
   return idx >= 0
})

function clearFilters() {
   Object.values(filters.value).forEach( fv => fv.value = null )
   searchStore.masterFiles.filters = []
   searchStore.executeSearch("masterfiles")
}

function onFilter(event) {
   searchStore.masterFiles.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         searchStore.masterFiles.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
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
}
</stype>