import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { tasksApi } from '@/services/api'

export const useTasksStore = defineStore('tasks', () => {
  const tasks = ref([
    {
      id: 'task1',
      title: 'เขียน FastAPI endpoint สำหรับ Telegram webhook',
      agent: 'แพร',
      agentEmoji: '👩‍💻',
      status: 'working',
      time: '2m ago',
      tokens: 12450,
    },
    {
      id: 'task2',
      title: 'ค้นหา best practice สำหรับ multi-agent routing',
      agent: 'NOVA',
      agentEmoji: '🔍',
      status: 'working',
      time: '5m ago',
      tokens: 8820,
    },
    {
      id: 'task3',
      title: 'เขียน README สำหรับ project',
      agent: 'LYRA',
      agentEmoji: '✍️',
      status: 'done',
      time: '18m ago',
      tokens: 5400,
    },
    {
      id: 'task4',
      title: 'สร้าง Docker Compose config',
      agent: 'FORGE',
      agentEmoji: '🛠️',
      status: 'done',
      time: '1h ago',
      tokens: 3200,
    },
    {
      id: 'task5',
      title: 'Deploy to VPS — timeout error',
      agent: 'FORGE',
      agentEmoji: '🛠️',
      status: 'error',
      time: '2h ago',
      tokens: 1100,
    },
  ])

  const selectedTaskId = ref('task1')
  const terminalLines = ref([
    { type: 'prompt', text: '$ task received: FastAPI webhook endpoint' },
    { type: 'info', text: '→ reading project structure...' },
    { type: 'output', text: 'found: /app/main.py, /app/routes/' },
    { type: 'prompt', text: '$ writing /app/routes/telegram.py' },
    { type: 'output', text: '@router.post("/webhook")' },
    { type: 'output', text: 'async def telegram_webhook(req: Request):' },
    { type: 'output', text: '  data = await req.json()' },
    { type: 'output', text: '  # route to agent...' },
    { type: 'info', text: '→ adding error handling...' },
    { type: 'cursor', text: '' },
  ])

  const taskStats = ref({
    agent: '👩‍💻 แพร',
    status: 'WORKING',
    tokensIn: '4,200',
    tokensOut: '8,250',
    cost: '~$0.018',
    elapsed: '2m 14s',
  })

  const todayCount = computed(() => tasks.value.length)
  const pendingCount = computed(
    () => tasks.value.filter((t) => t.status === 'working').length
  )

  function selectTask(id) {
    selectedTaskId.value = id
  }

  async function sendTask(agentId, taskText) {
    const newTask = {
      id: `task_${Date.now()}`,
      title: taskText,
      agent: agentId,
      agentEmoji: '⚡',
      status: 'working',
      time: 'just now',
      tokens: 0,
    }
    tasks.value.unshift(newTask)
    terminalLines.value.push({ type: 'prompt', text: `$ new task → ${agentId}: ${taskText}` })

    try {
      await tasksApi.send(agentId, taskText)
    } catch {
      // Optimistic update, handle error silently
    }
  }

  async function fetchTasks() {
    try {
      const data = await tasksApi.getAll()
      if (data?.length) tasks.value = data
    } catch {
      // Use mock data if backend not available
    }
  }

  return {
    tasks,
    selectedTaskId,
    terminalLines,
    taskStats,
    todayCount,
    pendingCount,
    selectTask,
    sendTask,
    fetchTasks,
  }
})
