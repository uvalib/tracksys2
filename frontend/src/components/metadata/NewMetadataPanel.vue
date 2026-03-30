<template>
   <form @submit="createMetadata" id="create-metadata">
      <Panel header="General Information">
         <FormField id="mdtype" label="Metadata Type" :error="errors.lastName" :required="true">
            <Select id="mdtype" v-model="type" placeholder="Select metadata type" @value-change="typeChanged"
               :options="metadataTypes" optionLabel="label" optionValue="value"
            />   
         </FormField>
         <template v-if="type == 'SirsiMetadata'">
            <div class="split">
               <FormField id="catkey" label="Catalog Key">
                  <InputText id="catkey" v-model="catalogKey" type="text" @update:modelValue="needsValidation=true"/>  
               </FormField>
               <FormField id="barcode" label="Barcode">
                  <InputText id="barcode" v-model="barcode" type="text" @update:modelValue="needsValidation=true"/>  
               </FormField>
               <DPGButton @click="sirsiLookup" label="Lookup" severity="secondary" :loading="metadataStore.sirsiMatch.searching"/>
            </div>
            <Message v-if="needsValidation" severity="error" size="small" variant="simple">Lookup a new match for changes in barcode or catalog key</Message>
            <Message v-if="metadataStore.sirsiMatch.error" severity="error" size="small" variant="simple">{{metadataStore.sirsiMatch.error}}</Message>
            <div>
               <dl>
                  <DataDisplay label="Title" :value="title" blankValue="Unknown"/>
                  <DataDisplay label="Call Number" :value="callNumber" blankValue="Unknown"/>
               </dl>
               <Message v-if="errors.title" severity="error" size="small" variant="simple">{{ errors.title}}</Message>
            </div>
            <div v-if="metadataStore.sirsiMatch.metadataExists" class="md-exists">
               <p>
                  TrackSys already contains a metadata record for this item. Details can be found
                  <router-link :to="`/metadata/${metadataStore.sirsiMatch.existingID}`">here</router-link>.
               </p>
            </div>
         </template>
        <template v-if="type == 'XmlMetadata'">
            <FormField id="xmltitle" label="Title" :error="errors.title" :required="true">
               <InputText id="xmltitle" v-model="title" type="text" />   
            </FormField>
            <FormField id="xmlauthor" label="Author">
               <InputText id="xmlauthor" v-model="author" type="text" />   
            </FormField>
         </template>
         <template v-if="type == 'ExternalMetadata'">
            <p class="note"><b>IMPORTANT</b>: Only URIs containing /resources/ or /archival_objects/ are supported.</p>
            <p class="note">Examples:</p>
            <ul class="note">
               <li>/repositories/uva-sc/resources/a_brief_survey_of_printing_history_and_practice_ma</li>
               <li class="note">/repositories/3/resources/811</li>
            </ul>
            <div class="split">
               <FormField id="exturi" label="External URI" :error="errors.externalURI" :required="true">
                  <div style="display: flex; flex-flow: row nowrap; gap: 10px">
                     <InputText id="exturi" type="text" v-model="externalURI"  @update:modelValue="needsValidation=true"/>  
                     <DPGButton @click="validateASMetadata" label="Validate" severity="secondary" :loading="metadataStore.asMatch.searching"/>
                  </div>  
               </FormField>
            </div>
            <Message v-if="needsValidation" severity="error" size="small" variant="simple">Changes to external URI need to be validated</Message>
            <Message v-if="metadataStore.asMatch.error" severity="error" size="small" variant="simple">{{metadataStore.asMatch.error}}</Message>
            <dl>
               <DataDisplay label="Title" :value="metadataStore.asMatch.title" blankValue="Unknown"/>
               <DataDisplay label="ID" :value="metadataStore.asMatch.id" blankValue="Unknown"/>
            </dl>
         </template>
         <template v-if="type">
            <div class="split">
               <FormField id="iscoll" label="Collection">
                  <Select id="iscoll" v-model="isCollection"  :options="yesNo" optionLabel="label" optionValue="value" />   
               </FormField>
               <FormField id="ispersonal" label="Personal Item">
                  <Select id="ispersonal" v-model="personalItem"  :options="yesNo" optionLabel="label" optionValue="value" />   
               </FormField>
               <FormField id="iscoll" label="Manuscript">
                  <Select id="iscoll" v-model="manuscript"  :options="yesNo" optionLabel="label" optionValue="value" />   
               </FormField>
            </div>
            <div class="split">
               <FormField id="ocrhint" label="OCR Hint">
                  <Select id="ocrhint" v-model="ocrHint"  :options="ocrHints" optionLabel="label" optionValue="value"  placeholder="Select a hint"/>   
               </FormField>
               <FormField id="ocrlang" label="OCR Language">
                  <Select id="ocrlang" v-model="ocrLanguageHint" :disabled="isLanguageDisabled"  
                     :options="ocrLanguages" optionLabel="label" optionValue="value"  placeholder="Select a language"
                  />   
               </FormField>
            </div>
         </template>
      </Panel>
      <Panel v-if="type && type != 'ExternalMetadata'" header="Digital Library Information">
         <div class="split" v-if="props.collection == false">
            <FormField id="collid" label="Collection ID">
               <InputText id="collid" v-model="collectionID" type="text" />   
            </FormField>
            <FormField id="cfacet" label="Collection Facet">
               <Select id="cfacet" v-model="collectionFacet"  :options="collectionFacets" optionLabel="label" optionValue="value"  placeholder="Select a facet"/>   
            </FormField>
         </div>
         <div class="split">
            <FormField id="indpla" label="In DPLA">
               <Select id="indpla" v-model="inDPLA"  :options="yesNo" optionLabel="label" optionValue="value" />   
            </FormField>
            <FormField id="availpolicy" label="Availability Policy" :error="errors.lastName" :required="true">
               <Select id="availpolicy" v-model="availabilityPolicy" 
                  :options="availabilityPolicies" optionLabel="label" optionValue="value"  placeholder="Select a policy"
               />   
            </FormField>
         </div>
         <div class="use-right" v-if="type == 'SirsiMetadata'">
            <FormField id="uright" label="Use Right"  :error="errors.useRight" :required="true">
               <Select id="uright" name="useRight"  v-model="useRight" 
                  :options="useRights" optionLabel="label" optionValue="value"  placeholder="Select a right"
               />   
            </FormField>
            <p>{{ rightStatement }}</p>
         </div>
      </Panel>
      <div class="acts">
         <DPGButton @click="cancelCreate" label="Cancel" severity="secondary"/>
         <DPGButton :label="createLabel" type="submit" /> 
      </div>
   </Form>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import DataDisplay from '@/components/DataDisplay.vue'
import Panel from 'primevue/panel'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import { useSystemStore } from "@/stores/system"
import { useMetadataStore } from "@/stores/metadata"

import { useForm } from 'vee-validate'
import * as yup from 'yup'
import FormField from '@/components/FormField.vue'
import Message from 'primevue/message'

const schema = yup.object().shape({
   type: yup.string().required('Metadata type is required'),
   title: yup.string().required('Title is required'),
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

const { values, errors, resetForm, handleSubmit, defineField, setValues } = useForm({
   validationSchema: schema
})

const [type] = defineField('type')
const [externalURI] = defineField('externalURI')
const [title] = defineField('title')
const [author] = defineField('author')
const [callNumber] = defineField('callNumber')
const [catalogKey] = defineField('catalogKey')
const [barcode] = defineField('barcode')
const [personalItem] = defineField('personalItem')
const [manuscript] = defineField('manuscript')
const [ocrHint] = defineField('ocrHint')
const [ocrLanguageHint] = defineField('ocrLanguageHint')
const [availabilityPolicy] = defineField('availabilityPolicy')
const [useRight] = defineField('useRight')
const [inDPLA] = defineField('inDPLA')
const [collectionID] = defineField('collectionID')
const [collectionFacet] = defineField('collectionFacet')
const [isCollection] = defineField('isCollection')

const emit = defineEmits( ['canceled', 'created' ])

// this indicates that extURI, barcode or catkey have changed and need to be validated
const needsValidation = ref(false) 

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
   resetForm({ 
      values: {
         type: null,
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
         availabilityPolicy: 1,
         useRight: 1,
         inDPLA: false,
         collectionID: "",
         collectionFacet: "",
         isCollection: props.collection,
      }
   })
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
   let ur = systemStore.useRights.find( r => r.id == values.useRight)
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
   if ( values.ocrHint == 0) return true
   let hint = systemStore.ocrHints.find( h => h.id == values.ocrHint)
   return !hint.ocrCandidate
})

const typeChanged = (() => {
   const newType = type.value
   resetData()
   needsValidation.value = false
   metadataStore.sirsiMatch.error = ""
   setValues({type: newType})
})

const validateASMetadata = ( async () => {
   await metadataStore.validateArchivesSpaceURI(externalURI.value.trim())
   if (metadataStore.asMatch.error == "") {
      needsValidation.value = false
      setValues({
         externalURI: metadataStore.asMatch.validatedURL,
         title: metadataStore.asMatch.title
      })
   }
})

const sirsiLookup= (async () => {
   await metadataStore.sirsiLookup(barcode.value, catalogKey.value)
   if ( metadataStore.sirsiMatch.error == "") {
      needsValidation.value = false
      setValues({
         title: metadataStore.sirsiMatch.title,
         callNumber: metadataStore.sirsiMatch.callNumber,
         author: metadataStore.sirsiMatch.creatorName,
         catalogKey: metadataStore.sirsiMatch.catalogKey,
         barcode: metadataStore.sirsiMatch.barcode,
      })
   }
})

const cancelCreate = (() => {
   emit("canceled")
})

const createMetadata = handleSubmit( async (values) => {
   if (needsValidation.value ) return
   if (values.type == "ExternalMetadata") {
      values.externalSystemID = 1
   }
   await metadataStore.create( values )
   emit("created")
})
</script>

<style lang="scss" scoped>
dl {
   margin: 0;
   display: inline-grid;
   grid-template-columns: max-content 1fr;
   grid-column-gap: 5px;
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

:deep(.p-panel-content), #create-metadata {
   display: flex;
   flex-direction: column;
   gap: 15px;
}

.split {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: flex-end;
   gap: 10px;
   .form-field {
      flex-grow: 1;
   }
}
</style>