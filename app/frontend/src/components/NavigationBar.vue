<template>
  <nav class="navigation-bar">
    <div class="nav-container">
      <div 
        class="nav-item" 
        :class="{ active: step >= 1, completed: step > 1 }"
      >
        <div class="step-number">1</div>
        <div class="step-label">File Upload &amp; Naming</div>
      </div>
      
      <div class="nav-divider"></div>
      
      <div 
        class="nav-item" 
        :class="{ active: step >= 2, completed: step > 2 }"
      >
        <div class="step-number">2</div>
        <div class="step-label">Task Progress Monitoring</div>
      </div>
    </div>
  </nav>
</template>

<script setup>
defineProps({
  step: {
    type: Number,
    required: true,
    validator: (value) => [1, 2].includes(value)
  }
})
</script>

<style scoped>
.navigation-bar {
  width: 100%;
  background-color: #4A90E2;
  padding: 20px 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 40px;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  opacity: 0.5;
  transition: all 0.3s ease;
}

.nav-item.active {
  opacity: 1;
}

.nav-item.completed .step-number {
  background-color: #50C878;
}

.step-number {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.3);
  color: #FFFFFF;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  transition: all 0.3s ease;
}

.nav-item.active .step-number {
  background-color: #FFFFFF;
  color: #4A90E2;
  transform: scale(1.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.step-label {
  color: #FFFFFF;
  font-size: 14px;
  font-weight: 500;
  text-align: center;
  white-space: nowrap;
}

.nav-divider {
  width: 60px;
  height: 2px;
  background-color: rgba(255, 255, 255, 0.3);
  margin: 0 10px;
}
</style>
