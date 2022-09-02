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
import SearchResults from '@/components/SearchResults.vue'
import { computed } from 'vue'

const searchStore = useSearchStore()

const scopeFields = computed( () => {
   let scope = searchStore.scope
   let allFields = searchStore.searchFields
   let fields = allFields[scope]
   if (fields) {
      return fields
   }
   return [{label: 'All fields', value: "all"}]
})

function scopeChanged() {
   searchStore.field = "all"
}

function doSearch() {
   if (searchStore.query.length > 0) {
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