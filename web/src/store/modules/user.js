import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, info, logout } from '@/api/base'

const TOKEN_KEY = 'go-admin-token'

export const useUserStore = defineStore('user', () => {
  // state
  const token = ref(localStorage.getItem(TOKEN_KEY) || '')
  const userInfo = ref(null)

  // getters
  const isLoggedIn = computed(() => !!token.value)

  // actions (mutations + actions merged)
  function setToken(t) {
    token.value = t
    if (t) localStorage.setItem(TOKEN_KEY, t)
    else localStorage.removeItem(TOKEN_KEY)
  }

  async function doLogin(payload) {
    const { data } = await login(payload)
    setToken(data.token)
    // 不在这里设置 userInfo，让路由守卫的 fetchInfo() + generateRoutes() 统一处理
    // 否则首次登录后 userInfo 已存在，generateRoutes() 被跳过，菜单栏为空
  }

  async function fetchInfo() {
    const { data } = await info()
    userInfo.value = data
    return data
  }

  async function doLogout() {
    // 主动退出时调 logout 通知后端；token 已过期时不调用，避免触发 401 循环
    if (token.value) {
      try { await logout() } catch (_) { /* ignore */ }
    }
    setToken('')
    userInfo.value = null
  }

  return {
    token, userInfo, isLoggedIn,
    setToken, doLogin, fetchInfo, doLogout
  }
})
