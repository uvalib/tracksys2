<template>
   <div class="home">
      <h2>Search</h2>
      <FormKit type="form" id="global-search" :actions="false" @submit="doSearch">
         <FormKit type="select" label="" v-model="searchStore.scope" outer-class="select-wrap" @change="scopeChanged"
            :options="{ all: 'All items', orders: 'Orders', masterfiles: 'Master Files', metadata: 'Metadata', components: 'Components'}"
         />
         <FormKit type="select" label="" v-model="searchStore.field" :options="scopeFields" outer-class="select-wrap"/>
         <FormKit label="" type="search" placeholder="Find Tracksys items..." v-model="searchStore.query" outer-class="searchbar" />
         <FormKit type="submit" label="Search" wrapper-class="submit-button" />
      </FormKit>
      <SearchResults v-if="searchStore.searched" />
   </div>
</template>

<script setup>
import { useSearchStore } from '../stores/search'
import SearchResults from '@/components/results/SearchResults.vue'
import { computed, onBeforeMount } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const searchStore = useSearchStore()
const route = useRoute()
const router = useRouter()

const scopeFields = computed( () => {
   let scope = searchStore.scope
   console.log("new scope "+scope)
   let allFields = searchStore.searchFields
   let fields = allFields[scope]
   if (fields) {
      return fields
   }
   return [{label: 'All fields', value: "all"}]
})


onBeforeMount( () => {
   let paramsDetected = false
   if ( route.query.q ) {
      searchStore.query = route.query.q
      paramsDetected = true
   }
   if ( route.query.scope ) {
      searchStore.scope = route.query.scope
      paramsDetected = true
      if ( route.query.filters ) {
         searchStore.setFilter(route.query.scope, route.query.filters)
      }
   }
   if ( route.query.field ) {
      searchStore.field = route.query.field
      paramsDetected = true
   }

   if (paramsDetected) {
      searchStore.executeSearch(route.query.scope)
   }
})

function scopeChanged() {
   searchStore.field = "all"
}

function doSearch() {
   if (searchStore.query.length > 0) {
      let query = Object.assign({}, route.query)
      query.q = searchStore.query
      delete query.scope
      if (searchStore.scope != "all") {
         query.scope = searchStore.scope
      }
      delete query.field
      if (searchStore.field != "all") {
         query.field = searchStore.field
      }
      router.push({query})

      searchStore.executeSearch("all")
   }
}
</script>

<style scoped lang="scss">
   .home {
      padding-top: 10px;
      padding-bottom: 50px;
      min-height:600px;
      :deep(#global-search) {
         display: flex;
         flex-flow: row nowrap;
         justify-content: center;
         align-items: center;
         width: 50%;
         margin: 0 auto;
         span  {
            font-weight: bold;
            display: inline-block;
         }
         .searchbar {
            flex-grow: 1;
            margin: 0;
            input {
               margin: 0;
            }
         }
         .submit-button button {
            @include primary-button();
            font-size: 0.95em;
            padding: 6px 15px;
            margin-left: 10px;
         }
         .formkit-outer.select-wrap {
            margin: 0 10px 0 0;
            select {
               margin: 0;
            }
         }
      }
   }
</style>