import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { tasksApi, statsApi } from '@/services/api'
import { useAgentsStore } from '@/stores/agents'

function relativeTime(value) {
  if (!value) return 'just now'
  const date = new Date(value)
  const diffMs = Date.now() - date.getTime()
  const diffSeconds = Math.max(1, Math.floor(diffMs / 1000))

  if (diffSeconds < 60) return `${diffSeconds}s ago`
  const diffMinutes = Math.floor(diffSeconds / 60)
  if (diffMinutes < 60) return `${diffMinutes}m ago`
  const diffHours = Math.floor(diffMinutes / 60)
  if (diffHours < 24) return `${diffHours}h ago`
  const diffDays = Math.floor(diffHours / 24)
  return `${diffDays}d ago`
}

function durationBetween(startedAt, completedAt, fallbackUpdatedAt) {
  if (!startedAt) return 'n/a'
  const start = new Date(startedAt).getTime()
  const end = completedAt
    ? new Date(completedAt).getTime()
    : fallbackUpdatedAt
      ? new Date(fallbackUpdatedAt).getTime()
      : Date.now()
  const diff = Math.max(0, end - start)
  const totalSeconds = Math.floor(diff / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  if (minutes === 0) return `${seconds}s`
  return `${minutes}m ${seconds}s`
}

function numberOrNull(value) {
  if (typeof value === 'number' && Number.isFinite(value)) return value
  if (typeof value === 'string' && value.trim() !== '' && !Number.isNaN(Number(value))) {
    return Number(value)
  }
  return null
}

function firstTextContent(content) {
  if (typeof content === 'string') return content
  if (!Array.isArray(content)) return ''
  return content
    .map((item) => {
      if (typeof item === 'string') return item
      if (item?.type === 'text' && typeof item.text === 'string') return item.text
      return ''
    })
    .filter(Boolean)
    .join('\n')
}

function parseEventPayload(event) {
  if (!event?.payload) return {}
  if (typeof event.payload === 'object') return event.payload
  try {
    return JSON.parse(event.payload)
  } catch {
    return {}
  }
}

function findDeepValue(input, keys) {
  const targets = new Set(keys)
  const queue = [input]
  while (queue.length > 0) {
    const current = queue.shift()
    if (!current || typeof current !== 'object') continue
    for (const [key, value] of Object.entries(current)) {
      if (targets.has(key) && value !== undefined && value !== null && value !== '') {
        return value
      }
      if (value && typeof value === 'object') queue.push(value)
    }
  }
  return null
}

function eventTimestamp(event, payload) {
  return (
    findDeepValue(payload, ['timestamp', 'ts', 'endedAt', 'startedAt', 'createdAt']) ||
    event.created_at ||
    null
  )
}

function terminalLinesFromEvents(events, fallbackTask) {
  if (!events.length && fallbackTask) {
    return [
      { type: 'prompt', text: `$ queued task -> ${fallbackTask.agentName || fallbackTask.agentId}` },
      { type: 'output', text: fallbackTask.title },
    ]
  }

  const lines = []
  for (const event of [...events].reverse()) {
    const payload = parseEventPayload(event)
    if (event.event === 'task.created') {
      lines.push({ type: 'prompt', text: `$ task created -> ${fallbackTask?.title || 'untitled'}` })
      continue
    }
    if (event.event === 'task.dispatched') {
      lines.push({
        type: 'info',
        text: `→ dispatched run ${event.run_id || payload.run_id || 'unknown'} to ${event.session_key || payload.session_key || 'session'}`,
      })
      continue
    }
    if (event.event === 'agent') {
      const stream = payload.stream
      const phase = payload.data?.phase
      const text = payload.data?.text || payload.data?.delta
      if (stream === 'lifecycle' && phase === 'start') {
        lines.push({ type: 'info', text: '→ agent started processing' })
      } else if (stream === 'lifecycle' && phase === 'end') {
        lines.push({ type: 'info', text: '→ agent finished run' })
      } else if (text) {
        lines.push({ type: 'output', text })
      }
      continue
    }
    if (event.event === 'chat') {
      const role = payload.message?.role || 'assistant'
      const text = firstTextContent(payload.message?.content)
      if (text) {
        lines.push({
          type: role === 'assistant' ? 'output' : 'prompt',
          text: role === 'assistant' ? text : `$ ${text}`,
        })
      }
      continue
    }
    if (event.status === 'error') {
      lines.push({ type: 'error', text: payload.message || fallbackTask?.errorMessage || 'Task failed' })
    }
  }

  return lines.length > 0 ? lines : [{ type: 'info', text: 'No task output yet' }]
}

function normalizeTask(task) {
  return {
    id: task.id,
    title: task.title || 'Untitled task',
    agentId: task.agent_id || task.agent?.id || '',
    agentName: task.agent?.name || task.agent_name || task.agent_id || 'Unknown',
    agentEmoji: task.agent?.avatar || '◆',
    status: task.status || 'pending',
    time: relativeTime(task.updated_at || task.created_at),
    tokens: Number(task.tokens || 0),
    createdAt: task.created_at || null,
    updatedAt: task.updated_at || null,
    startedAt: task.started_at || null,
    completedAt: task.completed_at || null,
    sessionKey: task.session_key || '',
    runId: task.run_id || '',
    errorMessage: task.error_message || '',
    elapsed: durationBetween(task.started_at, task.completed_at, task.updated_at),
    agent: task.agent || null,
    raw: task,
  }
}

function buildTaskDetail(task, events) {
  const detail = {
    id: task?.id || '',
    title: task?.title || 'No task selected',
    agent: task?.agentName || 'n/a',
    agentEmoji: task?.agentEmoji || '◆',
    status: task?.status || 'idle',
    prompt: task?.title || '',
    sessionKey: task?.sessionKey || 'n/a',
    runId: task?.runId || 'n/a',
    createdAt: task?.createdAt || null,
    startedAt: task?.startedAt || null,
    completedAt: task?.completedAt || null,
    elapsed: task?.elapsed || 'n/a',
    model: 'n/a',
    provider: 'n/a',
    tokenIn: null,
    tokenOut: null,
    totalTokens: numberOrNull(task?.tokens) ?? 0,
    eventCount: events.length,
    latestMessage: '',
    errorMessage: task?.errorMessage || '',
  }

  for (const event of events) {
    const payload = parseEventPayload(event)
    const model = findDeepValue(payload, ['model', 'modelId', 'model_id'])
    const provider = findDeepValue(payload, ['provider', 'providerId', 'provider_id'])
    const tokenIn = findDeepValue(payload, ['inputTokens', 'promptTokens', 'tokensIn'])
    const tokenOut = findDeepValue(payload, ['outputTokens', 'completionTokens', 'tokensOut'])
    const totalTokens = findDeepValue(payload, ['totalTokens', 'token_usage', 'tokens'])
    const text = firstTextContent(payload.message?.content) || payload.data?.text || ''
    const timestamp = eventTimestamp(event, payload)

    if (model && detail.model === 'n/a') detail.model = String(model)
    if (provider && detail.provider === 'n/a') detail.provider = String(provider)
    if (detail.tokenIn === null) detail.tokenIn = numberOrNull(tokenIn)
    if (detail.tokenOut === null) detail.tokenOut = numberOrNull(tokenOut)
    if ((detail.totalTokens === 0 || detail.totalTokens === null) && numberOrNull(totalTokens) !== null) {
      detail.totalTokens = numberOrNull(totalTokens)
    }
    if (text) detail.latestMessage = text
    if (!detail.completedAt && payload.state === 'final' && timestamp) {
      detail.completedAt = new Date(Number(timestamp) || timestamp).toISOString()
    }
  }

  detail.elapsed = durationBetween(detail.startedAt, detail.completedAt, task?.updatedAt)
  return detail
}

export const useTasksStore = defineStore('tasks', () => {
  const tasks = ref([])
  const selectedTaskId = ref(null)
  const taskEvents = ref([])
  const loading = ref(false)
  const eventsLoading = ref(false)
  const sending = ref(false)
  const error = ref('')
  const stats = ref({
    active_agents: 0,
    total_agents: 0,
    tasks_today: 0,
    pending_tasks: 0,
    tokens_used: 0,
    success_rate: 0,
    total_tasks: 0,
    done_tasks: 0,
  })

  const selectedTask = computed(
    () => tasks.value.find((task) => task.id === selectedTaskId.value) || null
  )

  const terminalLines = computed(() => terminalLinesFromEvents(taskEvents.value, selectedTask.value))

  const selectedTaskDetail = computed(() =>
    buildTaskDetail(selectedTask.value, taskEvents.value)
  )

  const taskStats = computed(() => ({
    agent: `${selectedTaskDetail.value.agentEmoji} ${selectedTaskDetail.value.agent}`,
    status: selectedTaskDetail.value.status.toUpperCase(),
    tokensIn:
      selectedTaskDetail.value.tokenIn === null ? 'n/a' : selectedTaskDetail.value.tokenIn.toLocaleString(),
    tokensOut:
      selectedTaskDetail.value.tokenOut === null ? 'n/a' : selectedTaskDetail.value.tokenOut.toLocaleString(),
    totalTokens:
      selectedTaskDetail.value.totalTokens === null ? 'n/a' : Number(selectedTaskDetail.value.totalTokens).toLocaleString(),
    elapsed: selectedTaskDetail.value.elapsed,
    model: selectedTaskDetail.value.model,
    provider: selectedTaskDetail.value.provider,
    runId: selectedTaskDetail.value.runId,
    sessionKey: selectedTaskDetail.value.sessionKey,
    events: selectedTaskDetail.value.eventCount,
  }))

  const todayCount = computed(() => Number(stats.value.tasks_today || tasks.value.length || 0))
  const pendingCount = computed(() => Number(stats.value.pending_tasks || 0))

  async function fetchTaskEvents(id) {
    if (!id) {
      taskEvents.value = []
      return
    }

    eventsLoading.value = true
    try {
      const data = await tasksApi.getEvents(id)
      taskEvents.value = Array.isArray(data) ? data : []
    } catch (err) {
      error.value = err.message || 'Failed to fetch task events'
      taskEvents.value = []
    } finally {
      eventsLoading.value = false
    }
  }

  async function fetchTasks() {
    loading.value = true
    error.value = ''
    try {
      const data = await tasksApi.getAll()
      tasks.value = Array.isArray(data) ? data.map(normalizeTask) : []
      if (!selectedTaskId.value && tasks.value.length > 0) {
        selectedTaskId.value = tasks.value[0].id
      }
      if (
        selectedTaskId.value &&
        !tasks.value.some((task) => task.id === selectedTaskId.value)
      ) {
        selectedTaskId.value = tasks.value[0]?.id || null
      }
      await fetchTaskEvents(selectedTaskId.value)
    } catch (err) {
      error.value = err.message || 'Failed to fetch tasks'
      tasks.value = []
      taskEvents.value = []
    } finally {
      loading.value = false
    }
  }

  async function fetchStats() {
    try {
      const data = await statsApi.getOverview()
      stats.value = data || stats.value
    } catch (err) {
      error.value = err.message || 'Failed to fetch stats'
    }
  }

  async function selectTask(id) {
    selectedTaskId.value = id
    await fetchTaskEvents(id)
  }

  async function sendTask(agentId, taskText) {
    if (!taskText.trim()) return

    const agentsStore = useAgentsStore()
    const resolvedAgentId =
      agentId === 'auto'
        ? agentsStore.agents.find((agent) => agent.status === 'working')?.id ||
          agentsStore.agents[0]?.id
        : agentId

    if (!resolvedAgentId) {
      throw new Error('No agent available for this task')
    }

    sending.value = true
    error.value = ''
    try {
      const created = await tasksApi.send(resolvedAgentId, taskText.trim())
      await fetchTasks()
      if (created?.id) {
        await selectTask(created.id)
      }
    } catch (err) {
      error.value = err.message || 'Failed to send task'
      throw err
    } finally {
      sending.value = false
    }
  }

  return {
    tasks,
    selectedTaskId,
    taskEvents,
    loading,
    eventsLoading,
    sending,
    error,
    stats,
    selectedTask,
    selectedTaskDetail,
    terminalLines,
    taskStats,
    todayCount,
    pendingCount,
    fetchTasks,
    fetchTaskEvents,
    fetchStats,
    selectTask,
    sendTask,
  }
})
