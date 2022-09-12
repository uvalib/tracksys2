<template>
   <div v-if="props.orders.length == 0">
      <h3>No related orders found</h3>
   </div>
   <DataTable v-else :value="props.orders" ref="relatedOrdersTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="false" :paginator="false"
   >
      <Column field="id" header="ID"/>
      <Column header="Customer" class="nowrap">
         <template #body="slotProps">
            {{slotProps.data.customer.lastName}}, {{slotProps.data.customer.firstName}}
         </template>
      </Column>
      <Column field="agency.name" header="Agency" class="nowrap" />
      <Column field="title" header="Order Title" />
      <Column field="notes" header="Staff Notes" />
      <Column field="specialInstructions" header="Special Instructions" />
      <Column header="" class="row-acts nowrap">
         <template #body="slotProps">
            <router-link :to="`/orders/${slotProps.data.id}`">View details</router-link>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const props = defineProps({
   orders: {
      type: Array,
      required: true
   }
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
}
</stype>