import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { getMockEnabled, mockResolve } from '@/api/mock'
import router from '@/router'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 30000
})

request.interceptors.request.use(async cfg => {
  const userStore = useUserStore()
  if (userStore.token) {
    cfg.headers.Authorization = 'Bearer ' + userStore.token
  }

  // 整体 Mock 模式：命中注册表 → 挂自定义 adapter 直接返回标准响应 { code, msg, data }，
  // 复用下方响应拦截器的统一错误处理；未命中放行真实请求
  // 开关动态读取（运行时可通过 setMockEnabled 切换）
  if (getMockEnabled()) {
    const hit = mockResolve(cfg.method, cfg.url)
    if (hit) {
      const body = await hit.handler({ query: cfg.params, body: cfg.data, params: hit.params })
      if (body !== undefined) {
        // 自定义 adapter 跳过网络层（不可直接返回伪响应对象：
        // axios 会把拦截器返回值当作 config 继续走 dispatchRequest）
        cfg.adapter = async () => ({
          data: body,
          status: 200,
          statusText: 'OK',
          headers: {},
          config: cfg,
          request: {}
        })
      }
    } else {
      console.warn('[mock] 未注册的接口，放行真实请求:', cfg.method, cfg.url)
    }
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
    const resp = err.response
    if (resp && resp.data && resp.data.msg) {
      ElMessage.error(resp.data.msg)
    } else {
      ElMessage.error(err.message || 'network error')
    }
    return Promise.reject(err)
  }
)

export default request
