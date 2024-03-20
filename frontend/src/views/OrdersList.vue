<template>
   <h2>Orders</h2>
   <div class="orders">
      <DataTable :value="ordersStore.orders" ref="ordersTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll"
         :sortField="ordersStore.searchOpts.sortField" :sortOrder="sortOrder" @sort="onSort($event)"
         :lazy="true" :paginator="true" @page="onPage($event)"  paginatorPosition="top"
         :rows="ordersStore.searchOpts.limit" :totalRecords="ordersStore.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[30,50,100]" :first="ordersStore.searchOpts.start"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
         v-model:filters="columnFilters" filterDisplay="menu" @filter="getOrders()"
      >
         <template #paginatorstart>
            <DPGButton v-if="(userStore.isAdmin || userStore.isSupervisor)" class="p-button-secondary" label="Create Order" @click="createOrder()"/>
         </template>
         <template #paginatorend>
            <div class="filters">
               <label for="orders-filter">Filter:</label>
               <Dropdown id="orders-filter" v-model="statusFilter" @change="getOrders()"
                  :options="filters" optionLabel="name" optionValue="code" />
               <ToggleButton v-model="assignedToMe" class="left-pad right-pad" onIcon="" offIcon="" onLabel="Assigned to Me" offLabel="Assigned to Me" @change="ownerToggled()" />
               <IconField iconPosition="left">
                  <InputIcon class="pi pi-search" />
                  <InputText v-model="ordersStore.searchOpts.query" placeholder="Search Orders" @input="queryOrders()"/>
               </IconField>
            </div>
         </template>
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
               <span class="fee-waived" v-if="slotProps.data.feeWaived">Waived</span>
               <span class="fee" v-else-if="slotProps.data.fee !== undefined">${{parseFloat(slotProps.data.fee).toFixed(2)}}</span>
            </template>
         </Column>
         <Column field="lastName" header="Customer" filterField="customer" :showFilterMatchModes="false">
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Last name"/>
            </template>
            <template #body="slotProps">
               <div class="nowrap">{{slotProps.data.customer.lastName}}, {{slotProps.data.customer.firstName}}</div>
               <div class="dimmed" v-if="slotProps.data.customer.academicStatus">({{slotProps.data.customer.academicStatus.name}})</div>
            </template>
         </Column>
         <Column field="agency.name" header="Agency" filterField="agency" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <Dropdown v-model="filterModel.value" :options="systemStore.agencies" optionLabel="name" optionValue="id" placeholder="Select agency" />
            </template>
         </Column>
         <Column field="processor" header="Processor" class="nowrap" filterField="processor" :showFilterMatchModes="false">
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Last name"/>
            </template>
            <template #body="slotProps">
               <span v-if="slotProps.data.processor">{{slotProps.data.processor.lastName}}, {{slotProps.data.processor.firstName}}</span>
               <span v-else class="dimmed">None</span>
            </template>
         </Column>
      </DataTable>
   </div>
</template>

<script setup>
import { onBeforeMount, onMounted, ref, computed } from 'vue'
import { useOrdersStore } from '@/stores/orders'
import { useUserStore } from '@/stores/user'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dropdown from 'primevue/dropdown'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import { useRoute, useRouter } from 'vue-router'
import ToggleButton from 'primevue/togglebutton'
import { FilterMatchMode } from 'primevue/api'
import { useSystemStore } from '@/stores/system'
import { usePinnable } from '@/composables/pin'

usePinnable("p-paginator-top")

const systemStore = useSystemStore()
const route = useRoute()
const router = useRouter()
const ordersStore = useOrdersStore()
const userStore = useUserStore()

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

const columnFilters = ref( {
   'customer': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
   'processor': {value: null, matchMode: FilterMatchMode.STARTS_WITH},
   'agency': {value: null, matchMode: FilterMatchMode.EQUALS},
})
const statusFilter = ref("active")
const assignedToMe = ref(false)

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
   ordersStore.searchOpts.filters = []
   if ( route.query.filters ) {
      let filters = JSON.parse(route.query.filters)
      filters.forEach( filter => {
         let bits = filter.split("|")
         ordersStore.searchOpts.filters.push( {field: bits[0], match: bits[1], value: bits[2]} )
         if ( bits[0] == "status") {
            statusFilter.value = bits[2]
         } else if ( bits[0] == "customer") {
            columnFilters.value.customer.value = bits[2]
         } else if ( bits[0] == "processor") {
            columnFilters.value.processor.value = bits[2]
         } else if ( bits[0] == "agency") {
            columnFilters.value.agency.value = bits[2]
         }
      })
   }
   if ( route.query.sort  ) {
      let bits = route.query.sort.split(" ")
      ordersStore.searchOpts.sortField = bits[0].trim()
      ordersStore.searchOpts.sortOrder = bits[1].trim()
   }
})

onMounted(() => {
   ordersStore.getOrders()
   document.title = `Orders`
})

const createOrder = (() => {
   router.push("/orders/new")
})

const displayStatus = ( (id) => {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
})

const queryOrders = (() => {
   getOrders()
})

const ownerToggled = (() => {
   if ( assignedToMe.value == true) {
      ordersStore.setTargetOwner( userStore.ID )
   } else {
      ordersStore.clearTargetOwner()
   }
   getOrders()
})

const setQueryParams = (() => {
   let query = Object.assign({}, route.query)
   delete query.q
   if (ordersStore.searchOpts.query) {
      query.q = ordersStore.searchOpts.query
   }
   query.filters = ordersStore.filtersAsQueryParam
   query.sort = `${ordersStore.searchOpts.sortField} ${ordersStore.searchOpts.sortOrder}`
   router.push({query})
})

const getOrders = (() => {
   ordersStore.searchOpts.filters = [{field: "status", value: statusFilter.value, match: FilterMatchMode.EQUALS}]
   Object.entries(columnFilters.value).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         ordersStore.searchOpts.filters.push({field: key, match: data.matchMode, value: data.value})
      }
   })
   setQueryParams()
   ordersStore.getOrders()
})

const onPage = ((event) => {
   ordersStore.searchOpts.start = event.first
   ordersStore.searchOpts.limit = event.rows
   ordersStore.getOrders()
})

const onSort = ((event) => {
   ordersStore.searchOpts.sortField = event.sortField
   ordersStore.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      ordersStore.searchOpts.sortOrder = "desc"
   }
   getOrders( )
})
</script>

<style scoped lang="scss">
:deep(td.nowrap) {
   white-space: nowrap;
}
.orders {
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
      span.fee-waived {
         background: var(--uvalib-blue-alt);
         padding: 3px 10px;
         border-radius: 5px;
         color: white;
         font-weight: bold;
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