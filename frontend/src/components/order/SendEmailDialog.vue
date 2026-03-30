<template>
   <DPGButton @click="show" :label="buttonLabel" severity="secondary"/>
   <Dialog v-model:visible="isOpen" :modal="true" :header="buttonLabel">
      <div class="email">
         <div class="choice">
            <Checkbox id="tocustomer" v-model="sendToCustomer" binary />
            <label for="tocustomer">Send to customer email: {{ordersStore.detail.customer.email}}</label>
         </div>
         <div class="choice">
            <Checkbox id="usealtemail" v-model="sendToAlt" binary />
            <label for="usealtemail">Send to alternate email</label>
         </div>
         <div class="choice leftpad">
            <label for="altemail">Alternate email:</label>
            <InputText id="altemail"  v-model="altEmail" fluid/>
         </div>
         <p class="error">{{error}}</p>
      </div>
      <template #footer>
         <DPGButton @click="hide" label="Cancel" severity="secondary"/>
         <DPGButton autofocus @click="sendClicked" label="Send"/>
      </template>
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import {useOrdersStore} from '@/stores/orders'
import {useUserStore} from '@/stores/user'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Checkbox from 'primevue/checkbox'

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
.choice {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: center;
   gap: 10px;
   label {
      white-space: nowrap;
      flex-grow: 1;
   }
}
.leftpad {
   margin-left: 30px;
}

.error {
   padding: 0;
   margin: 0;
   text-align: center;
   color: var(--uvalib-red-emergency);
}
</style>
