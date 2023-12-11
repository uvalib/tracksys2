<template>
   <h2>ArchivesSpace Reviews</h2>
   <div class="reviews">
   </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { usePinnable } from '@/composables/pin'

usePinnable("p-paginator-top")

onMounted(() => {
})

</script>

<style scoped lang="scss">
.reviews {
   min-height: 600px;
   text-align: left;
   padding: 0 25px;
}
</style>