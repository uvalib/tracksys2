<template>
   <div class="similar">
      <h3>Similar Images</h3>
      <div class="hits">
         <div v-if="searchStore.similarImages.total > 0" class="summary">
            {{searchStore.similarImages.total}} matche(s) found    
         </div>
         <UTable :data="searchStore.similarImages.hits" :columns="columns">
            <template #empty><h4>No matching images found</h4></template>
            <template #id-cell="{ row }">
                <router-link :to="`/masterfiles/${row.original.id}`">{{row.original.id}}</router-link>
            </template>
            <template #metadataPID-cell="{ row }">
               <router-link :to="`/metadata/${row.original.metadataID}`">{{row.original.metadataPID}}: {{ row.original.metadataTitle }}</router-link>
            </template>
            <template #unitID-cell="{ row }">
               <router-link :to="`/units/${row.original.unitID}`">{{row.original.unitID}}</router-link>
            </template>
            <template #thumbnailURL-cell="{ row }">
               <a :href="row.original.imageURL" target="_blank">
                     <img :src="row.original.thumbnailURL" />
               </a>
            </template>
         </UTable>
      </div>
   </div>
</template>

<script setup>
import { useSearchStore } from '../../stores/search'

const searchStore = useSearchStore()

const columns = [
   {
      accessorKey: 'id',
      header: 'ID'
   },
   {
      accessorKey: 'pid',
      header: 'PID'
   },
   {
      accessorKey: 'metadataPID',
      header: 'Metadata'
   },
   {
      accessorKey: 'unitID',
      header: 'Unit ID'
   },
   {
      accessorKey: 'filename',
      header: 'Filename'
   },
   {
      accessorKey: 'title',
      header: 'Title'
   },
   {
      accessorKey: 'description',
      header: 'Description'
   },
   {
      accessorKey: 'thumbnailURL',
      header: 'Thumb'
   },
]
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