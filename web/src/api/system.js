import request from '@/utils/request'

// user
export const userList   = (params) => request.get('/api/v1/system/user/list', { params })
export const userCreate = (data)   => request.post('/api/v1/system/user', data)
export const userUpdate = (data)   => request.put('/api/v1/system/user', data)
export const userDelete = (id)     => request.delete('/api/v1/system/user/' + id)

// role
export const roleList       = ()        => request.get('/api/v1/system/role/list')
export const roleCreate     = (data)    => request.post('/api/v1/system/role', data)
export const roleUpdate     = (data)    => request.put('/api/v1/system/role', data)
export const roleDelete     = (id)      => request.delete('/api/v1/system/role/' + id)
export const roleAuth            = (data)    => request.post('/api/v1/system/role/auth', data)
export const roleAuthDetail      = (id)      => request.get('/api/v1/system/role/auth/' + id)
export const roleSetDefaultRouter = (id, data) => request.put('/api/v1/system/role/' + id + '/default-router', data)

// menu
export const menuTree   = ()        => request.get('/api/v1/system/menu/tree')
export const menuCreate = (data)    => request.post('/api/v1/system/menu', data)
export const menuUpdate = (data)    => request.put('/api/v1/system/menu', data)
export const menuDelete = (id)      => request.delete('/api/v1/system/menu/' + id)

// api
export const apiList    = (params)  => request.get('/api/v1/system/api/list', { params })
export const apiCreate  = (data)    => request.post('/api/v1/system/api', data)
export const apiUpdate  = (data)    => request.put('/api/v1/system/api', data)
export const apiDelete  = (id)      => request.delete('/api/v1/system/api/' + id)

// plugin
export const pluginList = ()        => request.get('/api/v1/plugin/list')
