import { request } from '@/utils'

export default {
  list: () => request.get('/storage/config'),
  create: data => request.post('/storage/config', data),
  update: data => request.patch(`/storage/config/${data.id}`, data),
  delete: id => request.delete(`/storage/config/${id}`),
  setDefault: id => request.patch(`/storage/config/${id}/default`),
  test: id => request.post(`/storage/config/${id}/test`),
}
