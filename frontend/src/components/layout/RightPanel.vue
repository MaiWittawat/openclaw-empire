<template>
  <div class="right-panel">
    <!-- TABS -->
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

    <!-- OUTPUT TAB -->
    <div v-if="activeTab === 'output'" class="panel-content">
      <div class="mono-label">แพร — TASK #015</div>
      <div class="terminal" ref="terminalEl">
        <div
          v-for="(line, i) in tasksStore.terminalLines"
          :key="i"
          class="term-line"
          :class="line.type"
        >
          <template v-if="line.type === 'cursor'">
            <span class="term-cursor"></span>
          </template>
          <template v-else>{{ line.text }}</template>
        </div>
      </div>

      <div class="mono-label" style="margin-top:8px">TASK STATS</div>
      <div class="px-box" style="padding:12px;font-size:12px;display:flex;flex-direction:column;gap:8px;font-family:'Space Mono',monospace">
        <div class="stat-row">
          <span style="color:var(--text2)">Agent</span>
          <span>{{ tasksStore.taskStats.agent }}</span>
        </div>
        <div class="stat-row">
          <span style="color:var(--text2)">Status</span>
          <span style="color:var(--yellow)">▶ {{ tasksStore.taskStats.status }}</span>
        </div>
        <div class="stat-row">
          <span style="color:var(--text2)">Tokens in</span>
          <span style="color:var(--cyan)">{{ tasksStore.taskStats.tokensIn }}</span>
        </div>
        <div class="stat-row">
          <span style="color:var(--text2)">Tokens out</span>
          <span style="color:var(--accent2)">{{ tasksStore.taskStats.tokensOut }}</span>
        </div>
        <div class="stat-row">
          <span style="color:var(--text2)">Cost</span>
          <span style="color:var(--yellow)">{{ tasksStore.taskStats.cost }}</span>
        </div>
        <div class="stat-row">
          <span style="color:var(--text2)">Elapsed</span>
          <span>{{ tasksStore.taskStats.elapsed }}</span>
        </div>
      </div>
    </div>

    <!-- SEND TAB -->
    <div v-if="activeTab === 'send'" class="panel-content">
      <div class="task-input-title">▶ NEW TASK</div>
      <div class="task-form">
        <div>
          <div style="font-size:13px;color:var(--text2);margin-bottom:4px">ASSIGN TO</div>
          <select class="px-select" v-model="selectedAgent">
            <option value="auto">⚡ AUTO ROUTE</option>
            <option value="parae">👩‍💻 แพร — Tech Lead</option>
            <option value="nova">🔍 NOVA — Researcher</option>
            <option value="lyra">✍️ LYRA — Writer</option>
            <option value="forge">🛠️ FORGE — DevOps</option>
          </select>
        </div>
        <div>
          <div style="font-size:13px;color:var(--text2);margin-bottom:4px">TASK</div>
          <textarea class="px-textarea" v-model="taskInput" placeholder="บอก agent ว่าต้องการอะไร..."></textarea>
        </div>
        <div style="display:flex;gap:8px">
          <button class="px-btn px-btn-primary" style="flex:1" @click="sendTask">▶ SEND</button>
          <button class="px-btn px-btn-ghost" @click="taskInput = ''">✕</button>
        </div>
      </div>

      <div class="mono-label" style="margin-top:8px">QUICK TASKS</div>
      <div style="display:flex;flex-direction:column;gap:6px">
        <div
          v-for="qt in quickTasks"
          :key="qt"
          class="px-box quick-task-item"
          @click="taskInput = qt"
        >
          ⚡ {{ qt }}
        </div>
      </div>
    </div>

    <!-- TELEGRAM TAB -->
    <div v-if="activeTab === 'tg'" class="panel-content">
      <div class="mono-label">TELEGRAM FEED</div>
      <div class="tg-feed">
        <div
          v-for="(msg, i) in tgMessages"
          :key="i"
          class="tg-msg"
          :style="msg.isOut ? 'flex-direction:row-reverse' : ''"
        >
          <div style="font-size:20px">{{ msg.avatar }}</div>
          <div class="tg-bubble" :class="{ out: msg.isOut }" :style="msg.bubbleStyle">
            <div class="tg-name" :style="msg.isOut ? 'text-align:right' : ''">{{ msg.name }}</div>
            <span v-if="msg.highlight" :style="`color:${msg.highlight}`">{{ msg.text }}</span>
            <span v-else>{{ msg.text }}</span>
            <div class="tg-time">{{ msg.time }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, watch } from 'vue'
import { useTasksStore } from '@/stores/tasks'

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

const quickTasks = ['Review latest code', 'Status report', 'Run tests']

const tgMessages = [
  { avatar: '👤', name: 'CEO (You)', text: 'แพร ช่วยเขียน FastAPI endpoint สำหรับ Telegram webhook ด้วยนะ', time: '14:32', isOut: false },
  { avatar: '👩‍💻', name: 'แพร', text: 'รับค่ะ กำลังอ่าน codebase อยู่ จะเริ่มเขียนใน /app/routes/telegram.py เลยนะคะ', time: '14:32', isOut: true },
  { avatar: '👤', name: 'CEO (You)', text: 'NOVA หาข้อมูล best practice multi-agent routing ให้ด้วย', time: '14:35', isOut: false },
  { avatar: '🔍', name: 'NOVA', text: 'กำลังค้นหาค่ะ จะสรุปให้ใน 5 นาที', time: '14:35', isOut: true },
  { avatar: '👩‍💻', name: 'แพร', text: '▶ กำลังเขียน endpoint...', time: '14:36', isOut: true, highlight: 'var(--yellow)', bubbleStyle: 'border-color:var(--yellow)' },
]

async function sendTask() {
  if (!taskInput.value.trim()) return
  await tasksStore.sendTask(selectedAgent.value, taskInput.value.trim())
  taskInput.value = ''
  activeTab.value = 'output'
  await nextTick()
  if (terminalEl.value) terminalEl.value.scrollTop = terminalEl.value.scrollHeight
}

watch(() => tasksStore.terminalLines.length, async () => {
  await nextTick()
  if (terminalEl.value) terminalEl.value.scrollTop = terminalEl.value.scrollHeight
})
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
.panel-tab:hover { color: var(--text); background: rgba(255,255,255,0.03); }
.panel-tab.active { color: var(--accent2); border-bottom-color: var(--accent); background: rgba(124,92,252,0.05); }

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.mono-label {
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
  height: 180px;
  overflow-y: auto;
}
.terminal::before {
  content: '● ● ●';
  display: block;
  color: rgba(255,255,255,0.15);
  font-size: 10px;
  margin-bottom: 10px;
  letter-spacing: 3px;
}
.term-line { margin: 3px 0; line-height: 1.6; }
.term-line.prompt { color: var(--accent2); }
.term-line.output { color: var(--text); opacity: 0.8; }
.term-line.error  { color: var(--red); }
.term-line.info   { color: var(--cyan); }

.stat-row { display: flex; justify-content: space-between; }

.task-input-title {
  font-family: 'Space Mono', monospace;
  font-size: 10px;
  color: var(--accent2);
  letter-spacing: 2px;
}
.task-form { display: flex; flex-direction: column; gap: 10px; }

.quick-task-item {
  padding: 9px 12px;
  cursor: pointer;
  font-size: 11px;
  color: var(--text2);
  transition: all 0.15s;
}
.quick-task-item:hover { color: var(--text); background: rgba(255,255,255,0.05); }

.tg-feed { display: flex; flex-direction: column; gap: 10px; }
.tg-msg { display: flex; gap: 8px; align-items: flex-start; }
.tg-bubble {
  background: rgba(255,255,255,0.04);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 10px;
  font-size: 12px;
  line-height: 1.6;
  flex: 1;
}
.tg-bubble.out { background: rgba(124,92,252,0.08); border-color: rgba(124,92,252,0.2); }
.tg-name { font-size: 10px; color: var(--accent2); margin-bottom: 3px; font-family: 'Space Mono', monospace; }
.tg-time { font-size: 10px; color: var(--text2); margin-top: 4px; text-align: right; }
</style>
