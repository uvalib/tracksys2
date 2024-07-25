<template>
   <DataTable :value="props.units" ref="relatedUnitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="false" :paginator="true" :rows="15" :rowsPerPageOptions="[15,30,50]" removableSort
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
      v-model:filters="filters" filterDisplay="menu" paginatorPosition="top"
   >
      <template #empty><h3>No units found</h3></template>
      <template #paginatorstart>
         <div class="unit-acts">
            <HathiTrustUpdateDialog v-if="props.hathiTrust" :orderID="props.orderID" />
            <AddUnitDialog v-if="props.canAdd"/>
            <DPGButton label="Download Units CSV" severity="secondary" @click="downloadCSV" />
         </div>
      </template>
      <Column field="id" header="ID" :sortable="true">
         <template #body="slotProps">
            <router-link :to="`/units/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="metadata.title" header="Title" :sortable="true" v-if="showMetadata">
         <template #body="slotProps">
            <router-link v-if="slotProps.data.metadata" :to="`/metadata/${slotProps.data.metadataID}`">{{slotProps.data.metadata.title}}</router-link>
            <span v-else class="empty">N/A</span>
         </template>
      </Column>
      <Column field="status" header="Status" filterField="status" :showFilterMatchModes="false">
         <template #filter="{filterModel}">
            <Select v-model="filterModel.value" :options="statusValues" optionLabel="label" optionValue="value" placeholder="Select a status" />
         </template>
      </Column>
      <Column field="metadata.callNumber" header="Call Number" v-if="showMetadata" class="nowrap" />
      <Column field="intendedUse.name" header="Intended Use" filterField="intendedUse.id" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Select v-model="filterModel.value" :options="systemStore.intendedUses" optionLabel="name" optionValue="id" placeholder="Select an intended use" />
         </template>
      </Column>
      <Column field="datePatronDeliverablesReady"  header="Date Patron Deliverables Ready" :sortable="true">
         <template #body="slotProps">
            <span v-if="slotProps.data.datePatronDeliverablesReady">
               {{$formatDate(slotProps.data.datePatronDeliverablesReady)}}
            </span>
            <span v-else class="empty">N/A</span>
         </template>
      </Column>
      <Column field="dateDLDeliverablesReady" header="Date DL Deliverables Ready" :sortable="true">
         <template #body="slotProps">
            <span v-if="slotProps.data.dateDLDeliverablesReady">
               {{$formatDate(slotProps.data.dateDLDeliverablesReady)}}
            </span>
            <span v-else class="empty">N/A</span>
         </template>
      </Column>
      <Column field="reorder" header="Reorder?" filterField="reorder" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Select v-model="filterModel.value" :options="yesNo" optionLabel="label" optionValue="value" placeholder="Select a reorder status" />
         </template>
         <template #body="slotProps">
            {{$formatBool(slotProps.data.reorder)}}
         </template>
      </Column>
      <Column field="masterFilesCount" header="Master Files Count" :sortable="true"/>
   </DataTable>
</template>

<script setup>
import { ref, computed } from 'vue'
import { FilterMatchMode } from '@primevue/core/api'
import Select from 'primevue/select'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import AddUnitDialog from '@/components/order/AddUnitDialog.vue'
import HathiTrustUpdateDialog from '@/components/HathiTrustUpdateDialog.vue'
import { useSystemStore } from '@/stores/system'
import { usePinnable } from '@/composables/pin'

usePinnable("p-datatable-paginator-top")

const props = defineProps({
   orderID: {
      type: Number,
      default: -1
   },
   units: {
      type: Array,
      required: true
   },
   canAdd: {
      type: Boolean,
      default: false
   },
   showMetadata: {
      type: Boolean,
      default: true
   },
   hathiTrust: {
      type: Boolean,
      default: false
   },
})

const systemStore = useSystemStore()
const relatedUnitsTable = ref()
const filters = ref( {
   'status': {value: null, matchMode: FilterMatchMode.EQUALS},
   'intendedUse.id': {value: null, matchMode: FilterMatchMode.EQUALS},
   'reorder': {value: null, matchMode: FilterMatchMode.EQUALS},
})

const statusValues = computed(() => {
   let out = []
   out.push( {label: "Approved", value: 'approved'} )
   out.push( {label: "Canceled", value: 'canceled'} )
   out.push( {label: "Done", value: 'done'} )
   out.push( {label: "Error", value: 'error'} )
   out.push( {label: "Unapproved", value: 'unapproved'} )
   return out
})

const yesNo = computed(() => {
   let out = []
   out.push( {label: "No", value: false} )
   out.push( {label: "Yes", value: true} )
   return out
})

const downloadCSV = (() => {
   relatedUnitsTable.value.exportCSV()
})

</script>

<stype scoped lang="scss">
   .unit-acts{
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: center;
      gap: 10px;
   }
</stype>