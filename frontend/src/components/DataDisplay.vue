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
.spacer.label {
   margin-top: 15px;
   font-size: 1.15em;
   color: #aaa;
   margin-bottom: 5px;
   font-weight: normal;
}
</style>