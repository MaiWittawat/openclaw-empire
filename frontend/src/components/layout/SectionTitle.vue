<template>
  <div class="section-title">{{ text }}</div>
</template>

<script setup>
defineProps({ text: String })
</script>

<style scoped>
.section-title {
  font-family: 'Space Mono', monospace;
  font-size: 9px;
  color: var(--text2);
  letter-spacing: 3px;
  text-transform: uppercase;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 10px;
}
.section-title::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border);
}
</style>
