<template>
   <dt>{{props.label}}</dt>
   <dd v-if="props.value">
      <slot>
         {{props.value}}
      </slot>
   </dd>
   <dd v-else><span class="empty">{{blankValue}}</span></dd>
</template>

<script setup>
const props = defineProps({
   label: {
      type: String,
      required: true
   },
   value: {
      type: String,
      required: true
   },
   blankValue: {
      type: String,
      default: "Empty"
   },
})
</script>

<style scoped lang="scss">
.empty {
   color: #ccc;
}
</style>