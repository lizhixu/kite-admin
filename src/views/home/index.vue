<template>
  <AppPage show-footer>
    <!-- Hero Banner -->
    <n-card class="hero-card">
      <div class="hero-inner">
        <div class="hero-text">
          <n-text depth="3" style="font-size: 13px">
            {{ todayStr }} · {{ greeting }}
          </n-text>
          <div class="hero-title">
            <n-avatar round :size="44" :src="userStore.avatar" />
            <span>{{ userStore.nickName ?? userStore.username }}，欢迎回来</span>
          </div>

          <div class="hero-meta">
            <div class="meta-item">
              <i class="i-fe:user-check meta-icon" />
              <div>
                <div class="meta-label">当前角色</div>
                <div class="meta-value">{{ userStore.currentRole?.name ?? '未分配' }}</div>
              </div>
            </div>
            <n-divider vertical style="height: 36px" />
            <div class="meta-item">
              <i class="i-fe:clock meta-icon" />
              <div>
                <div class="meta-label">上次登录</div>
                <div class="meta-value">
                  {{ lastLogin ? formatTime(lastLogin.createTime) : '首次登录' }}
                  <n-text v-if="lastLogin?.ip" depth="3" style="font-size: 12px; margin-left: 6px">
                    {{ lastLogin.ip }}
                  </n-text>
                </div>
              </div>
            </div>
            <n-divider vertical style="height: 36px" />
            <div class="meta-item">
              <i class="i-fe:log-in meta-icon" />
              <div>
                <div class="meta-label">登录次数</div>
                <div class="meta-value">{{ loginCount }} 次</div>
              </div>
            </div>
          </div>
        </div>
        <img :src="welcomeImg" class="hero-img" alt="welcome" />
      </div>
    </n-card>

    <!-- Bottom row -->
    <n-grid :cols="24" :x-gap="12" class="mt-12">
      <n-gi :span="14">
        <n-card title="最近消息" size="small" class="msg-card">
          <template #header-extra>
            <n-button text type="primary" @click="notificationStore.showInbox = true">
              <template #icon><i class="i-fe:chevron-right" /></template>
            </n-button>
          </template>
          <div v-if="recentMessages.length" class="msg-list">
            <div
              v-for="msg in recentMessages"
              :key="msg.id"
              class="msg-item"
              @click="openMessage(msg)"
            >
              <n-badge :dot="!msg.isRead" :type="msg.isRead ? 'default' : 'error'">
                <span />
              </n-badge>
              <n-text class="msg-title" :depth="msg.isRead ? 3 : 1">{{ msg.title || '无标题' }}</n-text>
              <n-text depth="3" class="msg-time">{{ formatTime(msg.createTime) }}</n-text>
            </div>
          </div>
          <n-empty v-else description="暂无消息" class="py-32" />
        </n-card>
      </n-gi>
      <n-gi :span="10">
        <n-card title="快捷入口" size="small" class="quick-card">
          <template #header-extra>
            <n-button text type="primary" @click="showConfig = true">
              <template #icon><i class="i-fe:settings" /></template>
            </n-button>
          </template>
          <div class="quick-grid">
            <div
              v-for="action in quickActions"
              :key="action.code"
              class="action-tile"
              @click="$router.push(action.path)"
            >
              <i :class="action.icon" class="action-icon" />
              <span class="action-label">{{ action.name }}</span>
            </div>
          </div>
          <n-empty v-if="!quickActions.length" description="暂未配置快捷入口" class="py-16" />
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 配置弹窗 -->
    <n-modal v-model:show="showConfig" preset="card" title="配置快捷入口" style="width: 480px">
      <n-text depth="3" style="display: block; margin-bottom: 12px">
        选择常用菜单作为快捷入口（最多 {{ MAX_QUICK }} 个）
      </n-text>
      <n-checkbox-group v-model:value="tempSelected">
        <div class="config-list">
          <n-checkbox
            v-for="menu in allMenus"
            :key="menu.code"
            :value="menu.code"
            :disabled="!tempSelected.includes(menu.code) && tempSelected.length >= MAX_QUICK"
            class="config-item"
          >
            <i :class="`${menu.icon}?mask`" class="text-14 config-icon" />
            <span>{{ menu.name }}</span>
          </n-checkbox>
        </div>
      </n-checkbox-group>
      <template #action>
        <n-space justify="end">
          <n-button @click="showConfig = false">取消</n-button>
          <n-button type="primary" @click="saveConfig">确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </AppPage>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import { useUserStore, useNotificationStore, usePermissionStore } from '@/store'
import { request } from '@/utils'
import msgApi from '@/views/message/api'
import welcomeImg from '@/assets/images/welcome.svg'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const MAX_QUICK = 9
const STORAGE_KEY = 'home_quick_actions'

const userStore = useUserStore()
const notificationStore = useNotificationStore()
const permissionStore = usePermissionStore()

const lastLogin = ref(null)
const loginCount = ref(0)
const recentMessages = ref([])
const showConfig = ref(false)
const selectedCodes = ref([])
const tempSelected = ref([])

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 6) return '夜深了'
  if (h < 9) return '早上好'
  if (h < 12) return '上午好'
  if (h < 14) return '中午好'
  if (h < 18) return '下午好'
  return '晚上好'
})

const todayStr = computed(() => dayjs().format('YYYY年MM月DD日 dddd'))

/** 递归提取所有有 path 的叶子菜单 */
function collectLeafMenus(perms) {
  const result = []
  for (const p of perms) {
    if (p.type !== 'MENU') continue
    if (p.path && p.show !== false && p.enable !== false) {
      result.push({ code: p.code, name: p.name, icon: p.icon, path: p.path })
    }
    if (p.children?.length) {
      result.push(...collectLeafMenus(p.children))
    }
  }
  return result
}

/** 所有可用菜单（有权限的） */
const allMenus = computed(() => collectLeafMenus(permissionStore.permissions))

/** 当前选中的快捷入口（过滤掉已无权限的） */
const quickActions = computed(() =>
  selectedCodes.value
    .map(code => allMenus.value.find(m => m.code === code))
    .filter(Boolean)
    .map(m => ({ ...m, icon: `${m.icon}?mask` })),
)

/** 打开配置弹窗时，复制当前选中到临时变量 */
watch(showConfig, (val) => {
  if (val) tempSelected.value = [...selectedCodes.value]
})

function saveConfig() {
  selectedCodes.value = [...tempSelected.value]
  localStorage.setItem(STORAGE_KEY, JSON.stringify(selectedCodes.value))
  showConfig.value = false
}

function formatTime(time) {
  if (!time) return ''
  const d = dayjs(time)
  return d.isBefore(dayjs().subtract(1, 'day')) ? d.format('YYYY-MM-DD HH:mm') : d.fromNow()
}

onMounted(async () => {
  // 读取本地存储的配置
  try {
    const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]')
    if (Array.isArray(saved) && saved.length) {
      // 过滤掉已无权限的入口
      const validCodes = allMenus.value.map(m => m.code)
      selectedCodes.value = saved.filter(code => validCodes.includes(code))
      if (selectedCodes.value.length !== saved.length) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(selectedCodes.value))
      }
    } else {
      // 默认选前 6 个有权限的菜单
      selectedCodes.value = allMenus.value.slice(0, MAX_QUICK).map(m => m.code)
    }
  } catch {
    selectedCodes.value = allMenus.value.slice(0, MAX_QUICK).map(m => m.code)
  }

  try {
    const { data } = await request.get('/loginlog/list', {
      params: { username: userStore.username, pageNo: 1, pageSize: 2 },
    })
    const logs = data?.pageData || []
    loginCount.value = data?.total || 0
    lastLogin.value = logs.length > 1 ? logs[1] : logs[0] || null
  } catch { /* ignore */ }

  try {
    const { data } = await msgApi.getMyMessages({ pageNo: 1, pageSize: 20 })
    recentMessages.value = data?.pageData || []
  } catch { /* ignore */ }
})

function openMessage(msg) {
  notificationStore.detailMessage = msg
  notificationStore.showInbox = true
  if (!msg.isRead) {
    msgApi.markRead(msg.id).then(() => {
      msg.isRead = true
      notificationStore.fetchUnreadCount()
    })
  }
}
</script>

<style scoped>
:deep(.n-card) {
  border-radius: 8px;
}
.hero-card :deep(.n-card__content) {
  padding: 0;
}
.hero-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 28px;
  gap: 32px;
  min-height: 200px;
}
.hero-text {
  flex: 1;
  min-width: 0;
}
.hero-title {
  display: flex;
  align-items: center;
  gap: 14px;
  margin: 10px 0 24px;
  font-size: 22px;
  font-weight: 600;
  color: var(--home-title-color);
}
.hero-meta {
  display: flex;
  align-items: center;
  gap: 28px;
}
.meta-item {
  display: flex;
  align-items: center;
  gap: 12px;
}
.meta-icon {
  font-size: 20px;
  color: var(--home-icon-color);
  flex-shrink: 0;
}
.meta-label {
  font-size: 12px;
  color: var(--home-secondary-color);
  line-height: 1.4;
  margin-bottom: 4px;
}
.meta-value {
  font-size: 14px;
  color: var(--home-text-color);
  font-weight: 500;
  line-height: 1.4;
}
.hero-img {
  width: 280px;
  height: auto;
  flex-shrink: 0;
  transition: filter 0.3s;
}

/* Dark mode */
:global(.dark) .hero-img {
  filter: brightness(0.8) saturate(0.9);
}

.msg-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.msg-card :deep(.n-card__content) {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.msg-list {
  height: 248px;
  overflow-y: auto;
}
.msg-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s;
}
.msg-item + .msg-item {
  border-top: 1px solid var(--home-border-color);
}
.msg-item:hover {
  background: var(--home-hover-color);
}
.msg-title {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
  line-height: 1.5;
}
.msg-time {
  flex-shrink: 0;
  font-size: 12px;
  line-height: 1.5;
}

.quick-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.quick-card :deep(.n-card__content) {
  flex: 1;
  display: flex;
  align-items: center;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  width: 100%;
}

.action-tile {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}
.action-tile:hover {
  background: var(--home-hover-color);
}
.action-icon {
  font-size: 20px;
  color: var(--home-icon-color);
}
.action-label {
  font-size: 12px;
  color: var(--home-text-color);
  text-align: center;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.config-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.config-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
}
.config-icon {
  color: var(--home-icon-color);
}

/* Light mode */
:global(:root) {
  --home-title-color: rgba(0, 0, 0, 0.88);
  --home-text-color: rgba(0, 0, 0, 0.85);
  --home-secondary-color: rgba(0, 0, 0, 0.45);
  --home-icon-color: rgba(0, 0, 0, 0.45);
  --home-border-color: rgba(0, 0, 0, 0.06);
  --home-hover-color: rgba(0, 0, 0, 0.03);
}

/* Dark mode */
:global(.dark) {
  --home-title-color: rgba(255, 255, 255, 0.9);
  --home-text-color: rgba(255, 255, 255, 0.88);
  --home-secondary-color: rgba(255, 255, 255, 0.5);
  --home-icon-color: rgba(255, 255, 255, 0.5);
  --home-border-color: rgba(255, 255, 255, 0.08);
  --home-hover-color: rgba(255, 255, 255, 0.06);
}

@media (max-width: 1280px) {
  .hero-img { width: 220px; }
}
@media (max-width: 960px) {
  .hero-img { display: none; }
  .hero-meta { flex-wrap: wrap; gap: 12px 20px; }
}
</style>
