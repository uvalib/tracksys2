<template>
   <FormKit type="form" id="create-metadata" :actions="false" @submit="createMetadata">
      <Panel header="General Information" class="margin-bottom">
         <p><b>IMPORTANT</b>: This is not fully implemented, do not use!!</p>
         <FormKit label="Metadata Type" type="select" outer-class="first" :options="metadataTypes" v-model="info.type" @change="typeChanged"/>
         <template v-if="info.type == 'SirsiMetadata'">
            <div class="split">
               <FormKit label="Catalog Key" type="text" v-model="info.catalogKey"/>
               <span class="sep"/>
               <FormKit label="Barcode" type="text" v-model="info.barcode"/>
               <span class="sep"/>
               <DPGButton @click="lookupSirsiMetadata" label="Lookup" class="p-button-secondary"/>
            </div>
            <dl>
               <DataDisplay label="Title" :value="info.title" blankValue="Unknown"/>
               <DataDisplay label="Call Number" :value="info.callNumber" blankValue="Unknown"/>
            </dl>
         </template>
         <template v-if="info.type == 'XmlMetadata'">
            <FormKit label="Title" type="text" v-model="info.title" required/>
            <FormKit label="Author" type="text" v-model="info.author"/>
         </template>
         <template v-if="info.type == 'ExternalMetadata'">
            <FormKit label="External URI" type="text" v-model="info.date"/>
            <p class="note"><b>IMPORTANT</b>: Only URIs containing /resources/, /accessions/ or /archival_objects/ are supported.</p>
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
         <div class="split">
            <FormKit label="In DPLA" type="select" :options="yesNo" v-model="info.inDPLA"/>
            <span class="sep"/>
            <FormKit label="Availability Policy" outer-class="first" type="select" :options="availabilityPolicies" v-model="info.availabilityPolicy" required/>
            <span class="sep"/>
            <FormKit label="Right Statement" outer-class="first" type="select" :options="useRights" v-model="info.useRight" required/>
         </div>
         <FormKit label="Use Right Rationale" type="textarea" :rows="2" v-model="info.useRightRationale"/>
      </Panel>
      <p><b>IMPORTANT</b>: This is not fully implemented, do not use!!</p>
      <div class="acts">
         <DPGButton @click="cancelCreate" label="Cancel" class="p-button-secondary"/>
         <FormKit type="submit" label="Create Metadata" wrapper-class="submit-button" />
      </div>
   </FormKit>
</template>

<script setup>
import { ref, computed } from 'vue'
import DataDisplay from '../DataDisplay.vue'
import Panel from 'primevue/panel'
import { useSystemStore } from "@/stores/system"
import { useMetadataStore } from "@/stores/metadata"

const emit = defineEmits( ['canceled', 'created' ])

const systemStore = useSystemStore()
const metadataStore = useMetadataStore()

const info = ref({
   type: "XmlMetadata",
   externSystemID: null,
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
   useRightRationale: "",
   inDPLA: false,
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
const metadataTypes = computed(() => {
   let out = []
   out.push( {label: "Sirsi", value: "SirsiMetadata"} )
   out.push( {label: "XML", value: "XmlMetadata"} )
   out.push( {label: "ArchivesSpace", value: "ExternalMetadata"} )
   return out
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
   info.value.externSystemID = null
   if (info.value.type == "ExternalMetadata") {
      info.value.externSystemID = 1
   }
}
function lookupSirsiMetadata() {

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

div.p-panel {
   font-size: 0.8em;
}
dl {
   margin: 10px 0;
   display: inline-grid;
   grid-template-columns: max-content 1fr;
   grid-column-gap: 0px;
   font-size: 0.8em;
   text-align: left;
   box-sizing: border-box;
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