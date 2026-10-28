<template>
   <h2>
      <div class="row-left gap-2 items-center">
         <span>Master File {{route.params.id}}</span>
         <span v-if="masterFiles.details.originalID>0" class="text-sm">(Cloned from {{ masterFiles.details.originalID }})</span>
      </div>
      <div class="row-right gap-2">
         <UButton label="Previous" @click="prevImage()" v-if="masterFiles.prevID > 0"/>
         <UButton label="Next" @click="nextImage()" v-if="masterFiles.nextID > 0"/>
         <UButton label="Download Image" @click="downloadImage()" v-if="masterFiles.throwAway == false"/>
         <UButton label="Download PDF" @click="downloadPDF()" v-if="masterFiles.details.originalID==0"/>
         <UButton label="Replace" @click="replaceMasterFile()" v-if="masterFiles.details.originalID==0" :loading="masterFiles.replacing"/>
         <UButton label="OCR" @click="masterFiles.ocr()" v-if="masterFiles.isOCRCandidate  && (userStore.isAdmin || userStore.isSupervisor)"/>
         <UButton label="Edit" @click="editMasterFile()"/>
         <UButton label="Republish IIIF" @click="masterFiles.regenerateIIIF()" v-if="userStore.isAdmin && masterFiles.details.originalID==0 && masterFiles.throwAway == false" />
      </div>
   </h2>
   <div class="flex flex-row flex-wrap justify-start gap-6 p-4  text-sm" v-if="systemStore.working==false">
      <a :href="masterFiles.viewerURL" target="_blank">
         <img :src="masterFiles.thumbURL" />
      </a>
      <div class="column gap-6 flex-1">
         <UCard title="General Information">
            <dl>
               <DataDisplay label="PID" :value="masterFiles.details.pid" />
               <DataDisplay label="Filename" :value="masterFiles.details.filename" />
               <DataDisplay label="Title" :value="masterFiles.details.title" />
               <DataDisplay label="Description" :value="masterFiles.details.description" />
               <DataDisplay label="Orientation" :value="orientationName" />
               <DataDisplay label="Created" :value="$formatDateTime(masterFiles.details.createdAt)" blankValue="N/A" />
               <DataDisplay label="Date Archived" :value="$formatDate(masterFiles.details.dateArchived)" blankValue="N/A" />
               <DataDisplay label="Date DL Ingest" :value="$formatDate(masterFiles.details.dateDLIngest)" blankValue="N/A" />
               <DataDisplay label="Date DL Update" :value="$formatDate(masterFiles.details.dateDLUpdate)" blankValue="N/A" />
               <DataDisplay label="Tags" :value="tagList" blankValue="N/A" />
            </dl>
            <div class="tags">
               <TagsDialog />
            </div>
         </UCard>
         <UCard title="Location" v-if="location">
            <dl>
               <DataDisplay label="Container Type" :value="location.containerType.name"/>
               <DataDisplay label="Container ID" :value="location.containerID"/>
               <DataDisplay v-if="location.containerType.hasFolders" label="Folder" :value="location.folderID"/>
               <DataDisplay label="Notes" :value="location.notes"/>
            </dl>
         </UCard>
         <UCard title="Related Information">
            <dl>
               <template v-if="masterFiles.details.metadata">
                  <DataDisplay  label="Metadata" :value="masterFiles.details.metadata.id">
                     <router-link :to="`/metadata/${masterFiles.details.metadata.id}`">
                        {{masterFiles.details.metadata.pid}}: {{masterFiles.details.metadata.title}}
                     </router-link>
                  </DataDisplay>
                  <DataDisplay v-if="masterFiles.details.metadata.callNumber" label="Call Number" :value="masterFiles.details.metadata.callNumber">
                     <router-link :to="`/metadata/${masterFiles.details.metadata.id}`">{{masterFiles.details.metadata.callNumber}}</router-link>
                  </DataDisplay>
               </template>
               <DataDisplay v-else label="Metadata" value="" />
               <DataDisplay label="Unit ID" :value="masterFiles.details.unitID">
                  <router-link :to="`/units/${masterFiles.details.unitID}`">{{masterFiles.details.unitID}}</router-link>
               </DataDisplay>
               <DataDisplay label="Order ID" :value="masterFiles.orderID">
                  <router-link :to="`/orders/${masterFiles.orderID}`">{{masterFiles.orderID}}</router-link>
               </DataDisplay>
               <DataDisplay v-if="masterFiles.details.componentID" label="Component ID" :value="masterFiles.details.componentID">
                  <router-link :to="`/components/${masterFiles.details.componentID}`">{{masterFiles.details.componentID}}</router-link>
               </DataDisplay>
               <DataDisplay v-if="masterFiles.details.originalID>0" label="Cloned From" :value="masterFiles.details.originalID">
                  <router-link :to="`/masterfiles/${masterFiles.details.originalID}`">{{masterFiles.details.originalID}}</router-link>
               </DataDisplay>
            </dl>
         </UCard>
         <UCard title="Transcription"  v-if="masterFiles.details.transcription">
            <pre class="text-xs! whitespace-pre-wrap">{{masterFiles.details.transcription}}</pre>
         </UCard>
      </div>
      <div class="column gap-6 flex-1">
         <UCard title="Technical Information">
            <dl>
               <DataDisplay label="MD5" :value="masterFiles.details.md5" />
               <DataDisplay label="Filesize" :value="masterFiles.details.filesize" />
               <DataDisplay label="Format" :value="masterFiles.details.techMetadata.imageFormat" />
               <DataDisplay label="Height x Width" :value="masterFiles.details.techMetadata.height">
                  {{masterFiles.details.techMetadata.height}} x {{masterFiles.details.techMetadata.width}}
               </DataDisplay>
               <DataDisplay label="Resolution" :value="masterFiles.details.techMetadata.resolution" />
               <DataDisplay label="Depth" :value="masterFiles.details.techMetadata.depth" />
               <DataDisplay label="Compression" :value="masterFiles.details.techMetadata.compression" />
               <DataDisplay label="Color Space" :value="masterFiles.details.techMetadata.colorSpace" />
               <DataDisplay label="Color Profile" :value="masterFiles.details.techMetadata.colorProfile" />
               <DataDisplay label="Equipment" :value="masterFiles.details.techMetadata.equipment" />
               <DataDisplay label="Model" :value="masterFiles.details.techMetadata.model" />
               <DataDisplay label="ISO" :value="masterFiles.details.techMetadata.iso" />
               <DataDisplay label="Exposure Bias" :value="masterFiles.details.techMetadata.exposureBias" />
               <DataDisplay label="Exposure Time" :value="masterFiles.details.techMetadata.exposureTime" />
               <DataDisplay label="Aperture" :value="masterFiles.details.techMetadata.aperture" />
               <DataDisplay label="Focal Length" :value="masterFiles.details.techMetadata.focalLength" />
               <DataDisplay label="Software" :value="masterFiles.details.techMetadata.software" />
            </dl>
         </UCard>
         <UCard title="Audit Information" v-if="masterFiles.throwAway == false">
            <div class="no-audit" v-if="!masterFiles.details.audit">Not Audited</div>
            <dl v-else>
               <DataDisplay label="Audited" :value="$formatDateTime(masterFiles.details.audit.auditedAt)" />
               <DataDisplay label="Archive Exists" :value="$formatBool(masterFiles.details.audit.archiveExists)" />
               <DataDisplay label="Checksum Match" :value="$formatBool(masterFiles.details.audit.checksumMatch)" />
               <DataDisplay v-if="masterFiles.details.audit.checksumMatch==false" label="Audit Checksum" :value="masterFiles.details.audit.auditChecksum" />
               <DataDisplay label="IIIF Exists" :value="$formatBool(masterFiles.details.audit.iiifExists)" />
            </dl>
            <div class="audit-toolbar">
               <UButton @click="auditNow" severity="secondary" label="Audit Now"/>
            </div>
         </UCard>
      </div>
   </div>

   <UModal v-model:open="pdfStore.downloading" :modal="true" :dismissible="false" title="Generating PDF">
      <template #body>
         <div>PDF generation in progress...</div>
         <UProgress v-model="pdfStore.percent" status />
      </template>
   </UModal>
</template>

<script setup>
import { onBeforeMount, computed } from 'vue'
import { useMasterFilesStore } from '@/stores/masterfiles'
import { useSystemStore } from '@/stores/system'
import { useUserStore } from '@/stores/user'
import { usePDFStore } from '@/stores/pdf'
import { useRouter, useRoute,onBeforeRouteUpdate } from 'vue-router'
import DataDisplay from '../components/DataDisplay.vue'
import TagsDialog from '../components/masterfile/TagsDialog.vue'
import { useConfirm } from "@/composables/useConfirm"

const route = useRoute()
const router = useRouter()
const masterFiles = useMasterFilesStore()
const systemStore = useSystemStore()
const userStore = useUserStore()
const pdfStore = usePDFStore()

const orientationName = computed( () => {
   let names = ["Normal", "Flip Y Axis", "Rotate 90&deg;", "Rotate 180&deg;", "Rotate 270&deg;"]
   return names[masterFiles.details.techMetadata.orientation]
})
const tagList = computed( () => {
   let out = []
   masterFiles.details.tags.forEach( t => out.push(t.tag))
   return out.join(", ")
})
const location = computed(() => {
   if (masterFiles.details.locations == null) return null
   if (masterFiles.details.locations.length == 0) return null
   return masterFiles.details.locations[0]
})

onBeforeRouteUpdate( async (to) => {
   let mfID = to.params.id
   masterFiles.getDetails(mfID)
})

onBeforeMount(() => {
   let mfID = route.params.id
   masterFiles.getDetails(mfID)
   document.title = `Master File #${mfID}`
})

const prevImage = (() => {
   router.push(`/masterfiles/${masterFiles.prevID}`)
})

const nextImage = (() => {
   router.push(`/masterfiles/${masterFiles.nextID}`)
})

const downloadImage = (() => {
   masterFiles.downloadFromArchive( userStore.computeID )
})

const downloadPDF = (async () => {
   if ( masterFiles.hasText == false ) {
      pdfStore.requestPDF( masterFiles.details.unitID, [masterFiles.details.id], false, false )
   } else {
      const msg = 'This master file has transcription or OCR text. Include it with the PDF?'
      const resp = await useConfirm("Include Text", msg, "Include")
      if (resp) {
         pdfStore.requestPDF( masterFiles.details.unitID, [masterFiles.details.id], false, true )
      } else {
         pdfStore.requestPDF( masterFiles.details.unitID, [masterFiles.details.id], false, false )
      }
   }
})

const replaceMasterFile = (async () => {
   let unitDir = `${masterFiles.details.unitID}`.padStart(9, '0')
   const msg = `Replace this master file with ./finalization/unit_update/${unitDir}/${masterFiles.details.filename}?`
   const resp = await useConfirm("Confirm Replace Master File", msg, "Replace")
   if (resp) {
      masterFiles.replace()
   } 
})

const editMasterFile = (() => {
   router.push(`/masterfiles/${route.params.id}/edit`)
})

const auditNow = (() => {
   masterFiles.audit()
})
</script>

<style scoped lang="scss">
.no-audit {
   text-align: center;
   font-style: italic;
}
</style>