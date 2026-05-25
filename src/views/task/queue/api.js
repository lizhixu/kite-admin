import { request } from '@/utils'

export default {
  // 队列查询 / 编辑（运维参数）
  getOne: id => request.get(`/queue/${id}`),
  read: (params = {}) => request.get('/queue/page', { params }),
  update: data => request.patch(`/queue/${data.id}`, data),
  delete: id => request.delete(`/queue/${id}`),
  toggle: id => request.patch(`/queue/${id}/toggle`),

  // 汇总 / 已注册 handlers
  stats: () => request.get('/queue/stats'),
  getHandlers: () => request.get('/queue/handlers'),

  // 任务（jobs）—— params 支持 status / from / to / pageNo / pageSize
  getJobs: (queueId, params = {}) => request.get(`/queue/${queueId}/jobs`, { params }),
  addJob: (queueId, data) => request.post(`/queue/${queueId}/job`, data),
  bulkAddJobs: (queueId, items) => request.post(`/queue/${queueId}/jobs/bulk`, { items }),
  // params 支持 status / before（ISO 时间，清理 N 天前完成记录）
  clearJobs: (queueId, params = {}) => request.delete(`/queue/${queueId}/jobs`, { params }),
  deleteJob: jobId => request.delete(`/queue/job/${jobId}`),

  // Kick：单条 / 整队列
  kickJob: jobId => request.post(`/queue/job/${jobId}/kick`),
  kickAll: queueId => request.post(`/queue/${queueId}/kick`),
}
