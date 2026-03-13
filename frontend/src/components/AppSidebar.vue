<template>
  <nav class="sidebar">
    <div class="sidebar-label">NAVIGATION</div>

    <RouterLink
      v-for="item in navItems"
      :key="item.to"
      :to="item.to"
      class="nav-item"
      :class="{ active: route.path === item.to }"
    >
      <span class="nav-icon">{{ item.icon }}</span>
      {{ item.label }}
      <span v-if="item.badge" class="nav-badge" :style="item.badgeStyle">
        {{ item.badge }}
      </span>
    </RouterLink>

    <div class="sidebar-label" style="margin-top:16px">SYSTEM</div>

    <RouterLink to="/memory" class="nav-item" :class="{ active: route.path === '/memory' }">
      <span class="nav-icon">◉</span> MEMORY
    </RouterLink>
    <RouterLink to="/settings" class="nav-item" :class="{ active: route.path === '/settings' }">
      <span class="nav-icon">⚙</span> SETTINGS
    </RouterLink>

    <div class="sidebar-footer">
      <div class="sys-stats">
        <div>MEM <span style="color:var(--green)">2.4 GB</span></div>
        <div>CPU <span style="color:var(--yellow)">34%</span></div>
        <div>UPT <span style="color:var(--cyan)">2d 14h</span></div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { RouterLink, useRoute } from 'vue-router'

const route = useRoute()
const navItems = [
  { to: '/dashboard', icon: '⊞', label: 'DASHBOARD' },
  { to: '/agents', icon: '◈', label: 'AGENTS', badge: '4' },
  { to: '/tasks', icon: '▤', label: 'TASKS', badge: '12' },
  { to: '/telegram', icon: '✉', label: 'TELEGRAM', badge: '3', badgeStyle: 'background:var(--red)' },
]
</script>

<style scoped>
.sidebar {
  background: var(--sidebar-bg);
  border-right: 1px solid var(--border);
  padding: 16px 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
  backdrop-filter: blur(20px);
  transition: background 0.3s;
}
.sidebar-label {
  font-family: 'Space Mono', monospace;
  font-size: 9px;
  color: var(--text2);
  padding: 10px 16px 4px;
  letter-spacing: 2px;
  text-transform: uppercase;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 16px;
  cursor: pointer;
  color: var(--text2);
  font-size: 12px;
  transition: all 0.15s;
  border-left: 2px solid transparent;
  border-radius: 0 6px 6px 0;
  margin-right: 8px;
  text-decoration: none;
}
.nav-item:hover { color: var(--text); background: rgba(255,255,255,0.04); }
.nav-item.active {
  color: var(--accent2);
  background: rgba(124,92,252,0.1);
  border-left-color: var(--accent);
}
.nav-icon { font-size: 14px; width: 18px; text-align: center; }
.nav-badge {
  margin-left: auto;
  background: rgba(124,92,252,0.2);
  color: var(--accent2);
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 10px;
  font-family: 'Space Mono', monospace;
  border: 1px solid rgba(124,92,252,0.3);
}
.sidebar-footer {
  margin-top: auto;
  padding: 16px;
  border-top: 1px solid var(--border);
}
.sys-stats {
  font-size: 10px;
  color: var(--text2);
  line-height: 2;
  font-family: 'Space Mono', monospace;
}
</style>
