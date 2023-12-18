<template>
   <h2>ArchivesSpace Reviews</h2>
   <div class="reviews">
      <DataTable :value="archivesSpace.reviews" ref="asReviewsTable" dataKey="id"
            stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm" :lazy="false"
            :resizableColumns="true" paginatorPosition="top"
            :paginator="true" :rows="15" :rowsPerPageOptions="[15,30,50]" removableSort
            :sortField="archivesSpace.searchOpts.sortField" :sortOrder="sortOrder"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
            v-model:filters="filter" filterDisplay="menu"
            :globalFilterFields="['metadata.title']"
         >
         <template #paginatorstart >
         </template>
         <template #paginatorend>
            <span class="js-search p-input-icon-right">
               <i class="pi pi-search" />
               <InputText v-model="filter['global'].value" placeholder="Search reviews..."/>
            </span>
            <DPGButton label="Clear" class="p-button-secondary clear" @click="clearSearch()"/>
         </template>
         <Column field="metadata.pid" header="PID"  class="nowrap" :sortable="true">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.metadata.id}`">
                  {{ slotProps.data.metadata.pid }}
               </router-link>
            </template>
         </Column>
         <Column field="metadata.title" header="Title" class="long-text" />
         <Column field="submittedAt" header="Requested" :sortable="true">
            <template #body="slotProps">
               {{formatDate(slotProps.data.submittedAt)}}
            </template>
         </Column>
         <Column field="submitter" header="Requested By" filterField="submitter.lastName" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Last name"/>
            </template>
            <template #body="slotProps">
               {{slotProps.data.submitter.lastName}}, {{slotProps.data.submitter.firstName}}
            </template>
         </Column>
         <Column field="status" header="Status" filterField="status" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <Dropdown v-model="filterModel.value" :options="statuses" optionLabel="name" optionValue="id" placeholder="Select status" />
            </template>
         </Column>
         <Column field="reviewer" header="Reviewer" filterField="reviewer.lastName" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Last name"/>
            </template>
            <template #body="slotProps">
               <span v-if="slotProps.data.reviewer">
                  {{slotProps.data.reviewer.lastName}}, {{slotProps.data.reviewer.firstName}}
               </span>
               <span v-else class="empty">N/A</span>
            </template>
         </Column>
         <Column field="reviewStartedAt" header="Review" :sortable="true">
            <template #body="slotProps">
               <span v-if="slotProps.data.reviewStartedAt">{{formatDate(slotProps.data.reviewStartedAt)}}</span>
               <span v-else class="empty">N/A</span>
            </template>
         </Column>
         <Column  header="Acts" class="acts">
            <template #body="slotProps">
               <DPGButton label="View images" class="p-button-secondary first" @click="viewClicked(slotProps.data)"/>
               <DPGButton label="Claim for review" class="p-button-secondary first" @click="reviewClicked(slotProps.data)"
                  :disabled="slotProps.data.status != 'requested' || user.ID == slotProps.data.submitter.id"/>
               <DPGButton label="Resubmit" class="p-button-secondary first" @click="resubmitClicked(slotProps.data)"
                  :disabled="!canResubmit(slotProps.data)" />
               <DPGButton label="Publish" class="p-button-secondary first" @click="publishClicked(slotProps.data)"
                  :disabled="slotProps.data.status != 'review' || user.ID != slotProps.data.reviewer.id" />
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
import Dropdown from 'primevue/dropdown'
import { FilterMatchMode } from 'primevue/api'
import { useConfirm } from "primevue/useconfirm"
import { usePinnable } from '@/composables/pin'
import { useArchivesSpaceStore } from '@/stores/archivesspace'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'

usePinnable("p-paginator-top")

const archivesSpace = useArchivesSpaceStore()
const user = useUserStore()
const confirm = useConfirm()

const filter = ref( {
      'global': {value: null, matchMode: FilterMatchMode.CONTAINS},
      'status': {value: null, matchMode: FilterMatchMode.EQUALS},
      'submitter.lastName': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
      'reviewer.lastName': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
})

const statuses = ref( [
   {id: "requested", name: "Requested"},
   {id: "rejected", name:  "Rejected"},
])

const sortOrder = computed(() => {
   if (archivesSpace.searchOpts.sortOrder == "desc") {
      return -1
   }
   return 1
})

const canResubmit = ( (data) => {
   if ( data.status != 'rejected' ) return false
   if ( !data.reviewer ) return false
   if ( data.reviewer.id == user.ID) return false
   return true
})

onMounted(() => {
   archivesSpace.getReviews()
})

const reviewClicked = ( (item) => {
   confirm.require({
      message: 'Are you sure you want claim this item for review?',
      header: 'Confirm Rreview',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await archivesSpace.claimForReview( item, user.ID )
      }
   })
})

const viewClicked = ( (item) => {
   window.open(`${archivesSpace.viewerBaseURL}/${item.metadata.pid}`, '_blank').focus()
})
const formatDate = (  (date ) => {
   if (date) {
      return dayjs(date).format("YYYY-MM-DD")
   }
   return ""
})

const clearSearch = (() => {
   filter.value['global'].value = ""
})

</script>

<style scoped lang="scss">
.reviews {
   min-height: 600px;
   text-align: left;
   padding: 0 25px;
   button.clear {
      margin-left: 10px;
   }
}
:deep(td.long-text) {
   white-space: break-spaces;
   max-width: 25%;
}
td.acts {
   vertical-align: top;

   button.p-button.first {
      margin: 0;
   }
   button.p-button.p-button-secondary {
      font-size: 0.75em;
      padding: 3px 6px;
      display: block;
      width: 100%;
      margin-top: 5px;
   }
}
</style>