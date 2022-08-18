<template>
   <h2>Staff Members</h2>
   <div class="staff">
      <DataTable :value="staffStore.staff" ref="staffTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="true" :paginator="true" @page="onPage($event)"
         sortField="lastName" :sortOrder="1" @sort="onSort($event)"
         :rows="staffStore.searchOpts.limit" :totalRecords="staffStore.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[10,30,100]"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <template #header>
            <div class="filter-controls">
               <DPGButton label="Add" @click="addStaff()"/>
               <span>
                  <span class="p-input-icon-right">
                     <i class="pi pi-search" />
                     <InputText v-model="filter" placeholder="Staff Search" @input="applyFilter()"/>
                  </span>
                  <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()"/>
               </span>
            </div>
         </template>
         <Column field="lastName" header="Last Name" :sortable="true"/>
         <Column field="firstName" header="First Name"/>
         <Column field="computingID" header="UVA Computing ID" :sortable="true"/>
         <Column field="email" header="Email" :sortable="true"/>
         <Column field="role" header="Role" />
         <Column field="active" header="Active?"></Column>
         <Column header="" class="row-acts">
            <template #body="slotProps">
               <DPGButton label="Edit" class="p-button-text"  @click="edit(slotProps.data)" />
            </template>
         </Column>
      </DataTable>
      <Dialog v-model:visible="showEdit" :style="{width: '450px'}" header="Staff Member Details" :modal="true" position="top">
         <FormKit type="form" id="staff-detail" :actions="false" @submit="submitChanges">
            <FormKit label="Last Name" type="text" v-model="staffDetails.lastName" validation="required" autofocus />
            <FormKit label="First Name" type="text" v-model="staffDetails.firstName" validation="required" />
            <FormKit label="UVA Computing ID" type="text" v-model="staffDetails.computingID" validation="required" />
            <FormKit label="Email" type="email" v-model="staffDetails.email" validation="required" />
            <FormKit type="select" label="Role" v-model="staffDetails.roleID" :options="{ 0: 'Admin', 1: 'Supervisor', 2: 'Student', 3: 'Viewer' }" />
            <FormKit type="select" label="Active" v-model="staffDetails.active" :options="{ false: 'No', true: 'Yes' }" />
            <div class="form-controls">
               <FormKit type="button" label="Cancel" wrapper-class="cancel-button" @click="showEdit = false" />
               <FormKit type="submit" label="Save" wrapper-class="submit-button" />
            </div>
         </FormKit>
      </Dialog>
   </div>
</template>

<script setup>
import { onMounted, ref} from 'vue'
import { useStaffStore } from '@/stores/staff'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'

const staffStore = useStaffStore()
const filter = ref("")
const showEdit = ref(false)
const staffDetails = ref({
   id: 0,
   lastName: "",
   firstName: "",
   email: "",
   computingID: "",
   roleID: 0,
   active: false}
)

function addStaff() {
   staffDetails.value = {
      id: 0,
      lastName: "",
      firstName: "",
      email: "",
      computingID: "",
      roleID: 0,
      active: false
   }
   showEdit.value = true
}

function submitChanges() {
   let active = false
   if (staffDetails.value.active == "true") {
      active = true
   }
   let roles = ['Admin', 'Supervisor', 'Student', 'Viewer']
   staffDetails.value.active = active
   staffDetails.value.roleID = parseInt(staffDetails.value.roleID, 10)
   staffDetails.value.role = roles[ staffDetails.value.roleID ]
   staffStore.addOrUpdateStaff(staffDetails.value)
   showEdit.value = false
}

function edit(data) {
   staffDetails.value = {...data} // clone the data so edits dont change the store
   if (staffDetails.value.active) {
      staffDetails.value.active = "true"
   } else {
      staffDetails.value.active = "false"
   }
   showEdit.value = true
}

function onPage(event) {
   staffStore.searchOpts.start = event.first
   staffStore.searchOpts.limit = event.rows
   staffStore.getStaff( filter.value  )
}

function onSort(event) {
   staffStore.searchOpts.sortField = event.sortField
   staffStore.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      staffStore.searchOpts.sortOrder = "desc"
   }
   staffStore.getStaff( filter.value )
}

function applyFilter() {
   staffStore.getStaff( filter.value )
}

function clearSearch() {
   filter.value = ""
   staffStore.getStaff( filter.value )
}

onMounted(() => {
   staffStore.getStaff( filter.value  )
})
</script>

<style scoped lang="scss">
#staff-detail {
   .form-controls {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      margin-top: 5px;
      text-align: right;
      padding: 10px 0;
      :deep(.cancel-button button) {
         @include base-button();
         width: auto;
         margin-right: 10px;
      }
      :deep(.submit-button button) {
         @include primary-button();
         width: auto;
      }
   }
}
   .staff {
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