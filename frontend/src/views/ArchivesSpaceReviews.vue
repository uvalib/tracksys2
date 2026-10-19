<template>
   <h2>ArchivesSpace Reviews</h2>
   <div class="row-between bg-white border-b border-t border-brand-grey-100 py-3 px-2 items-center">
      <span>{{ table?.tableApi?.getFilteredRowModel().rows.length }} Reviews<span v-if="filters.length > 0"> (filtered)</span></span>    
   </div>
   <UTable ref="table" :data="archivesSpace.reviews" :columns="columns" v-model:column-filters="filters" v-model:sorting="sorting" sticky virtualize class="h-150">
      <template #pid-cell="{ row }">
         <router-link :to="`/metadata/${row.original.metadata.id}`">
            {{ row.original.metadata.pid }}
         </router-link>
      </template>
      <template #submittedAt-cell="{ row }">
         <span>{{ toCalendarDate(parseAbsolute(row.original.submittedAt)).toString() }}</span>
      </template>
      <template #submitter-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('submitter')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Submitter..." />
         </FilterPopover>
      </template>
      <template #status-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('status')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <USelect v-model="pendingFilter" :items="statuses" placeholder="Status..." />
         </FilterPopover>
      </template>
      <template #notes-cell="{ row }">
         <UButton v-if="row.original.notes" class="notes" size="xs" label="View" color="secondary" @click="notesClicked(row.original)"/>
         <UButton v-else label="Add" class="notes" size="xs" color="secondary" @click="notesClicked(row.original)"/>   
      </template>
      <template #reviewer-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied('reviewer')" @opened="filterOpened(column)"
            @clear="clearFilter(column)" @apply="applyFilter(column)">
            <UInput v-model="pendingFilter" placeholder="Reviewer..." />
         </FilterPopover>
      </template>
      <template #reviewStartedAt-cell="{ row }">
         <span v-if="row.original.reviewStartedAt">
            {{ toCalendarDate(parseAbsolute(row.original.reviewStartedAt)).toString() }}
         </span>
      </template>
      <template #acts-cell="{ row }">
         <UDropdownMenu :items="getActions(row.original)">
            <UButton icon="i-lucide-ellipsis-vertical"
               color="secondary" aria-label="Actions" 
            />
         </UDropdownMenu>
      </template>
   </UTable>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useArchivesSpaceStore } from '@/stores/archivesspace'
import NotesModal from '@/components/archivesspace/NotesModal.vue'
import RejectModal from '@/components/archivesspace/RejectModal.vue'
import { useUserStore } from '@/stores/user'
import { useRoute } from 'vue-router'
import { toCalendarDate, parseAbsolute } from '@internationalized/date'
import { useConfirm } from "@/composables/useConfirm"

const archivesSpace = useArchivesSpaceStore()
const user = useUserStore()
const route = useRoute()
const table = ref()
const notesModal = useOverlay().create(NotesModal)
const rejectModal = useOverlay().create(RejectModal)

const pendingFilter = ref()
const filters = ref([])
const sorting = ref([])
const columns = [
   { accessorKey: 'metadata.pid', id: 'pid', header: 'PID'},
   { accessorKey: 'metadata.title', id: 'title', header: 'Title'},
   { accessorKey: 'submittedAt', header: 'Requested'},
   { 
      id: 'submitter', 
      header: 'Requested By',
      accessorFn: (row) => `${row.submitter.lastName}, ${row.submitter.firstName}`
   },
   { accessorKey: 'status', header: 'Status'},
   { accessorKey: 'notes', header: 'Notes'},
   { 
      id: 'reviewer', 
      header: 'Reviewer',
      accessorFn: (row) => row.reviewer ? `${row.reviewer.lastName}, ${row.reviewer.firstName}` : ''
   },
   { accessorKey: 'reviewStartedAt', header: 'Review'},
   { accessorKey: 'acts', header: 'Actions'},
]

const isFilterApplied = ((name) => {
   return filters.value.findIndex( f => f.id == name) > -1
})
const filterOpened = ((column) => {
   pendingFilter.value = column.getFilterValue()  
})
const clearFilter = ((column) => {
   column.setFilterValue(null)
})
const applyFilter = ((column) => {
   column.setFilterValue(pendingFilter.value)
})
const statuses = ref( [
   {value: "requested", label: "Requested"},
   {value: "review", label:  "In Review"},
   {value: "rejected", label:  "Rejected"},
])

const getActions = ((item) => {
   let acts = [ [ { label: 'View', onSelect: () => viewClicked(item) } ] ]
   if ( canReview(item) ) {
      acts.push( [ {label: "Claim for review", onSelect: () => reviewClicked(item) } ])
   }
   if ( canResubmit(item) ) {
      acts.push( [ {label: "Resubmit", onSelect: () => resubmitClicked(item) } ])
   }
   if ( canCancel(item) ) {
      acts.push( [ {label: "Cancel Submission", onSelect: () => cancelClicked(item) } ])
   }
   if ( canPublish(item) ) {
      acts.push( [ {label: "Reject", onSelect: () => rejectClicked(item) } ])
      acts.push( [ {label: "Publish", onSelect: () => publishClicked(item) } ])
   }
   return acts  
})

const canReview = ( (data) => {
   if (data.status == 'requested' ) {
      return user.ID != data.submitter.id
   }
   if (data.status == 'review' ) {
      return user.ID != data.submitter.id && user.ID != data.reviewer.id
   }
   return false
})

const canPublish = ((data) => {
   return data.status == 'review' && user.ID == data.reviewer.id
})
const canCancel = ((data) => {
   return data.status != 'review'
})
const canResubmit = ( (data) => {
   // must be in a a rejected state and the current user is not the reviewer
   if ( data.status != 'rejected' ) return false
   if ( !data.reviewer ) return false
   if ( data.reviewer.id == user.ID) return false
   return true
})

onMounted(() => {
   if ( route.query.view == "reject" ) {
      filters.value = [ {"id": "status", "value": "rejected" }]
   } else if ( route.query.view == "review" ) {
      filters.value = [ {"id": "status", "value": "review" }]
   } else if ( route.query.view == "request" ) {
      filters.value = [ {"id": "status", "value": "requested" }]
   }
   archivesSpace.getReviews()
   document.title = "ArchivesSpace Reviews"
})

const notesClicked= ((item) => {
   notesModal.open({item: item})
})

const reviewClicked = ( async (item) => {
   const msg = 'Are you sure you want claim this item for review?'
   const resp = await useConfirm("Confirm Review", msg, "Claim")
   if (resp) {
      archivesSpace.claimForReview( item.metadata, user.ID )
   } 
})

const resubmitClicked = ( async (item) => {
   const msg = 'Are you sure you want resubmit this item for review?'
   const resp = await useConfirm("Confirm Resubmit", msg, "Resubmit")
   if (resp) {
      archivesSpace.resubmit( item.metadata )
   } 
})

const cancelClicked = ( async (item) => {
   const msg = 'Are you sure you want cancel this submission? All review data will be lost.'
   const resp = await useConfirm("Confirm Cancel", msg, "Cancel Submission")
   if (resp) {
      archivesSpace.cancel( item )
   } 
})

const rejectClicked = ( async (item) => {
  rejectModal.open({userID: user.ID, item: item})
})

const publishClicked = ( (item) => {
   console.log(item)
   confirm.require({
      message: 'Are you sure you want publish this item to ArchivesSpace? After publication, the images will be visble to all ArchivesSpace users within a few minutes.',
      header: 'Confirm Publish',
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Publish'
      },
      accept: async () => {
         await archivesSpace.publish( user.ID, item )
      }
   })
})

const viewClicked = ( (item) => {
   window.open(`${archivesSpace.viewerBaseURL}/${item.metadata.pid}?unit=${item.unitID}`, '_blank').focus()
})
</script>

<style scoped lang="scss">
</style>