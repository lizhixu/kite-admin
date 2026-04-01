import { request } from '@/utils'

export default {
    getList: () => request.get('/plugins/tasks'),
    create: (data) => request.post('/plugins/tasks', data),
    updateStatus: (id, status) => request.patch(`/plugins/tasks/${id}/status`, { status }),
    delete: (id) => request.delete(`/plugins/tasks/${id}`),
}
