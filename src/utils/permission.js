import { usePermissionStore } from '@/store'

function walkPermissions(permissions = [], matcher) {
  for (const permission of permissions) {
    if (matcher(permission))
      return true
    if (permission.children?.length && walkPermissions(permission.children, matcher))
      return true
  }
  return false
}

export function hasPermissionCode(code) {
  if (!code)
    return false
  const permissionStore = usePermissionStore()
  return walkPermissions(permissionStore.permissions, permission =>
    permission.code === code && permission.enable !== false,
  )
}

export function hasAnyPermissionCode(codes = []) {
  return codes.some(code => hasPermissionCode(code))
}
