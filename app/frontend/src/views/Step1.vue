<template>
  <div class="step1-container">
    <Header />
    <NavigationBar :step="1" />
    
    <div class="content-wrapper">
      <div class="content-card">
        <div class="card-layout">
          <!-- left: upload file -->
          <div class="left-section">
            <h2 class="section-title">Upload File</h2>
            <div class="new-task-form">
              <!-- file selection -->
              <div class="form-group">
                <div v-if="!taskStore.fileName" class="file-placeholder">
                  <img src="../assets/images/Excel_icon.png" alt="Excel" class="placeholder-icon" />
                  <p class="placeholder-text">No file selected</p>
                </div>
                <div v-else class="file-selected">
                  <img src="../assets/images/Excel_icon.png" alt="Excel" class="file-icon" />
                  <div class="file-details">
                    <p class="file-name">{{ taskStore.fileName }}</p>
                    <p class="file-path" :title="taskStore.filePath">{{ truncatePath(taskStore.filePath, 30) }}</p>
                  </div>
                  <button class="btn-remove" @click="handleFileRemove" title="Remove file">×</button>
                </div>
                <button class="btn-select-file" @click="handleFileSelect">
                  Select File
                </button>
              </div>

              <!-- task naming -->
              <div class="form-group">
                <label for="task-name" class="form-label">Export Folder Name</label>
                <div class="input-wrapper">
                  <input
                    id="task-name"
                    ref="taskNameInput"
                    v-model="localTaskName"
                    type="text"
                    class="task-input"
                    placeholder="example: YYYY-MM"
                    @input="handleTaskNameInput"
                    @focus="showHistoryDropdown = true"
                    @keydown="handleKeyDown"
                    @blur="handleInputBlur"
                  />
                  <!-- history dropdown -->
                  <div
                    v-if="showHistoryDropdown && filteredHistory.length > 0"
                    class="history-dropdown"
                  >
                    <div
                      v-for="(item, index) in filteredHistory"
                      :key="index"
                      class="history-item"
                      :class="{ 'selected': selectedHistoryIndex === index }"
                      @mousedown.prevent="selectHistoryItem(item)"
                      @mouseenter="selectedHistoryIndex = index"
                    >
                      {{ item }}
                    </div>
                  </div>
                </div>
                <p v-if="taskNameError" class="error-text">{{ taskNameError }}</p>
              </div>

              <!-- export path selection -->
              <div class="form-group">
                <label class="form-label">
                  Export Path
                  <span class="required-mark">*</span>
                </label>
                <button class="btn-select-path" @click="handleSelectPath" title="Choose the export path (required)">
                  <span class="folder-icon">📁</span>
                  <span class="btn-text">Select Export Path</span>
                </button>
                <div
                  v-if="taskStore.exportPath"
                  class="export-path-display"
                  :title="taskStore.exportPath"
                  @click="handleSelectPath"
                >
                  <span class="path-icon">📂</span>
                  <span class="path-text">{{ truncatePath(taskStore.exportPath, 35) }}</span>
                  <span class="path-edit-hint">Click to change</span>
                </div>
                <p v-if="!taskStore.exportPath" class="hint-text">Please select an export path for the deduplicated files and match results</p>
              </div>
            </div>
          </div>
          
          <!-- right: process downloaded files -->
          <div class="right-section">
            <div class="section-header">
              <h2 class="section-title">Process Downloaded Files</h2>
              <button class="btn-refresh" @click="loadExistingTasks" title="Refresh task list">
                <span class="refresh-icon">🔄</span>
              </button>
            </div>
            <div class="existing-tasks-container">
              <div v-if="loadingTasks" class="loading-state">
                <div class="loading-spinner">⏳</div>
                <p>Loading...</p>
              </div>
              <div v-else-if="existingTasks.length === 0" class="empty-state">
                <div class="empty-icon">📂</div>
                <p class="empty-text">No downloaded tasks yet</p>
              </div>
              <div v-else class="tasks-list">
                <div
                  v-for="task in existingTasks"
                  :key="task"
                  class="task-item"
                  :class="{ 'selected': selectedExistingTask === task }"
                  @click="handleSelectExistingTask(task)"
                >
                  <div class="task-icon">📁</div>
                  <div class="task-name">{{ task }}</div>
                  <div v-if="selectedExistingTask === task" class="task-check">✓</div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- footer: next button and invoice management button -->
        <div class="footer-section">
          <button class="btn-invoice-mgmt" @click="handleGoToInvoiceMgmt" title="Invoice management">
            <span class="invoice-icon">📄</span>
            <span class="btn-text">Invoice Management</span>
          </button>
          <button v-if="canProceed" class="btn-next" @click="handleNext">
            <span class="arrow-icon">→</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useTaskStore } from '../stores/taskStore'
import Header from '../components/Header.vue'
import NavigationBar from '../components/NavigationBar.vue'

const taskStore = useTaskStore()

// local state
const localTaskName = ref(taskStore.taskName)
const taskNameError = ref('')
const showHistoryDropdown = ref(false)
const selectedHistoryIndex = ref(-1)
const taskNameInput = ref(null)

// existing-task state
const existingTasks = ref([])
const loadingTasks = ref(false)
const selectedExistingTask = ref('')

// computed properties
const filteredHistory = computed(() => {
  if (!localTaskName.value.trim()) {
    return taskStore.taskNameHistory.slice(0, 3)
  }
  const searchText = localTaskName.value.toLowerCase()
  return taskStore.taskNameHistory
    .filter(item => item.toLowerCase().includes(searchText))
    .slice(0, 3)
})

const canProceed = computed(() => {
  // the export path is required (for both new and existing tasks)
  if (!taskStore.exportPath) {
    return false
  }
  
  // if an existing task is selected, only the export path is needed
  if (selectedExistingTask.value) {
    return true
  }
  
  // a new task requires all fields
  return taskStore.filePath && 
         taskStore.fileName && 
         localTaskName.value.trim() !== '' &&
         !taskNameError.value
})

// methods
async function loadExistingTasks() {
  try {
    loadingTasks.value = true
    const tasks = await window.go?.main?.App?.Tasks?.()
    if (tasks) {
      existingTasks.value = tasks
    }
  } catch (error) {
    console.error('failed to load task list:', error)
    alert('Failed to load the task list')
  } finally {
    loadingTasks.value = false
  }
}

function handleSelectExistingTask(taskName) {
  if (selectedExistingTask.value === taskName) {
    // deselect
    selectedExistingTask.value = ''
    taskStore.setOpenExistingTask(false)
  } else {
    // select the existing task
    selectedExistingTask.value = taskName
    taskStore.setFile('', '')
    localTaskName.value = ''
    taskStore.setTaskName(taskName)  // set the task name to the selected task
    // keep the export path so the user can reuse the previously selected one
    taskStore.setOpenExistingTask(true)
  }
}

async function handleFileSelect() {
  try {
    const path = await window.go?.main?.App?.SelectFile?.('Select an Excel file')
    if (path) {
      const name = path.split(/[/\\]/).pop()
      taskStore.setFile(path, name)
      selectedExistingTask.value = ''
      taskStore.setOpenExistingTask(false)
    }
  } catch (error) {
    console.error('failed to select file:', error)
    alert('Failed to select a file; make sure the backend API bindings are generated')
  }
}

function handleFileRemove() {
  taskStore.setFile('', '')
}

function handleTaskNameInput() {
  const value = localTaskName.value.trim()
  
  if (value === '') {
    taskNameError.value = 'Task name cannot be empty'
  } else if (value.length > 50) {
    taskNameError.value = 'Task name cannot exceed 50 characters'
  } else {
    taskNameError.value = ''
    taskStore.setTaskName(value)
  }
  
  if (value) {
    selectedExistingTask.value = ''
  }
  
  if (value && filteredHistory.value.length > 0) {
    showHistoryDropdown.value = true
  }
  selectedHistoryIndex.value = -1
}

function handleInputBlur() {
  setTimeout(() => {
    showHistoryDropdown.value = false
    selectedHistoryIndex.value = -1
  }, 200)
}

function handleKeyDown(event) {
  if (!showHistoryDropdown.value || filteredHistory.value.length === 0) {
    return
  }

  if (event.key === 'ArrowDown') {
    event.preventDefault()
    selectedHistoryIndex.value = Math.min(
      selectedHistoryIndex.value + 1,
      filteredHistory.value.length - 1
    )
  } else if (event.key === 'ArrowUp') {
    event.preventDefault()
    selectedHistoryIndex.value = Math.max(selectedHistoryIndex.value - 1, -1)
  } else if (event.key === 'Enter' && selectedHistoryIndex.value >= 0) {
    event.preventDefault()
    selectHistoryItem(filteredHistory.value[selectedHistoryIndex.value])
  } else if (event.key === 'Escape') {
    showHistoryDropdown.value = false
    selectedHistoryIndex.value = -1
  }
}

function selectHistoryItem(item) {
  localTaskName.value = item
  handleTaskNameInput()
  showHistoryDropdown.value = false
  selectedHistoryIndex.value = -1
  if (taskNameInput.value) {
    taskNameInput.value.focus()
  }
}

async function handleSelectPath() {
  try {
    const path = await window.go?.main?.App?.SelectDirectory?.('Choose where to save the export results')
    if (path) {
      taskStore.setExportPath(path)
      // keep the selected task so the user can still use an existing task after choosing the export path
    }
  } catch (error) {
    console.error('failed to select directory:', error)
    alert('Failed to select a directory; make sure the backend API bindings are generated')
  }
}

function truncatePath(path, maxLength = 40) {
  if (!path) return ''
  if (path.length <= maxLength) return path
  return '...' + path.slice(-maxLength)
}

function handleNext() {
  if (!canProceed.value) {
    return
  }
  
  if (selectedExistingTask.value) {
    taskStore.setTaskName(selectedExistingTask.value)
    taskStore.setFile('', '')
  }
  
  taskStore.setStep(2)
  
  if (!selectedExistingTask.value) {
    taskStore.addTaskNameToHistory(taskStore.taskName)
  }
}

function handleGoToInvoiceMgmt() {
  taskStore.setStep(3)
}

onMounted(() => {
  localTaskName.value = taskStore.taskName
  taskStore.loadTaskNameHistory()
  loadExistingTasks()
})
</script>

<style scoped>
.step1-container {
  height: 100vh;
  background-color: #F8F9FA;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content-wrapper {
  flex: 1;
  padding: 30px 20px;
  overflow: auto;
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.content-card {
  max-width: 1200px;
  width: 100%;
  background-color: #FFFFFF;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.card-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  margin-bottom: 30px;
  align-items: start;
}

.left-section,
.right-section {
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #2C3E50;
  margin: 0;
}

.btn-refresh {
  background: none;
  border: 1px solid #D0D0D0;
  border-radius: 6px;
  padding: 6px 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 16px;
}

.btn-refresh:hover {
  border-color: #2C8CFF;
  background-color: #F0F7FF;
  transform: rotate(90deg);
}

/* new task form */
.new-task-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #2C3E50;
}

.required-mark {
  color: #DC3545;
  margin-left: 4px;
  font-weight: 600;
}

.hint-text {
  margin-top: 8px;
  font-size: 12px;
  color: #DC3545;
  display: flex;
  align-items: center;
  gap: 4px;
}

/* file selection area */
.file-placeholder {
  border: 2px dashed #D0D0D0;
  border-radius: 8px;
  padding: 30px 20px;
  background-color: #F8F9FA;
  text-align: center;
  transition: all 0.3s ease;
}

.placeholder-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 8px;
  display: block;
}

.placeholder-text {
  font-size: 14px;
  color: #6C757D;
  margin: 0;
}

.file-selected {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background-color: #F0FDF4;
  border: 1px solid #50C878;
  border-radius: 8px;
}

.file-icon {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

.file-details {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: #2C3E50;
  margin: 0 0 4px 0;
  word-break: break-all;
}

.file-path {
  font-size: 12px;
  color: #6C757D;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.btn-remove {
  background: none;
  border: none;
  font-size: 24px;
  color: #DC3545;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s ease;
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-remove:hover {
  background-color: #FFE5E5;
  transform: scale(1.1);
}

.btn-select-file {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  background-color: #4A90E2;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-select-file:hover {
  background-color: #357ABD;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(74, 144, 226, 0.3);
}

/* task name input */
.input-wrapper {
  position: relative;
}

.task-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #D0D0D0;
  border-radius: 6px;
  font-size: 14px;
  color: #333333;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.task-input:focus {
  outline: none;
  border-color: #2C8CFF;
  box-shadow: 0 0 0 3px rgba(44, 140, 255, 0.1);
}

.task-input::placeholder {
  color: #999;
}

/* history dropdown */
.history-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background-color: #F5F7FA;
  border: 1px solid #D0D0D0;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  max-height: 150px;
  overflow-y: auto;
}

.history-item {
  padding: 10px 16px;
  font-size: 14px;
  color: #333333;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.history-item:hover,
.history-item.selected {
  background-color: #E8EBF0;
}

.error-text {
  font-size: 12px;
  color: #DC3545;
  margin: 0;
}

/* path selection button */
.btn-select-path {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px 16px;
  background-color: #2C8CFF;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 100%;
}

.btn-select-path:hover {
  background-color: #1A7AE6;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(44, 140, 255, 0.3);
}

/* export path display */
.export-path-display {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background-color: #F5F7FA;
  border: 1px solid #E0E0E0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.export-path-display:hover {
  background-color: #E8EBF0;
  border-color: #2C8CFF;
}

.path-text {
  flex: 1;
  font-size: 13px;
  color: #666666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.path-edit-hint {
  font-size: 12px;
  color: #2C8CFF;
  flex-shrink: 0;
}

/* existing task list */
.existing-tasks-container {
  border: 1px solid #E0E0E0;
  border-radius: 8px;
  background-color: #F8F9FA;
  padding: 16px;
  height: 360px;
  overflow-y: auto;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #6C757D;
}

.loading-spinner {
  font-size: 32px;
  margin-bottom: 12px;
  animation: spin 2s linear infinite;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
  opacity: 0.5;
}

.empty-text {
  font-size: 14px;
  margin: 0;
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background-color: #FFFFFF;
  border: 2px solid #E0E0E0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.task-item:hover {
  border-color: #2C8CFF;
  background-color: #F0F7FF;
  transform: translateX(4px);
}

.task-item.selected {
  border-color: #4A90E2;
  background-color: #E8F4FD;
  box-shadow: 0 2px 8px rgba(74, 144, 226, 0.2);
}

.task-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.task-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  color: #2C3E50;
  word-break: break-all;
}

.task-check {
  font-size: 18px;
  color: #4A90E2;
  flex-shrink: 0;
}

/* footer buttons */
.footer-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20px;
  border-top: 1px solid #E0E0E0;
}

.btn-next {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-color: #4A90E2;
  color: #FFFFFF;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(74, 144, 226, 0.3);
}

.btn-next:hover {
  background-color: #357ABD;
  transform: translateX(5px) scale(1.05);
  box-shadow: 0 4px 12px rgba(74, 144, 226, 0.4);
}

.btn-invoice-mgmt {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background-color: #6C757D;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(108, 117, 125, 0.2);
}

.btn-invoice-mgmt:hover {
  background-color: #5A6268;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(108, 117, 125, 0.3);
}
</style>
