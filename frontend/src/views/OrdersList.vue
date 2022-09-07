<template>
   <h2>Orders</h2>
   <div class="orders">
      <div class="toolbar">
         <span>
            <label for="orders-filter">Filter:</label>
            <Dropdown id="orders-filter" v-model="ordersStore.searchOpts.filter" @change="getOrders()"
               :options="filters" optionLabel="name" optionValue="code" />
         </span>
         <span>
            <span class="p-input-icon-right">
               <i class="pi pi-search" />
               <InputText v-model="ordersStore.searchOpts.query" placeholder="Orders Search" @input="queryOrders()"/>
            </span>
            <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()"/>
         </span>
      </div>
      <DataTable :value="ordersStore.orders" ref="ordersTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll"
         :sortField="ordersStore.searchOpts.sortField" :sortOrder="sortOrder" @sort="onSort($event)"
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
         <Column field="masterFileCount" header="Master Files" :sortable="true" />
         <Column field="fee" header="Fee" :sortable="true" />
         <Column field="lastName" header="Customer" class="nowrap" >
            <template #body="slotProps">
               {{slotProps.data.customer.lastName}}, {{slotProps.data.customer.firstName}}
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
import { onBeforeMount, onMounted, ref, computed } from 'vue'
import { useOrdersStore } from '@/stores/orders'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const ordersStore = useOrdersStore()

const filters = ref([
   {name: "Active", code: "active"},
   {name: "Await Approval", code: "await"},
   {name: "Deferred", code: "deferred"},
   {name: "Complete", code: "complete"},
   {name: "Canceled", code: "canceled"},
   {name: "Due Today", code: "due_today"},
   {name: "Due within a Week", code: "due_week"},
   {name: "Overdue", code: "overdue"}
])

const sortOrder = computed(() => {
   if (ordersStore.searchOpts.sortOrder == "desc") {
      return -1
   }
   return 1
})

onBeforeMount( () => {
   if ( route.query.q ) {
      ordersStore.searchOpts.query = route.query.q
   }
   if ( route.query.filter ) {
      ordersStore.searchOpts.filter = route.query.filter
   }
   if ( route.query.sort  ) {
      let bits = route.query.sort.split(" ")
      ordersStore.searchOpts.sortField = bits[0].trim()
      ordersStore.searchOpts.sortOrder = bits[1].trim()
   }
})



function displayStatus( id) {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
}

function clearSearch() {
   ordersStore.searchOpts.query = ""
   getOrders()
}

function queryOrders() {
   if (ordersStore.searchOpts.query.length > 3) {
      getOrders()
   }
}

function setQueryParams() {
   let query = Object.assign({}, route.query)
   delete query.q
   if (ordersStore.searchOpts.query) {
      query.q = ordersStore.searchOpts.query
   }
   query.filter = ordersStore.searchOpts.filter
   query.sort = `${ordersStore.searchOpts.sortField} ${ordersStore.searchOpts.sortOrder}`
   router.push({query})
}

function getOrders() {
   setQueryParams()
   ordersStore.getOrders()
}

function onPage(event) {
   ordersStore.searchOpts.start = event.first
   ordersStore.searchOpts.limit = event.rows
   ordersStore.getOrders()
}

function onSort(event) {
   ordersStore.searchOpts.sortField = event.sortField
   ordersStore.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      ordersStore.searchOpts.sortOrder = "desc"
   }
   getOrders( )
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

      .toolbar {
         padding: 10px 0;
         display: flex;
         flex-flow: row nowrap;
         justify-content: space-between;
         label {
            font-weight: bold;
            margin-right: 5px;
            display: inline-block;
         }
         button.p-button {
            margin-left: 5px;
         }
      }

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
         span.status.completed {
            background: var( --uvalib-brand-blue-light);
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