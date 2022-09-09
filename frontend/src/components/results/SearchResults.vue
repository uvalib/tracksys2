<template>
   <div class="summary">
      <div class="title">Search summary</div>
      <div class="content">
         <table>
            <tr v-if="searchStore.scope=='all' || searchStore.scope=='orders'">
               <td class="label">Orders:</td><td class="count">{{searchStore.orders.total}} matches</td>
            </tr>
            <tr v-if="searchStore.scope=='all' || searchStore.scope=='metadata'">
               <td class="label">Metadata:</td><td class="count">{{searchStore.metadata.total}} matches</td>
            </tr>
            <tr v-if="searchStore.scope=='all' || searchStore.scope=='masterfiles'">
               <td class="label">Master Files:</td><td class="count">{{searchStore.masterFiles.total}} matches</td>
            </tr>
            <tr v-if="searchStore.scope=='all' || searchStore.scope=='components'">
               <td class="label">Components:</td><td class="count">{{searchStore.components.total}} matches</td>
            </tr>
         </table>
         <div class="actions">
            <DPGButton label="Reset search" class="p-button-secondary" @click="resetSearch()"/>
         </div>
      </div>
   </div>
   <TabView class="results">
      <TabPanel :header="`Orders`" v-if="searchStore.orders.hits.length > 0">
         <OrdersResults />
      </TabPanel>
      <TabPanel :header="`Metadata`" v-if="searchStore.metadata.hits.length > 0">
         <MetadataResults />
      </TabPanel>
      <TabPanel :header="`Master Files`" v-if="searchStore.masterFiles.hits.length > 0">
         <MasterFilesResults />
      </TabPanel>
      <TabPanel :header="`Components`" v-if="searchStore.components.hits.length > 0">
         <ComponentsResults />
      </TabPanel>
   </TabView>
</template>

<script setup>
import { useSearchStore } from '@/stores/search'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import MetadataResults from './MetadataResults.vue'
import OrdersResults from './OrdersResults.vue'
import MasterFilesResults from './MasterFilesResults.vue'
import ComponentsResults from './ComponentsResults.vue'
import { useRoute, useRouter } from 'vue-router'

const searchStore = useSearchStore()
const route = useRoute()
const router = useRouter()

function resetSearch() {
   searchStore.resetSearch()
   let query = Object.assign({}, route.query)
   delete query.q
   delete query.scope
   delete query.field
   delete query.filters
   router.push({query})
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