<template>
   <h2>
      <div>
         <span v-if="metadataStore.detail.isCollection" class="collection-tag">Collection </span>
         <span>Metadata {{route.params.id}}</span>
      </div>
      <div class="row-right gap-2">
         <CollectionAddDialog v-if="canAddToCollection" :metadataID="metadataStore.detail.id"/>
         <UButton label="Delete" color="error" class="edit" @click="deleteMetadata()" v-if="canDelete"/>
         <template  v-if="metadataStore.detail.type == 'XmlMetadata' && systemStore.working==false">
            <UButton label="Download XML"  @click="downloadXMLClicked()" />
            <XmlUploadModal  v-if="userStore.isAdmin || userStore.isSupervisor" @upload="xmlUploader" />
         </template>
         <UButton label="Edit" @click="editMetadata()"  v-if="canEdit"/>
      </div>
   </h2>
   <div class="column gap-6 p-6 flex-1">
      <div class="row-left gap-6" v-if="systemStore.working==false">
         <div v-if="metadataStore.detail.thumbURL" class="thumb">
            <a :href="metadataStore.detail.viewerURL" target="_blank">
               <img :src="metadataStore.detail.thumbURL" />
            </a>
         </div>
         <div class="column flex-1">
            <UCard title="General Information" class="w-full">
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
                           <div class="flex flex-row flex-wrap gap-2">
                              <UButton v-for="l in sortedFolders" color="secondary" size="xs" icon="i-lucide-folder-open" :label="l.folderID" @click="folderClicked(l)" />
                           </div>
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
                        <UButton  v-if="metadataStore.detail.inDL" variant="link" class="rounded-none p-1" trailing-icon="i-lucide-external-link"
                           :to="`${metadataStore.detail.externalSystem.publicURL}${metadataStore.detail.externalURI}`" 
                           target="_blank" :label="metadataStore.detail.externalURI"/>
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
                     <DataDisplay label="Legacy URL" :value="metadataStore.detail.externalURI">
                        <UButton variant="link" class="rounded-none p-1" trailing-icon="i-lucide-external-link"
                           :to="`${metadataStore.detail.externalSystem.publicURL}${metadataStore.detail.externalURI}`" 
                           target="_blank" :label="metadataStore.detail.externalURI"/>
                     </DataDisplay>
                     <DataDisplay label="Search Link" :value="metadataStore.jstor.searchURL">
                        <UButton variant="link" class="rounded-none p-1" trailing-icon="i-lucide-external-link"
                           :to="metadataStore.jstor.searchURL" target="_blank" label="JSTORForum search"/>
                     </DataDisplay>
                     <DataDisplay label="Title" :value="metadataStore.detail.title"/>
                  </template>

                  <template v-if="externalSystem == 'Apollo'">
                     <DataDisplay label="Type" :value="externalSystem"/>
                     <DataDisplay label="URL" :value="metadataStore.apollo.itemURL">
                        <UButton variant="link" class="rounded-none p-1" trailing-icon="i-lucide-external-link"
                           :to="metadataStore.apollo.itemURL" target="_blank" :label="metadataStore.apollo.itemURL"/>
                     </DataDisplay>
                     <DataDisplay label="Collection PID" :value="metadataStore.apollo.collectionPID">
                        <UButton variant="link" class="rounded-none p-1" trailing-icon="i-lucide-external-link"
                           :to="metadataStore.apollo.collectionURL" target="_blank" :label="metadataStore.apollo.collectionURL"/>
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
               </dl>

               <template v-if="externalSystem == 'ArchivesSpace'">
                  <div v-if="metadataStore.hasMasterFiles == false"  class="text-right font-bold!">
                     <p>Not published to ArchivesSpace - no master files.</p>
                  </div>
                  <div v-else-if="metadataStore.asReviewInProgress" class="text-right font-bold!">
                     ArchivesSpace review has been requested.
                  </div>
                  <div v-else-if="metadataStore.archivesSpace.publishedAt" class="text-right">
                     <UButton label="Unpublish" class="as-publish" @click="unpublishAS()" :loading="publishing"/>
                  </div>
                  <div v-else class="row-right gap-2">
                     <UButton label="Publish now" class="as-publish" @click="publishToAS()" :loading="publishing"/>
                     <UButton label="Submit for review" class="as-publish" @click="submitForASReview()" :loading="publishing"/>
                  </div>
                  <p class="text-right text-brand-red-B" v-if="metadataStore.archivesSpace.error">{{metadataStore.archivesSpace.error}}</p>
               </template>
            </UCard>
         </div>
         <div class="column flex-1" >
            <UCard title="Digital Library Information" class="w-full">
               <dl>
                  <DataDisplay label="PID" :value="metadataStore.detail.pid"/>
                  <DataDisplay label="Virgo" :value="$formatBool(metadataStore.detail.inDL)">
                     <UButton  v-if="metadataStore.detail.inDL && metadataStore.detail.virgoURL" 
                        variant="link" class="rounded-none p-1" trailing-icon="i-lucide-external-link"
                        :to="metadataStore.detail.virgoURL" target="_blank" label="Yes" />
                     <span v-else>No</span>
                  </DataDisplay>
                  <DataDisplay label="DPLA" :value="$formatBool(metadataStore.detail.inDPLA)"/>
                  <DataDisplay label="HathiTrust" :value="$formatBool(metadataStore.detail.inHathiTrust)">
                     <span v-if="metadataStore.detail.inHathiTrust">Yes</span>
                     <span v-else>No</span>
                  </DataDisplay>
                  <DataDisplay  v-if="metadataStore.detail.apTrustSubmissionID" label="APTrust" :value="metadataStore.detail.apTrustSubmissionID">
                     <UButton variant="link" class="rounded-none p-1" trailing-icon="i-lucide-external-link"
                        :to="`${systemStore.apTrustURL}/submissions/${metadataStore.detail.apTrustSubmissionID}`" target="_blank" label="Yes"/>
                  </DataDisplay>
                  <template v-if="metadataStore.detail.type == 'SirsiMetadata'">
                     <DataDisplay label="Use Right" :value="metadataStore.detail.useRightName">
                        <a :href="metadataStore.detail.useRightURI" target="_blank">
                           {{metadataStore.detail.useRightName}}
                        </a>
                     </DataDisplay>
                     <DataDisplay label="Use Right Statement" :value="metadataStore.detail.useRightStatement"/>
                  </template>
                  <DataDisplay label="Creator Death Date" :value="metadataStore.detail.creatorDeathDate"/>
                  <DataDisplay label="Availability Policy" :value="availabilityPolicy"/>
                  <DataDisplay label="Collection ID" :value="metadataStore.detail.collectionID"/>
                  <DataDisplay label="Collection Facet" :value="metadataStore.detail.collectionFacet"/>
                  <DataDisplay v-if="metadataStore.detail.supplementalURL" label="Supplemental System" :value="metadataStore.detail.supplementalURL">
                     <UButton variant="link" class="rounded-none p-1" trailing-icon="i-lucide-external-link"
                        :to="metadataStore.detail.supplementalURL" target="_blank" :label="metadataStore.detail.supplementalSystem"/>
                  </DataDisplay>
                  <template v-if="metadataStore.canPublishToVirgo && metadataStore.detail.dateDLIngest">
                     <DataDisplay :spacer="true"/>
                     <DataDisplay label="Virgo Ingest" :value="$formatDateTime(metadataStore.detail.dateDLIngest)"/>
                     <DataDisplay label="Virgo Update" :value="$formatDateTime(metadataStore.detail.dateDLUpdate)"/>
                  </template>
               </dl>
               <div class="row-right gap-2">
                  <HathiTrustDialog v-if="metadataStore.detail.inHathiTrust" />
                  <UButton v-if="canSubmitAPTrust" label="Submit to APTrust" color="secondary" @click="apTrustSubmitClicked" 
                     :disabled="metadataStore.detail.apTrustSubmissionPending" :loading="metadataStore.detail.apTrustSubmissionPending"/>
                  <UButton v-if="metadataStore.canPublishToVirgo" label="Publish to Virgo" autofocus color="secondary" @click="publishVirgoClicked()" :loading="publishing"/>
               </div>
            </UCard>
         </div>
      </div>

      <UAccordion v-if="metadataStore.detail.type=='XmlMetadata'" class="w-full"
         :items="[{label: 'XML Metadata', content: metadataStore.detail.xmlMetadata}]">
         <template #content="{ item }">
            <pre class="text-sm border-l border-r border-b border-brand-grey-100 p-4">{{item.content}}</pre>
         </template>
      </UAccordion>

      <div class="row-left gap-4  flex-1">
         <UCard title="Related Information" class="w-full">
            <UTabs :items="relatedTabs" variant="link"  >
               <template #collection>
                  <WaitSpinner v-if="collectionStore.working" :overlay="true" message="Please wait..." />
                  <CollectionRecords v-else :collectionID="metadataStore.detail.id"/>
               </template>
               <template #orders>
                  <RelatedOrders :orders="metadataStore.related.orders" />
               </template>
               <template #units>
                  <RelatedUnits :units="metadataStore.related.units" :showMetadata="false"/>  
               </template>
            </UTabs>
         </UCard>
      </div>
   </div>
</template>

<script setup>
import { onBeforeMount, computed, ref } from 'vue'
import { useRouter, useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useMetadataStore } from '@/stores/metadata'
import { useUserStore } from '@/stores/user'
import { useCollectionsStore } from '@/stores/collections'
import DataDisplay from '../components/DataDisplay.vue'
import RelatedOrders from '@/components/related/RelatedOrders.vue'
import RelatedUnits from '@/components/related/RelatedUnits.vue'
import { useConfirm } from "@/composables/useConfirm"
import CollectionRecords from '@/components/related/CollectionRecords.vue'
import HathiTrustDialog from '@/components/metadata/HathiTrustDialog.vue'
import CollectionAddDialog from '@/components/metadata/CollectionAddDialog.vue'
import WaitSpinner from '@/components/WaitSpinner.vue'
import XmlUploadModal from '@/components/metadata/XmlUploadModal.vue'

const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()
const userStore = useUserStore()
const collectionStore = useCollectionsStore()

const publishing = ref(false)

const relatedTabs = computed( () => {
   let out = []
   if (metadataStore.detail.isCollection) {
      out.push({ label: 'Collection Members', slot: 'collection' })
   } else {
      out.push({ label: 'Orders', slot: 'orders' })
   }
   out.push({ label: 'Units', slot: 'units' })
   return out
})

const canSubmitAPTrust = computed(() => {
   if ( userStore.isAdmin == false ) return false
   if ( metadataStore.detail.apTrustSubmissionID ) return false
   if ( metadataStore.detail.isCollection == false ) return false
   return true
})

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

const canDelete = computed(() => {
   if (userStore.isAdmin == false) return false
   if (metadataStore.related.units.length > 0) return false
   if (metadataStore.related.orders.length > 0) return false
   if (metadataStore.related.masterFiles.length > 0) return false
   if (metadataStore.related.collection != null) return false
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
   
   if (metadataStore.detail.isCollection) {
      collectionStore.setCollection( metadataStore.detail )
      collectionStore.getItems()
   }
})

const folderClicked = (async (locInfo) => {
   await metadataStore.getLocationUnit(locInfo)
   if (systemStore.error == "") {
      router.push("/units/"+metadataStore.locationUnit.id)
   }
})

const deleteMetadata = (async () => {
   const msg = 'Are you sure you want delete this metadata? All data will be lost. This cannot be reversed.'
   const resp = await useConfirm("Confirm Delete Metadata", msg, "Delete")
   if (resp) {
      metadataStore.deleteMetadata()
   } 
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

const xmlUploader = (( uploadFile ) => {
   metadataStore.uploadXML( uploadFile )
})

const apTrustSubmitClicked = (async () => {
   const msg = 'Are you sure you want to submit this metadata record to APTrust?'
   const resp = await useConfirm("Confirm APTrust Submission", msg, "Submit")
   if (resp) {
      metadataStore.apTrustSubmit()
   } 
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
   const msg = 'Are you sure you want remove this item from ArchivesSpace? After this, the digitial content will no longer be publicly visible.'
   const resp = await useConfirm("Confirm Unpublish", msg, "Unpublish")
   if (resp) {
      publishing.value = true
      await metadataStore.unpublishFromArchivesSpace()
      publishing.value = false
      if (systemStore.error == "") {
         systemStore.toastMessage('Unpublish Success', 'This item has successfully been removed from ArchivesSpace')
      }
   }
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
})

</script>

<style scoped lang="scss">
</style>