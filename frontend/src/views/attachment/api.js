import { request } from '@/utils'

export default {
    getList: (params) => request.get('/attachment', { params }),
    upload: (formData, onProgress) =>
        request.post('/attachment/upload', formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
            onUploadProgress: onProgress,
        }),
    delete: (id) => request.delete(`/attachment/${id}`),
    getConfig: () => request.get('/attachment/config'),
    saveConfig: (data) => request.post('/attachment/config', data),
    getGroups: () => request.get('/attachment/group'),
    createGroup: (data) => request.post('/attachment/group', data),
    updateGroup: (id, data) => request.patch(`/attachment/group/${id}`, data),
    deleteGroup: (id) => request.delete(`/attachment/group/${id}`),
}
