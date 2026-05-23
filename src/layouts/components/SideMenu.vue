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

const menu = ref(null)
const expandedKeys = ref([])

watch(
  [activeKey, () => permissionStore.menus.length],
  async () => {
    await nextTick()
    if (!activeKey.value || !permissionStore.menus.length)
      return
    const ancestors = getAncestors(activeKey.value)
    if (ancestors && ancestors.length)
      expandedKeys.value = ancestors
    menu.value?.showOption()
  },
  { immediate: true },
)

function handleExpandedKeysUpdate(keys) {
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
