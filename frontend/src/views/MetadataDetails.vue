<template>
   <h2>Metadata {{route.params.id}}</h2>
   <div class="details" v-if="systemStore.working==false">
      <div v-if="metadataStore.thumbURL" class="thumb">
         <a :href="metadataStore.viewerURL" target="_blank">
            <img :src="metadataStore.thumbURL" />
         </a>
      </div>
      <Panel header="General Information">
         <dl v-if="metadataStore.detail.type != 'ExternalMetadata'">
            <DataDisplay label="Type" :value="metadataStore.detail.type"/>
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

            <DataDisplay label="Manuscript/Unpublished Item" :value="formatBoolean(metadataStore.other.isManuscript)"/>
            <DataDisplay label="Personal Item" :value="formatBoolean(metadataStore.other.isPersonalItem)"/>
            <DataDisplay label="OCR Hint" :value="ocrHint"/>
            <DataDisplay label="OCR Language Hint" :value="metadataStore.other.ocrLanguageHint"/>
            <DataDisplay label="Preservation Tier" :value="preservationTier"/>
         </dl>
         <dl v-if="metadataStore.detail.externalSystem == 'ArchivesSpace'">
            <DataDisplay label="Type" :value="metadataStore.detail.externalSystem"/>
            <DataDisplay label="URL" :value="metadataStore.detail.externalURL">
               <a class="supplemental" :href="metadataStore.detail.externalURL" target="_blank">
                  {{metadataStore.detail.externalURL}}
                  <i class="icon fas fa-external-link"></i>
               </a>
            </DataDisplay>
            <DataDisplay label="Repository" :value="metadataStore.archivesSpace.repo"/>
            <DataDisplay label="Collection Title" :value="metadataStore.archivesSpace.collectionTitle"/>
            <DataDisplay label="ID" :value="metadataStore.archivesSpace.id"/>
            <DataDisplay label="Language" :value="metadataStore.archivesSpace.language"/>
            <DataDisplay label="Dates" :value="metadataStore.archivesSpace.dates"/>
            <DataDisplay label="Title" :value="metadataStore.archivesSpace.title"/>
            <DataDisplay label="Level" :value="metadataStore.archivesSpace.level"/>
            <DataDisplay label="Created By" :value="metadataStore.archivesSpace.createdBy"/>
            <DataDisplay label="Create Date" :value="metadataStore.archivesSpace.createDate"/>
         </dl>
         <dl v-if="metadataStore.detail.externalSystem == 'JSTOR Forum'">
            <DataDisplay label="Type" :value="metadataStore.detail.externalSystem"/>
            <DataDisplay label="URL" :value="metadataStore.detail.externalURL">
               <a class="supplemental" :href="metadataStore.detail.externalURL" target="_blank">
                  {{metadataStore.detail.externalURL}}
                  <i class="icon fas fa-external-link"></i>
               </a>
            </DataDisplay>
            <DataDisplay label="Collection" :value="metadataStore.jstor.collection"/>
            <DataDisplay label="Title" :value="metadataStore.jstor.title"/>
            <DataDisplay label="Description" :value="metadataStore.jstor.desc"/>
            <DataDisplay label="Creator" :value="metadataStore.jstor.creator"/>
            <DataDisplay label="Date" :value="metadataStore.jstor.date"/>
            <DataDisplay label="Width" :value="metadataStore.jstor.width"/>
            <DataDisplay label="Height" :value="metadataStore.jstor.height"/>
            <DataDisplay label="Artstor ID" :value="metadataStore.jstor.id"/>
            <DataDisplay label="Forum ID" :value="metadataStore.jstor.ssid"/>
         </dl>
         <dl v-if="metadataStore.detail.externalSystem == 'Apollo'">
            <DataDisplay label="Type" :value="metadataStore.detail.externalSystem"/>
            <DataDisplay label="URL" :value="metadataStore.apollo.itemURL">
               <a class="supplemental" :href="metadataStore.apollo.itemURL" target="_blank">
                  {{metadataStore.apollo.itemURL}}
                  <i class="icon fas fa-external-link"></i>
               </a>
            </DataDisplay>
            <DataDisplay label="Collection PID" :value="metadataStore.apollo.collectionPID">
               <a class="supplemental" :href="metadataStore.apollo.collectionURL" target="_blank">
                  {{metadataStore.apollo.collectionPID}}
                  <i class="icon fas fa-external-link"></i>
               </a>
            </DataDisplay>
            <DataDisplay label="Collection Title" :value="metadataStore.apollo.collectionTitle"/>
            <DataDisplay label="Collection Barcode" :value="metadataStore.apollo.collectionBarcode"/>
            <DataDisplay label="Collection Catalog Key" :value="metadataStore.apollo.collectionCatalogKey"/>
            <DataDisplay label="Item PID" :value="metadataStore.apollo.pid"/>
            <DataDisplay label="Item Type" :value="metadataStore.apollo.type"/>
            <DataDisplay label="Item Title" :value="metadataStore.apollo.title"/>
         </dl>
      </Panel>
      <Panel header="Digital Library Information">
         <dl>
            <DataDisplay label="PID" :value="metadataStore.dl.pid"/>
            <DataDisplay label="In Digital Library" :value="formatBoolean(metadataStore.dl.inDL)"/>
            <DataDisplay label="DPLA" :value="formatBoolean(metadataStore.dl.inDPLA)"/>
            <DataDisplay label="Right Statement" :value="useRight"/>
            <DataDisplay label="Rights Rationale" :value="metadataStore.dl.useRightRationale"/>
            <DataDisplay label="Creator Death Date" :value="metadataStore.dl.creatorDeathDate"/>
            <DataDisplay label="Availability Policy" :value="availabilityPolicy"/>
            <DataDisplay label="Collection ID" :value="metadataStore.dl.collectionID"/>
            <DataDisplay v-if="metadataStore.detail.supplementalURL" label="Supplemental System" :value="metadataStore.detail.supplementalURL">
               <a :href="metadataStore.detail.supplementalURL" target="_blank" class="supplemental">
                  {{metadataStore.detail.supplementalSystem}}<i class="icon fas fa-external-link"></i>
               </a>
            </DataDisplay>
            <template v-if="metadataStore.detail.type != 'ExternalMetadata'">
               <DataDisplay :spacer="true"/>
               <DataDisplay label="Date DL Ingest" :value="formatDate(metadataStore.dl.dateDLIngest)"/>
               <DataDisplay label="Date DL Update" :value="formatDate(metadataStore.dl.dateDLUpdate)"/>
            </template>
         </dl>
         <div v-if="canPublish" class="publish">
            <DPGButton label="Publsh to Virgo" autofocus class="p-button-secondary" @click="publishClicked()"/>
         </div>
      </Panel>
   </div>
   <template v-if="systemStore.working==false">
      <div class="more-detail">
         <Accordion v-if="metadataStore.detail.type=='XmlMetadata'">
            <AccordionTab header="XML Metadata">
               <pre class="xml">{{metadataStore.detail.xmlMetadata}}</pre>
            </AccordionTab>
         </Accordion>
      </div>
      <div class="details">
         <Panel header="Related Information">
            <TabView class="related">
               <TabPanel header="Orders">
                  <RelatedOrders :orders="metadataStore.related.orders" />
               </TabPanel>
               <TabPanel header="Units">
                  <RelatedUnits :units="metadataStore.related.units" />
               </TabPanel>
            </TabView>
         </Panel>
      </div>
   </template>
</template>

<script setup>
import { onBeforeMount, computed } from 'vue'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useMetadataStore } from '@/stores/metadata'
import Panel from 'primevue/panel'
import Accordion from 'primevue/accordion';
import AccordionTab from 'primevue/accordiontab'
import DataDisplay from '../components/DataDisplay.vue'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import dayjs from 'dayjs'
import RelatedOrders from '../components/related/RelatedOrders.vue'
import RelatedUnits from '../components/related/RelatedUnits.vue'

const route = useRoute()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()

const canPublish = computed(() => {
   if (metadataStore.dl.dateDLIngest) {
      return true
   } else {
      if (metadataStore.detail.type == 'XmlMetadata' || metadataStore.detail.type == 'SirsilMetadata') {
         let canPub = false
         metadataStore.related.units.forEach( u => {
            if (u.inDL)  {
               canPub = true
            }
         })
         return canPub
      }
   }
   return false
})

const availabilityPolicy = computed(() => {
   if ( metadataStore.dl.availability ) {
      return metadataStore.dl.availability.name
   }
   return ""
})

const preservationTier = computed(() => {
   if ( metadataStore.other.preservationTier ) {
      return metadataStore.other.preservationTier.name
   }
   return ""
})

const ocrHint = computed(() => {
   if ( metadataStore.other.ocrHint ) {
      return metadataStore.other.ocrHint.name
   }
   return ""
})

const useRight = computed(() => {
   if (metadataStore.dl.useRight ) {
      return metadataStore.dl.useRight.name
   }
   return ""
})

onBeforeRouteUpdate(async (to) => {
   let mdID = to.params.id
   metadataStore.getDetails( mdID )
})

onBeforeMount(() => {
   let mdID = route.params.id
   metadataStore.getDetails( mdID )
   document.title = `Metadata #${mdID}`
})

async function publishClicked() {
   await metadataStore.publish()
   systemStore.toastMessage('Publish Success', 'This item has successfully been published to Virgo')
}

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
   .more-detail {
      padding: 0 35px 10px 35px;
      text-align: left;
      .xml {
         font-weight: normal;
         font-size: 0.85em;
         padding: 10px;
         margin: 0;
         border-top: 0;
         white-space: pre-wrap;       /* Since CSS 2.1 */
         white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
         white-space: -pre-wrap;      /* Opera 4-6 */
         white-space: -o-pre-wrap;    /* Opera 7 */
         word-wrap: break-word;       /* Internet Explorer 5.5+ */
      }
   }
.details {
   padding:  0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   :deep(p-tabview) {
      margin: 0 !important;
   }
   a.virgo, a.supplemental {
      display: inline-block;
      margin-left: 10px;
   }
   a.supplemental {
      margin-left: 0px;
   }
   :deep(div.p-panel) {
      margin: 10px;
      flex: 40%;
      text-align: left;
   }
   .thumb {
      margin: 10px;
   }
   .empty {
      color: #ccc;
   }
   .publish {
      padding: 15px 0 0 0;
      text-align: right;
   }
}
</style>