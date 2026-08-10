/**
 * 资源管理 适配层（中间件）
 *
 * 定位：前端页面与后端接口之间的统一抽象层。
 * 页面只允许导入本模块暴露的 4 个标准方法，禁止直接拼接后端 URL。
 *
 * 职责：
 *  1. 统一方法契约：listResources / getResourceDetail / createResource / updateResource
 *  2. Token 注入 + 错误码转友好提示：由 utils/request 拦截器统一完成（401 跳登录、非 0 码弹 msg）
 *  3. 双环境切换：Mock 模拟数据 / 真实后端，开关仅在本文件维护
 *
 * 环境切换：
 *  - 方式一（环境变量）：.env 中配置 VITE_RESOURCE_MOCK=true 开启 Mock
 *  - 方式二（直接改码）：把下方 USE_MOCK 改为 true
 *  切换后页面代码无需任何改动。
 */
import request from '@/utils/request'

// ============ 环境开关（仅在此层维护） ============
const USE_MOCK = import.meta.env.VITE_RESOURCE_MOCK === 'true'

// ============ Mock 数据实现 ============
let mockSeq = 5
const mockDb = [
  { id: 1, name: 'web-server-01', type: '通用型 g6', status: 1, status_at: '2026-08-01 10:00:00', created_at: '2026-08-01 09:00:00', updated_at: '2026-08-01 10:00:00' },
  { id: 2, name: 'db-master-01', type: '内存型 r6', status: 1, status_at: '2026-08-02 11:30:00', created_at: '2026-08-02 10:00:00', updated_at: '2026-08-02 11:30:00' },
  { id: 3, name: 'cache-node-02', type: '计算型 c6', status: 2, status_at: '2026-08-03 14:20:00', created_at: '2026-08-03 09:00:00', updated_at: '2026-08-03 14:20:00' },
  { id: 4, name: 'build-runner-03', type: '通用型 g6', status: 3, status_at: '2026-08-04 16:00:00', created_at: '2026-08-04 08:00:00', updated_at: '2026-08-04 16:00:00' },
  { id: 5, name: 'cdn-edge-07', type: '突发性能型 t6', status: 1, status_at: '2026-08-05 09:15:00', created_at: '2026-08-05 09:00:00', updated_at: '2026-08-05 09:15:00' },
]

const delay = (ms = 300) => new Promise(r => setTimeout(r, ms))
const now = () => new Date().toISOString().replace('T', ' ').slice(0, 19)

const mockImpl = {
  // 关键词 / 类型 / 状态 多条件分页查询
  async listResources(params = {}) {
    await delay()
    const { keyword = '', type = '', status = 0, page = 1, limit = 10 } = params
    const kw = String(keyword).trim().toLowerCase()
    const filtered = mockDb.filter(r =>
      (!kw || r.name.toLowerCase().includes(kw)) &&
      (!type || r.type === type) &&
      (!status || r.status === Number(status))
    )
    return {
      code: 0, msg: 'ok',
      data: { list: filtered.slice((page - 1) * limit, page * limit), total: filtered.length }
    }
  },

  async getResourceDetail(id) {
    await delay()
    const row = mockDb.find(r => r.id === Number(id))
    if (!row) return { code: 1, msg: '资源不存在' }
    return { code: 0, msg: 'ok', data: { ...row } }
  },

  async createResource(data) {
    await delay()
    const row = {
      id: ++mockSeq, ...data,
      status_at: now(), created_at: now(), updated_at: now()
    }
    mockDb.unshift(row)
    return { code: 0, msg: 'ok', data: { ...row } }
  },

  async updateResource(id, data) {
    await delay()
    const row = mockDb.find(r => r.id === Number(id))
    if (!row) return { code: 1, msg: '资源不存在' }
    Object.assign(row, data, { updated_at: now() })
    // 与后端一致：仅状态变更时刷新状态更新时间
    if (data.status && row.status !== data.status) row.status_at = now()
    return { code: 0, msg: 'ok', data: { ...row } }
  },

  // 扩展方法（页面批量删除用，不在题目 4 方法契约内）
  async batchDeleteResources(ids) {
    await delay()
    ids.forEach(id => {
      const i = mockDb.findIndex(r => r.id === Number(id))
      if (i > -1) mockDb.splice(i, 1)
    })
    return { code: 0, msg: 'ok' }
  }
}

// ============ 真实后端实现 ============
// Token 注入 / 错误码转友好提示由 utils/request 拦截器统一完成
const realImpl = {
  listResources(params) {
    return request.get('/api/plugin/resource/v1/list', { params })
  },
  getResourceDetail(id) {
    return request.get('/api/plugin/resource/v1/' + id)
  },
  createResource(data) {
    return request.post('/api/plugin/resource/v1', data)
  },
  updateResource(id, data) {
    return request.put('/api/plugin/resource/v1/' + id, data)
  },
  batchDeleteResources(ids) {
    return request.delete('/api/plugin/resource/v1/batch', { data: { ids } })
  }
}

// 对外唯一出口：页面只依赖这一个对象
export const resourceAdapter = USE_MOCK ? mockImpl : realImpl
