<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :title="title">
      <UButton @click="show" :label="props.label" color="secondary"/>
      <template  #body>
         <UForm v-if="mode=='unit'" :schema="schema"  :state="state" class="flex flex-col gap-4" @submit="createUnit" ref="form">
            <UCard title="Unit Metadata">
               <UFormField name="metadataID" :required="true">
                  <div class="row-left gap-2!">
                     <UInput v-model="metadataSearch"  @keydown.stop.prevent.enter="lookupMetadata" class="flex-1"  />
                     <UButton @click="lookupMetadata" label="Lookup" color="secondary" :loading="metadataStore.working" size="sm"/>
                     <UButton @click="createMetadata" label="Create" color="secondary" size="sm"/>
                  </div>
               </UFormField>
               <div class="scroller">
                  <UTable :data="metadataStore.searchHits" :columns="columns" :ui="{ td: 'text-xs!' }" v-model:sorting="sorting"
                     v-model:row-selection="rowSelection" :rowSelectionOptions="{enableMultiRowSelection: false}" @select="metadataSelected" 
                     
                  >
                     <template #pid-header="{ column }">
                        <UButton size="xs" color="neutral" variant="ghost" label="PID" :icon="sortIcon(column)" @click="column.toggleSorting(column.getIsSorted() === 'asc')"/>
                     </template>
                     <template #title-header="{ column }">
                        <UButton size="xs" color="neutral" variant="ghost" label="Title" :icon="sortIcon(column)" @click="column.toggleSorting(column.getIsSorted() === 'asc')"/>
                     </template>
                     <template #call_number-header="{ column }">
                        <UButton size="xs" color="neutral" variant="ghost" :icon="sortIcon(column)" @click="column.toggleSorting(column.getIsSorted() === 'asc')">
                           <span class="text-start">Call Number</span>
                        </UButton>
                     </template>
                     <template #barcode-header="{ column }">
                        <UButton size="xs" color="neutral" variant="ghost" label="Barcode" :icon="sortIcon(column)" @click="column.toggleSorting(column.getIsSorted() === 'asc')"/>
                     </template>
                  </UTable>
               </div>
            </UCard>
            <UCard title="Digitization Information">
               <div class="column-close">
                  <UFormField name="intendedUseID" label="Intended Use" :required="true">
                     <USelect v-model="state.intendedUseID"  :items="intendedUses" placeholder="Select an intended use" class="w-full" />   
                  </UFormField>
                  <UFormField name="sourceURL" label="Source URL">
                     <UInput v-model="state.sourceURL"  class="w-full" />   
                  </UFormField>
                  <UFormField name="specialInstructions" label="Special Instructions">
                     <UTextarea v-model="state.specialInstructions"  class="w-full" />   
                  </UFormField>
                  <UFormField name="staffNotes" label="Staff Notes">
                     <UTextarea v-model="state.staffNotes"  class="w-full" />   
                  </UFormField>
                  <label class="row-left gap-2! items-center">
                     <UCheckbox size="lg" v-model="state.completeScan" />
                     <span>Complete Scan</span>
                  </label>
                  <label class="row-left gap-2! items-center">
                     <UCheckbox size="lg" v-model="state.throwAway" />
                     <span>Throw Away</span>
                  </label>
                  <label class="row-left gap-2! items-center">
                     <UCheckbox size="lg" v-model="state.includeInDL" />
                     <span>Include in Virgo</span>
                  </label>
               </div>
            </UCard>
            <div class="row-right gap_2">
               <UButton @click="isOpen=false" label="Cancel" color="secondary"/>
               <UButton label="Add Unit" type="submit" />
            </div>
         </UForm>
         <NewMetadataPanel v-else @canceled="metadataCreateCanceled" @created="metadataCreated" />
      </template>
   </UModal>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useSystemStore } from '@/stores/system'
import { useMetadataStore } from '@/stores/metadata'
import { useOrdersStore } from '@/stores/orders'
import NewMetadataPanel from '@/components/metadata/NewMetadataPanel.vue'
import * as yup from 'yup'

const props = defineProps({
   label: {
      type: String,
      default: "Add Unit",
   },
   item: {
      type: Object,
      default: null,
   },
})

const columns = [
   {
      accessorKey: 'pid',
      header: 'PID'
   },
   {
      accessorKey: 'title',
      header: 'Title'
   },
   {
      accessorKey: 'call_number',
      header: 'Call Number'
   },
   {
      accessorKey: 'barcode',
      header: 'Barcode'
   },
]
const sorting = ref([])


const systemStore = useSystemStore()
const metadataStore = useMetadataStore()
const ordersStore = useOrdersStore()

const isOpen = ref(false)
const mode = ref("unit")
const searched = ref(false)
const metadataSearch = ref("")
const selectedMetadata = ref()
const form = ref()

const schema = yup.object({
   metadataID: yup.number().min(1, "Metadata ID is required"),
   intendedUseID: yup.number().min(1,'Intended use is required'),
})

const state = ref({
   itemID: null,
   metadataID: null,
   intendedUseID: null,
   sourceURL: "",
   specialInstructions: "",
   staffNotes: "",
   completeScan: false,
   throwAway: false,
   includeInDL: false
})

const rowSelection = ref({})
const hideCheckbox = ref({select: false})

const title = computed(() => {
   if (mode.value == "unit") return "Add Unit"
   return "Create Metadata"
})

const intendedUses = computed(() => {
   let out = []
   systemStore.intendedUses.forEach( a => {
      if (a.name == "Digital Collection Building") {
         out.push( {label: `${a.name}: Highest Possible resolution TIFF`, value: a.id} )
      } else if (a.deliverableFormat == 'pdf') {
         out.push( {label: `${a.name}: PDF`, value: a.id} )
      } else {
         let dpi = "DPI"
         if (a.deliverableResolution == "Highest Possible") {
            dpi = "resolution"
         }
         out.push( {label: `${a.name}: ${a.deliverableResolution} ${dpi} ${a.deliverableFormat.toUpperCase()}`, value: a.id} )
      }
   })
   return out
})

const sortIcon = ((column) => {
   const isSorted = column.getIsSorted()
   if (isSorted) {
      if (isSorted === 'asc') return 'i-lucide-arrow-up-narrow-wide'
      return 'i-lucide-arrow-down-wide-narrow'
   }
   return 'i-lucide-arrow-up-down'
})

const metadataSelected = (( _e, row) => {
   const selID = row.original.id 
   if (state.value.metadataID == selID) {
      state.value.metadataID = null   
   } else {
      state.value.metadataID = selID       
   }
   row.toggleSelected(!row.getIsSelected())
   form.value.validate({name: 'metadataID'})
})

const createMetadata = (() => {
   mode.value = "metadata"
})

const metadataCreateCanceled = (() => {
   mode.value = "unit"
})

const metadataCreated = (() => {
   console.log(metadataStore.searchHits[0])
   metadataSearch.value = metadataStore.detail.pid
   searched.value = true
   selectedMetadata.value = metadataStore.searchHits[0]
   setValues({metadataID: metadataStore.searchHits[0].id})
   mode.value = "unit"
})

const lookupMetadata = ( async () => {
   await metadataStore.lookup( metadataSearch.value )
   searched.value = true
})

const createUnit = (async () => {
   await ordersStore.addUnit(state.value)
   if (systemStore.error == "") {
      isOpen.value = false
   }
})

const show = (() => {
   mode.value = "unit"
   isOpen.value = true
   metadataSearch.value = ""
   selectedMetadata.value = null
   metadataStore.resetSearch()

   state.value = {
      itemID: 0,
      metadataID: 0,
      intendedUseID: null,
      sourceURL: "",
      specialInstructions: "",
      staffNotes: "",
      completeScan: false,
      throwAway: false,
      includeInDL: false
   }

   if (props.item) {
      metadataSearch.value = props.item.title
      var si = `Title: ${props.item.title}`
      si += `\nPages to Digitize: ${props.item.pages}`
      if (props.item.callNumber && props.item.callNumber != "") {
         si += `\nCall Number: ${props.item.callNumber}`
         metadataSearch.value = props.item.callNumber
      }
      if (props.item.author && props.item.author != "") {
         si += `\nAuthor: ${props.item.author}`
      }
      if (props.item.year && props.item.year != "") {
         si += `\nYear: ${props.item.year}`
      }
      if (props.item.location && props.item.location != "") {
         si += `\nLocation: ${props.item.location}`
      }
      if (props.item.description && props.item.description != "") {
         si += `\nDescription: ${props.item.description}`
      }
      state.itemID = props.item.id
      state.specialInstructions = si
      state.sourceURL = props.item.sourceURL
      state.intendedUseID = props.item.intendedUse.id
      lookupMetadata()
   }
})

</script>

<style lang="scss" scoped>
   div.no-results {
      margin: 15px 0 0 0;
      font-size: 1.2em;
      font-weight: 600;
      text-align: center;
   }
   .scroller {
      max-height: 250px;
      overflow-y: scroll;
      font-size: 0.8em;
   }
</style>
