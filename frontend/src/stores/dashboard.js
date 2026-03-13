import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { agentApi, taskApi, statsApi, telegramApi } from '@/api/index.js'

export const useDashboardStore = defineStore('dashboard', () => {
  // State
  const agents = ref([
    { id: 'parae', name: 'แพร', role: 'Tech Lead', emoji: '👩‍💻', status: 'working', tokens: 12450, tokenPct: 62, color: '#a78bfa' },
    { id: 'nova', name: 'NOVA', role: 'Researcher', emoji: '🔍', status: 'working', tokens: 8820, tokenPct: 44, color: '#34d399' },
    { id: 'forge', name: 'FORGE', role: 'DevOps', emoji: '🛠️', status: 'idle', tokens: 3200, tokenPct: 16, color: '#60a5fa' },
    { id: 'lyra', name: 'LYRA', role: 'Writer', emoji: '✍️', status: 'done', tokens: 5400, tokenPct: 27, color: '#f472b6' },
  ])
  const tasks = ref([
    { id: 1, title: 'เขียน FastAPI endpoint สำหรับ Telegram webhook', agent: '👩‍💻 แพร', agentId: 'parae', time: '2m ago', status: 'working', tokens: 12450 },
    { id: 2, title: 'ค้นหา best practice สำหรับ multi-agent routing', agent: '🔍 NOVA', agentId: 'nova', time: '5m ago', status: 'working', tokens: 8820 },
    { id: 3, title: 'เขียน README สำหรับ project', agent: '✍️ LYRA', agentId: 'lyra', time: '18m ago', status: 'done', tokens: 5400 },
    { id: 4, title: 'สร้าง Docker Compose config', agent: '🛠️ FORGE', agentId: 'forge', time: '1h ago', status: 'done', tokens: 3200 },
    { id: 5, title: 'Deploy to VPS — timeout error', agent: '🛠️ FORGE', agentId: 'forge', time: '2h ago', status: 'error', tokens: 1100 },
  ])
  const stats = ref({ activeAgents: 2, totalAgents: 4, tasksToday: 17, pendingTasks: 3, tokensUsed: 284000, costToday: 0.42, successRate: 94, doneTasks: 16 })
  const tokenChart = ref([40, 55, 35, 70, 60, 85, 100])
  const telegramFeed = ref([
    { id: 1, from: 'CEO (You)', emoji: '👤', text: 'แพร ช่วยเขียน FastAPI endpoint สำหรับ Telegram webhook ด้วยนะ', time: '14:32', isOut: false },
    { id: 2, from: 'แพร', emoji: '👩‍💻', text: 'รับค่ะ กำลังอ่าน codebase อยู่ จะเริ่มเขียนใน /app/routes/telegram.py เลยนะคะ', time: '14:32', isOut: true },
    { id: 3, from: 'CEO (You)', emoji: '👤', text: 'NOVA หาข้อมูล best practice multi-agent routing ให้ด้วย', time: '14:35', isOut: false },
    { id: 4, from: 'NOVA', emoji: '🔍', text: 'กำลังค้นหาค่ะ จะสรุปให้ใน 5 นาที', time: '14:35', isOut: true },
  ])
  const terminalLines = ref([
    { type: 'prompt', text: '$ task received: FastAPI webhook endpoint' },
    { type: 'info', text: '→ reading project structure...' },
    { type: 'output', text: 'found: /app/main.py, /app/routes/' },
    { type: 'prompt', text: '$ writing /app/routes/telegram.py' },
    { type: 'output', text: '@router.post("/webhook")' },
    { type: 'output', text: 'async def telegram_webhook(req: Request):' },
    { type: 'output', text: '  data = await req.json()' },
    { type: 'info', text: '→ adding error handling...' },
  ])
  const selectedAgent = ref(null)
  const selectedTask = ref(1)
  const loading = ref(false)
  const theme = ref(localStorage.getItem('theme') || 'dark')

  // Getters
  const activeAgents = computed(() => agents.value.filter(a => a.status === 'working'))

  // Actions
  function toggleTheme() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    localStorage.setItem('theme', theme.value)
    document.documentElement.setAttribute('data-theme', theme.value === 'light' ? 'light' : '')
  }

  function selectAgent(id) {
    selectedAgent.value = id === selectedAgent.value ? null : id
  }

  function selectTask(id) {
    selectedTask.value = id
  }

  async function sendTask(agentId, taskText) {
    if (!taskText.trim()) return false
    try {
      loading.value = true
      // Try real API first
      try {
        await agentApi.sendTask(agentId, taskText)
      } catch (e) {
        // Fallback: add to local state
        console.warn('API not available, using local state')
      }
      terminalLines.value.push({ type: 'prompt', text: `$ new task → ${agentId}: ${taskText}` })
      tasks.value.unshift({ id: Date.now(), title: taskText, agent: agentId, agentId, time: 'now', status: 'working', tokens: 0 })
      return true
    } finally {
      loading.value = false
    }
  }

  async function fetchOverview() {
    try {
      const data = await statsApi.getOverview()
      if (data) stats.value = { ...stats.value, ...data }
    } catch (e) {
      console.warn('Using local stats data')
    }
  }

  async function fetchAgents() {
    try {
      const data = await agentApi.getAll()
      if (data && data.length) agents.value = data
    } catch (e) {
      console.warn('Using local agents data')
    }
  }

  async function fetchTasks() {
    try {
      const data = await taskApi.getHistory()
      if (data && data.length) tasks.value = data
    } catch (e) {
      console.warn('Using local tasks data')
    }
  }

  async function fetchTelegram() {
    try {
      const data = await telegramApi.getFeed()
      if (data && data.length) telegramFeed.value = data
    } catch (e) {
      console.warn('Using local telegram data')
    }
  }

  function init() {
    // Apply saved theme
    if (theme.value === 'light') {
      document.documentElement.setAttribute('data-theme', 'light')
    }
    // Fetch data (will fallback to local if API not ready)
    fetchOverview()
    fetchAgents()
    fetchTasks()
    fetchTelegram()
  }

  return {
    agents, tasks, stats, tokenChart, telegramFeed, terminalLines,
    selectedAgent, selectedTask, loading, theme, activeAgents,
    toggleTheme, selectAgent, selectTask, sendTask, init,
    fetchOverview, fetchAgents, fetchTasks,
  }
})
