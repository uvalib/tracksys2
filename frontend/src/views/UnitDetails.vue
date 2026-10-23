<template>
   <h2>
      <span>{{pageTitle}}</span>
      <div class="row-right gap-2">
         <UButton v-if="detail.reorder == false && unitsStore.hasMasterFiles && (userStore.isAdmin || userStore.isSupervisor)"
            label="Export CSV" @click="exportUnitCSV()" :loading="unitsStore.exportingCSV"
         />
         <UButton label="Audit Master Files" @click="auditUnit()"
            v-if="detail.throwAway == false && detail.reorder == false && unitsStore.hasMasterFiles && (userStore.isAdmin || userStore.isSupervisor)"
         />
         <UButton label="Edit" @click="editUnit()"/>
         <UButton label="Delete" @click="deleteUnit()" v-if="canDelete"/>
      </div>
   </h2>
   <div v-if="detail.lastError" class="row-between border-t border-b border-brand-red-B items-center bg-brand-red-100 p-2">
      <div>
         <div>Recent Error:</div>
         <div>{{detail.lastError.error}}</div>
      </div>
      <UButton :to="`/jobs/${detail.lastError.jobID}`" size="lg" icon="i-lucide-bug" color="error" class="rounded-full text-white!"/>
   </div>
   <div class="column gap-6 p-6 flex-1" v-if="systemStore.working==false">
      <div class="row-left gap-6">
         <UCard title="General Information" class="flex-1">
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
         </UCard>
         <UCard title="Digitization Information" class="flex-1">
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

            <div class="row-left gap-4 border-t border-brand-grey-100 pt-4" v-if="detail.status != 'finalizing' && (userStore.isAdmin || userStore.isSupervisor)">
               <CreateProjectDialog v-if="!detail.projectID && unitsStore.masterFiles.length == 0"/>
               <template v-if="detail.intendedUseID != 110 && unitsStore.hasMasterFiles">
                  <UButton v-if="detail.datePatronDeliverablesReady == null" @click="generateDeliverablesClicked" size="sm" color="secondary" label="Generate Deliverables"/>
                  <UButton v-else @click="generateDeliverablesClicked" color="secondary" label="Regenerate Deliverables" />
               </template>
               <UButton v-if="unitsStore.hasMasterFiles" @click="unitsStore.regenerateIIIF()" color="secondary" label="Republish IIIF"" />
               <UButton v-if="unitsStore.canDownload" @click="downloadClicked" color="secondary" label="Download from Archive" />
               <UButton label="OCR Master Files" color="secondary"  @click="unitOCRClicked()" v-if="unitsStore.canOCR && (userStore.isAdmin || userStore.isSupervisor)" />
               <UButton label="Download PDF" color="secondary" @click="unitPDFClicked()" v-if="unitsStore.canPDF" />
               <UButton v-if="unitsStore.canComplete" @click="completeClicked" color="secondary" label="Complete Unit" />
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
         </UCard>
      </div>
   
      <div class="row-left gap-6" v-if="systemStore.working==false">
         <UCard title="Related Information" class="flex-1">
            <dl>
               <DataDisplay v-if="detail.metadata" label="Metadata" :value="detail.metadata.pid">
                  <router-link :to="`/metadata/${detail.metadata.id}`">{{detail.metadata.pid}}: {{detail.metadata.title}}</router-link>
               </DataDisplay>
               <DataDisplay v-else label="Metadata" value=""/>
               <DataDisplay label="Order" :value="`${detail.orderID}`">
                  <router-link :to="`/orders/${detail.orderID}`">{{detail.orderID}}</router-link>
               </DataDisplay>
               <DataDisplay v-if="detail.relatedUnits.length > 1" label="Order Units" :value="`${detail.relatedUnits}`">
                  <div class="flex flex-row flex-wrap justify-start gap-2 pb-4">
                     <template v-for="(uid) in detail.relatedUnits.slice(0,30)" :key="`related-${uid}`">
                        <router-link :to="`/units/${uid}`" v-if="uid != detail.id">{{uid}}</router-link>
                        <span class="current-unit" v-else>{{uid}}</span>
                     </template>
                     <RelatedUnitsDialog  v-if="detail.relatedUnits.length > 30" :units="detail.relatedUnits.map(u => `${u}`)" :currentUnitID="detail.id" :orderID="detail.orderID"/>
                  </div>
               </DataDisplay>
               <DataDisplay v-if="detail.projectID"  label="Project" :value="`${detail.projectID}`">
                  <a :href="`${systemStore.projectsURL}/projects/${detail.projectID}`" target="_blank">
                     {{detail.projectID}}<i class="icon fas fa-external-link"></i>
                  </a>
               </DataDisplay>
            </dl>
         </UCard>
         <UCard title="Attachments" class="flex-1">
            <UTable v-if="detail.attachments && detail.attachments.length > 0" :data="detail.attachments" :columns="attachmentCols"> 
                <template #acts-cell="{ row }">
                  <div class="row-left gap-2">
                     <UButton label="Download" size="sm" color="secondary" icon="i-lucide-cloud-download"  @click="downloadAttachment(row.original)"/>
                     <UButton label="Delete" size="sm" color="secondary" icon="i-lucide-trash" @click="deleteAttachment(row.original)"/>
                  </div>  
                </template>
            </UTable>
            <div v-else>
               No attachments are associated with this unit.
            </div>
            <p v-if='detail.status != "approved" &&  detail.status != "done"'>
               Attachments cannot be aded to unapproved units.
            </p>
            <div v-else class="row-right border-t border-brand-grey-100 pt-4">
               <AddAttachmentDialog  v-if="systemStore.working==false"/>
            </div>
         </UCard>
      </div>
      
      <MasterFilesList v-if="unitsStore.masterFiles.length > 0" />
      <UCard  v-else title="Master Files">
         <template v-if="cloneStore.uiVisible == false">
            <p>No master files are associated with this unit.</p>
            <UButton label="Clone Existing Master Files" color="secondary" @click="cloneClicked()" />
         </template>
         <CloneMasterFiles v-else @canceled="cloneCanceled()" @cloned="cloneCompleted()" />
      </UCard>
   </div>

   <UModal v-model:open="pdfStore.downloading" :modal="true" :dismissible="false" title="Generating PDF">
      <template #body>
         <div>PDF generation in progress...</div>
         <UProgress v-model="pdfStore.percent" status />
      </template>
   </UModal>
   
   <UModal v-model:open="showPDFOptions" :modal="true" :dismissible="false" title="PDF Options">
      <template #body>
         <div v-if="bundleSubmitted == false">
            <div class="column gap-4" v-if="unitsStore.hasMultipleMetadataRecords">
               <div>This unit has multiple metadata records. How would you like the PDF to be generated?</div>
               <USelect v-model="pdfType" :items="bundleOpts" placeholder="Select an option"/>
               <div class="text-sm rounded-md p-2 border border-brand-grey-100" v-if="pdfType=='single'">
                  A singe PDF will be generated that contains images from all metadata records.
                  When generation is complete, it will automatically be downloaded.
               </div>
               <div class="text-sm rounded-md p-2 border border-brand-grey-100" v-if="pdfType=='bundle'">
                  One PDF will be generated for each metadata record. These will be bundled into
                  a single zipfile that will be placed on Digiserv Delivery. Check the job status logs for bundle
                  generation progress and a download link.
               </div>
               <div v-if="unitsStore.hasText" class="column gap-1">
                  <div>This unit has transcription or OCR text. Include it with the PDF?</div>
                  <label class="row-left gap-2 items-center">
                     <UCheckbox size="lg" v-model="pdfWithText" />
                     <span>Include Text</span>
                  </label>
               </div>
            </div>
         </div>
         <div v-else>
            Bundle generation has started. Check the job status logs for progress and a download link.
         </div>   
      </template>
      <template #footer>
         <template v-if="bundleSubmitted == false">
            <UButton label="Cancel" color="secondary" @click="showPDFOptions = false"/>
            <UButton label="Download" @click="downloadPDF()" :disabled="pdfType==''"/>
         </template>
         <UButton v-else label="Close" color="secondary" @click="showPDFOptions = false"/>   
      </template>
   </UModal>
</template>

<script setup>
import { onBeforeMount, computed, ref } from 'vue'
import { useRoute, onBeforeRouteUpdate, useRouter } from 'vue-router'
import { useCloneStore } from '@/stores/clone'
import { useSystemStore } from '@/stores/system'
import { useUnitsStore } from '@/stores/units'
import { useUserStore } from '@/stores/user'
import { usePDFStore } from '@/stores/pdf'
import { storeToRefs } from "pinia"
import DataDisplay from '../components/DataDisplay.vue'
import CreateProjectDialog from '../components/unit/CreateProjectDialog.vue'
import AddAttachmentDialog from '../components/unit/AddAttachmentDialog.vue'
import MasterFilesList from '../components/unit/MasterFilesList.vue'
import CloneMasterFiles from '../components/unit/CloneMasterFiles.vue'
import { useConfirm } from "@/composables/useConfirm"
import RelatedUnitsDialog from '../components/unit/RelatedUnitsDialog.vue'

const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()
const unitsStore = useUnitsStore()
const userStore = useUserStore()
const pdfStore = usePDFStore()
const cloneStore = useCloneStore()

const attachmentCols = [
   { accessorKey: 'filename', header: 'File Name'},
   { accessorKey: 'description', header: 'Description'},
   { accessorKey: 'acts', header: 'Actions'}
]

const showPDFOptions = ref(false)
const bundleSubmitted = ref(false)
const pdfType = ref("")
const pdfWithText = ref(false)
const bundleOpts = ref([
    { label: 'Single: One PDF for all', value: 'single' },
    { label: 'Bundle: Multiple PDFs in a zipfile', value: 'bundle' },
]);

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
   bundleSubmitted.value = false
   if ( unitsStore.hasMultipleMetadataRecords || unitsStore.hasText) {
      pdfType.value = ""
      pdfWithText.value = false
      showPDFOptions.value = true
   } else {
      pdfStore.requestPDF( unitsStore.detail.id, [], false, false )
   }
})

const downloadPDF = (async () => {
   let bundle = (pdfType.value == "bundle")
   await pdfStore.requestPDF( unitsStore.detail.id, [], bundle, pdfWithText.value )
   if ( bundle ) {
      console.log("BUNDLE! Error: ["+systemStore.error+"]")
      if ( systemStore.error == "") {
         bundleSubmitted.value = true
      }
   } else {
      console.log("NO BUNDLE. JUST CLOSE")
      showPDFOptions.value = false
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

const exportUnitCSV = (() => {
   unitsStore.exportCSV()
})

const deleteUnit = (async () => {
   const msg = 'Are you sure you want delete the selected unit? This cannot be reversed.'
   const resp = await useConfirm("Confirm Delete Unit", msg, "Delete")
   if (resp) {
      await unitsStore.deleteUnit(unitsStore.detail.id)
   }
})

const editUnit = (() => {
   router.push(`/units/${route.params.id}/edit`)
})

const downloadAttachment = ((item) => {
   let url = `${systemStore.jobsURL}/units/${unitsStore.detail.id}/attachments/${item.filename}`
   window.open(url)
})

const deleteAttachment = ( async (item) => {
   const msg = 'Are you sure you want delete the selected attachment? All data will be lost. This cannot be reversed.'
   const resp = await useConfirm("Confirm Delete Attachment", msg, "Delete")
   if (resp) {
      await unitsStore.deleteAttachment(item)
   }
})

const displayStatus = (( id) => {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
})

</script>

<style scoped lang="scss">
</style>