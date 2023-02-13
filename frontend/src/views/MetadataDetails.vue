<template>
   <h2><span v-if="metadataStore.detail.isCollection" class="collection-tag">Collection</span>Metadata {{route.params.id}}</h2>
   <div class="unit-acts">
      <template  v-if="metadataStore.detail.type == 'XmlMetadata' && systemStore.working==false">
         <DPGButton label="Delete" class="edit" @click="deleteMetadata()" v-if="canDelete"/>
         <DPGButton label="Download XML"  @click="downloadXMLClicked()" />
         <FileUpload mode="basic" name="xml" accept=".xml" :customUpload="true" @uploader="xmlUploader"
            :auto="true" chooseLabel="Upload XML" uploadIcon="" v-if="userStore.isAdmin || userStore.isSupervisor"/>
      </template>
      <DPGButton label="Edit" @click="editMetadata()"  v-if="canEdit"/>
   </div>
   <div class="details" v-if="systemStore.working==false">
      <div v-if="metadataStore.thumbURL" class="thumb">
         <a :href="metadataStore.viewerURL" target="_blank">
            <img :src="metadataStore.thumbURL" />
         </a>
      </div>
      <Panel header="General Information">
         <dl v-if="metadataStore.detail.type != 'ExternalMetadata'">
            <DataDisplay label="Type" :value="metadataStore.detail.type"/>
            <DataDisplay label="Catalog Key" :value="metadataStore.detail.catalogKey" v-if="metadataStore.detail.type == 'SirsiMetadata'">
               <span>{{metadataStore.detail.catalogKey}}</span>
               <a class="virgo" :href="metadataStore.virgoURL" target="_blank">VIRGO<i class="icon fas fa-external-link"></i></a>
            </DataDisplay>
            <DataDisplay label="Barcode" :value="metadataStore.detail.barcode" v-if="metadataStore.detail.type == 'SirsiMetadata'"/>
            <DataDisplay label="Call Number" :value="metadataStore.detail.callNumber" v-if="metadataStore.detail.type == 'SirsiMetadata'"/>
            <DataDisplay label="Title" :value="metadataStore.detail.title"/>
            <template  v-if="metadataStore.detail.type == 'SirsiMetadata'">
               <DataDisplay label="Creator Name" :value="metadataStore.detail.creatorName"/>
               <DataDisplay label="Creator Name Type" :value="metadataStore.detail.creatorNameType"/>
               <DataDisplay label="Year" :value="metadataStore.detail.year"/>
               <DataDisplay label="Place of Publication" :value="metadataStore.detail.publicationPlace"/>
               <DataDisplay label="Location" :value="metadataStore.detail.location"/>
            </template>

            <DataDisplay label="Manuscript/Unpublished Item" :value="formatBoolean(metadataStore.detail.isManuscript)"/>
            <DataDisplay label="Personal Item" :value="formatBoolean(metadataStore.detail.isPersonalItem)"/>
            <DataDisplay label="OCR Hint" :value="ocrHint"/>
            <DataDisplay label="OCR Language Hint" :value="metadataStore.detail.ocrLanguageHint"/>
            <DataDisplay label="Preservation Tier" :value="preservationTier"/>
            <DataDisplay v-if="metadataStore.related.collection" label="Collection" :value="metadataStore.related.collection.id">
               <router-link :to="`/metadata/${metadataStore.related.collection.id}`">
                  {{ metadataStore.related.collection.title }}
               </router-link>
            </DataDisplay>
         </dl>
         <template v-if="externalSystem == 'ArchivesSpace'">
            <dl>
               <DataDisplay label="Type" :value="externalSystem"/>
               <DataDisplay label="URL" :value="metadataStore.detail.externalURI">
                  <a class="supplemental" :href="`${metadataStore.detail.externalSystem.publicURL}${metadataStore.detail.externalURI}`" target="_blank">
                     {{metadataStore.detail.externalURI}}
                     <i class="icon fas fa-external-link"></i>
                  </a>
               </DataDisplay>
               <DataDisplay label="Repository" :value="metadataStore.archivesSpace.repo"/>
               <DataDisplay label="Collection Title" :value="metadataStore.archivesSpace.collectionTitle"/>
               <DataDisplay label="ID" :value="metadataStore.archivesSpace.id"/>
               <DataDisplay label="Language" :value="metadataStore.archivesSpace.language"/>
               <DataDisplay label="Dates" :value="metadataStore.archivesSpace.dates"/>
               <DataDisplay label="Title" :value="metadataStore.detail.title"/>
               <DataDisplay label="Level" :value="metadataStore.archivesSpace.level"/>
               <DataDisplay label="Created By" :value="metadataStore.archivesSpace.createdBy"/>
               <DataDisplay label="Create Date" :value="metadataStore.archivesSpace.createDate"/>
               <DataDisplay v-if="metadataStore.archivesSpace.publishedAt" label="Published Date" :value="metadataStore.archivesSpace.publishedAt"/>
               <DataDisplay v-else-if="metadataStore.hasMasterFiles==false" label="Published Date" value="No Master Files. Not Published."/>
               <DataDisplay v-else label="Published Date" value="placeholder">
                  <DPGButton label="Publish Now" class="as-publish" @click="publishToAS" :loading="publishing"/>
               </DataDisplay>
            </dl>
            <p class="error" v-if="metadataStore.archivesSpace.error">{{metadataStore.archivesSpace.error}}</p>
         </template>
         <dl v-if="externalSystem == 'JSTOR Forum'">
            <DataDisplay label="Type" :value="externalSystem"/>
            <DataDisplay label="URL" :value="metadataStore.detail.externalURI">
               <a class="supplemental" :href="`${metadataStore.detail.externalSystem.publicURL}${metadataStore.detail.externalURI}`" target="_blank">
                  {{metadataStore.detail.externalURI}}
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
         <dl v-if="externalSystem == 'Apollo'">
            <DataDisplay label="Type" :value="externalSystem"/>
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
            <DataDisplay label="PID" :value="metadataStore.detail.pid"/>
            <DataDisplay label="In Digital Library" :value="formatBoolean(metadataStore.detail.inDL)"/>
            <DataDisplay label="DPLA" :value="formatBoolean(metadataStore.detail.inDPLA)"/>
            <!-- <template v-if="metadataStore.detail.type == 'SirsiMetadata'">
               <DataDisplay label="Right Statement" :value="useRight"/>
               <DataDisplay label="Rights Rationale" :value="metadataStore.detail.useRightRationale"/>
            </template> -->
            <DataDisplay label="Creator Death Date" :value="metadataStore.detail.creatorDeathDate"/>
            <DataDisplay label="Availability Policy" :value="availabilityPolicy"/>
            <DataDisplay label="Collection ID" :value="metadataStore.detail.collectionID"/>
            <DataDisplay label="Collection Facet" :value="metadataStore.detail.collectionFacet"/>
            <DataDisplay v-if="metadataStore.detail.supplementalURL" label="Supplemental System" :value="metadataStore.detail.supplementalURL">
               <a :href="metadataStore.detail.supplementalURL" target="_blank" class="supplemental">
                  {{metadataStore.detail.supplementalSystem}}<i class="icon fas fa-external-link"></i>
               </a>
            </DataDisplay>
            <template v-if="metadataStore.detail.type != 'ExternalMetadata'">
               <DataDisplay :spacer="true"/>
               <DataDisplay label="Date DL Ingest" :value="formatDate(metadataStore.detail.dateDLIngest)"/>
               <DataDisplay label="Date DL Update" :value="formatDate(metadataStore.detail.dateDLUpdate)"/>
            </template>
         </dl>
         <div v-if="canPublish" class="publish">
            <DPGButton label="Publish to Virgo" autofocus class="p-button-secondary" @click="publishClicked()" :loading="publishing"/>
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
               <TabPanel header="Collection" v-if="metadataStore.detail.isCollection">
                  <CollectionRecords :collectionID="metadataStore.detail.id"/>
               </TabPanel>
               <TabPanel header="Orders">
                  <RelatedOrders :orders="metadataStore.related.orders" />
               </TabPanel>
               <TabPanel header="Units">
                  <RelatedUnits :units="metadataStore.related.units" :showMetadata="false"/>
               </TabPanel>
               <TabPanel v-if="metadataStore.related.masterFiles.length > 0" header="Master Files">
                  <RelatedMasterFiles :masterFiles="metadataStore.related.masterFiles" />
               </TabPanel>
            </TabView>
         </Panel>
      </div>
   </template>
</template>

<script setup>
import { onBeforeMount, computed, ref } from 'vue'
import { useRouter, useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useMetadataStore } from '@/stores/metadata'
import { useUserStore } from '@/stores/user'
import Panel from 'primevue/panel'
import Accordion from 'primevue/accordion';
import AccordionTab from 'primevue/accordiontab'
import DataDisplay from '../components/DataDisplay.vue'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import dayjs from 'dayjs'
import RelatedOrders from '../components/related/RelatedOrders.vue'
import RelatedUnits from '../components/related/RelatedUnits.vue'
import RelatedMasterFiles from '../components/related/RelatedMasterFiles.vue'
import FileUpload from 'primevue/fileupload'
import { useConfirm } from "primevue/useconfirm"
import CollectionRecords from '../components/related/CollectionRecords.vue'

const confirm = useConfirm()
const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()
const userStore = useUserStore()

const publishing = ref(false)

const canDelete = computed(() => {
   if (!userStore.isAdmin && !userStore.isSupervisor) return false
   if (metadataStore.related.units.length > 0) return false
   if (metadataStore.related.orders.length > 0) return false
   return true
})

const canEdit = computed(() => {
   if (metadataStore.detail.type != 'ExternalMetadata') return true
   if (!metadataStore.detail.externalSystem) return true
   return metadataStore.detail.externalSystem.name == "ArchivesSpace"
})

const externalSystem = computed(() => {
   if (!metadataStore.detail.externalSystem) return ""
   return metadataStore.detail.externalSystem.name
})

const canPublish = computed(() => {
   if (metadataStore.detail.dateDLIngest) {
      return true
   } else {
      if (metadataStore.detail.type == 'XmlMetadata' || metadataStore.detail.type == 'SirsiMetadata') {
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
   if ( metadataStore.detail.availabilityPolicy ) {
      return metadataStore.detail.availabilityPolicy.name
   }
   return ""
})

const preservationTier = computed(() => {
   if ( metadataStore.detail.preservationTier ) {
      return metadataStore.detail.preservationTier.name
   }
   return ""
})

const ocrHint = computed(() => {
   if ( metadataStore.detail.ocrHint ) {
      return metadataStore.detail.ocrHint.name
   }
   return ""
})

// const useRight = computed(() => {
//    if (metadataStore.detail.useRight ) {
//       return metadataStore.detail.useRight.name
//    }
//    return ""
// })

onBeforeRouteUpdate(async (to) => {
   let mdID = to.params.id
   document.title = `Metadata #${mdID}`
   await metadataStore.getDetails( mdID )
})

onBeforeMount( async () => {
   let mdID = route.params.id
   document.title = `Metadata #${mdID}`
   await metadataStore.getDetails( mdID )
})

function deleteMetadata() {
   confirm.require({
      message: 'Are you sure you want delete this metadata? All data will be lost. This cannot be reversed.',
      header: 'Confirm Delete Metadata',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         await metadataStore.deleteMetadata()
      }
   })
}

function editMetadata() {
   router.push(`/metadata/${metadataStore.detail.id}/edit`)
}

function downloadXMLClicked() {
   const fileURL = window.URL.createObjectURL(new Blob([metadataStore.detail.xmlMetadata], { type: 'application/xml' }))
   const fileLink = document.createElement('a')
   fileLink.href =  fileURL
   fileLink.setAttribute('download', `${metadataStore.detail.pid}.xml`)
   document.body.appendChild(fileLink)
   fileLink.click()
   window.URL.revokeObjectURL(fileURL)
}
function xmlUploader( event ) {
   metadataStore.uploadXML( event.files[0] )
}

async function publishClicked() {
   publishing.value = true
   await metadataStore.publish()
   publishing.value = false
   if (systemStore.error == "") {
      systemStore.toastMessage('Publish Success', 'This item has successfully been published to Virgo')
   }
}

async function publishToAS() {
   publishing.value = true
   await metadataStore.publishToArchivesSpace(userStore.ID)
   publishing.value = false
   if (systemStore.error == "") {
      systemStore.toastMessage('Publish Success', 'This item has successfully been published to ArchivesSpace')
   }
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
div.unit-acts {
   position: absolute;
   right:15px;
   top: 15px;
   button.p-button {
      margin-right: 5px;
      font-size: 0.9em;
   }
   :deep(.p-fileupload.p-fileupload-basic.p-component) {
      margin-right: 5px;
      font-size: 0.9em;
      display: inline-block;
      .p-button-label {
         font-size: 0.9em;
         outline: 0;
      }
   }
}
.collection-tag {
   display: inline-block;
   margin-right: 10px;
}
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
   button.as-publish {
      font-size: 0.85em;
      padding: 5px 15px;
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
   p.error {
      color: var(--uvalib-red-emergency);
      text-align: center;
      padding: 0;
      margin: 15px 0 0 0;
   }
   .publish {
      padding: 15px 0 0 0;
      text-align: right;
   }
}
</style>