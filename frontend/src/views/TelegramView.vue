<template>
  <div>
    <SectionTitle text="◈ TELEGRAM FEED" />
    <div class="tg-container px-box">
      <div class="tg-feed" ref="feedEl">
        <div
          v-for="(msg, i) in messages"
          :key="i"
          class="tg-msg"
          :style="msg.isOut ? 'flex-direction:row-reverse' : ''"
        >
          <div class="tg-avatar">{{ msg.avatar }}</div>
          <div class="tg-bubble" :class="{ out: msg.isOut }" :style="msg.bubbleStyle || ''">
            <div class="tg-name" :style="msg.isOut ? 'text-align:right' : ''">{{ msg.name }}</div>
            <span v-if="msg.highlight" :style="`color:${msg.highlight}`">{{ msg.text }}</span>
            <span v-else>{{ msg.text }}</span>
            <div class="tg-time">{{ msg.time }}</div>
          </div>
        </div>
      </div>

      <!-- Send Message -->
      <div class="tg-input-row">
        <input
          class="px-input"
          v-model="newMessage"
          placeholder="พิมพ์ข้อความ..."
          @keydown.enter="sendMessage"
          style="flex:1"
        />
        <select class="px-select" v-model="targetAgent" style="width:140px">
          <option value="all">📢 Broadcast</option>
          <option value="parae">👩‍💻 แพร</option>
          <option value="nova">🔍 NOVA</option>
          <option value="lyra">✍️ LYRA</option>
          <option value="forge">🛠️ FORGE</option>
        </select>
        <button class="px-btn px-btn-primary" @click="sendMessage">▶</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import SectionTitle from '@/components/layout/SectionTitle.vue'

const feedEl = ref(null)
const newMessage = ref('')
const targetAgent = ref('all')

const messages = ref([
  { avatar: '👤', name: 'CEO (You)', text: 'แพร ช่วยเขียน FastAPI endpoint สำหรับ Telegram webhook ด้วยนะ', time: '14:32', isOut: false },
  { avatar: '👩‍💻', name: 'แพร', text: 'รับค่ะ กำลังอ่าน codebase อยู่ จะเริ่มเขียนใน /app/routes/telegram.py เลยนะคะ', time: '14:32', isOut: true },
  { avatar: '👤', name: 'CEO (You)', text: 'NOVA หาข้อมูล best practice multi-agent routing ให้ด้วย', time: '14:35', isOut: false },
  { avatar: '🔍', name: 'NOVA', text: 'กำลังค้นหาค่ะ จะสรุปให้ใน 5 นาที', time: '14:35', isOut: true },
  { avatar: '👩‍💻', name: 'แพร', text: '▶ กำลังเขียน endpoint...', time: '14:36', isOut: true, highlight: 'var(--yellow)', bubbleStyle: 'border-color:var(--yellow)' },
])

async function sendMessage() {
  if (!newMessage.value.trim()) return
  const now = new Date()
  const time = `${now.getHours()}:${String(now.getMinutes()).padStart(2, '0')}`
  messages.value.push({
    avatar: '👤',
    name: 'CEO (You)',
    text: newMessage.value.trim(),
    time,
    isOut: false,
  })
  newMessage.value = ''
  await nextTick()
  if (feedEl.value) feedEl.value.scrollTop = feedEl.value.scrollHeight
}
</script>

<style scoped>
.tg-container { display: flex; flex-direction: column; height: calc(100vh - 200px); padding: 0; overflow: hidden; }
.tg-feed { flex: 1; overflow-y: auto; padding: 16px; display: flex; flex-direction: column; gap: 12px; }
.tg-input-row { display: flex; gap: 8px; padding: 12px 16px; border-top: 1px solid var(--border); }
.tg-msg { display: flex; gap: 8px; align-items: flex-start; }
.tg-avatar { font-size: 22px; flex-shrink: 0; }
.tg-bubble { background: rgba(255,255,255,0.04); border: 1px solid var(--border); border-radius: 8px; padding: 8px 12px; font-size: 13px; line-height: 1.6; max-width: 80%; }
.tg-bubble.out { background: rgba(124,92,252,0.08); border-color: rgba(124,92,252,0.2); }
.tg-name { font-size: 10px; color: var(--accent2); margin-bottom: 3px; font-family: 'Space Mono', monospace; }
.tg-time { font-size: 10px; color: var(--text2); margin-top: 4px; text-align: right; }
</style>
