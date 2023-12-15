<template>
   <h2>ArchivesSpace Reviews</h2>
   <div class="reviews">
      <DataTable :value="archivesSpace.reviews" ref="apTrustStatusTable" dataKey="id"
            stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm" :lazy="false"
            :resizableColumns="true" columnResizeMode="fit" paginatorPosition="top"
            :paginator="true" :rows="15" :rowsPerPageOptions="[15,30,50]" removableSort
            :sortField="archivesSpace.searchOpts.sortField" :sortOrder="sortOrder"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
         >
         <Column field="metadata.pid" header="PID"  class="nowrap" :sortable="true">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.metadata.id}`">
                  {{ slotProps.data.metadata.pid }}
               </router-link>
            </template>
         </Column>
         <Column field="metadata.title" header="Title" />
         <Column field="submittedAt" header="Requested" :sortable="true">
            <template #body="slotProps">
               {{formatDate(slotProps.data.submittedAt)}}
            </template>
         </Column>
         <Column field="metadata.submitter" header="Requested By">
            <template #body="slotProps">
               {{slotProps.data.submitter.lastName}}, {{slotProps.data.submitter.firstName}}
            </template>
         </Column>
         <Column field="status" header="Status" />
         <Column field="metadata.reviewer" header="Reviewer">
            <template #body="slotProps">
               <span v-if="slotProps.data.reviewer">
                  {{slotProps.data.reviewer.lastName}}, {{slotProps.data.reviewer.firstName}}
               </span>
               <span v-else class="empty">N/A</span>
            </template>
         </Column>
         <Column  header="Acts" />
      </DataTable>
   </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { usePinnable } from '@/composables/pin'
import { useArchivesSpaceStore } from '@/stores/archivesspace'
import dayjs from 'dayjs'

usePinnable("p-paginator-top")

const archivesSpace = useArchivesSpaceStore()

onMounted(() => {
   archivesSpace.getReviews()
})

const sortOrder = computed(() => {
   if (archivesSpace.searchOpts.sortOrder == "desc") {
      return -1
   }
   return 1
})

const formatDate = (  (date ) => {
   if (date) {
      return dayjs(date).format("YYYY-MM-DD")
   }
   return ""
})

</script>

<style scoped lang="scss">
.reviews {
   min-height: 600px;
   text-align: left;
   padding: 0 25px;
}
</style>