<template>
   <h2>Orders</h2>
   <div class="orders">
      <DataTable :value="ordersStore.orders" ref="ordersTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll"
         sortField="id" :sortOrder="1" @sort="onSort($event)"
         :lazy="true" :paginator="true" @page="onPage($event)"
         :rows="ordersStore.searchOpts.limit" :totalRecords="ordersStore.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[10,30,100]"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <Column field="id" header="ID" :sortable="true" />
         <Column field="status" header="Status" >
            <template #body="slotProps">
               <span :class="`status ${slotProps.data.status}`">{{displayStatus(slotProps.data.status)}}</span>
            </template>
         </Column>
         <Column field="dateSubmitted" header="Request Submitted" :sortable="true" class="nowrap" />
         <Column field="dateDue" header="Date Due" :sortable="true" class="nowrap" />
         <Column field="title" header="Title" :sortable="true" />
         <Column field="unitCount" header="Units" :sortable="true" />
         <Column field="fee" header="Fee" :sortable="true" />
         <Column field="lastName" header="Customer" class="nowrap" >
            <template #body="slotProps">
               <router-link :to="`/customers/${slotProps.data.customer.id}`">
                  {{slotProps.data.customer.lastName}}, {{slotProps.data.customer.firstName}}
               </router-link>
            </template>
         </Column>
         <Column field="agency.name" header="Agency" />

         <Column header="" class="row-acts">
            <template #body="slotProps">
               <router-link :to="`/orders/${slotProps.data.id}`">View</router-link>
            </template>
         </Column>
      </DataTable>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useOrdersStore } from '@/stores/orders'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const ordersStore = useOrdersStore()

function displayStatus( id) {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
}
function onPage(event) {
   ordersStore.searchOpts.start = event.first
   ordersStore.searchOpts.limit = event.rows
   ordersStore.getJobs()
}

function onSort(event) {
   ordersStore.searchOpts.sortField = event.sortField
   ordersStore.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      ordersStore.searchOpts.sortOrder = "desc"
   }
   ordersStore.getOrders( )
}

onMounted(() => {
   ordersStore.getOrders()
})
</script>

<style scoped lang="scss">
   .orders {
      min-height: 600px;
      text-align: left;
      padding: 0 25px;

      .p-datatable {
         font-size: 0.85em;
         span.status {
            border-radius: 5px;
            color: white;
            padding: 2px 15px 4px 15px;
            box-sizing: border-box;
            width: 100%;
            display:inline-block;
            text-align: center;
            font-weight: bold;
         }
         span.status.await_fee {
            background: var(--uvalib-grey-dark);
         }
         span.status.canceled {
            background: var(--uvalib-red-darker);
         }
         span.status.deferred {
            background: var(--uvalib-blue-alt-dark);
         }
         span.status.approved {
            background: var(--uvalib-green-dark);
         }
         span.status.requested {
            background: var( --uvalib-brand-orange);
         }
         :deep(td.nowrap) {
            white-space: nowrap;
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