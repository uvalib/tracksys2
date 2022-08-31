<template>
   <ConfirmDialog position="top"/>
   <h2>Job Statuses</h2>
   <div class="job-status">
      <div class="toolbar">
         <DPGButton label="Delete selected" :disabled="selectedJobs.length == 0"  class="p-button-secondary" @click="deletAllClicked"/>
      </div>
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
         <Column header="" class="row-acts">
            <template #body="slotProps">
               <router-link :to="`/jobs/${slotProps.data.id}`">View</router-link>
               <DPGButton label="Delete"  class="p-button-text" @click="deleteJob(slotProps.data.id)"/>
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
import { useConfirm } from "primevue/useconfirm"

const jobsStore = useJobsStore()
const confirm = useConfirm()

const selectedJobs = ref([])
const selectAll = ref(false)

function deletAllClicked() {
   confirm.require({
      message: 'Are you sure you want delete the selected job status records? All data will be lost. This cannot be reversed.',
      header: 'Confirm Delete All',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: () => {
         let payload = []
         selectedJobs.value.forEach( j => payload.push(j.id))
         jobsStore.deleteJobs( payload )
      }
   })
}
function deleteJob(id) {
   confirm.require({
      message: 'Are you sure you want delete this job status?',
      header: 'Confirm Delete All',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: () => {
         jobsStore.deleteJobs( [id] )
      }
   })
}

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
   // jobsStore.getJobs()
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
      .toolbar {
         padding: 10px 0;
         font-size: 0.8em;
      }
      .p-datatable {
         font-size: 0.85em;
         :deep(td), :deep(th) {
            padding: 10px;
         }
         :deep(.row-acts) {
            text-align: center;
            padding: 0;
            a {
               display: inline-block;
               margin-right: 15px;
            };
         }
      }
      :deep(.error-row) {
         background-color: #944 !important;
         color: #fff;
         .row-acts {
            a, button.p-button-text {
               color: #fff !important;
            }
         }
         &:hover {
            background-color: #a44 !important;
         }
      }
      :deep(tr.p-highlight.error-row) {
         color: white !important;
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