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

// 从用户信息中提取默认首页路由
// 遍历用户角色的 default_router，取首个非空值；全部为空则返回 /dashboard
function getDefaultRouter(info) {
  if (!info || !info.roles || !info.roles.length) return '/dashboard'
  for (const role of info.roles) {
    if (role.default_router) return role.default_router
  }
  return '/dashboard'
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
        const info = await userStore.fetchInfo()
        const dynamic = await permStore.generateRoutes()
        // 将动态路由逐条挂到 Layout 下，确保路径与菜单 SubTree 一致
        for (const r of dynamic) {
          router.addRoute('Layout', r)
        }
        // 登录后跳转到角色设定的默认首页，为空则默认 /dashboard
        if (to.path === '/') {
          const defaultRouter = getDefaultRouter(info)
          return next({ path: defaultRouter, replace: true })
        }
        return next({ ...to, replace: true })
      } catch (e) {
        await userStore.doLogout()
        return next('/login')
      }
    }
    // 已登录用户访问根路径时，重定向到默认首页
    if (to.path === '/') {
      const defaultRouter = getDefaultRouter(userStore.userInfo)
      return next({ path: defaultRouter, replace: true })
    }
    return next()
  }
  if (whiteList.includes(to.path)) return next()
  next('/login')
})
