<template>
  <div>
    <SectionTitle text="◈ ALL TASKS" />

    <!-- Filter bar -->
    <div class="filter-bar px-box" style="padding: 12px 16px; margin-bottom: 12px; display: flex; gap: 8px; flex-wrap: wrap;">
      <button
        v-for="f in filters"
        :key="f.value"
        class="filter-btn"
        :class="{ active: activeFilter === f.value }"
        @click="activeFilter = f.value"
      >
        {{ f.label }}
      </button>
    </div>

    <div class="task-list">
      <TaskItem
        v-for="task in filteredTasks"
        :key="task.id"
        :task="task"
        :isActive="tasksStore.selectedTaskId === task.id"
        @select="tasksStore.selectTask"
      />
    </div>

    <div v-if="filteredTasks.length === 0" class="empty-state">
      <div style="font-size: 32px">📭</div>
      <div style="font-family: 'Space Mono', monospace; font-size: 11px; color: var(--text2); margin-top: 8px">NO TASKS FOUND</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import SectionTitle from '@/components/layout/SectionTitle.vue'
import TaskItem from '@/components/tasks/TaskItem.vue'
import { useTasksStore } from '@/stores/tasks'

const tasksStore = useTasksStore()
const activeFilter = ref('all')

const filters = [
  { label: 'ALL', value: 'all' },
  { label: 'WORKING', value: 'working' },
  { label: 'DONE', value: 'done' },
  { label: 'ERROR', value: 'error' },
]

const filteredTasks = computed(() =>
  activeFilter.value === 'all'
    ? tasksStore.tasks
    : tasksStore.tasks.filter(t => t.status === activeFilter.value)
)
</script>

<style scoped>
.task-list { display: flex; flex-direction: column; gap: 6px; }
.filter-btn {
  font-family: 'Space Mono', monospace;
  font-size: 10px;
  padding: 5px 12px;
  border: 1px solid var(--border);
  border-radius: 20px;
  background: transparent;
  color: var(--text2);
  cursor: pointer;
  letter-spacing: 1px;
  transition: all 0.15s;
}
.filter-btn:hover { color: var(--text); border-color: var(--border2); }
.filter-btn.active { color: var(--accent2); border-color: var(--accent); background: rgba(124,92,252,0.1); }
.empty-state { text-align: center; padding: 40px; }
</style>
