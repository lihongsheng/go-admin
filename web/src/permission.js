import router from '@/router'
import pinia from '@/store'
import { useUserStore } from '@/store/modules/user'
import { usePermissionStore } from '@/store/modules/permission'
import { getInstallStatus } from '@/api/install'

let installChecked = false
let installed = false

async function ensureInstallStatus(force = false) {
  if (installChecked && !force) return installed
  try {
    const { data } = await getInstallStatus()
    installed = !!data.installed
  } catch (_) {
    installed = false
  }
  installChecked = true
  return installed
}

// 由安装向导在 install_done 后调用，强制刷新缓存
export function refreshInstallStatus() {
  installChecked = false
  installed = false
  return ensureInstallStatus(true)
}

const whiteList = ['/login', '/install', '/404']

router.beforeEach(async (to, _from, next) => {
  const ok = await ensureInstallStatus()
  if (!ok) {
    if (to.path !== '/install') return next('/install')
    return next()
  }
  if (to.path === '/install') return next('/')

  // Pinia 必须在 app.use(pinia) 之后才能使用，在路由守卫中通过 pinia 实例获取 store
  const userStore = useUserStore(pinia)
  const permStore = usePermissionStore(pinia)

  if (userStore.token) {
    if (!userStore.userInfo) {
      try {
        await userStore.fetchInfo()
        const dynamic = await permStore.generateRoutes()
        // 将动态路由逐条挂到 Layout 下，确保路径与菜单 SubTree 一致
        for (const r of dynamic) {
          router.addRoute('Layout', r)
        }
        return next({ ...to, replace: true })
      } catch (e) {
        await userStore.doLogout()
        return next('/login')
      }
    }
    return next()
  }
  if (whiteList.includes(to.path)) return next()
  next('/login')
})
