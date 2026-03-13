<template>
  <div
    class="task-item px-box"
    :class="{ active: isActive }"
    @click="$emit('select', task.id)"
  >
    <div class="task-dot" :class="task.status"></div>
    <div class="task-content">
      <div class="task-title">{{ task.title }}</div>
      <div class="task-meta">
        <span>{{ task.agentEmoji }} {{ task.agent }}</span>
        <span>🕐 {{ task.time }}</span>
        <span :style="statusColor">{{ statusLabel }}</span>
      </div>
    </div>
    <div class="task-tokens">{{ task.tokens.toLocaleString() }} tok</div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  task: { type: Object, required: true },
  isActive: { type: Boolean, default: false },
})
defineEmits(['select'])

const statusLabel = computed(() => {
  const map = { working: 'IN PROGRESS', done: 'DONE', error: 'ERROR', pending: 'PENDING' }
  return map[props.task.status] || props.task.status.toUpperCase()
})
const statusColor = computed(() => {
  const map = { working: 'color:var(--yellow)', done: 'color:var(--green)', error: 'color:var(--red)' }
  return map[props.task.status] || ''
})
</script>

<style scoped>
.task-item {
  padding: 11px 14px;
  display: flex;
  gap: 12px;
  align-items: flex-start;
  cursor: pointer;
  transition: all 0.15s;
}
.task-item:hover { background: rgba(255,255,255,0.04); }
.task-item.active { border-color: rgba(167,139,250,0.35); background: rgba(124,92,252,0.05); }

.task-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; margin-top: 5px; }
.task-dot.done    { background: var(--green); box-shadow: 0 0 5px var(--green); }
.task-dot.working { background: var(--yellow); box-shadow: 0 0 5px var(--yellow); animation: pulse-dot 1s ease-in-out infinite; }
.task-dot.pending { background: var(--text2); }
.task-dot.error   { background: var(--red); box-shadow: 0 0 5px var(--red); }

.task-content { flex: 1; min-width: 0; }
.task-title { font-size: 12px; color: var(--text); margin-bottom: 4px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.task-meta { font-size: 10px; color: var(--text2); display: flex; gap: 10px; flex-wrap: wrap; }
.task-tokens { font-size: 10px; color: var(--accent2); flex-shrink: 0; font-family: 'Space Mono', monospace; }
</style>
