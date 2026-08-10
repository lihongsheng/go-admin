/**
 * Mock 注册表核心（独立模块，无依赖，避免与各 mock 模块形成循环导入）
 * 由 mock/index.js 聚合导出，各 mock 模块从本文件 import mockRoute 注册路由
 */

// 模拟网络延迟（适配层统一处理，页面感知不到）
export const mockDelay = () => new Promise(r => setTimeout(r, 150 + Math.random() * 250))

// 当前时间字符串，mock 数据用
export const mockNow = () => new Date().toISOString().replace('T', ' ').slice(0, 19)

const routes = []

function patternToRegExp(pattern) {
  const segs = pattern.split('/').map(s =>
    s.startsWith(':') ? '([^/]+)' : s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  )
  return new RegExp('^' + segs.join('/') + '$')
}

// 注册 mock 路由：handler({ query, body, params }) → { code, msg, data }
// 返回 undefined 表示不处理（放行真实请求）
export function mockRoute(method, pattern, handler) {
  routes.push({ method: method.toUpperCase(), pattern, re: patternToRegExp(pattern), handler })
}

// 供 request 拦截器调用：命中返回 { handler, params }，未命中返回 null
export function mockResolve(method, url) {
  const path = String(url || '').split('?')[0]
  const m = routes.find(r => r.method === method.toUpperCase() && r.re.test(path))
  if (!m) return null
  const names = m.pattern.split('/').filter(s => s.startsWith(':')).map(s => s.slice(1))
  const values = m.re.exec(path).slice(1)
  const params = {}
  names.forEach((n, i) => { params[n] = values[i] })
  return { handler: m.handler, params }
}
