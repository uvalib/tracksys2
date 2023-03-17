<template>
   <div v-if="hasFilter" class="filters">
      <div class="filter-head">Filters</div>
      <div class="content">
            <ul>
               <li v-for="(vf,idx) in selectedFilters" :key="`order-filter=${idx}`">
                  <label>{{vf.filter}}:</label>
                  <span>{{vf.value}}</span>
               </li>
            </ul>
         <div class="filter-acts">
            <DPGButton label="Clear all" class="p-button-secondary" @click="clearFilters"/>
         </div>
      </div>
   </div>
   <div v-if="searchStore.orders.total == 0">
      <h3>No matching orders found</h3>
   </div>
   <DataTable v-else :value="searchStore.orders.hits" ref="orderHitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
      :lazy="true" :paginator="searchStore.orders.total > 15" @page="onPage($event)"
      :rows="searchStore.orders.limit" :totalRecords="searchStore.orders.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[15,30,100]" :first="searchStore.orders.start"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <template #header>
         <div class="results-toolbar">
            <div class="matches">{{searchStore.orders.total}} matches found</div>
            <DPGButton label="Download Results CSV" class="p-button-secondary" @click="downloadCSV"/>
         </div>
      </template>
      <Column field="id" header="ID">
         <template #body="slotProps">
            <router-link :to="`/orders/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="status" header="Status" class="nowrap" filterField="order_status" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="orderStatuses" optionLabel="name" optionValue="code" placeholder="Select a status" />
         </template>
         <template #body="slotProps">
            <span :class="`status ${slotProps.data.status}`">{{displayStatus(slotProps.data.status)}}</span>
         </template>
      </Column>
      <Column header="Customer" class="nowrap" filterField="last_name" :showFilterMatchModes="false">
         <template #body="slotProps">
            {{slotProps.data.customer.lastName}}, {{slotProps.data.customer.firstName}}
         </template>
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Last name"/>
         </template>
      </Column>
      <Column field="agency.name" header="Agency" class="nowrap"  filterField="agencies.name" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Agency name"/>
         </template>
      </Column>
      <Column field="title" header="Order Title" filterField="order_title" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
         </template>
      </Column>
      <Column field="notes" header="Staff Notes" filterField="orders.staff_notes" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Notes"/>
         </template>
      </Column>
      <Column field="specialInstructions" header="Special Instructions" filterField="orders.special_instructions" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
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
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const orderHitsTable = ref()

const filters = ref( {
   'order_status': {value: null, matchMode: FilterMatchMode.EQUALS},
   'last_name': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'agencies.name': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'order_title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'orders.staff_notes': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'orders.special_instructions': {value: null, matchMode: FilterMatchMode.CONTAINS},
})

const orderStatuses = ref([
   {name: "Requested", code: "requested"},
   {name: "Approved", code: "approved"},
   {name: "Await Fee", code: "await_fee"},
   {name: "Completed", code: "completed"},
   {name: "Canceled", code: "canceled"},
   {name: "Deferred", code: "deferred"},
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
   searchStore.orders.filters.forEach( fv => {
      filters.value[fv.field].value = fv.value
   })
})

function downloadCSV() {
   orderHitsTable.value.exportCSV()
}

function displayStatus( id) {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
}

function clearFilters() {
   Object.values(filters.value).forEach( fv => fv.value = null )
   searchStore.orders.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.executeSearch("orders")
}

function onFilter(event) {
   searchStore.orders.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         searchStore.orders.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   let query = Object.assign({}, route.query)
   query.filters = searchStore.filtersAsQueryParam("orders")
   delete query.filters
   if ( searchStore.orders.filters.length > 0) {
      query.filters = searchStore.filtersAsQueryParam("orders")
   }
   router.push({query})
   searchStore.executeSearch("orders")
}

function onPage(event) {
   searchStore.orders.start = event.first
   searchStore.orders.limit = event.rows
   searchStore.executeSearch("orders")
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