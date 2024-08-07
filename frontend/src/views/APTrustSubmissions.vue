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
         :loading="apTrust.working"
      >
         <template #paginatorstart></template>
         <template #paginatorend>
            <IconField iconPosition="left">
               <InputIcon class="pi pi-search" />
               <InputText v-model="apTrust.searchOpts.query" placeholder="Search Submissions" @input="apTrust.getSubmissions(false)"/>
            </IconField>
         </template>
         <Column field="pid" header="PID" :sortable="true"  class="nowrap">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.metadataID}`">{{slotProps.data.pid}}</router-link>
            </template>
         </Column>
         <Column field="title" header="Title" :sortable="true" />
         <Column field="requestedAt" header="Requested" :sortable="true" class="nowrap">
            <template #body="slotProps">{{ $formatDateTime(slotProps.data.requestedAt) }}</template>
         </Column>
         <Column field="processedAt" header="Processed" :sortable="true" class="nowrap" >
            <template #body="slotProps">{{ $formatDateTime(slotProps.data.processedAt) }}</template>
         </Column>
         <Column field="success" header="Status" :sortable="true" class="apt-status">
            <template #body="slotProps">
               <template v-if="slotProps.data.processedAt">
                  <span v-if="slotProps.data.success" class="pi pi-check-circle success"></span>
                  <span v-else class="pi pi-times-circle fail"></span>
               </template>
               <span v-else class="pi pi-spin pi-cog"></span>
               <span class="pi pi-info-circle info" @click="infoClicked(slotProps.data)"
                  v-tooltip.left="{ value: 'View submisson status', showDelay: 250 }"></span>
            </template>
         </Column>
      </DataTable>
      <Dialog v-model:visible="showDialog" :header="`Submission Status for ${tgtPID}`" :modal="true" position="top"  >
         <APTrustPanel :readonly="true"/>
      </Dialog>
   </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import { useAPTrustStore } from '@/stores/aptrust'
import { usePinnable } from '@/composables/pin'
import APTrustPanel from '@/components/aptrust/APTrustPanel.vue'
import Dialog from 'primevue/dialog'

usePinnable("p-datatable-paginator-top")

const apTrust = useAPTrustStore()
const showDialog = ref(false)
const tgtPID = ref("")

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

const infoClicked = ((submission) => {
   tgtPID.value = submission.pid
   showDialog.value = true
   apTrust.getItemStatus( submission.metadataID )
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
      span.pi.info {
         display: inline-block;
         margin-left: 15px;
         color: var(--uvalib-blue-alt);
         font-weight: normal;
         cursor: pointer;
         &:hover {
            font-weight: bold;
         }
      }
   }
}
:deep(td.nowrap) {
   white-space: nowrap;
}
</style>