<template>
   <h2>Not Found</h2>
   <div class="notfound">
      <div>Sorry, the item you requested could not be found!</div>
      <UButton label="Return to Home Page" @click="homeClicked" />
   </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
function homeClicked() {
   router.push("/")
}
</script>

<style scoped lang="scss">
.notfound {
   font-size: 1.25em;
   padding: 25px;
   align-items: center; 
   display: flex;
   flex-direction: column;
   gap: 20px;
}
</style>
