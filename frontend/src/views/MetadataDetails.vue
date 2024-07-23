<template>
   <h2>
      <div>
         <span v-if="metadataStore.detail.isCollection" class="collection-tag">Collection</span>
         <span>Metadata {{route.params.id}}</span>
      </div>
      <div class="actions">
         <CollectionAddDialog v-if="canAddToCollection" :metadataID="metadataStore.detail.id" />
         <template  v-if="metadataStore.detail.type == 'XmlMetadata' && systemStore.working==false">
            <DPGButton label="Delete" class="edit" @click="deleteMetadata()" v-if="canDelete"/>
            <DPGButton label="Download XML"  @click="downloadXMLClicked()" />
            <FileUpload mode="basic" name="xml" accept=".xml" :customUpload="true" @uploader="xmlUploader"
               :auto="true" chooseLabel="Upload XML" uploadIcon="" v-if="userStore.isAdmin || userStore.isSupervisor"/>
         </template>
         <DPGButton label="Edit" @click="editMetadata()"  v-if="canEdit"/>
      </div>
   </h2>
   <div class="details" v-if="systemStore.working==false">
      <div v-if="metadataStore.detail.thumbURL" class="thumb">
         <a :href="metadataStore.detail.viewerURL" target="_blank">
            <img :src="metadataStore.detail.thumbURL" />
         </a>
      </div>
      <div class="column">
         <Panel header="General Information">
            <dl>
               <template v-if="metadataStore.detail.type != 'ExternalMetadata'">
                  <DataDisplay label="Type" :value="metadataStore.detail.type"/>
                  <DataDisplay label="Catalog Key" :value="metadataStore.detail.catalogKey" v-if="metadataStore.detail.type == 'SirsiMetadata'" />
                  <DataDisplay label="Barcode" :value="metadataStore.detail.barcode" v-if="metadataStore.detail.type == 'SirsiMetadata'"/>
                  <DataDisplay label="Call Number" :value="metadataStore.detail.callNumber" v-if="metadataStore.detail.type == 'SirsiMetadata'"/>
                  <DataDisplay label="Title" :value="metadataStore.detail.title"/>
                  <template  v-if="metadataStore.detail.type == 'SirsiMetadata'">
                     <DataDisplay label="Creator Name" :value="metadataStore.detail.creatorName"/>
                     <DataDisplay label="Creator Name Type" :value="metadataStore.detail.creatorNameType"/>
                     <DataDisplay label="Year" :value="metadataStore.detail.year"/>
                     <DataDisplay label="Place of Publication" :value="metadataStore.detail.publicationPlace"/>
                     <DataDisplay label="Location" :value="metadataStore.detail.location"/>
                     <DataDisplay v-if="metadataStore.detail.folders" label="Folders" :value="metadataStore.detail.folders.length">
                        <span v-for="(l,idx) in sortedFolders">
                           <DPGButton class="folder" severity="secondary" icon="pi pi-folder-open" :label="l.folderID" @click="folderClicked(l)" />
                        </span>
                     </DataDisplay>
                  </template>
                  <DataDisplay v-if="metadataStore.related.collection" label="Collection" :value="metadataStore.related.collection.id">
                     <router-link :to="`/metadata/${metadataStore.related.collection.id}`">
                        {{ metadataStore.related.collection.title }}
                     </router-link>
                  </DataDisplay>
               </template>

               <template v-if="externalSystem == 'ArchivesSpace'">
                  <DataDisplay label="Type" :value="externalSystem"/>
                  <DataDisplay label="URL" :value="metadataStore.detail.externalURI">
                     <a class="supplemental" :href="`${metadataStore.detail.externalSystem.publicURL}${metadataStore.detail.externalURI}`" target="_blank">
                        {{metadataStore.detail.externalURI}}
                        <i class="icon fas fa-external-link"></i>
                     </a>
                  </DataDisplay>
                  <DataDisplay label="Repository" :value="metadataStore.archivesSpace.repo"/>
                  <DataDisplay label="Collection Title" :value="metadataStore.archivesSpace.collectionTitle"/>
                  <DataDisplay label="Collection ID" :value="metadataStore.archivesSpace.collectionID"/>
                  <DataDisplay label="Language" :value="metadataStore.archivesSpace.language"/>
                  <DataDisplay label="Dates" :value="metadataStore.archivesSpace.dates"/>
                  <DataDisplay label="Title" :value="metadataStore.detail.title"/>
                  <DataDisplay label="Level" :value="metadataStore.archivesSpace.level"/>
                  <DataDisplay label="Created By" :value="metadataStore.archivesSpace.createdBy"/>
                  <DataDisplay label="Create Date" :value="metadataStore.archivesSpace.createDate"/>
                  <DataDisplay v-if="metadataStore.archivesSpace.publishedAt" label="Published Date" :value="metadataStore.archivesSpace.publishedAt"/>
                  <template  v-if="metadataStore.related.collection" >
                     <DataDisplay :spacer="true"/>
                     <DataDisplay label="TrackSys Collection" :value="metadataStore.related.collection.id">
                        <router-link :to="`/metadata/${metadataStore.related.collection.id}`">
                           {{ metadataStore.related.collection.title }}
                        </router-link>
                     </DataDisplay>
                  </template>
               </template>

               <template v-if="externalSystem == 'JSTOR Forum'">
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
               </template>

               <template v-if="externalSystem == 'Apollo'">
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
               </template>

               <DataDisplay :spacer="true"/>
               <DataDisplay label="Manuscript/Unpublished Item" :value="$formatBool(metadataStore.detail.isManuscript)"/>
               <DataDisplay label="Personal Item" :value="$formatBool(metadataStore.detail.isPersonalItem)"/>
               <DataDisplay label="OCR Hint" :value="ocrHint"/>
               <DataDisplay label="OCR Language Hint" :value="metadataStore.detail.ocrLanguageHint"/>
               <DataDisplay label="Preservation Tier" :value="preservationTier"/>
            </dl>

            <template v-if="externalSystem == 'ArchivesSpace'">
               <div v-if="metadataStore.hasMasterFiles == false"  class="as-toolbar">
                  <p>Not published to ArchivesSpace - no master files.</p>
               </div>
               <div v-else-if="metadataStore.asReviewInProgress" class="as-review">
                  ArchivesSpace review has been requested.
               </div>
               <div v-else-if="metadataStore.archivesSpace.publishedAt" class="as-toolbar">
                  <DPGButton label="Unpublish" class="as-publish" @click="unpublishAS()" :loading="publishing"/>
               </div>
               <div v-else class="as-toolbar">
                  <DPGButton label="Publish now" class="as-publish" @click="publishToAS()" :loading="publishing"/>
                  <DPGButton label="Submit for review" class="as-publish" @click="submitForASReview()" :loading="publishing"/>
               </div>
               <p class="error" v-if="metadataStore.archivesSpace.error">{{metadataStore.archivesSpace.error}}</p>
            </template>
         </Panel>
         <Panel header="APTrust Information" v-if="apTrustPreservation">
            <APTrustPanel />
         </Panel>
      </div>
      <div class="column">
         <Panel header="Digital Library Information">
            <dl>
               <DataDisplay label="PID" :value="metadataStore.detail.pid"/>
               <DataDisplay label="Virgo" :value="$formatBool(metadataStore.detail.inDL)">
                  <a v-if="metadataStore.detail.inDL" class="virgo no-pad" :href="metadataStore.detail.virgoURL" target="_blank">Yes<i class="icon fas fa-external-link"></i></a>
                  <span v-else>No</span>
               </DataDisplay>
               <DataDisplay label="DPLA" :value="$formatBool(metadataStore.detail.inDPLA)"/>
               <DataDisplay label="HathiTrust" :value="$formatBool(metadataStore.detail.inHathiTrust)">
                  <div class="hathi" v-if="metadataStore.detail.inHathiTrust" @click="showHathiDialog = true">
                     <span>Yes</span>
                     <i class="icon fas fa-info-circle" aria-label="HathiTrust status"></i>
                  </div>
                  <span v-else>No</span>
               </DataDisplay>
               <template v-if="metadataStore.detail.type == 'SirsiMetadata'">
                  <DataDisplay label="Use Right" :value="metadataStore.detail.useRightName">
                     <a :href="metadataStore.detail.useRightURI" target="_blank" class="supplemental">
                        {{metadataStore.detail.useRightName}}<i class="icon fas fa-external-link"></i>
                     </a>
                  </DataDisplay>
                  <DataDisplay label="Use Right Statement" :value="metadataStore.detail.useRightStatement"/>
               </template>
               <DataDisplay label="Creator Death Date" :value="metadataStore.detail.creatorDeathDate"/>
               <DataDisplay label="Availability Policy" :value="availabilityPolicy"/>
               <DataDisplay label="Collection ID" :value="metadataStore.detail.collectionID"/>
               <DataDisplay label="Collection Facet" :value="metadataStore.detail.collectionFacet"/>
               <DataDisplay v-if="metadataStore.detail.supplementalURL" label="Supplemental System" :value="metadataStore.detail.supplementalURL">
                  <a :href="metadataStore.detail.supplementalURL" target="_blank" class="supplemental">
                     {{metadataStore.detail.supplementalSystem}}<i class="icon fas fa-external-link"></i>
                  </a>
               </DataDisplay>
               <template v-if="metadataStore.canPublishToVirgo && metadataStore.detail.dateDLIngest">
                  <DataDisplay :spacer="true"/>
                  <DataDisplay label="Virgo Ingest" :value="$formatDateTime(metadataStore.detail.dateDLIngest)"/>
                  <DataDisplay label="Virgo Update" :value="$formatDateTime(metadataStore.detail.dateDLUpdate)"/>
               </template>
            </dl>
            <div v-if="metadataStore.canPublishToVirgo" class="publish">
               <DPGButton label="Publish to Virgo" autofocus class="p-button-secondary" @click="publishVirgoClicked()" :loading="publishing"/>
            </div>
         </Panel>
      </div>
   </div>
   <template v-if="systemStore.working==false">
      <div class="more-detail">
         <Accordion value="none" v-if="metadataStore.detail.type=='XmlMetadata'">
            <AccordionPanel value="xml">
               <AccordionHeader>XML Metadata</AccordionHeader>
               <AccordionContent>
                  <pre class="xml">{{metadataStore.detail.xmlMetadata}}</pre>
               </AccordionContent>
            </AccordionPanel>
         </Accordion>
      </div>
      <div class="details">
         <Panel header="Related Information">
            <Tabs value="units" :lazy="true">
               <TabList>
                  <Tab value="collection" v-if="metadataStore.detail.isCollection">Collection Members</Tab>
                  <Tab value="orders">Orders</Tab>
                  <Tab value="units">Units</Tab>
                  <Tab value="masterfiles" v-if="metadataStore.related.masterFiles.length > 0">Master Files</Tab>
               </TabList>
               <TabPanels>
                  <TabPanel value="collection" v-if="metadataStore.detail.isCollection">
                     <WaitSpinner v-if="collectionStore.working" :overlay="true" message="Please wait..." />
                     <CollectionRecords v-else :collectionID="metadataStore.detail.id"/>
                  </TabPanel>
                  <TabPanel value="orders">
                     <RelatedOrders :orders="metadataStore.related.orders" />
                  </TabPanel>
                  <TabPanel value="units">
                     <RelatedUnits :units="metadataStore.related.units" :showMetadata="false"/>
                  </TabPanel>
                  <TabPanel v-if="metadataStore.related.masterFiles.length > 0" value="masterfiles">
                     <RelatedMasterFiles :masterFiles="metadataStore.related.masterFiles" />
                  </TabPanel>
               </TabPanels>
            </Tabs>
         </Panel>
      </div>
   </template>
   <HathiTrustDialog v-if="showHathiDialog == true" @closed="showHathiDialog=false" />
   <LocationUnitsDialog v-if="showLocUnitsDialog == true" :folder="targetFolder" @closed="showLocUnitsDialog=false" @unit="unitPicked"/>
</template>

<script setup>
import { onBeforeMount, computed, ref } from 'vue'
import { useRouter, useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useMetadataStore } from '@/stores/metadata'
import { useUserStore } from '@/stores/user'
import { useCollectionsStore } from '@/stores/collections'
import { useAPTrustStore } from '@/stores/aptrust'
import Panel from 'primevue/panel'
import Accordion from 'primevue/accordion';
import AccordionPanel from 'primevue/accordionpanel'
import AccordionHeader from 'primevue/accordionheader'
import AccordionContent from 'primevue/accordioncontent'
import DataDisplay from '../components/DataDisplay.vue'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import APTrustPanel from '@/components/aptrust/APTrustPanel.vue'
import RelatedOrders from '@/components/related/RelatedOrders.vue'
import RelatedUnits from '@/components/related/RelatedUnits.vue'
import RelatedMasterFiles from '@/components/related/RelatedMasterFiles.vue'
import FileUpload from 'primevue/fileupload'
import { useConfirm } from "primevue/useconfirm"
import CollectionRecords from '@/components/related/CollectionRecords.vue'
import HathiTrustDialog from '@/components/metadata/HathiTrustDialog.vue'
import CollectionAddDialog from '@/components/metadata/CollectionAddDialog.vue'
import LocationUnitsDialog from '@/components/metadata/LocationUnitsDialog.vue'
import WaitSpinner from '@/components/WaitSpinner.vue'

const confirm = useConfirm()
const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()
const userStore = useUserStore()
const collectionStore = useCollectionsStore()
const apTrust = useAPTrustStore()

const publishing = ref(false)
const showHathiDialog = ref(false)
const showLocUnitsDialog = ref(false)
const targetFolder = ref("")

const canAddToCollection = computed(() => {
   if ( systemStore.working ) return false
   if ( metadataStore.detail.isCollection ) return false
   if ( metadataStore.related.collection != null ) return false
   let hasDigitalCollectonUnit = false
   metadataStore.related.units.some( u => {
      if (u.intendedUse && u.intendedUse.id == 110) {
         hasDigitalCollectonUnit = true
         return hasDigitalCollectonUnit
      }
   })
   return hasDigitalCollectonUnit
})

const sortedFolders = computed(() => {
   return metadataStore.detail.folders.sort( (a,b) => {
      if ( parseInt(a.folderID,10) < parseInt(b.folderID,10)) return -1
      if ( parseInt(a.folderID,10) > parseInt(b.folderID,10)) return 1
      return 0
   })
})

const apTrustPreservation = computed( () => {
   if ( metadataStore.detail.preservationTier && metadataStore.detail.preservationTier.id > 1 ) return true
   return false
})

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

const availabilityPolicy = computed(() => {
   if ( metadataStore.detail.availabilityPolicy ) {
      return metadataStore.detail.availabilityPolicy.name
   }
   return ""
})
const preservationTier = computed(() => {
   if ( metadataStore.detail.preservationTier ) {
      return `${metadataStore.detail.preservationTier.name}: ${metadataStore.detail.preservationTier.description}`
   }
   return ""
})

const ocrHint = computed(() => {
   if ( metadataStore.detail.ocrHint ) {
      return metadataStore.detail.ocrHint.name
   }
   return ""
})

onBeforeRouteUpdate((to) => {
   let mdID = to.params.id
   loadDetails(mdID)
})

onBeforeMount( async () => {
   let mdID = route.params.id
   loadDetails(mdID)
})

const loadDetails = ( async (mdID) => {
   document.title = `Metadata #${mdID}`
   await metadataStore.getDetails( mdID )
   apTrust.clearItemStatus()
   if ( metadataStore.detail.inAPTrust) {
      apTrust.getItemStatus(metadataStore.detail.id)
   }
   if (metadataStore.detail.isCollection) {
      collectionStore.setCollection( metadataStore.detail )
      collectionStore.getItems()
   }
})

const folderClicked = (async (locInfo) => {
   await metadataStore.getLocationUnits(locInfo)
   if (systemStore.error == "") {
      if ( metadataStore.locationUnits && metadataStore.locationUnits.length == 1) {
         router.push("/units/"+metadataStore.locationUnits[0].id)
      } else {
         targetFolder.value = locInfo.folderID
         showLocUnitsDialog.value = true
      }
   }
})

const unitPicked = ((unitID) => {
   showLocUnitsDialog.value = false
   router.push("/units/"+unitID)
})

const deleteMetadata = (() => {
   confirm.require({
      message: 'Are you sure you want delete this metadata? All data will be lost. This cannot be reversed.',
      header: 'Confirm Delete Metadata',
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Delete'
      },
      accept: async () => {
         await metadataStore.deleteMetadata()
      }
   })
})

const editMetadata = (() => {
   router.push(`/metadata/${metadataStore.detail.id}/edit`)
})

const downloadXMLClicked = (() => {
   const fileURL = window.URL.createObjectURL(new Blob([metadataStore.detail.xmlMetadata], { type: 'application/xml' }))
   const fileLink = document.createElement('a')
   fileLink.href =  fileURL
   fileLink.setAttribute('download', `${metadataStore.detail.pid}.xml`)
   document.body.appendChild(fileLink)
   fileLink.click()
   window.URL.revokeObjectURL(fileURL)
})

const xmlUploader = (( event ) => {
   metadataStore.uploadXML( event.files[0] )
})

const publishVirgoClicked = (async () => {
   publishing.value = true
   await metadataStore.publish()
   publishing.value = false
   if (systemStore.error == "") {
      systemStore.toastMessage('Publish Success', 'This item has successfully been published to Virgo')
   }
})

const unpublishAS = ( async () => {
   confirm.require({
      message: 'Are you sure you want remove this item from ArchivesSpace? After this, the digitial content will no longer be publicly visible.',
      header: 'Confirm Unpublish',
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Unpublish'
      },
      accept: async () => {
         publishing.value = true
         await metadataStore.unpublishFromArchivesSpace()
         publishing.value = false
         if (systemStore.error == "") {
            systemStore.toastMessage('Unpublish Success', 'This item has successfully been removed from ArchivesSpace')
         }
      }
   })
})

const publishToAS = ( async () => {
   publishing.value = true
   await metadataStore.publishToArchivesSpace(userStore.ID)
   publishing.value = false
   if (systemStore.error == "") {
      systemStore.toastMessage('Publish Success', 'This item has successfully been published to ArchivesSpace')
   }
})

const submitForASReview = ( async () => {
   publishing.value = true
   await metadataStore.requestArchivesSpaceReview(userStore.ID)
   publishing.value = false
   if (systemStore.error == "") {
      systemStore.toastMessage('Submnission Success', 'This item has successfully been submitted for ArchivesSpace review')
   }
   console.log("PUBLISHING: "+publishing.value)
})

</script>

<style scoped lang="scss">
.collection-tag {
   display: inline-block;
   margin-right: 10px;
}
.more-detail {
   padding: 0 35px 10px 35px;
   text-align: left;
   .xml {
      font-weight: normal;
      padding: 0;
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

   .column {
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      flex: 1;
   }

   div.hathi {
      color: var(--uvalib-brand-blue-light);
      cursor: pointer;
      font-weight: bold;
      i {
         font-size: 1.2em;
      }
      &:hover {
         text-decoration: underline;
      }
   }
   .as-review {
      text-align: center;
      margin: 25px 0 10px 0;
      font-weight: bold;
   }
   .as-toolbar {
      text-align: right;
      p {
         text-align: center;
         font-weight: bold;
         margin: 10px 0 0 0;
      }
      button.as-publish {
         font-size: 0.85em;
         padding: 5px 15px;
         margin-right: 10px;
      }
   }
   a.virgo, a.supplemental {
      display: inline-block;
      margin-left: 10px;
   }
   a.virgo.no-pad {
      margin-left: 0
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
   button.folder {
      font-size: 0.85em;
      margin: 2px 4px;
      padding: 3px 8px;
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