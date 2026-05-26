import { request } from '@/utils'

export default {
  getSystemConfig: () => request.get('/system/config'),
  saveSystemConfig: data => request.put('/system/config', data),
}
