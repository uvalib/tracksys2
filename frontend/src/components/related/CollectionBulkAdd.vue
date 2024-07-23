<template>
   <div class="bulk-toolbar">
      <span class="acts">
         <DPGButton label="Add Selected" class="p-button-secondary" @click="submitClicked()" :disabled="selectedRecords.length == 0" />
         <DPGButton label="Cancel" class="p-button-secondary" @click="cancelClicked()" />
      </span>
      <span class="collection-search">
         <IconField iconPosition="left">
            <InputIcon class="pi pi-search" />
            <InputText v-model="query" placeholder="Search" @keypress="searchKeyPressed($event)"/>
         </IconField>
         <DPGButton label="Search" @click="queryMetadata()" :disabled="query.length == 0"/>
         <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()" :disabled="query.length == 0"/>
      </span>
   </div>
   <div class="metadata-hits">
      <div v-if="collectionStore.searching" class="searching">
         <WaitSpinner :overlay="false" message="Searching..." />
      </div>
      <DataTable v-else :value="collectionStore.metadataHits" ref="relatedOrdersTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :rowsPerPageOptions="[15,30,50]" paginatorPosition="top" removableSort
         :lazy="false" :paginator="true" :alwaysShowPaginator="true" :rows="15"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
         v-model:selection="selectedRecords" :selectAll="selectAll" @select-all-change="onSelectAllChange" @row-select="onRowSelect" @row-unselect="onRowUnselect"
      >
         <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
         <Column field="id" header="ID" />
         <Column field="pid" header="PID" class="nowrap"/>
         <Column field="type" header="Type" />
         <Column field="title" header="Title" />
         <Column field="barcode" header="Barcode" :sortable="true">
            <template #body="slotProps">
               <span v-if="slotProps.data.barcode">{{ slotProps.data.barcode }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
         <Column field="catalogKey" header="Catalog Key" :sortable="true">
            <template #body="slotProps">
               <span v-if="slotProps.data.catalogKey">{{ slotProps.data.catalogKey }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
         <Column field="callNumber" header="Call Number" :sortable="true">
            <template #body="slotProps">
               <span v-if="slotProps.data.callNumber">{{ slotProps.data.callNumber }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
      </DataTable>
   </div>
</template>

<script setup>
import { ref } from 'vue'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useConfirm } from "primevue/useconfirm"
import { useCollectionsStore } from '@/stores/collections'
import { useSystemStore } from '@/stores/system'
import WaitSpinner from "@/components/WaitSpinner.vue"

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
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Add'
      },
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
.collection-search {
   display: flex;
   flex-flow: row nowrap;
}
td.nowrap, th {
   white-space: nowrap;
}
.none {
   color: var(--uvalib-grey-light);
   font-style: italic;
}
div.searching {
   font-size: 0.9em;
   text-align: center;
}
p.help {
   font-size: 1.2em;
   text-align: center;
}
</stype>