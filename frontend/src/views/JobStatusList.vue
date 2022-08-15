<template>
   <h2>Job Statuses</h2>
   <div class="job-status">
      <DataTable :value="jobsStore.jobs" ref="jobsTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll"
         :lazy="true" :paginator="true" @page="onPage($event)" :rowClass="rowClass"
         :rows="jobsStore.searchOpts.limit" :totalRecords="jobsStore.totalJobs"
         v-model:selection="selectedJobs" :selectAll="selectAll" @select-all-change="onSelectAllChange" @row-select="onRowSelect" @row-unselect="onRowUnselect"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[10,30,100]"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
         <Column field="name" header="Job Type"></Column>
         <Column field="associatedObject" header="Associated Object"></Column>
         <Column field="status" header="Status"></Column>
         <Column field="warnings" header="Warnings"></Column>
         <Column field="startedAt" header="Started"></Column>
         <Column field="finishedAt" header="Finished"></Column>
         <Column header="">
            <template #body="slotProps">
               <router-link :to="`/jobs/${slotProps.data.id}`">View Log</router-link>
            </template>
         </Column>
      </DataTable>
   </div>
</template>

<script setup>
import { onMounted, ref} from 'vue'
import { useJobsStore } from '@/stores/jobs'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const jobsStore = useJobsStore()

const selectedJobs = ref()
const selectAll = ref(false)

function rowClass(rowData) {
   if (rowData.status ==  "failure"){
      return "error-row"
   }
   if (rowData.status ==  "running"){
      return "running-row"
   }
   if (rowData.status ==  "warn"){
      return "warn-row"
   }
   return ""
}

function onPage(event) {
   jobsStore.searchOpts.start = event.first
   jobsStore.searchOpts.limit = event.rows
   jobsStore.getJobs()
}

function onRowSelect() {
   selectAll.value = selectedJobs.value === jobsStore.searchOpts.limit
}
function onRowUnselect() {
   selectAll.value  = false
}
function onSelectAllChange(event) {
   selectAll.value = event.checked
   if (selectAll.value) {
      selectedJobs.value = jobsStore.jobs
   }
   else {
      selectedJobs.value = []
   }
}

onMounted(() => {
   jobsStore.getJobs()
})
</script>

<style scoped lang="scss">
   .job-status {
      min-height: 600px;
      text-align: left;
      padding: 0 25px;
      .p-datatable {
         font-size: 0.85em;
         :deep(td), :deep(th) {
            padding: 10px;
         }
      }
      :deep(.error-row) {
         background-color: #944 !important;
         color: #fff;
         &:hover {
            background-color: #a44 !important;
         }
      }
      :deep(.running-row)  {
         background-color: var(--uvalib-blue-alt-light) !important;
         &:hover {
            background-color: #def !important;
         }
      }
      :deep(.warn-row)  {
         background-color: var(--uvalib-yellow-light) !important;
      }
   }
</style>