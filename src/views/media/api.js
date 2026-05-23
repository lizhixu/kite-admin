import { request } from '@/utils'

export default {
  // 媒体
  read: (params = {}) => request.get('/media/page', { params }),
  delete: id => request.delete(`/media/${id}`),
  bulkDelete: ids => request.post('/media/bulk/delete', { ids }),
  upload: (file, configId, onProgress) => {
    const fd = new FormData()
    fd.append('file', file)
    if (configId)
      fd.append('configId', configId)
    return request.post('/media/upload', fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (e) => {
        if (e.total)
          onProgress?.(Math.round((e.loaded / e.total) * 100))
      },
    })
  },

  // 存储配置
  listConfigs: () => request.get('/storage/config'),
  createConfig: data => request.post('/storage/config', data),
  updateConfig: data => request.patch(`/storage/config/${data.id}`, data),
  deleteConfig: id => request.delete(`/storage/config/${id}`),
  setDefault: id => request.patch(`/storage/config/${id}/default`),
  testConfig: id => request.post(`/storage/config/${id}/test`),
}
