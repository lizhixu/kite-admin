<template>
  <div v-if="message" class="message-detail">
    <n-descriptions :column="column" label-placement="left" bordered size="small">
      <n-descriptions-item label="类型">
        <n-tag :type="typeTagType(message.type)" size="small" :bordered="false">
          {{ typeLabel(message.type) }}
        </n-tag>
      </n-descriptions-item>
      <n-descriptions-item v-if="showTarget" label="目标">
        <n-tag :type="targetTagType" size="small" :bordered="false">
          {{ targetLabel }}
        </n-tag>
      </n-descriptions-item>
      <n-descriptions-item label="发送人">
        {{ message.senderName }}
      </n-descriptions-item>
      <n-descriptions-item label="发送时间">
        {{ formatTime(message.createTime) }}
      </n-descriptions-item>
    </n-descriptions>
    <n-divider style="margin: 0" />
    <div class="content" v-html="renderedContent" />
  </div>
</template>

<script setup>
import { NDescriptions, NDescriptionsItem, NDivider, NTag } from 'naive-ui'
import { renderMarkdown } from '@/utils/markdown'

const props = defineProps({
  message: { type: Object, default: null },
  showTarget: { type: Boolean, default: false },
  column: { type: Number, default: 2 },
})

const renderedContent = computed(() => renderMarkdown(props.message?.content || ''))

const targetTagType = computed(() => {
  const map = { ALL: 'default', ROLE: 'warning', USER: 'info' }
  return map[props.message?.targetType] || 'default'
})

const targetLabel = computed(() => {
  const map = { ALL: '全员广播', ROLE: '按角色', USER: '指定用户' }
  return map[props.message?.targetType] || props.message?.targetType
})

function typeLabel(type) {
  const map = { SYSTEM: '系统消息', NOTICE: '通知公告', ANNOUNCEMENT: '公告' }
  return map[type] || type
}

function typeTagType(type) {
  const map = { SYSTEM: 'info', NOTICE: 'warning', ANNOUNCEMENT: 'success' }
  return map[type] || 'default'
}

function formatTime(time) {
  if (!time) return '-'
  const d = new Date(time)
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}
</script>

<style scoped>
.content :deep(h1),
.content :deep(h2),
.content :deep(h3),
.content :deep(h4) {
  margin: 1em 0 0.5em;
  font-weight: 600;
  line-height: 1.4;
}
.content :deep(h1) {
  font-size: 1.5em;
}
.content :deep(h2) {
  font-size: 1.3em;
}
.content :deep(h3) {
  font-size: 1.15em;
}
.content :deep(p) {
  margin: 0.5em 0;
  line-height: 1.8;
}
.content :deep(ul),
.content :deep(ol) {
  padding-left: 1.5em;
  margin: 0.5em 0;
}
.content :deep(li) {
  margin: 0.25em 0;
  line-height: 1.6;
}
.content :deep(blockquote) {
  margin: 0.8em 0;
  padding: 0.5em 1em;
  border-left: 4px solid #18a058;
  background: #f8fdf9;
  color: #555;
}
.content :deep(code) {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 0.9em;
  font-family: 'Fira Code', 'Consolas', monospace;
}
.content :deep(pre) {
  background: #f6f8fa;
  border-radius: 6px;
  padding: 16px;
  overflow-x: auto;
  margin: 0.8em 0;
}
.content :deep(pre code) {
  background: none;
  padding: 0;
  font-size: 0.85em;
  line-height: 1.6;
}
.content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 0.8em 0;
}
.content :deep(th),
.content :deep(td) {
  border: 1px solid #e0e0e0;
  padding: 8px 12px;
  text-align: left;
}
.content :deep(th) {
  background: #f6f8fa;
  font-weight: 600;
}
.content :deep(hr) {
  border: none;
  border-top: 1px solid #e0e0e0;
  margin: 1.2em 0;
}
.content :deep(img) {
  max-width: 100%;
  border-radius: 4px;
}
.content :deep(a) {
  color: #18a058;
  text-decoration: none;
}
.content :deep(a:hover) {
  text-decoration: underline;
}
</style>
