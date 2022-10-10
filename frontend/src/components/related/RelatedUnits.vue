<template>
  <div v-if="props.units.length == 0">
      <h3>No units found</h3>
   </div>
   <DataTable v-else :value="props.units" ref="relatedUnitsTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="false" :paginator="props.units.length > 15" :rows="15" :rowsPerPageOptions="[15,30,50]" removableSort
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
   >
      <Column field="id" header="ID" :sortable="true"/>
      <Column field="intendedUse.name" header="Intended Use"/>
      <Column header="Date Patron Deliverables Ready">
         <template #body="slotProps">
            <span v-if="slotProps.data.datePatronDeliverablesReady">
               {{formatDate(slotProps.data.datePatronDeliverablesReady)}}
            </span>
            <span v-else class="empty">N/A</span>
         </template>
      </Column>
      <Column header="Date DL Deliverables Ready">
         <template #body="slotProps">
            <span v-if="slotProps.data.dateDLDeliverablesReady">
               {{formatDate(slotProps.data.dateDLDeliverablesReady)}}
            </span>
            <span v-else class="empty">N/A</span>
         </template>
      </Column>
      <Column field="reorder" header="Reorder?">
         <template #body="slotProps">
            {{formatBoolean(slotProps.data.reorder)}}
         </template>
      </Column>
      <Column field="masterFilesCount" header="Master Files Count" :sortable="true"/>
      <Column header="" class="row-acts nowrap">
         <template #body="slotProps">
            <router-link :to="`/units/${slotProps.data.id}`">View details</router-link>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import dayjs from 'dayjs'

const props = defineProps({
   units: {
      type: Array,
      required: true
   }
})

function formatBoolean( flag) {
   if (flag) return "Yes"
   return "No"
}

function formatDate( date ) {
   if (date) {
      return dayjs(date).format("YYYY-MM-DD")
   }
   return ""
}

</script>

<stype scoped lang="scss">
.results {
   margin: 20px;
   font-size: 0.9em;
   td.nowrap, th {
      white-space: nowrap;
   }
   th, td {
      font-size: 0.85em;
   }
}
.empty {
   color: #ccc;
}
</stype>