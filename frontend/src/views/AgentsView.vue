<template>
  <div>
    <SectionTitle text="◈ ALL AGENTS" />
    <div class="agents-grid">
      <AgentCard
        v-for="agent in agentsStore.agents"
        :key="agent.id"
        :agent="agent"
        :isSelected="agentsStore.selectedAgentId === agent.id"
        @select="agentsStore.selectAgent"
      />
    </div>

    <!-- Agent Detail -->
    <div v-if="agentsStore.selectedAgent" style="margin-top: 20px">
      <SectionTitle :text="`◈ ${agentsStore.selectedAgent.name} — DETAIL`" />
      <div class="px-box" style="padding: 20px">
        <div class="detail-grid">
          <div>
            <div class="detail-label">AGENT ID</div>
            <div class="detail-value">{{ agentsStore.selectedAgent.id }}</div>
          </div>
          <div>
            <div class="detail-label">ROLE</div>
            <div class="detail-value">{{ agentsStore.selectedAgent.role }}</div>
          </div>
          <div>
            <div class="detail-label">STATUS</div>
            <div class="detail-value" :style="statusStyle">{{ agentsStore.selectedAgent.status.toUpperCase() }}</div>
          </div>
          <div>
            <div class="detail-label">TOKENS USED</div>
            <div class="detail-value">{{ agentsStore.selectedAgent.tokens.toLocaleString() }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import SectionTitle from '@/components/layout/SectionTitle.vue'
import AgentCard from '@/components/agents/AgentCard.vue'
import { useAgentsStore } from '@/stores/agents'

const agentsStore = useAgentsStore()

const statusStyle = computed(() => {
  const map = { working: 'color:var(--yellow)', done: 'color:var(--green)', error: 'color:var(--red)', idle: 'color:var(--text2)' }
  return map[agentsStore.selectedAgent?.status] || ''
})
</script>

<style scoped>
.agents-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}
.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}
.detail-label { font-family: 'Space Mono', monospace; font-size: 9px; color: var(--text2); letter-spacing: 2px; margin-bottom: 4px; }
.detail-value { font-size: 14px; font-weight: 600; color: var(--text); }
</style>
