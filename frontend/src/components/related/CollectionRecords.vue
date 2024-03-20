<template>
   <div class="collection">
      <template v-if="collectionStore.bulkAdd == false">
         <DataTable :value="collectionStore.records" ref="collectionRecordsTable" dataKey="id"
            stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
            :sortField="collectionStore.searchOpts.sortField" :sortOrder="sortOrder" @sort="onSort($event)"
            :lazy="true" :paginator="true" @page="onCollectionPage($event)" paginatorPosition="top"
            :rows="collectionStore.searchOpts.limit" :totalRecords="collectionStore.totalRecords"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            :rowsPerPageOptions="[15,30,100]" :first="collectionStore.searchOpts.start"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
         >
            <template #empty><h3>No items found</h3></template>
            <template #paginatorstart>
               <DPGButton label="Add Item(s)" class="p-button-secondary" @click="bulkAddClicked()" v-if="userStore.isAdmin"/>
               <DPGButton label="Export" class="p-button-secondary pad-left" @click="exportCollection" :disabled="collectionStore.totalRecords == 0"/>
            </template>
            <template #paginatorend>
               <IconField iconPosition="left">
                  <InputIcon class="pi pi-search" />
                  <InputText v-model="collectionStore.searchOpts.query" placeholder="Collection Search" @input="queryCollection()"/>
               </IconField>
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
            <Column v-if="collectionStore.inAPTrust" header="APTrust" class="apt-status" field="aptStatus" :sortable="true">
               <template #body="slotProps">
                  <template v-if="slotProps.data.apTrustSubmission">
                     <template v-if="slotProps.data.apTrustSubmission.processedAt">
                        <span v-if="slotProps.data.apTrustSubmission.success" class="pi pi-check-circle success"></span>
                        <span v-else class="pi pi-times-circle fail"></span>
                     </template>
                     <span v-else class="pi pi-spin pi-cog"></span>
                  </template>
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
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import { useCollectionsStore } from '@/stores/collections'
import { useConfirm } from "primevue/useconfirm"
import { useUserStore } from '@/stores/user'
import CollectionBulkAdd from './CollectionBulkAdd.vue'
import { usePinnable } from '@/composables/pin'
import { computed } from 'vue'

usePinnable("p-paginator-top")

const userStore = useUserStore()
const collectionStore = useCollectionsStore()
const confirm = useConfirm()

const props = defineProps({
   collectionID: {
      type: Number,
      required: true
   }
})

const sortOrder = computed(() => {
   if (collectionStore.searchOpts.sortOrder == "desc") {
      return -1
   }
   return 1
})

const onSort = ((event) => {
   collectionStore.searchOpts.sortField = event.sortField
   collectionStore.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      collectionStore.searchOpts.sortOrder = "desc"
   }
   collectionStore.getItems()
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
   collectionStore.getItems( false )
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
   h3  {
      text-align: center;
   }
   .none {
      color: var(--uvalib-grey-light);
      font-style: italic;
   }
   td.nowrap, th {
      white-space: nowrap;
   }
   td.apt-status {
      text-align: center;
      span.pi {
         font-size: 1.2em;
         font-weight: bold;
      }
      span.pi.success {
         color: var(--uvalib-green);
      }
      span.pi.fail {
         color: var(--uvalib-red-dark);
      }
   }
   .none{
      text-align: center;
   }
   td.row-acts {
      text-align: center;
      width: 25px;
   }
   button.pad-left {
      margin-left: 10px;
   }
}
</stype>