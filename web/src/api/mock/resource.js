/**
 * 资源管理 mock（统一 mock 方案）
 * 注册 /api/plugin/resource/v1/* 全部路由，覆盖 api/resource.js 封装的
 * resourceList / resourceDetail / resourceCreate / resourceUpdate
 * 及 resourceBatchDelete。结构同真实后端。
 */
import { mockRoute, mockDelay, mockNow } from './registry'

let nextId = 10
const mockResources = [
  { id: 1, name: 'web-server-01', type: '通用型 g6', status: 1, status_at: '2026-08-01 10:00:00', created_at: '2026-08-01 09:00:00', updated_at: '2026-08-01 10:00:00' },
  { id: 2, name: 'db-master-01', type: '内存型 r6', status: 1, status_at: '2026-08-02 11:30:00', created_at: '2026-08-02 10:00:00', updated_at: '2026-08-02 11:30:00' },
  { id: 3, name: 'cache-node-02', type: '计算型 c6', status: 2, status_at: '2026-08-03 14:20:00', created_at: '2026-08-03 09:00:00', updated_at: '2026-08-03 14:20:00' },
  { id: 4, name: 'build-runner-03', type: '通用型 g6', status: 3, status_at: '2026-08-04 16:00:00', created_at: '2026-08-04 08:00:00', updated_at: '2026-08-04 16:00:00' },
  { id: 5, name: 'cdn-edge-07', type: '突发性能型 t6', status: 1, status_at: '2026-08-05 09:15:00', created_at: '2026-08-05 09:00:00', updated_at: '2026-08-05 09:15:00' },
]

mockRoute('GET', '/api/plugin/resource/v1/list', async ({ query }) => {
  await mockDelay()
  const { keyword = '', type = '', status = 0, page = 1, limit = 10 } = query || {}
  const kw = String(keyword).trim().toLowerCase()
  const filtered = mockResources.filter(r =>
    (!kw || r.name.toLowerCase().includes(kw)) &&
    (!type || r.type === type) &&
    (!Number(status) || r.status === Number(status))
  )
  const start = (Number(page) - 1) * Number(limit)
  return { code: 0, msg: 'ok', data: { list: filtered.slice(start, start + Number(limit)), total: filtered.length } }
})

mockRoute('GET', '/api/plugin/resource/v1/:id', async ({ params }) => {
  await mockDelay()
  const row = mockResources.find(r => r.id === Number(params.id))
  if (!row) return { code: 1, msg: '资源不存在' }
  return { code: 0, msg: 'ok', data: { ...row } }
})

mockRoute('POST', '/api/plugin/resource/v1', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const row = {
    id: nextId++, name: body.name, type: body.type, status: body.status ?? 1,
    status_at: mockNow(), created_at: mockNow(), updated_at: mockNow()
  }
  mockResources.unshift(row)
  return { code: 0, msg: 'ok', data: { ...row } }
})

mockRoute('PUT', '/api/plugin/resource/v1/:id', async ({ params, body }) => {
  await mockDelay()
  body = body || {}
  const row = mockResources.find(r => r.id === Number(params.id))
  if (!row) return { code: 1, msg: '资源不存在' }
  Object.assign(row, { name: body.name ?? row.name, type: body.type ?? row.type, updated_at: mockNow() })
  // 与真实后端一致：仅状态变更时刷新状态更新时间
  if (body.status && row.status !== Number(body.status)) {
    row.status = Number(body.status)
    row.status_at = mockNow()
  }
  return { code: 0, msg: 'ok', data: { ...row } }
})

mockRoute('DELETE', '/api/plugin/resource/v1/batch', async ({ body }) => {
  await mockDelay()
  const ids = (body && body.ids) || []
  if (!ids.length) return { code: 1, msg: '请选择要删除的资源' }
  ids.forEach(id => {
    const i = mockResources.findIndex(r => r.id === Number(id))
    if (i > -1) mockResources.splice(i, 1)
  })
  return { code: 0, msg: 'ok' }
})
