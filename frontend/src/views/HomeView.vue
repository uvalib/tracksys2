<template>
   <h2>Home</h2>
   <div class="actions" v-if="(userStore.isAdmin || userStore.isSupervisor)" >
      <DPGButton label="Create Metadata" class="create" @click="createMetadata()"/>
      <DPGButton label="Create Order" class="create" @click="createOrder()"/>
   </div>
   <div class="home">
      <div class="stats">
         <div>
            <label>Orders due in one week:</label>
            <router-link to="/orders?filter=due_week&sort=id+desc">{{dashboard.dueInOneWeek}}</router-link>
         </div>
         <span class="sep"></span>
         <div>
            <label>Overdue orders:</label>
            <router-link to="/orders?filter=overdue&sort=id+desc">{{dashboard.overdue}}</router-link>
         </div>
         <span class="sep"></span>
         <div>
            <label>Orders ready for delivery:</label>
            <router-link to="/orders?filter=ready&sort=id+desc">{{dashboard.readyForDelivery}}</router-link>
         </div>
      </div>
      <div class="search">
         <FormKit type="form" id="global-search" :actions="false" @submit="doSearch">
            <FormKit type="select" label="" v-model="selectedScope" outer-class="select-wrap" @change="scopeChanged"
               :options="{ all: 'All items', orders: 'Orders', metadata: 'Metadata', masterfiles: 'Master Files', components: 'Components'}"
            />
            <FormKit type="select" label="" v-model="selectedField" :options="scopeFields" outer-class="select-wrap"/>
            <FormKit label="" type="search" placeholder="Find Tracksys items..." v-model="newQuery" outer-class="searchbar" />
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
   <Dialog v-model:visible="showCreateMetadata" :modal="true" header="Create Metadata" @hide="createMetadataClosed" :style="{width: '750px'}">
      <NewMetadataPanel @canceled="createMetadataClosed" @created="metadataCreated" />
   </Dialog>
</template>

<script setup>
import { useSearchStore } from '../stores/search'
import { useDashboardStore } from '../stores/dashboard'
import { useUserStore } from '../stores/user'
import { useSystemStore } from '../stores/system'
import { useMetadataStore } from '../stores/metadata'
import SearchResults from '@/components/results/SearchResults.vue'
import { ref, computed, onBeforeMount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Dialog from 'primevue/dialog'
import NewMetadataPanel from '../components/order/NewMetadataPanel.vue'

const searchStore = useSearchStore()
const route = useRoute()
const router = useRouter()
const dashboard = useDashboardStore()
const userStore = useUserStore()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()

const unitID = ref("")
const unitError = ref("")
const showCreateMetadata = ref(false)

const selectedScope = ref("all")
const selectedField = ref("all")
const newQuery = ref("")

const scopeFields = computed( () => {
   let scope = selectedScope.value
   let allFields = searchStore.searchFields
   let fields = allFields[scope]
   if (fields) {
      return fields
   }
   return [{label: 'All fields', value: "all"}]
})

onBeforeMount( () => {
   document.title = `Tracksys2`
   dashboard.getStatistics()

   let paramsChanged = false

   newQuery.value = ""
   selectedField.value = "all"
   selectedScope.value = "all"

   // detect and set scope first as it affects all other aspects of the search
   if ( route.query.scope ) {
      selectedScope.value = route.query.scope
      if (searchStore.scope != route.query.scope ) {
         // console.log("SCOPE CHANGE "+searchStore.scope+" vs new q "+route.query.scope)
         paramsChanged = true
         searchStore.scope = route.query.scope
      }

      // if scope anything but all, ensure view matches it
      if ( route.query.scope != "all" ) {
         searchStore.setActiveView(route.query.scope)
      }
   } else {
      searchStore.scope = "all"
      selectedScope.value = "all"
   }

   // view is set next because it controls which filters get applied
   if ( route.query.view ) {
      searchStore.view = route.query.view
   }
   if ( route.query.q  ) {
      // paramsDetected = true
      newQuery.value = route.query.q
      if (searchStore.query != route.query.q) {
         // console.log("QUERY CHANGE "+searchStore.query+" vs new q "+route.query.q)
         searchStore.query = route.query.q
         paramsChanged = true
      }
   } else {
      newQuery.value = ""
   }

   if ( route.query.field ) {
      selectedField.value = route.query.field
      if ( searchStore.field != route.query.field ) {
         searchStore.field = route.query.field
         paramsChanged = true
      }
   }

   if ( route.query.filters ) {
      searchStore.setFilter(route.query.filters)
   }

   if (paramsChanged) {
      searchStore.executeSearch()
   }
})

function resetSearch() {
   searchStore.resetSearch()
   selectedScope.value = "all"
   selectedField.value = "all"
   newQuery.value = ""
   let query = Object.assign({}, route.query)
   delete query.q
   delete query.scope
   delete query.field
   delete query.filters
   delete query.view
   router.push({query})
}

function scopeChanged() {
   selectedField.value = "all"
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
   if (newQuery.value.length > 0) {
      unitError.value = ""

      // this is only called when clicking search. reset everything.
      searchStore.resetSearch()

      // promote local changes to the store. these will be used in the search. This promotion is necessary
      // because the UI would change before search is clicked otherwise.
      searchStore.scope = selectedScope.value
      searchStore.field = selectedField.value
      searchStore.query = newQuery.value
      if ( searchStore.scope != "all") {
         // if the scope is narrowed to a single type, the view must be too.
         // In that case, there is only 1 result. Set the active result index to 0.
         console.log("doSearch set active view")
         searchStore.setActiveView(selectedScope.value)
      }

      // convert the search store into query params so it can be shared / bookmarked
      let query = Object.assign({}, route.query)
      query.q = searchStore.query
      query.scope = searchStore.scope
      query.field = searchStore.field
      delete query.view
      delete query.filters
      let filterQP = searchStore.filtersAsQueryParam(searchStore.scope)
      if (filterQP != "") {
         query.filters = filterQP
      }

      router.push({query})

      // do the search last. This will pick a view and upodate the URL to include it.
      searchStore.executeSearch()
   }
}

function createOrder() {
   router.push("/orders/new")
}

function createMetadata() {
   showCreateMetadata.value = true
}
function metadataCreated() {
   systemStore.toastMessage("Metadata Created", `Metadata ${metadataStore.dl.pid}: ${metadataStore.detail.title} has been created.`)
   showCreateMetadata.value = false
}
function createMetadataClosed() {
   showCreateMetadata.value = false
}
</script>

<style scoped lang="scss">
h2 {
   margin-bottom: 0 !important;
}
div.actions {
   position: absolute;
   right:15px;
   top: 15px;
   button.p-button {
      margin-right: 5px;
      font-size: 0.9em;
   }
}
   .home {
      margin-top: 0px;
      padding-bottom: 50px;
      min-height:600px;
      .stats {
         margin: 0px auto 40px auto;
         display: flex;
         flex-flow: row wrap;
         justify-content: center;
         background: #fafafa;
         padding: 15px;
         border-bottom: 1px solid var(--uvalib-grey-light);

         label {
            font-weight: 600;
            margin-right: 5px;
         }
         .sep {
            display: inline-block;
            width: 50px;
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
         width: 70%;
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
            width: 125px;
         }
         .reset-button button {
            @include base-button();
            font-size: 0.95em;
            padding: 6px 15px;
            margin-left: 10px;
            display: inline-block;
            width: 125px;
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