<template>
   <h2>Staff Members</h2>

   <div class="row-between p-2 sticky z-50 bg-white border-b border-brand-grey-100" :style="{top: system.headerHeight}">
      <UButton  v-if="(userStore.isAdmin || userStore.isSupervisor)" label="Add Staff" color="secondary" @click="addStaff" />
      <UInput v-model="filter" placeholder="Search staff..." @update:modelValue="applyFilter"/>
   </div>
   <UTable :data="staffStore.staff" :columns="columns" :meta="meta" v-model:global-filter="filter" 
      sticky virtualize class="h-[calc(100vh-185px)]" 
   >
      <template #acts-cell="{ row }">
         <UButton label="Edit" color="secondary" size="xs" @click="edit(row.original)"/>
      </template>  
   </UTable>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useStaffStore } from '@/stores/staff'
import { useUserStore } from '../stores/user'
import { useSystemStore } from '../stores/system'

import StaffEditModal from '@/components/StaffEditModal.vue'

const staffStore = useStaffStore()
const userStore = useUserStore()
const system = useSystemStore()
const staffModal = useOverlay().create(StaffEditModal)

const filter = ref("")

const columns = [
   { accessorKey: 'id', header: "ID"},
   { accessorKey: 'lastName', header: "Last Name"},
   { accessorKey: 'firstName', header: "First Name"},
   { accessorKey: 'computingID', header: "UVA Computing ID"},
   { accessorKey: 'email', header: "Email"},
   { accessorKey: 'role', header: "Role"},
   { accessorKey: 'active', header: "Active"},
   { accessorKey: 'acts', header: ""},
]

const meta = {
  class: {
    tr: (row) => {
      if (!row.original.active) {
        return 'bg-brand-grey-200/60'
      } 
      return ''
    }
  }
}

const addStaff = (() => {
   staffModal.open({mode: "add", title: "Add Staff Member"})
})

const edit = ((data) => {   
   staffModal.open({mode: "edit", title: "Edit Staff Member", staff: {...data}})
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