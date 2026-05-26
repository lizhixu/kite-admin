import { request } from '@/utils'

export default {
  // 媒体
  page: (params = {}) => request.get('/media/page', { params }),
  delete: id => request.delete(`/media/${id}`),
  bulkDelete: ids => request.post('/media/bulk/delete', { ids }),
  move: (ids, folderId) => request.post('/media/move', { ids, folderId }),
  upload: (file, { configId, folderId } = {}, onProgress) => {
    const fd = new FormData()
    fd.append('file', file)
    if (configId != null)
      fd.append('configId', configId)
    if (folderId)
      fd.append('folderId', folderId)
    return request.post('/media/upload', fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (e) => {
        if (e.total)
          onProgress?.(Math.round((e.loaded / e.total) * 100))
      },
    })
  },

  // 文件夹
  listFolders: configId => request.get('/media/folder/tree', { params: { configId } }),
  resolveFolder: (configId, path, autoCreate = false) =>
    request.get('/media/folder/resolve', { params: { configId, path, autoCreate: autoCreate ? '1' : '0' } }),
  createFolder: data => request.post('/media/folder', data),
  renameFolder: (id, name) => request.patch(`/media/folder/${id}`, { name }),
  deleteFolder: (id, cascade = false) =>
    request.delete(`/media/folder/${id}`, { params: cascade ? { cascade: 1 } : {} }),

  // 存储配置（媒体库页面只需要列表用于选择）
  listConfigs: () => request.get('/storage/config'),
}
