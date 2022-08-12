<template>
   <h2>Job Statuses</h2>
   <div class="job-status">
      <EasyDataTable
         ref="jobsTable"
         v-model:server-options="jobsStore.searchOpts"
         :server-items-length="jobsStore.totalJobs"
         table-class-name="jobs-table"
         :body-row-class-name="rowClass"
         :headers="headers"
         :hide-footer="true"
         :items="jobsStore.jobs"
         border-cell
      >
         <template #item-actions="item">
            <router-link :to="`/jobs/${item.id}`">View Log</router-link>
         </template>
      </EasyDataTable>
      <div class="pagination">
         <button @click="firstPage" :disabled="isFirstPage">First</button>
         <button @click="prevPage" :disabled="isFirstPage">Prev</button>
         <span class="info">{{pageRange}}</span>
         <button @click="nextPage" :disabled="isLastPage">Next</button>
         <button @click="lastPage" :disabled="isLastPage">Last</button>
      </div>
   </div>
</template>

<script setup>
import { onMounted, ref, computed} from 'vue'
import { useJobsStore } from '@/stores/jobs'

const jobsStore = useJobsStore()

const headers = [
  { text: "Job Type", value: "name" },
  { text: "Associated Object", value: "associatedObject" },
  { text: "Status", value: "status" },
  { text: "Warnings", value: "warnings" },
  { text: "Started", value: "startedAt" },
  { text: "Finished", value: "finishedAt" },
  { text: "", value: "actions"},
]

const jobsTable = ref()

const isFirstPage = computed(()=>{
   return jobsStore.searchOpts.page == 1
})
const isLastPage = computed(()=>{
    return jobsStore.searchOpts.page == jobsStore.totalJobs
})
const pageRange = computed(()=>{
   var firstPage = ((jobsStore.searchOpts.page-1)*jobsStore.searchOpts.rowsPerPage)+1
   var lastPage = firstPage + jobsStore.searchOpts.rowsPerPage -1
   return `${firstPage}-${lastPage} of ${jobsStore.totalJobs}`
})

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

function nextPage() {
   jobsTable.value.updatePage(jobsStore.searchOpts.page+1)
   jobsStore.getJobs()
}
function prevPage() {
   jobsTable.value.updatePage(jobsStore.searchOpts.page-1)
   jobsStore.getJobs()
}
function firstPage() {
   jobsTable.value.updatePage(1)
   jobsStore.getJobs()
}
function lastPage() {
   let lp = Math.floor(jobsStore.totalJobs / jobsStore.searchOpts.rowsPerPage)
   jobsTable.value.updatePage(lp)
   jobsStore.getJobs()
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
   }
   .pagination {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      align-items: center;
      padding: 10px 0;
      .info {
         margin: 0 5px 0 10px;
         font-size: 0.8em;
         display: inline-block;
      }
      button {
         margin-left: 5px;
      }
   }
   :deep(.jobs-table) {
      --easy-table-header-background-color: var(--uvalib-grey-lightest);
      --easy-table-body-row-hover-background-color: #FaFaFa;
   }
   :deep(.even-row) {
      --easy-table-body-row-background-color: #fafafa;
   }
   :deep(.running-row)  {
      --easy-table-body-row-background-color: var(--uvalib-blue-alt-light);
      --easy-table-body-row-hover-background-color: #def;
   }
   :deep(.warn-row)  {
      --easy-table-body-row-background-color: var(--uvalib-yellow-light);
      --easy-table-body-row-hover-background-color: var(--uvalib-yellow-light);
   }
   :deep(.error-row)  {
      --easy-table-body-row-background-color: #944;
      --easy-table-body-row-font-color: #fff;
      --easy-table-body-row-hover-background-color: #a66;
      --easy-table-body-row-hover-font-color: #fff;
      a {
         color: white !important;
      }
   }
</style>