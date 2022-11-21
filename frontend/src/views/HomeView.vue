<template>
   <h2>Home</h2>
   <div class="home">
      <div class="stats">
         <div><label>Orders due in one week:</label><span>{{dashboard.dueInOneWeek}}</span></div>
         <span class="sep"></span>
         <div><label>Overdue orders:</label><span>{{dashboard.overdue}}</span></div>
         <span class="sep"></span>
         <div><label>Orders ready for delivery:</label><span>{{dashboard.readyForDelivery}}</span></div>
      </div>
      <div class="search">
         <FormKit type="form" id="global-search" :actions="false" @submit="doSearch">
            <FormKit type="select" label="" v-model="searchStore.scope" outer-class="select-wrap" @change="scopeChanged"
               :options="{ all: 'All items', orders: 'Orders', masterfiles: 'Master Files', metadata: 'Metadata', components: 'Components'}"
            />
            <FormKit type="select" label="" v-model="searchStore.field" :options="scopeFields" outer-class="select-wrap"/>
            <FormKit label="" type="search" placeholder="Find Tracksys items..." v-model="searchStore.query" outer-class="searchbar" />
            <FormKit type="submit" label="Search" wrapper-class="submit-button" />
            <FormKit type="button" v-if="searchStore.searched"  label="Reset search" @click="resetSearch()" wrapper-class="reset-button"/>
         </FormKit>
         <FormKit type="form" id="unit-search" :actions="false" @submit="doUnitSearch" outer-class="select-wrap" >
            <FormKit label="" type="search" placeholder="Find Unit by ID..." v-model="unitID" outer-class="searchbar" />
            <FormKit type="submit" label="Find Unit" wrapper-class="submit-button" />
         </FormKit>
         <p class="error" v-if="unitError">{{unitError}}</p>
         <SearchResults v-if="searchStore.searched" />
      </div>
   </div>
</template>

<script setup>
import { useSearchStore } from '../stores/search'
import { useDashboardStore } from '../stores/dashboard'
import SearchResults from '@/components/results/SearchResults.vue'
import { ref, computed, onBeforeMount } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const searchStore = useSearchStore()
const route = useRoute()
const router = useRouter()
const dashboard = useDashboardStore()

const unitID = ref("")
const unitError = ref("")

const scopeFields = computed( () => {
   let scope = searchStore.scope
   let allFields = searchStore.searchFields
   let fields = allFields[scope]
   if (fields) {
      return fields
   }
   return [{label: 'All fields', value: "all"}]
})

function resetSearch() {
   searchStore.resetSearch()
   let query = Object.assign({}, route.query)
   delete query.q
   delete query.scope
   delete query.field
   delete query.filters
   router.push({query})
}

onBeforeMount( () => {
   document.title = `Tracksys2`
   dashboard.getStatistics()

   let paramsDetected = false
   if ( route.query.q && searchStore.query != route.query.q ) {
      searchStore.query = route.query.q
      paramsDetected = true
   }

   if ( route.query.scope ) {
      if (searchStore.scope != route.query.scope) {
         searchStore.scope = route.query.scope
         paramsDetected = true
         if ( route.query.filters ) {
            searchStore.setFilter(route.query.scope, route.query.filters)
         }
      }
   } else {
      searchStore.scope = "all"
   }

   if ( route.query.field &&  searchStore.field != route.query.field) {
      searchStore.field = route.query.field
      paramsDetected = true
   }

   if (paramsDetected) {
      searchStore.executeSearch(searchStore.scope)
   }
})

function scopeChanged() {
   searchStore.field = "all"
}

async function doUnitSearch() {
   resetSearch()
   unitError.value = ""
   await searchStore.unitExists(unitID.value)
   if ( searchStore.unitValid == false) {
      unitError.value = `${unitID.value} is not a valid unit ID, please try again.`
      return
   }

   router.push(`/units/${unitID.value}`)
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

      let filterQP = searchStore.filtersAsQueryParam(searchStore.scope)
      if (filterQP != "") {
         query.filters = filterQP
      } else {
         delete query.filters
      }

      router.push({query})

      searchStore.executeSearch(searchStore.scope)
   }
}
</script>

<style scoped lang="scss">
   .home {
      margin-top: 15px;
      padding-bottom: 50px;
      min-height:600px;
      .stats {
         margin: 40px auto;
         display: flex;
         flex-flow: row wrap;
         justify-content: center;
         label {
            font-weight: 600;
            margin-right: 5px;
         }
         .sep {
            display: inline-block;
            width: 25px;
         }
      }

      :deep(#unit-search) {
         width: 25% !important;
         align-items: flex-end !important;
         margin-top: 25px !important;
         div.searchbar {
            width: 40%;
         }
      }
      p.error {
         color: var(--uvalib-red-emergency);
      }
      :deep(#global-search), :deep(#unit-search) {
         display: flex;
         flex-flow: row nowrap;
         justify-content: center;
         align-items: center;
         width: 75%;
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
            display: inline-block;
         }
         .reset-button button {
            @include base-button();
            font-size: 0.95em;
            padding: 6px 15px;
            margin-left: 20px;
            display: inline-block;
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