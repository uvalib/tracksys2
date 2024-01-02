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
               <ul class="acts">
                  <li><DPGButton label="View Images" icon="pi pi-external-link" iconPos="right"  class="p-button-secondary first" @click="viewClicked(slotProps.data)"/></li>
                  <li v-if="canReview(slotProps.data)"><DPGButton label="Claim for Review" class="p-button-secondary first" @click="reviewClicked(slotProps.data)"/></li>
                  <li v-if="canResubmit(slotProps.data)"><DPGButton label="Resubmit" class="p-button-secondary" @click="resubmitClicked(slotProps.data)"/></li>
                  <li v-if="canPublish(slotProps.data)"><DPGButton label="Reject" class="p-button-secondary" @click="rejectClicked(slotProps.data)"/></li>
                  <li v-if="canPublish(slotProps.data)"><DPGButton label="Publish Now" class="p-button-secondary" @click="publishClicked(slotProps.data.metadata)"/></li>
               </ul>
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
import { useRoute } from 'vue-router'

usePinnable("p-paginator-top")

const archivesSpace = useArchivesSpaceStore()
const user = useUserStore()
const confirm = useConfirm()
const route = useRoute()

const filter = ref( {
      'global': {value: null, matchMode: FilterMatchMode.CONTAINS},
      'status': {value: null, matchMode: FilterMatchMode.EQUALS},
      'submitter.lastName': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
      'reviewer.lastName': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
})

const statuses = ref( [
   {id: "requested", name: "Requested"},
   {id: "review", name:  "In Review"},
   {id: "rejected", name:  "Rejected"},
])

const sortOrder = computed(() => {
   if (archivesSpace.searchOpts.sortOrder == "desc") {
      return -1
   }
   return 1
})

const canReview = ( (data) => {
   return (data.status == 'requested' || data.status == 'review')  && user.ID != data.submitter.id
})

const canPublish = ((data) => {
   return data.status == 'review' && user.ID == data.reviewer.id
})
const canResubmit = ( (data) => {
   if ( data.status != 'rejected' ) return false
   if ( !data.reviewer ) return false
   if ( data.reviewer.id == user.ID) return false
   return true
})

onMounted(() => {
   if ( route.query.view == "reject" ) {
      filter.value.status.value = "rejected"
   } else if ( route.query.view == "review" ) {
      filter.value.status.value = "review"
   } else if ( route.query.view == "request" ) {
      filter.value.status.value = "requested"
   }
   archivesSpace.getReviews()
})

const reviewClicked = ( (item) => {
   confirm.require({
      message: 'Are you sure you want claim this item for review?',
      header: 'Confirm Review',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await archivesSpace.claimForReview( item, user.ID )
      }
   })
})

const resubmitClicked = ( (item) => {
   confirm.require({
      message: 'Are you sure you want resubmit this item for review?',
      header: 'Confirm Resubmit',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await archivesSpace.resubmit( item )
      }
   })
})

const rejectClicked = ( (item) => {

})

const publishClicked = ( (item) => {
   console.log(item)
   confirm.require({
      message: 'Are you sure you want publish this item to ArchivesSpace? After publication, the images will be visble to all ArchivesSpace users within a few minutes.',
      header: 'Confirm Publish',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await archivesSpace.publish( user.ID, item )
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
:deep(td.acts) {
   width: 130px;

   ul.acts {
      list-style: none;
      margin: 0;
      padding: 0;
      li {
         width: max-content;
         padding: 0;
         button.p-button.p-button-secondary {
            font-size: 0.75em;
            width: 130px;
            margin: 5px 0 0 0;
            padding: 0.4em 1em;
            .p-button-icon {
               color: #bbb;
            }
         }
      }
   }
}
</style>