<template>
   <div v-if="searchStore.components.total == 0">
      <h3>No matching components found</h3>
   </div>
   <DataTable v-else :value="searchStore.components.hits" ref="componentHitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="true" :paginator="true" @page="onPage($event)"
      :rows="searchStore.components.limit" :totalRecords="searchStore.components.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10,30,100]"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <Column field="id" header="ID"/>
      <Column field="pid" header="PID" class="nowrap"/>
      <Column field="title" header="Title"/>
      <Column field="label" header="Label"/>
      <Column field="description" header="Content Description"/>
      <Column field="date" header="Date" class="nowrap"/>
      <Column field="eadID" header="Finding Aide" class="nowrap"/>
   </DataTable>
</template>

<script setup>
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const searchStore = useSearchStore()

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
}
</stype>