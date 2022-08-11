<template>
   <h2>Job Statuses</h2>
   <div class="job-status">
      <EasyDataTable
         v-model:server-options="serverOptions"
         :server-items-length="jobsStore.totalJobs"
         table-class-name="jobs-table"
         :body-row-class-name="rowClass"
         :headers="headers"
         :items="jobsStore.jobs"
         :hide-rows-per-page="true"
         border-cell
      >
         <template #item-actions="item">
            <router-link :to="`/jobs/${item.id}`">View Log</router-link>
         </template>
      </EasyDataTable>
   </div>
</template>

<script setup>
import { onMounted} from 'vue'
import { useJobsStore } from '@/stores/jobs'

const jobsStore = useJobsStore()
const serverOptions = {
   page: 1,
   rowsPerPage: 30,
   sortBy: 'startedAt',
   sortType: 'desc',
}
const headers = [
  { text: "Job Type", value: "name" },
  { text: "Associated Object", value: "associatedObject" },
  { text: "Status", value: "status" },
  { text: "Warnings", value: "warnings" },
  { text: "Started", value: "startedAt" },
  { text: "Finished", value: "finishedAt" },
  { text: "", value: "actions"},
]

function rowClass(rowData) {
   if (rowData.status ==  "failure"){
      return "error-row"
   }
   return ""
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
   :deep(.jobs-table) {
      --easy-table-header-background-color: var(--uvalib-grey-lightest);
      --easy-table-body-row-hover-background-color: #FaFaFa;
   }
   :deep(.even-row) {
      --easy-table-body-row-background-color: #fafafa;
   }
   :deep(.error-row)  {
      --easy-table-body-row-background-color: #A77;
      --easy-table-body-row-font-color: #fff;
      --easy-table-body-row-hover-background-color: #B77;
      --easy-table-body-row-hover-font-color: #fff;
      a {
         color: white !important;
      }
   }
   :deep(.even-row.error-row)  {
      --easy-table-body-row-background-color: #944;
      --easy-table-body-row-font-color: #fff;
      --easy-table-body-row-hover-background-color: #A44;
      --easy-table-body-row-hover-font-color: #fff;
      a {
         color: white !important;
      }
   }
</style>