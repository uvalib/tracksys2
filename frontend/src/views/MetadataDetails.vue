<template>
   <h2>Metadata {{route.params.id}}</h2>
   <div class="details" v-if="systemStore.working==false">
      <Panel header="General Information">
         <dl>
            <DataDisplay label="Catalog Key" :value="metadataStore.detail.catalogKey">
               <span>{{metadataStore.detail.catalogKey}}</span>
               <a class="virgo" :href="metadataStore.virgoURL" target="_blank">VIRGO<i class="icon fas fa-external-link"></i></a>
            </DataDisplay>
            <DataDisplay label="Barcode" :value="metadataStore.detail.barcode"/>
            <DataDisplay label="Call Number" :value="metadataStore.detail.callNumber"/>
            <DataDisplay label="Title" :value="metadataStore.detail.title"/>
            <DataDisplay label="Creator Name" :value="metadataStore.detail.creatorName"/>
            <DataDisplay label="Creator Name Type" :value="metadataStore.detail.creatorNameType"/>
            <DataDisplay label="Year" :value="metadataStore.detail.year"/>
            <DataDisplay label="Place of Publication" :value="metadataStore.detail.publicationPlace"/>
            <DataDisplay label="Location" :value="metadataStore.detail.location"/>
         </dl>
      </Panel>
      <div class="right">
         <Panel header="Digital Library Information">
            <dl>
               <DataDisplay label="PID" :value="metadataStore.dl.pid"/>
               <DataDisplay label="In Digital Library?" :value="formatBoolean(metadataStore.dl.inDL)"/>
               <DataDisplay label="DPLA" :value="formatBoolean(metadataStore.dl.inDPLA)"/>
               <DataDisplay label="Right Statement" :value="metadataStore.dl.useRight.name"/>
               <DataDisplay label="Rights Rationale" :value="metadataStore.dl.useRightRationale"/>
               <DataDisplay label="Creator Death Date" :value="metadataStore.dl.creatorDeathDate"/>
               <DataDisplay label="Availability Policy" :value="metadataStore.dl.availability.name"/>
               <DataDisplay label="Collection Facet" :value="metadataStore.dl.collectionFacet"/>
               <DataDisplay label="Date DL Ingest" :value="formatDate(metadataStore.dl.dateDLIngest)"/>
               <DataDisplay label="Date DL Update" :value="formatDate(metadataStore.dl.dateDLUpdate)"/>
            </dl>
         </Panel>
         <Panel header="Administrative Information">
            <dl>
               <DataDisplay label="Manuscript/Unpublished Item?" :value="formatBoolean(metadataStore.other.isManuscript)"/>
               <DataDisplay label="OCR Hint" :value="metadataStore.other.ocrHint.name"/>
               <DataDisplay label="OCR Language Hint" :value="metadataStore.other.ocrHLanguage"/>
               <DataDisplay label="Preservation Tier" :value="metadataStore.other.preservationTier.name"/>
            </dl>
         </Panel>

<!-- Ocr Hint	Modern Font
Ocr Language Hint	eng
Date Created	February 17, 2022 13:47
Checked Out?	No
Preservation Tier -->
      </div>
      <div v-if="metadataStore.thumbURL" class="thumb">
         <a :href="metadataStore.viewerURL" target="_blank">
            <img :src="metadataStore.thumbURL" />
         </a>
      </div>
   </div>
</template>

<script setup>
import { onBeforeMount } from 'vue'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useMetadataStore } from '@/stores/metadata'
import Panel from 'primevue/panel'
import DataDisplay from '../components/DataDisplay.vue'
import dayjs from 'dayjs'

const route = useRoute()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()

onBeforeRouteUpdate(async (to) => {
   let mdID = to.params.id
   metadataStore.getDetails( mdID )
})

onBeforeMount(() => {
   let mdID = route.params.id
   metadataStore.getDetails( mdID )
})

function formatBoolean( flag) {
   if (flag) return "Yes"
   return "No"
}

function formatDate( date ) {
   if (date) {
      return dayjs(date).format("YYYY-MM-DD HH:mm")
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
   a.virgo {
      display: inline-block;
      margin-left: 10px;
      i.icon {
         display: inline-block;
         margin-left: 5px;
      }
   }
   .right {
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