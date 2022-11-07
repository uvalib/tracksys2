<template>
   <h2>Master File {{route.params.id}}</h2>
   <div class="masterfile-acts">
      <DPGButton label="Previous" @click="prevImage()" v-if="masterFiles.prevID > 0"/>
      <DPGButton label="Next" @click="nextImage()" v-if="masterFiles.nextID > 0"/>
      <DPGButton label="Download Image" @click="downloadImage()"/>
      <DPGButton label="Download PDF" @click="downloadPDF()"/>
      <DPGButton label="Edit" @click="editMasterFile()"/>
   </div>
   <div class="details" v-if="systemStore.working==false">
      <div class="thumb">
         <a :href="masterFiles.viewerURL" target="_blank">
            <img :src="masterFiles.thumbURL" />
         </a>
      </div>
      <div class="left">
         <Panel header="General Information">
            <dl>
               <DataDisplay label="PID" :value="masterFiles.details.pid" />
               <DataDisplay label="Filename" :value="masterFiles.details.filename" />
               <DataDisplay label="Title" :value="masterFiles.details.title" />
               <DataDisplay label="Description" :value="masterFiles.details.description" />
               <DataDisplay label="Orientation" :value="orientationName" />
               <DataDisplay label="Date Archived" :value="formatDate(masterFiles.details.dateArchived)" blankValue="N/A" />
               <DataDisplay label="Date DL Ingest" :value="formatDate(masterFiles.details.dateDLIngest)" blankValue="N/A" />
               <DataDisplay label="Date DL Update" :value="formatDate(masterFiles.details.dateDLUpdate)" blankValue="N/A" />
               <DataDisplay label="Tags" :value="tagList" blankValue="N/A" />
            </dl>
            <div class="tags">
               <TagsDialog />
            </div>
         </Panel>
         <Panel header="Related Information">
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
               <Fieldset legend="Location" v-if="location">
                  <dl>
                     <DataDisplay label="Container Type" :value="location.containerType.name"/>
                     <DataDisplay label="Container ID" :value="location.containerID"/>
                     <DataDisplay label="Folder" :value="location.folderID"/>
                     <DataDisplay label="Notes" :value="location.notes"/>
                  </dl>
               </Fieldset>
         </Panel>
      </div>
      <Panel header="Technical Information">
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
      </Panel>
   </div>
   <div class="details" v-if="masterFiles.details.transcription">
      <Panel header="Transcription">
         <pre>{{masterFiles.details.transcription}}</pre>
      </Panel>
   </div>
</template>

<script setup>
import { onBeforeMount, computed } from 'vue'
import { useMasterFilesStore } from '@/stores/masterfiles'
import { useSystemStore } from '@/stores/system'
import { useUserStore } from '@/stores/user'
import { useRouter, useRoute,onBeforeRouteUpdate } from 'vue-router'
import Panel from 'primevue/panel'
import dayjs from 'dayjs'
import DataDisplay from '../components/DataDisplay.vue'
import Fieldset from 'primevue/fieldset'
import TagsDialog from '../components/masterfile/TagsDialog.vue'

const route = useRoute()
const router = useRouter()
const masterFiles = useMasterFilesStore()
const systemStore = useSystemStore()
const userStore = useUserStore()

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

function prevImage() {
   router.push(`/masterfiles/${masterFiles.prevID}`)
}
function nextImage() {
   router.push(`/masterfiles/${masterFiles.nextID}`)
}

onBeforeRouteUpdate( async (to) => {
   let mfID = to.params.id
   masterFiles.getDetails(mfID)
})

onBeforeMount(() => {
   let mfID = route.params.id
   masterFiles.getDetails(mfID)
   document.title = `Master File #${mfID}`
})

function downloadImage() {
   masterFiles.downloadFromArchive( userStore.computeID )
}
function downloadPDF() {
   let url = `${systemStore.pdfURL}/${masterFiles.details.pid}`
   window.open(url)
}

function editMasterFile() {
   router.push(`/masterfiles/${route.params.id}/edit`)
}

function formatDate( date ) {
   if (date) {
      return dayjs(date).format("YYYY-MM-DD")
   }
   return ""
}
</script>

<style scoped lang="scss">
div.masterfile-acts {
   position: absolute;
   right:15px;
   top: 15px;
   button.p-button {
      margin-right: 5px;
      font-size: 0.9em;
   }
}
:deep(.p-fieldset) {
   margin-top: 15px;
   .p-fieldset-content {
      padding: 10px;
   }
   .p-fieldset-legend {
      background: none;
      border: 0;
      padding: 0 10px;
   }
}
.tags {
   font-size: 0.8em;
   text-align: right;
}
.details {
   padding:  0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   .left {
      display: flex;
      flex-flow: column;
      justify-content: flex-start;
      width: 50%;
   }
   :deep(div.p-panel) {
      margin: 10px;
      flex-grow: 1;
      text-align: left;
   }
   .thumb {
      margin: 10px;
   }
   .empty {
      color: #ccc;
   }
}
</style>