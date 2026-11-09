<template>
   <UPopover>
      <UButton color="secondary" size="sm" icon="i-lucide-info" label="What is searched?" />
      <template #content>
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
      </template>
   </UPopover>
</template>

<script setup>
</script>

<style lang="scss" scoped>
.content {
   max-width: 400px;
   display: flex;
   flex-direction: column;
   gap: 10px;
   padding: 15px;
   font-size: 0.85em;
   table {
      border: 1px solid var(--uvalib-grey-light);
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
   b {
      font-weight: bold
   }
}
</style>
