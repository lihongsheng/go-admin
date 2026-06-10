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
    userInfo.value = data.user
  }

  async function fetchInfo() {
    const { data } = await info()
    userInfo.value = data
    return data
  }

  async function doLogout() {
    try { await logout() } catch (_) { /* ignore */ }
    setToken('')
    userInfo.value = null
  }

  return {
    token, userInfo, isLoggedIn,
    setToken, doLogin, fetchInfo, doLogout
  }
})
