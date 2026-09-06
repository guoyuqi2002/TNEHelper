<template>
  <div class="step2-container">
    <Header />
    <NavigationBar :step="2" />
    
    <div class="content-wrapper" ref="contentWrapper">
      <div class="content-card">
        <!-- initialization overlay -->
        <InitializingOverlay :visible="taskStore.isInitializing" />
        <!-- exporting overlay -->
        <InitializingOverlay :visible="isExporting" :text="exportingText" />
        
        <!-- progress table -->
        <div class="progress-section" :class="{ 'loading': taskStore.isInitializing }">
          <h2 class="section-title">Task Progress</h2>
          <ProgressTable
            :progress1="taskStore.progress1"
            :progress2="taskStore.progress2"
            :progress3="taskStore.progress3"
            :progress-text1="taskStore.progressText1"
            :progress-text2="taskStore.progressText2"
            :progress-text3="taskStore.progressText3"
            :is-task1-running="taskStore.isTask1Running"
            :is-task2-running="taskStore.isTask2Running"
            :is-task3-running="taskStore.isTask3Running"
            :has-error1="taskStore.errors.download !== null"
            :has-error2="taskStore.errors.deduplicate !== null"
            :has-error3="taskStore.errors.recognize !== null"
            :is-open-existing-task="taskStore.isOpenExistingTask"
            @run="handleRunTask"
          />
        </div>
        
        <!-- log area -->
        <div class="logs-section">
          <div class="logs-header">
            <h3 class="logs-title">Execution Log</h3>
            <div class="logs-buttons">
              <button
                v-if="taskStore.allLogs.length > 0"
                class="btn-save-logs"
                @click="handleSaveLogs"
              >
                View All Attachments
              </button>
            </div>
          </div>

          <div class="logs-box">
            <div v-if="showEmptyHint" class="empty-logs">No logs yet</div>
            <pre v-else class="logs-pre" ref="logPre"></pre>
          </div>
        </div>
      </div>
    </div>
    
    <!-- back-to-top button -->
    <button
      v-if="taskStore.allLogs.length > 10"
      class="btn-scroll-to-top"
      @click="scrollToPageTop"
      title="Scroll to the top of the app"
    >
      Back to Top
    </button>
    
    <!-- error modal -->
    <ErrorModal
      :visible="showErrorModal"
      :error-type="currentErrorType"
      :error-message="currentErrorMessage"
      :error-log="currentErrorLog"
      :file-name="taskStore.fileName"
      @confirm="handleErrorConfirm"
      @close="showErrorModal = false"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useTaskStore } from '../stores/taskStore'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import Header from '../components/Header.vue'
import NavigationBar from '../components/NavigationBar.vue'
import ProgressTable from '../components/ProgressTable.vue'
import ErrorModal from '../components/ErrorModal.vue'
import InitializingOverlay from '../components/InitializingOverlay.vue'

const taskStore = useTaskStore()

// state
const showErrorModal = ref(false)
const currentErrorType = ref('download')
const currentErrorMessage = ref('')
const currentErrorLog = ref('')
const logsContainer = ref(null)
const contentWrapper = ref(null)

// export state
const isExporting = ref(false)
const exportingText = ref('')
const exportProgress = ref({ current: 0, total: 0 })

let progressInterval = null

// write logs directly into the <pre> element 
const showEmptyHint = ref(true)
const logPre = ref(null)

// queued batched flushing
const logQueue = []
let scrollThrottleTimer = null
let uiLineCount = 0
const MAX_UI_LINES = 1500
const FLUSH_INTERVAL_MS = 100

function formatTime() {
  return new Date().toLocaleTimeString('en-GB', { hour12: false })
}

function ensurePreReady() {
  if (!showEmptyHint.value) return Promise.resolve()
  showEmptyHint.value = false
  return nextTick()
}

// enqueueUILog displays text in the log box and flushes it at throttled intervals
function enqueueUILog(line) {
  ensurePreReady().then(() => {
    logQueue.push(line)
    if (!scrollThrottleTimer) {
      scrollThrottleTimer = setTimeout(flushUILogs, FLUSH_INTERVAL_MS)
    }
  })
}

function flushUILogs() {
  scrollThrottleTimer = null
  const el = logPre.value
  if (!el || logQueue.length === 0) return

  const nearBottom = (el.scrollHeight - el.scrollTop - el.clientHeight) < 80

  const batch = logQueue.splice(0, logQueue.length).join('')
  el.appendChild(document.createTextNode(batch))

  for (let i = 0; i < batch.length; i++) {
    if (batch.charCodeAt(i) === 10) uiLineCount++
  }

  if (uiLineCount > MAX_UI_LINES) {
    const text = el.textContent || ''
    const lines = text.split('\n')
    const kept = lines.slice(-MAX_UI_LINES).join('\n')
    el.textContent = kept
    uiLineCount = Math.min(MAX_UI_LINES, kept.split('\n').length)
  }

  if (nearBottom) {
    el.scrollTop = el.scrollHeight
  }
}

function scrollLogsToBottomHard() {
  const el = logPre.value
  if (!el) return
  el.scrollTop = el.scrollHeight
}

// logs for export (only written to allLogs) 
function addLogForExportOnly(message, type = 'info') {
  const ts = formatTime()
  if (typeof taskStore.addLogForExportOnly === 'function') {
    taskStore.addLogForExportOnly(message, type, ts)
  } else {
    taskStore.allLogs.push({ time: ts, message, type })
  }
}

// ====== error state watch ======
watch(
  () => [
    taskStore.errors.download,
    taskStore.errors.deduplicate,
    taskStore.errors.recognize,
    taskStore.errorLogs.download,
    taskStore.errorLogs.deduplicate,
    taskStore.errorLogs.recognize
  ],
  ([downloadError, deduplicateError, recognizeError, downloadLog, deduplicateLog, recognizeLog]) => {
    if (downloadError) {
      currentErrorType.value = 'download'
      currentErrorMessage.value = downloadError
      currentErrorLog.value = downloadLog || ''
      showErrorModal.value = true
    } else if (deduplicateError) {
      currentErrorType.value = 'deduplicate'
      currentErrorMessage.value = deduplicateError
      currentErrorLog.value = deduplicateLog || ''
      showErrorModal.value = true
    } else if (recognizeError) {
      currentErrorType.value = 'recognize'
      currentErrorMessage.value = recognizeError
      currentErrorLog.value = recognizeLog || ''
      showErrorModal.value = true
    }
  },
  { deep: true }
)

// fetch the task state
async function fetchTaskState() {
  try {
    const state = await window.go?.main?.App?.GetTaskState?.()
    if (state) {
      taskStore.updateTaskState(state)
    }
  } catch (error) {
    console.error('failed to get task state:', error)
    enqueueUILog(`[${formatTime()}] failed to get task state: ${error.message || error}\n`)
    addLogForExportOnly(`failed to get task state: ${error.message || error}`, 'error')
  }
}

// start progress polling
function startProgressPolling() {
  if (!progressInterval) {
    progressInterval = setInterval(() => {
      fetchTaskState()
    }, 3000)
  }
}

// stop progress polling
function stopProgressPolling() {
  if (progressInterval) {
    clearInterval(progressInterval)
    progressInterval = null
  }
}

// run a task
async function handleRunTask(taskId) {
  // when opening an existing task, running the download task is not allowed
  if (taskId === 1 && taskStore.isOpenExistingTask) {
    showErrorModal.value = true
    currentErrorType.value = 'prerequisite'
    currentErrorMessage.value = 'The files of this task already exist; no need to download them again. Run the deduplication or recognition task directly.'
    currentErrorLog.value = ''
    return
  }
  
  try {
    startProgressPolling()
    taskStore.setTaskRunning(taskId, true)
    
    if (taskId === 1) {
      taskStore.clearError('download')
      enqueueUILog(`[${formatTime()}] Starting download task...\n`)
      addLogForExportOnly('Starting download task...', 'info')
      await window.go?.main?.App?.RunTask1?.()
      enqueueUILog(`[${formatTime()}] Download task finished\n`)
      addLogForExportOnly('Download task finished', 'info')
    } else if (taskId === 2) {
      taskStore.clearError('deduplicate')
      enqueueUILog(`[${formatTime()}] Starting deduplication task...\n`)
      addLogForExportOnly('Starting deduplication task...', 'info')
      await window.go?.main?.App?.RunTask2?.()
      enqueueUILog(`[${formatTime()}] Deduplication task finished\n`)
      addLogForExportOnly('Deduplication task finished', 'info')
    } else if (taskId === 3) {
      if (taskStore.progress3 === 100) {
        enqueueUILog(`[${formatTime()}] Recognition task already finished, re-exporting match results...\n`)
        addLogForExportOnly('Recognition task already finished, re-exporting match results...', 'info')
      } else {
        taskStore.clearError('recognize')
        enqueueUILog(`[${formatTime()}] Starting recognition task...\n`)
        addLogForExportOnly('Starting recognition task...', 'info')
        await window.go?.main?.App?.RunTask3?.()
        enqueueUILog(`[${formatTime()}] Recognition task finished\n`)
        addLogForExportOnly('Recognition task finished', 'info')
      }
    }
    
    await fetchTaskState()
    
    // auto-export after deduplication finishes
    if (taskId === 2 && taskStore.progress2 === 100) {
      if (!taskStore.exportPath) {
        enqueueUILog(`[${formatTime()}] export path not set, cannot auto-export\n`)
        addLogForExportOnly('export path not set, cannot auto-export', 'error')
      } else if (!taskStore.taskName) {
        enqueueUILog(`[${formatTime()}] task name not set, cannot auto-export\n`)
        addLogForExportOnly('task name not set, cannot auto-export', 'error')
      } else {
        try {
          exportProgress.value = { current: 0, total: 0 }
          exportingText.value = 'Exporting deduplicated attachment files...'
          isExporting.value = true
          
          const fullExportPath = taskStore.exportPath + "/" + taskStore.taskName
          await window.go?.main?.App?.Expand5081?.(fullExportPath)
          
          // hide the overlay and show a success hint
          isExporting.value = false
          exportProgress.value = { current: 0, total: 0 }
          enqueueUILog(`[${formatTime()}] deduplication export complete\n`)
          addLogForExportOnly('deduplication export complete', 'info')
          nextTick(() => scrollLogsToBottomHard())
        } catch (error) {
          isExporting.value = false
          exportProgress.value = { current: 0, total: 0 }
          enqueueUILog(`[${formatTime()}] export failed: ${error.message || error}\n`)
          addLogForExportOnly(`export failed: ${error.message || error}`, 'error')
        }
      }
    }
    
    // auto-export match results after recognition finishes
    if (taskId === 3 && taskStore.progress3 === 100) {
      if (!taskStore.exportPath) {
        enqueueUILog(`[${formatTime()}] export path not set, cannot auto-export match results\n`)
        addLogForExportOnly('export path not set, cannot auto-export match results', 'error')
      } else if (!taskStore.taskName) {
        enqueueUILog(`[${formatTime()}] task name not set, cannot auto-export match results\n`)
        addLogForExportOnly('task name not set, cannot auto-export match results', 'error')
      } else {
        try {
          exportProgress.value = { current: 0, total: 0 }
          exportingText.value = 'Exporting match results Excel...'
          isExporting.value = true
          
          // export into the task sub-directory, consistent with Expand5081
          const fullExportPath = taskStore.exportPath + "/" + taskStore.taskName
          await window.go?.main?.App?.ExportExcel?.(fullExportPath)
          
          // hide the overlay and show a success hint
          isExporting.value = false
          exportProgress.value = { current: 0, total: 0 }
          enqueueUILog(`[${formatTime()}] match results export complete\n`)
          addLogForExportOnly('match results export complete', 'info')
          nextTick(() => scrollLogsToBottomHard())
        } catch (error) {
          isExporting.value = false
          exportProgress.value = { current: 0, total: 0 }
          enqueueUILog(`[${formatTime()}] failed to export match results: ${error.message || error}\n`)
          addLogForExportOnly(`failed to export match results: ${error.message || error}`, 'error')
        }
      }
    }
    
  } catch (error) {
    console.error(`task ${taskId} failed:`, error)
    const errorMsg = error.message || error.toString()
    
    if (taskId === 1) {
      taskStore.setError('download', `Download of ${taskStore.fileName} failed; click Run to continue`, errorMsg)
    } else if (taskId === 2) {
      taskStore.setError('deduplicate', 'Deduplication failed; click Run to continue', errorMsg)
    } else if (taskId === 3) {
      taskStore.setError('recognize', 'Recognition failed; click Run to continue', errorMsg)
    }

    enqueueUILog(`[${formatTime()}] task ${taskId} failed: ${errorMsg}\n`)
    addLogForExportOnly(`task ${taskId} failed: ${errorMsg}`, 'error')
  } finally {
    taskStore.setTaskRunning(taskId, false)
    await fetchTaskState()
    stopProgressPolling()
  }
}

// handle error confirmation
function handleErrorConfirm() {
  showErrorModal.value = false
}

// save logs (export everything)
function handleSaveLogs() {
  try {
    const date = new Date()
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const fileName = `${year}.${month}.${day}_task-log.txt`
    
    // export all logs from allLogs
    const logContent = taskStore.allLogs
      .map(log => `[${log.time}] ${log.message}`)
      .join('\n')
    
    const blob = new Blob([logContent], { type: 'text/plain;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)

    enqueueUILog(`[${formatTime()}] logs saved: ${fileName} (${taskStore.allLogs.length} entries)\n`)
    addLogForExportOnly(`logs saved: ${fileName} (${taskStore.allLogs.length} entries)`, 'info')
  } catch (error) {
    enqueueUILog(`[${formatTime()}] failed to save logs: ${error.message || error}\n`)
    addLogForExportOnly(`failed to save logs: ${error.message || error}`, 'error')
  }
}

// ====== resize fallback: force a repaint (very useful for Wails/WebView2) ======
function forceRepaint() {
  const root = document.getElementById('app')
  if (!root) return
  void document.documentElement.clientWidth
  root.style.transform = 'translateZ(0)'
  requestAnimationFrame(() => {
    root.style.transform = ''
  })
}

function handleResize() {
  requestAnimationFrame(() => {
    requestAnimationFrame(forceRepaint)
  })
}

// scroll the page to the top
function scrollToPageTop() {
  if (contentWrapper.value) {
    contentWrapper.value.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

onMounted(async () => {
  // listen for window resize to fix rendering issues when switching between maximized and windowed mode
  window.addEventListener('resize', handleResize)
  
  // set up event listeners
  EventsOn('process:log', (message) => {
    // export progress message: only updates the overlay text, not the log box
    if (message.includes('__EXPORT_PROGRESS__')) {
      const match = message.match(/__EXPORT_PROGRESS__:(\d+):(\d+):(excel|attachment)/)
      if (match) {
        const current = parseInt(match[1], 10)
        const total = parseInt(match[2], 10)
        const type = match[3]

        exportProgress.value = { current, total }
        exportingText.value =
          type === 'excel'
            ? `Exporting match results Excel... (${current}/${total})`
            : `Exporting deduplicated attachment files... (${current}/${total})`
        return
      }
    }

    const ts = formatTime()
    enqueueUILog(`[${ts}] ${message}\n`)
    addLogForExportOnly(message, 'info')
  })
  
  EventsOn('process:error', (error) => {
    const ts = formatTime()
    enqueueUILog(`[${ts}] ${error}\n`)
    addLogForExportOnly(error, 'error')

    if (error.includes('download')) {
      taskStore.setError('download', `Download of ${taskStore.fileName} failed; click Run to continue`, error)
    } else if (error.includes('dedup') || error.includes('deduplicate')) {
      taskStore.setError('deduplicate', 'Deduplication failed; click Run to continue', error)
    } else if (error.includes('recognition') || error.includes('recognize')) {
      taskStore.setError('recognize', 'Recognition failed; click Run to continue', error)
    }
  })
  
  // initialize the task
  taskStore.setInitializing(true)
  try {
    enqueueUILog(`[${formatTime()}] Initializing task, please wait...\n`)
    addLogForExportOnly('Initializing task, please wait...', 'info')

    await new Promise((resolve) => setTimeout(resolve, 100))

    if (taskStore.isOpenExistingTask) {
      await window.go?.main?.App?.OpenTask?.(taskStore.taskName)
      enqueueUILog(`[${formatTime()}] opened task "${taskStore.taskName}"\n`)
      addLogForExportOnly(`opened task "${taskStore.taskName}"`, 'info')

      await new Promise((resolve) => setTimeout(resolve, 300))
      await fetchTaskState()
      
      if (taskStore.taskState.Total > 0) {
        enqueueUILog(`[${formatTime()}] read ${taskStore.taskState.Total} records\n`)
        enqueueUILog(
          `[${formatTime()}] current progress - download: ${taskStore.taskState.NDone1}/${taskStore.taskState.Total}, deduplication: ${taskStore.taskState.NDone2}/${taskStore.taskState.Total}, recognition: ${taskStore.taskState.NDone3}/${taskStore.taskState.Total}\n`
        )
        addLogForExportOnly(`read ${taskStore.taskState.Total} records`, 'info')
      } else {
        enqueueUILog(`[${formatTime()}] warning: could not read task data; check whether the task folder exists\n`)
        addLogForExportOnly('warning: could not read task data; check whether the task folder exists', 'error')
      }
    } else {
      await window.go?.main?.App?.CreateTask?.(taskStore.taskName, taskStore.filePath)
      enqueueUILog(`[${formatTime()}] task "${taskStore.taskName}" created\n`)
      addLogForExportOnly(`task "${taskStore.taskName}" created`, 'info')

      await new Promise((resolve) => setTimeout(resolve, 300))
      await fetchTaskState()
    }
    
    taskStore.setInitializing(false)
    enqueueUILog(`[${formatTime()}] task is ready; click Run to start\n`)
    addLogForExportOnly('task is ready; click Run to start', 'info')
    nextTick(() => scrollLogsToBottomHard())
  } catch (error) {
    console.error('failed to initialize task:', error)
    
    const errMsg = error.message || (typeof error === 'string' ? error : JSON.stringify(error))

    taskStore.setInitializing(false)

    enqueueUILog(`[${formatTime()}] failed to initialize task: ${errMsg}\n`)
    addLogForExportOnly(`failed to initialize task: ${errMsg}`, 'error')

    alert(`failed to initialize task: ${errMsg}\n\nPlease check:\n1. Whether the correct 5081 file was selected\n2. Whether the file is corrupted\n3. The logs for detailed error information`)
    taskStore.setStep(1)
  }
})

onUnmounted(() => {
  stopProgressPolling()
  EventsOff('process:log')
  EventsOff('process:error')
  window.removeEventListener('resize', handleResize)
  
  // clean up the throttle timer
  if (scrollThrottleTimer) {
    clearTimeout(scrollThrottleTimer)
    scrollThrottleTimer = null
  }
})
</script>

<style scoped>
.step2-container {
  height: 100vh;
  background-color: #F8F9FA;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content-wrapper {
  flex: 1;
  padding: 40px 20px;
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
  display: flex;
  flex-direction: column;
  gap: 30px;
  position: relative;
}

.progress-section.loading {
  opacity: 0.5;
  pointer-events: none;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #2C3E50;
  margin: 0 0 20px 0;
}

.progress-section {
  flex: 0 0 auto;
}

.logs-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.logs-title {
  font-size: 18px;
  font-weight: 600;
  color: #2C3E50;
  margin: 0;
}

.logs-buttons {
  display: flex;
  gap: 8px;
}

/* log box */
.logs-box {
  flex: 1;
  min-height: 300px;
  max-height: 500px;
  background-color: #1e1e1e;
  border-radius: 8px;
  overflow: auto;
  position: relative;
}

/* empty state */
.empty-logs {
  color: #6c757d;
  text-align: center;
  padding: 40px 20px;
  font-style: italic;
}

/* live log */
.logs-pre {
  margin: 0;
  padding: 16px;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #d4d4d4;
}

.btn-scroll-to-top {
  position: fixed;
  bottom: 30px;
  right: 30px;
  padding: 12px 20px;
  background-color: #3498DB;
  color: #FFFFFF;
  border: none;
  border-radius: 50px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.4);
  z-index: 1000;
  min-width: 120px;
}

.btn-scroll-to-top:hover {
  background-color: #2980B9;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(52, 152, 219, 0.5);
}

.btn-save-logs {
  padding: 8px 16px;
  background-color: #50C878;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-save-logs:hover {
  background-color: #45B869;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(80, 200, 120, 0.3);
}
</style>
