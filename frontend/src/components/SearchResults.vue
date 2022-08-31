<template>
<TabView class="results">
   <TabPanel :header="`Orders (${searchStore.orders.total})`"  v-if="searchStore.scope=='all' || searchStore.scope=='orders'">
      <OrdersResults />
   </TabPanel>
   <TabPanel :header="`Metadata (${searchStore.metadata.total})`" v-if="searchStore.scope=='all' || searchStore.scope=='metadata'">
      <MetadataResults />
   </TabPanel>
   <TabPanel :header="`Master Files (${searchStore.masterFiles.total})`"  v-if="searchStore.scope=='all' || searchStore.scope=='masterfiles'">
      <MasterFilesResults />
   </TabPanel>
   <TabPanel :header="`Components (${searchStore.components.total})`"  v-if="searchStore.scope=='all' || searchStore.scope=='components'">
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