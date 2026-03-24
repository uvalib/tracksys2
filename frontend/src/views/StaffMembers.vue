<template>
   <h2>Staff Members</h2>
   <div class="staff">
      <DataTable :value="staffStore.staff" ref="staffTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="true" :paginator="true" @page="onPage($event)"
         sortField="lastName" :sortOrder="1" @sort="onSort($event)"
         :rows="staffStore.searchOpts.limit" :totalRecords="staffStore.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[30,50,100]" paginatorPosition="top"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <template #paginatorstart  v-if="(userStore.isAdmin || userStore.isSupervisor)" >
            <DPGButton label="Add Staff" severity="secondary" @click="addStaff()"/>
         </template>
         <template #paginatorend>
            <IconField iconPosition="left">
               <InputIcon class="pi pi-search" />
               <InputText v-model="filter" placeholder="Search Staff" @input="applyFilter()" />
            </IconField>
         </template>
         <Column field="id" header="ID" :sortable="true"/>
         <Column field="lastName" header="Last Name" :sortable="true"/>
         <Column field="firstName" header="First Name"/>
         <Column field="computingID" header="UVA Computing ID" :sortable="true"/>
         <Column field="email" header="Email" :sortable="true"/>
         <Column field="role" header="Role" />
         <Column field="active" header="Active?"></Column>
         <Column header="" class="row-acts">
            <template #body="slotProps">
               <DPGButton label="Edit" class="edit-btn" severity="secondary"  @click="edit(slotProps.data)" />
            </template>
         </Column>
      </DataTable>
      <Dialog v-model:visible="showEdit" :style="{width: '450px'}" header="Staff Member Details" :modal="true" position="top" :closable="false">
         <form @submit="submitChanges" id="staff-detail">
            <FormField id="lname" label="Last Name" :error="errors.lastName" :required="true">
               <InputText id="lname" v-model="lastName" type="text" autofocus/>   
            </FormField>
            <FormField id="fname" label="First Name" :error="errors.firstName" :required="true">
               <InputText id="fname" v-model="firstName" type="text" />   
            </FormField>
            <FormField id="cid" label="UVA Computing ID" :error="errors.computingID" :required="true">
               <InputText id="cid" v-model="computingID" type="text" />   
            </FormField>
            <FormField id="email" label="Email" :error="errors.email" :required="true">
               <InputText id="email" v-model="email" type="text" />   
            </FormField>
            <FormField id="role" label="Role" :error="errors.roleID" :required="true">
               <Select id="role" v-model="roleID"  :options="roles" optionLabel="label" optionValue="id" placeholder="Select a role" />   
            </FormField>
             <FormField id="active" label="Active">
               <Select id="active" v-model="active"  :options="[{label: 'No', val: 'false'},{label: 'Yes', val: 'true'}]" optionLabel="label" optionValue="val" />   
            </FormField>
            <div class="notes">
               <b>IMPORTANT:</b>
               <span>All new staff must be added to a group named lb-digiserv. This can be done here: </span>
               <a href="https://mygroups.virginia.edu/groups/" target="_blank">MyGroups</a>
            </div>
            <div class="form-controls">
               <DPGButton label="Cancel" severity="secondary" @click="showEdit=false"/>
               <DPGButton label="Save" type="submit" />
            </div>
         </form>
      </Dialog>
   </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useStaffStore } from '@/stores/staff'
import { useUserStore } from '../stores/user'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Dialog from 'primevue/dialog'

import { usePinnable } from '@/composables/pin'
usePinnable("p-datatable-paginator-top")

import { useForm } from 'vee-validate'
import * as yup from 'yup'
import FormField from '@/components/FormField.vue'

const { errors, resetForm, handleSubmit, defineField } = useForm({
  validationSchema: yup.object().shape({
      lastName: yup.string().required('Last name is required'),
      firstName: yup.string().required('First name is required'),
      computingID: yup.string().required('ComputingID is required'),
      email: yup.string().email("Email is invalid").required("Email is required"),
      roleID: yup.string().required("Role is required"),
   })
})

const [lastName] = defineField('lastName')
const [firstName] = defineField('firstName')
const [email] = defineField('email')
const [computingID] = defineField('computingID')
const [roleID] = defineField('roleID')
const [active] = defineField('active')

const staffStore = useStaffStore()
const userStore = useUserStore()

const filter = ref("")
const showEdit = ref(false)

const roles = computed( () => {
   return [
      {label:"Admin", id: 0}, {label:"Supervisor", id: 1}, 
      {label:"Student", id: 2}, {label:"Viewer", id: 3}
   ]
})

const addStaff = (() => {
   resetForm({ 
      values: {
         id: 0,
         lastName: "",
         firstName: "",
         email: "",
         computingID: "",
         roleID: null,
         active: "true"
      } 
   })
   showEdit.value = true
})


const submitChanges = handleSubmit(values => {
   let roles = ['Admin', 'Supervisor', 'Student', 'Viewer']
   values.roleID = parseInt(values.roleID, 10)
   values.role = roles[ values.roleID ]
   if (values.active == "true") {
      values.active = true
   } else {
      values.active = false
   }
   staffStore.addOrUpdateStaff(values)
   showEdit.value = false
}) 

const edit = ((data) => {   
   let editData  = {...data}
   editData.active = "false"
   if (editData.active) {
      editData.active = "true"
   }
   resetForm({ values: editData })
   showEdit.value = true
})

const onPage = ((event) => {
   staffStore.searchOpts.start = event.first
   staffStore.searchOpts.limit = event.rows
   staffStore.getStaff( filter.value  )
})

const onSort = ((event) => {
   staffStore.searchOpts.sortField = event.sortField
   staffStore.searchOpts.sortOrder = "asc"
   if (event.sortOrder == -1) {
      staffStore.searchOpts.sortOrder = "desc"
   }
   staffStore.getStaff( filter.value )
})

const applyFilter = (() => {
   staffStore.getStaff( filter.value )
})

onMounted(() => {
   staffStore.getStaff( filter.value  )
   document.title = `Staff Members`
})
</script>

<style scoped lang="scss">
#staff-detail {
   display: flex;
   flex-direction: column;
   gap: 15px;

   .form-controls {
      display: flex;
      flex-flow: row nowrap;
      gap: 10px;
      justify-content: flex-end;
   }
}
.notes {
   font-size: 0.9em;
   margin: 15px 0 10px 0;
   border:  1px solid var(--uvalib-teal);
   padding: 10px 20px;
   border-radius: 5px;
   background-color: var(--uvalib-teal-lightest);
   b {
      padding-right: 10px;
   }
   a {
      color: var(--uvalib-blue-alt-dark);
   }
}
.staff {
   min-height: 600px;
   text-align: left;
   padding: 0 25px 25px 25px;

   .js-search {
      margin-right: 10px;
   }

   :deep(td.row-acts) {
      text-align: center;
      .edit-btn {
         font-size: 0.85em;
         padding: 2px 12px;
      }
   }
}
</style>