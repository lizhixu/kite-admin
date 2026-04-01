import { request } from '@/utils'

export default {
    getList: () => request.get('/plugin/list'),
    install: (code) => request.post(`/plugin/install/${code}`),
    uninstall: (code) => request.post(`/plugin/uninstall/${code}`),
    toggleEnable: (code, enable) => request.patch(`/plugin/enable/${code}`, { enable }),
}
