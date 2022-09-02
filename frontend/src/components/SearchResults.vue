<template>
   <div class="summary">
      <div class="title">Search summary</div>
      <div class="content">
         <table>
            <tr><td class="label">Orders:</td><td class="count">{{searchStore.orders.total}} matches</td></tr>
            <tr><td class="label">Metadata:</td><td class="count">{{searchStore.metadata.total}} matches</td></tr>
            <tr><td class="label">Master Files:</td><td class="count">{{searchStore.masterFiles.total}} matches</td></tr>
            <tr><td class="label">Components:</td><td class="count">{{searchStore.components.total}} matches</td></tr>
         </table>
         <div class="actions">
            <DPGButton label="Reset search" class="p-button-secondary" @click="resetSearch()"/>
         </div>
      </div>
   </div>
   <TabView class="results">
      <TabPanel :header="`Orders`">
         <OrdersResults />
      </TabPanel>
      <TabPanel :header="`Metadata`">
         <MetadataResults />
      </TabPanel>
      <TabPanel :header="`Master Files`">
         <MasterFilesResults />
      </TabPanel>
      <TabPanel :header="`Components`">
         <ComponentsResults />
      </TabPanel>
   </TabView>
</template>

<script setup>
import { useSearchStore } from '../stores/search'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import MetadataResults from './results/MetadataResults.vue'
import OrdersResults from './results/OrdersResults.vue'
import MasterFilesResults from './results/MasterFilesResults.vue'
import ComponentsResults from './results/ComponentsResults.vue'

const searchStore = useSearchStore()

function resetSearch() {
   searchStore.$reset()
}

</script>

<stype scoped lang="scss">
.summary {
   text-align: left;
   margin: 25px 20px 0 20px;
   border: 1px solid #e9ecef;
   .content {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      .actions {
         font-size: 0.9em;
         padding: 15px;
      }
   }
   .title {
      padding: 5px 10px;
      background: var(--uvalib-grey-lightest);
      border-bottom: 1px solid #e9ecef;
   }
   table  {
      font-size: 0.9em;
      border-collapse: collapse;
      margin: 15px;
      td {
         padding: 2px 5px;
      }
      td.label {
         font-weight: bold;
         text-align: right;
      }
   }
}
.results {
   margin: 20px;
   font-size: 0.9em;

   td.nowrap,
   th {
      white-space: nowrap;
   }

   th,
   td {
      font-size: 0.85em;
   }
}
</stype>