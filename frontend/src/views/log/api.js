import { request } from '@/utils'

export default {
  getLogs: (params) => request.get('/syslog/list', { params }),
}
