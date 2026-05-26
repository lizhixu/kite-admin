import { request } from '@/utils'

export default {
  getLogs: params => request.get('/loginlog/list', { params }),
}
