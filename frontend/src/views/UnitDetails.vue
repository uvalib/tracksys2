<template>
   <h2>{{pageTitle}}</h2>
   <div class="unit-acts">
      <DPGButton label="OCR"  @click="unitOCRClicked()" v-if="canOCR" />
      <DPGButton label="PDF" @click="unitPDFClicked()" v-if="detail.metadata && unitsStore.masterFiles.length > 0 && detail.reorder==false" />
      <DPGButton label="Edit" @click="editUnit()"/>
   </div>
   <div v-if="detail.lastError" class="last-error">
      <span>Recent Error:</span>
      <router-link :to="`/jobs/${detail.lastError.jobID}`">{{detail.lastError.error}}</router-link>
   </div>
   <div class="details" v-if="systemStore.working==false">
      <Panel header="General Information">
         <dl>
            <DataDisplay label="Status" :value="detail.status">
               <span :class="`status ${detail.status}`">{{displayStatus(detail.status)}}</span>
            </DataDisplay>
            <DataDisplay label="Patron Source URL" :value="detail.patronSourceURL">
               <a :href="detail.patronSourceURL" target="_blank">{{detail.patronSourceURL}}<i class="icon fas fa-external-link"></i></a>
            </DataDisplay>
            <DataDisplay label="Special Instructions" :value="detail.specialInstructions"/>
            <DataDisplay label="Staff Notes" :value="detail.staffNotes"/>
            <DataDisplay label="Complete Scan" :value="flagString(detail.completeScan)">
               <span :class="`flag ${flagString(detail.completeScan)}`">{{displayFlag(detail.completeScan)}}</span>
            </DataDisplay>
            <DataDisplay label="Throwaway" :value="flagString(detail.throwAway)">
               <span :class="`flag ${flagString(detail.throwAway)}`">{{displayFlag(detail.throwAway)}}</span>
            </DataDisplay>
         </dl>
      </Panel>
      <Panel header="Digitization Information">
         <dl>
            <template v-if="detail.intendedUse">
               <DataDisplay label="Intended Use" :value="detail.intendedUse.name"/>
               <DataDisplay label="Deliverable Format" :value="detail.intendedUse.deliverableFormat"/>
               <DataDisplay label="Deliverable Resolution" :value="detail.intendedUse.deliverableResolution"/>
            </template>
            <template v-else>
               <DataDisplay label="Intended Use" value=""/>
               <DataDisplay label="Deliverable Format" value=""/>
               <DataDisplay label="Deliverable Resolution" value=""/>
            </template>
            <DataDisplay label="OCR Master Files" :value="flagString(detail.ocrMasterFiles)">
               <span :class="`flag ${flagString(detail.ocrMasterFiles)}`">{{displayFlag(detail.ocrMasterFiles)}}</span>
            </DataDisplay>
            <DataDisplay label="Remove Watermark" :value="flagString(detail.removeWatermark)">
               <span :class="`flag ${flagString(detail.removeWatermark)}`">{{displayFlag(detail.removeWatermark)}}</span>
            </DataDisplay>
            <DataDisplay label="Date Archived" :value="formatDate(detail.dateArchived)" />
            <DataDisplay label="Date Patron Deliverables Ready" :value="formatDate(detail.datePatronDeliverablesReady)" />
            <DataDisplay label="In Digital Library" :value="flagString(detail.includeInDL)">
               <span :class="`flag ${flagString(detail.includeInDL)}`">{{displayFlag(detail.includeInDL)}}</span>
            </DataDisplay>
            <DataDisplay label="Date DL Deliverables Ready" :value="formatDate(detail.dateDLDeliverablesReady)" />
         </dl>

         <div class="acts-wrap" v-if="detail.status != 'finalizing' && (userStore.isAdmin || userStore.isSupervisor)">
            <div class="acts" v-if="!detail.projectID && unitsStore.masterFiles.length == 0">
               <CreateProjectDialog />
            </div>
            <div class="acts">
               <DPGButton v-if="detail.reorder && !detail.datePatronDeliverablesReady" @click="generateDeliverablesClicked"
                  class="p-button-secondary" label="Generate Deliverables" />
               <DPGButton v-if="detail.intendedUseID != 110 && detail.datePatronDeliverablesReady" @click="generateDeliverablesClicked"
                  class="p-button-secondary" label="Regenerate Deliverables" />
               <DPGButton v-if="unitsStore.masterFiles.length > 0 && detail.status != 'error'" @click="regenerateIIIFClicked"
                  class="p-button-secondary" label="Regenerate IIIF Manifest" />
               <template v-if="detail.status == 'done'">
                  <DPGButton v-if="detail.dateArchived" @click="downloadClicked" class="p-button-secondary" label="Download Unit From Archive" />
               </template>
               <template v-else>
                  <DPGButton v-if="detail.reorder && detail.datePatronDeliverablesReady" @click="completeClicked"
                     class="p-button-secondary" label="Complete Unit" />
               </template>
            </div>
         </div>

         <div class="note"
            v-if="(unitsStore.detail.status != 'approved' || unitsStore.detail.order.status != 'approved') && detail.projectID && unitsStore.masterFiles.length == 0"
         >
            Cannot create project, unit or order has not been approved.
         </div>
         <div class="note" v-if="detail.status == 'finalizing'">
            Unit finalization is in-progress. No other actions available at this time.
         </div>
         <div class="note" v-if="detail.status == 'error' && detail.reorder == false">
            Unit has failed finalization finalization. Correct the errors and use the 'Retry Finalization' button on the project page to restart finalization.
         </div>
      </Panel>
   </div>
   <div class="details" v-if="systemStore.working==false">
      <Panel header="Related Information" class="small">
         <dl>
            <DataDisplay v-if="detail.metadata" label="Metadata" :value="detail.metadata.pid">
               <router-link :to="`/metadata/${detail.metadata.id}`">{{detail.metadata.pid}}: {{detail.metadata.title}}</router-link>
            </DataDisplay>
            <DataDisplay v-else label="Metadata" value=""/>
            <DataDisplay label="Order" :value="`${detail.orderID}`">
               <router-link :to="`/orders/${detail.orderID}`">{{detail.orderID}}</router-link>
            </DataDisplay>
            <DataDisplay label="Order Units" :value="`${detail.relatedUnits}`">
               <div class="related-unit-ids">
                  <template v-for="(uid,idx) in detail.relatedUnits" :key="`related-${uid}`">
                     <template v-if="idx > 0"><span class="sep"></span></template>
                     <router-link :to="`/units/${uid}`" v-if="uid != detail.id">{{uid}}</router-link>
                     <span class="current-unit" v-else>{{uid}}</span>
                  </template>
               </div>
            </DataDisplay>
            <DataDisplay v-if="detail.projectID"  label="Project" :value="`${detail.projectID}`">
               <a :href="`${systemStore.projectsURL}/projects/${detail.projectID}`" target="_blank">
                  {{detail.projectID}}<i class="icon fas fa-external-link"></i>
               </a>
            </DataDisplay>
         </dl>
      </Panel>
      <Panel header="Attachments">
         <div v-if="detail.attachments && detail.attachments.length > 0">
            <DataTable :value="detail.attachments" ref="attachmentsTable" dataKey="id"
               stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm attachments-data"
            >
               <Column field="filename" header="File Name" />
               <Column field="description" header="Description" />
               <Column header="" class="row-acts">
                  <template #body="slotProps">
                     <DPGButton label="Download" class="p-button-secondary" @click="downloadAttachment(slotProps.data)"/>
                     <DPGButton label="Delete" class="p-button-secondary" @click="deleteAttachment(slotProps.data)"/>
                  </template>
               </Column>
            </DataTable>
         </div>
         <div v-else>
            No attachments are associated with this unit.
         </div>
         <p v-if='detail.status != "approved" &&  detail.status != "done"'>
            Attachments cannot be aded to unapproved units.
         </p>
         <div v-else class="toolbar">
            <AddAttachmentDialog  v-if="systemStore.working==false"/>
         </div>
      </Panel>
   </div>
   <MasterFilesList />
</template>

<script setup>
import { onBeforeMount, computed } from 'vue'
import { useRoute, onBeforeRouteUpdate, useRouter } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useUnitsStore } from '@/stores/units'
import { useUserStore } from '@/stores/user'
import Panel from 'primevue/panel'
import { storeToRefs } from "pinia"
import DataDisplay from '../components/DataDisplay.vue'
import dayjs from 'dayjs'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import CreateProjectDialog from '../components/unit/CreateProjectDialog.vue'
import AddAttachmentDialog from '../components/unit/AddAttachmentDialog.vue'
import MasterFilesList from '../components/unit/MasterFilesList.vue'
import { useConfirm } from "primevue/useconfirm"

const confirm = useConfirm()
const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()
const unitsStore = useUnitsStore()
const userStore = useUserStore()

const { detail } = storeToRefs(unitsStore)

const pageTitle = computed(() => {
   let t = 'Unit'
   if ( unitsStore.detail ) {
      t += ` ${unitsStore.detail.id}`
      if (unitsStore.detail.reorder) {
         t += " : REORDER"
      }
   }
   return t
})
const canOCR = computed(() => {
   if ( !unitsStore.detail.metadata ) return false
   if ( !unitsStore.detail.metadata.ocrHint) return false
   let isCandidate = unitsStore.detail.metadata.ocrHint.ocrCandidate
   return ( unitsStore.masterFiles.length > 0 && unitsStore.detail.reorder == false && isCandidate && (userStore.isAdmin || userStore.isSupervisor) )
})

onBeforeRouteUpdate(async (to) => {
   let uID = to.params.id
   unitsStore.getDetails( uID )
   unitsStore.getMasterFiles(uID)
})

onBeforeMount(() => {
   let uID = route.params.id
   unitsStore.getDetails( uID )
   unitsStore.getMasterFiles(uID)
   document.title = `Unit #${uID}`
})

function flagString( flag ) {
   if ( flag === true) {
      return "true"
   }
   return "false"
}

async function completeClicked() {
   let update = {
      status: "done",
      patronSourceURL: unitsStore.detail.patronSourceURL,
      specialInstructions: unitsStore.detail.specialInstructions,
      staffNotes: unitsStore.detail.staffNotes,
      completeScan: unitsStore.detail.completeScan,
      throwAway: unitsStore.detail.throwAway,
      orderID: unitsStore.detail.orderID,
      metadataID: unitsStore.detail.metadataID,
      intendedUseID: unitsStore.detail.intendedUse.id,
      ocrMasterFiles: unitsStore.detail.ocrMasterFiles,
      removeWatermark: unitsStore.detail.removeWatermark,
      includeInDL: unitsStore.detail.includeInDL,
   }
   await unitsStore.submitEdit( update )
   systemStore.toastMessage("Unit Updated", "Unit has been marked as done.")
}

function unitOCRClicked() {
   unitsStore.startUnitOCR()
}
function unitPDFClicked() {
   let url = `${systemStore.pdfURL}/${unitsStore.detail.metadata.pid}?unit=${unitsStore.detail.id}`
   window.open(url)
}
function regenerateIIIFClicked() {
   unitsStore.regenerateIIIF()
}
function downloadClicked() {
   unitsStore.downloadFromArchive( userStore.computeID )
}
function generateDeliverablesClicked() {
   unitsStore.generateDeliverables()
}

function editUnit() {
   router.push(`/units/${route.params.id}/edit`)
}

function downloadAttachment(item) {
   let url = `${systemStore.jobsURL}/units/${unitsStore.detail.id}/attachments/${item.filename}`
   window.open(url)
}

function deleteAttachment(item) {
   confirm.require({
      message: 'Are you sure you want delete the selected attachment? All data will be lost. This cannot be reversed.',
      header: 'Confirm Delete Attachment',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await unitsStore.deleteAttachment(item)
      }
   })
}

function displayStatus( id) {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
}

function displayFlag( f ) {
   if (f) return "Yes"
   return "No"
}

function formatDate( dateStr ) {
   if (dateStr) {
      let d = dayjs(dateStr)
      return d.format("YYYY-MM-DD HH:mm A")
   }
   return ""
}

</script>

<style scoped lang="scss">
div.unit-acts {
   position: absolute;
   right:15px;
   top: 15px;
   button.p-button {
      margin-right: 5px;
      font-size: 0.9em;
   }
}
.related-unit-ids {
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
.sep {
      margin-right: 5px;
      display: inline-block;
   }
   .current-unit {
      font-weight: bold;
      background: var(--uvalib-teal-lightest);
      padding: 2px 4px;
      border-radius: 5px;
   }
}
.last-error {
   background: var(--uvalib-red-darker);
   padding: 10px;
   color: white;
   span {
      display: inline-block;
      font-weight: bold;
      margin-right: 10px;
   }
   a {
      color: white !important;
   }
}
.details {
   padding: 0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   div.p-panel.p-component.small {
      flex: 25%;
   }
   div.note {
      margin-top: 15px;
      font-size: 0.9em;
   }
   p.note {
      padding: 0;
      margin: 10px 0 0 0;
   }
   .attachments-data {
      font-size: 0.8em;
   }
   :deep(div.p-panel) {
      margin: 10px;
      flex: 45%;
      text-align: left;
   }
   .toolbar {
      padding: 15px 0 0 0;
      text-align: right;
      font-size: 0.8em;
   }
   .acts-wrap {
      border-top: 1px solid var(--uvalib-grey-light);
      padding-top: 15px;
      .acts {
         padding: 5px 0;
         font-size: 0.8em;
         display: flex;
         flex-flow: row wrap;
         justify-content: flex-start;
         :deep(button.p-button) {
            margin-right: 10px;
         }
      }
   }
   .row-acts {
      font-size: 0.9em;
      width: 75px;
      button.p-button {
         padding: 4px 10px;
         width: 100%;
         margin-bottom: 8px;
      }
   }
}
</style>