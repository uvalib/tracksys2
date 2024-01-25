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
         <Column field="metadataStatus" header="Metadata Status" class="nowrap" ></Column>

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
         <Column field="packageStatus" header="Package Status" class="nowrap" ></Column>

         <Column field="finished_at" header="Finished" :sortable="true" class="nowrap" >
            <template #body="slotProps">
               <span v-if="slotProps.data.finishedAt">{{ $formatDate(slotProps.data.finishedAt, false) }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
      </DataTable>
   </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import { useHathiTrustStore } from '@/stores/hathitrust'
import { usePinnable } from '@/composables/pin'

usePinnable("p-paginator-top")

const hathiTrust = useHathiTrustStore()

const sortOrder = computed(() => {
   if (hathiTrust.searchOpts.sortOrder == "desc") {
      return -1
   }
   return 1
})

onMounted(() => {
   hathiTrust.getSubmissions()
   document.title = `HathiTrust Submissions`
})

const clearSearch = (() => {
   hathiTrust.searchOpts.query = ""
   hathiTrust.getSubmissions()
})

const onPage = ((event) => {
   hathiTrust.searchOpts.start = event.first
   hathiTrust.searchOpts.limit = event.rows
   hathiTrust.getSubmissions()
})

const onSort = ((event) => {
   hathiTrust.searchOpts.sortField = event.sortField
   hathiTrust.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      hathiTrust.searchOpts.sortOrder = "desc"
   }
   hathiTrust.getSubmissions()
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
   .none {
      color: var(--uvalib-grey-light);
      font-style: italic;
   }
}
</style>