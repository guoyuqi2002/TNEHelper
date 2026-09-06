<script>
import { ref, onMounted } from 'vue';
import SplashScreen from './components/SplashScreen.vue';
import Tab1TaxFile from './components/Tab1TaxFile.vue';
import Tab2ReimburseCheck from './components/Tab2ReimburseCheck.vue';

export default {
    components: { SplashScreen, Tab1TaxFile, Tab2ReimburseCheck },
    setup() {
        const isLoading = ref(true);
        const activeTab = ref('tab1');
        const isTab2Imported = ref(false);

        // auto-hide the splash screen
        onMounted(() => {
            setTimeout(() => {
                isLoading.value = false;
            }, 2500);
        });

        // receive the Tab2 import-success event
        const handleTab2ImportSuccess = () => {
            isTab2Imported.value = true;
        };

        return {
            isLoading,
            activeTab,
            isTab2Imported,
            handleTab2ImportSuccess
        };
    }
};
</script>

<template>
    <div>
        <!-- splash screen -->
        <SplashScreen v-if="isLoading" />

        <!-- main view -->
        <div v-else class="main-container fade-in">
            <!-- header -->
            <header class="header">
                <svg class="logo" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <rect x="3" y="5" width="18" height="14" rx="2" stroke="currentColor" stroke-width="2" />
                    <path d="M7 9.5h4M7 13h10" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
                    <circle cx="16.5" cy="9.5" r="1.5" fill="currentColor" />
                </svg>
                <h1 class="app-title">Expense Reimbursement Checker</h1>
            </header>

            <!-- tab navigation (two side-by-side panels) -->
            <el-tabs v-model="activeTab" class="custom-tabs">
                <el-tab-pane label="Tax File Processing" name="tab1">
                    <Tab1TaxFile />
                </el-tab-pane>
                <el-tab-pane label="Reimbursement Check" name="tab2">
                    <Tab2ReimburseCheck
                        @import-success="handleTab2ImportSuccess"
                        :is-imported="isTab2Imported"
                    />
                </el-tab-pane>
            </el-tabs>

            <!-- content area -->
            <div class="content-card">
                <!-- no conditional rendering of tab components needed here; they are already inside el-tab-pane -->
            </div>
        </div>
    </div>
</template>

<style scoped>
.main-container {
    padding: 25px 30px;
}
.fade-in {
    animation: fadeIn 0.5s ease-in;
}
@keyframes fadeIn {
    from { opacity: 0; transform: scale(0.98); }
    to { opacity: 1; transform: scale(1); }
}

.header {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 25px;
}
.logo {
    width: 40px;
    height: 40px;
    margin-right: 15px;
    color: #28a745;
}
.app-title {
    font-size: 1.7rem;
    font-weight: 600;
    color: #333;
}

/* tab navigation: side-by-side layout */
.custom-tabs {
    margin-bottom: 20px;
    width: 100%;
}
.custom-tabs.el-tabs__nav {
    display: flex;
    width: 100%;
    gap: 10px; /* gap between the two tabs */
}
.custom-tabs.el-tabs__item {
    flex: 1; /* each tab takes 50% width */
    text-align: center;
    color: #fff;
    background-color: #28a745;
    border-radius: 8px;
    padding: 10px 0;
    opacity: 1;
}
.custom-tabs.el-tabs__item.is-active {
    background-color: #1e7e34; /* darker shade for the active tab */
}
.custom-tabs.el-tabs__nav-wrap {
    background: transparent!important;
    padding: 0;
    border-radius: 0;
}

.content-card {
    background-color: #f8f9fa;
    border-radius: 12px;
    padding: 30px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    min-height: 500px;
}
</style>
