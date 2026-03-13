import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { agentsApi } from '@/services/api'

export const useAgentsStore = defineStore('agents', () => {
  const agents = ref([
    {
      id: 'parae',
      name: 'แพร',
      role: 'Tech Lead',
      emoji: '👩‍💻',
      color: '#a78bfa',
      status: 'working',
      tokens: 12450,
      tokenMax: 20000,
    },
    {
      id: 'nova',
      name: 'NOVA',
      role: 'Researcher',
      emoji: '🔍',
      color: '#34d399',
      status: 'working',
      tokens: 8820,
      tokenMax: 20000,
    },
    {
      id: 'forge',
      name: 'FORGE',
      role: 'DevOps',
      emoji: '🛠️',
      color: '#60a5fa',
      status: 'idle',
      tokens: 3200,
      tokenMax: 20000,
    },
    {
      id: 'lyra',
      name: 'LYRA',
      role: 'Writer',
      emoji: '✍️',
      color: '#f472b6',
      status: 'done',
      tokens: 5400,
      tokenMax: 20000,
    },
  ])

  const selectedAgentId = ref(null)

  const selectedAgent = computed(() =>
    agents.value.find((a) => a.id === selectedAgentId.value)
  )

  const activeCount = computed(
    () => agents.value.filter((a) => a.status === 'working').length
  )

  function selectAgent(id) {
    selectedAgentId.value = id
  }

  async function fetchAgents() {
    try {
      const data = await agentsApi.getAll()
      if (data?.length) agents.value = data
    } catch {
      // Use mock data if backend not available
    }
  }

  return { agents, selectedAgent, selectedAgentId, activeCount, selectAgent, fetchAgents }
})
