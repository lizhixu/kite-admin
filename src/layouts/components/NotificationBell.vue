<template>
  <NBadge :value="store.unreadCount" :max="99" :offset="[4, 0]" class="notification-badge mr-16" @click="showInbox = true">
    <i class="i-fe:bell cursor-pointer text-18" />
  </NBadge>

  <NDrawer v-model:show="showInbox" width="560" placement="right" @after-leave="resetView">
    <NDrawerContent>
      <template #header>
        <div
          class="flex items-center"
          style="width: 100%; gap: 8px"
        >
          <NButton v-if="view === 'detail'" text @click="view = 'list'">
            <template #icon><i class="i-fe:arrow-left" /></template>
          </NButton>
          <span class="flex-1 text-16 font-bold">{{ view === 'detail' ? detailMsg?.title : '我的消息' }}</span>
          <NButton
            v-if="view === 'list'"
            text
            type="primary"
            size="tiny"
            :disabled="store.unreadCount === 0"
            @click="handleMarkAllRead"
          >
            全部已读
          </NButton>
        </div>
      </template>

      <!-- 列表视图 -->
      <template v-if="view === 'list'">
        <NSpin :show="inboxLoading">
          <div
            v-if="store.inboxMessages.length === 0"
            class="py-48 text-center"
            style="opacity: 0.5"
          >
            暂无消息
          </div>

          <div v-else class="inbox-list">
            <div
              v-for="msg in store.inboxMessages"
              :key="msg.id"
              class="inbox-item"
              :class="{ unread: !msg.isRead }"
              @click="openDetail(msg)"
            >
              <div class="inbox-item-header">
                <span class="inbox-title">{{ msg.title }}</span>
                <NBadge v-if="!msg.isRead" dot type="error" processing />
              </div>
              <div class="inbox-meta">
                <NTag :type="typeTagType(msg.type)" size="tiny" :bordered="false">
                  {{ typeLabel(msg.type) }}
                </NTag>
                <span class="inbox-time">{{ formatTime(msg.createTime) }}</span>
              </div>
              <div class="inbox-snippet">{{ truncate(msg.content, 60) }}</div>
            </div>
          </div>

          <div v-if="store.inboxTotal > store.inboxPageSize" class="py-12 flex justify-center">
            <NPagination
              v-model:page="store.inboxPage"
              :page-size="store.inboxPageSize"
              :item-count="store.inboxTotal"
              @update:page="loadInbox"
            />
          </div>
        </NSpin>
      </template>

      <!-- 详情视图 -->
      <template v-else-if="view === 'detail'">
        <MessageDetail :message="detailMsg" :column="1" />
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<script setup>
import {
  NBadge,
  NButton,
  NDrawer,
  NDrawerContent,
  NPagination,
  NSpin,
  NTag,
} from 'naive-ui'
import { MessageDetail } from '@/components'
import { useNotificationStore } from '@/store/modules/notification'
import { stripMarkdown } from '@/utils/markdown'

const store = useNotificationStore()

const showInbox = computed({
  get: () => store.showInbox,
  set: (val) => { store.showInbox = val },
})

const view = ref('list')
const detailMsg = ref(null)

watch(() => store.detailMessage, (val) => {
  if (val) {
    detailMsg.value = val
    view.value = 'detail'
  }
})
const inboxLoading = ref(false)

watch(() => store.showInbox, (val) => {
  if (val)
    loadInbox()
})

function resetView() {
  view.value = 'list'
  detailMsg.value = null
  store.detailMessage = null
}

async function loadInbox() {
  inboxLoading.value = true
  try {
    await store.fetchInbox()
  }
  catch { /* ignore */ }
  inboxLoading.value = false
}

async function handleMarkAllRead() {
  await store.markAllAsRead()
}

async function openDetail(msg) {
  detailMsg.value = msg
  view.value = 'detail'
  if (!msg.isRead)
    await store.markAsRead(msg.id)
}

// ====== Helpers ======
function typeLabel(type) {
  const map = { SYSTEM: '系统消息', NOTICE: '通知公告', ANNOUNCEMENT: '公告' }
  return map[type] || type
}

function typeTagType(type) {
  const map = { SYSTEM: 'info', NOTICE: 'warning', ANNOUNCEMENT: 'success' }
  return map[type] || 'default'
}

function truncate(str, len) {
  if (!str)
    return ''
  const plain = stripMarkdown(str)
  return plain.length > len ? `${plain.slice(0, len)}...` : plain
}

function formatTime(time) {
  if (!time)
    return ''
  const d = new Date(time)
  const now = new Date()
  const diff = now - d
  if (diff < 60000)
    return '刚刚'
  if (diff < 3600000)
    return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000)
    return `${Math.floor(diff / 3600000)}小时前`
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
.notification-badge :deep(.n-badge-sup) {
  font-size: 10px;
  padding: 0 4px;
  height: 14px;
  min-width: 14px;
  line-height: 14px;
}

.inbox-list {
  display: flex;
  flex-direction: column;
}
.inbox-item {
  padding: 12px 16px;
  cursor: pointer;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  transition: background 0.2s;
}
.inbox-item:hover {
  background: rgba(0, 0, 0, 0.03);
}
:root.dark .inbox-item {
  border-bottom-color: rgba(255, 255, 255, 0.09);
}
:root.dark .inbox-item:hover {
  background: rgba(255, 255, 255, 0.05);
}
.inbox-item.unread .inbox-title {
  font-weight: 600;
}
.inbox-item-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
}
.inbox-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
}
.inbox-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}
.inbox-time {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.4);
}
.inbox-snippet {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.5);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
:root.dark .inbox-time {
  color: rgba(255, 255, 255, 0.4);
}
:root.dark .inbox-snippet {
  color: rgba(255, 255, 255, 0.5);
}
</style>
