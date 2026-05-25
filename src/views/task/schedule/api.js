import { request } from '@/utils'

export default {
  create: data => request.post('/task', data),
  read: (params = {}) => request.get('/task/page', { params }),
  update: data => request.patch(`/task/${data.id}`, data),
  delete: id => request.delete(`/task/${id}`),
  toggle: id => request.patch(`/task/${id}/toggle`),
  run: id => request.post(`/task/${id}/run`),

  // 批量
  bulkDelete: ids => request.post('/task/bulk/delete', { ids }),
  bulkToggle: (ids, enabled) => request.post('/task/bulk/toggle', { ids, enabled }),

  // 仪表盘与可视化
  stats: () => request.get('/task/stats'),
  previewNext: (spec, n = 5) => request.get('/task/preview-next', { params: { spec, n } }),

  // 日志与内置函数
  getLogs: (params = {}) => request.get('/task/log/page', { params }),
  getFuncs: () => request.get('/task/funcs'),
}
