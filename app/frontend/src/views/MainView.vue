<template>
  <div class="main-view">
    <Transition name="page" mode="out-in">
      <Step1 v-if="taskStore.currentStep === 1" key="step1" />
      <Step2 v-else-if="taskStore.currentStep === 2" key="step2" />
      <InvoiceMgmt v-else-if="taskStore.currentStep === 3" key="invoice" />
    </Transition>
  </div>
</template>

<script setup>
import { useTaskStore } from '../stores/taskStore'
import Step1 from './Step1.vue'
import Step2 from './Step2.vue'
import InvoiceMgmt from './InvoiceMgmt.vue'

const taskStore = useTaskStore()
</script>

<style scoped>
.main-view {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.page-enter-active,
.page-leave-active {
  transition: all 0.3s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}
</style>
