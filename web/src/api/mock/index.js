/**
 * 前端整体 Mock 层（拦截器中间件，适配层体系的一部分）
 *
 * 职责：
 *  - 动态开关：getMockEnabled() / setMockEnabled()，localStorage 持久化，
 *    未设置时取环境变量 VITE_USE_MOCK 作为默认值（编译期默认 / 运行时可控）
 *  - 路由注册表核心见 mock/registry.js（独立模块避免循环导入）
 *  - 由 utils/request 请求拦截器调用：命中注册表直接返回标准响应 { code, msg, data }，未命中放行真实请求
 *  - 统一 mock 方案：包括资源管理（/api/plugin/resource/v1/*，见 mock/resource.js）在内的
 *    全部接口都在本注册表维护；页面经 api 封装的请求在 mock 开启时由本层自动覆盖
 *
 * 启用方式：
 *  运行时：布局右上角 Mock 开关 / setMockEnabled(true|false)（立即生效，localStorage 持久化）
 *  默认值：.env 配置 VITE_USE_MOCK=true
 *
 * 页面代码零改动；所有 mock 数据只在 api/mock/ 目录内维护。
 */

// 环境变量仅作为默认值；运行时可被 setMockEnabled 覆盖并持久化
const DEFAULT_MOCK = import.meta.env.VITE_USE_MOCK === 'true'
const STORAGE_KEY = 'go-admin-mock-enabled'

export function getMockEnabled() {
  const v = localStorage.getItem(STORAGE_KEY)
  if (v === null) return DEFAULT_MOCK
  return v === '1'
}

export function setMockEnabled(enabled) {
  localStorage.setItem(STORAGE_KEY, enabled ? '1' : '0')
}

export { mockRoute, mockResolve, mockDelay, mockNow } from './registry'

// ---------- 各模块 mock 注册（副作用导入） ----------
// 统一 mock 方案：所有接口（含资源管理 /api/plugin/resource/v1/*）都在本注册表注册
import './auth'
import './user'
import './role'
import './menu'
import './mch'
import './plugin'
import './resource'
import './example'
