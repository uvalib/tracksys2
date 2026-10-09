<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" title="Assign Order Processor">
      <UButton @click="show" label="Assign Order Processor" color="secondary"/>
      <template  #body>
         <div class="candidate-scroller">
            <div class="val" v-for="(c,idx) in staffStore.staff" :key="c.id" :class="{selected: idx == selectedIdx}"
               @click="selectProcessor(idx)"
            >
               <span class="candidate">{{c.lastName}}, {{c.firstName}}</span> ({{c.computingID}})
            </div>
         </div>
         <p class="error">{{error}}</p>
      </template>
      <template #footer>
         <UButton @click="hide" label="Cancel" color="secondary"/>
         <UButton @click="assignClicked" label="Assign"/>
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useOrdersStore } from '@/stores/orders'
import { useStaffStore } from '@/stores/staff'

const ordersStore = useOrdersStore()
const staffStore = useStaffStore()

const isOpen = ref(false)
const selectedIdx = ref(-1)
const error = ref("")

function selectProcessor(idx) {
   selectedIdx.value = idx
}
function assignClicked() {
   error.value = ""
   if ( selectedIdx.value == -1) {
      error.value = "Please select a new order processor"
      return
   }
   let staffID = staffStore.staff[selectedIdx.value].id
   ordersStore.setProcessor( staffID )
   hide()
}
function hide() {
   isOpen.value=false
}
function show() {
   staffStore.getAdmins()
   isOpen.value = true
   error.value = ""
   selectedIdx.value = -1
}
</script>

<style lang="scss" scoped>
.error {
   padding: 0;
   margin: 0;
   text-align: center;
   color: var(--uvalib-red-emergency);
}

.candidate-scroller {
   max-height: 300px;
   overflow: scroll;
   padding: 0;
   margin:  0;
   border: 1px solid var(--uvalib-grey-light);
   border-radius: 4px;
   .val {
      padding: 2px 10px 3px 10px;
      cursor: pointer;
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      &:hover  {
         background: var(--uvalib-blue-alt-light);
      }
   }
   .val.selected {
      background: var(--uvalib-blue-alt);
      color: white;
   }
   .candidate {
      font-weight: bold;
   }
}
</style>
