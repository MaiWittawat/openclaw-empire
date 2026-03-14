import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { agentsApi } from '@/services/api'

const STATUS_ORDER = {
  working: 0,
  pending: 1,
  idle: 2,
  done: 3,
  error: 4,
}

function normalizeAgent(agent) {
  return {
    id: agent.id,
    name: agent.name || agent.id,
    role: agent.role || 'OpenClaw Agent',
    emoji: agent.avatar || '◆',
    color: agent.color || '#0ea5e9',
    status: agent.status || 'idle',
    tokens: Number(agent.tokens || 0),
    tokenMax: Number(agent.max_tokens || 20000) || 20000,
    sessionKey: agent.session_key || '',
    lastSeenAt: agent.last_seen_at || null,
    updatedAt: agent.updated_at || null,
    createdAt: agent.created_at || null,
    raw: agent,
  }
}

export const useAgentsStore = defineStore('agents', () => {
  const agents = ref([])
  const selectedAgentId = ref(null)
  const loading = ref(false)
  const error = ref('')

  const sortedAgents = computed(() =>
    [...agents.value].sort((left, right) => {
      const statusDelta = (STATUS_ORDER[left.status] ?? 99) - (STATUS_ORDER[right.status] ?? 99)
      if (statusDelta !== 0) return statusDelta
      return left.name.localeCompare(right.name)
    })
  )

  const selectedAgent = computed(() =>
    sortedAgents.value.find((agent) => agent.id === selectedAgentId.value) || null
  )

  const activeCount = computed(
    () => sortedAgents.value.filter((agent) => agent.status === 'working').length
  )

  function selectAgent(id) {
    selectedAgentId.value = id
  }

  async function fetchAgents() {
    loading.value = true
    error.value = ''
    try {
      const data = await agentsApi.getAll()
      agents.value = Array.isArray(data) ? data.map(normalizeAgent) : []
      if (!selectedAgentId.value && agents.value.length > 0) {
        selectedAgentId.value = agents.value[0].id
      }
      if (
        selectedAgentId.value &&
        !agents.value.some((agent) => agent.id === selectedAgentId.value)
      ) {
        selectedAgentId.value = agents.value[0]?.id || null
      }
    } catch (err) {
      error.value = err.message || 'Failed to fetch agents'
      agents.value = []
    } finally {
      loading.value = false
    }
  }

  return {
    agents: sortedAgents,
    selectedAgent,
    selectedAgentId,
    activeCount,
    loading,
    error,
    selectAgent,
    fetchAgents,
  }
})
