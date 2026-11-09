<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" :title="buttonLabel">
      <UButton @click="show" :label="buttonLabel" color="secondary"/>
      <template #body>
         <div class="email">
            <div class="row-left">
               <UCheckbox id="tocustomer" size="lg" v-model="sendToCustomer" binary />
               <label for="tocustomer">Send to customer email: {{ordersStore.detail.customer.email}}</label>
            </div>
            <div class="row-left">
               <UCheckbox id="usealtemail" v-model="sendToAlt" binary />
               <label for="usealtemail">Send to alternate email</label>
            </div>
            <div class="row-left">
               <label for="altemail">Alternate email:</label>
               <UInput id="altemail" class="flex-1" v-model="altEmail" fluid/>
            </div>
            <p class="error">{{error}}</p>
         </div>
      </template>
      <template #footer>
         <UButton @click="hide" label="Cancel" color="secondary"/>
         <UButton autofocus @click="sendClicked" label="Send"/>
      </template>
   </UModal>
</template>

<script setup>
import { ref, computed } from 'vue'
import {useOrdersStore} from '@/stores/orders'
import {useUserStore} from '@/stores/user'

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
      ordersStore.sendEmail(user.computeID, sendToCustomer.value, sendToAlt.value, altEmail.value)
   } else {
      ordersStore.resendFeeEstimate( user.computeID, sendToCustomer.value, sendToAlt.value, altEmail.value)
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
.email {
   display: flex;
   flex-direction: column;
   gap: 15px;
}
.error {
   padding: 0;
   margin: 0;
   text-align: center;
   color: var(--uvalib-red-emergency);
}
</style>
