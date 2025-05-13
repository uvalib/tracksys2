<template>
   <h2>
      <span>Home</span>
      <div class="actions" v-if="(userStore.isAdmin || userStore.isSupervisor)" >
         <DPGButton v-if="userStore.isAdmin" label="Create Agency" class="create" @click="showCreateAgencyClicked()"/>
         <DPGButton v-if="userStore.isAdmin" label="Create Collection Facet" class="create" @click="showCreateCollectionDialog()"/>
         <DPGButton label="Create Metadata" class="create" @click="createMetadata()"/>
         <DPGButton label="Create Order" class="create" @click="createOrder()"/>
      </div>
   </h2>
   <div class="home">
      <div class="stats">
         <div>
            <label>Orders due in one week:</label>
            <router-link  v-if="dashboard.dueInOneWeek" to='/orders?filters=["status|equals|due_week"]&sort=id+desc'>{{dashboard.dueInOneWeek}}</router-link>
            <span v-else>0</span>
         </div>
         <span class="sep"></span>
         <div>
            <label>Overdue orders:</label>
            <router-link  v-if="dashboard.overdue" to='/orders?filters=["status|equals|overdue"]&sort=id+desc'>{{dashboard.overdue}}</router-link>
            <span v-else>0</span>
         </div>
         <span class="sep"></span>
         <div>
            <label>Orders ready for delivery:</label>
            <router-link v-if="dashboard.readyForDelivery" to='/orders?filters=["status|equals|ready"]&sort=id+desc'>{{dashboard.readyForDelivery}}</router-link>
            <span v-else>0</span>
         </div>
      </div>
      <div class="stats archivesspace">
         <div>
            <label>ArchivesSpace Requests:</label>
            <router-link v-if="dashboard.asRequests" to='/archivesspace?view=request'>{{dashboard.asRequests}}</router-link>
            <span v-else>0</span>
         </div>
         <span class="sep"></span>
         <div>
            <label>ArchivesSpace Reviews:</label>
            <router-link v-if="dashboard.asReviews" to='/archivesspace?view=review'>{{dashboard.asReviews}}</router-link>
            <span v-else>0</span>
         </div>
         <span class="sep"></span>
         <div>
            <label>ArchivesSpace Rejections:</label>
            <router-link v-if="dashboard.asRejections" to='/archivesspace?view=reject'>{{dashboard.asRejections}}</router-link>
            <span v-else>0</span>
         </div>
      </div>
      <div class="search">
         <div class="text-search">
            <!-- <div class="search-ctl-group">
               <Select v-model="selectedScope" :options="scopes" optionLabel="label" optionValue="value" />
            </div> -->
            <InputText placeholder="Find TrackSys items..." v-model="newQuery" class="searchbar"  @keyup.enter="doSearch" />
            <div class="search-ctl-group">
               <DPGButton label="Search" class="submit-button" @click="doSearch"/>
               <DPGButton v-if="searchStore.searched || searchStore.similarSearch == true" label="Reset Search" severity="secondary" @click="resetSearch"/>
            </div>
         </div>

         <div class="image-search" v-if="userStore.isAdmin">
            <label>Search for similar images</label>
            <p class="hint">Set a similarity threshold then upload the search image</p>
            <div class="slide">
               <div class="labels">
                  <span>More Similar</span>
                  <span>Less Similar</span>
               </div>
               <Slider class="w-14rem" :min="5" :max="20" v-model="searchStore.distance" @change="slideChanged"/>
            </div>
            <FileUpload mode="basic" name="imageSearch" url="/upload_search_image" accept="image/*" :maxFileSize="55000000"
               @upload="imageUploaded" @before-upload="beforeUpload" :auto="true" chooseLabel="Upload Image" />
         </div>

         <p class="error" v-if="unitError">{{unitError}}</p>
         <template v-if="systemStore.working == false">
            <SearchResults v-if="searchStore.searched" />
            <SimilarImages v-if="searchStore.similarSearch" />
         </template>
      </div>
   </div>
   <Dialog v-model:visible="showCreateAgency" :modal="true" header="Create Agency"@hide="createAgencyClosed" :style="{width: '450px'}">
      <div class="agency">
         <label>Name</label>
         <input type="text" v-model="newAgencyName" autofocus/>
         <label>Description</label>
         <textarea rows="4" v-model="newAgencyDesc"/>
      </div>
      <template #footer>
         <div class="acts">
            <DPGButton @click="createAgencyClosed()" label="Cancel" severity="secondary"/>
            <DPGButton @click="createAgency()" label="Create" :disabled="newAgencyName.length == 0"/>
         </div>
      </template>
   </Dialog>
   <Dialog v-model:visible="showCreateMetadata" :modal="true" header="Create Metadata" @hide="createMetadataClosed" :style="{width: '750px'}">
      <NewMetadataPanel @canceled="createMetadataClosed" @created="metadataCreated" />
   </Dialog>
   <Dialog v-model:visible="showCreateCollection" :modal="true" header="Create Collection Facet" @hide="createCollectionClosed" :style="{width: '450px'}">
      <p>Enter the name of the new collection facet</p>
      <input type="text" v-model="newCollectionFacet" autofocus/>
      <template #footer>
         <div class="acts">
            <DPGButton @click="createCollectionClosed" label="Cancel" severity="secondary"/>
            <DPGButton @click="createCollection()" label="Create" :disabled="newCollectionFacet.length == 0"/>
         </div>
      </template>
   </Dialog>
</template>

<script setup>
import { useSearchStore } from '../stores/search'
import { useDashboardStore } from '../stores/dashboard'
import { useUserStore } from '../stores/user'
import { useSystemStore } from '../stores/system'
import { useMetadataStore } from '../stores/metadata'
import SearchResults from '@/components/results/SearchResults.vue'
import SimilarImages from '@/components/results/SimilarImages.vue'
import { ref, computed, onBeforeMount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Dialog from 'primevue/dialog'
import NewMetadataPanel from '@/components/metadata/NewMetadataPanel.vue'
import FileUpload from 'primevue/fileupload'
import Slider from 'primevue/slider'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'

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
const showCreateCollection = ref(false)
const newCollectionFacet = ref("")

const showCreateAgency = ref(false)
const newAgencyName = ref("")
const newAgencyDesc = ref("")

const selectedScope = ref("all")
const newQuery = ref("")

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
      searchStore.imageSearch( searchStore.searchPHash )
   }
})

const beforeUpload = (() => {
   systemStore.working = true
})

const imageUploaded = ((e) => {
   searchStore.imageSearch( e.xhr.responseText )
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
      unitError.value = ""

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

const showCreateAgencyClicked = ( () => {
   newAgencyDesc.value = ""
   newAgencyName.value = ""
   showCreateAgency.value = true
})
const createAgency = ( async () => {
   await systemStore.createAgency(newAgencyName.value, newAgencyDesc.value)
   showCreateAgency.value = false
})
const createAgencyClosed = ( () => {
   showCreateAgency.value = false
})

const createOrder = (() => {
   router.push("/orders/new")
})

const showCreateCollectionDialog = (() => {
   newCollectionFacet.value = ""
   showCreateCollection.value = true
})

const createCollection = ( async () => {
   await systemStore.createCollectionFacet(newCollectionFacet.value)
   showCreateCollection.value = false
})

const createCollectionClosed = (() => {
   showCreateCollection.value = false
})

const createMetadata = (() => {
   showCreateMetadata.value = true
})

const metadataCreated = (() => {
   systemStore.toastMessage("Metadata Created", `Metadata ${metadataStore.detail.pid}: ${metadataStore.detail.title} has been created.`)
   showCreateMetadata.value = false
})

const createMetadataClosed = (() => {
   showCreateMetadata.value = false
})
</script>

<style scoped lang="scss">
.home {
   margin-top: 0px;
   padding-bottom: 50px;
   min-height:600px;
   .image-search {
      width: 275px;
      margin: 50px auto 0 auto;
      border: 1px solid var(--uvalib-grey-light);
      padding: 15px 25px 25px 25px;
      border-radius: 5px;

      label {
         font-weight: 600;
         margin-bottom: 15px;
         display: inline-block;
      }
      .p-fileupload.p-fileupload-basic {
         margin-top: 5px;
      }
      .hint {
         font-size: 0.8em;
         margin: 5px 0 10px 0;;

      }
      .slide {
         margin: 15px 0 20px 0;
         .labels {
            margin: 10px 0 15px 0;
            font-size: 0.85em;
            display: flex;
            flex-flow: row nowrap;
            justify-content: space-between;
         }
      }
   }
  .stats.archivesspace {
      margin: 0 auto 40px auto;
      border-bottom: 1px solid var(--uvalib-grey-light);
      padding: 10px 0 10px 0;
   }
   .stats {
      margin: 0 auto 0 auto;
      display: flex;
      flex-flow: row wrap;
      justify-content: center;
      background: #fafafa;
      padding: 15px 0 5px;

      label {
         font-weight: 600;
         margin-right: 5px;
      }
      .sep {
         display: inline-block;
         width: 50px;
      }
   }

   p.error {
      color: var(--uvalib-red-emergency);
   }
   div.text-search {
      display: flex;
      flex-flow: row nowrap;
      justify-content: center;
      align-items: center;
      width: 70%;
      margin: 0 auto;
      gap: 10px;

      .searchbar {
         margin: 0;
      }
      select {
         margin: 0;
         width: max-content;
      }
      .search-ctl-group {
         display: flex;
         flex-flow: row nowrap;
         gap: 10px;
      }
   }
}
div.agency {
   label {
      display: block;
      margin: 10px 0 5px 0;
   }
   textarea {
      width: 100%;
      border-color: var(--uvalib-grey-light);
      border-radius: 5px;
      font-family: "franklin-gothic-urw", arial, sans-serif;
      -webkit-font-smoothing: antialiased;
      -moz-osx-font-smoothing: grayscale;
      color: var(--color-primary-text);
      padding: 5px 10px;
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   gap: 10px;
}
</style>