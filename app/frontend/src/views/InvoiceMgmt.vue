<template>
  <div class="invoice-mgmt-container">
    <Header />
    
    <div class="content-wrapper">
      <div class="content-card">
        <div class="header-section">
          <button class="btn-back" @click="handleGoBack" title="Go back">
            <span class="back-icon">←</span>
            <span class="btn-text">Back</span>
          </button>
          <h2 class="section-title">Invoice Management</h2>
        </div>
        
        <!-- imported invoice list -->
        <div class="invoice-list-section">
          <div class="section-header" @click="toggleTableExpanded">
            <h3 class="subsection-title">Imported Invoice List</h3>
            <span class="toggle-icon">{{ isTableExpanded ? '▼' : '▶' }}</span>
          </div>
          
          <div v-if="isTableExpanded" class="table-container">
            <table v-if="invoiceIndices.length > 0" class="invoice-table">
              <thead>
                <tr>
                  <th>Year</th>
                  <th>Month</th>
                  <th>Branch</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, index) in invoiceIndices" :key="index">
                  <td>{{ item.Year }}</td>
                  <td>{{ item.Month }}</td>
                  <td>{{ item.Branch }}</td>
                  <td>
                    <button class="btn-delete" @click="handleDelete(item.Raw)" title="Delete">
                      Delete
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
            <p v-else class="empty-message">No imported invoice data yet</p>
          </div>
        </div>
        
        <!-- loading overlay -->
        <div v-if="isLoading" class="loading-overlay">
          <div class="loading-spinner">
            <div class="spinner"></div>
            <p class="loading-text">Parsing file, please wait...</p>
          </div>
        </div>
        
        <!-- add invoice form -->
        <div class="add-invoice-section">
          <h3 class="subsection-title">Add Invoices</h3>
          
          <div class="form-grid">
            <div class="form-group">
              <label class="form-label">Year</label>
              <select v-model="formData.year" class="form-select">
                <option v-for="year in yearOptions" :key="year" :value="year">
                  {{ year }}
                </option>
              </select>
            </div>
            
            <div class="form-group">
              <label class="form-label">Month</label>
              <select v-model="formData.month" class="form-select">
                <option v-for="month in monthOptions" :key="month" :value="month">
                  {{ month }}
                </option>
              </select>
            </div>
            
            <div class="form-group">
              <label class="form-label">Branch</label>
              <input
                v-model="formData.branch"
                type="text"
                class="form-input"
                placeholder="Enter the branch name"
              />
            </div>
          </div>
          
          <div class="file-selection">
            <button class="btn-select-file" @click="handleSelectFile">
              <span class="file-icon">📁</span>
              <span class="btn-text">{{ formData.filePath ? 'Reselect File' : 'Select File' }}</span>
            </button>
            
            <div v-if="formData.filePath" class="file-info">
              <span class="file-icon">📄</span>
              <span class="file-name">{{ formData.fileName }}</span>
            </div>
          </div>
          
          <div class="action-buttons">
            <button class="btn-add" :disabled="!canAdd" @click="handleAdd">
              <span class="add-icon">+</span>
              <span class="btn-text">Add</span>
            </button>
          </div>
          
          <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
          <p v-if="successMessage" class="success-message">{{ successMessage }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useTaskStore } from '../stores/taskStore'
import Header from '../components/Header.vue'

const taskStore = useTaskStore()

// state
const isTableExpanded = ref(true)
const invoiceIndices = ref([])
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

// form data
const formData = ref({
  year: new Date().getFullYear().toString(),
  month: String(new Date().getMonth() + 1).padStart(2, '0'),
  branch: '',
  filePath: '',
  fileName: ''
})

// computed properties
const yearOptions = computed(() => {
  const currentYear = new Date().getFullYear()
  return [currentYear, currentYear - 1, currentYear - 2]
})

const monthOptions = computed(() => {
  return Array.from({ length: 12 }, (_, i) => String(i + 1).padStart(2, '0'))
})

const canAdd = computed(() => {
  return formData.value.year &&
         formData.value.month &&
         formData.value.branch.trim() !== '' &&
         formData.value.filePath !== ''
})

// methods
function handleGoBack() {
  taskStore.setStep(1)
}

function toggleTableExpanded() {
  isTableExpanded.value = !isTableExpanded.value
}

async function loadInvoiceIndices() {
  try {
    const data = await window.go?.main?.App?.GovInvoiceIndices?.()
    if (data) {
      invoiceIndices.value = data
    }
  } catch (error) {
    console.error('failed to load invoice indices:', error)
    errorMessage.value = 'failed to load invoice indices: ' + error.message
  }
}

async function handleSelectFile() {
  try {
    const filePath = await window.go?.main?.App?.SelectFile?.('Select an invoice Excel file')
    if (filePath) {
      formData.value.filePath = filePath
      formData.value.fileName = filePath.split(/[/\\]/).pop()
      errorMessage.value = ''
    }
  } catch (error) {
    console.error('failed to select file:', error)
    errorMessage.value = 'failed to select a file'
  }
}

async function handleDelete(rawIndex) {
  if (!confirm(`Are you sure you want to delete the invoice data of ${rawIndex}?`)) {
    return
  }
  
  try {
    await window.go?.main?.App?.DeleteGovInvoiceIndex?.(rawIndex)
    successMessage.value = 'deleted successfully'
    errorMessage.value = ''
    await loadInvoiceIndices()
  } catch (error) {
    console.error('failed to delete:', error)
    errorMessage.value = 'failed to delete: ' + error.message
    successMessage.value = ''
  }
}

async function handleAdd() {
  if (!canAdd.value) return
  
  errorMessage.value = ''
  successMessage.value = ''
  isLoading.value = true
  
  try {
    const rawIndex = `${formData.value.year}-${formData.value.month}-${formData.value.branch}`
    await window.go?.main?.App?.ImportGovInvoices?.(formData.value.filePath, rawIndex)
    
    successMessage.value = 'added successfully'
    formData.value.branch = ''
    formData.value.filePath = ''
    formData.value.fileName = ''
    
    await loadInvoiceIndices()
  } catch (error) {
    console.error('failed to add:', error)
    errorMessage.value = 'failed to add: ' + error
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadInvoiceIndices()
})
</script>

<style scoped>
.invoice-mgmt-container {
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
}

.header-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
  padding-bottom: 16px;
  border-bottom: 2px solid #E0E0E0;
}

.section-title {
  font-size: 24px;
  font-weight: 600;
  color: #2C3E50;
  margin: 0;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
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

.btn-back:hover {
  background-color: #5A6268;
  transform: translateX(-2px);
}

.back-icon {
  font-size: 16px;
}

/* invoice list */
.invoice-list-section {
  margin-bottom: 32px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #F5F7FA;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s ease;
  user-select: none;
}

.section-header:hover {
  background-color: #E8EBF0;
}

.subsection-title {
  font-size: 18px;
  font-weight: 600;
  color: #2C3E50;
  margin: 0;
}

.toggle-icon {
  font-size: 14px;
  color: #666666;
}

.table-container {
  margin-top: 16px;
  overflow-x: auto;
}

.invoice-table {
  width: 100%;
  border-collapse: collapse;
  background-color: #FFFFFF;
}

.invoice-table thead {
  background-color: #F5F7FA;
}

.invoice-table th {
  padding: 12px 16px;
  text-align: left;
  font-size: 14px;
  font-weight: 600;
  color: #2C3E50;
  border-bottom: 2px solid #E0E0E0;
}

.invoice-table td {
  padding: 12px 16px;
  font-size: 14px;
  color: #333333;
  border-bottom: 1px solid #E0E0E0;
}

.invoice-table tbody tr:hover {
  background-color: #F5F7FA;
}

.invoice-table tbody tr:last-child td {
  border-bottom: none;
}

.btn-delete {
  padding: 6px 16px;
  background-color: #DC3545;
  color: #FFFFFF;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-delete:hover {
  background-color: #C82333;
  transform: translateY(-1px);
}

.empty-message {
  padding: 24px;
  text-align: center;
  color: #999999;
  font-size: 14px;
}

/* add invoice form */
.add-invoice-section {
  padding: 24px;
  background-color: #F5F7FA;
  border-radius: 8px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #2C3E50;
}

.form-select,
.form-input {
  padding: 10px 12px;
  border: 1px solid #D0D0D0;
  border-radius: 6px;
  font-size: 14px;
  color: #333333;
  transition: all 0.2s ease;
  background-color: #FFFFFF;
}

.form-select:focus,
.form-input:focus {
  outline: none;
  border-color: #2C8CFF;
  box-shadow: 0 0 0 3px rgba(44, 140, 255, 0.1);
}

.file-selection {
  margin-bottom: 16px;
}

.btn-select-file {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background-color: #2C8CFF;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-select-file:hover {
  background-color: #1A7AE6;
  transform: translateY(-1px);
}

.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  padding: 10px 12px;
  background-color: #FFFFFF;
  border: 1px solid #E0E0E0;
  border-radius: 6px;
}

.file-name {
  font-size: 14px;
  color: #333333;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.btn-add {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background-color: #28A745;
  color: #FFFFFF;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add:hover:not(:disabled) {
  background-color: #218838;
  transform: translateY(-1px);
}

.btn-add:disabled {
  background-color: #CCCCCC;
  cursor: not-allowed;
}

.add-icon {
  font-size: 18px;
  font-weight: bold;
}

.error-message {
  margin-top: 12px;
  padding: 10px 12px;
  background-color: #F8D7DA;
  color: #721C24;
  border: 1px solid #F5C6CB;
  border-radius: 6px;
  font-size: 14px;
}

.success-message {
  margin-top: 12px;
  padding: 10px 12px;
  background-color: #D4EDDA;
  color: #155724;
  border: 1px solid #C3E6CB;
  border-radius: 6px;
  font-size: 14px;
}

/* loading overlay */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(4px);
}

.loading-spinner {
  background-color: #FFFFFF;
  border-radius: 12px;
  padding: 40px 60px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #E0E0E0;
  border-top: 4px solid #2C8CFF;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-text {
  font-size: 16px;
  font-weight: 500;
  color: #2C3E50;
  margin: 0;
}
</style>
