<template>
  <div>
    <SectionTitle text="◈ SETTINGS" />

    <div class="settings-grid">
      <!-- Backend -->
      <div class="settings-section px-box">
        <div class="settings-title">🔌 BACKEND CONNECTION</div>
        <div class="settings-form">
          <div class="setting-row">
            <label>API Base URL</label>
            <input class="px-input" v-model="settings.apiUrl" placeholder="http://localhost:8000" />
          </div>
          <div class="setting-row">
            <label>API Key</label>
            <input class="px-input" v-model="settings.apiKey" type="password" placeholder="sk-..." />
          </div>
          <button class="px-btn px-btn-primary" @click="saveSettings">💾 SAVE</button>
        </div>
      </div>

      <!-- Theme -->
      <div class="settings-section px-box">
        <div class="settings-title">🎨 APPEARANCE</div>
        <div class="settings-form">
          <div class="setting-row">
            <label>Theme</label>
            <button class="px-btn px-btn-ghost" @click="themeStore.toggle()">
              {{ themeStore.isDark ? '🌙 Dark' : '☀️ Light' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Agents Config -->
      <div class="settings-section px-box">
        <div class="settings-title">🤖 AGENT CONFIG</div>
        <div class="settings-form">
          <div class="setting-row">
            <label>Max Tokens / Task</label>
            <input class="px-input" v-model="settings.maxTokens" type="number" />
          </div>
          <div class="setting-row">
            <label>Auto Route</label>
            <input type="checkbox" v-model="settings.autoRoute" style="width:auto;accent-color:var(--accent)" />
          </div>
        </div>
      </div>
    </div>

    <div v-if="saved" class="save-toast">✅ Settings saved!</div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import SectionTitle from '@/components/layout/SectionTitle.vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const saved = ref(false)

const settings = ref({
  apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8000',
  apiKey: '',
  maxTokens: 20000,
  autoRoute: true,
})

function saveSettings() {
  localStorage.setItem('settings', JSON.stringify(settings.value))
  saved.value = true
  setTimeout(() => saved.value = false, 2000)
}
</script>

<style scoped>
.settings-grid { display: flex; flex-direction: column; gap: 16px; }
.settings-section { padding: 20px; }
.settings-title { font-family: 'Space Mono', monospace; font-size: 10px; color: var(--accent2); letter-spacing: 2px; margin-bottom: 16px; }
.settings-form { display: flex; flex-direction: column; gap: 12px; }
.setting-row { display: flex; flex-direction: column; gap: 4px; }
.setting-row label { font-size: 10px; color: var(--text2); font-family: 'Space Mono', monospace; letter-spacing: 1px; }
.save-toast {
  position: fixed; bottom: 24px; right: 24px;
  background: rgba(52,211,153,0.15);
  border: 1px solid var(--green);
  color: var(--green);
  padding: 10px 20px;
  border-radius: 8px;
  font-family: 'Space Mono', monospace;
  font-size: 12px;
}
</style>
