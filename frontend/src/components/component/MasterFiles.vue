<template>
   <Panel header="Related Master Files"  v-if="component.loadingMasterFiles">
      <p>Loading master files...</p>
   </Panel>
   <template v-else>
      <Panel header="Related Master Files" v-if="component.relatedMasterFiles.length > 0">
         <DataTable :value="component.relatedMasterFiles" ref="componentMasterFilesTable" dataKey="id"
            showGridlines stripedRows responsiveLayout="scroll" class="p-datatable-sm"
            :lazy="false" :paginator="component.relatedMasterFiles.length > 15" :rows="15" :rowsPerPageOptions="[15,30,50,100]"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
         >
            <Column field="id" header="ID">
               <template #body="slotProps">
                  <router-link :to="`/masterfiles/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
               </template>
            </Column>
            <Column field="pid" header="PID"/>
            <Column field="filename" header="File Name"/>
            <Column field="title" header="Title"/>
            <Column field="description" header="Description"/>
            <Column field="metadata.tiitle" header="Metadata">
               <template #body="slotProps">
                  <router-link :to="`/metadata/${slotProps.data.metadata.id}`">{{slotProps.data.metadata.title}}</router-link>
               </template>
            </Column>
            <Column field="thumbnailURL" header="Thumb" class="thumb">
               <template #body="slotProps">
                  <a :href="slotProps.data.viewerURL" target="_blank">
                     <img :src="slotProps.data.thumbnailURL" :class="{exemplar: slotProps.data.exemplar}"/>
                  </a>
               </template>
            </Column>
         </DataTable>
      </Panel>
      <Panel header="Related Master Files" v-else>
         <p>No master files are associated with this component.</p>
      </Panel>
   </template>
</template>

<script setup>
import { useComponentsStore } from '@/stores/components'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const component = useComponentsStore()
</script>

<style scoped lang="scss">
</style>