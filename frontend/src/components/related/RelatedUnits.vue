<template>
   <DataTable :value="props.units" ref="relatedUnitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="false" :paginator="true" :rows="15" :rowsPerPageOptions="[15,30,50]" removableSort
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
      v-model:filters="filters" filterDisplay="menu" paginatorPosition="top"
   >
      <template #paginatorstart>
         <HathiTrustUpdateDialog v-if="props.hathiTrust" />
         <AddUnitDialog v-if="props.orderStatus != 'completed' && props.orderStatus != 'canceled'"/>
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
      <Column field="metadata.callNumber" header="Call Number" v-if="showMetadata" class="nowrap" />
      <Column field="intendedUse.name" header="Intended Use" filterField="intendedUse.id" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="systemStore.intendedUses" optionLabel="name" optionValue="id" placeholder="Select an intended use" />
         </template>
      </Column>
      <Column field="datePatronDeliverablesReady"  header="Date Patron Deliverables Ready" :sortable="true">
         <template #body="slotProps">
            <span v-if="slotProps.data.datePatronDeliverablesReady">
               {{formatDate(slotProps.data.datePatronDeliverablesReady)}}
            </span>
            <span v-else class="empty">N/A</span>
         </template>
      </Column>
      <Column field="dateDLDeliverablesReady" header="Date DL Deliverables Ready" :sortable="true">
         <template #body="slotProps">
            <span v-if="slotProps.data.dateDLDeliverablesReady">
               {{formatDate(slotProps.data.dateDLDeliverablesReady)}}
            </span>
            <span v-else class="empty">N/A</span>
         </template>
      </Column>
      <Column field="reorder" header="Reorder?" filterField="reorder" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="yesNo" optionLabel="label" optionValue="value" placeholder="Select a reorder status" />
         </template>
         <template #body="slotProps">
            {{formatBoolean(slotProps.data.reorder)}}
         </template>
      </Column>
      <Column field="masterFilesCount" header="Master Files Count" :sortable="true"/>
   </DataTable>
</template>

<script setup>
import { ref,computed } from 'vue'
import { FilterMatchMode } from 'primevue/api'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import dayjs from 'dayjs'
import AddUnitDialog from '@/components/order/AddUnitDialog.vue'
import HathiTrustUpdateDialog from '@/components/order/HathiTrustUpdateDialog.vue'
import { usePinnable } from '@/composables/pin'

usePinnable("p-paginator-top")

const props = defineProps({
   units: {
      type: Array,
      required: true
   },
   orderStatus: {
      type: String,
      default: "requested"
   },
   showMetadata: {
      type: Boolean,
      default: true
   },
   hathiTrust: {
      type: Boolean,
      default: false
   }
})


const yesNo = computed(() => {
   let out = []
   out.push( {label: "No", value: false} )
   out.push( {label: "Yes", value: true} )
   return out
})

const filters = ref( {
   'intendedUse.id': {value: null, matchMode: FilterMatchMode.EQUALS},
   'reorder': {value: null, matchMode: FilterMatchMode.EQUALS},
})

const formatBoolean = ( (flag) => {
   if (flag) return "Yes"
   return "No"
})

const formatDate = (  (date ) => {
   if (date) {
      return dayjs(date).format("YYYY-MM-DD")
   }
   return ""
})

</script>

<stype scoped lang="scss">
td.nowrap {
   white-space: nowrap;
}

.empty {
   color: #ccc;
}
</stype>