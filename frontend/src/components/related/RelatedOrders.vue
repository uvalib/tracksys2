<template>
   <DataTable :value="props.orders" ref="relatedOrdersTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="false" :paginator="true" :rows="15" :rowsPerPageOptions="[15,30,50]" removableSort
      paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="{first} - {last} of {totalRecords}" paginatorPosition="top"
      v-model:filters="filters" filterDisplay="menu"
   >
      <template #empty><h3>No related orders found</h3></template>
      <template #paginatorstart></template>
      <Column field="id" header="ID" :sortable="true">
         <template #body="slotProps">
            <router-link :to="`/orders/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
         </template>
      </Column>
      <Column  field="customer.lastName" header="Customer" class="nowrap" :sortable="true" filterField="customer.lastName" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Last name"/>
         </template>
         <template #body="slotProps">
            {{slotProps.data.customer.lastName}}, {{slotProps.data.customer.firstName}}
         </template>
      </Column>
      <Column field="agency.name" header="Agency" class="nowrap" :sortable="true" filterField="agency.id" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <Dropdown v-model="filterModel.value" :options="systemStore.agencies" optionLabel="name" optionValue="id" placeholder="Select agency" />
         </template>
      </Column>
      <Column field="title" header="Order Title" filterField="title" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
         </template>
      </Column>
      <Column field="notes" header="Staff Notes" filterField="staffNotes" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Staff notes"/>
         </template>
      </Column>
      <Column field="specialInstructions" header="Special Instructions" filterField="specialInstructions" :showFilterMatchModes="false" >
         <template #filter="{filterModel}">
            <InputText type="text" v-model="filterModel.value" placeholder="Special instructions"/>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { ref } from 'vue'
import { FilterMatchMode } from '@primevue/core/api'
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import { useSystemStore } from '@/stores/system'
import { usePinnable } from '@/composables/pin'

usePinnable("p-datatable-paginator-top")

const systemStore = useSystemStore()

const props = defineProps({
   orders: {
      type: Array,
      required: true
   }
})

const filters = ref( {
   'agency.id': {value: null, matchMode: FilterMatchMode.EQUALS},
   'customer.lastName': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
   'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'specialInstructions': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'staffNotes': {value: null, matchMode: FilterMatchMode.CONTAINS},
})
</script>

<stype scoped lang="scss">
// .results {
//    margin: 20px;
//    font-size: 0.9em;
//    h3 {
//       text-align: center;
//    }
//    td.nowrap, th {
//       white-space: nowrap;
//    }
//    th, td {
//       font-size: 0.85em;
//    }
// }
</stype>