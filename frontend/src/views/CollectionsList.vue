<template>
   <h2>
      <span>Collections</span>
      <div class="actions" v-if="userStore.isAdmin" >
         <DPGButton label="Create Collection" class="create" @click="createCollection()"/>
      </div>
   </h2>
   <div class="collections">
      <DataTable :value="collectionStore.collections" ref="collectionRecordsTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="false" :paginator="collectionStore.totalCollections > 15"
         :rows="pageSize"  :rowsPerPageOptions="[15,30,100]" :totalRecords="collectionStore.totalRecords"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <Column field="id" header="ID">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
            </template>
         </Column>
         <Column field="pid" header="PID" class="nowrap"/>
         <Column field="title" header="Title" />
         <Column field="creatorName" header="Creator Name" />
         <Column field="barcode" header="Barcode" />
         <Column field="callNumber" header="Call Number" />
         <Column field="catalogKey" header="Catalog Key" />
         <Column field="recordCount" header="Items" />
      </DataTable>
   </div>
   <Dialog v-model:visible="showCreateCollection" :modal="true" header="Create Collection" @hide="createCollectionClosed" :style="{width: '750px'}">
      <NewMetadataPanel @canceled="createCollectionClosed" @created="collectionCreated" :collection="true"/>
   </Dialog>
</template>

<script setup>
import { onMounted, ref} from 'vue'
import { useCollectionsStore } from '../stores/collections'
import { useUserStore } from '../stores/user'
import { useSystemStore } from '../stores/system'
import { useMetadataStore } from '../stores/metadata'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import NewMetadataPanel from '@/components/NewMetadataPanel.vue'

const collectionStore = useCollectionsStore()
const userStore = useUserStore()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()

const pageSize = ref(15)
const showCreateCollection = ref(false)

onMounted(() => {
   collectionStore.getCollections()
   document.title = `Collections`
})

const createCollection = (() => {
   showCreateCollection.value = true
})

const createCollectionClosed = (() => {
   showCreateCollection.value = false
})

const collectionCreated = (() => {
   collectionStore.addCollection(  metadataStore.detail )
   systemStore.toastMessage("Collection Created", `A new collection metadata record has been created.`)
   showCreateCollection.value = false
})
</script>

<style scoped lang="scss">
.collections {
   min-height: 600px;
   text-align: left;
   padding: 25px 40px;
}
</style>