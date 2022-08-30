<template>
<TabView class="results">
   <TabPanel :header="`Orders (${searchStore.orders.length})`"  v-if="searchStore.scope=='all' || searchStore.scope=='orders'">
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
         <Column field="notes" header="Staff Notes" />
         <Column field="specialInstructions" header="Special Instructions" />
      </DataTable>
   </TabPanel>
   <TabPanel :header="`Metadata (${searchStore.metadata.length})`" v-if="searchStore.scope=='all' || searchStore.scope=='metadata'">
      <div v-if="searchStore.metadata.length == 0">
         <h3>No matching metadata records found</h3>
      </div>
      <DataTable v-else :value="searchStore.metadata" ref="metadataHitsTable"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      >
         <Column field="id" header="ID"  :sortable="true"/>
         <Column field="pid" header="PID" class="nowrap" :sortable="true"/>
         <Column field="type" header="Type" :sortable="true"/>
         <Column field="title" header="Title"  :sortable="true"/>
         <Column field="creatorName" header="Creator Name" :sortable="true"/>
         <Column field="barcode" header="Barcode" class="nowrap" :sortable="true"/>
         <Column field="callNumber" header="Call Number" class="nowrap" :sortable="true"/>
         <Column field="catalogKey" header="Catalog Key" class="nowrap" :sortable="true"/>
      </DataTable>
   </TabPanel>
   <TabPanel :header="`Master Files (${searchStore.masterFiles.length})`"  v-if="searchStore.scope=='all' || searchStore.scope=='masterfiles'">
      <div v-if="searchStore.masterFiles.length == 0">
         <h3>No matching master files found</h3>
      </div>
      <DataTable v-else :value="searchStore.masterFiles" ref="masterFileHitsTable"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      >
         <Column field="id" header="ID" :sortable="true"/>
         <Column field="pid" header="PID" class="nowrap" :sortable="true" />
         <Column field="filename" header="Filename" :sortable="true"/>
         <Column field="title" header="Title" :sortable="true"/>
         <Column field="description" header="Description" />
      </DataTable>
   </TabPanel>
   <TabPanel :header="`Components (${searchStore.components.length})`"  v-if="searchStore.scope=='all' || searchStore.scope=='components'">
      <div v-if="searchStore.components.length == 0">
         <h3>No matching components found</h3>
      </div>
      <DataTable v-else :value="searchStore.components" ref="componentHitsTable"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      >
         <Column field="id" header="ID" :sortable="true"/>
         <Column field="pid" header="PID" :sortable="true" class="nowrap"/>
         <Column field="title" header="Title" :sortable="true"/>
         <Column field="label" header="Label" :sortable="true"/>
         <Column field="description" header="Content Description"/>
         <Column field="date" header="Date" :sortable="true" class="nowrap"/>
         <Column field="eadID" header="Finding Aide" :sortable="true" class="nowrap"/>
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