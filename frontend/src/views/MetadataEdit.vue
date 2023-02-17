<template>
   <h2>{{pageHeader}}</h2>
   <div class="edit-form">
      <FormKit type="form" id="maetadata-edit" :actions="false" @submit="submitChanges" v-if="systemStore.working == false">
         <Panel header="General Information" class="margin-bottom">
            <template v-if="metadataStore.detail.type == 'SirsiMetadata'">
               <div class="split">
                  <FormKit label="Catalog Key" type="text" v-model="edited.catalogKey"/>
                  <span class="sep"/>
                  <FormKit label="Barcode" type="text" v-model="edited.barcode"/>
                  <span class="sep"/>
                  <DPGButton @click="sirsiLookup" label="Lookup" class="p-button-secondary" :loading="metadataStore.sirsiMatch.searching"/>
               </div>
               <p v-if="metadataStore.sirsiMatch.error" class="error">{{metadataStore.sirsiMatch.error}}</p>
               <dl>
                  <DataDisplay label="Title" :value="edited.title" blankValue="Unknown"/>
                  <DataDisplay label="Call Number" :value="edited.callNumber" blankValue="Unknown"/>
               </dl>
            </template>
            <template v-if="metadataStore.detail.type == 'XmlMetadata'">
               <p class="note"><b>Note</b>:
                  To update the XML, use the 'Download XML' button on the details page to get a copy of the data.
                  Edit it with a standalone XML editor, then upload the result using the'Upload XML' button.
               </p>
            </template>
            <template v-if="metadataStore.detail.type == 'ExternalMetadata'">
               <p class="note"><b>IMPORTANT</b>: Only URIs containing /resources/, /accessions/ or /archival_objects/ are supported.</p>
               <div class="split">
                  <FormKit label="External URI" type="text" v-model="edited.externalURI" required @input="uriChanged"/>
                  <span class="sep"/>
                  <DPGButton @click="validateASMetadata" label="Validate" class="p-button-secondary" :loading="metadataStore.asMatch.searching"/>
               </div>
               <p class="error" v-if="metadataStore.asMatch.error">Validation Failed: {{metadataStore.asMatch.error}}</p>
               <dl>
                  <DataDisplay label="Title" :value="edited.title" blankValue="Unknown"/>
               </dl>
               <p class="error" v-if="metadataStore.archivesSpace.error">{{metadataStore.archivesSpace.error}}</p>
            </template>
            <div class="split">
               <template v-if="userStore.isAdmin">
                  <FormKit label="Collection" type="select" :options="yesNo" v-model="edited.isCollection"/>
                  <span class="sep"/>
               </template>
               <FormKit label="Personal Item" type="select" :options="yesNo" v-model="edited.personalItem"/>
               <span class="sep"/>
               <FormKit label="Manuscript" type="select" :options="yesNo" v-model="edited.manuscript"/>
            </div>
            <div class="split">
               <FormKit label="OCR Hint" type="select" :options="ocrHints" v-model="edited.ocrHint" placeholder="Select a hint"/>
               <span class="sep"/>
               <FormKit label="OCR Language" type="select" :options="ocrLanguages" :disabled="isLanguageDisabled"
                  v-model="edited.ocrLanguageHint" placeholder="Select a language"/>
               <span class="sep"/>
               <FormKit label="Preservation Tier" type="select" :options="preservationTiers" v-model="edited.preservationTier" placeholder="Select a tier"/>
            </div>
         </Panel>
         <Panel v-if="metadataStore.detail.type != 'ExternalMetadata'" header="Digital Library Information">
            <div class="split">
               <FormKit label="Collection ID" type="text" v-model="edited.collectionID"/>
               <span class="sep"/>
               <FormKit label="Collection Facet" type="select" :options="collectionFacets" v-model="edited.collectionFacet" placeholder="Select a facet"/>
            </div>
            <div class="split">
               <FormKit label="In DPLA" type="select" :options="yesNo" v-model="edited.inDPLA"/>
               <span class="sep"/>
               <FormKit label="Availability Policy" outer-class="first" type="select" :options="availabilityPolicies" v-model="edited.availabilityPolicy" required/>
            </div>
            <div class="use-right" v-if="metadataStore.detail.type == 'SirsiMetadata'">
               <FormKit label="Use Right" outer-class="first" type="select" :options="useRights" v-model="edited.useRight" required/>
               <p>{{ rightStatement }}</p>
            </div>
         </Panel>
         <div class="acts">
            <DPGButton label="Cancel" class="p-button-secondary" @click="cancelEdit()"/>
            <FormKit type="submit" label="Save" :disabled="!validated" :wrapper-class="submitClass"/>
         </div>
      </FormKit>
   </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import { useMetadataStore } from '@/stores/metadata'
import { useSystemStore } from '@/stores/system'
import { useUserStore } from '@/stores/user'
import Panel from 'primevue/panel'
import DataDisplay from '@/components/DataDisplay.vue'

const route = useRoute()
const router = useRouter()
const metadataStore = useMetadataStore()
const systemStore = useSystemStore()
const userStore = useUserStore()

const validated = ref(false)
const edited = ref({
   externalURI: "",
   catalogKey: "",
   barcode: "",
   title: "",
   callNumber: "",
   author: "",
   personalItem: false,
   manuscript: false,
   ocrHint: 0,
   ocrLanguageHint: "",
   preservationTier: 0,
   availabilityPolicy: 1,
   useRight: 1,
   inDPLA: false,
   collectionID: "",
   collectionFacet: "",
   isCollection: false
})

const pageHeader = computed( () => {
   let baseHdr = `Metadata ${route.params.id}`
   if ( systemStore.working){
      return baseHdr
   }
   if ( metadataStore.detail.isCollection) {
      return "Collection "+ baseHdr
   }
   if ( metadataStore.detail.type == "SirsiMetadata") {
      return "Sirsi "+ baseHdr
   }
   if ( metadataStore.detail.type == "XmlMetadata") {
      return "XML "+ baseHdr
   }
   return "ArchivesSpace "+baseHdr

})
const collectionFacets = computed(() => {
   let out = []
   systemStore.collectionFacets.forEach( o => {
      out.push({label: o.name, value: o.name})
   })
   return out
})
const rightStatement = computed(() => {
   let ur = systemStore.useRights.find( r => r.id == edited.value.useRight)
   if (ur) {
      return ur.statement
   }
   return "Unknown"
})
const useRights = computed(() => {
   let out = []
   systemStore.useRights.forEach( o => {
      out.push({label: o.name, value: o.id})
   })
   return out
})
const yesNo = computed(() => {
   let out = []
   out.push( {label: "No", value: false} )
   out.push( {label: "Yes", value: true} )
   return out
})
const isLanguageDisabled = computed(() => {
   if ( edited.value.ocrHint == 0) return true
   let hint = systemStore.ocrHints.find( h => h.id == edited.value.ocrHint)
   return !hint.ocrCandidate
})
const availabilityPolicies = computed(() => {
   let out = []
   systemStore.availabilityPolicies.forEach( o => {
      out.push({label: o.name, value: o.id})
   })
   return out
})
const preservationTiers = computed(() => {
   let out = []
   systemStore.preservationTiers.forEach( o => {
      out.push({label: o.name, value: o.id})
   })
   return out
})
const ocrLanguages = computed(() => {
   let out = []
   systemStore.ocrLanguageHints.forEach( o => {
      out.push({label: o.language, value: o.code})
   })
   return out
})
const ocrHints = computed(() => {
   let out = []
   systemStore.ocrHints.forEach( o => {
      out.push({label: o.name, value: o.id})
   })
   return out
})
const submitClass = computed(() => {
   let c = "submit-button"
   if (validated.value === false ) {
      c += " disabled"
   }
   return c
})

onMounted( async () =>{
   let mdID = route.params.id
   await metadataStore.getDetails(mdID)
   document.title = `Edit | Metadata ${mdID}`

   edited.value.externalURI = metadataStore.detail.externalURI
   edited.value.title = metadataStore.detail.title
   edited.value.callNumber = metadataStore.detail.callNumber

   // NOTE: catalogKey and barcode may be null as they are optional data members
   // but the lookup code does not handle null values. make sure they are empty string instead
   edited.value.catalogKey = ""
   if (  metadataStore.detail.catalogKey ) {
      edited.value.catalogKey = metadataStore.detail.catalogKey
   }
   edited.value.barcode = ""
   if (  metadataStore.detail.barcode ) {
      edited.value.barcode = metadataStore.detail.barcode
   }

   edited.value.personalItem = metadataStore.detail.isPersonalItem
   edited.value.manuscript = metadataStore.detail.isManuscript
   edited.value.ocrHint = 0
   if (metadataStore.detail.ocrHint) {
      edited.value.ocrHint = metadataStore.detail.ocrHint.id
   }
   edited.value.ocrLanguageHint = metadataStore.detail.ocrLanguageHint
   edited.value.preservationTier = 0
   if (metadataStore.detail.preservationTier) {
      edited.value.ocrHint = metadataStore.detail.preservationTier.id
   }
   edited.value.availabilityPolicy = 0
   if (metadataStore.detail.availabilityPolicy) {
      edited.value.availabilityPolicy = metadataStore.detail.availabilityPolicy.id
   }
   edited.value.useRight=1
   if (metadataStore.detail.useRight) {
      edited.value.useRight = metadataStore.detail.useRight.id
   }
   edited.value.inDPLA = metadataStore.detail.inDPLA
   edited.value.isCollection = metadataStore.detail.isCollection
   edited.value.collectionID = metadataStore.detail.collectionID
   edited.value.collectionFacet = metadataStore.detail.collectionFacet
   if ( metadataStore.detail.type == "ExternalMetadata") {
      validated.value = metadataStore.archivesSpace.error != ""
   }
   if ( metadataStore.detail.type == "SirsiMetadata" && (edited.value.catalogKey != "" || edited.value.barcode !="")) {
      validated.value = true
   }
   if ( metadataStore.detail.type == "XmlMetadata" ) {
      validated.value = true
   }
})

async function sirsiLookup() {
   validated.value = false
   await metadataStore.sirsiLookup(edited.value.barcode, edited.value.catalogKey)
   edited.value.title = metadataStore.sirsiMatch.title
   edited.value.callNumber = metadataStore.sirsiMatch.callNumber
   edited.value.author = metadataStore.sirsiMatch.creatorName
   edited.value.catalogKey = metadataStore.sirsiMatch.catalogKey
   edited.value.barcode = metadataStore.sirsiMatch.barcode
   if ( metadataStore.sirsiMatch.error == "") {
      validated.value = true
   }
}
function uriChanged() {
   validated.value = false
}
async function validateASMetadata() {
   await metadataStore.validateArchivesSpaceURI(edited.value.externalURI)
   if (metadataStore.asMatch.error == "") {
      validated.value = true
      edited.value.externalURI = metadataStore.asMatch.validatedURL
      edited.value.title = metadataStore.asMatch.title
   }
}
function cancelEdit() {
   router.push(`/metadata/${route.params.id}`)
}

async function submitChanges() {
   await metadataStore.submitEdit( edited.value )
   if (systemStore.showError == false) {
      router.push(`/metadata/${metadataStore.detail.id}`)
   }
}
</script>


<style lang="scss" scoped>
.edit-form {
   width: 60%;
   margin: 30px auto 0 auto;
   text-align: left;
   p.note {
      margin: 0;
      padding: 10px;
      border: 1px solid var(--uvalib-teal-light);
      background: var(--uvalib-teal-lightest);
      border-radius: 3px;
      text-align: left;
   }
   p.error {
      color: var(--uvalib-red-emergency);
      padding: 0;
      margin: 15px 0 0 0;
   }

   .margin-bottom {
      margin-bottom: 15px;
   }

   .split {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: flex-end;
      :deep(.formkit-outer) {
         flex-grow: 1;
      }
      .p-button {
         margin-bottom: 0.3em;
         font-size: 0.8em;
      }
      .sep {
         display: inline-block;
         width: 10px;
      }
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 25px 0;
   button {
      margin-right: 10px;
   }
}
</style>
