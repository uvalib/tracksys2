<template>
   <h2>{{ title }}</h2>
   <div class="published">
      <DataTable :value="publishedStore.records" ref="publishedTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll"
         :lazy="true" :paginator="true" @page="onPage($event)"  paginatorPosition="top"
         :rows="publishedStore.searchOpts.limit" :totalRecords="publishedStore.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[25,50,100]" :first="publishedStore.searchOpts.start"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
         v-model:filters="columnFilters" filterDisplay="menu" @filter="onFilter($event)"
      >
         <Column field="id" header="ID">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
            </template>
         </Column>
         <Column field="pid" header="PID" class="nowrap" />
         <Column v-if="route.params.type != 'archivesspace'" field="type" header="Type" filterField="type" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <Dropdown v-model="filterModel.value" :options="mdTypes" optionLabel="name" optionValue="code" placeholder="Select a type" />
            </template>
            <template #body="slotProps">
               <div v-if="slotProps.data.type != 'ExternalMetadata'">{{slotProps.data.type}}</div>
               <div v-else>ArchivesSpace</div>
            </template>
         </Column>
         <Column field="title" header="Title" filterField="title" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
            </template>
         </Column>
         <Column field="creatorName" header="Creator Name" filterField="creator_name" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Creator name"/>
            </template>
            <template #body="slotProps">
               <span v-if="slotProps.data.creatorName">{{ slotProps.data.creatorName }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
         <Column  v-if="route.params.type == 'virgo'" field="barcode" header="Barcode" class="nowrap" filterField="barcode" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Barcode"/>
            </template>
            <template #body="slotProps">
               <span v-if="slotProps.data.barcode">{{ slotProps.data.barcode }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
         <Column v-if="route.params.type == 'virgo'" field="callNumber" header="Call Number" class="nowrap" filterField="call_number" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Call number"/>
            </template>
            <template #body="slotProps">
               <span v-if="slotProps.data.callNumber">{{ slotProps.data.callNumber }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
         <Column v-if="route.params.type == 'virgo'" field="catalogKey" header="Catalog Key" class="nowrap" filterField="catalog_key" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Catalog key"/>
            </template>
            <template #body="slotProps">
               <span v-if="slotProps.data.catalogKey">{{ slotProps.data.catalogKey }}</span>
               <span v-else class="none">N/A</span>
            </template>
         </Column>
      </DataTable>
   </div>
</template>

<script setup>
import { onBeforeMount, onMounted, ref, computed } from 'vue'
import { usePublishedStore } from '@/stores/published'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import { useRoute, useRouter, onBeforeRouteUpdate } from 'vue-router'
import { FilterMatchMode } from 'primevue/api'
import { usePinnable } from '@/composables/pin'

usePinnable("p-paginator-top")

const route = useRoute()
const router = useRouter()
const publishedStore = usePublishedStore()

const columnFilters = ref( {
   'type': {value: null, matchMode: FilterMatchMode.EQUALS},
   'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'creator_name': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'barcode': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
   'call_number': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
   'catalog_key': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
})

const mdTypes = computed(() => {
   let out = []
   out.push({name: "Sirsi", code: "SirsiMetadata"})
   out.push({name: "XML", code: "XmlMetadata"})
   return out
})

const title = computed(() => {
   if (route.params.type == "archivesspace") {
      return "Items Published to ArchivesSpace"
   }
   if (route.params.type == "dpla") {
      return "Items Published to DPLA"
   }
   return "Items Published to Virgo"
})

onBeforeRouteUpdate( (to) => {
   publishedStore.searchOpts.filters = []
   publishedStore.getRecords(to.params.type)
})

onBeforeMount( () => {
   publishedStore.searchOpts.filters = []
   if ( route.query.filters ) {
      let filters = JSON.parse(route.query.filters)
      filters.forEach( filter => {
         let bits = filter.split("|")
         publishedStore.searchOpts.filters.push( {field: bits[0], match: bits[1], value: bits[2]} )
      })
   }
})

onMounted(() => {
   publishedStore.getRecords(route.params.type)
   document.title = `Published`
})

const onFilter = ((event) => {
   publishedStore.searchOpts.filters = []
   Object.entries(event.filters).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         publishedStore.searchOpts.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })

   let query = Object.assign({}, route.query)
   if ( publishedStore.searchOpts.filters.length > 0 ) {
      query.filters = publishedStore.filtersAsQueryParam
   } else {
      delete query.filters
   }
   router.push({query})
   publishedStore.getRecords(route.params.type)
})

const onPage = ((event) => {
   publishedStore.searchOpts.start = event.first
   publishedStore.searchOpts.limit = event.rows
   publishedStore.getRecords(route.params.type)
})

</script>

<style scoped lang="scss">
.published {
   min-height: 600px;
   text-align: left;
   padding: 0 25px;

   .filters {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      align-items: center;
   }
   .left-pad {
      margin-left: 10px;
   }
   .right-pad {
      margin-right: 10px;
   }

   .p-datatable {
      font-size: 0.85em;
      span.status {
         width: 100%;
      }
      .dimmed {
         display:inline-block;
         color: #ccc;
      }
      span.dimmed {
         margin-left: 3px;
      }
      :deep(td), :deep(th) {
         padding: 10px;
      }
      :deep(.row-acts) {
         text-align: center;
         padding: 0;
         a {
            display: inline-block;
            margin: 0;
            padding: 5px 10px;
         };
      }
   }
}
</style>