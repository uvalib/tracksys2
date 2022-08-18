<template>
   <h2>Customers</h2>
   <div class="customers">
      <DataTable :value="customersStore.customers" ref="customerTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="true" :paginator="true" @page="onPage($event)"
         sortField="lastName" :sortOrder="1" @sort="onSort($event)"
         :rows="customersStore.searchOpts.limit" :totalRecords="customersStore.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[10,30,100]"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <template #header>
            <div class="filter-controls">
               <DPGButton label="Add" @click="addCustomer()"/>
               <span>
                  <span class="p-input-icon-right">
                     <i class="pi pi-search" />
                     <InputText v-model="filter" placeholder="Customer Search" @input="applyFilter()"/>
                  </span>
                  <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()"/>
               </span>
            </div>
         </template>
         <Column field="lastName" header="Last Name" :sortable="true"/>
         <Column field="firstName" header="First Name"/>
         <Column field="email" header="Email" :sortable="true"/>
         <Column field="academicStatus.name" header="Acedemic Status"/>
         <Column header="" class="row-acts">
            <template #body="slotProps">
               <DPGButton label="Edit" class="p-button-text"  @click="edit(slotProps.data)" />
            </template>
         </Column>
      </DataTable>
   </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useCustomersStore } from '@/stores/customers'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'

const customersStore = useCustomersStore()

const filter = ref("")

function addCustomer() {
   // console.log("ADD CUSTOMER")
}

function onPage(event) {
   customersStore.searchOpts.start = event.first
   customersStore.searchOpts.limit = event.rows
   customersStore.getCustomers( filter.value  )
}

function onSort(event) {
   customersStore.searchOpts.sortField = event.sortField
   customersStore.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      customersStore.searchOpts.sortOrder = "desc"
   }
   customersStore.getCustomers( filter.value )
}

function applyFilter() {
   customersStore.getCustomers( filter.value )
}

function clearSearch() {
   filter.value = ""
   customersStore.getCustomers( filter.value )
}

onMounted(() => {
   customersStore.getCustomers( filter.value  )
})
</script>

<style scoped lang="scss">
.customers {
   min-height: 600px;
   text-align: left;
   padding: 0 25px 25px 25px;
   .filter-controls {
      display: flex;
      flex-flow: row wrap;
      justify-content: space-between;
      button.p-button-secondary.p-button {
         margin-left: 5px;
      }
   }
   :deep(.row-acts) {
      text-align: center;
   }
}
</style>