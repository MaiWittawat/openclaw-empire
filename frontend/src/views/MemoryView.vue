<template>
  <div>
    <SectionTitle text="◈ AGENT MEMORY" />

    <div class="search-bar">
      <input class="px-input" v-model="query" placeholder="🔍 ค้นหา memory..." />
    </div>

    <div class="memory-list">
      <div v-for="item in filteredMemory" :key="item.id" class="memory-item px-box">
        <div class="memory-header">
          <strong>{{ item.agent }}</strong>
          <span class="memory-time">{{ item.time }}</span>
        </div>
        <div class="memory-key">{{ item.key }}</div>
        <div class="memory-value">{{ item.value }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import SectionTitle from '@/components/layout/SectionTitle.vue'

const query = ref('')

const memory = ref([
  { id: 1, agent: '👩‍💻 แพร', key: 'PROJECT STACK', value: 'FastAPI + PostgreSQL + pgvector + Vue 3 + Tailwind', time: '2h ago' },
  { id: 2, agent: '🔍 NOVA', key: 'RESEARCH TOPIC', value: 'Multi-agent routing best practices — LangGraph, AutoGen', time: '5h ago' },
  { id: 3, agent: '🛠️ FORGE', key: 'INFRA', value: 'Docker Compose: postgres, pgvector, backend, frontend', time: '1d ago' },
  { id: 4, agent: '✍️ LYRA', key: 'TONE', value: 'Technical but approachable, Thai/English mixed, concise', time: '1d ago' },
  { id: 5, agent: '👩‍💻 แพร', key: 'API STRUCTURE', value: '/agents, /tasks, /stats, /telegram, /memory endpoints', time: '3h ago' },
])

const filteredMemory = computed(() =>
  query.value
    ? memory.value.filter(m =>
        m.key.toLowerCase().includes(query.value.toLowerCase()) ||
        m.value.toLowerCase().includes(query.value.toLowerCase())
      )
    : memory.value
)
</script>

<style scoped>
.search-bar { margin-bottom: 16px; }
.memory-list { display: flex; flex-direction: column; gap: 10px; }
.memory-item { padding: 14px 16px; }
.memory-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px; }
.memory-header strong { font-size: 11px; color: var(--accent2); font-family: 'Space Mono', monospace; letter-spacing: 1px; }
.memory-time { font-size: 10px; color: var(--text2); font-family: 'Space Mono', monospace; }
.memory-key { font-size: 10px; color: var(--text2); letter-spacing: 1.5px; text-transform: uppercase; margin-bottom: 4px; font-family: 'Space Mono', monospace; border-left: 2px solid var(--accent); padding-left: 8px; }
.memory-value { font-size: 13px; color: var(--text); line-height: 1.6; }
</style>
