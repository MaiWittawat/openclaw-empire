<template>
  <div>
    <!-- STATS -->
    <div>
      <SectionTitle text="◈ OVERVIEW" />
      <div class="stats-row">
        <StatCard
          icon="🟢"
          label="ACTIVE AGENTS"
          :value="agentsStore.activeCount"
          :sup="`/${agentsStore.agents.length}`"
          sub='<span class="up">↑</span> 1 just started'
          color="var(--green)"
        />
        <StatCard
          icon="📋"
          label="TASKS TODAY"
          :value="tasksStore.todayCount"
          sup=""
          :sub='`<span style="color:var(--yellow)">●</span> ${tasksStore.pendingCount} pending`'
          color="var(--cyan)"
        />
        <StatCard
          icon="⚡"
          label="TOKENS USED"
          value="284"
          sup="K"
          sub="≈ $0.42 today"
          color="var(--yellow)"
        />
        <StatCard
          icon="✅"
          label="SUCCESS RATE"
          value="94"
          sup="%"
          sub='<span class="up">↑</span> 16 of 17 done'
          color="var(--accent2)"
        />
      </div>
    </div>

    <!-- AGENTS -->
    <div>
      <SectionTitle text="◈ AGENTS" />
      <div class="agents-grid">
        <AgentCard
          v-for="agent in agentsStore.agents"
          :key="agent.id"
          :agent="agent"
          :isSelected="agentsStore.selectedAgentId === agent.id"
          @select="agentsStore.selectAgent"
        />
      </div>
    </div>

    <!-- TASK HISTORY -->
    <div>
      <SectionTitle text="◈ TASK HISTORY" />
      <div class="task-list">
        <TaskItem
          v-for="task in tasksStore.tasks"
          :key="task.id"
          :task="task"
          :isActive="tasksStore.selectedTaskId === task.id"
          @select="tasksStore.selectTask"
        />
      </div>
    </div>

    <!-- TOKEN CHART -->
    <div>
      <SectionTitle text="◈ TOKEN USAGE (7 DAYS)" />
      <div class="px-box" style="padding:16px">
        <div class="token-legend">
          <span>👩‍💻 แพร: <span style="color:var(--accent2)">89K</span></span>
          <span>🔍 NOVA: <span style="color:var(--cyan)">72K</span></span>
          <span>✍️ LYRA: <span style="color:var(--yellow)">61K</span></span>
          <span>🛠️ FORGE: <span style="color:var(--green)">48K</span></span>
        </div>
        <div class="mini-chart">
          <div class="chart-bar" style="height:40%;background:var(--border2)"></div>
          <div class="chart-bar" style="height:55%"></div>
          <div class="chart-bar" style="height:35%;background:var(--border2)"></div>
          <div class="chart-bar" style="height:70%"></div>
          <div class="chart-bar" style="height:60%;background:var(--border2)"></div>
          <div class="chart-bar" style="height:85%"></div>
          <div class="chart-bar" style="height:100%;background:var(--accent2);box-shadow:0 0 8px var(--accent2)"></div>
        </div>
        <div class="chart-labels">
          <span>Mon</span><span>Tue</span><span>Wed</span><span>Thu</span><span>Fri</span><span>Sat</span>
          <span style="color:var(--accent2)">Today</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import SectionTitle from '@/components/layout/SectionTitle.vue'
import StatCard from '@/components/agents/StatCard.vue'
import AgentCard from '@/components/agents/AgentCard.vue'
import TaskItem from '@/components/tasks/TaskItem.vue'
import { useAgentsStore } from '@/stores/agents'
import { useTasksStore } from '@/stores/tasks'

const agentsStore = useAgentsStore()
const tasksStore = useTasksStore()

onMounted(() => {
  agentsStore.fetchAgents()
  tasksStore.fetchTasks()
})
</script>

<style scoped>
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}
.agents-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}
.task-list { display: flex; flex-direction: column; gap: 6px; }

.token-legend {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
  font-size: 11px;
  color: var(--text2);
  font-family: 'Space Mono', monospace;
  flex-wrap: wrap;
}
.mini-chart {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  height: 50px;
  padding: 0 2px;
}
.chart-bar {
  flex: 1;
  background: rgba(124,92,252,0.35);
  border-radius: 3px 3px 0 0;
  transition: all 0.3s;
  min-width: 8px;
}
.chart-bar:hover { background: rgba(167,139,250,0.7); }
.chart-labels {
  display: flex;
  justify-content: space-between;
  font-size: 9px;
  color: var(--text2);
  margin-top: 8px;
  padding: 0 2px;
  font-family: 'Space Mono', monospace;
}

:deep(.up) { color: var(--green); }
:deep(.down) { color: var(--red); }
</style>
