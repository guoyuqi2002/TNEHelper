<script>
import { ref } from 'vue';

export default {
  setup() {
    const taxFilePath = ref('');
    const isProcessing = ref(false);

    // select a file
    const selectFile = () => {
      // in the real project this is replaced by actual file-picking logic
      taxFilePath.value = "tax-system-file_202601.xlsx";
    };

    // start the import
    const startImport = () => {
      isProcessing.value = true;
      // simulate an import delay
      setTimeout(() => {
        isProcessing.value = false;
        ElMessage.success('Tax system file imported successfully!');
      }, 1500);
    };

    return {
      taxFilePath,
      isProcessing,
      selectFile,
      startImport
    };
  }
};
</script>

<template>
  <div class="tab-content">
    <h2 class="tab-title">Tax System File Upload</h2>
    <div class="upload-section">
      <div class="step">
        <div class="step-label">
          <span class="step-number">1</span>
          <span>Upload the tax system file</span>
        </div>
        <div class="step-action">
          <p v-if="taxFilePath" class="file-name">Selected: {{ taxFilePath }}</p>
          <button
            class="btn-select"
            @click="selectFile"
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
      :disabled="!taxFilePath || isProcessing"
    >
      {{ isProcessing ? 'Importing...' : 'Start Import' }}
    </button>
  </div>
</template>

<style scoped>
.tab-content {
  max-width: 600px;
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
</style>
