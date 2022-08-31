<template>
   <div v-if="searchStore.masterFiles.total == 0">
      <h3>No matching master files found</h3>
   </div>
   <DataTable v-else :value="searchStore.masterFiles.hits" ref="masterFileHitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="true" :paginator="true" @page="onPage($event)"
      :rows="searchStore.masterFiles.limit" :totalRecords="searchStore.masterFiles.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10,30,100]"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <Column field="id" header="ID"/>
      <Column field="pid" header="PID" class="nowrap" />
      <Column field="filename" header="Filename"/>
      <Column field="title" header="Title"/>
      <Column field="description" header="Description" />
   </DataTable>
</template>

<script setup>
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const searchStore = useSearchStore()

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
</stype>