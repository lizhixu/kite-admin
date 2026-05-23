<template>
  <div class="json-viewer" :class="{ 'json-viewer--dark': appStore.isDark }">
    <div v-if="showToolbar" class="json-viewer__toolbar">
      <NSpace size="small" align="center">
        <NTag size="tiny" :type="parsed.kind === 'json' ? 'success' : 'default'" :bordered="false">
          {{ parsed.kind === 'json' ? 'JSON' : 'TEXT' }}
        </NTag>
        <span class="json-viewer__meta">{{ sizeLabel }}</span>
        <template v-if="parsed.prefix">
          <NTag size="tiny" type="info" :bordered="false">前缀</NTag>
        </template>
      </NSpace>
      <NSpace size="small">
        <NButton v-if="parsed.kind === 'json'" size="tiny" quaternary @click="toggleExpand">
          {{ expanded ? '全部折叠' : '全部展开' }}
        </NButton>
        <NButton size="tiny" quaternary @click="copyRaw">复制</NButton>
      </NSpace>
    </div>

    <div v-if="parsed.prefix" class="json-viewer__prefix">{{ parsed.prefix }}</div>

    <VueJsonPretty
      v-if="parsed.kind === 'json'"
      :data="parsed.data"
      :deep="deep"
      :show-length="true"
      :show-line="true"
      :show-icon="true"
      :show-double-quotes="true"
      :virtual="parsed.virtual"
      :height="parsed.virtual ? 400 : undefined"
      :item-height="20"
      :theme="appStore.isDark ? 'dark' : 'light'"
    />
    <pre v-else class="json-viewer__raw">{{ raw }}</pre>
  </div>
</template>

<script setup>
import { NButton, NSpace, NTag } from 'naive-ui'
import { computed, ref } from 'vue'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import { useAppStore } from '@/store'

const props = defineProps({
  raw: { type: [String, Object, Array, Number, Boolean, null], default: '' },
  showToolbar: { type: Boolean, default: true },
  initialDeep: { type: Number, default: 2 },
})

const appStore = useAppStore()

const expanded = ref(false)
const deep = ref(props.initialDeep)

function toggleExpand() {
  expanded.value = !expanded.value
  deep.value = expanded.value ? 99 : props.initialDeep
}

const parsed = computed(() => tryParse(props.raw))

const sizeLabel = computed(() => {
  const s = typeof props.raw === 'string' ? props.raw : JSON.stringify(props.raw ?? '')
  const len = s?.length || 0
  if (len < 1024)
    return `${len} B`
  if (len < 1024 * 1024)
    return `${(len / 1024).toFixed(1)} KB`
  return `${(len / 1024 / 1024).toFixed(2)} MB`
})

function tryParse(input) {
  if (input == null || input === '')
    return { kind: 'text', prefix: '', data: null, virtual: false }
  if (typeof input === 'object')
    return { kind: 'json', prefix: '', data: input, virtual: estimateBig(input) }

  const str = String(input)
  // direct parse
  try {
    const obj = JSON.parse(str)
    if (obj && typeof obj === 'object')
      return { kind: 'json', prefix: '', data: obj, virtual: str.length > 200 * 1024 }
  }
  catch {}

  // try to extract JSON substring; supports "Body: {...}" style
  const start = firstBracket(str)
  const end = lastBracket(str)
  if (start !== -1 && end > start) {
    const candidate = str.substring(start, end + 1)
    try {
      const obj = JSON.parse(candidate)
      if (obj && typeof obj === 'object') {
        return {
          kind: 'json',
          prefix: str.substring(0, start).trim(),
          data: obj,
          virtual: candidate.length > 200 * 1024,
        }
      }
    }
    catch {}
  }

  return { kind: 'text', prefix: '', data: null, virtual: false }
}

function firstBracket(s) {
  const a = s.indexOf('{')
  const b = s.indexOf('[')
  if (a === -1)
    return b
  if (b === -1)
    return a
  return Math.min(a, b)
}

function lastBracket(s) {
  return Math.max(s.lastIndexOf('}'), s.lastIndexOf(']'))
}

function estimateBig(obj) {
  try {
    return JSON.stringify(obj).length > 200 * 1024
  }
  catch {
    return false
  }
}

async function copyRaw() {
  const text = typeof props.raw === 'string' ? props.raw : JSON.stringify(props.raw, null, 2)
  try {
    await navigator.clipboard.writeText(text)
    window.$message?.success('已复制到剪贴板')
  }
  catch {
    window.$message?.error('复制失败')
  }
}
</script>

<style scoped>
.json-viewer {
  border: 1px solid var(--n-border-color, #e5e7eb);
  border-radius: 6px;
  background: var(--n-color, #fafafa);
  font-size: 13px;
  overflow: hidden;
}

.json-viewer--dark {
  background: #1e1e1e;
  border-color: #333;
}

.json-viewer__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px;
  border-bottom: 1px solid var(--n-divider-color, #eee);
  background: var(--n-card-color, #f5f5f5);
}

.json-viewer--dark .json-viewer__toolbar {
  background: #252525;
  border-color: #333;
}

.json-viewer__meta {
  color: var(--n-text-color-3, #999);
  font-size: 12px;
}

.json-viewer__prefix {
  padding: 6px 10px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--n-text-color-2, #666);
  border-bottom: 1px dashed var(--n-divider-color, #ddd);
  white-space: pre-wrap;
  word-break: break-all;
}

.json-viewer :deep(.vjs-tree) {
  padding: 8px 10px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 12.5px;
  line-height: 1.6;
}

.json-viewer__raw {
  margin: 0;
  padding: 10px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 12.5px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 480px;
  overflow: auto;
}
</style>
