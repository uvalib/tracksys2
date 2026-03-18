<template>
   <div class="form-field">
      <label :for="props.id">
         <span>{{ props.label }}</span>
         <span v-if="props.required">*<span class="required">(required)</span></span>
      </label>
      <slot></slot>
      <Message v-if="props.required && props.error" severity="error" size="small" variant="simple">{{ props.error }}</Message>
   </div>
</template>

<script setup>
import Message from 'primevue/message'

const props = defineProps({
   id: {
      type: String,
      required: true,
   },
   required: {
      type: Boolean,
      default: false
   },
   label: {
      type: String,
      required: true,
   },
   error: {
      type: String,
      default: ""
   },
})
</script>

<style scoped lang="scss">
.form-field {
   display: flex;
   flex-direction: column;
   gap: 5px;
}
.required {
   margin-left: 5px;
   color: var(--uvalib-grey-dark);
   font-size: .9em;
}
</style>