<template>
   <h2>{{pageHeader}}</h2>
   <div class="edit-form">
      <form id="maetadata-edit" @submit="submitChanges" v-if="systemStore.working == false">
         <Panel header="General Information">
            <template v-if="metadataStore.detail.type == 'SirsiMetadata'">
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
               <dl>
                  <DataDisplay label="Title" :value="title" blankValue="Unknown"/>
                  <DataDisplay label="Call Number" :value="callNumber" blankValue="Unknown"/>
               </dl>
            </template>
            <template v-if="metadataStore.detail.type == 'XmlMetadata'">
               <p class="note"><b>Note</b>:
                  To update the XML, use the 'Download XML' button on the details page to get a copy of the data.
                  Edit it with a standalone XML editor, then upload the result using the'Upload XML' button.
               </p>
            </template>
            <template v-if="metadataStore.detail.type == 'ExternalMetadata'">
               <p class="note"><b>IMPORTANT</b>: Only URIs containing /resources/ or /archival_objects/ are supported.</p>
               <div class="split">
                  <FormField id="exturi" label="External URI" :error="errors.externalURI" :required="true">
                     <div style="display: flex; flex-flow: row nowrap; gap: 10px">
                        <InputText id="exturi" type="text" v-model="externalURI" @update:modelValue="needsValidation=true" />   
                        <DPGButton @click="validateASMetadata" label="Validate" severity="secondary" :loading="metadataStore.asMatch.searching"/>
                     </div>
                  </FormField>
                 
               </div>
               <Message v-if="needsValidation" severity="error" size="small" variant="simple">Changes to external URI need to be validated</Message>
               <Message v-if="metadataStore.asMatch.error" severity="error" size="small" variant="simple">{{metadataStore.asMatch.error}}</Message>
               <dl>
                  <DataDisplay label="Title" :value="title" blankValue="Unknown"/>
               </dl>
            </template>
            <div class="split">
               <FormField v-if="userStore.isAdmin" id="iscoll" label="Collection">
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
         </Panel>
         <Panel v-if="metadataStore.detail.type != 'ExternalMetadata'" header="Digital Library Information">
            <div class="split">
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
               <FormField id="inhathi" label="In HathiTrust">
                  <Select id="inhathi" v-model="inHathiTrust" :options="yesNo" optionLabel="label" optionValue="value" :disabled="submittedToHathiTrust"/>   
               </FormField>
               <FormField id="availpolicy" label="Availability Policy" :error="errors.lastName" :required="true">
                  <Select id="availpolicy" v-model="availabilityPolicy" 
                     :options="availabilityPolicies" optionLabel="label" optionValue="value"  placeholder="Select a policy"
                  />   
               </FormField>
            </div>
            <div class="use-right" v-if="metadataStore.detail.type == 'SirsiMetadata'">
               <FormField id="uright" label="Use Right"  :error="errors.useRight" :required="true">
                  <Select id="uright" name="useRight"  v-model="useRight" 
                     :options="useRights" optionLabel="label" optionValue="value"  placeholder="Select a right"
                  />   
               </FormField>
               <p>{{ rightStatement }}</p>
            </div>
         </Panel>
         <div class="acts">
            <DPGButton label="Cancel" severity="secondary" @click="cancelEdit()"/>
            <DPGButton label="Save" type="submit" /> 
         </div>
      </form>
   </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import { useMetadataStore } from '@/stores/metadata'
import { useSystemStore } from '@/stores/system'
import { useUserStore } from '@/stores/user'
import Panel from 'primevue/panel'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import DataDisplay from '@/components/DataDisplay.vue'

import { useForm } from 'vee-validate'
import * as yup from 'yup'
import FormField from '@/components/FormField.vue'
import Message from 'primevue/message'

const schema = yup.object().shape({
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

const [externalURI] = defineField('externalURI')
const [title] = defineField('title')
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
const [inHathiTrust] = defineField('inHathiTrust')
const [collectionID] = defineField('collectionID')
const [collectionFacet] = defineField('collectionFacet')
const [isCollection] = defineField('isCollection')

const route = useRoute()
const router = useRouter()
const metadataStore = useMetadataStore()
const systemStore = useSystemStore()
const userStore = useUserStore()

const originalUseRight = ref(1)

// this indicates that extURI, barcode or catkey have changed and need to be validated
const needsValidation = ref(false) 

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
   let out = [{label: "None", value: "none"}]
   systemStore.collectionFacets.forEach( o => {
      out.push({label: o.name, value: o.name})
   })
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
   if (typeof values.ocrHint === 'undefined') return true
   if ( values.ocrHint == 0) return true
   let hint = systemStore.ocrHints.find( h => h.id == values.ocrHint)
   return !hint.ocrCandidate
})
const availabilityPolicies = computed(() => {
   let out = []
   systemStore.availabilityPolicies.forEach( o => {
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

const submittedToHathiTrust = computed( () => {
   if ( metadataStore.detail.inHathiTrust == false) return false
   return ( metadataStore.hathiTrustStatus.metadataSublittedAt != null)
})

onMounted( async () =>{
   let mdID = route.params.id
   await metadataStore.getDetails(mdID)
   document.title = `Edit | Metadata ${mdID}`
   needsValidation.value = false

   // NOTE: catalogKey and barcode may be null as they are optional data members
   // but the lookup code does not handle null values. make sure they are empty string instead
   let md = {
      type: metadataStore.detail.type,
      externalURI: metadataStore.detail.externalURI,
      title: metadataStore.detail.title,
      callNumber: metadataStore.detail.callNumber,
      author: metadataStore.detail.creatorName,
      catalogKey: "",
      barcode: "",
      personalItem: metadataStore.detail.isPersonalItem,
      manuscript: metadataStore.detail.isManuscript,
      ocrHint: 0,
      ocrLanguageHint: metadataStore.detail.ocrLanguageHint,
      inDPLA: metadataStore.detail.inDPLA,
      inHathiTrust: metadataStore.detail.inHathiTrust,
      isCollection: metadataStore.detail.isCollection,
      collectionID: metadataStore.detail.collectionID,
      collectionFacet: metadataStore.detail.collectionFacet,
   }

   if (metadataStore.detail.type == "ExternalMetadata") {
      md.externalSystemID = 1
   }

   if (  metadataStore.detail.catalogKey ) {
      md.catalogKey = metadataStore.detail.catalogKey
   }
   if (  metadataStore.detail.barcode ) {
      md.barcode = metadataStore.detail.barcode
   }
   if (metadataStore.detail.ocrHint) {
      md.ocrHint = metadataStore.detail.ocrHint.id
   }
   if (metadataStore.detail.availabilityPolicy) {
      md.availabilityPolicy = metadataStore.detail.availabilityPolicy.id
   }

   // set use right based on current metadata settings and preserve original setting
   if (metadataStore.detail.useRightName) {
      let ur = systemStore.useRights.find( r => r.name == metadataStore.detail.useRightName)
      md.useRight = ur.id
   }
   originalUseRight.value = md.useRight
   resetForm({ values: md })
   
})

const sirsiLookup = ( async () => {
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

const validateASMetadata = ( async () => {
   await metadataStore.validateArchivesSpaceURI(values.externalURI)
   if (metadataStore.asMatch.error == "") {
      needsValidation.value = false
      setValues({
         externalURI: metadataStore.asMatch.validatedURL,
         title: metadataStore.asMatch.title
      })
   }
})

const cancelEdit = (() => {
   metadataStore.resetArchivesSpaceErrors()
   router.push(`/metadata/${route.params.id}`)
})

const submitChanges = handleSubmit( async (values) => {
   if ( needsValidation.value ) {
      return
   }

   if (metadataStore.detail.type == 'SirsiMetadata') {
      // SEE IF UR changed from CNE / UND to sotething else, or if UR chanegd from something valid to CNE/UND
      // in these cases, send the new ID. Otehrwise send a 0 so backend ignores the request.
      let origCNE = (originalUseRight.value == 1 || values.useRight == 11)
      let updatedCNE = ( values.useRight == 1 ||  values.useRight == 11)
      if (origCNE && updatedCNE ) {
         // no change from CNE... don't send
         values.useRight = 0
      } else {
         if ( originalUseRight.value == values.useRight) {
            // no change. do not send
            values.useRight = 0
         }
      }
   }
   await metadataStore.submitEdit( values )
   if (systemStore.showError == false) {
      router.push(`/metadata/${metadataStore.detail.id}`)
   }
})
</script>

<style lang="scss" scoped>
.edit-form {
   width: 60%;
   margin: 30px auto 0 auto;
   text-align: left;

   :deep(.p-panel-content) {
      display: flex;
      flex-direction: column;
      gap: 15px;   
      dl {
         margin: 0 !important;
      } 
   }

   p.note {
      margin: 0;
      padding: 10px;
      border: 1px solid var(--uvalib-teal-light);
      background: var(--uvalib-teal-lightest);
      border-radius: 3px;
      text-align: left;
   }

   .split {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: flex-end;
      gap: 10px;
   }
}
</style>
