<template>
   <div class="signedout">
      <template v-if="$route.query.expired">
         <h2>Session Expired</h2>
         <div class="message">
            <p>Your session has expired and you have been signed out.</p>
         </div>
      </template>
      <template v-else>
         <h2>Signed Out</h2>
         <div class="message">
            <p>You have successfully signed out of your account.</p>
         </div>
      </template>
      <div class="buttons"><DPGButton @click="signinClick">Sign in with NetBadge</DPGButton></div>
   </div>
</template>

<script setup>
function signinClick() {
   window.location.href = "/authenticate"
}
</script>

<style scoped lang="scss">
.signedout {
   padding: 0;
   min-height: 600px;
   .message {
      margin-top: 40px;
      font-size: 1.15em;
   }
   .buttons {
      margin-top: 35px;
      font-size: 1.25em;
   }
}
</style>
