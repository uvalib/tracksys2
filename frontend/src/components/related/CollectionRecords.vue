<template>
   <div class="collection">
      <template v-if="collectionStore.bulkAdd == false">
         <DataTable :value="collectionStore.records" ref="collectionRecordsTable" dataKey="id"
            removableSort stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
            :lazy="true" :paginator="collectionStore.totalRecords > 15" @page="onCollectionPage($event)" paginatorPosition="top"
            :rows="collectionStore.searchOpts.limit" :totalRecords="collectionStore.totalRecords"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            :rowsPerPageOptions="[15,30,100]" :first="collectionStore.searchOpts.start"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
         >
            <template #paginatorstart>
               <div class="toolbar">
                  <span class="left">
                     <DPGButton label="Add Item(s)" class="p-button-secondary" @click="bulkAddClicked()" v-if="userStore.isAdmin"/>
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
            </template>
            <Column field="id" header="ID" :sortable="true">
               <template #body="slotProps">
                  <router-link :to="`/metadata/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
               </template>
            </Column>
            <Column field="pid" header="PID" class="nowrap"/>
            <Column field="title" header="Title" :sortable="true"/>
            <Column field="callNumber" header="Call Number" :sortable="true">
               <template #body="slotProps">
                  <span v-if="slotProps.data.callNumber">{{ slotProps.data.callNumber }}</span>
                  <span v-else class="none">N/A</span>
               </template>
            </Column>
            <Column field="barcode" header="Barcode" :sortable="true">
               <template #body="slotProps">
                  <span v-if="slotProps.data.barcode">{{ slotProps.data.barcode }}</span>
                  <span v-else class="none">N/A</span>
               </template>
            </Column>
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
import { onMounted, onUnmounted, ref } from 'vue'
import { useCollectionsStore } from '@/stores/collections'
import { useConfirm } from "primevue/useconfirm"
import { useUserStore } from '@/stores/user'
import CollectionBulkAdd from './CollectionBulkAdd.vue'

const userStore = useUserStore()
const collectionStore = useCollectionsStore()
const confirm = useConfirm()

const toolbarTop = ref(0)
const toolbarHeight = ref(0)
const toolbarWidth = ref(0)
const toolbar = ref(null)

const props = defineProps({
   collectionID: {
      type: Number,
      required: true
   }
})

onMounted(() => {
   let tb = null
   let tbs = document.getElementsByClassName("p-paginator-top")
   if ( tbs ) {
      tb = tbs[0]
   }
   if ( tb) {
      toolbar.value = tb
      toolbarHeight.value = tb.offsetHeight
      toolbarWidth.value = tb.offsetWidth
      toolbarTop.value = tb.getBoundingClientRect().top
      window.addEventListener("scroll", scrollHandler)
   }
})

onUnmounted(() => {
   window.removeEventListener("scroll", scrollHandler)
})

const scrollHandler = (( ) => {
   if ( toolbar.value) {
      if ( window.scrollY <= toolbarTop.value ) {
         if ( toolbar.value.classList.contains("sticky") ) {
            toolbar.value.classList.remove("sticky")
            let dts = document.getElementsByClassName("p-datatable-wrapper")
            if ( dts ) {
               dts[0].style.top = `0px`
            }
         }
      } else {
         if ( toolbar.value.classList.contains("sticky") == false ) {
            let dts = document.getElementsByClassName("p-datatable-wrapper")
            if ( dts ) {
               dts[0].style.top = `${toolbarHeight.value}px`
            }
            toolbar.value.classList.add("sticky")
            toolbar.value.style.width = `${toolbarWidth.value}px`
         }
      }
   }
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
   .none {
      color: var(--uvalib-grey-light);
      font-style: italic;
   }
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
   padding: 0;
   .search {
      button.p-button {
         margin-left: 10px;
      }
   }
   button.p-button {
      margin-right: 10px;
   }
}
</stype>