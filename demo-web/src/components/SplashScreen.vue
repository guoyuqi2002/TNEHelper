<script>
export default {
  setup() {
    return {};
  }
};
</script>

<template>
  <div class="splash-screen">
    <svg class="splash-logo" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="3" y="5" width="18" height="14" rx="2" stroke="currentColor" stroke-width="2" />
      <path d="M7 9.5h4M7 13h10" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
      <circle cx="16.5" cy="9.5" r="1.5" fill="currentColor" />
    </svg>
    <h2 class="splash-title">Expense Reimbursement Checker</h2>
    <div class="splash-loading"></div>
  </div>
</template>

<style scoped>
.splash-screen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #fff;
  z-index: 9999;
}
.splash-logo {
  width: 120px;
  height: 120px;
  color: #28a745;
  margin-bottom: 24px;
  animation: pulse 2s infinite ease-in-out;
}
.splash-title {
  font-size: 24px;
  color: #333;
  font-weight: 600;
}
.splash-loading {
  margin-top: 20px;
  border: 8px solid #f3f3f3;
  border-top: 8px solid #28a745;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 2s linear infinite;
}
@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
}
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
