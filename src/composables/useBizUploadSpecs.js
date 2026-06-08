import { request } from '@/utils'

// 全局缓存：业务上传白名单（一个会话拉一次）
let cache = null
let inflight = null

export async function getBizUploadSpecs() {
  if (cache)
    return cache
  if (inflight)
    return inflight
  inflight = request
    .get('/upload/specs')
    .then((res) => {
      cache = res?.data || {}
      inflight = null
      return cache
    })
    .catch((err) => {
      inflight = null
      throw err
    })
  return inflight
}

export function getCachedBizUploadSpec(bizType) {
  return cache?.[bizType] || null
}

export function resetBizUploadSpecsCache() {
  cache = null
  inflight = null
}
