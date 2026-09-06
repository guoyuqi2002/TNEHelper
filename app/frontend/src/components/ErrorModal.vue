<template>
  <div v-if="visible" class="modal-overlay" @click.self="handleClose">
    <div class="modal-container">
      <div class="modal-header">
        <h3 class="modal-title">Error</h3>
        <button class="close-btn" @click="handleClose">×</button>
      </div>
      
      <div class="modal-body">
        <div class="error-icon">⚠️</div>
        <div class="error-message">
          <p v-if="errorType === 'download'">
            Download of {{ fileName }} was interrupted; click Run to continue
          </p>
          <p v-else-if="errorType === 'deduplicate'">
            Deduplication was interrupted; click Run to continue
          </p>
          <p v-else-if="errorType === 'recognize'">
            Recognition was interrupted; click Run to continue
          </p>
          <p v-else-if="errorType === 'prerequisite'">
            {{ errorMessage }}
          </p>
          <pre v-if="errorLog" class="error-log">{{ errorLog }}</pre>
          <p v-else-if="errorType !== 'prerequisite'">
            {{ errorMessage }}
          </p>
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="btn-confirm" @click="handleConfirm">OK</button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  errorType: {
    type: String,
    default: 'download'
  },
  errorMessage: {
    type: String,
    default: ''
  },
  errorLog: {
    type: String,
    default: ''
  },
  fileName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['confirm', 'close'])

function handleConfirm() {
  emit('confirm', props.errorType)
  handleClose()
}

function handleClose() {
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.modal-container {
  background-color: #FFFFFF;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  max-width: 500px;
  width: 90%;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  animation: slideUp 0.3s ease;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #E0E0E0;
}

.modal-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #2C3E50;
}

.close-btn {
  background: none;
  border: none;
  font-size: 28px;
  color: #999;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background-color: #F5F5F5;
  color: #333;
}

.modal-body {
  padding: 30px 24px;
  text-align: center;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-message {
  color: #2C3E50;
  font-size: 16px;
  line-height: 1.6;
}

.error-message p {
  margin: 0;
}

.error-log {
  margin-top: 16px;
  padding: 12px;
  background-color: #F5F5F5;
  border: 1px solid #E0E0E0;
  border-radius: 4px;
  text-align: left;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #D32F2F;
  white-space: pre-wrap;
  word-wrap: break-word;
  max-height: 300px;
  overflow-y: auto;
  line-height: 1.5;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #E0E0E0;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-confirm {
  padding: 10px 24px;
  background-color: #4A90E2;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.btn-confirm:hover {
  background-color: #357ABD;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(74, 144, 226, 0.3);
}
</style>
