// 共享的媒体相关工具
export function humanSize(bytes) {
  if (bytes == null)
    return '-'
  if (bytes < 1024)
    return `${bytes} B`
  if (bytes < 1024 * 1024)
    return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024)
    return `${(bytes / 1024 / 1024).toFixed(1)} MB`
  return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`
}

export function iconForMime(mime = '') {
  if (mime.startsWith('image/'))
    return 'i-fe:image'
  if (mime.startsWith('video/'))
    return 'i-fe:film'
  if (mime.startsWith('audio/'))
    return 'i-fe:music'
  if (mime.includes('pdf'))
    return 'i-fe:file-text'
  if (mime.includes('zip') || mime.includes('compressed'))
    return 'i-fe:archive'
  return 'i-fe:file'
}

// 把后端返回的扁平 folder 列表构建成 NTree 节点（含根节点）
export function buildFolderTree(folders, rootLabel = '全部文件') {
  const byId = new Map()
  for (const f of folders)
    byId.set(f.id, { ...f, children: [] })

  const root = { id: 0, key: 0, label: rootLabel, path: '', children: [] }
  for (const node of byId.values()) {
    const treeNode = {
      id: node.id,
      key: node.id,
      label: node.name,
      path: node.path,
      parentId: node.parentId,
      configId: node.configId,
      children: node.children,
    }
    if (node.parentId && byId.has(node.parentId)) {
      const parent = byId.get(node.parentId)
      parent.children.push(treeNode)
    }
    else {
      root.children.push(treeNode)
    }
  }
  return [root]
}
