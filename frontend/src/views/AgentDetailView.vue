<template>
  <div>
    <SectionTitle :text="`◈ AGENT — ${route.params.id?.toUpperCase()}`" />
    <div class="px-box" style="padding: 24px; text-align: center; color: var(--text2)">
      <div style="font-size: 48px; margin-bottom: 12px">🤖</div>
      <div style="font-family: 'Space Mono', monospace; font-size: 12px">
        Detailed agent view for <span style="color: var(--accent2)">{{ route.params.id }}</span>
      </div>
      <RouterLink to="/agents" style="color: var(--accent2); font-size: 11px; margin-top: 12px; display: block">
        ← Back to Agents
      </RouterLink>
    </div>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import SectionTitle from '@/components/layout/SectionTitle.vue'
const route = useRoute()
</script>
