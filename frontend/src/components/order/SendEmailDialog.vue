<template>
   <DPGButton @click="show" :label="buttonLabel" class="p-button-secondary"/>
   <Dialog v-model:visible="isOpen" :modal="true" :header="buttonLabel">
      <div class="choice border">
        <input type="checkbox" v-model="sendToCustomer"/>
        <span>Send to customer email: {{ordersStore.detail.customer.email}}</span>
      </div>
      <div class="choice">
         <input type="checkbox" v-model="sendToAlt"/>
         <label>Send to alternate email:<input type="text" v-model="altEmail"/></label>
      </div>
      <p class="error">{{error}}</p>
      <template #footer>
         <DPGButton @click="hide" label="Cancel" class="p-button-secondary"/>
         <span class="spacer"></span>
         <DPGButton autofocus @click="sendClicked" label="Send"/>
      </template>
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import {useOrdersStore} from '@/stores/orders'
import {useUserStore} from '@/stores/user'
import Dialog from 'primevue/dialog'

const ordersStore = useOrdersStore()
const user = useUserStore()

const props = defineProps({
   mode: {
      type: String,
      default: "order",
   },
})

const isOpen = ref(false)
const error = ref("")
const sendToCustomer = ref(true)
const sendToAlt = ref(false)
const altEmail = ref("")

const buttonLabel = computed(() => {
   if (props.mode == "fee") {
      return "Resend Fee Estimate"
   }
   return "Send Email"
})

function sendClicked() {
   error.value = ""
   if (sendToAlt.value && altEmail.value == "") {
      error.value = "An alternate email address is required."
      return
   }
   if (props.mode == "order") {
      ordersStore.sendEmail(user.ID, sendToCustomer.value, sendToAlt.value, altEmail.value)
   } else {
      ordersStore.resendFeeEstimate( user.ID, sendToCustomer.value, sendToAlt.value, altEmail.value)
   }
   hide()
}

function hide() {
   isOpen.value=false
}
function show() {
   isOpen.value = true
   error.value = ""
}
</script>

<style lang="scss" scoped>
   .choice {
      padding: 5px;
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: flex-start;
      span, label {
         position: relative;
         top: 3px;
      }

      input[type=checkbox] {
         width: 18px;
         height: 18px;
         margin-right: 10px;
      }
      input[type=text] {
         margin-top: 5px;
      }
   }
   .choice.border {
      padding-bottom: 10px;
      margin-bottom: 10px;
   }
.error {
   padding: 0;
   margin: 0;
   text-align: center;
   color: var(--uvalib-red-emergency);
}
</style>
