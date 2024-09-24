<template>
   <div class="similar">
      <h3>Similar Images</h3>
      <div class="hits">
         <DataTable :value="searchStore.similarImages.hits" ref="similarHitsTable" dataKey="id"
            stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm" :rowStyle="rowStyle"
            :totalRecords="searchStore.similarImages.total"
         >
            <template #empty><h4>No matching images found</h4></template>
            <template #header v-if="searchStore.similarImages.total > 0">
               <div class="results-toolbar">
                  <div class="matches" v-if="searchStore.similarImages.total > 50">{{searchStore.similarImages.total}} matches found, only showing the top 50</div>
                  <div class="matches" v-else>{{searchStore.similarImages.total}} matches found</div>
               </div>
            </template>
            <Column field="id" header="ID">
               <template #body="slotProps">
                  <router-link :to="`/masterfiles/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
               </template>
            </Column>
            <Column field="pid" header="PID" class="nowrap" />
            <Column field="metadataPID" header="Metadata">
               <template #body="slotProps">
                  <router-link :to="`/metadata/${slotProps.data.metadataID}`">{{slotProps.data.metadataPID}}: {{ slotProps.data.metadataTitle }}</router-link>
               </template>
            </Column>
            <Column field="unitID" header="Unit ID" class="nowrap">
               <template #body="slotProps">
                  <router-link :to="`/units/${slotProps.data.unitID}`">{{slotProps.data.unitID}}</router-link>
               </template>
            </Column>
            <Column field="filename" header="Filename" class="nowrap" />
            <Column field="title" header="Title" />
            <Column field="description" header="Description" class="nowrap" />
            <Column field="thumbnailURL" header="Thumb">
               <template #body="slotProps">
                  <a :href="slotProps.data.imageURL" target="_blank">
                     <img :src="slotProps.data.thumbnailURL" />
                  </a>
               </template>
            </Column>
         </DataTable>
      </div>
   </div>
</template>

<script setup>
import { useSearchStore } from '../../stores/search'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const searchStore = useSearchStore()

const rowStyle = (data) => {
    if (data.originalID) {
        return { background: '#f5f5f5' };
    }
}
</script>

<style scoped lang="scss">
.similar {
   h3 {
      text-align: left;
      font-weight: 600;
      border-bottom: 1px solid #dee2e6;
      padding: 25px 10px 10px 10px;
      margin-bottom:0;
   }
   h4 {
      text-align: center;
      font-size: 1.1em;
   }
   .hits {
      padding: 1rem;
   }
   .matches {
      margin: 0px 0 20px 0;
      text-align: left;
   }
}
</style>