<template>
   <h2>Orders</h2>
   <div class="row-between p-2 bg-white border-b border-brand-grey-100">
      <UButton v-if="(userStore.isAdmin || userStore.isSupervisor)" color="secondary" label="Create Order" @click="createOrder()"/>
      <div class="row-right gap-4">
         <UPagination color="neutral" variant="ghost" show-edges
            v-model:page="ordersStore.searchOpts.currPage" :items-per-page="ordersStore.searchOpts.limit" 
            :total="ordersStore.total" @update:page="pageChanged"
         />
         <USelect v-model="ordersStore.searchOpts.limit" :items="[15,30,100]" @change="perPageChanged" />
         <span>Total Orders: {{ ordersStore.total }}</span>
      </div>
      <div class="row-left items-center gap-4">
         <div class="row-left gap-2! items-center">
            <label for="orders-filter">Filter:</label>
            <USelect id="orders-filter" v-model="filters.status.value" @change="getOrders" :items="statuses" />
         </div>
         <div class="row-left gap-2! items-center">
            <label>Assigned to me:</label>
            <USwitch v-model="assignedToMe" @update:modelValue="ownerToggled"/>
         </div>
      </div>
   </div>
    <UTable :data="ordersStore.orders" :columns="columns" v-model:sorting="sorting" sticky class="h-[calc(100vh-185px)]" >
      <template #id-header="{ column }">
         <UButton color="neutral" variant="ghost" label="ID" :icon="sortIcon(column)"  @click="sortClicked(column)"/>
      </template>
      <template #id-cell="{ row }">
         <router-link :to="`/orders/${row.original.id}`">{{row.original.id}}</router-link>
      </template>
      <template #status-cell="{ row }">
         <span :class="`status ${row.original.status}`">{{displayStatus(row.original.status)}}</span>
      </template>
      <template #dateSubmitted-header="{ column }">
         <UButton color="neutral" variant="ghost" :icon="sortIcon(column)" @click="sortClicked(column)"> 
            <span class="text-start">Request Submitted</span>
         </UButton>
      </template>
      <template #dateDue-header="{ column }">
         <UButton color="neutral" variant="ghost" :icon="sortIcon(column)"  @click="sortClicked(column)">
            <span class="text-start">Date Due</span>
         </UButton>
      </template>
      <template #title-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @clear="clearFilter(column.id)" @apply="applyFilters">
            <UInput v-model="filters.title.value" class="w-full" placeholder="Title..." />   
         </FilterPopover>
      </template>
      <template #specialInstructions-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @clear="clearFilter(column.id)" @apply="applyFilters">
            <UInput v-model="filters.specialInstructions.value" class="w-full" placeholder="Instructions..." />   
         </FilterPopover>
      </template>
      <template #unitCount-header="{ column }">
         <UButton color="neutral" variant="ghost" label="Units" :icon="sortIcon(column)"  @click="sortClicked(column)"/>
      </template>
      <template #masterFileCount-header="{ column }">
         <UButton color="neutral" variant="ghost" :icon="sortIcon(column)"  @click="sortClicked(column)">
            <span class="text-start">Master Files</span>
         </UButton>
      </template>
      <template #fee-header="{ column }">
         <UButton color="neutral" variant="ghost" label="Fee" :icon="sortIcon(column)"  @click="sortClicked(column)"/>
      </template>
      <template #fee-cell="{ row }">
         <span class="fee-waived" v-if="row.original.feeWaived">Waived</span>
         <span class="fee" v-else-if="row.original.fee !== undefined">${{parseFloat(row.original.fee).toFixed(2)}}</span>
      </template>
      <template #customer-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @clear="clearFilter(column.id)" @apply="applyFilters">
            <USelectMenu v-model="filters.customer.value" :items="customers" class="w-full" 
               valueKey="value" labelKey="label" placeholder="Select customer..." virtualize/>   
         </FilterPopover>
      </template>
      <template #customer-cell="{ row }">
         <div class="nowrap">{{row.original.customer.lastName}}, {{row.original.customer.firstName}}</div>
         <div class="dimmed" v-if="row.original.customer.academicStatus">({{row.original.customer.academicStatus.name}})</div>
      </template>
      <template #agency-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @clear="clearFilter(column.id)" @apply="applyFilters">
             <USelectMenu v-model="filters.agency.value" :items="system.agencies" class="w-full" 
               valueKey="id" labelKey="name" descriptionKey="" placeholder="Select an agency..." virtualize/>   
         </FilterPopover>
      </template>
      <template #agency-cell="{ row }">
         <div v-if="row.original.agency">{{ row.original.agency.name }}</div>
      </template>
      <template #processor-header="{ column }">
         <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @clear="clearFilter(column.id)" @apply="applyFilters">
            <USelectMenu v-model="filters.processor.value" :items="staff" class="w-full" 
               valueKey="value" labelKey="label" placeholder="Select staff..." virtualize/>   
         </FilterPopover>
      </template>
      <template #processor-cell="{ row }">
         <span v-if="row.original.processor">{{row.original.processor.lastName}}, {{row.original.processor.firstName}}</span>
      </template>
    </UTable>
</template>

<script setup>
import { onBeforeMount, onMounted, ref, computed } from 'vue'
import { useOrdersStore } from '@/stores/orders'
import { useSystemStore } from '@/stores/system'
import { useStaffStore } from '@/stores/staff'
import { useCustomersStore } from '@/stores/customers'
import { useUserStore } from '@/stores/user'
import { useRoute, useRouter } from 'vue-router'
import FilterPopover from '@/components/FilterPopover.vue'

const route = useRoute()
const router = useRouter()
const ordersStore = useOrdersStore()
const userStore = useUserStore()
const customersStore = useCustomersStore()
const staffStore = useStaffStore()
const system = useSystemStore()

const sorting = ref([
  {
    id: 'id',
    desc: true
  }
])
const columns = [
   { accessorKey: 'id',header: 'ID' },
   { accessorKey: 'status', header: 'Status' },
   { accessorKey: 'dateSubmitted', header: 'Request Submitted' },
   { accessorKey: 'dateDue', header: 'Date Due' },
   { accessorKey: 'title', header: 'Title' },
   { accessorKey: 'specialInstructions', header: 'Special Instructions' },
   { accessorKey: 'unitCount', header: 'Units' },
   { accessorKey: 'masterFileCount', header: 'Master Files' },
   { accessorKey: 'fee', header: 'Fee' },
   { accessorKey: 'customer', header: 'Customer' },
   { accessorKey: 'agency', header: 'Agency' },
   { accessorKey: 'processor', header: 'Processor' },
]

const statuses = ref([
   {label: "Active", value: "active"},
   {label: "Await Approval", value: "await"},
   {label: "Deferred", value: "deferred"},
   {label: "Complete", value: "complete"},
   {label: "Canceled", value: "canceled"},
   {label: "Due in a Week", value: "due_week"},
   {label: "Overdue", value: "overdue"},
   {label: "Ready for Delivery", value: "ready"}
])

const filters = ref({
   status: {value: "active", mode: "equals"},
   title: {value: null, mode: "contains"},
   specialInstructions: {value: null, mode: "contains"},
   customer: {value: null, mode: "equals"},
   processor: {value: null, mode: "equals"},
   agency: {value: null, mode: "equals"},
})

const assignedToMe = ref(false)

const sortIcon = ((column) => {
   const isSorted = column.getIsSorted()
   if (isSorted) {
      if (isSorted === 'asc') return 'i-lucide-arrow-up-narrow-wide'
      return 'i-lucide-arrow-down-wide-narrow'
   }
   return 'i-lucide-arrow-up-down'
})
const sortClicked = (async (column) => {
   const isSorted = column.getIsSorted()
   console.log(isSorted)
   ordersStore.searchOpts.sortField = column.id
   ordersStore.searchOpts.sortOrder = "asc"
   if (isSorted === "asc") {
      ordersStore.searchOpts.sortOrder = "desc"   
   }
   ordersStore.resetSearch()
   await getOrders( )
   column.toggleSorting(column.getIsSorted() === 'asc')
})

const customers = computed(() => {
   let out = []
   customersStore.customers.forEach( a => {
      out.push( {label: `${a.lastName}, ${a.firstName}`, value: a.id} )
   })
   return out
})
const staff = computed(() => {
   let out = []
   staffStore.staff.forEach( a => {
      out.push( {label: `${a.lastName}, ${a.firstName}`, value: a.id} )
   })
   return out
})

onBeforeMount( () => {
   ordersStore.searchOpts.filters = []
   if ( route.query.filters ) {
      let qpFilters = JSON.parse(route.query.filters)
      qpFilters.forEach( filter => {
         let bits = filter.split("|")
         ordersStore.searchOpts.filters.push( {field: bits[0], match: bits[1], value: bits[2]} )
         if ( bits[0] == "status") {
            filters.value.status.value = bits[2]
         } else if ( bits[0] == "customer") {
            filters.value.customer.value = bits[2]
         } else if ( bits[0] == "processor") {
            filters.value.processor.value = bits[2]
         } else if ( bits[0] == "agency") {
            filters.value.agency.value = bits[2]
         }  else if ( bits[0] == "title") {
            filters.value.title.value = bits[2]
         } else if ( bits[0] == "specialInstructions") {
            filters.value.specialInstructions.value = bits[2]
         }
      })
   }
})

onMounted( () => {
   customersStore.getCustomers()
   staffStore.getStaff()
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
   query.filters = ordersStore.filtersAsQueryParam
   query.sort = `${ordersStore.searchOpts.sortField} ${ordersStore.searchOpts.sortOrder}`
   router.push({query})
})

const getOrders = (async () => {
   ordersStore.searchOpts.filters = []
   Object.entries(filters.value).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         ordersStore.searchOpts.filters.push({field: key, match: data.mode, value: data.value})
      }
   })
   setQueryParams()
   await ordersStore.getOrders()
})
const isFilterApplied = ((name) => {
   return filters.value[name].value != null 
})
const clearFilter = ((name) => {
   filters.value[name].value = null 
   ordersStore.resetSearch()
   getOrders()
})
const applyFilters = (() => {
   ordersStore.resetSearch()
   getOrders()  
})

const perPageChanged = (() => {
   ordersStore.searchOpts.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   ordersStore.searchOpts.start = (ordersStore.searchOpts.currPage-1) * ordersStore.searchOpts.limit
   ordersStore.getOrders()
})
</script>

<style scoped lang="scss">
span.fee-waived {
   background: var(--uvalib-blue-alt);
   padding: 3px 10px;
   border-radius: 5px;
   color: white;
   font-weight: normal;
}
</style>