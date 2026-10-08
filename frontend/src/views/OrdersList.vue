<template>
   <h2>Orders</h2>
   <div class="row-between p-2 sticky z-50 bg-white border-b border-brand-grey-100" :style="{top: headerHeight}">
      <UButton v-if="(userStore.isAdmin || userStore.isSupervisor)" color="secondary" label="Create Order" @click="createOrder()"/>
      <div class="row-right">
         <UPagination color="neutral" variant="ghost"
            v-model:page="ordersStore.searchOpts.currPage" :items-per-page="ordersStore.searchOpts.limit" 
            :total="ordersStore.total" @update:page="pageChanged"
         />
         <USelect v-model="ordersStore.searchOpts.limit" :items="[15,30,100]" @change="perPageChanged" />
      </div>
      <div class="row-left items-center">
         <div class="row-left gap-2! items-center">
            <label for="orders-filter">Filter:</label>
            <USelect id="orders-filter" v-model="statusFilter" @change="getOrders" :items="statuses" />
         </div>
         <div class="row-left gap-2! items-center">
            <label>Assigned to me:</label>
            <USwitch v-model="assignedToMe" @update:modelValue="ownerToggled"/>
         </div>
         <UInput v-model="ordersStore.searchOpts.query" placeholder="Search Job Status" @update:modelValue="queryOrders"/>
      </div>
   </div>
    <UTable :data="ordersStore.orders" :columns="columns" v-model:sorting="sorting" >
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
         <UButton color="neutral" variant="ghost" label="Title" :icon="sortIcon(column)"  @click="sortClicked(column)"/>
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
      <template #customer-header>
         <FilterPopover label="Customer"  :applied="isFilterApplied('customer')" @clear="clearFilter('customer')" @apply="applyFilters">
            <UInput v-model="columnFilters.customer.value" placeholder="Last name..." />
         </FilterPopover>
      </template>
      <template #customer-cell="{ row }">
         <div class="nowrap">{{row.original.customer.lastName}}, {{row.original.customer.firstName}}</div>
         <div class="dimmed" v-if="row.original.customer.academicStatus">({{row.original.customer.academicStatus.name}})</div>
      </template>
      <template #agency-header>
         <FilterPopover label="Agency" :applied="isFilterApplied('agency')" @clear="clearFilter('agency')" @apply="applyFilters">
            <UInput v-model="columnFilters.agency.value" placeholder="Agency name..." />
         </FilterPopover>
      </template>
      <template #agency-cell="{ row }">
         <div v-if="row.original.agency">{{ row.original.agency.name }}</div>
      </template>
      <template #processor-header>
         <FilterPopover label="Processor" :applied="isFilterApplied('processor')" @clear="clearFilter('processor')" @apply="applyFilters">
            <UInput v-model="columnFilters.processor.value" placeholder="Last name..." />
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
import { useUserStore } from '@/stores/user'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const ordersStore = useOrdersStore()
const userStore = useUserStore()

const headerHeight = computed(() => {
   let hdr = document.getElementById('uva-header')
   return `${hdr.clientHeight}px`
})

const sorting = ref([
  {
    id: 'id',
    desc: true
  }
])
const columns = [
   {
      accessorKey: 'id',
      header: 'ID',
   },
   {
      accessorKey: 'status',
      header: 'Status'
   },
   {
      accessorKey: 'dateSubmitted',
      header: 'Request Submitted',
   },
   {
      accessorKey: 'dateDue',
      header: 'Date Due',
   },
   {
      accessorKey: 'title',
      header: 'Title',
   },
   {
      accessorKey: 'specialInstructions',
      header: 'Special Instructions',
   },
   {
      accessorKey: 'unitCount',
      header: 'Units',
   },
   {
      accessorKey: 'masterFileCount',
      header: 'Master Files',
   },
   {
      accessorKey: 'fee',
      header: 'Fee',
   },
   {
      accessorKey: 'customer',
      header: 'Customer'
   },
   {
      accessorKey: 'agency',
      header: 'Agency'
   },
   {
      accessorKey: 'processor',
      header: 'Processor'
   },
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

const columnFilters = ref({
   customer: {value: null, mode: "contains"},
   processor: {value: null, mode: "contains"},
   agency: {value: null, mode: "contains"},
})

const statusFilter = ref("active")
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

const getOrders = (async () => {
   ordersStore.searchOpts.filters = [{field: "status", value: statusFilter.value, match: 'equals'}]
   Object.entries(columnFilters.value).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         ordersStore.searchOpts.filters.push({field: key, match: data.mode, value: data.value})
      }
   })
   setQueryParams()
   await ordersStore.getOrders()
})
const isFilterApplied = ((name) => {
   return columnFilters.value[name].value != null 
})
const clearFilter = ((name) => {
   columnFilters.value[name].value = null 
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