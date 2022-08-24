<template>
<TabView class="results">
   <TabPanel header="Metadata" v-if="searchStore.scope=='all' || searchStore.scope=='metadata'">
      <div v-if="searchStore.metadata.length == 0">
         <h3>No matching metadata records found</h3>
      </div>
      <DataTable v-else :value="searchStore.metadata" ref="metadataHitsTable"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      >
         <Column field="id" header="ID"  :sortable="true"/>
         <Column field="type" header="Type" :sortable="true"/>
         <Column field="title" header="Title"  :sortable="true"/>
         <Column field="creatorName" header="Creator Name" :sortable="true"/>
         <Column field="barcode" header="Barcode"  :sortable="true"/>
         <Column field="callNumber" header="Call Number" class="nowrap" :sortable="true"/>
         <Column field="catalogKey" header="Catalog Key" class="nowrap" :sortable="true"/>
      </DataTable>
   </TabPanel>
   <TabPanel header="Orders"  v-if="searchStore.scope=='all' || searchStore.scope=='orders'">
      <div v-if="searchStore.orders.length == 0">
         <h3>No matching orders found</h3>
      </div>
      <DataTable v-else :value="searchStore.orders" ref="orderHitsTable"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      >
         <Column field="id" header="ID" :sortable="true"/>
         <Column header="Customer" class="nowrap">
            <template #body="slotProps">
              {{slotProps.data.customer.lastName}}, {{slotProps.data.customer.firstName}}
            </template>
         </Column>
         <Column field="agency.name" header="Agency" class="nowrap" :sortable="true"/>
         <Column field="title" header="Order Title" :sortable="true"/>
         <Column field="staffNotes" header="Staff Notes" />
         <Column field="specialInstructions" header="Special Instructions" />
      </DataTable>
   </TabPanel>
   <TabPanel header="Master Files"  v-if="searchStore.scope=='all' || searchStore.scope=='masterfiles'">
      <div v-if="searchStore.masterFiles.length == 0">
         <h3>No matching master files found</h3>
      </div>
      <DataTable v-else :value="searchStore.masterFiles" ref="masterFileHitsTable"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      >
         <Column field="id" header="ID" :sortable="true"/>
         <Column field="pid" header="ID" :sortable="true"/>
         <Column field="filename" header="Filename" :sortable="true"/>
         <Column field="title" header="Title" :sortable="true"/>
         <Column field="description" header="Description" />
      </DataTable>
   </TabPanel>
</TabView>
</template>

<script setup>
import { useSearchStore } from '../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'

const searchStore = useSearchStore()
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