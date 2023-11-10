<template>
   <template v-if="props.spacer">
      <dt class="spacer"/>
      <dd />
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
   font-size: 1.15em;
   font-weight: bold;
   margin-bottom: 10px;
}
</style>