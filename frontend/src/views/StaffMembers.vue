<template>
   <h2>Staff Members</h2>
   <div class="staff">
      <DataTable :value="staffStore.staff" ref="staffTable" dataKey="id"
         stripedRows showGridlines responsiveLayout="scroll"
         :lazy="true" :paginator="true" @page="onPage($event)"
         :rows="staffStore.searchOpts.limit" :totalRecords="staffStore.total"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         :rowsPerPageOptions="[10,30,100]"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
      >
         <Column field="lastName" header="Last Name"></Column>
         <Column field="firstName" header="First Name"></Column>
         <Column field="computingID" header="UVA Computing ID"></Column>
         <Column field="email" header="Email"></Column>
         <Column field="role" header="Role"></Column>
         <Column field="active" header="Active?"></Column>
         <!-- <Column header="" class="row-acts">
            <template #body="slotProps">
               <router-link :to="`/jobs/${slotProps.data.id}`">View</router-link>
               <DPGButton label="Delete"  class="p-button-text" @click="deleteJob(slotProps.data.id)"/>
            </template>
         </Column> -->
      </DataTable>
   </div>
</template>

<script setup>
import { onMounted} from 'vue'
import { useStaffStore } from '@/stores/staff'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const staffStore = useStaffStore()

function onPage(event) {
   staffStore.searchOpts.start = event.first
   staffStore.searchOpts.limit = event.rows
   staffStore.getStaff()
}

onMounted(() => {
   staffStore.getStaff()
})
</script>

<style scoped lang="scss">
   .staff {
      min-height: 600px;
      text-align: left;
      padding: 0 25px 25px 25px;
      .p-datatable {
         font-size: 0.85em;
         :deep(td), :deep(th) {
            padding: 10px;
         }
         :deep(.row-acts) {
            text-align: center;
            padding: 0;
            a {
               display: inline-block;
               margin-right: 15px;
            };
         }
      }
   }
</style>