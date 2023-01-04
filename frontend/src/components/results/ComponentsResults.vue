<template>
   <div v-if="hasFilter" class="filters">
      <div class="filter-head">Filters</div>
      <div class="content">
            <ul>
               <li v-for="(vf,idx) in selectedFilters" :key="`component-filter=${idx}`">
                  <label>{{vf.filter}}:</label>
                  <span>{{vf.value}}</span>
               </li>
            </ul>
         <div class="filter-acts">
            <DPGButton label="Clear all" class="p-button-secondary" @click="clearFilters()"/>
         </div>
      </div>
   </div>
   <div v-if="searchStore.components.total == 0">
      <h3>No matching components found</h3>
   </div>
   <DataTable v-else :value="searchStore.components.hits" ref="componentHitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
      :lazy="true" :paginator="searchStore.components.total > 15" @page="onPage($event)"
      :rows="searchStore.components.limit" :totalRecords="searchStore.components.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[15,30,100]"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <template #header>
         <div class="results-toolbar">
            <div class="matches">{{searchStore.components.total}} matches found</div>
            <DPGButton label="Download Results CSV" class="p-button-secondary" @click="downloadCSV"/>
         </div>
      </template>
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
      <Column field="description" header="Content Description" filterField="content_desc" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Description"/>
         </template>
      </Column>
      <Column field="date" header="Date" class="nowrap" filterField="date" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Date"/>
         </template>
      </Column>
      <Column field="eadID" header="EAD ID" />
      <Column field="masterFileCount" header="Master Files"/>
      <Column header="" class="row-acts nowrap">
         <template #body="slotProps">
            <router-link :to="`/components/${slotProps.data.id}`">View</router-link>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { FilterMatchMode } from 'primevue/api'
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const componentHitsTable = ref()

const filters = ref( {
   'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'label': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'content_desc': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'date': {value: null, matchMode: FilterMatchMode.CONTAINS},
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
   query.scope = "components"
   router.push({query})
   searchStore.executeSearch("components")
}

function onPage(event) {
   searchStore.components.start = event.first
   searchStore.components.limit = event.rows
   searchStore.executeSearch("components")
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
}
</stype>