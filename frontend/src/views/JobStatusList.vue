<template>
   <h2>Job Statuses</h2>
   <div class="job-status">
      <div class="toolbar">
         <DPGButton label="Delete selected" :disabled="selectedJobs.length == 0"  class="p-button-secondary" @click="deletAllClicked"/>
         <span>
            <span class="p-input-icon-right">
               <i class="pi pi-search" />
               <InputText v-model="jobsStore.searchOpts.query" placeholder="Job Status Search" @input="queryJobs()"/>
            </span>
            <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()" :disabled="jobsStore.searchOpts.query.length == 0"/>
         </span>
      </div>
      <DataTable :value="jobsStore.jobs" ref="jobsTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll"
         :lazy="true" :paginator="true" @page="onPage($event)" :rowClass="rowClass"
         :rows="jobsStore.searchOpts.limit" :totalRecords="jobsStore.totalJobs"
         v-model:selection="selectedJobs" :selectAll="selectAll" @select-all-change="onSelectAllChange" @row-select="onRowSelect" @row-unselect="onRowUnselect"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[10,30,100]" :first="jobsStore.searchOpts.start"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
         <Column field="name" header="Job Type"></Column>
         <Column field="associatedObject" header="Associated Object">
            <template #body="slotProps">
               <template v-if="getAssociatedObjectLink(slotProps.data.associatedObject)">
                  <router-link :to="getAssociatedObjectLink(slotProps.data.associatedObject)">{{slotProps.data.associatedObject}}</router-link>
               </template>
               <template v-else>
                  {{slotProps.data.associatedObject}}
               </template>
            </template>
         </Column>
         <Column field="status" header="Status"></Column>
         <Column field="warnings" header="Warnings"></Column>
         <Column field="startedAt" header="Started"></Column>
         <Column field="finishedAt" header="Finished"></Column>
         <Column header="" class="row-acts">
            <template #body="slotProps">
               <router-link :to="`/jobs/${slotProps.data.id}`">View</router-link>
               <span class="sep">|</span>
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
import InputText from 'primevue/inputtext'
import { useConfirm } from "primevue/useconfirm"

const jobsStore = useJobsStore()
const confirm = useConfirm()

const selectedJobs = ref([])
const selectAll = ref(false)

function queryJobs() {
   jobsStore.getJobs()
}
function clearSearch() {
   jobsStore.searchOpts.query = ""
   jobsStore.getJobs()
}
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

function getAssociatedObjectLink( objName ) {
   if (objName.split(" ").length != 2) {
      return ""
   }
   let objType = objName.split(" ")[0].toLowerCase().trim()
   let objID =  objName.split(" ")[1].toLowerCase().trim()
   if (objType == "unit") {
      return `/units/${objID}`
   }
   if (objType == "order") {
      return `/orders/${objID}`
   }
   if (objType == "metadata") {
      return `/metadata/${objID}`
   }
   return ""
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
   document.title = `Job Statuses`
})
</script>

<style scoped lang="scss">
   .job-status {
      min-height: 600px;
      text-align: left;
      padding: 0 25px;
      .sep {
         display: inline-block;
         margin: 0 10px;
      }
      .toolbar {
         padding: 0 0 10px 0;
         display: flex;
         flex-flow: row nowrap;
         justify-content: space-between;
         label {
            font-weight: bold;
            margin-right: 5px;
            display: inline-block;
         }
         button.p-button {
            margin-left: 5px;
         }
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
            };
         }
      }
      :deep(.error-row) {
         background-color: #944 !important;
         color: #fff;
         a {
            color: #fff !important;
         }
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