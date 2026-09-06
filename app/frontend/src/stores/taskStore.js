import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useTaskStore = defineStore('task', () => {
  // base state
  const currentStep = ref(1)
  const filePath = ref('')
  const fileName = ref('')
  const taskName = ref('')
  const exportPath = ref('')
  const taskNameHistory = ref([])
  const isOpenExistingTask = ref(false)

  // task state
  const taskState = ref({
    NDone1: 0,
    NDone2: 0,
    NDone3: 0,
    Total: 0
  })

  // error state
  const errors = ref({
    download: null,
    deduplicate: null,
    recognize: null
  })

  const errorLogs = ref({
    download: null,
    deduplicate: null,
    recognize: null
  })

  // logs and running state
  const allLogs = ref([])      // kept in full for export
  const isTask1Running = ref(false)
  const isTask2Running = ref(false)
  const isTask3Running = ref(false)
  const isInitializing = ref(false)

  // progress computed properties
  const progress1 = computed(() => {
    if (taskState.value.Total === 0) return 0
    return Math.round((taskState.value.NDone1 / taskState.value.Total) * 100)
  })

  const progress2 = computed(() => {
    if (taskState.value.Total === 0) return 0
    return Math.round((taskState.value.NDone2 / taskState.value.Total) * 100)
  })

  const progress3 = computed(() => {
    if (taskState.value.Total === 0) return 0
    return Math.round((taskState.value.NDone3 / taskState.value.Total) * 100)
  })

  const progressText1 = computed(() => `${taskState.value.NDone1}/${taskState.value.Total}`)
  const progressText2 = computed(() => `${taskState.value.NDone2}/${taskState.value.Total}`)
  const progressText3 = computed(() => `${taskState.value.NDone3}/${taskState.value.Total}`)

  // base setters
  function setStep(step) {
    currentStep.value = step
  }

  function setFile(path, name) {
    filePath.value = path
    fileName.value = name
  }

  function setTaskName(name) {
    taskName.value = name
  }

  function setExportPath(path) {
    exportPath.value = path
  }

  function setInitializing(value) {
    isInitializing.value = value
  }

  function setOpenExistingTask(value) {
    isOpenExistingTask.value = value
  }

  // history management
  function addTaskNameToHistory(name) {
    if (!name || name.trim() === '') return

    const trimmedName = name.trim()
    const filtered = taskNameHistory.value.filter(item => item !== trimmedName)
    filtered.unshift(trimmedName)

    if (filtered.length > 3) filtered.length = 3
    taskNameHistory.value = filtered

    try {
      localStorage.setItem('taskNameHistory', JSON.stringify(taskNameHistory.value))
    } catch (e) {
      console.error('failed to save history:', e)
    }
  }

  function loadTaskNameHistory() {
    try {
      const saved = localStorage.getItem('taskNameHistory')
      if (saved) taskNameHistory.value = JSON.parse(saved).slice(0, 3)
    } catch (e) {
      console.error('failed to load history:', e)
    }
  }

  // task state update
  function updateTaskState(state) {
    taskState.value.NDone1 = state.NDone1 || 0
    taskState.value.NDone2 = state.NDone2 || 0
    taskState.value.NDone3 = state.NDone3 || 0
    taskState.value.Total = state.Total || 0
  }

  // error handling
  function setError(type, error, errorLog = null) {
    if (type === 'download') {
      errors.value.download = error
      errorLogs.value.download = errorLog || error
    } else if (type === 'deduplicate') {
      errors.value.deduplicate = error
      errorLogs.value.deduplicate = errorLog || error
    } else if (type === 'recognize') {
      errors.value.recognize = error
      errorLogs.value.recognize = errorLog || error
    }
  }

  function clearError(type) {
    if (type === 'download') {
      errors.value.download = null
      errorLogs.value.download = null
    } else if (type === 'deduplicate') {
      errors.value.deduplicate = null
      errorLogs.value.deduplicate = null
    } else if (type === 'recognize') {
      errors.value.recognize = null
      errorLogs.value.recognize = null
    }
  }

  // ========== log management ==========
  // the UI no longer renders from logs, so this method only writes to allLogs
  function addLogForExportOnly(message, type = 'info', timestamp = null) {
    const ts = timestamp || new Date().toLocaleTimeString('en-GB', { hour12: false })
    allLogs.value.push({
      time: ts,
      message,
      type
    })
  }

  function addLog(message, type = 'info') {
    addLogForExportOnly(message, type)
  }

  // task running state
  function setTaskRunning(taskId, running) {
    if (taskId === 1) {
      isTask1Running.value = running
    } else if (taskId === 2) {
      isTask2Running.value = running
    } else if (taskId === 3) {
      isTask3Running.value = running
    }
  }

  // initialization
  loadTaskNameHistory()

  return {
    // State
    currentStep,
    filePath,
    fileName,
    taskName,
    exportPath,
    taskNameHistory,
    taskState,
    errors,
    errorLogs,
    allLogs,
    isTask1Running,
    isTask2Running,
    isTask3Running,
    isInitializing,
    isOpenExistingTask,

    // Computed
    progress1,
    progress2,
    progress3,
    progressText1,
    progressText2,
    progressText3,

    // Actions
    setStep,
    setFile,
    setTaskName,
    setExportPath,
    addTaskNameToHistory,
    loadTaskNameHistory,
    updateTaskState,
    setError,
    clearError,
    addLog,
    addLogForExportOnly,
    setTaskRunning,
    setInitializing,
    setOpenExistingTask
  }
})
