<template>
   <h2>Collections</h2>
   <div class="row-between bg-white border-b border-brand-grey-100 py-3 px-2 items-center">
      <template v-if="userStore.isAdmin">
         <UModal v-model:open="showCreateCollection" :modal="true" :dismissible="false" :close="false" title="Create Collection">
            <UButton color="secondary" size="sm" @click="showCreateCollection = true" label="Create Collection" />
            <template #body>
               <NewMetadataPanel @canceled="showCreateCollection = false" @created="collectionCreated" :collection="true"/>
            </template>
         </UModal>
      </template>
      <span v-else></span>
      <span>{{ collectionStore.collections.length }} Collections</span>
   </div>
   <UTable ref="table" :data="collectionStore.collections" :columns="columns" sticky virtualize class="h-[calc(100vh-185px)]">
      <template #id-cell="{ row }">
         <router-link :to="`/metadata/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
   </UTable>
</template>

<script setup>
import { onMounted, ref} from 'vue'
import { useCollectionsStore } from '../stores/collections'
import { useUserStore } from '../stores/user'
import { useSystemStore } from '../stores/system'
import { useMetadataStore } from '../stores/metadata'
import NewMetadataPanel from '@/components/metadata/NewMetadataPanel.vue'

const collectionStore = useCollectionsStore()
const userStore = useUserStore()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()

const table = ref()
const showCreateCollection = ref(false)

const columns = [
   { accessorKey: 'id', header: 'ID'},
   { accessorKey: 'pid', header: '{ID'},
   { accessorKey: 'title', header: 'Title'},
   { accessorKey: 'creatorName', header: 'Creator Name'},
   { accessorKey: 'barcode', header: 'Barcode'},
   { accessorKey: 'callNumber', header: 'Call NumberD'},
   { accessorKey: 'catalogKey', header: 'Catalog Key'},
   { accessorKey: 'collectionID', header: 'ID'},
   { accessorKey: 'recordCount', header: 'Items'}
]

onMounted(() => {
   collectionStore.getCollections()
   document.title = `Collections`
})

const collectionCreated = (() => {
   collectionStore.addCollection(  metadataStore.detail )
   systemStore.toastMessage("Collection Created", `A new collection metadata record has been created.`)
   showCreateCollection.value = false
})
</script>

<style scoped lang="scss">
.collections {
   margin: 20px;
}
</style>