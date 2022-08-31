<template>
   <div v-if="searchStore.metadata.total == 0">
      <h3>No matching metadata records found</h3>
   </div>
   <DataTable v-else :value="searchStore.metadata.hits" ref="metadataTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll"
      :lazy="true" :paginator="true" @page="onMetadataPage($event)"
      :rows="searchStore.metadata.limit" :totalRecords="searchStore.metadata.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10,30,100]"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <Column field="id" header="ID" />
      <Column field="pid" header="PID" class="nowrap"/>
      <Column field="type" header="Type"/>
      <Column field="title" header="Title" />
      <Column field="creatorName" header="Creator Name"/>
      <Column field="barcode" header="Barcode" class="nowrap"/>
      <Column field="callNumber" header="Call Number" class="nowrap"/>
      <Column field="catalogKey" header="Catalog Key" class="nowrap"/>
   </DataTable>
</template>

<script setup>
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const searchStore = useSearchStore()

function onMetadataPage(event) {
   searchStore.metadata.start = event.first
   searchStore.metadata.limit = event.rows
   searchStore.executeSearch("metadata")
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