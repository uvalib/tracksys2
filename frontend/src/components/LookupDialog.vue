<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :title="dialogTitle">
      <UButton @click="show()" size="sm" color="secondary" :label="props.label" :disabled="props.disabled"/>
      <template #body>
         <div class="row-left gap-2" v-if="mode=='lookup'">
            <UInput v-model="query"  @keydown.stop.prevent.enter="lookupRecords" :placeholder="searchPlaceholder" class="grow"/>
            <UButton @click="lookupRecords" label="Lookup" color="secondary"/>
            <UButton @click="createMetadata" label="Create" color="secondary" v-if="props.create"/>
         </div>
         <NewMetadataPanel v-else @canceled="metadataCreateCanceled" @created="metadataCreated" />
         <div v-if="searched && mode=='lookup'" class="pt-4 max-h-100 overflow-y-scroll text-xs">
            <UTable :data="lookupResults" :columns="columns" :ui="{th: 'text-xs p-2', td: 'text-xs p-2'}"
               :rowSelectionOptions="{enableMultiRowSelection: false}" @select="hitSelected" 
            />
         </div>
      </template>
      <template #footer  v-if="mode=='lookup'">
         <UButton @click="hide" label="Cancel" color="secondary"/>
         <UButton @click="okClicked" label="Select" :disabled="!selectedHit" />
      </template>
   </UModal>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useMetadataStore } from '@/stores/metadata'
import { useOrdersStore } from '@/stores/orders'
import { useComponentsStore } from '@/stores/components'
import NewMetadataPanel from '@/components/metadata/NewMetadataPanel.vue'

const emit = defineEmits( ['selected', 'selectedObject' ])
const props = defineProps({
   target: {
      type: String,
      required: true
   },
   create: {
      type: Boolean,
      default: false
   },
   label: {
      type: String,
      default: "Lookup"
   },
   disabled: {
      type: Boolean,
      default: false
   },
})

const metadataStore = useMetadataStore()
const ordersStore = useOrdersStore()
const componentsStore = useComponentsStore()

const isOpen = ref(false)
const mode = ref("lookup")
const searched = ref(false)
const query = ref("")
const selectedHit = ref()

const lookupResults = computed(() => {
   if ( props.target == "metadata") return metadataStore.searchHits
   if ( props.target == "component") return componentsStore.searchHits
   if ( props.target == "orders") return ordersStore.lookupHits
   return []
})

const columns = computed(() => {
   let cols = [ ]
   if ( props.target == "metadata") {
      cols.push( {accessorKey: "id", header: "ID"} )
      cols.push( {accessorKey: "pid", header: "PID"} )
      cols.push( {accessorKey: "type", header: "Type"} )
      cols.push( {accessorKey: "title", header: "Title", cell: ({ row }) => `${truncateTitle(row.original.title)}`} )
      cols.push( {accessorKey: "callNumber", header: "Call Number"} )
      cols.push( {accessorKey: "barcode", header: "Barcode"} )
   } else if ( props.target == "component") {
      cols.push( {accessorKey: "id", header: "ID"} )
      cols.push( {accessorKey: "pid", header: "PID"} )
      cols.push( {accessorKey: "title", header: "Title", cell: ({ row }) => `${truncateTitle(row.original.title)}`} )
      cols.push( {accessorKey: "label", header: "Label", cell: ({ row }) => `${truncateTitle(row.original.label)}`} )
      cols.push( {accessorKey: "date", header: "Date"} )
      cols.push( {accessorKey: "eadID", header: "EAD ID"} )
      cols.push( {accessorKey: "masterFileCount", header: "Master Files"} )
   } else if ( props.target == "orders") {
      cols.push( {accessorKey: "id", header: "ID"} )
      cols.push( {accessorKey: "dateDue", header: "Date Due"} )
      cols.push( {accessorKey: "title", header: "Title", cell: ({ row }) => `${truncateTitle(row.original.title)}`} )
      cols.push( {accessorKey: "staff_notes", header: "Staff Notes", cell: ({ row }) => `${truncateTitle(row.original.staff_notes)}`} )
      cols.push( {accessorKey: "special_instructions", header: "Special Instructions", cell: ({ row }) => `${truncateTitle(row.original.special_instructions)}`} )
      cols.push( {accessorKey: "customer", header: "Customer"} )
      cols.push( {accessorKey: "agency", header: "Agency"} )
   }
   return cols
})

const hitSelected = (( _e, row) => {
   selectedHit.value = row.original
   row.toggleSelected(!row.getIsSelected())
})

const createMetadata = (() =>{
   mode.value = "create"
   searched.value = false
   query.value = ""
   selectedHit.value = null
})
const metadataCreateCanceled = (() =>{
   mode.value = "lookup"
})
const metadataCreated = (() =>{
   selectedHit.value = metadataStore.searchHits[0]
   emit("selected", selectedHit.value.id)
   emit("selectedObject", selectedHit.value)
   hide()
})

const searchPlaceholder = computed(() => {
   if (props.target == "metadata") {
      return "Lookup metadata..."
   }
   if (props.target == "component") {
      return "Lookup component..."
   }
   return "Lookup order..."
})

const dialogTitle = computed(() => {
   if (props.target == "metadata") {
      if (mode.value == "create") return "Create Metadata"
      return "Metadata Lookup"
   }
   if (props.target == "component") {
      return "Component Lookup"
   }
   return "Order Lookup"
})

const lookupRecords = ( async () => {
   selectedHit.value = null
   searched.value = false
   if (props.target == "metadata") {
      await metadataStore.lookup( query.value )
   } else if  (props.target == "component") {
      await componentsStore.lookup( query.value )
   } else {
      await ordersStore.lookup( query.value )
   }
   searched.value = true
})

const truncateTitle = ((t) => {
   if (t.length < 75) return t
   return t.slice(0,75)+"..."
})

const okClicked = (() => {
   emit("selected", selectedHit.value.id)
   emit("selectedObject", selectedHit.value)
   hide()
})

const hide = (() => {
   isOpen.value=false
})

const show = (() => {
   isOpen.value = true
   selectedHit.value = null
   searched.value = false
   query.value = ""
   mode.value = "lookup"
})
</script>

<style lang="scss" scoped>
</style>
