<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:50:35
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <n-menu
    ref="menu"
    class="side-menu"
    accordion
    :indent="18"
    :collapsed-icon-size="22"
    :collapsed-width="64"
    :expanded-keys="expandedKeys"
    :collapsed="appStore.collapsed"
    :options="permissionStore.menus"
    :value="activeKey"
    @update:value="handleMenuSelect"
    @update:expanded-keys="handleExpandedKeysUpdate"
  />
</template>

<script setup>
import { useAppStore, usePermissionStore } from '@/store'
import { isExternal } from '@/utils'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const permissionStore = usePermissionStore()

const activeKey = computed(() => route.meta?.parentKey || route.name)
const ancestorsOfActiveKey = computed(() => {
  if (!activeKey.value || !permissionStore.menus.length) return []
  return getAncestors(activeKey.value) || []
})

const menu = ref(null)
const expandedKeys = ref([])
const manuallyCollapsed = ref(new Set())

watch(
  [activeKey, () => permissionStore.menus.length],
  async () => {
    await nextTick()
    if (!activeKey.value || !permissionStore.menus.length)
      return
    // Route changed — clear manual collapse state and expand ancestors for new route
    manuallyCollapsed.value = new Set()
    const ancestors = getAncestors(activeKey.value)
    expandedKeys.value = ancestors && ancestors.length
      ? ancestors.filter(k => !manuallyCollapsed.value.has(k))
      : []
    menu.value?.showOption()
  },
  { immediate: true },
)

function handleExpandedKeysUpdate(keys) {
  const prev = new Set(expandedKeys.value)
  const next = new Set(keys)

  // Detect newly expanded keys — if any, this is an expand action;
  // the accordion auto-collapsed groups should NOT be treated as manual collapse
  let hasExpand = false
  for (const k of next) {
    if (!prev.has(k)) {
      hasExpand = true
      manuallyCollapsed.value.delete(k)
    }
  }

  // Only track manual collapse when the user collapsed a group
  // without expanding another (pure collapse action)
  if (!hasExpand) {
    for (const k of prev) {
      if (!next.has(k) && ancestorsOfActiveKey.value.includes(k))
        manuallyCollapsed.value.add(k)
    }
  }

  expandedKeys.value = keys
}

/** 查找 targetKey 的所有祖先节点的 key 列表 */
function getAncestors(targetKey, menus = permissionStore.menus, parents = []) {
  for (const menu of menus) {
    if (menu.key === targetKey) return parents
    if (menu.children) {
      const result = getAncestors(targetKey, menu.children, [...parents, menu.key])
      if (result) return result
    }
  }
  return null
}

function handleMenuSelect(key, item) {
  // 点击无子菜单的节点时：只保留其祖先展开，其余菜单折叠
  if (!item.children || item.children.length === 0) {
    const ancestors = getAncestors(key)
    if (expandedKeys.value.length > 0) {
      expandedKeys.value = expandedKeys.value.filter(k => ancestors?.includes(k))
    }
  }

  if (isExternal(item.originPath)) {
    $dialog.confirm({
      type: 'info',
      title: `请选择打开方式`,
      positiveText: '外链打开',
      negativeText: '在本站内嵌打开',
      confirm() {
        window.open(item.originPath)
      },
      cancel: () => {
        router.push(item.path)
      },
    })
  }
  else {
    if (!item.path)
      return
    router.push(item.path)
  }
}
</script>

<style>
.side-menu:not(.n-menu--collapsed) {
  .n-menu-item-content {
    &::before {
      left: 8px;
      right: 8px;
    }
    &.n-menu-item-content--selected::before {
      border-left: 4px solid rgb(var(--primary-color));
    }
  }
}
</style>
