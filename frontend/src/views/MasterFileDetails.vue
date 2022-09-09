<template>
   <h2>Master File {{route.params.id}}</h2>
   <div class="details" v-if="systemStore.working==false">
      <div class="left">
         <Panel header="General Information">
            <dl>
               <DataDisplay label="PID" :value="masterFiles.details.pid" />
               <DataDisplay label="Filename" :value="masterFiles.details.filename" />
               <DataDisplay label="Title" :value="masterFiles.details.title" />
               <DataDisplay label="Description" :value="masterFiles.details.description" />
               <DataDisplay label="Date Archived" :value="formatDate(masterFiles.details.dateArchived)" blankValue="N/A" />
               <DataDisplay label="Date DL Ingest" :value="formatDate(masterFiles.details.dateDLIngest)" blankValue="N/A" />
               <DataDisplay label="Date DL Update" :value="formatDate(masterFiles.details.dateDLUpdate)" blankValue="N/A" />
               <DataDisplay label="Tags" :value="tagList" blankValue="N/A" />
            </dl>
         </Panel>
         <Panel header="Related Information">
            <dl>
               <DataDisplay label="Metadata" :value="masterFiles.details.metadata.id">
                  <router-link :to="`/metadata/${masterFiles.details.metadata.id}`">
                     {{masterFiles.details.metadata.pid}}: {{masterFiles.details.metadata.title}}
                  </router-link>
               </DataDisplay>
               <DataDisplay label="Unit ID" :value="masterFiles.details.unitID">
                  <router-link :to="`/units/${masterFiles.details.unitID}`">{{masterFiles.details.unitID}}</router-link>
               </DataDisplay>
               <DataDisplay label="Order ID" :value="masterFiles.orderID">
                  <router-link :to="`/orders/${masterFiles.orderID}`">{{masterFiles.orderID}}</router-link>
               </DataDisplay>
               <DataDisplay label="Component ID" :value="masterFiles.details.componentID">
                  <router-link :to="`/components/${masterFiles.details.componentID}`">{{masterFiles.details.componentID}}</router-link>
               </DataDisplay>
               <DataDisplay v-if="masterFiles.details.originalID>0" label="Cloned From" :value="masterFiles.details.originalID">
                  <router-link :to="`/masterfiles/${masterFiles.details.originalID}`">{{masterFiles.details.originalID}}</router-link>
               </DataDisplay>
            </dl>
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
      <div class="thumb">
         <a :href="masterFiles.viewerURL" target="_blank">
            <img :src="masterFiles.thumbURL" />
         </a>
      </div>
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
import { useRoute,onBeforeRouteUpdate } from 'vue-router'
import Panel from 'primevue/panel'
import dayjs from 'dayjs'
import DataDisplay from '../components/DataDisplay.vue'


const route = useRoute()
const masterFiles = useMasterFilesStore()
const systemStore = useSystemStore()

const tagList = computed( () => {
   let out = []
   masterFiles.details.tags.forEach( t => out.push(t.tag))
   return out.join(", ")
})

onBeforeRouteUpdate( async (to) => {
   let mfID = to.params.id
   masterFiles.getDetails(mfID)
})

onBeforeMount(() => {
   let mfID = route.params.id
   masterFiles.getDetails(mfID)
})

function formatDate( date ) {
   if (date) {
      return dayjs(date).format("YYYY-MM-DD")
   }
   return ""
}
</script>

<style scoped lang="scss">
.details {
   padding:  0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   .left {
      flex-grow: 1;
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