<template>
  <n-popover
    trigger="click"
    placement="bottom"
    :style="{ width: '360px' }"
    @update:show="onShow"
  >
    <template #trigger>
      <span class="mr-16">
        <n-badge :value="store.unreadCount" :max="99" :offset="[-6, -2]" class="notification-badge">
          <i
            class="cursor-pointer i-fe:bell text-18"
            style="vertical-align: middle;"
          />
        </n-badge>
      </span>
    </template>

    <div class="notification-popover">
      <div class="flex items-center justify-between px-12 py-8">
        <span class="font-bold text-14">消息通知</span>
        <n-space size="small">
          <n-button
            text
            type="primary"
            size="tiny"
            :disabled="store.unreadCount === 0"
            @click="handleMarkAllRead"
          >
            全部已读
          </n-button>
        </n-space>
      </div>

      <n-divider style="margin: 0" />

      <div v-if="unreadMessages.length === 0" class="py-24 text-center" style="color: #999">
        暂无未读消息
      </div>

      <n-list v-else hoverable clickable :show-divider="false" class="notification-list">
        <n-list-item
          v-for="msg in unreadMessages"
          :key="msg.id"
          @click="handleClickMessage(msg)"
        >
          <n-thing>
            <template #header>
              <div class="flex items-center gap-8">
                <span :class="{ 'font-bold': !msg.isRead }" class="text-13">
                  {{ msg.title }}
                </span>
                <n-badge v-if="!msg.isRead" dot type="error" processing />
              </div>
            </template>
            <template #header-extra>
              <span class="text-12" style="color: #999">
                {{ formatTime(msg.createTime) }}
              </span>
            </template>
            <template #description>
              <span class="text-12" style="color: #666">
                {{ truncate(msg.content, 50) }}
              </span>
            </template>
          </n-thing>
        </n-list-item>
      </n-list>

      <n-divider style="margin: 0" />

      <div class="text-center py-8">
        <n-button text type="primary" size="tiny" @click="goToMessageList">
          查看全部消息
        </n-button>
      </div>
    </div>
  </n-popover>
</template>

<script setup>
import { NBadge, NButton, NDivider, NList, NListItem, NPopover, NSpace, NThing } from 'naive-ui'
import { useNotificationStore } from '@/store/modules/notification'
import { stripMarkdown } from '@/utils/markdown'

const store = useNotificationStore()
const router = useRouter()
const unreadMessages = computed(() => store.messages.filter(m => !m.isRead))

function onShow(visible) {
  if (visible) {
    store.fetchRecentMessages()
  }
}

function handleMarkAllRead() {
  store.markAllAsRead()
}

function handleClickMessage(msg) {
  if (!msg.isRead) {
    store.markAsRead(msg.id)
  }
}

function goToMessageList() {
  router.push('/message/list')
}

function truncate(str, len) {
  if (!str) return ''
  const plain = stripMarkdown(str)
  return plain.length > len ? plain.slice(0, len) + '...' : plain
}

function formatTime(time) {
  if (!time) return ''
  const d = new Date(time)
  const now = new Date()
  const diff = now - d
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  return `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

onMounted(() => {
  store.fetchUnreadCount()
  store.connectSSE()
})

onBeforeUnmount(() => {
  store.disconnectSSE()
})
</script>

<style scoped>
.notification-popover {
  margin: -12px -16px;
}
.notification-list {
  max-height: 360px;
  overflow-y: auto;
}
.notification-list :deep(.n-list-item) {
  padding: 8px 12px;
}
.notification-badge :deep(.n-badge-sup) {
  font-size: 10px;
  padding: 0 4px;
  height: 14px;
  min-width: 14px;
  line-height: 14px;
}
</style>