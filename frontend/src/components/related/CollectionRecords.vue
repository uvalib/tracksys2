<template>
   <WaitSpinner v-if="collectionStore.working" :overlay="true" message="Please wait..." />
   <div class="collection">
      <template v-if="collectionStore.bulkAdd == false">
         <div class="toolbar">
            <span class="left">
               <LookupDialog label="Add Item" @selected="addItem" target="metadata" :create="false" :collection="collectionStore.collectionID" v-if="userStore.isAdmin"/>
               <DPGButton label="Bulk Add" class="p-button-secondary" @click="bulkAddClicked()" v-if="userStore.isAdmin"/>
               <DPGButton label="Export" class="p-button-secondary" @click="exportCollection" :disabled="collectionStore.totalRecords == 0"/>
            </span>
            <span class="search">
               <span class="p-input-icon-right">
                  <i class="pi pi-search" />
                  <InputText v-model="collectionStore.searchOpts.query" placeholder="Collection Search" @input="queryCollection()"/>
               </span>
               <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()" :disabled="collectionStore.searchOpts.query.length == 0"/>
            </span>
         </div>
         <div v-if="collectionStore.working == false && collectionStore.totalRecords == 0" class="none">
            <h3>No items found</h3>
         </div>
         <DataTable v-if="collectionStore.totalRecords>0" :value="collectionStore.records" ref="collectionRecordsTable" dataKey="id"
            stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
            :lazy="true" :paginator="collectionStore.totalRecords > 15" @page="onCollectionPage($event)"
            :rows="collectionStore.searchOpts.limit" :totalRecords="collectionStore.totalRecords"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            :rowsPerPageOptions="[15,30,100]" :first="collectionStore.searchOpts.start"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
         >
            <Column field="id" header="ID">
               <template #body="slotProps">
                  <router-link :to="`/metadata/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
               </template>
            </Column>
            <Column field="pid" header="PID" class="nowrap"/>
            <Column field="title" header="Title" />
            <Column header="" class="row-acts nowrap" v-if="userStore.isAdmin">
               <template #body="slotProps">
                  <DPGButton icon="pi pi-times" class="p-button-rounded p-button-text p-button-secondary" @click="deleteItem(slotProps.data)"/>
               </template>
            </Column>
         </DataTable>
      </template>
      <CollectionBulkAdd v-else />
   </div>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import { onBeforeMount } from 'vue'
import { useCollectionsStore } from '@/stores/collections'
import { useConfirm } from "primevue/useconfirm"
import { useUserStore } from '@/stores/user'
import LookupDialog from '@/components/LookupDialog.vue'
import WaitSpinner from '@/components/WaitSpinner.vue'
import CollectionBulkAdd from './CollectionBulkAdd.vue'

const userStore = useUserStore()
const collectionStore = useCollectionsStore()

const confirm = useConfirm()

const props = defineProps({
   collectionID: {
      type: Number,
      required: true
   }
})

onBeforeMount( async () => {
   collectionStore.setCollection( props.collectionID )
   await collectionStore.getItems()
})

const bulkAddClicked = (() => {
   collectionStore.toggleBulkAdd()
})

const onCollectionPage = ((event) => {
   collectionStore.searchOpts.start = event.first
   collectionStore.searchOpts.limit = event.rows
   collectionStore.getItems()
})

const queryCollection = (() => {
   collectionStore.getItems()
})
const clearSearch = (() => {
   collectionStore.searchOpts.query = ""
   collectionStore.getItems()
})
const addItem = (( metadataID ) => {
  collectionStore.addItems( [metadataID] )
})
const deleteItem = (( item ) => {
   confirm.require({
      message: `Remove "${item.pid} : ${item.title}" from this collection?`,
      header: 'Confirm Remove Item',
      icon: 'pi pi-question-circle',
      rejectClass: 'p-button-secondary',
      accept: () => {
         collectionStore.removeItem(item)
      }
   })
})
const exportCollection = (() => {
   collectionStore.exportCSV()
})

</script>

<stype scoped lang="scss">
.collection  {
   margin: 0;
   td.nowrap, th {
      white-space: nowrap;
   }
   .none{
      text-align: center;
   }
   td.row-acts {
      text-align: center;
      width: 25px;
   }
}
.toolbar {
   padding: 0 0 10px 0;
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   .search {
      margin-left: auto;
   }
   label {
      font-weight: bold;
      margin-right: 5px;
      display: inline-block;
   }
   button.p-button {
      margin-left: 5px;
   }
}
</stype>