<template>
   <h2>Not Found</h2>
   <div class="notfound">
      <div>Sorry, the item you requested could not be found!</div>
      <DPGButton label="Return to Home Page"  @click="homeClicked()" />
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
   min-height: 600px;
   font-size: 1.25em;
   margin: 0;
   padding-top: 50px;
   .p-button {
      margin-top: 25px;
      font-size: 0.8em;
   }

}
</style>
