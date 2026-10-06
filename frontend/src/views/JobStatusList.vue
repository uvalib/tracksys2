<template>
   <h2>Job Statuses</h2>

   <div class="row-between p-2 sticky z-50 bg-white border-b-1 border-brand-grey-100" :style="{top: headerHeight}">
      <UButton label="Delete selected" :disabled="selectedJobs.length == 0"  color="secondary" @click="deletAllClicked"/>
      <div class="row-left">
         <UPagination color="neutral" variant="ghost"
            v-model:page="jobsStore.searchOpts.currPage" :items-per-page="jobsStore.searchOpts.limit" 
            :total="jobsStore.totalJobs" @update:page="pageChanged"
         />
         <USelect v-model="jobsStore.searchOpts.limit" :items="[15,30,100]" @change="perPageChanged" />
      </div>
      <UInput v-model="jobsStore.searchOpts.query" placeholder="Search Job Status" @update:modelValue="queryJobs"/>
   </div>

   <UTable :data="jobsStore.jobs" :columns="columns" :meta="meta" v-if="jobsStore.jobs.length > 0">
      <template #select-header="">
         <UCheckbox size="xl" :modelValue="allCheckboxValue" @update:modelValue="toggleAll"/>
      </template>
      <template #select-cell="{ row }">
         <UCheckbox size="xl" :modelValue="selectedJobs.includes(row.original.id)" @update:modelValue="toggleJobSelected(row.original)"/>
      </template>
      <template #associatedObject-cell="{ row }">
         <template v-if="getAssociatedObjectLink(row.original.associatedObject)">
            <router-link :to="getAssociatedObjectLink(row.original.associatedObject)">{{row.original.associatedObject}}</router-link>
         </template>
         <template v-else>
            {{ row.original.associatedObject }}
         </template>
      </template>
      <template #startedAt-cell="{ row }">
         {{ $formatDateTime(row.original.startedAt) }}
      </template>
      <template #finishedAt-cell="{ row }">
         {{ $formatDateTime(row.original.finishedAt) }}
      </template>
      <template #actions-cell="{ row }">
         <div class="row-right gap-2!">
            <UButton :to="`/jobs/${row.original.id}`"  icon="i-lucide-eye"  label="View" class="text-white!" size="xs"/>
            <UButton label="Delete" color="error" icon="i-lucide-trash" size="xs" @click="deleteJob(row.original.id)"/>
         </div>
      </template>
   </UTable>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useJobsStore } from '@/stores/jobs'
import { useConfirm } from "@/composables/useConfirm"

const jobsStore = useJobsStore()

const selectedJobs = ref([])


const columns = [
    {
      accessorKey: 'select'
   },
   {
      accessorKey: 'name',
      header: 'Job Type'
   },
   {
      accessorKey: 'associatedObject',
      header: 'Associated Object'
   },
   {
      accessorKey: 'status',
      header: 'Status'
   },
   {
      accessorKey: 'warnings',
      header: 'Warnings'
   },
   {
      accessorKey: 'startedAt',
      header: 'Started'
   },
   {
      accessorKey: 'finishedAt',
      header: 'Finished'
   },
   {
      accessorKey: 'actions',
      header: 'Actions'
   }
]

const meta = {
  class: {
    tr: (row) => {
      if (row.original.status === 'failure') {
        return 'bg-red-600/25'
      } 
      if (row.original.status === 'running') {
        return 'bg-sky-600/20'
      }
      return ''
    }
  }
}

const allCheckboxValue = computed( () => {
   if (allJobsSelected.value) return true
   if ( someJobsSelected.value ) return "indeterminate"
   return false
})
const allJobsSelected = computed(() =>{
   return selectedJobs.value.length == jobsStore.searchOpts.limit
})
const someJobsSelected = computed(() =>{
   return selectedJobs.value.length > 0
})
const toggleAll = (() => {
   if ( selectedJobs.value.length < jobsStore.searchOpts.limit ) {
      selectedJobs.value = jobsStore.jobs.map( js => js.id)
   } else {
      selectedJobs.value = []   
   }
})
const toggleJobSelected = ((js) => {
   const idx = selectedJobs.value.indexOf(j => j.id == js.id) 
   if ( idx > -1 ) {
      selectedJobs.value.slice(idx,1)
   } else {
      selectedJobs.value.push(js.id)
   }

})

const headerHeight = computed(() => {
   let hdr = document.getElementById('uva-header')
   return `${hdr.clientHeight}px`
})

const queryJobs = (() => {
   jobsStore.getJobs(false)
})

const deletAllClicked = (async () => {
   const msg = 'Are you sure you want delete the selected job status records? All data will be lost. This cannot be reversed.'
   const resp = await useConfirm("Confirm Delete All", msg, "Delete")
   if (resp) {
      jobsStore.deleteJobs( selectedJobs.value )
   } 
})

const deleteJob = (async (id) => {
   const msg = `Are you sure you want delete this job status?`
   const resp = await useConfirm("Confirm Delete", msg, "Delete")
   if (resp) {
      jobsStore.deleteJobs( [id] )
   } 
})

const getAssociatedObjectLink = (( objName ) => {
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
   if (objType == "masterfile") {
      return `/masterfiles/${objID}`
   }
   return ""
})

const perPageChanged = (() => {
   jobsStore.searchOpts.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   jobsStore.searchOpts.start  = (jobsStore.searchOpts.currPage-1) * jobsStore.searchOpts.limit
   jobsStore.getJobs()
})

const onRowSelect = (() => {
   selectAll.value = selectedJobs.value === jobsStore.searchOpts.limit
})

const onRowUnselect = (() => {
   selectAll.value  = false
})

const onSelectAllChange = ((event) => {
   selectAll.value = event.checked
   if (selectAll.value) {
      selectedJobs.value = jobsStore.jobs
   }
   else {
      selectedJobs.value = []
   }
})

onMounted(() => {
   jobsStore.getJobs()
   document.title = `Job Statuses`
})
</script>

<style scoped lang="scss">
.none {
   font-style: italic;
   color: var(--uvalib-grey-light);
}
.job-search {
   display: flex;
   flex-flow: row nowrap;
}
.job-status {
   min-height: 600px;
   text-align: left;
   padding: 0 25px;

   .js-search {
      margin-right: 10px;
   }
   .sep {
      display: inline-block;
      margin: 0 10px;
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
   // :deep(.error-row) {
   //    background-color: #944 !important;
   //    color: #fff;
   //    a {
   //       color: #fff !important;
   //    }
   //    .row-acts {
   //       a, button.p-button-text {
   //          color: #fff !important;
   //       }
   //    }
   //    &:hover {
   //       background-color: #a44 !important;
   //    }
   // }
   // :deep(tr.p-highlight.error-row) {
   //    color: white !important;
   // }
   // :deep(.running-row)  {
   //    background-color: var(--uvalib-blue-alt-light) !important;
   //    &:hover {
   //       background-color: #def !important;
   //    }
   // }
   // :deep(.warn-row)  {
   //    background-color: var(--uvalib-yellow-light) !important;
   // }
}
</style>