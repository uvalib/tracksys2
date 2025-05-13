<template>
   <DataTable :value="searchStore.orders.hits" ref="orderHitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
      :lazy="true" :paginator="true" @page="onPage($event)"
      :rows="searchStore.orders.limit" :totalRecords="searchStore.orders.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[15,30,100]" :first="searchStore.orders.start"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"  paginatorPosition="top"
   >
      <template #empty><h3>No matching orders found</h3></template>
      <template #paginatorstart>
         <div class="acts">
            <DPGButton label="Download Results CSV" severity="secondary" @click="downloadCSV" v-if="searchStore.orders.total>0" />
            <DPGButton v-if="hasFilter" label="Clear All Filters" severity="secondary" @click="clearFilters"/>
         </div>
      </template>
      <Column field="id" header="ID">
         <template #body="slotProps">
            <router-link :to="`/orders/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="status" header="Status" class="nowrap" filterField="status" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Select v-model="filterModel.value" :options="orderStatuses" optionLabel="name" optionValue="code" placeholder="Select a status" />
         </template>
         <template #body="slotProps">
            <span :class="`status ${slotProps.data.status}`">{{displayStatus(slotProps.data.status)}}</span>
         </template>
      </Column>
      <Column field="customer" header="Customer" class="nowrap" filterField="customer" :showFilterMatchModes="false">
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Customer"/>
         </template>
      </Column>
      <Column field="agency" header="Agency" class="nowrap"  filterField="agency" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Agency name"/>
         </template>
      </Column>
      <Column field="title" header="Order Title" filterField="title" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
         </template>
      </Column>
      <Column field="staff_notes" header="Staff Notes" filterField="staff_notes" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Notes"/>
         </template>
      </Column>
      <Column field="special_instructions" header="Special Instructions" filterField="special_instructions" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { FilterMatchMode } from '@primevue/core/api'
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import { useRoute, useRouter } from 'vue-router'
import { usePinnable } from '@/composables/pin'

usePinnable("p-datatable-paginator-top")

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const orderHitsTable = ref()

const filters = ref( {
   'status': {value: null, matchMode: FilterMatchMode.EQUALS},
   'customer': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'agency': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'staff_notes': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'special_instructions': {value: null, matchMode: FilterMatchMode.CONTAINS},
})

const orderStatuses = ref([
   {name: "Requested", code: "requested"},
   {name: "Approved", code: "approved"},
   {name: "Await Fee", code: "await_fee"},
   {name: "Completed", code: "completed"},
   {name: "Canceled", code: "canceled"},
   {name: "Deferred", code: "deferred"},
])

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

<style scoped lang="scss">
.acts{
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: center;
   gap: 10px;
}
</style>