<script>
import { ref, watch } from 'vue';

export default {
  props: {
    isImported: {
      type: Boolean,
      default: false
    }
  },
  emits: ['import-success'],
  setup(props, { emit }) {
    const reimburseFilePath = ref('');
    const isProcessing = ref(false);
    const savePath = ref('');

    // select a reimbursement file
    const selectReimburseFile = () => {
      reimburseFilePath.value = "reimbursement-summary_202601.xlsx";
    };

    // start the import
    const startImport = () => {
      isProcessing.value = true;
      setTimeout(() => {
        isProcessing.value = false;
        ElMessage.success('Reimbursement file imported successfully!');
        emit('import-success');
      }, 1500);
    };

    // select the save path
    const selectSavePath = () => {
      savePath.value = "D:/reimburse-vouchers/";
      ElMessage.success(`Save path selected: ${savePath.value}`);
    };

    // check against the tax system
    const checkWithTaxSystem = () => {
      ElMessage.info('Checking against the tax system... (simulated)');
    };

    return {
      reimburseFilePath,
      isProcessing,
      savePath,
      selectReimburseFile,
      startImport,
      selectSavePath,
      checkWithTaxSystem
    };
  }
};
</script>

<template>
  <div class="tab-content">
    <!-- pre-import section -->
    <div v-if="!isImported" class="import-section">
      <h2 class="tab-title">Reimbursement File Check</h2>
      <div class="upload-section">
        <div class="step">
          <div class="step-label">
            <span class="step-number">1</span>
            <span>Upload the reimbursement summary sheet</span>
          </div>
          <div class="step-action">
            <p v-if="reimburseFilePath" class="file-name">Selected: {{ reimburseFilePath }}</p>
            <button
              class="btn-select"
              @click="selectReimburseFile"
              :disabled="isProcessing"
            >
              Select File
            </button>
          </div>
        </div>
      </div>
      <button
        class="btn-import"
        @click="startImport"
        :disabled="!reimburseFilePath || isProcessing"
      >
        {{ isProcessing ? 'Importing...' : 'Start Import' }}
      </button>
    </div>

    <!-- post-import section (with slide-in animation) -->
    <div v-else class="post-import-section slide-in">
      <div class="action-row">
        <!-- download reimbursement summary vouchers -->
        <div class="action-card">
          <h3>1. Download Reimbursement Vouchers</h3>
          <button class="btn-action" @click="selectSavePath">
            Select Save Path
          </button>
          <p v-if="savePath" class="path-text">Current path: {{ savePath }}</p>
        </div>

        <!-- reimbursement information check -->
        <div class="action-card">
          <h3>2. Reimbursement Check</h3>
          <button class="btn-action" @click="checkWithTaxSystem">
            Check Against Tax System
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tab-content {
  max-width: 800px;
  margin: 0 auto;
}
.tab-title {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
}
.upload-section {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 2rem;
}
.step {
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
}
.step-label {
  display: flex;
  align-items: center;
  font-weight: 500;
}
.step-number {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: #28a745;
  color: #fff;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 10px;
}
.step-action {
  display: flex;
  align-items: center;
}
.file-name {
  margin: 0 15px 0 0;
  font-style: italic;
  color: #666;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.btn-select {
  padding: 8px 16px;
  border: 1px solid #28a745;
  background-color: #fff;
  color: #28a745;
  border-radius: 6px;
  cursor: pointer;
}
.btn-select:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.btn-import {
  width: 100%;
  padding: 12px;
  border: none;
  border-radius: 8px;
  background-color: #28a745;
  color: #fff;
  font-size: 1rem;
  cursor: pointer;
}
.btn-import:disabled {
  background-color: #8ee49d;
  cursor: not-allowed;
}

/* post-import section styles */
.post-import-section {
  margin-top: 2rem;
}
.slide-in {
  animation: slideIn 0.5s ease-out;
}
@keyframes slideIn {
  from { transform: translateX(50px); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}
.action-row {
  display: flex;
  justify-content: space-between;
  gap: 20px;
}
.action-card {
  flex: 1;
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  text-align: center;
}
.btn-action {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  background-color: #28a745;
  color: #fff;
  cursor: pointer;
  margin: 10px 0;
}
.path-text {
  color: #666;
  font-size: 0.9rem;
  margin: 0;
}
</style>
