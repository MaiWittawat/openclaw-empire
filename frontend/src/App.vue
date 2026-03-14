<template>
  <div id="app-root">
    <AppHeader />
    <div class="app-layout">
      <AppSidebar />
      <main class="main-content">
        <RouterView />
      </main>
      <RightPanel />
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'
import { useAgentsStore } from '@/stores/agents'
import { useTasksStore } from '@/stores/tasks'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import RightPanel from '@/components/layout/RightPanel.vue'

const themeStore = useThemeStore()
const agentsStore = useAgentsStore()
const tasksStore = useTasksStore()

onMounted(async () => {
  themeStore.init()
  await Promise.all([
    agentsStore.fetchAgents(),
    tasksStore.fetchTasks(),
    tasksStore.fetchStats(),
  ])
})
</script>

<style scoped>
#app-root {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
.app-layout {
  display: grid;
  grid-template-columns: 220px 1fr 300px;
  flex: 1;
  min-height: calc(100vh - 53px);
}
.main-content {
  padding: 20px;
  overflow-y: auto;
  background: transparent;
  display: flex;
  flex-direction: column;
  gap: 20px;
  background-image:
    linear-gradient(rgba(124,92,252,0.025) 1px, transparent 1px),
    linear-gradient(90deg, rgba(124,92,252,0.025) 1px, transparent 1px);
  background-size: 40px 40px;
}
</style>
