<template>
   <template v-if="props.spacer">
      <dt v-if="props.label" class="spacer label">{{props.label}}</dt>
      <dt v-else class="spacer"/>
      <dd></dd>
   </template>
   <template v-else>
      <dt>{{props.label}}:</dt>
      <dd v-if="props.value">
         <slot>
            {{props.value}}
         </slot>
      </dd>
      <dd v-else><span class="empty">{{blankValue}}</span></dd>
   </template>
</template>

<script setup>
const props = defineProps({
   label: {
      type: String,
      default: ""
   },
   value: {
      type: [String, Number],
      default: ""
   },
   blankValue: {
      type: String,
      default: "Empty"
   },
   spacer: {
      type: Boolean,
      default: false
   }
})
</script>

<style scoped lang="scss">
.empty {
   color: #ccc;
}
.spacer {
   height: 20px;
}
dt.spacer.label {
   margin-top: 25px;
   margin-bottom: 5px !important;
   font-size: 1.15em;
   color: #aaa;
   margin-bottom: 5px;
   font-weight: normal !important;
   text-align: left !important;
}
</style>