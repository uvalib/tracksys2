<template>
   <h2>
      <span>{{pageTitle}}</span>
      <div class="actions" >
         <DPGButton label="Audit Master Files" @click="auditUnit()" v-if="detail.reorder == false && unitsStore.hasMasterFiles && (userStore.isAdmin || userStore.isSupervisor)"/>
         <DPGButton label="Edit" @click="editUnit()"/>
         <DPGButton label="Delete" @click="deleteUnit()" v-if="canDelete"/>
      </div>
   </h2>
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
               <span :class="`flag ${flagString(detail.completeScan)}`">{{$formatBool(detail.completeScan)}}</span>
            </DataDisplay>
            <DataDisplay label="Throwaway" :value="flagString(detail.throwAway)">
               <span :class="`flag ${flagString(detail.throwAway)}`">{{$formatBool(detail.throwAway)}}</span>
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
               <span :class="`flag ${flagString(detail.ocrMasterFiles)}`">{{$formatBool(detail.ocrMasterFiles)}}</span>
            </DataDisplay>
            <DataDisplay label="Remove Watermark" :value="flagString(detail.removeWatermark)">
               <span :class="`flag ${flagString(detail.removeWatermark)}`">{{$formatBool(detail.removeWatermark)}}</span>
            </DataDisplay>
            <DataDisplay label="Virgo" :value="flagString(detail.includeInDL)" v-if="unitsStore.canPublishToVirgo">
               <span :class="`flag ${flagString(detail.includeInDL)}`">{{$formatBool(detail.includeInDL)}}</span>
            </DataDisplay>
            <DataDisplay label="Date Archived" :value="$formatDate(detail.dateArchived)" />
            <DataDisplay label="Date Virgo Deliverables Ready" :value="$formatDate(detail.dateDLDeliverablesReady)" v-if="detail.includeInDL" />
            <DataDisplay label="Date Patron Deliverables Ready" :value="$formatDate(detail.datePatronDeliverablesReady)" />
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
               <DPGButton v-if="unitsStore.canDownload" @click="downloadClicked" class="p-button-secondary" label="Download from Archive" />
               <DPGButton label="OCR Master Files" class="p-button-secondary"  @click="unitOCRClicked()" v-if="unitsStore.canOCR && (userStore.isAdmin || userStore.isSupervisor)" />
               <DPGButton label="Download PDF" class="p-button-secondary" @click="unitPDFClicked()" v-if="unitsStore.canPDF" />

               <DPGButton v-if="unitsStore.canComplete" @click="completeClicked" class="p-button-secondary" label="Complete Unit" />
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
                  <template v-for="(uid,idx) in detail.relatedUnits.slice(0,30)" :key="`related-${uid}`">
                     <template v-if="idx > 0"><span class="sep"></span></template>
                     <router-link :to="`/units/${uid}`" v-if="uid != detail.id">{{uid}}</router-link>
                     <span class="current-unit" v-else>{{uid}}</span>
                  </template>
                  <template v-if="detail.relatedUnits.length > 30">
                     <span class="sep"></span>
                     <span>...</span>
                  </template>
               </div>
               <div class="more-link" v-if="detail.relatedUnits.length > 30">
                  <RelatedUnitsDialog  :units="detail.relatedUnits" :currentUnitID="detail.id" :orderID="detail.orderID"/>
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
   <div class="details" v-if="unitsStore.loadingMasterFiles == false">
      <MasterFilesList v-if="unitsStore.masterFiles.length > 0" />
      <Panel header="Master Files" v-else>
         <template v-if="cloneStore.uiVisible == false">
            <p>No master files are associated with this unit.</p>
            <DPGButton label="Clone Existing Master Files" class="p-button-secondary" @click="cloneClicked()" />
         </template>
         <CloneMasterFiles v-else @canceled="cloneCanceled()" @cloned="cloneCompleted()" />
      </Panel>
   </div>
   <Dialog v-model:visible="pdfStore.downloading" :modal="true" header="Generating PDF" :style="{width: '350px'}">
      <div class="download">
         <p>PDF generation in progress...</p>
         <ProgressBar :value="pdfStore.percent"/>
      </div>
   </Dialog>
</template>

<script setup>
import { onBeforeMount, computed } from 'vue'
import { useRoute, onBeforeRouteUpdate, useRouter } from 'vue-router'
import { useCloneStore } from '@/stores/clone'
import { useSystemStore } from '@/stores/system'
import { useUnitsStore } from '@/stores/units'
import { useUserStore } from '@/stores/user'
import { usePDFStore } from '@/stores/pdf'
import Panel from 'primevue/panel'
import { storeToRefs } from "pinia"
import DataDisplay from '../components/DataDisplay.vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import CreateProjectDialog from '../components/unit/CreateProjectDialog.vue'
import AddAttachmentDialog from '../components/unit/AddAttachmentDialog.vue'
import MasterFilesList from '../components/unit/MasterFilesList.vue'
import CloneMasterFiles from '../components/unit/CloneMasterFiles.vue'
import ProgressBar from 'primevue/progressbar'
import Dialog from 'primevue/dialog'
import { useConfirm } from "primevue/useconfirm"
import RelatedUnitsDialog from '../components/unit/RelatedUnitsDialog.vue'

const confirm = useConfirm()
const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()
const unitsStore = useUnitsStore()
const userStore = useUserStore()
const pdfStore = usePDFStore()
const cloneStore = useCloneStore()

const { detail } = storeToRefs(unitsStore)

const canDelete = computed(() => {
   return userStore.isAdmin && unitsStore.masterFiles.length == 0
})

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

onBeforeRouteUpdate(async (to) => {
   cloneStore.show( false )
   let uID = to.params.id
   await unitsStore.getDetails( uID )
   unitsStore.getMasterFiles(uID)
})

onBeforeMount( async () => {
   cloneStore.show( false )
   let uID = route.params.id
   await unitsStore.getDetails( uID )
   unitsStore.getMasterFiles(uID)
   document.title = `Unit #${uID}`
})

const flagString = (( flag ) => {
   if ( flag === true) {
      return "true"
   }
   return "false"
})

const completeClicked = ( async () => {
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
})

const unitOCRClicked = (() => {
   unitsStore.startUnitOCR()
})

const unitPDFClicked = (() => {
   if (unitsStore.hasText == false ) {
      pdfStore.requestPDF( unitsStore.detail.id )
   } else {
      confirm.require({
         message: `This unit has transcription or OCR text. Include it with the PDF?`,
         header: 'Include Text',
         icon: 'pi pi-question-circle',
         rejectClass: 'p-button-secondary',
         accept: () => {
            pdfStore.requestPDF( unitsStore.detail.id, [], true )
         },
         reject: () => {
            pdfStore.requestPDF( unitsStore.detail.id )
         }
      })
   }
})

const cloneClicked = (() => {
   cloneStore.show( true )
})
const cloneCanceled = (() => {
   cloneStore.show( false )
})
const cloneCompleted = (() => {
   cloneStore.show( false )
   systemStore.toastMessage("Clone Success", 'All master files have been cloned.')
})

const downloadClicked = (() => {
   unitsStore.downloadFromArchive( userStore.computeID )
})

const generateDeliverablesClicked = (() => {
   unitsStore.generateDeliverables()
})

const auditUnit = (() => {
   unitsStore.audit()
})

const deleteUnit = (() => {
   confirm.require({
      message: 'Are you sure you want delete the selected unit? This cannot be reversed.',
      header: 'Confirm Delete Unit',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await unitsStore.deleteUnit(unitsStore.detail.id)
      }
   })
})

const editUnit = (() => {
   router.push(`/units/${route.params.id}/edit`)
})

const downloadAttachment = ((item) => {
   let url = `${systemStore.jobsURL}/units/${unitsStore.detail.id}/attachments/${item.filename}`
   window.open(url)
})

const deleteAttachment = ((item) => {
   confirm.require({
      message: 'Are you sure you want delete the selected attachment? All data will be lost. This cannot be reversed.',
      header: 'Confirm Delete Attachment',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await unitsStore.deleteAttachment(item)
      }
   })
})

const displayStatus = (( id) => {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
})

</script>

<style scoped lang="scss">
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
   .toolbar, .more-link {
      padding: 15px 0 0 0;
      text-align: right;
      font-size: 0.8em;
   }
   .more-link {
      padding-top: 0;
      justify-content: flex-end;
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
         gap: 10px;
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
.download {
   padding: 5px 15px 15px 15px;
   p {
      margin:0 0 15px 0;
   }
}
</style>