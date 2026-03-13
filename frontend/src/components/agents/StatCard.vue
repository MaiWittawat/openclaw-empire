<template>
  <div class="stat-card px-box" :style="`--color:${color}`">
    <span class="stat-icon">{{ icon }}</span>
    <div class="stat-label">{{ label }}</div>
    <div class="stat-value">{{ value }}<sup>{{ sup }}</sup></div>
    <div class="stat-sub" v-html="sub"></div>
  </div>
</template>

<script setup>
defineProps({
  icon: String,
  label: String,
  value: [String, Number],
  sup: String,
  sub: String,
  color: String,
})
</script>

<style scoped>
.stat-card {
  padding: 18px 16px 14px;
  position: relative;
  overflow: hidden;
}
.stat-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--color, var(--accent)), transparent);
  border-radius: 12px 12px 0 0;
}
.stat-icon { font-size: 18px; margin-bottom: 10px; display: block; }
.stat-label { font-size: 10px; font-weight: 500; color: var(--text2); margin-bottom: 6px; letter-spacing: 0.8px; text-transform: uppercase; }
.stat-value { font-family: 'Inter', sans-serif; font-weight: 700; font-size: 32px; color: var(--text); line-height: 1; letter-spacing: -1.5px; }
.stat-value sup { font-size: 14px; font-weight: 500; color: var(--color, var(--accent2)); letter-spacing: 0; vertical-align: super; }
.stat-sub { font-size: 11px; font-weight: 400; color: var(--text2); margin-top: 6px; display: flex; align-items: center; gap: 4px; }
</style>
