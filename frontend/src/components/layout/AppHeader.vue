<template>
  <header class="app-header">
    <div class="logo">AI <span>Empire</span></div>
    <div class="header-right">
      <div class="ticker">[ SYS:ONLINE ] &nbsp;[ TASKS:{{ tasksStore.tasks.length }} ] &nbsp;[ AGENTS:{{ agentsStore.agents.length }} ]</div>

      <!-- Theme Toggle -->
      <div class="theme-toggle-wrap" @click="themeStore.toggle()" title="Toggle light/dark">
        <span class="theme-icon">{{ themeStore.isDark ? '🌙' : '☀️' }}</span>
        <div class="toggle-track">
          <div class="toggle-thumb" :class="{ light: !themeStore.isDark }"></div>
        </div>
      </div>

      <!-- Status Pill -->
      <div class="status-pill">
        <div class="status-dot"></div>
        GATEWAY OK
      </div>
    </div>
  </header>
</template>

<script setup>
import { useThemeStore } from '@/stores/theme'
import { useAgentsStore } from '@/stores/agents'
import { useTasksStore } from '@/stores/tasks'

const themeStore = useThemeStore()
const agentsStore = useAgentsStore()
const tasksStore = useTasksStore()
</script>

<style scoped>
.app-header {
  background: var(--header-bg);
  border-bottom: 1px solid var(--border);
  backdrop-filter: blur(24px);
  padding: 12px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  top: 0;
  z-index: 100;
  transition: background 0.3s;
}

.logo {
  font-family: 'Inter', sans-serif;
  font-weight: 700;
  font-size: 16px;
  color: var(--text);
  letter-spacing: -0.5px;
}
.logo span { color: var(--accent2); }

.header-right {
  display: flex;
  align-items: center;
  gap: 14px;
}

.ticker {
  font-family: 'Space Mono', monospace;
  color: var(--text2);
  font-size: 10px;
}

.theme-toggle-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}
.theme-icon { font-size: 13px; }
.toggle-track {
  width: 38px; height: 21px;
  background: var(--border2);
  border-radius: 20px;
  position: relative;
  border: 1px solid var(--border);
  transition: background 0.3s;
}
[data-theme="light"] .toggle-track { background: var(--accent); }
.toggle-thumb {
  position: absolute;
  width: 15px; height: 15px;
  border-radius: 50%;
  background: #fff;
  top: 2px; left: 2px;
  transition: transform 0.25s cubic-bezier(.4,0,.2,1);
  box-shadow: 0 1px 4px rgba(0,0,0,0.3);
}
.toggle-thumb.light { transform: translateX(17px); }

.status-pill {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 500;
  color: var(--green);
  background: rgba(52,211,153,0.1);
  border: 1px solid rgba(52,211,153,0.2);
  border-radius: 20px;
  padding: 4px 10px;
}
.status-dot {
  width: 6px; height: 6px;
  background: var(--green);
  border-radius: 50%;
  box-shadow: 0 0 5px var(--green);
  animation: pulse-dot 2s ease-in-out infinite;
}
</style>
