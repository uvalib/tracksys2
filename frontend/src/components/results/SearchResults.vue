<template>
   <TabView class="results" @tabChange="tabChanged()" v-model:activeIndex="searchStore.activeResultsIndex">
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
import { useSystemStore } from '@/stores/system'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import MetadataResults from './MetadataResults.vue'
import OrdersResults from './OrdersResults.vue'
import MasterFilesResults from './MasterFilesResults.vue'
import ComponentsResults from './ComponentsResults.vue'
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

function showTargetView() {
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
}

function tabChanged() {
   if (searchStore.scope == "all") {
      let tabs = ['orders', 'metadata', 'masterfiles', 'components']
      let query = Object.assign({}, route.query)
      query.view = tabs[searchStore.activeResultsIndex]
      searchStore.view = query.view
      let fp = searchStore.filtersAsQueryParam(query.view)
      if (fp != "") {
         query.filters = fp
      } else {
         delete query.filters
      }
      router.push({query})
   }
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