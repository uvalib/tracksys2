<template>
   <UPopover>
      <UButton color="neutral" variant="soft" size="sm" :label="customerInfo" />
      <template #content>
         <div class="px-4">
            <UTabs :items="tabs" variant="link" valueKey="slot">
               <template #customer>
                  <dl>
                     <DataDisplay label="Last Name" :value="props.customer.lastName"></DataDisplay>
                     <DataDisplay label="First Name" :value="props.customer.firstName"></DataDisplay>
                     <DataDisplay label="Email" :value="props.customer.email"></DataDisplay>
                     <DataDisplay label="Academic Status" :value="props.customer.academicStatus.name"></DataDisplay>
                  </dl>
               </template>
               <template v-for="info in addresses" #[info.name]>
                  <dl>
                     <DataDisplay label="Address 1" :value="info.data.address1"></DataDisplay>
                     <DataDisplay v-if="info.data.address2" label="Address 2" :value="info.data.address2"></DataDisplay>
                     <DataDisplay v-if="info.data.city" label="City" :value="info.data.city"></DataDisplay>
                     <DataDisplay v-if="info.data.state" label="State" :value="info.data.state"></DataDisplay>
                     <DataDisplay v-if="info.data.zip"  label="Zip" :value="info.data.zip"></DataDisplay>
                     <DataDisplay v-if="info.data.phone" label="Phone" :value="info.data.phone"></DataDisplay>
                  </dl>
               </template>
            </UTabs>
         </div>
      </template>
   </UPopover>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
   customer: {
      type: Object,
      required: true
   }
})

const customerInfo = computed(() => {
   let cust = `${props.customer.lastName}, ${props.customer.firstName}`
   if (props.customer.academicStatus.id != 0) {
      cust += ` (${props.customer.academicStatus.name})`
   }
   return cust
})

const addresses = computed(() => {
   let out = [ {name: "primary", data: props.customer.addresses[0]} ]
   if (props.customer.addresses > 1) {
         out.push( { label: `billling`, data:  props.customer.addresses[1]}  )
      }
   return out
})

const tabs = computed(() => {
   let out = []
   out.push( { label: `Customer`, slot: 'customer'} )
   if ( props.customer.addresses ) {
      out.push( { label: `Primary Address`, slot: 'primary'}  )
      if (props.customer.addresses > 1) {
         out.push( { label: `Billing Address`, slot: 'billing'}  )
      }
   }
   return out
})
</script>

<style lang="scss" scoped>
dl {
   margin: 0;
   display: inline-grid;
   grid-template-columns: max-content 1fr;
   grid-column-gap: 10px;
   font-size: 0.9em;
   text-align: left;
   box-sizing: border-box;

   dt {
      font-weight: bold;
      text-align: right;
   }

   dd {
      margin: 0 0 5px 0;
      word-break: break-word;
      -webkit-hyphens: auto;
      -moz-hyphens: auto;
      hyphens: auto;
      white-space: break-spaces;
      margin-inline-start: 5px;
   }
}
</style>