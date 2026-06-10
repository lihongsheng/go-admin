import { ref, watchEffect } from 'vue'

const DARK_KEY = 'go-admin-dark'

// 读取 localStorage 或跟随系统偏好
function getInitial() {
  const stored = localStorage.getItem(DARK_KEY)
  if (stored !== null) return stored === 'true'
  return window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
}

const isDark = ref(getInitial())

// 同步 html.dark class + localStorage
watchEffect(() => {
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem(DARK_KEY, String(isDark.value))
})

export function useDarkMode() {
  function toggle() {
    isDark.value = !isDark.value
  }
  return { isDark, toggle }
}
