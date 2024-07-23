<template>
   <div class="results">
      <Tabs :value="searchStore.view" @update:value="tabChanged" :lazy="true">
         <TabList>
            <Tab value="orders" :disabled="searchStore.orders.total==0">
               Orders ({{ searchStore.orders.total }} hits)
            </Tab>
            <Tab value="metadata" :disabled="searchStore.metadata.total==0">
               Metadata ({{ searchStore.metadata.total }} hits)
            </Tab>
            <Tab value="masterfiles" :disabled="searchStore.masterFiles.total==0">
               Master Files ({{ searchStore.masterFiles.total }} hits)
            </Tab>
            <Tab value="components" :disabled="searchStore.components.total==0">
               Components ({{ searchStore.components.total }} hits)
            </Tab>
            <Tab value="units" :disabled="searchStore.units.total==0">
               Units ({{ searchStore.units.total }} hits)
            </Tab>
         </TabList>
         <TabPanels>
            <TabPanel value="orders">
               <OrdersResults />
            </TabPanel>
            <TabPanel value="metadata">
               <MetadataResults />
            </TabPanel>
            <TabPanel value="masterfiles">
               <MasterFilesResults />
            </TabPanel>
            <TabPanel value="components">
               <ComponentsResults />
            </TabPanel>
            <TabPanel value="units">
               <UnitsResults />
            </TabPanel>
         </TabPanels>
      </Tabs>
   </div>
</template>

<script setup>
import { useSearchStore } from '@/stores/search'
import { useSystemStore } from '@/stores/system'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import MetadataResults from '@/components/results/MetadataResults.vue'
import OrdersResults from '@/components/results/OrdersResults.vue'
import MasterFilesResults from '@/components/results/MasterFilesResults.vue'
import ComponentsResults from '@/components/results/ComponentsResults.vue'
import UnitsResults from '@/components/results/UnitsResults.vue'
import { useRoute, useRouter } from 'vue-router'
import { onMounted, watch } from 'vue'

const searchStore = useSearchStore()
const systemStore = useSystemStore()
const route = useRoute()
const router = useRouter()

watch(() => systemStore.working, (newVal) => {
   if ( newVal == false) {
      showTargetView()
   }
})

onMounted( () => {
   showTargetView()
})

const showTargetView = (() => {
   let query = Object.assign({}, route.query)
   if (query.view) {
      if ( query.view != searchStore.view && searchStore.scope == "all") {
         router.push({query})
      }
      searchStore.setActiveView(query.view)
   } else if ( searchStore.view ) {
     if ( searchStore.scope == "all") {
         query.view = searchStore.view
         router.push({query})
     }
   }
})

const tabChanged =(( newTab ) => {
   searchStore.view = newTab
   if (searchStore.scope == "all") {
      let query = Object.assign({}, route.query)
      query.view = newTab
      let fp = searchStore.filtersAsQueryParam(query.view)
      if (fp != "") {
         query.filters = fp
      } else {
         delete query.filters
      }
      router.push({query})
   }
})
</script>

<stype scoped lang="scss">
   .results {
      margin: 25px 0;
   }
</stype>