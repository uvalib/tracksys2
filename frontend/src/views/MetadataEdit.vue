<template>
   <h2>{{pageHeader}}</h2>
   <UForm v-if="systemStore.working == false" ref="form" :schema="schema" :state="state" @submit="submitChanges" class="flex flex-col gap-6 w-full max-w-250 mx-auto! mt-6!"">
      <UCard title="General Information">
         <div class="column gap-4">
            <p v-if="metadataStore.detail.type == 'XmlMetadata'" class="bg-brand-teal-200 border border-brand-teal p-2 text-sm rounded-sm">
               <b>Note</b>:
               To update the XML, use the 'Download XML' button on the details page to get a copy of the data.
               Edit it with a standalone XML editor, then upload the result using the'Upload XML' button.
            </p>
            <template  v-if="metadataStore.detail.type == 'SirsiMetadata'" >
               <div class="flex flex-col gap-1">
                  <div class="row-left gap-2">
                     <UFormField name="catalogKey" label="Catalog Key" :required="true" class="grow">
                        <UInput v-model="state.catalogKey" class="w-full"/>  
                     </UFormField>
                     <UFormField name="barcode" label="Barcode" :required="true" class="grow">
                        <UInput v-model="state.barcode" class="w-full"/>  
                     </UFormField>
                     <div class="flex flex-row flex-nowrap gap-2 pt-5.25 items-start">
                        <UButton @click="sirsiLookup()" size="sm" label="Lookup" color="secondary" 
                           :disabled="metadataStore.sirsiMatch.searching" :loading="metadataStore.sirsiMatch.searching"
                        />
                        <UButton @click="clearSirsiInfo" size="sm" label="Clear" color="secondary"/>
                     </div>
                  </div>
                  <div class="validation" v-if="metadataStore.sirsiMatch.error">{{ metadataStore.sirsiMatch.error }}</div>
               </div>
               <div>
                  <dl class="text-sm">
                     <DataDisplay label="Title" :value="state.title" blankValue="Unknown"/>
                     <DataDisplay label="Call Number" :value="state.callNumber" blankValue="Unknown"/>
                  </dl>
               </div>
            </template>
            <template v-if="metadataStore.detail.type == 'ExternalMetadata'">
               <div class="flex flex-col gap-1">
                  <p class="bg-brand-teal-200 border border-brand-teal p-2 text-sm rounded-sm">
                     <span class="font-bold!">IMPORTANT</span>: 
                     Only URIs containing /resources/ or /archival_objects/ are supported.
                  </p>
               </div>
               <div class="row-left gap-2">
                  <UFormField name="externalURI" label="External URI" :required="true" class="grow">   
                     <div class="flex flex-row gap-4">
                        <UInput v-model="state.externalURI"  @update:modelValue="asNeedsValidation" class="w-full"/>  
                        <UButton @click="validateASMetadata()" label="Validate" color="secondary" :loading="metadataStore.asMatch.searching"/>
                     </div>  
                  </UFormField>
               </div>
               <div class="validation" v-if="needsValidation">Changes to external URI need to be validated</div>
               <div class="validation" v-if="metadataStore.asMatch.error">{{ metadataStore.asMatch.error }}</div>
               <dl>
                  <DataDisplay label="Title" :value="state.title" blankValue="Unknown"/>
               </dl>
            </template> 
            <div class="row-left gap-2">
               <UFormField v-if="userStore.isAdmin" name="isCollection" label="Collection" class="grow">   
                  <USelect v-model="state.isCollection"  :items="yesNo" class="w-full"/>   
               </UFormField>
               <UFormField name="personalItem" label="Personal Item" class="grow">   
                  <USelect v-model="state.personalItem" :items="yesNo"  class="w-full"/>   
               </UFormField>
               <UFormField nam="state" label="Manuscript" class="grow">   
                  <USelect v-model="state.manuscript" :items="yesNo" class="w-full"/>   
               </UFormField>
            </div>
            <div class="row-left gap-2">
               <UFormField name="ocrHint" label="OCR Hint" class="grow">   
                     <USelect v-model="state.ocrHint" :items="ocrHints" class="w-full" placeholder="Select a hint"/>   
               </UFormField>
               <UFormField name="ocrLanguageHint" label="OCR Language" class="grow">   
                     <USelectMenu v-model="state.ocrLanguageHint" :items="ocrLanguages" class="w-full" 
                     :disabled="isLanguageDisabled" placeholder="Select a language"/>   
               </UFormField>
            </div>
         </div>
      </UCard>
      <UCard v-if="metadataStore.detail.type != 'ExternalMetadata'" title="Digital Library Information">
         <div class="column gap-4">
            <div class="row-left gap-2">
               <UFormField name="collectionID" label="Collection ID" class="grow">   
                  <UInput v-model="state.collectionID" class="w-full" />   
               </UFormField>
               <UFormField name="collectionFacet" label="Collection Facet" class="grow">   
                  <USelect v-model="state.collectionFacet" :items="collectionFacets" placeholder="Select a facet" class="w-full" />   
               </UFormField>
            </div>
            <div class="row-left gap-2">
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
               <div class="text-sm p-3">{{ rightStatement }}</div>
            </div>    
         </div>
      </UCard>
      <div class="row-right gap-2">
         <UButton label="Cancel" color="secondary" @click="cancelEdit()"/>
         <UButton label="Save" type="submit" /> 
      </div>
   </UForm>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import { useMetadataStore } from '@/stores/metadata'
import { useSystemStore } from '@/stores/system'
import { useUserStore } from '@/stores/user'
import * as yup from 'yup'

const schema = yup.object({
   title: yup.string().required('Title is required'),
   callNumber:  yup.string().when('type', {
      is: (value) => value == 'SirsiMetadata',
      then: (schema) => schema.required("Call number is required"),
   }),
   catalogKey:  yup.string().when('type', {
      is: (value) => value == 'SirsiMetadata',
      then: (schema) => schema.required("CatKey is required"),
   }),
   barcode:  yup.string().when('type', {
      is: (value) => value == 'SirsiMetadata',
      then: (schema) => schema.required("Barcode is required"),
   }),
   availabilityPolicy:  yup.number().when('type', {
      is: (value) => value == 'SirsiMetadata',
      then: (schema) => schema.min(1).required("Availability is required"),
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

const state = ref({
   externalURI: "",
   title: "",
   callNumber: "",
   catalogKey: "",
   barcode: "",
   personalItem: false,
   manuscript: false,
   ocrHint: null,
   ocrLanguageHint: null,
   availabilityPolicy: 1,
   useRight: 1,
   inDPLA: false,
   inHathiTrust: false,
   collectionID: "",
   collectionFacet: "",
   isCollection: false
})

const route = useRoute()
const router = useRouter()
const metadataStore = useMetadataStore()
const systemStore = useSystemStore()
const userStore = useUserStore()

const form = ref()
const originalUseRight = ref(1)
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
   if (typeof state.value.ocrHint === 'undefined') return true
   if ( state.value.ocrHint == null) return true
   let hint = systemStore.ocrHints.find( h => h.id == state.value.ocrHint)
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

const clearSirsiInfo = (() => {
   state.value.title = ""
   state.value.callNumber = ""
   state.value.author = ""
   state.value.catalogKey = ""
   state.value.barcode = ""
   metadataStore.sirsiMatch.error = ""
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
   // but the lookup code does not handle null state.value. make sure they are empty string instead
   state.value = {
      type: metadataStore.detail.type,
      externalURI: metadataStore.detail.externalURI,
      title: metadataStore.detail.title,
      callNumber: metadataStore.detail.callNumber,
      author: metadataStore.detail.creatorName,
      catalogKey: "",
      barcode: "",
      personalItem: metadataStore.detail.isPersonalItem,
      manuscript: metadataStore.detail.isManuscript,
      ocrHint: null,
      ocrLanguageHint: metadataStore.detail.ocrLanguageHint,
      inDPLA: metadataStore.detail.inDPLA,
      inHathiTrust: metadataStore.detail.inHathiTrust,
      isCollection: metadataStore.detail.isCollection,
      collectionID: metadataStore.detail.collectionID,
      collectionFacet: metadataStore.detail.collectionFacet,
   }

   if (metadataStore.detail.type == "ExternalMetadata") {
      state.value.externalSystemID = 1
   }

   if (  metadataStore.detail.catalogKey ) {
      state.value.catalogKey = metadataStore.detail.catalogKey
   }
   if (  metadataStore.detail.barcode ) {
      state.value.barcode = metadataStore.detail.barcode
   }
   if (metadataStore.detail.ocrHint) {
      state.value.ocrHint = metadataStore.detail.ocrHint.id
   }
   if (metadataStore.detail.availabilityPolicy) {
      state.value.availabilityPolicy = metadataStore.detail.availabilityPolicy.id
   }

   // set use right based on current metadata settings and preserve original setting
   if (metadataStore.detail.useRightName) {
      let ur = systemStore.useRights.find( r => r.name == metadataStore.detail.useRightName)
      state.value.useRight = ur.id
   }
   originalUseRight.value = state.value.useRight   
})

const sirsiLookup = ( async () => {
   await metadataStore.sirsiLookup(state.value.barcode, state.value.catalogKey)
   if ( metadataStore.sirsiMatch.error == "") {
      needsValidation.value = false
      state.value.title = metadataStore.sirsiMatch.title
      state.value.callNumber = metadataStore.sirsiMatch.callNumber
      state.value.author = metadataStore.sirsiMatch.creatorName
      state.value.catalogKey = metadataStore.sirsiMatch.catalogKey
      state.value.barcode = metadataStore.sirsiMatch.barcode
      form.value.validate()
   }
})

const validateASMetadata = ( async () => {
   await metadataStore.validateArchivesSpaceURI(state.value.externalURI)
   if (metadataStore.asMatch.error == "") {
      needsValidation.value = false
      setValues({
         externalURI: metadataStore.asMatch.validatedURL,
         title: metadataStore.asMatch.title
      })
   }
})

const asNeedsValidation = () => {
   needsValidation.value = true 
   state.value.title = ""
}

const cancelEdit = (() => {
   metadataStore.resetArchivesSpaceErrors()
   router.push(`/metadata/${route.params.id}`)
})

const submitChanges = ( async () => {
   if ( needsValidation.value ) {
      return
   }

   if (metadataStore.detail.type == 'SirsiMetadata') {
      // SEE IF UR changed from CNE / UND to sotething else, or if UR chanegd from something valid to CNE/UND
      // in these cases, send the new ID. Otehrwise send a 0 so backend ignores the request.
      let origCNE = (originalUseRight.value == 1 || state.value.useRight == 11)
      let updatedCNE = ( state.value.useRight == 1 ||  state.value.useRight == 11)
      if (origCNE && updatedCNE ) {
         // no change from CNE... don't send
         state.value.useRight = 0
      } else {
         if ( originalUseRight.value == state.value.useRight) {
            // no change. do not send
            state.value.useRight = 0
         }
      }
   }
   await metadataStore.submitEdit( state.value )
   if (systemStore.showError == false) {
      router.push(`/metadata/${metadataStore.detail.id}`)
   }
})
</script>

<style lang="scss" scoped>
.validation {
   color: var(--uvalib-red-dark);
}
</style>
