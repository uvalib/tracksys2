<template>
   <DataTable :value="props.masterFiles" ref="relatedMasterFilesTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="false" :paginator="props.masterFiles.length > 15" :rows="15" :rowsPerPageOptions="[15,30,50]" removableSort
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="{first} - {last} of {totalRecords}"
      v-model:filters="filters" filterDisplay="menu"
   >
      <Column field="id" header="ID" :sortable="true">
         <template #body="slotProps">
            <router-link :to="`/masterfiles/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column field="pid" header="PID" class="nowrap" :sortable="true"/>
      <Column field="title" header="Title" class="nowrap" :sortable="true" filterField="title" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
         </template>
      </Column>
      <Column field="description" header="Description" :sortable="true" filterField="description" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Description"/>
         </template>
      </Column>
      <Column field="thumbnailURL" header="Thumb" class="thumb">
         <template #body="slotProps">
            <a :href="slotProps.data.viewerURL" target="_blank">
               <img :src="slotProps.data.thumbnailURL"/>
            </a>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { ref } from 'vue'
import InputText from 'primevue/inputtext'
import { FilterMatchMode } from 'primevue/api'

const props = defineProps({
   masterFiles: {
      type: Array,
      required: true
   }
})

const filters = ref( {
   'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'description': {value: null, matchMode: FilterMatchMode.CONTAINS},
})
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
   :deep(td.thumb) {
      width: 160px !important;
      text-align: center !important;
   }
}
</stype>