<template>
   <div class="home">
      <FormKit type="form" id="global-search" :actions="false" @submit="doSearch">
         <span>Search</span>
         <FormKit type="select" label="" v-model="searchStore.scope"
            :options="{ all: 'Everything', orders: 'Orders', masterfiles: 'Master Files', metadata: 'Metadata' }"
         />
         <FormKit label="" type="search" placeholder="Find Tracksys items..." v-model="searchStore.query" outer-class="searchbar" />
         <FormKit type="submit" label="Search" wrapper-class="submit-button" />
      </FormKit>
      <SearchResults v-if="searchStore.hasResults && !systemStore.working"/>
   </div>
</template>

<script setup>
import { useSearchStore } from '../stores/search'
import { useSystemStore } from '../stores/system'
import SearchResults from '@/components/SearchResults.vue'

const searchStore = useSearchStore()
const systemStore = useSystemStore()

function doSearch() {
   if (searchStore.query.length > 0) {
      searchStore.globalSearch()
   }
}
</script>

<style scoped lang="scss">
   .home {
      padding-top: 50px;
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
            margin-right: 10px;
         }
         .searchbar {
            flex-grow: 1;
            margin: 0 5px;
         }
         .submit-button button {
            @include primary-button();
            font-size: 0.95em;
         padding: 6px 15px;
         }
         .dpg-form-input {
            margin-bottom: 0 !important;
         }
      }
   }
</style>