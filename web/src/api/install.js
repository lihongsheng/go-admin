import request from '@/utils/request'

export const getInstallStatus = () => request.get('/install/status')
export const checkDB = (db)     => request.post('/install/check-db', db)
export const doInstall = (body) => request.post('/install/init', body)
