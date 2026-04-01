import { request } from '@/utils'

export default {
    getList: () => request.get('/plugins/announcement'),
    create: (data) => request.post('/plugins/announcement', data),
    delete: (id) => request.delete(`/plugins/announcement/${id}`),
}
