import { request } from '@/utils'

export default {
  read: (params = {}) => request.get('/message/page', { params }),
  create: data => request.post('/message', data),
  delete: id => request.delete(`/message/${id}`),
  bulkDelete: ids => request.post('/message/bulk/delete', { ids }),

  getMyMessages: (params = {}) => request.get('/message/mine', { params }),
  getUnreadCount: () => request.get('/message/unread/count'),
  markRead: id => request.patch(`/message/${id}/read`),
  markAllRead: () => request.patch('/message/read/all'),

  getEmailConfig: () => request.get('/email/config'),
  saveEmailConfig: data => request.put('/email/config', data),
  testEmailConfig: () => request.post('/email/config/test'),

  getEmailTemplates: () => request.get('/email-template/list'),
  getEmailTemplate: id => request.get(`/email-template/${id}`),
  saveEmailTemplate: (id, data) => request.put(`/email-template/${id}`, data),
  previewEmailTemplate: (id, vars) => request.post(`/email-template/${id}/preview`, { vars }),
}
