export function mediaAccessUrl(item) {
  return item?.accessUrl || item?.url || ''
}

export function normalizeMediaItem(item) {
  if (!item)
    return item
  const accessUrl = mediaAccessUrl(item)
  return {
    ...item,
    accessUrl,
    url: accessUrl,
  }
}

export function normalizeMediaItems(items) {
  return (items || []).map(normalizeMediaItem)
}
