<template>
  <div v-if="visible" class="initializing-overlay">
    <div class="initializing-content">
      <div class="spinner"></div>
      <p class="initializing-text">{{ text }}</p>
    </div>
  </div>
</template>

<script setup>
defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  text: {
    type: String,
    default: 'Loading file, please wait...'
  }
})
</script>

<style scoped>
.initializing-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  z-index: 1000;
}

.initializing-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #E0E0E0;
  border-top-color: #2C8CFF;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.initializing-text {
  font-size: 16px;
  color: #2C3E50;
  font-weight: 500;
}
</style>
