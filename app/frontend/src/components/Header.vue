<template>
  <header class="header">
    <div class="header-container">
      <div class="logo-text">
        <svg class="logo-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="3" y="5" width="18" height="14" rx="2" stroke="currentColor" stroke-width="2" />
          <path d="M7 9.5h4M7 13h10" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
          <circle cx="16.5" cy="9.5" r="1.5" fill="currentColor" />
        </svg>
        <span>TNE Assistant</span>
      </div>
    </div>
  </header>
</template>

<script setup>
</script>

<style scoped>
.header {
  width: 100%;
  background-color: #FFFFFF;
  border-bottom: 1px solid #E0E0E0;
  padding: 20px 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 24px;
  font-weight: 600;
  color: #1F2937;
  letter-spacing: 0.5px;
  transition: opacity 0.3s ease;
}

.logo-icon {
  width: 32px;
  height: 32px;
  color: #2563EB;
}

.logo-text:hover {
  opacity: 0.8;
}
</style>
