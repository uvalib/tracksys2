<template>
   <h2>Master File {{route.params.id}}</h2>
   <div class="details" v-if="systemStore.working==false">
      <div class="left">
         <Panel header="General Information">
            <dl>
               <dt>PID</dt>
               <dd>{{masterFiles.details.pid}}</dd>
               <dt>Filename</dt>
               <dd>{{masterFiles.details.filename}}</dd>
               <dt v-if="masterFiles.details.title">Title</dt>
               <dd v-else><span class="empty">Empty</span></dd>
               <dd>{{masterFiles.details.title}}</dd>
               <dt>Description</dt>
               <dd v-if="masterFiles.details.description">{{masterFiles.details.description}}</dd>
               <dd v-else><span class="empty">Empty</span></dd>
               <dt>Date Archived</dt>
               <dd v-if="masterFiles.details.dateArchived">{{dayjs(masterFiles.details.dateArchived).format("YYYY-MM-DD")}}</dd>
               <dd v-else><span class="empty">N/A</span></dd>
               <dt>Date DL Ingest</dt>
               <dd v-if="masterFiles.details.dateDLIngest">{{dayjs(masterFiles.details.dateDLIngest).format("YYYY-MM-DD")}}</dd>
               <dd v-else><span class="empty">N/A</span></dd>
               <dt>Date DL Update</dt>
               <dd v-if="masterFiles.details.dateDLUpdate">{{dayjs(masterFiles.details.dateDLUpdate).format("YYYY-MM-DD")}}</dd>
               <dd v-else><span class="empty">N/A</span></dd>
               <dt>Tags</dt>
               <dd v-if="masterFiles.details.length > 0">{{tagList}}</dd>
               <dd v-else><span class="empty">N/A</span></dd>
            </dl>
         </Panel>
         <Panel header="Related Information">
            <dl>
               <dt>Metadata</dt>
               <dd>
                  <router-link :to="`/metadata/${masterFiles.details.metadata.id}`">
                     {{masterFiles.details.metadata.pid}}: {{masterFiles.details.metadata.title}}
                  </router-link>
               </dd>
               <dt>Unit ID</dt>
               <dd><router-link :to="`/units/${masterFiles.details.unitID}`">{{masterFiles.details.unitID}}</router-link></dd>
               <dt>Order ID</dt>
               <dd><router-link :to="`/orders/${masterFiles.orderID}`">{{masterFiles.orderID}}</router-link></dd>
               <dt>Component ID</dt>
               <dd v-if="masterFiles.details.componentID"><router-link :to="`/units/${masterFiles.details.componentID}`">{{masterFiles.details.componentID}}</router-link></dd>
               <dd v-else><span class="empty">N/A</span></dd>
            </dl>
         </Panel>
      </div>
      <Panel header="Technical Information">
         <dl>
            <dt>MD5</dt>
            <dd v-if="masterFiles.details.md5">{{masterFiles.details.md5}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
            <dt>Filesize</dt>
            <dd>{{masterFiles.details.filesize}}</dd>
            <dt>Image format</dt>
            <dd>{{masterFiles.details.techMetadata.imageFormat}}</dd>
            <dt>Height x Width</dt>
            <dd>{{masterFiles.details.techMetadata.height}} x {{masterFiles.details.techMetadata.width}}</dd>
            <dt>Resolution</dt>
            <dd>{{masterFiles.details.techMetadata.resolution}}</dd>
            <dt>Depth</dt>
            <dd>{{masterFiles.details.techMetadata.depth}}</dd>
            <dt>Compression</dt>
            <dd>{{masterFiles.details.techMetadata.compression}}</dd>
            <dt>Color space</dt>
            <dd>{{masterFiles.details.techMetadata.colorSpace}}</dd>
            <dt>Color profile</dt>
            <dd v-if="masterFiles.details.techMetadata.colorProfile">{{masterFiles.details.techMetadata.colorProfile}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
            <dt>Equipment</dt>
            <dd v-if="masterFiles.details.techMetadata.equipment">{{masterFiles.details.techMetadata.equipment}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
            <dt>Model</dt>
            <dd v-if="masterFiles.details.techMetadata.model">{{masterFiles.details.techMetadata.model}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
            <dt>ISO</dt>
            <dd v-if="masterFiles.details.techMetadata.iso">{{masterFiles.details.techMetadata.iso}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
            <dt>Exposure bias</dt>
            <dd v-if="masterFiles.details.techMetadata.exposureBias">{{masterFiles.details.techMetadata.exposureBias}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
            <dt>Exposure time</dt>
            <dd v-if="masterFiles.details.techMetadata.exposureTime">{{masterFiles.details.techMetadata.exposureTime}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
            <dt>Aperture</dt>
            <dd v-if="masterFiles.details.techMetadata.aperture">{{masterFiles.details.techMetadata.aperture}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
            <dt>Focal length</dt>
            <dd v-if="masterFiles.details.techMetadata.focalLength">{{masterFiles.details.techMetadata.focalLength}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
            <dt>Software</dt>
            <dd v-if="masterFiles.details.techMetadata.software">{{masterFiles.details.techMetadata.software}}</dd>
            <dd v-else><span class="empty">Empty</span></dd>
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
import { useRoute } from 'vue-router'
import Panel from 'primevue/panel'
import dayjs from 'dayjs'

const route = useRoute()
const masterFiles = useMasterFilesStore()
const systemStore = useSystemStore()

const tagList = computed( () => {
   let out = []
   masterFiles.details.tags.forEach( t => out.push(t.tag))
   return out.join(", ")
})

onBeforeMount(() => {
   let mfID = route.params.id
   masterFiles.getDetails(mfID)
})
</script>

<style scoped lang="scss">
.details {
   padding:  0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
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