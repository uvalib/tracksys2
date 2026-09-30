<template>
   <div class="results">
      <!-- NOTE: to make tabs work with the sybolic name, set valueKey to teh slot field of the items -->
      <UTabs :items="tabs" variant="link" v-model="searchStore.view" @update:modelValue="tabChanged" valueKey="slot">
         <template #orders><OrdersResults /></template>
         <template #metadata><MetadataResults /></template>
         <template #masterfiles><MasterFilesResults /></template>
         <template #components><ComponentsResults /></template>
         <template #units><UnitsResults /></template>
      </UTabs>
   </div>
</template>

<script setup>
import { useSearchStore } from '@/stores/search'
import MetadataResults from '@/components/results/MetadataResults.vue'
import OrdersResults from '@/components/results/OrdersResults.vue'
import MasterFilesResults from '@/components/results/MasterFilesResults.vue'
import ComponentsResults from '@/components/results/ComponentsResults.vue'
import UnitsResults from '@/components/results/UnitsResults.vue'
import { useRoute, useRouter } from 'vue-router'
import { onMounted, computed } from 'vue'

const searchStore = useSearchStore()
const route = useRoute()
const router = useRouter()

const tabs = computed(() => {
   let out = []
   out.push( { label: `Orders (${searchStore.orders.total}) hits`, slot: 'orders', disabled: searchStore.orders.total == 0 } )
   out.push( { label: `Metadata (${searchStore.metadata.total}) hits`, slot: 'metadata', disabled: searchStore.metadata.total == 0 } )
   out.push( { label: `Master Files (${searchStore.masterFiles.total}) hits`, slot: 'masterfiles', disabled: searchStore.masterFiles.total == 0 } )
   out.push( { label: `Components (${searchStore.components.total}) hits`, slot: 'components', disabled: searchStore.components.total == 0 } )
   out.push( { label: `Units (${searchStore.units.total}) hits`, slot: 'units', disabled: searchStore.units.total == 0 } )
   return out
})

onMounted( () => {
   let query = Object.assign({}, route.query)
   if ( query.view ) {
      searchStore.setActiveView(query.view)
   } 
})

const tabChanged =(() => {
   if (searchStore.scope == "all") {
      let query = Object.assign({}, route.query)
      query.view = searchStore.view
      let fp = searchStore.filtersAsQueryParam(searchStore.view)
      if (fp != "") {
         query.filters = fp
      } else {
         delete query.filters
      }
      router.push({query})
   }
})
</script>

<style scoped lang="scss">
   .results {
      margin: 25px 0;
   }
</style>