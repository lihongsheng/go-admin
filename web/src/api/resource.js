import request from '@/utils/request'

// 资源管理插件 API（统一 mock 方案：mock 由 api/mock/resource.js 注册表覆盖）
export const resourceList = (params) => request.get('/api/plugin/resource/v1/list', { params })

export const resourceDetail = (id) => request.get('/api/plugin/resource/v1/' + id)

export const resourceCreate = (data) => request.post('/api/plugin/resource/v1', data)

export const resourceUpdate = (id, data) => request.put('/api/plugin/resource/v1/' + id, data)

export const resourceBatchDelete = (ids) => request.delete('/api/plugin/resource/v1/batch', { data: { ids } })
