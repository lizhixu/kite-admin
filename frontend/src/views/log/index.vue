<template>
  <CommonPage>
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1200"
      :columns="columns"
      :get-data="api.getLogs"
    >
      <MeQueryItem label="操作人" :label-width="50">
        <n-input
          v-model:value="queryItems.username"
          type="text"
          placeholder="请输入操作人"
          clearable
        />
      </MeQueryItem>

      <MeQueryItem label="请求方法" :label-width="70">
        <n-select
          v-model:value="queryItems.method"
          clearable
          :options="[
            { label: 'GET', value: 'GET' },
            { label: 'POST', value: 'POST' },
            { label: 'PUT', value: 'PUT' },
            { label: 'DELETE', value: 'DELETE' },
            { label: 'PATCH', value: 'PATCH' },
          ]"
        />
      </MeQueryItem>

      <MeQueryItem label="状态码" :label-width="60">
        <n-input
          v-model:value="queryItems.statusCode"
          type="text"
          placeholder="请输入状态码"
          clearable
        />
      </MeQueryItem>
    </MeCrud>

    <n-drawer v-model:show="drawerVisible" :width="600">
      <n-drawer-content title="日志详情" closable>
        <n-space vertical :size="16">
          <n-card title="基本信息" size="small" :bordered="true">
            <p><strong>操作人：</strong> {{ currentLog.username || '未知' }}</p>
            <p><strong>请求路径：</strong> {{ currentLog.path }}</p>
            <p><strong>状态码：</strong>
              <n-tag :type="currentLog.statusCode === 200 ? 'success' : 'error'" :bordered="false" size="small">{{ currentLog.statusCode }}</n-tag>
            </p>
            <p><strong>操作时间：</strong> {{ currentLog.createTime ? formatDateTime(currentLog.createTime) : '无' }}</p>
          </n-card>

          <n-card title="详细传参" size="small" :bordered="true">
            <n-code v-if="currentLog.params" :code="formatJson(currentLog.params)" language="json" word-wrap />
            <span v-else>无传参</span>
          </n-card>

          <n-card title="响应内容" size="small" :bordered="true">
            <n-code v-if="currentLog.response" :code="formatJson(currentLog.response)" language="json" word-wrap />
            <span v-else>无响应</span>
          </n-card>
        </n-space>
      </n-drawer-content>
    </n-drawer>
  </CommonPage>
</template>

<script setup>
import { NTag, NButton, NDrawer, NDrawerContent, NCard, NSpace, NCode } from 'naive-ui'
import { h, ref, onMounted } from 'vue'
import { MeCrud, MeQueryItem } from '@/components'
import { formatDateTime } from '@/utils'
import api from './api'

defineOptions({ name: 'SysLog' })

const $table = ref(null)
const queryItems = ref({})
const drawerVisible = ref(false)
const currentLog = ref({})

function handleViewDetails(row) {
  currentLog.value = row
  drawerVisible.value = true
}

function formatJson(str) {
  if (!str) return ''
  try {
    // If it is a string containing JSON or an object directly
    if (typeof str === 'object') return JSON.stringify(str, null, 2)
    // Extract JSON part if mixed with text like "Body: {...}"
    if (str.includes('{') && str.includes('}')) {
      const firstBrace = str.indexOf('{')
      const lastBrace = str.lastIndexOf('}')
      const jsonStr = str.substring(firstBrace, lastBrace + 1)
      const obj = JSON.parse(jsonStr)
      // replace the json part with formatted json, keeping text before/after
      return str.substring(0, firstBrace) + JSON.stringify(obj, null, 2) + str.substring(lastBrace + 1)
    }
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch (e) {
    return str
  }
}

onMounted(() => {
  $table.value?.handleSearch()
})

const methodColorMap = {
  GET: 'info',
  POST: 'success',
  PUT: 'warning',
  DELETE: 'error',
  PATCH: 'warning',
}

const columns = [
  { title: 'ID', key: 'id', width: 60, ellipsis: { tooltip: true } },
  { title: '操作人', key: 'username', width: 100, ellipsis: { tooltip: true } },
  {
    title: '请求方法',
    key: 'method',
    width: 100,
    render: ({ method }) =>
      h(
        NTag,
        { type: methodColorMap[method] || 'default', bordered: false },
        { default: () => method },
      ),
  },
  { title: '请求路径', key: 'path', width: 250, ellipsis: { tooltip: true } },
  { title: 'IP地址', key: 'ip', width: 130, ellipsis: { tooltip: true } },
  {
    title: '状态码',
    key: 'statusCode',
    width: 80,
    render: ({ statusCode }) =>
      h(
        NTag,
        { type: statusCode === 200 ? 'success' : 'error', bordered: false },
        { default: () => statusCode },
      ),
  },
  { title: '操作时间', key: 'createTime', width: 180, render: row => h('span', formatDateTime(row.createTime)) },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    fixed: 'right',
    render: row =>
      h(
        NButton,
        { size: 'small', type: 'primary', text: true, onClick: () => handleViewDetails(row) },
        { default: () => '查看详情' }
      ),
  },
]
</script>
