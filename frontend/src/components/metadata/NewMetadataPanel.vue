<template>
   <UForm :schema="schema" :state="state" class="flex flex-col gap-4" @submit="createMetadata" id="create-metadata">
      <UCard title="General Information">
         <div class="section">
            <UFormField label="Metadata Type" name="type">
               <USelect v-model="state.type" :items="metadataTypes" placeholder="Select metadata type" class="w-full" @update:modelValue="typeChanged" />
            </UFormField>
            <template v-if="state.type == 'SirsiMetadata'">
               <div class="split">
                  <UFormField name="catkey" label="Catalog Key">
                     <UInput v-model="state.catalogKey" />  
                  </UFormField>
                  <UFormField name="barcode" label="Barcode">
                     <UInput v-model="state.barcode"/>  
                  </UFormField>
                  <UButton @click="sirsiLookup()" size="sm" label="Lookup" color="secondary" 
                     :disabled="metadataStore.sirsiMatch.searching" :loading="metadataStore.sirsiMatch.searching"
                  />
               </div>
               <div class="validation" v-if="needsValidation">Lookup a new match for changes in barcode or catalog key</div>
               <div class="validation" v-if="metadataStore.sirsiMatch.error">{{ metadataStore.sirsiMatch.error }}</div>
               <div>
                  <dl>
                     <DataDisplay label="Title" :value="state.title" blankValue="Unknown"/>
                     <DataDisplay label="Call Number" :value="state.callNumber" blankValue="Unknown"/>
                  </dl>
               </div>
               <div v-if="metadataStore.sirsiMatch.metadataExists" class="md-exists">
                  <p>
                     TrackSys already contains a metadata record for this item. Details can be found
                     <router-link :to="`/metadata/${metadataStore.sirsiMatch.existingID}`">here</router-link>.
                  </p>
               </div>
            </template>
            <template v-if="state.type == 'XmlMetadata'">
               <UFormField name="title" label="Title" :required="true">
                  <UInput v-model="state.title" class="w-full" />   
               </UFormField>
               <UFormField name="author" label="Author">
                  <UInput v-model="state.author" class="w-full"/>   
               </UFormField>
            </template>
            <template v-if="state.type == 'ExternalMetadata'">
               <div class="flex flex-col gap-1">
                  <p class="note"><b>IMPORTANT</b>: Only URIs containing /resources/ or /archival_objects/ are supported.</p>
                  <p class="note">Examples:</p>
                  <ul>
                     <li>/repositories/uva-sc/resources/a_brief_survey_of_printing_history_and_practice_ma</li>
                     <li class="note">/repositories/3/resources/811</li>
                  </ul>
               </div>
               <div class="split">
                  <UFormField name="externalURI" label="External URI" :required="true" class="grow">   
                     <div class="flex flex-row gap-4">
                        <UInput v-model="state.externalURI"  @update:modelValue="needsValidation=true" class="w-full"/>  
                        <UButton @click="validateASMetadata()" label="Validate" color="secondary" :loading="metadataStore.asMatch.searching"/>
                     </div>  
                  </UFormField>
               </div>
               <div class="validation" v-if="needsValidation">Changes to external URI need to be validated</div>
               <div class="validation" v-if="metadataStore.asMatch.error">{{ metadataStore.asMatch.error }}</div>
               <dl>
                  <DataDisplay label="Title" :value="metadataStore.asMatch.title" blankValue="Unknown"/>
                  <DataDisplay label="ID" :value="metadataStore.asMatch.id" blankValue="Unknown"/>
               </dl>
            </template>
            <template v-if="state.type">
               <div class="split">
                  <UFormField name="isCollection" label="Collection" class="grow">   
                     <USelect v-model="state.isCollection"  :items="yesNo" class="w-full"/>   
                  </UFormField>
                  <UFormField name="personalItem" label="Personal Item" class="grow">   
                     <USelect v-model="state.personalItem" :items="yesNo"  class="w-full"/>   
                  </UFormField>
                  <UFormField nam="state" label="Manuscript" class="grow">   
                     <USelect v-model="state.manuscript" :items="yesNo" class="w-full"/>   
                  </UFormField>
               </div>
               <div class="split">
                  <UFormField name="ocrHint" label="OCR Hint" class="grow">   
                      <USelect v-model="state.ocrHint" :items="ocrHints" class="w-full" placeholder="Select a hint"/>   
                  </UFormField>
                  <UFormField name="ocrLanguageHint" label="OCR Language" class="grow">   
                      <USelectMenu v-model="state.ocrLanguageHint" :items="ocrLanguages" class="w-full" 
                        :disabled="isLanguageDisabled" placeholder="Select a language"/>   
                  </UFormField>
               </div>
            </template>
         </div>
      </UCard>
      <UCard v-if="state.type && state.type != 'ExternalMetadata'" title="Digital Library Information">
         <div class="section">
            <div class="split" v-if="props.collection == false">
               <UFormField name="collectionID" label="Collection ID" class="grow">   
                  <UInput v-model="state.collectionID" class="w-full" />   
               </UFormField>
               <UFormField name="collectionFacet" label="Collection Facet" class="grow">   
                  <USelect v-model="state.collectionFacet" :items="collectionFacets" placeholder="Select a facet" class="w-full" />   
               </UFormField>
            </div>
            <div class="split">
               <UFormField name="inDPLA" label="In DPLA" class="grow">   
                  <USelect v-model="state.inDPLA" :items="yesNo" class="w-full" />   
               </UFormField>
               <UFormField availabilityPolicy label="Availability Policy" :required="true" class="grow">   
                  <USelect v-model="state.availabilityPolicy" 
                     :items="availabilityPolicies" placeholder="Select a policy" class="w-full" 
                  />   
               </UFormField>
            </div>
            <div class="use-right" v-if="state.type == 'SirsiMetadata'">
               <UFormField name="useRight" label="Use Right" :required="true">
                  <USelect v-model="state.useRight" :items="useRights" placeholder="Select a right" class="w-full" />   
               </UFormField>
               <p>{{ rightStatement }}</p>
            </div>
         </div>
      </UCard>
      <div class="acts">
         {{ needsValidation }}
         <UButton @click="cancelCreate" label="Cancel" color="secondary"/>
         <UButton :label="createLabel" type="submit" :disabled="needsValidation"/> 
      </div>
   </UForm>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import DataDisplay from '@/components/DataDisplay.vue'
import { useSystemStore } from "@/stores/system"
import { useMetadataStore } from "@/stores/metadata"
import * as yup from 'yup'

const schema = yup.object().shape({
   type: yup.string().required('Metadata type is required'),
   title: yup.string().required('Title is required'),
   callNumber:  yup.number().when('type', {
      is: (value) => value == 'SirsiMetadata',
      then: (schema) => schema.required("Call number is required"),
   }),
   availabilityPolicy:  yup.number().when('type', {
      is: (value) => value == 'SirsiMetadata',
      then: (schema) => schema.min(1).required("Availability policy is required"),
   }),
   useRight:  yup.number().when('type', {
      is: (value) => value == 'SirsiMetadata',
      then: (schema) => schema.min(1).required("Use right is required"),
   }),
   externalURI: yup.string().when('type', {
      is: (value) => value == 'ExternalMetadata',
      then: (schema) => schema.required("External URI is required"),
   })
})

// this indicates that extURI, barcode or catkey have changed and need to be validated
const needsValidation = ref(false) 

// form state
const state = ref({
   type: null, 
   externalURI: "",
   title: "",
   author: "",
   callNumber: "",
   catalogKey: "",
   barcode: "",
   personalItem: false,
   manuscript: false,
   ocrHint: null,
   ocrLanguageHint: null,
   availabilityPolicy: 1,
   useRight:1,
   inDPLA: false,
   collectionID: "",
   collectionFacet: "",
   isCollection: false
})

const emit = defineEmits( ['canceled', 'created' ])
const props = defineProps({
   collection: {
      type: Boolean,
      default: false
   },
})

const systemStore = useSystemStore()
const metadataStore = useMetadataStore()

onMounted(() => {
   resetData()
})

const resetData = (() => {
   state.value = {
      type: null, 
      externalURI: "",
      title: "",
      author: "",
      callNumber: "",
      catalogKey: "",
      barcode: "",
      personalItem: false,
      manuscript: false,
      ocrHint: null,
      ocrLanguageHint: null,
      availabilityPolicy: 1,
      useRight:1,
      inDPLA: false,
      collectionID: "",
      collectionFacet: "",
      isCollection: props.collection,
   }
   needsValidation.value = false
})

const createLabel = computed(() => {
   if ( props.collection) return "Create Collection"
   return "Create Metadata"
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
   let ur = systemStore.useRights.find( r => r.id == state.value.useRight)
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
   if ( !state.value.ocrHint ) return true
   let hint = systemStore.ocrHints.find( h => h.id == state.value.ocrHint)
   return !hint.ocrCandidate
})

const typeChanged = (() => {
   const updatedType = state.value.type
   resetData()
   state.value.type = updatedType
   if ( state.value.type == 'ExternalMetadata' || state.value.type=='SirsiMetadata') {
      needsValidation.value = true
   }
})

const validateASMetadata = ( async () => {
   await metadataStore.validateArchivesSpaceURI(state.value.externalURI.trim())
   if (metadataStore.asMatch.error == "") {
      needsValidation.value = false
      state.value.externalURI = metadataStore.asMatch.validatedURL
      state.value.title = metadataStore.asMatch.title
      state.value.externalSystemID = 1
      state.value.callNumber = metadataStore.asMatch.id
   }
})

const sirsiLookup = (async () => {
   await metadataStore.sirsiLookup(state.value.barcode, state.value.catalogKey)
   if ( metadataStore.sirsiMatch.error == "") {
      needsValidation.value = false
      state.value.title = metadataStore.sirsiMatch.title
      state.value.callNumber = metadataStore.sirsiMatch.callNumber
      state.value.author = metadataStore.sirsiMatch.creatorName
      state.value.catalogKey = metadataStore.sirsiMatch.catalogKey
      state.value.barcode = metadataStore.sirsiMatch.barcode
   }
})

const cancelCreate = (() => {
   emit("canceled")
})

const createMetadata = (async () => {
   await metadataStore.create( values )
   emit("created")
})
</script>

<style lang="scss" scoped>
.acts {
   display: flex;
   flex-flow: row nowrap;
   gap: 10px;
   justify-content: flex-end;
}
.section {
   display: flex;
   flex-direction: column;
   gap: 10px;
   .validation {
      color: var(--uvalib-red-emergency);
   }
}
dl {
   margin: 0;
   display: inline-grid;
   grid-template-columns: max-content 1fr;
   grid-column-gap: 5px;
   text-align: left;
   box-sizing: border-box;
}
.note {
   margin: 0;
   font-size: 0.9em;
   b {
      font-weight: bold;
   }
}
ul {
   display: block;
   list-style-type: disc;
   margin-block-start: 0em;
   margin-block-end: 0.5em;
   padding-inline-start: 30px;
   unicode-bidi: isolate;
   font-size: 0.8em;
}
.use-right {
   p {
      margin: 10px 0 0 0;
      padding: 0;
      font-size: 0.9em;;
   }
}
.split {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: flex-end;
   gap: 10px;
}
</style>