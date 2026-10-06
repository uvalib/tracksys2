<template>
   <div class="acts results-toolbar" style="padding: 5px 0 15px 0;">
      <div class="acts">
         <DPGButton label="Download Results CSV" severity="secondary" @click="emit('csv')" v-if="total>0" />
         <DPGButton v-if="hasFilter" label="Clear All Filters" severity="secondary" @click="emit('clear-filter')"/>
         <DPGButton v-if="count < total" 
            label="Load More Results" @click="emit('more')"/>
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
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   align-items: center;
   gap: 10px;
}
</style>