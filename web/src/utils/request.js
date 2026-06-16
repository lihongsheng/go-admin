import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import router from '@/router'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000
})

request.interceptors.request.use(cfg => {
  const userStore = useUserStore()
  if (userStore.token) {
    cfg.headers.Authorization = 'Bearer ' + userStore.token
  }
  return cfg
})

let isLoggingOut = false

request.interceptors.response.use(
  res => {
    const data = res.data
    if (!data || typeof data !== 'object') return data
    if (data.code === 0) return data
    if (data.code === 9001) {
      if (router.currentRoute.value.path !== '/install') {
        router.replace('/install')
      }
      return Promise.reject(data)
    }
    if (data.code === 401) {
      if (!isLoggingOut) {
        isLoggingOut = true
        const userStore = useUserStore()
        // 先清 token，避免 logout() 的 401 再次触发本拦截器
        userStore.setToken('')
        userStore.userInfo = null
        router.replace('/login')
        // 延迟重置标志，防止同一批次其他 401 重复跳转
        setTimeout(() => { isLoggingOut = false }, 1000)
      }
      return Promise.reject(data)
    }
    ElMessage.error(data.msg || 'error')
    return Promise.reject(data)
  },
  err => {
    ElMessage.error(err.message || 'network error')
    return Promise.reject(err)
  }
)

export default request
