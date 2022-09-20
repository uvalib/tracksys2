<template>
   <TabView class="results" @tabChange="tabChanged()" v-model:activeIndex="activeTab">
      <TabPanel :header="`Orders`" v-if="searchStore.scope=='all' || searchStore.scope=='orders'" >
         <OrdersResults />
      </TabPanel>
      <TabPanel :header="`Metadata`" v-if="searchStore.scope=='all' || searchStore.scope=='metadata'">
         <MetadataResults />
      </TabPanel>
      <TabPanel :header="`Master Files`" v-if="searchStore.scope=='all' || searchStore.scope=='masterfiles'">
         <MasterFilesResults />
      </TabPanel>
      <TabPanel :header="`Components`" v-if="searchStore.scope=='all' || searchStore.scope=='components'">
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
import { ref } from 'vue'

const searchStore = useSearchStore()
const route = useRoute()
const router = useRouter()

const activeTab = ref(0)

function tabChanged() {
   let tabs = ['orders', 'metadata', 'masterfiles', 'components']
   let query = Object.assign({}, route.query)
   query.scope = tabs[activeTab.value]
   let fp = searchStore.filtersAsQueryParam(query.scope)
   if (fp != "") {
      query.filters = fp
   } else {
      delete query.filters
   }
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