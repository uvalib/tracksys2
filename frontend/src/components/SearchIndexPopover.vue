<template>
<DPGButton severity="secondary" size="small" icon="pi pi-info-circle" type="button" label="What is searched?" @click="toggle" />
<Popover ref="po">
   <div class="content">
      <table>
         <tbody>
            <tr>
               <td class="label">Orders</td>
               <td>ID, Title, Agency, Customer, Staff Notes, Special Instructions</td>
            </tr>
            <tr>
               <td class="label">Metadata</td>
               <td>ID, PID, Title, Barcode, Call Number, Catalog Key, Creator Name</td>
            </tr>
            <tr>
               <td class="label">Units</td>
               <td>ID, Staff Notes, Special Instructions</td>
            </tr>
            <tr>
               <td class="label">Master Files</td>
               <td>ID, PID, Title, Description, Filename, Call Number, Tag</td>
            </tr>
            <tr>
               <td class="label">Components</td>
               <td>ID, PID, Title, Label, Description, Data, Finding Aid</td>
            </tr>
         </tbody>
      </table>
      <div class="important">
         <b>IMPORTANT</b>
         <div>
            The search data is updated nightly, so changes will not appear in the search results
            until the morning after the changes are made.
         </div>
      </div>
   </div>
</Popover>
</template>

<script setup>
import { ref } from 'vue'
import Popover from 'primevue/popover'

const po = ref()

const toggle = ((event) => {
   po.value.toggle(event);
})
</script>

<style lang="scss" scoped>
.content {
   max-width: 400px;
   display: flex;
   flex-direction: column;
   gap: 10px;
   table {
      border: 1px solid var(--uvalib-grey-light);
      border-radius: 0.3rem;
      margin-bottom: 10px;
      padding: 5px;
      border-collapse: collapse;
      td {
         vertical-align: top;
         padding: 10px;
         border: 1px solid var(--uvalib-grey-light);
      }
      td.label {
         text-align: right;
         font-weight: 600;
      }
   }
}
</style>
