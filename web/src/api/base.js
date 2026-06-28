import request from '@/utils/request'

export const captcha = ()    => request.get('/api/v1/base/captcha')
export const login  = (data) => request.post('/api/v1/base/login', data)
export const logout = ()     => request.post('/api/v1/base/logout')
export const info   = ()     => request.get('/api/v1/base/info')
export const menu         = ()     => request.get('/api/v1/base/menu')
