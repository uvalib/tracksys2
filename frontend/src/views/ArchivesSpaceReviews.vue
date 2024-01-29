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
               {{$formatDate(slotProps.data.submittedAt)}}
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
            <template #body="slotProps">
               <span class="as-status">{{ `${slotProps.data.status[0].toUpperCase()}${slotProps.data.status.slice(1)}` }}</span>
            </template>
         </Column>
         <Column field="notes" header="Notes" >
            <template #body="slotProps">
               <DPGButton v-if="slotProps.data.notes"  class="notes" label="Click to View" severity="secondary" @click="notesClicked(slotProps.data)"/>
               <DPGButton v-else  label="Click to Add" class="notes" severity="secondary" @click="notesClicked(slotProps.data)"/>
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
               <span v-if="slotProps.data.reviewStartedAt">{{$formatDate(slotProps.data.reviewStartedAt)}}</span>
               <span v-else class="empty">N/A</span>
            </template>
         </Column>
         <Column  header="Acts" class="acts">
            <template #body="slotProps">
               <ul class="acts">
                  <li><DPGButton label="View Images" icon="pi pi-external-link" iconPos="right" severity="secondary" class="first" @click="viewClicked(slotProps.data)"/></li>
                  <li v-if="canReview(slotProps.data)"><DPGButton label="Claim for Review" severity="secondary" @click="reviewClicked(slotProps.data)"/></li>
                  <li v-if="canResubmit(slotProps.data)"><DPGButton label="Resubmit" severity="secondary" @click="resubmitClicked(slotProps.data)"/></li>
                  <li v-if="canCancel(slotProps.data)"><DPGButton label="Cancel Submission" severity="danger" @click="cancelClicked(slotProps.data)"/></li>
                  <li v-if="canPublish(slotProps.data)"><DPGButton label="Reject" severity="danger" @click="rejectClicked(slotProps.data)"/></li>
                  <li v-if="canPublish(slotProps.data)"><DPGButton label="Publish Now" severity="primary" @click="publishClicked(slotProps.data.metadata)"/></li>
               </ul>
            </template>
         </Column>
      </DataTable>
   </div>
   <Dialog v-model:visible="rejectRequested" :modal="true" header="Reject Submission">
      <div>Reject submission {{ tgtReview.metadata.pid }} - {{ tgtReview.metadata.title }}</div>
      <label class="reject-note">Please add some notes about the rejection:</label>
      <textarea v-model="reason" autofocus rows="5" ref="reasontxt" :class="{'invalid': rejectError}"></textarea>
      <template #footer>
         <DPGButton label="Cancel" severity="secondary" @click="rejectCanceled()"/>
         <DPGButton label="Reject" class="reject" severity="danger" @click="rejectSubmitted()"/>
      </template>
   </Dialog>
   <Dialog v-model:visible="showNotes" :modal="true" header="Submission Notes">
      <div>{{ tgtReview.metadata.pid }} - {{ tgtReview.metadata.title }}</div>
      <textarea class="notes" v-if="editNote" v-model="newNotes" ref="noteedit" rows="10"></textarea>
      <div v-else class="note-text">{{ newNotes }}</div>
      <template #footer>
         <template v-if="editNote">
            <DPGButton label="Cancel" severity="secondary" @click="cancelNoteEdit()"/>
            <DPGButton label="Submit" severity="primary" class="left-margin" @click="submitNoteEdit()"/>
         </template>
         <template v-else>
            <DPGButton label="Edit" severity="primary" @click="editNoteClicked()"/>
            <DPGButton label="Close" severity="secondary" class="left-margin" @click="closeNotesClicked()"/>
         </template>
      </template>
   </Dialog>
</template>

<script setup>
import { onMounted, ref, computed, nextTick } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Dialog from 'primevue/dialog'
import { FilterMatchMode } from 'primevue/api'
import { useConfirm } from "primevue/useconfirm"
import { usePinnable } from '@/composables/pin'
import { useArchivesSpaceStore } from '@/stores/archivesspace'
import { useUserStore } from '@/stores/user'
import { useRoute } from 'vue-router'

usePinnable("p-paginator-top")

const archivesSpace = useArchivesSpaceStore()
const user = useUserStore()
const confirm = useConfirm()
const route = useRoute()

const tgtReview = ref()

const rejectRequested = ref(false)
const reason = ref("")
const reasontxt = ref()
const rejectError = ref(false)

const showNotes = ref(false)
const newNotes = ref("")
const editNote = ref(false)
const noteedit = ref()

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
      filter.value.status.value = "rejected"
   } else if ( route.query.view == "review" ) {
      filter.value.status.value = "review"
   } else if ( route.query.view == "request" ) {
      filter.value.status.value = "requested"
   }
   archivesSpace.getReviews()
})

const notesClicked= ((item) => {
   tgtReview.value = item
   showNotes.value = true
   newNotes.value = item.notes
   editNote.value = item.notes == ""
})

const editNoteClicked = ( () => {
   editNote.value = true
   nextTick( () => {
      noteedit.value.focus()
   })
})

const closeNotesClicked = ( () => {
   showNotes.value = false
   tgtReview.value = null
})

const cancelNoteEdit = ( () => {
   editNote.value = false
   tgtReview.value = null
})

const submitNoteEdit = ( async () => {
   await archivesSpace.updateNotes( tgtReview.value, newNotes.value )
   editNote.value = false
   tgtReview.value = null
   showNotes.value = false
})

const reviewClicked = ( (item) => {
   confirm.require({
      message: 'Are you sure you want claim this item for review?',
      header: 'Confirm Review',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await archivesSpace.claimForReview( item.metadata, user.ID )
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
         await archivesSpace.resubmit( item.metadata )
      }
   })
})

const cancelClicked = ( (item) => {
   confirm.require({
      message: 'Are you sure you want cancel this submission? All review data will be lost.',
      header: 'Confirm Cancel',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await archivesSpace.cancel( item )
      }
   })
})

const rejectClicked = ( (item) => {
   rejectRequested.value = true
   tgtReview.value = item
   reason.value = ""
   rejectError.value = false
})

const rejectCanceled = ( () => {
   rejectRequested.value = false
   tgtReview.value = null
})

const rejectSubmitted = ( async () => {
   if ( reason.value == "") {
      reasontxt.value.focus()
      rejectError.value = true
   } else {
      await archivesSpace.reject( user.ID, tgtReview.value.metadata, reason.value )
      rejectRequested.value = false
      tgtReview.value = null
   }
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
   button.notes {
      font-size: 0.75em;
      padding: 5px 10px;
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
         button.p-button {
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
textarea {
   width: 100%;
   border-color: var(--uvalib-grey-light);
   border-radius: 5px;
   font-family: "franklin-gothic-urw", arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   color: var(--color-primary-text);
   padding: 5px 10px;
   &:focus {
      @include be-accessible();
   }
}
div.note-text {
   height: 250px;
   overflow-y: scroll;
   margin-top: 15px;
   padding: 5px 10px;
   border: 1px solid var(--uvalib-grey-light);
   border-radius: 5px;
   background-color: white;
}
textarea.notes {
   margin-top: 15px;
}
textarea.invalid {
   border-color: var(--uvalib-red);
   border-width: 2px;
}
label.reject-note {
   margin: 15px 0 10px 0;
   display: block;
}
button.reject, button.left-margin {
   margin-left: 10px;
}
</style>