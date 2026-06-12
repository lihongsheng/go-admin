import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import router from '@/router'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000
})

request.interceptors.request.use(cfg => {
  // Pinia store 可在拦截器运行时直接使用（此时 app 已挂载）
  const userStore = useUserStore()
  if (userStore.token) {
    cfg.headers.Authorization = 'Bearer ' + userStore.token
  }
  return cfg
})

request.interceptors.response.use(
  res => {
    const data = res.data
    if (!data || typeof data !== 'object') return data
    if (data.code === 0) return data
    if (data.code === 9001) {
      // 未安装 → 跳安装向导
      if (router.currentRoute.value.path !== '/install') {
        router.replace('/install')
      }
      return Promise.reject(data)
    }
    if (data.code === 401) {
      const userStore = useUserStore()
      userStore.doLogout()
      router.replace('/login')
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
