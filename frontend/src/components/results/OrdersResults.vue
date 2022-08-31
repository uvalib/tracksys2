<template>
   <div v-if="searchStore.orders.total == 0">
      <h3>No matching orders found</h3>
   </div>
   <DataTable v-else :value="searchStore.orders.hits" ref="orderHitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="true" :paginator="true" @page="onPage($event)"
      :rows="searchStore.orders.limit" :totalRecords="searchStore.orders.total"
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10,30,100]"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <Column field="id" header="ID"/>
      <Column header="Customer" class="nowrap">
         <template #body="slotProps">
            {{slotProps.data.customer.lastName}}, {{slotProps.data.customer.firstName}}
         </template>
      </Column>
      <Column field="agency.name" header="Agency" class="nowrap"/>
      <Column field="title" header="Order Title"/>
      <Column field="notes" header="Staff Notes" />
      <Column field="specialInstructions" header="Special Instructions" />
   </DataTable>
</template>

<script setup>
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const searchStore = useSearchStore()

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
}
</stype>