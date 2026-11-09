<template>
   <div class="row-between z-5 bg-white border-b border-brand-grey-100 pt-1 pb-3 px-2 items-center">
      <div class="row-between gap-2">
         <UButton label="Download Results CSV" color="secondary" @click="emit('csv')" v-if="total>0" />
         <UButton v-if="hasFilter" label="Clear All Filters" color="secondary" @click="emit('clear-filter')"/>
         <UButton v-if="count < total" label="Load More Results" @click="emit('more')"/>
      </div>
      <div>Showing results 1 - {{ count }}</div>
   </div>
</template>

<script setup>
const emit = defineEmits( ['csv', 'more', 'clear-filter' ])
const props = defineProps({
   total: {
      type: Number,
      required: true,
   },
   count: {
      type: Number,
      required: true,
   },
   hasFilter: {
      type: Boolean,
      required: true,
   },
})
</script>

<style scoped lang="scss">
</style>