<template>
   <h2>Collections</h2>
   <div class="collections">
      <DataTable :value="collectionStore.collections" ref="collectionRecordsTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="false" :paginator="collectionStore.totalCollections > 15"
         :rows="pageSize"  :rowsPerPageOptions="[15,30,100]" :totalRecords="collectionStore.totalRecords"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <Column field="id" header="ID">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
            </template>
         </Column>
         <Column field="pid" header="PID" class="nowrap"/>
         <Column field="title" header="Title" />
         <Column field="creatorName" header="Creator Name" />
         <Column field="barcode" header="Barcode" />
         <Column field="callNumber" header="Call Number" />
         <Column field="catalogKey" header="Catalog Key" />
         <Column field="recordCount" header="Items" />
      </DataTable>
   </div>
</template>

<script setup>
import { onMounted, ref} from 'vue'
import { useCollectionsStore } from '@/stores/collections'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const collectionStore = useCollectionsStore()

const pageSize = ref(15)

onMounted(() => {
   collectionStore.getCollections()
   document.title = `Collections`
})
</script>

<style scoped lang="scss">
.collections {
   min-height: 600px;
   text-align: left;
   padding: 25px 40px;
}
</style>