<template>
  <div class="progress-table-container">
    <table class="progress-table">
      <thead>
        <tr>
          <th class="col-process">Step</th>
          <th class="col-progress">Progress</th>
          <th class="col-action">Action</th>
        </tr>
      </thead>
      <tbody>
        <!-- task 1: download -->
        <tr>
          <td class="process-name">Download</td>
          <td class="progress-cell">
            <div class="progress-wrapper">
              <div 
                class="progress-bar" 
                :style="{ width: progress1 + '%' }"
                :class="{
                  'progress-running': isTask1Running,
                  'progress-completed': progress1 === 100 && !isTask1Running,
                  'progress-error': hasError1
                }"
              ></div>
              <span class="progress-text">{{ progressText1 }}</span>
            </div>
          </td>
          <td class="action-cell">
            <button 
              class="btn-run"
              :class="{
                'btn-running': isTask1Running,
                'btn-error': hasError1,
                'btn-disabled': isOpenExistingTask
              }"
              :disabled="isTask1Running || isOpenExistingTask"
              @click="emit('run', 1)"
              :title="isOpenExistingTask ? 'existing task, no download needed' : ''"
            >
              {{ isTask1Running ? 'Running...' : isOpenExistingTask ? 'Done' : 'Run' }}
            </button>
          </td>
        </tr>
        
        <!-- task 2: deduplicate -->
        <tr>
          <td class="process-name">Deduplicate</td>
          <td class="progress-cell">
            <div class="progress-wrapper">
              <div 
                class="progress-bar" 
                :style="{ width: progress2 + '%' }"
                :class="{
                  'progress-running': isTask2Running,
                  'progress-completed': progress2 === 100 && !isTask2Running,
                  'progress-error': hasError2
                }"
              ></div>
              <span class="progress-text">{{ progressText2 }}</span>
            </div>
          </td>
          <td class="action-cell">
            <button 
              class="btn-run"
              :class="{
                'btn-running': isTask2Running,
                'btn-error': hasError2,
                'btn-rerun': progress2 === 100 && !isTask2Running
              }"
              :disabled="isTask2Running"
              @click="emit('run', 2)"
              :title="progress2 === 100 ? 'click to re-run the export' : 'click to start deduplication'"
            >
              {{ isTask2Running ? 'Running...' : progress2 === 100 ? 'Re-export Files' : 'Run' }}
            </button>
          </td>
        </tr>
        
        <!-- task 3: recognize -->
        <tr>
          <td class="process-name">Recognize</td>
          <td class="progress-cell">
            <div class="progress-wrapper">
              <div 
                class="progress-bar" 
                :style="{ width: progress3 + '%' }"
                :class="{
                  'progress-running': isTask3Running,
                  'progress-completed': progress3 === 100 && !isTask3Running,
                  'progress-error': hasError3
                }"
              ></div>
              <span class="progress-text">{{ progressText3 }}</span>
            </div>
          </td>
          <td class="action-cell">
            <button 
              class="btn-run"
              :class="{
                'btn-running': isTask3Running,
                'btn-error': hasError3,
                'btn-rerun': progress3 === 100 && !isTask3Running
              }"
              :disabled="isTask3Running"
              @click="emit('run', 3)"
              :title="progress3 === 100 ? 'click to re-run the export' : 'click to start recognition'"
            >
              {{ isTask3Running ? 'Running...' : progress3 === 100 ? 'Re-export Match Results' : 'Run' }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
defineProps({
  progress1: { type: Number, default: 0 },
  progress2: { type: Number, default: 0 },
  progress3: { type: Number, default: 0 },
  progressText1: { type: String, default: '0/0' },
  progressText2: { type: String, default: '0/0' },
  progressText3: { type: String, default: '0/0' },
  isTask1Running: { type: Boolean, default: false },
  isTask2Running: { type: Boolean, default: false },
  isTask3Running: { type: Boolean, default: false },
  hasError1: { type: Boolean, default: false },
  hasError2: { type: Boolean, default: false },
  hasError3: { type: Boolean, default: false },
  isOpenExistingTask: { type: Boolean, default: false }
})

const emit = defineEmits(['run'])
</script>

<style scoped>
.progress-table-container {
  background-color: #FFFFFF;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.progress-table {
  width: 100%;
  border-collapse: collapse;
}

.progress-table thead {
  background-color: #F8F9FA;
}

.progress-table th {
  padding: 16px 20px;
  text-align: left;
  font-size: 14px;
  font-weight: 600;
  color: #2C3E50;
  border-bottom: 2px solid #E0E0E0;
}

.col-process {
  width: 150px;
}

.col-progress {
  width: auto;
  min-width: 300px;
}

.col-action {
  width: 120px;
  text-align: center;
}

.progress-table tbody tr {
  border-bottom: 1px solid #E0E0E0;
  transition: background-color 0.2s ease;
}

.progress-table tbody tr:hover {
  background-color: #F8F9FA;
}

.progress-table tbody tr:last-child {
  border-bottom: none;
}

.process-name {
  padding: 20px;
  font-size: 16px;
  font-weight: 500;
  color: #2C3E50;
}

.progress-cell {
  padding: 20px;
}

.progress-wrapper {
  position: relative;
  width: 100%;
  height: 32px;
  background-color: #E9ECEF;
  border-radius: 16px;
  overflow: hidden;
}

.progress-bar {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background-color: #4A90E2;
  border-radius: 16px;
  transition: width 0.5s ease, background-color 0.3s ease;
  z-index: 1;
}

.progress-bar.progress-running {
  background-color: #4A90E2;
  animation: pulse 2s ease-in-out infinite;
}

.progress-bar.progress-completed {
  background-color: #50C878;
}

.progress-bar.progress-error {
  background-color: #DC3545;
}

.progress-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 14px;
  font-weight: 600;
  color: #2C3E50;
  z-index: 2;
  pointer-events: none;
}

.action-cell {
  padding: 20px;
  text-align: center;
}

.btn-run {
  padding: 8px 20px;
  background-color: #4A90E2;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 80px;
}

.btn-run:hover:not(:disabled) {
  background-color: #357ABD;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(74, 144, 226, 0.3);
}

.btn-run:active:not(:disabled) {
  transform: translateY(0);
}

.btn-run:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-run.btn-running {
  background-color: #6C757D;
}

.btn-run.btn-error {
  background-color: #DC3545;
}

.btn-run.btn-error:hover:not(:disabled) {
  background-color: #C82333;
  box-shadow: 0 2px 8px rgba(220, 53, 69, 0.3);
}

.btn-run.btn-disabled {
  background-color: #50C878;
  opacity: 0.8;
}

.btn-run.btn-rerun {
  background-color: #10B981;
  min-width: 140px;
}

.btn-run.btn-rerun:hover:not(:disabled) {
  background-color: #059669;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.3);
}
</style>
