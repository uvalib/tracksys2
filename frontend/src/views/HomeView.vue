<template>
   <h2>
      <span>Home</span>
      <div class="row-right gap-2!" v-if="(userStore.isAdmin || userStore.isSupervisor)" >
         <CreateAgencyModal v-if="userStore.isAdmin" />
         <CreateCollectionFacetModal v-if="userStore.isAdmin" />
         <CreateMetadataModal />
         <UButton label="Create Order" @click="createOrder()"/>
      </div>
   </h2>
   <div class="home">
      <div class="stats">
         <div class="row">
            <div class="value">
               <label>Orders due in one week:</label>
               <router-link  v-if="dashboard.dueInOneWeek" to='/orders?filters=["status|equals|due_week"]&sort=id+desc'>{{dashboard.dueInOneWeek}}</router-link>
               <span v-else>0</span>
            </div>
            <div class="value">
               <label>Overdue orders:</label>
               <router-link  v-if="dashboard.overdue" to='/orders?filters=["status|equals|overdue"]&sort=id+desc'>{{dashboard.overdue}}</router-link>
               <span v-else>0</span>
            </div>
            <div class="value">
               <label>Orders ready for delivery:</label>
               <router-link v-if="dashboard.readyForDelivery" to='/orders?filters=["status|equals|ready"]&sort=id+desc'>{{dashboard.readyForDelivery}}</router-link>
               <span v-else>0</span>
            </div>
         </div>
         <div class="row">
            <div class="value">
               <label>ArchivesSpace Requests:</label>
               <router-link v-if="dashboard.asRequests" to='/archivesspace?view=request'>{{dashboard.asRequests}}</router-link>
               <span v-else>0</span>
            </div>
            <div class="value">
               <label>ArchivesSpace Reviews:</label>
               <router-link v-if="dashboard.asReviews" to='/archivesspace?view=review'>{{dashboard.asReviews}}</router-link>
               <span v-else>0</span>
            </div>
            <div class="value">
               <label>ArchivesSpace Rejections:</label>
               <router-link v-if="dashboard.asRejections" to='/archivesspace?view=reject'>{{dashboard.asRejections}}</router-link>
               <span v-else>0</span>
            </div>
         </div>
      </div>

      <div class="search">
         <div class="text-search">
            <USelect v-model="selectedScope" :items="scopes"  />
            <div class="search-info">
               <div class="row-left">
                  <SearchIndexPopover />
                  <SearchHelpPopover />
               </div>
               <UInput placeholder="Find TrackSys items..." v-model="newQuery"  @keyup.enter="doSearch" />
            </div>

            <UButton label="Search"  @click="doSearch"/>
            <UButton v-if="searchStore.searched || searchStore.similarSearch == true" label="Reset Search" color="secondary" @click="resetSearch"/>
         </div>

         <UCard v-if="userStore.isAdmin" title="Search for similar images" class="w-fit mx-auto! my-5!">
            <div class="image-search">
               <p class="hint">Set a similarity threshold, select a search image then click 'Search Images'</p>
               <div class="labels">
                  <span>More Similar</span>
                  <span>Less Similar</span>
               </div>
               <USlider class="w-full" :min="5" :max="20" v-model="searchStore.distance" @change="slideChanged"/>
               <UFileUpload accept="image/*" label="Drop your image here" v-model="lookupImage" @change="uploadImageChanged"/>
               <UButton color="secondary" size="sm" label="Search Images" @click="startImageUpload" :disabled="!lookupImage"/>
            </div>
         </UCard>

         <template v-if="systemStore.working == false">
            <SearchResults v-if="searchStore.searched" />
            <SimilarImages v-if="searchStore.similarSearch" />
         </template>
      </div>
   </div>
</template>

<script setup>
import { useSearchStore } from '../stores/search'
import { useDashboardStore } from '../stores/dashboard'
import { useUserStore } from '../stores/user'
import { useSystemStore } from '../stores/system'
import SearchIndexPopover from '@/components/SearchIndexPopover.vue'
import SearchHelpPopover from '@/components/SearchHelpPopover.vue'
import SearchResults from '@/components/results/SearchResults.vue'
import SimilarImages from '@/components/results/SimilarImages.vue'
import CreateAgencyModal from '@/components/CreateAgencyModal.vue'
import CreateCollectionFacetModal from '@/components/CreateCollectionFacetModal.vue'
import { ref, computed, onBeforeMount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import CreateMetadataModal from '../components/CreateMetadataModal.vue'

const searchStore = useSearchStore()
const route = useRoute()
const router = useRouter()
const dashboard = useDashboardStore()
const userStore = useUserStore()
const systemStore = useSystemStore()

const selectedScope = ref("all")
const newQuery = ref("")
const lookupImage = ref()

const scopes = computed( () => {
   return [
      {label: "All items", value: "all"},
      {label: "Orders", value: "orders"},
      {label: "Metadata", value: "metadata"},
      {label: "Master Files", value: "masterfiles"},
      {label: "Components", value: "components"},
      {label: "Units", value: "units"},
   ]
})

onBeforeMount( () => {
   document.title = `Tracksys`
   dashboard.getStatistics()

   let paramsChanged = false

   newQuery.value = ""
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

   if ( route.query.filters ) {
      searchStore.setFilter(route.query.filters)
   }

   if (paramsChanged) {
      searchStore.executeSearch()
   }
})

const slideChanged = ( () => {
   if (searchStore.similarSearch == true && searchStore.searchPHash !== 0) {
      searchStore.imageSearch()
   }
})
const uploadImageChanged = (() => {
   searchStore.resetImageSearch()
})
const startImageUpload = ( async () => {
   searchStore.uploadSearchImage( lookupImage.value )
})

const resetSearch = (() => {
   searchStore.resetSearch()
   selectedScope.value = "all"
   newQuery.value = ""
   let query = Object.assign({}, route.query)
   delete query.q
   delete query.scope
   delete query.field
   delete query.filters
   delete query.view
   router.push({query})
})

const doSearch = (() => {
   if (newQuery.value.length > 0) {
      // this is only called when clicking search. reset everything.
      searchStore.resetSearch()

      // promote local changes to the store. these will be used in the search. This promotion is necessary
      // because the UI would change before search is clicked otherwise.
      searchStore.scope = selectedScope.value
      searchStore.query = newQuery.value
      if ( searchStore.scope != "all") {
         // if the scope is narrowed to a single type, the view must be too.
         // In that case, there is only 1 result. Set the active result index to 0.
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
})

const createOrder = (() => {
   router.push("/orders/new")
})

</script>

<style scoped lang="scss">
.home {
   .image-search {
      width: 275px;
      display: flex;
      flex-direction: column;
      gap: 15px;
      align-items: center;
      label {
         font-weight: bold;
      }
      .hint {
         font-size: 0.8em;
      }
      .labels {
         width: 100%;
         font-size: 0.85em;
         display: flex;
         flex-flow: row nowrap;
         justify-content: space-between;
      }
   }
   .stats {
      display: flex;
      flex-direction: column;
      gap: 5px;
      padding: 10px 0;
      border-bottom: 1px solid var(--uvalib-grey-light);
      margin-bottom: 15px;
      background: #f8f8f8;
      .row {
         display: flex;
         flex-flow: row wrap;
         justify-content: center;
         gap: 50px;
         label {
            font-weight: bold;
         }
         .value {
            display: flex;
            flex-flow: row nowrap;
            gap: 10px;
         }
      }
   }

   div.text-search {
      display: flex;
      flex-flow: row nowrap;
      justify-content: center;
      align-items: flex-end;
      width: 70%;
      margin: 0 auto;
      gap: 10px;
      .search-info {
         flex-grow: 1;
         display: flex;
         flex-direction: column;
         gap: 15px;
      }
   }
}
</style>