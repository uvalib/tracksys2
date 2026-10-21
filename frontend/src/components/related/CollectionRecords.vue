<template>
   <div class="row-between p-2 sticky z-5 bg-white border-b border-brand-grey-100" :style="{top: system.headerHeight}">
      <div class="row-left gap-2">
         <CollectionBulkAdd v-if="userStore.isAdmin" />
         <UButton label="Export CSV" size="sm" color="secondary" @click="exportCollectionCSV" :disabled="collectionStore.totalRecords == 0"/>
         <UButton label="Export Collection" size="sm" color="secondary" @click="exportCollection" :disabled="collectionStore.totalRecords == 0" v-if="userStore.isAdmin"/>
         <UButton label="Master Files Report" size="sm" color="secondary" @click="exportCollectionMasterfiles" :disabled="collectionStore.totalRecords == 0" v-if="userStore.isAdmin"/>
      </div>
      <div class="row-left items-center gap-2">
         <UPagination color="neutral" variant="ghost"
            v-model:page="collectionStore.searchOpts.currPage" :items-per-page="collectionStore.searchOpts.limit" 
            :total="collectionStore.totalRecords" @update:page="pageChanged"
         />
         <USelect v-model="collectionStore.searchOpts.limit" :items="[15,30,100]" @change="perPageChanged" />
         <span>Total: {{ collectionStore.totalRecords }}</span>
      </div>
      <UInput v-model="collectionStore.searchOpts.query" placeholder="Search collection..." @update:modelValue="queryCollection"/>
   </div>
   <UTable :data="collectionStore.records" :columns="columns" >
      <template #id-cell="{ row }">
         <router-link :to="`/metadata/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
      <template #type-cell="{ row }">
         <span v-if="row.original.type!='ExternalMetadata'">{{ row.original.type }}</span>
         <span v-else>{{externalSystemName(row.original.externalSystemID)}}</span>
      </template>
      <template #callNumber-cell="{ row }">
         <span v-if="row.original.callNumber">{{ row.original.callNumber }}</span>
         <span v-else>N/A</span>
      </template>
      <template #barcode-cell="{ row }">
         <span v-if="row.original.barcode">{{ row.original.barcode }}</span>
         <span v-else>N/A</span>
      </template>
   </UTable>
</template>

<script setup>
import { useCollectionsStore } from '@/stores/collections'
import { useConfirm } from "@/composables/useConfirm"
import { useUserStore } from '@/stores/user'
import { useSystemStore } from '@/stores/system'
import CollectionBulkAdd from './CollectionBulkAdd.vue'

const userStore = useUserStore()
const collectionStore = useCollectionsStore()
const system = useSystemStore()

const columns = [
   { accessorKey: "id", header: "ID" },
   { accessorKey: "pid", header: "PID" },
   { accessorKey: "type", header: "Type" },
   { accessorKey: "title", header: "Title" },
   { accessorKey: "callNumber", header: "Call Number" },
   { accessorKey: "barcode", header: "Barcode" },
]

const props = defineProps({
   collectionID: {
      type: Number,
      required: true
   }
})

const externalSystemName = ( (id) => {
   const sys = system.externalSystems.find( s => s.id == id)
   if ( sys ) {
      return sys.name
   }
   return "Unknown"
})

const onSort = ((event) => {
   collectionStore.searchOpts.sortField = event.sortField
   collectionStore.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      collectionStore.searchOpts.sortOrder = "desc"
   }
   collectionStore.getItems()
})

const perPageChanged = (() => {
   collectionStore.searchOpts.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   collectionStore.searchOpts.start  = (collectionStore.searchOpts.currPage-1) * collectionStore.searchOpts.limit
   collectionStore.getItems()
})

const queryCollection = (() => {
   collectionStore.getItems( false )
})
const deleteItem = ( async ( item ) => {
   const msg = `Remove "${item.pid} : ${item.title}" from this collection?`
   const resp = await useConfirm("Confirm Remove Item", msg, "Remove")
   if (resp) {
       collectionStore.removeItem(item)
   } 
})
const exportCollectionCSV = (() => {
   collectionStore.exportCSV()
})

const exportCollectionMasterfiles = (async () => {
   const msg = 'Generate a master files report for this collection?'
   const resp = await useConfirm("Master Files Report", msg, "Export")
   if (resp) {
      collectionStore.masterFilesReport()
   } 
})

const exportCollection = ( async () => {
   const msg = 'Export all metadata and image files for this collection?'
   const resp = await useConfirm("Export Collection", msg, "Export")
   if (resp) {
      collectionStore.exportAll()
   } 
})

</script>

<style scoped lang="scss">
.collection  {
   margin: 0;
   h3  {
      text-align: center;
   }
   .buttons {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: flex-start;
      gap: 10px;
   }
   td.apt-status {
      text-align: center;
      span.pi {
         font-size: 1.2em;
         font-weight: bold;
      }
      span.pi.success {
         color: var(--uvalib-green);
      }
      span.pi.fail {
         color: var(--uvalib-red-dark);
      }
   }
   td.row-acts {
      text-align: center;
      width: 25px;
   }
}
</style>