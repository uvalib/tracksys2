<template>
   <div v-if="collectionStore.working">
   </div>
   <div v-if="collectionStore.working == false && collectionStore.total == 0">
      <h3>No collection records found</h3>
   </div>
   <div v-else class="results">
      <div class="toolbar">
         <span>
            <span class="p-input-icon-right">
               <i class="pi pi-search" />
               <InputText v-model="collectionStore.searchOpts.query" placeholder="Collection Search" @input="queryCollection()"/>
            </span>
            <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()" :disabled="collectionStore.searchOpts.query.length == 0"/>
         </span>
      </div>
      <DataTable :value="collectionStore.records" ref="collectionRecordsTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="true" :paginator="collectionStore.total > 15" @page="onCollectionPage($event)"
         :rows="collectionStore.searchOpts.limit" :totalRecords="collectionStore.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[15,30,100]" :first="collectionStore.searchOpts.start"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <Column field="id" header="ID">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
            </template>
         </Column>
         <Column field="pid" header="PID" />
         <Column field="title" header="Title" />
      </DataTable>
   </div>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import { onBeforeMount } from 'vue'
import { useCollectionsStore } from '@/stores/collections'

const props = defineProps({
   collectionID: {
      type: Number,
      required: true
   }
})

const collectionStore = useCollectionsStore()

onBeforeMount( async () => {
   collectionStore.setCollection( props.collectionID )
   await collectionStore.getRecords()
})

function onCollectionPage(event) {
   collectionStore.searchOpts.start = event.first
   collectionStore.searchOpts.limit = event.rows
   collectionStore.getRecords()
}
function queryCollection() {
   collectionStore.getRecords()
}
function clearSearch() {
   collectionStore.searchOpts.query = ""
   collectionStore.getRecords()
}

</script>

<stype scoped lang="scss">
.results {
   margin: 0;
   font-size: 0.9em;
   td.nowrap, th {
      white-space: nowrap;
   }
   th, td {
      font-size: 0.85em;
   }
}
.toolbar {
   padding: 0 0 10px 0;
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
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