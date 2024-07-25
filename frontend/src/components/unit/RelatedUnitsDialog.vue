<template>
   <DPGButton @click="show" severity="secondary" label="More"/>
   <Dialog v-model:visible="isOpen" :modal="true" :header="`Units for Order ${props.orderID}`" :style="{width: '90%'}">
      <div class="related-unit-ids">
         <template v-for="(uid,idx) in props.units" :key="`related-${uid}`">
            <template v-if="idx > 0"><span class="sep"></span></template>
            <router-link :to="`/units/${uid}`" v-if="uid != currentUnitID">{{uid}}</router-link>
            <span class="current-unit" v-else>{{uid}}</span>
         </template>
      </div>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import Dialog from 'primevue/dialog'

const props = defineProps({
   units: {
      type: Array,
      required: true
   },
   orderID: {
      type: Number,
      reqtured: true
   },
   currentUnitID: {
      type: Number,
      reqtured: true
   }
})


const isOpen = ref(false)

function hide() {
   isOpen.value=false
}
function show() {
   isOpen.value = true

}
</script>

<style lang="scss" scoped>
.related-unit-ids {
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   a {
      color: var(--uvalib-brand-blue-light);
      font-weight: 600;
      text-decoration: none;
      &:hover {
         text-decoration: underline;
      }
   }
   .sep {
      margin-right: 5px;
      display: inline-block;
   }
   .current-unit {
      font-weight: bold;
      background: var(--uvalib-teal-lightest);
      padding: 2px 4px;
      border-radius: 5px;
   }
}
</style>
