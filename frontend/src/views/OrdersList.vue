<template>
   <h2>Orders</h2>
   <DPGButton label="New Order" class="create" @click="createOrder()"/>
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
         <Column field="id" header="ID" :sortable="true">
            <template #body="slotProps">
               <router-link :to="`/orders/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
            </template>
         </Column>
         <Column field="status" header="Status" >
            <template #body="slotProps">
               <span :class="`status ${slotProps.data.status}`">{{displayStatus(slotProps.data.status)}}</span>
            </template>
         </Column>
         <Column field="dateSubmitted" header="Request Submitted" :sortable="true" class="nowrap" />
         <Column field="dateDue" header="Date Due" :sortable="true" class="nowrap" />
         <Column field="title" header="Title" :sortable="true" />
         <Column field="specialInstructions" header="Special Instructions" />
         <Column field="unitCount" header="Units" :sortable="true" />
         <Column field="masterFileCount" header="Master Files" :sortable="true" />
         <Column field="fee" header="Fee" :sortable="true">
            <template #body="slotProps">
               <span class="fee" v-if="slotProps.data.fee > 0">${{parseFloat(slotProps.data.fee).toFixed(2)}}</span>
            </template>
         </Column>
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
   {name: "Due in a Week", code: "due_week"},
   {name: "Overdue", code: "overdue"},
   {name: "Ready for Delivery", code: "ready"}
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

function createOrder() {
   router.push("/orders/new")
}

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
   document.title = `Orders`
})
</script>

<style scoped lang="scss">
:deep(td.nowrap) {
   white-space: nowrap;
}
button.p-button.create {
   position: absolute;
   right:15px;
   top: 15px;
   font-size: 0.9em;
}
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
            width: 100%;
         }
         span.sta
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