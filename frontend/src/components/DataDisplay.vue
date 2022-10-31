<template>
   <dt v-if="props.spacer" class="spacer"/>
   <dt v-else>{{props.label}}:</dt>
   <dd v-if="props.spacer"></dd>
   <template v-else>
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
dt {
   font-weight: 600;
   text-align: right;
}
.empty {
   color: #ccc;
}
.spacer {
   height: 25px;
}
</style>