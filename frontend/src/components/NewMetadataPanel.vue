<template>
   <FormKit type="form" id="create-metadata" :actions="false" @submit="createMetadata">
      <Panel header="General Information" class="margin-bottom">
         <FormKit label="Metadata Type" type="select" outer-class="first" :options="metadataTypes" v-model="info.type" @change="typeChanged"/>
         <template v-if="info.type == 'SirsiMetadata'">
            <div class="split">
               <FormKit label="Catalog Key" type="text" v-model="info.catalogKey"/>
               <span class="sep"/>
               <FormKit label="Barcode" type="text" v-model="info.barcode"/>
               <span class="sep"/>
               <DPGButton @click="sirsiLookup" label="Lookup" class="p-button-secondary" :loading="metadataStore.sirsiMatch.searching"/>
            </div>
            <p v-if="metadataStore.sirsiMatch.error" class="error">{{metadataStore.sirsiMatch.error}}</p>
            <dl>
               <DataDisplay label="Title" :value="info.title" blankValue="Unknown"/>
               <DataDisplay label="Call Number" :value="info.callNumber" blankValue="Unknown"/>
            </dl>
            <div v-if="metadataStore.sirsiMatch.metadataExists" class="md-exists">
               <p>
                  TrackSys already contains a metadata record for this item. Details can be found
                  <router-link :to="`/metadata/${metadataStore.sirsiMatch.existingID}`">here</router-link>.
               </p>
            </div>
         </template>
         <template v-if="info.type == 'XmlMetadata'">
            <FormKit label="Title" type="text" v-model="info.title" required @input="xmlTitleChanged"/>
            <FormKit label="Author" type="text" v-model="info.author"/>
         </template>
         <template v-if="info.type == 'ExternalMetadata'">
            <p class="note"><b>IMPORTANT</b>: Only URIs containing /resources/, /accessions/ or /archival_objects/ are supported.</p>
            <p class="note">Examples:</p>
            <ul class="note">
               <li>/repositories/uva-sc/resources/a_brief_survey_of_printing_history_and_practice_ma</li>
               <li class="note">/repositories/3/resources/811</li>
            </ul>
            <div class="split">
               <FormKit label="External URI" type="text" v-model="info.externalURI" required/>
               <span class="sep"/>
               <DPGButton @click="validateASMetadata" label="Validate" class="p-button-secondary" :loading="metadataStore.asMatch.searching"/>
            </div>
            <p class="error" v-if="metadataStore.asMatch.error">Validation Failed: {{metadataStore.asMatch.error}}</p>
            <dl>
               <DataDisplay label="Title" :value="metadataStore.asMatch.title" blankValue="Unknown"/>
               <DataDisplay label="ID" :value="metadataStore.asMatch.id" blankValue="Unknown"/>
            </dl>
         </template>
         <div class="split">
            <FormKit label="Personal Item" type="select" :options="yesNo" v-model="info.personalItem"/>
            <span class="sep"/>
            <FormKit label="Manuscript" type="select" :options="yesNo" v-model="info.manuscript"/>
         </div>
         <div class="split">
            <FormKit label="OCR Hint" type="select" :options="ocrHints" v-model="info.ocrHint" placeholder="Select a hint"/>
            <span class="sep"/>
            <FormKit label="OCR Language" type="select" :options="ocrLanguages" :disabled="isLanguageDisabled"
               v-model="info.ocrLanguageHint" placeholder="Select a language"/>
            <span class="sep"/>
            <FormKit label="Preservation Tier" type="select" :options="preservationTiers" v-model="info.preservationTier" placeholder="Select a tier"/>
         </div>
      </Panel>
      <Panel v-if="info.type != 'ExternalMetadata'" header="Digital Library Information">
         <div class="split" v-if="props.collection == false">
            <FormKit label="Collection ID" type="text" v-model="info.collectionID"/>
            <span class="sep"/>
            <FormKit label="Collection Facet" type="select" :options="collectionFacets" v-model="info.collectionFacet" placeholder="Select a facet"/>
         </div>
         <div class="split">
            <FormKit label="In DPLA" type="select" :options="yesNo" v-model="info.inDPLA"/>
            <span class="sep"/>
            <FormKit label="Availability Policy" outer-class="first" type="select" :options="availabilityPolicies" v-model="info.availabilityPolicy" required/>
         </div>
         <div class="use-right" v-if="info.type == 'SirsiMetadata'">
            <FormKit label="Use Right" outer-class="first" type="select" :options="useRights" v-model="info.useRight" required/>
            <p>{{ rightStatement }}</p>
         </div>
      </Panel>
      <div class="acts">
         <DPGButton @click="cancelCreate" label="Cancel" class="p-button-secondary"/>
         <FormKit type="submit" :label="createLabel" :disabled="metadataStore.sirsiMatch.metadataExists" :wrapper-class="submitClass"/>
      </div>
   </FormKit>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import DataDisplay from './DataDisplay.vue'
import Panel from 'primevue/panel'
import { useSystemStore } from "@/stores/system"
import { useMetadataStore } from "@/stores/metadata"

const emit = defineEmits( ['canceled', 'created' ])

const props = defineProps({
   collection: {
      type: Boolean,
      default: false
   },
})

const systemStore = useSystemStore()
const metadataStore = useMetadataStore()

const validated = ref(false)
const info = ref({
   type: "XmlMetadata",
   externalSystemID: 0,
   externalURI: "",
   title: "",
   callNumber: "",
   author: "",
   catalogKey: "",
   barcode: "",
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
   isCollection: props.collection
})

onMounted(() => {
   validated.value = false
   info.value.type = "XmlMetadata"
   info.value.externalSystemID = 0
   info.value.externalURI = ""
   info.value.title = ""
   info.value.callNumber = ""
   info.value.author = ""
   info.value.catalogKey = ""
   info.value.barcode = ""
   info.value.personalItem = false
   info.value.manuscript = false
   info.value.ocrHint = 0
   info.value.ocrLanguageHint = ""
   info.value.preservationTier = 0
   info.value.availabilityPolicy = 1
   info.value.useRight = 1
   info.value.inDPLA = false
   info.value.collectionID = ""
   info.value.collectionFacet = ""
   info.value.isCollection = props.collection
})

const createLabel = computed(() => {
   if ( props.collection) return "Create Collection"
   return "Create Metadata"
})
const submitClass = computed(() => {
   let c = "submit-button"
   if (validated.value === false ) {
      c += " disabled"
   }
   return c
})
const availabilityPolicies = computed(() => {
   let out = []
   systemStore.availabilityPolicies.forEach( o => {
      out.push({label: o.name, value: o.id})
   })
   return out
})
const collectionFacets = computed(() => {
   let out = []
   systemStore.collectionFacets.forEach( o => {
      out.push({label: o.name, value: o.name})
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
const metadataTypes = computed(() => {
   let out = []
   out.push( {label: "Sirsi", value: "SirsiMetadata"} )
   out.push( {label: "XML", value: "XmlMetadata"} )
   out.push( {label: "ArchivesSpace", value: "ExternalMetadata"} )
   return out
})
const rightStatement = computed(() => {
   let ur = systemStore.useRights.find( r => r.id == info.value.useRight)
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
   if ( info.value.ocrHint == 0) return true
   let hint = systemStore.ocrHints.find( h => h.id == info.value.ocrHint)
   return !hint.ocrCandidate
})

function typeChanged() {
   info.value.externalSystemID = 0
   validated.value = false
   if (info.value.type == "ExternalMetadata") {
      info.value.externalSystemID = 1
   }
}
function xmlTitleChanged() {
   validated.value = ( info.value.title.length > 0)
}
async function validateASMetadata() {
   await metadataStore.validateArchivesSpaceURI(info.value.externalURI)
   if (metadataStore.asMatch.error == "") {
      validated.value = true
      info.value.externalURI = metadataStore.asMatch.validatedURL
      info.value.externalSystemID = 1
      info.value.title = metadataStore.asMatch.title
   }
}
async function sirsiLookup() {
   await metadataStore.sirsiLookup(info.value.barcode, info.value.catalogKey)
   info.value.title = metadataStore.sirsiMatch.title
   info.value.callNumber = metadataStore.sirsiMatch.callNumber
   info.value.author = metadataStore.sirsiMatch.creatorName
   info.value.catalogKey = metadataStore.sirsiMatch.catalogKey
   info.value.barcode = metadataStore.sirsiMatch.barcode
   if ( metadataStore.sirsiMatch.error == "") {
      validated.value = true
   }
}
function cancelCreate() {
   emit("canceled")
}
async function createMetadata() {
   await metadataStore.create(info.value)
   emit("created")
}
</script>

<style lang="scss" scoped>
.margin-bottom {
   margin-bottom: 15px;
}
p.error {
   color: var(--uvalib-red-emergency);
   text-align: right;
   margin:5px 0 0 0;
}
p.valid {
   margin:5px 0 0 0;
   color: var(--uvalib-green-dark);
   font-weight: bold;
   text-align: right;
}
div.p-panel {
   font-size: 0.8em;
}
dl {
   margin: 10px 0 25px 0;
   display: inline-grid;
   grid-template-columns: max-content 1fr;
   grid-column-gap: 0px;
   text-align: left;
   box-sizing: border-box;
}
p.note {
   margin: 5px 0 0 0;
}
ul.note {
   margin:5px 0 0 0;
}
.md-exists {
   text-align: center;
   p {
      padding: 0;
      margin: 0 0 15px 0;
      font-weight: bold;
      color: var(--uvalib-red-dark);
   }
   a {
      color: var(--uvalib-brand-blue-light);
      font-weight: 600;
      text-decoration: none;

      &:hover {
         text-decoration: underline;
      }
   }
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
   }
   .sep {
      display: inline-block;
      width: 10px;
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 20px 0 10px 0;
   button {
      margin-right: 10px;
   }
}
</style>