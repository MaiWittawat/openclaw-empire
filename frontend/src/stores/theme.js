import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(true)

  function init() {
    const saved = localStorage.getItem('theme')
    isDark.value = saved !== 'light'
    applyTheme()
  }

  function toggle() {
    isDark.value = !isDark.value
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
    applyTheme()
  }

  function applyTheme() {
    document.documentElement.setAttribute(
      'data-theme',
      isDark.value ? 'dark' : 'light'
    )
  }

  return { isDark, init, toggle }
})
