<template>
   <h2>HathiTrust Submissions</h2>
   <div class="submissions">
      <DataTable :value="hathiTrust.submissions" ref="hathiTrustTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll"
         :sortField="hathiTrust.searchOpts.sortField" :sortOrder="sortOrder" @sort="onSort($event)"
         :lazy="true" :paginator="true" @page="onPage($event)"  paginatorPosition="top"
         :rows="hathiTrust.searchOpts.limit" :totalRecords="hathiTrust.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[30,50,100]" :first="hathiTrust.searchOpts.start"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
         :loading="hathiTrust.working"
         v-model:filters="columnFilters" filterDisplay="menu" @filter="getSubmissions()"
      >
         <template #paginatorstart></template>
         <template #paginatorend>
            <span class="p-input-icon-right">
               <i class="pi pi-search" />
               <InputText v-model="hathiTrust.searchOpts.query" placeholder="Submission Search" @input="hathiTrust.getSubmissions()"/>
            </span>
            <DPGButton label="Clear" class="p-button-secondary pad" @click="clearSearch()"/>
         </template>
         <Column field="pid" header="PID" :sortable="true"  class="nowrap">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.metadataID}`">{{slotProps.data.metadata.pid}}</router-link>
            </template>
         </Column>
         <Column field="title" header="Title" :sortable="true">
            <template #body="slotProps">{{ slotProps.data.metadata.title }}</template>
         </Column>
         <Column field="barcode" header="Barcode" :sortable="true">
            <template #body="slotProps">{{ slotProps.data.metadata.barcode }}</template>
         </Column>
         <Column field="requested_at" header="Requested" :sortable="true" class="nowrap" >
            <template #body="slotProps">{{ $formatDate(slotProps.data.requestedAt, false) }}</template>
         </Column>
         <Column field="metadata_submitted_at" header="Metadata Submitted" :sortable="true" class="nowrap" >
            <template #body="slotProps">
               <span v-if="slotProps.data.metadataSubmittedAt">{{ $formatDate(slotProps.data.metadataSubmittedAt, false) }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
         <Column field="metadataStatus" header="Metadata Status" class="nowrap" filterField="metadataStatus" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <Dropdown v-model="filterModel.value" :options="statuses" optionLabel="name" optionValue="value" placeholder="Select status" />
            </template>
         </Column>

         <Column field="package_created_at" header="Package Created" :sortable="true" class="nowrap" >
            <template #body="slotProps">
               <span v-if="slotProps.data.packageCreatedAt">{{ $formatDate(slotProps.data.packageCreatedAt, false) }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
         <Column field="package_submitted_at" header="Package Submitted" :sortable="true" class="nowrap" >
            <template #body="slotProps">
               <span v-if="slotProps.data.packageSubmittedAt">{{ $formatDate(slotProps.data.packageSubmittedAt, false) }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
         <Column field="packageStatus" header="Package Status" class="nowrap" filterField="packageStatus" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <Dropdown v-model="filterModel.value" :options="statuses" optionLabel="name" optionValue="value" placeholder="Select status" />
            </template>
         </Column>

         <Column field="finished_at" header="Finished" :sortable="true" class="nowrap" >
            <template #body="slotProps">
               <span v-if="slotProps.data.finishedAt">{{ $formatDate(slotProps.data.finishedAt, false) }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
         <Column field="notes" header="Notes" >
            <template #body="slotProps">
               <DPGButton v-if="slotProps.data.notes"  class="notes" label="View" severity="secondary" @click="notesClicked(slotProps.data)"/>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
      </DataTable>
   </div>
   <Dialog v-model:visible="showNotes" :modal="true" header="Submission Notes">
      <div>{{ tgtSubmission.metadata.pid }} - {{ tgtSubmission.metadata.title }}</div>
      <div class="note-text">{{ tgtSubmission.notes }}</div>
   </Dialog>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import Dropdown from 'primevue/dropdown'
import { FilterMatchMode } from 'primevue/api'
import { useHathiTrustStore } from '@/stores/hathitrust'
import { usePinnable } from '@/composables/pin'

usePinnable("p-paginator-top")

const hathiTrust = useHathiTrustStore()

const showNotes = ref(false)
const tgtSubmission = ref()

const columnFilters = ref( {
   'metadataStatus': {value: null, matchMode: FilterMatchMode.EQUALS},
   'packageStatus': {value: null, matchMode: FilterMatchMode.EQUALS},
})

const statuses = ref([
   {name: "Submitted", value: "submitted"},
   {name: "Accepted", value: "accepted"},
   {name: "Failed", value: "failed"},
])

const sortOrder = computed(() => {
   if (hathiTrust.searchOpts.sortOrder == "desc") {
      return -1
   }
   return 1
})

onMounted(() => {
   getSubmissions()
   document.title = `HathiTrust Submissions`
})

const getSubmissions = (() => {
   hathiTrust.searchOpts.filters = []
   Object.entries(columnFilters.value).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         hathiTrust.searchOpts.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   hathiTrust.getSubmissions()
})

const notesClicked = ( (s) => {
   tgtSubmission.value = s
   showNotes.value = true
})

const clearSearch = (() => {
   hathiTrust.searchOpts.query = ""
   getSubmissions()
})

const onPage = ((event) => {
   hathiTrust.searchOpts.start = event.first
   hathiTrust.searchOpts.limit = event.rows
   getSubmissions()
})

const onSort = ((event) => {
   hathiTrust.searchOpts.sortField = event.sortField
   hathiTrust.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      hathiTrust.searchOpts.sortOrder = "desc"
   }
   getSubmissions()
})

</script>

<style scoped lang="scss">
.submissions {
   min-height: 600px;
   text-align: left;
   padding: 0 25px;
   button.pad {
      margin-left: 10px;
   }
   button.notes {
      font-size: 0.75em;
      padding: 5px 10px;
   }
   .none {
      color: var(--uvalib-grey-light);
      font-style: italic;
   }
}
div.note-text {
   height: 150px;
   overflow-y: scroll;
   margin: 15px 0 10px 0;
   padding: 5px 10px;
   border: 1px solid var(--uvalib-grey-light);
   border-radius: 5px;
   background-color: white;
}
</style>