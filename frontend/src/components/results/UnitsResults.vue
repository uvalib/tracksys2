<template>
   <DataTable :value="searchStore.units.hits" ref="unittHitsTable" dataKey="id"
      stripedRows showGridlines size="small"
      v-model:filters="filters" filterDisplay="menu" @filter="onFilter($event)"
      :lazy="true" :paginator="true" @page="onPage($event)" paginatorPosition="top"
      :rows="searchStore.units.limit" :totalRecords="searchStore.units.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[15,30,100]" :first="searchStore.units.start"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <template #empty><h3>No matching units found</h3></template>
      <template #paginatorstart>
         <div class="acts">
            <DPGButton label="Download Results CSV" severity="secondary" @click="downloadCSV" v-if="searchStore.units.total>0" />
            <DPGButton v-if="hasFilter" label="Clear All Filters" severity="secondary" @click="clearFilters"/>
         </div>
      </template>
      <Column field="id" header="ID">
         <template #body="slotProps">
            <router-link :to="`/units/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="status" header="Status" class="nowrap" filterField="status" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Select v-model="filterModel.value" :options="unitStatuses" optionLabel="name" optionValue="code" placeholder="Select a status" />
         </template>
         <template #body="slotProps">
            <span :class="`status ${slotProps.data.status}`">{{displayStatus(slotProps.data.status)}}</span>
         </template>
      </Column>
      <Column field="staff_notes" header="Staff Notes" filterField="staff_notes" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Notes"/>
         </template>
         <template #body="slotProps">
            <span v-if="slotProps.data.staff_notes">{{ slotProps.data.staff_notes }}</span>
            <span v-else class="none">N/A</span>
         </template>
      </Column>
      <Column field="special_instructions" header="Special Instructions" filterField="special_instructions" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Instructions"/>
         </template>
         <template #body="slotProps">
            <span v-if="slotProps.data.special_instructions">{{ slotProps.data.special_instructions }}</span>
            <span v-else class="none">N/A</span>
         </template>
      </Column>
      <Column field="date_dl_deliverables_ready" header="DL Deliverable Date" class="nowrap">
         <template #body="slotProps">
            <span v-if="slotProps.data.date_dl_deliverables_ready">{{ $formatDate(slotProps.data.date_dl_deliverables_ready) }}</span>
            <span v-else class="none">N/A</span>
         </template>
      </Column>
      <Column field="date_patron_deliverables_ready" header="Patron Deliverable Date" class="nowrap">
         <template #body="slotProps">
            <span v-if="slotProps.data.date_patron_deliverables_ready">{{ $formatDate(slotProps.data.date_patron_deliverables_ready) }}</span>
            <span v-else class="none">N/A</span>
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
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import { useRoute, useRouter } from 'vue-router'
import { usePinnable } from '@/composables/pin'

usePinnable("p-datatable-paginator-top")

const route = useRoute()
const router = useRouter()
const searchStore = useSearchStore()

const unittHitsTable = ref()

const filters = ref( {
   'status': {value: null, matchMode: FilterMatchMode.EQUALS},
   'staff_notes': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'special_instructions': {value: null, matchMode: FilterMatchMode.CONTAINS},
})

const unitStatuses = ref([
   {name: "Approved", code: "approved"},
   {name: "Unapproved", code: "unapproved"},
   {name: "Canceled", code: "canceled"},
   {name: "Done", code: "done"},
   {name: "Error", code: "error"},
])

const hasFilter = computed(() => {
   let idx = Object.values(filters.value).findIndex( fv => fv.value && fv.value != "")
   return idx >= 0
})

onMounted(() => {
   searchStore.units.filters.forEach( fv => {
      filters.value[fv.field].value = fv.value
   })
})

const displayStatus = ((id) => {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
})

const downloadCSV = (() => {
   unittHitsTable.value.exportCSV()
})

const clearFilters = (() => {
   Object.values(filters.value).forEach( fv => fv.value = null )
   searchStore.units.filters = []
   let query = Object.assign({}, route.query)
   delete query.filters
   router.push({query})
   searchStore.executeSearch("units")
})

const onFilter = ((event) => {
   searchStore.units.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         searchStore.units.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   let query = Object.assign({}, route.query)
   query.filters = searchStore.filtersAsQueryParam("units")
   delete query.filters
   if ( searchStore.units.filters.length > 0) {
      query.filters = searchStore.filtersAsQueryParam("units")
   }
   router.push({query})
   searchStore.executeSearch("units")
})

const onPage = ((event) => {
   searchStore.units.start = event.first
   searchStore.units.limit = event.rows
   searchStore.executeSearch("units")
})
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