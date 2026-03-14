<template>
  <div class="right-panel">
    <div class="panel-tabs">
      <div
        v-for="tab in tabs"
        :key="tab.id"
        class="panel-tab"
        :class="{ active: activeTab === tab.id }"
        @click="activeTab = tab.id"
      >
        {{ tab.label }}
      </div>
    </div>

    <div v-if="activeTab === 'output'" class="panel-content">
      <div class="mono-label">
        {{ taskTitleLabel }}
      </div>

      <div class="terminal" ref="terminalEl">
        <div
          v-for="(line, index) in tasksStore.terminalLines"
          :key="`${line.type}-${index}`"
          class="term-line"
          :class="line.type"
        >
          {{ line.text }}
        </div>
      </div>

      <div class="mono-label" style="margin-top:8px">TASK DETAIL</div>
      <div class="px-box stats-box">
        <div class="stat-row">
          <span class="stat-label">Prompt</span>
          <span class="stat-value wrap">{{ tasksStore.selectedTaskDetail.prompt || 'n/a' }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Agent</span>
          <span class="stat-value">{{ tasksStore.taskStats.agent }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Status</span>
          <span class="stat-value status" :class="tasksStore.selectedTaskDetail.status">
            {{ tasksStore.taskStats.status }}
          </span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Elapsed</span>
          <span class="stat-value">{{ tasksStore.taskStats.elapsed }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Total tokens</span>
          <span class="stat-value">{{ tasksStore.taskStats.totalTokens }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Tokens in</span>
          <span class="stat-value">{{ tasksStore.taskStats.tokensIn }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Tokens out</span>
          <span class="stat-value">{{ tasksStore.taskStats.tokensOut }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Model</span>
          <span class="stat-value wrap">{{ tasksStore.taskStats.model }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Provider</span>
          <span class="stat-value wrap">{{ tasksStore.taskStats.provider }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Run ID</span>
          <span class="stat-value wrap">{{ tasksStore.taskStats.runId }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Session</span>
          <span class="stat-value wrap">{{ tasksStore.taskStats.sessionKey }}</span>
        </div>
        <div class="stat-row">
          <span class="stat-label">Events</span>
          <span class="stat-value">{{ tasksStore.taskStats.events }}</span>
        </div>
        <div v-if="tasksStore.selectedTaskDetail.latestMessage" class="detail-block">
          <div class="stat-label">Latest reply</div>
          <div class="reply-preview">{{ tasksStore.selectedTaskDetail.latestMessage }}</div>
        </div>
        <div v-if="tasksStore.selectedTaskDetail.errorMessage" class="detail-block">
          <div class="stat-label">Error</div>
          <div class="reply-preview error-text">{{ tasksStore.selectedTaskDetail.errorMessage }}</div>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'send'" class="panel-content">
      <div class="task-input-title">▶ NEW TASK</div>
      <div class="task-form">
        <div>
          <div class="field-label">ASSIGN TO</div>
          <select class="px-select" v-model="selectedAgent">
            <option value="auto">⚡ AUTO ROUTE</option>
            <option
              v-for="agent in agentsStore.agents"
              :key="agent.id"
              :value="agent.id"
            >
              {{ agent.emoji }} {{ agent.name }} — {{ agent.role }}
            </option>
          </select>
        </div>
        <div>
          <div class="field-label">TASK</div>
          <textarea
            class="px-textarea"
            v-model="taskInput"
            placeholder="บอก agent ว่าต้องการอะไร..."
          ></textarea>
        </div>
        <div style="display:flex;gap:8px">
          <button class="px-btn px-btn-primary" style="flex:1" @click="sendTask">
            {{ tasksStore.sending ? 'SENDING...' : '▶ SEND' }}
          </button>
          <button class="px-btn px-btn-ghost" @click="taskInput = ''">✕</button>
        </div>
      </div>

      <div v-if="tasksStore.error" class="error-banner">
        {{ tasksStore.error }}
      </div>

      <div class="mono-label" style="margin-top:8px">QUICK TASKS</div>
      <div style="display:flex;flex-direction:column;gap:6px">
        <div
          v-for="quickTask in quickTasks"
          :key="quickTask"
          class="px-box quick-task-item"
          @click="taskInput = quickTask"
        >
          ⚡ {{ quickTask }}
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'tg'" class="panel-content">
      <div class="mono-label">TELEGRAM FEED</div>
      <div class="px-box tg-placeholder">
        Telegram integration is not wired to the new backend yet.
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { useAgentsStore } from '@/stores/agents'
import { useTasksStore } from '@/stores/tasks'

const agentsStore = useAgentsStore()
const tasksStore = useTasksStore()

const activeTab = ref('output')
const selectedAgent = ref('auto')
const taskInput = ref('')
const terminalEl = ref(null)

const tabs = [
  { id: 'output', label: 'OUTPUT' },
  { id: 'send', label: 'SEND' },
  { id: 'tg', label: 'TG FEED' },
]

const quickTasks = [
  'Summarize the latest agent activity',
  'Check why the most recent task is slow',
  'Prepare a short status report for today',
]

const taskTitleLabel = computed(() => {
  const detail = tasksStore.selectedTaskDetail
  if (!detail.id) return 'NO TASK SELECTED'
  return `${detail.agentEmoji} ${detail.agent} — ${detail.id}`
})

async function sendTask() {
  if (!taskInput.value.trim()) return
  try {
    await tasksStore.sendTask(selectedAgent.value, taskInput.value.trim())
    taskInput.value = ''
    activeTab.value = 'output'
  } catch {
    // Store already exposes the error banner.
  }
  await nextTick()
  if (terminalEl.value) terminalEl.value.scrollTop = terminalEl.value.scrollHeight
}

watch(
  () => tasksStore.terminalLines.length,
  async () => {
    await nextTick()
    if (terminalEl.value) terminalEl.value.scrollTop = terminalEl.value.scrollHeight
  }
)
</script>

<style scoped>
.right-panel {
  background: var(--right-bg);
  border-left: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  backdrop-filter: blur(20px);
  transition: background 0.3s;
}

.panel-tabs {
  display: flex;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.panel-tab {
  flex: 1;
  padding: 11px 8px;
  text-align: center;
  cursor: pointer;
  font-size: 10px;
  color: var(--text2);
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: all 0.15s;
  letter-spacing: 1px;
  font-family: 'Space Mono', monospace;
}

.panel-tab:hover {
  color: var(--text);
  background: rgba(255, 255, 255, 0.03);
}

.panel-tab.active {
  color: var(--accent2);
  border-bottom-color: var(--accent);
  background: rgba(124, 92, 252, 0.05);
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.mono-label,
.field-label,
.stat-label {
  font-family: 'Space Mono', monospace;
  font-size: 9px;
  color: var(--text2);
  letter-spacing: 2px;
}

.terminal {
  background: var(--term-bg);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 12px;
  font-family: 'Space Mono', monospace;
  font-size: 11px;
  height: 220px;
  overflow-y: auto;
}

.terminal::before {
  content: '● ● ●';
  display: block;
  color: rgba(255, 255, 255, 0.15);
  font-size: 10px;
  margin-bottom: 10px;
  letter-spacing: 3px;
}

.term-line {
  margin: 3px 0;
  line-height: 1.6;
  white-space: pre-wrap;
}

.term-line.prompt {
  color: var(--accent2);
}

.term-line.output {
  color: var(--text);
  opacity: 0.85;
}

.term-line.error {
  color: var(--red);
}

.term-line.info {
  color: var(--cyan);
}

.stats-box {
  padding: 12px;
  font-size: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-family: 'Space Mono', monospace;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.stat-value {
  color: var(--text);
  text-align: right;
}

.stat-value.wrap {
  max-width: 58%;
  word-break: break-word;
}

.stat-value.status.working {
  color: var(--yellow);
}

.stat-value.status.done {
  color: var(--green);
}

.stat-value.status.error {
  color: var(--red);
}

.detail-block {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: 4px;
}

.reply-preview {
  font-size: 12px;
  line-height: 1.6;
  color: var(--text);
  white-space: pre-wrap;
}

.reply-preview.error-text {
  color: var(--red);
}

.task-input-title {
  font-family: 'Space Mono', monospace;
  font-size: 10px;
  color: var(--accent2);
  letter-spacing: 2px;
}

.task-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.quick-task-item {
  padding: 9px 12px;
  cursor: pointer;
  font-size: 11px;
  color: var(--text2);
  transition: all 0.15s;
}

.quick-task-item:hover {
  color: var(--text);
  background: rgba(255, 255, 255, 0.05);
}

.error-banner,
.tg-placeholder {
  padding: 12px;
  border-radius: 8px;
  border: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.03);
  color: var(--text2);
  font-size: 12px;
  line-height: 1.6;
}

.error-banner {
  border-color: rgba(248, 113, 113, 0.25);
  color: var(--red);
}
</style>
