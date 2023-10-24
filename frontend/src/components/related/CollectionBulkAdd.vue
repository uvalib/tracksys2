<template>
   <div class="bulk-toolbar">
      <span class="acts">
         <DPGButton label="Add Selected" class="p-button-secondary" @click="submitClicked()" :disabled="selectedRecords.length == 0" />
         <DPGButton label="Cancel" class="p-button-secondary" @click="cancelClicked()" />
      </span>
      <span class="search">
         <span class="p-input-icon-right">
            <i class="pi pi-search" @click="queryMetadata()"/>
            <InputText v-model="query" placeholder="Metadata Search" @keypress="searchKeyPressed($event)"/>
         </span>
         <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()" :disabled="query.length == 0"/>
      </span>
   </div>
   <div class="metadata-hits" v-if="collectionStore.metadataHits.length > 0">
      <DataTable :value="collectionStore.metadataHits" ref="relatedOrdersTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :rowsPerPageOptions="[15,30,50]" paginatorPosition="top" removableSort
         :lazy="false" :paginator="true" :alwaysShowPaginator="true" :rows="15"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
         v-model:selection="selectedRecords" :selectAll="selectAll" @select-all-change="onSelectAllChange" @row-select="onRowSelect" @row-unselect="onRowUnselect"
      >
         <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
         <Column field="id" header="ID" />
         <Column field="pid" header="PID" />
         <Column field="title" header="Title" />
         <Column field="barcode" header="Barcode" :sortable="true"  />
         <Column field="catalogKey" header="Catalog Key" :sortable="true" />
         <Column field="callNumber" header="Call Number" :sortable="true" />
         <Column field="masterFilesCount" header="Master Files" />
      </DataTable>
   </div>
   <div v-else>
      <p class="help">Search for metadata records by title, barcode, catalog key or call number to include in this collection.</p>
   </div>
</template>

<script setup>
import { ref } from 'vue'
import InputText from 'primevue/inputtext'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useConfirm } from "primevue/useconfirm"
import { useCollectionsStore } from '@/stores/collections'
import { useSystemStore } from '@/stores/system'

const collectionStore = useCollectionsStore()
const systemStore = useSystemStore()
const confirm = useConfirm()

const query = ref("")
const selectedRecords = ref([])
const selectAll = ref(false)

const searchKeyPressed = ((event) => {
   if (event.keyCode == 13) {
      queryMetadata()
   }
})
const queryMetadata = (() => {
   collectionStore.metadataSearch( query.value )
})
const clearSearch = (() => {
   query.value = ""
   collectionStore.metadataHits = []
   selectedRecords.value = []
})
const cancelClicked = (() => {
   selectedRecords.value = []
   collectionStore.toggleBulkAdd()
})
const submitClicked = ( () => {
   let metadataIDs = []
   selectedRecords.value.forEach( s => metadataIDs.push(s.id))
   confirm.require({
      message: `Add the selected metadata records to the collection?`,
      header: 'Confirm Add',
      icon: 'pi pi-question-circle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await collectionStore.addRecords( metadataIDs )
         selectedRecords.value = []
         collectionStore.toggleBulkAdd()
         systemStore.toastMessage("Add Started", "The selected items are being added to the collection. Check the job status page for updates.")
      }
   })
})
const onRowSelect = (() => {
   selectAll.value = selectedRecords.value.length < collectionStore.metadataHits.length
})
const onRowUnselect = (() => {
   selectAll.value  = false
})
const onSelectAllChange = ((event) => {
   selectAll.value = event.checked
   if (selectAll.value) {
      selectedRecords.value = collectionStore.metadataHits
   }
   else {
      selectedRecords.value = []
   }
})

</script>

<stype scoped lang="scss">
.bulk-toolbar {
   padding: 0 0 10px 0;
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   button.p-button {
      margin-left: 5px;
   }
}
p.help {
   font-size: 1.2em;
   text-align: center;
}
</stype>