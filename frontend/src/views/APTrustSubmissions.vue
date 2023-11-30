<template>
    <h2>APTrust Submissions</h2>
    <div class="submissions">
      <DataTable :value="apTrust.submissions" ref="apTrustTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll"
         :sortField="apTrust.searchOpts.sortField" :sortOrder="sortOrder" @sort="onSort($event)"
         :lazy="true" :paginator="true" @page="onPage($event)"  paginatorPosition="top"
         :rows="apTrust.searchOpts.limit" :totalRecords="apTrust.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[30,50,100]" :first="apTrust.searchOpts.start"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <template #paginatorstart></template>
         <template #paginatorend>
            <span class="p-input-icon-right">
               <i class="pi pi-search" />
               <InputText v-model="apTrust.searchOpts.query" placeholder="Submission Search" @input="apTrust.getSubmissions()"/>
            </span>
            <DPGButton label="Clear" class="p-button-secondary pad" @click="clearSearch()"/>
         </template>
         <Column field="id" header="ID" :sortable="true" >
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
            </template>
         </Column>
         <Column field="pid" header="PID" :sortable="true"  class="nowrap"/>
         <Column field="title" header="Title" :sortable="true" />
         <Column field="requestedAt" header="Date Requested" :sortable="true" class="nowrap">
            <template #body="slotProps">{{ formatDate(slotProps.data.requestedAt) }}</template>
         </Column>
         <Column field="processedAt" header="Date Processed" :sortable="true" class="nowrap" >
            <template #body="slotProps">{{ formatDate(slotProps.data.processedAt) }}</template>
         </Column>
         <Column field="success" header="Status" :sortable="true" class="apt-status">
            <template #body="slotProps">
               <template v-if="slotProps.data.processedAt">
                  <span v-if="slotProps.data.success" class="pi pi-check-circle success"></span>
                  <span v-else class="pi pi-times-circle fail"></span>
               </template>
               <span v-else class="pi pi-spin pi-cog"></span>
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
import { useAPTrustStore } from '@/stores/aptrust'
import { usePinnable } from '@/composables/pin'
import dayjs from 'dayjs'

usePinnable("p-paginator-top")

const apTrust = useAPTrustStore()

const sortOrder = computed(() => {
   if (apTrust.searchOpts.sortOrder == "desc") {
      return -1
   }
   return 1
})

onMounted(() => {
   apTrust.getSubmissions()
   document.title = `APTrust Submissions`
})

const formatDate = ( ( dateStr ) => {
   if (dateStr) {
      let d = dayjs(dateStr)
      return d.format("YYYY-MM-DD HH:mm")
   }
   return ""
})

const clearSearch = (() => {
   apTrust.searchOpts.query = ""
   apTrust.getSubmissions()
})

const onPage = ((event) => {
   apTrust.searchOpts.start = event.first
   apTrust.searchOpts.limit = event.rows
   apTrust.getSubmissions()
})

const onSort = ((event) => {
   apTrust.searchOpts.sortField = event.sortField
   apTrust.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      apTrust.searchOpts.sortOrder = "desc"
   }
   apTrust.getSubmissions()
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
   :deep(td.apt-status) {
      text-align: center;
      padding: 0;
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
}
:deep(td.nowrap) {
   white-space: nowrap;
}
</style>