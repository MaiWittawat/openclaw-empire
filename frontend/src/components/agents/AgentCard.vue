<template>
  <div
    class="agent-card px-box"
    :class="[agent.status, { selected: isSelected }]"
    :style="`--agent-color:${agent.color}`"
    @click="$emit('select', agent.id)"
  >
    <div class="agent-avatar-wrap">
      <div class="px-char">{{ agent.emoji }}</div>
    </div>
    <div class="agent-header">
      <div class="agent-info">
        <div class="agent-name">{{ agent.name }}</div>
        <div class="agent-role">{{ agent.role }}</div>
      </div>
    </div>
    <div class="agent-status" :class="agent.status">
      <span v-if="agent.status === 'working'" class="status-blink">▶</span>
      <span v-else-if="agent.status === 'done'">✔</span>
      <span v-else-if="agent.status === 'error'">✖</span>
      <span v-else>◻</span>
      {{ statusLabel }}
    </div>
    <div class="token-bar-wrap">
      <div class="token-bar-label">
        <span>tokens</span>
        <span :style="tokenColor">{{ agent.tokens.toLocaleString() }}</span>
      </div>
      <div class="token-bar">
        <div
          class="token-bar-fill"
          :class="tokenFillClass"
          :style="`width:${tokenPercent}%;background:${agent.status !== 'working' ? `linear-gradient(90deg,${agent.color},${agent.color}88)` : ''}`"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  agent: { type: Object, required: true },
  isSelected: { type: Boolean, default: false },
})
defineEmits(['select'])

const statusLabel = computed(() => props.agent.status.toUpperCase())
const tokenPercent = computed(() => Math.round((props.agent.tokens / props.agent.tokenMax) * 100))
const tokenColor = computed(() => {
  if (tokenPercent.value > 80) return 'color:var(--red)'
  if (tokenPercent.value > 50) return 'color:var(--yellow)'
  return 'color:var(--green)'
})
const tokenFillClass = computed(() => {
  if (tokenPercent.value > 80) return 'critical'
  if (tokenPercent.value > 50) return 'high'
  return ''
})
</script>

<style scoped>
.agent-card {
  padding: 14px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  overflow: hidden;
}
.agent-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--agent-color, var(--accent)), transparent);
  border-radius: 10px 10px 0 0;
}
.agent-card:hover {
  border-color: rgba(167,139,250,0.3);
  background: rgba(255,255,255,0.05);
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.3);
}
.agent-card.selected {
  border-color: rgba(167,139,250,0.4);
  background: rgba(124,92,252,0.06);
}
.agent-card.working {
  box-shadow: 0 4px 20px rgba(251,191,36,0.1), 0 0 0 1px rgba(251,191,36,0.15);
}

.agent-avatar-wrap {
  display: flex;
  justify-content: center;
  margin-bottom: 8px;
}
.agent-card.working .agent-avatar-wrap { animation: walk 0.6s ease-in-out infinite; }
.px-char { font-size: 32px; line-height: 1; }

.agent-header { display: flex; align-items: flex-start; gap: 10px; margin-bottom: 10px; }
.agent-info { flex: 1; min-width: 0; }
.agent-name { font-family: 'Syne', sans-serif; font-weight: 700; font-size: 14px; color: var(--text); margin-bottom: 2px; }
.agent-role { font-size: 10px; color: var(--text2); letter-spacing: 1px; text-transform: uppercase; }

.agent-status {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  padding: 3px 9px;
  border-radius: 20px;
  margin-bottom: 10px;
  font-family: 'Space Mono', monospace;
  letter-spacing: 0.5px;
}
.agent-status.idle   { color: var(--text2); background: rgba(255,255,255,0.04); border: 1px solid var(--border); }
.agent-status.working{ color: var(--yellow); background: rgba(251,191,36,0.1);  border: 1px solid rgba(251,191,36,0.3); }
.agent-status.done   { color: var(--green);  background: rgba(52,211,153,0.1);  border: 1px solid rgba(52,211,153,0.3); }
.agent-status.error  { color: var(--red);    background: rgba(248,113,113,0.1); border: 1px solid rgba(248,113,113,0.3); }

.status-blink { animation: blink 1s step-end infinite; }

.token-bar-wrap { margin-top: 6px; }
.token-bar-label { display: flex; justify-content: space-between; font-size: 10px; color: var(--text2); margin-bottom: 5px; }
.token-bar { height: 4px; background: rgba(255,255,255,0.06); border-radius: 4px; overflow: hidden; }
.token-bar-fill { height: 100%; background: linear-gradient(90deg, var(--accent), var(--accent2)); border-radius: 4px; transition: width 0.5s ease; }
.token-bar-fill.high { background: linear-gradient(90deg, #f59e0b, #fbbf24); }
.token-bar-fill.critical { background: linear-gradient(90deg, var(--red), #fca5a5); }
</style>
