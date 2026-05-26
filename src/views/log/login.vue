<template>
  <CommonPage>
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1200"
      :columns="columns"
      :get-data="api.getLogs"
    >
      <MeQueryItem label="用户名" :label-width="50">
        <n-input
          v-model:value="queryItems.username"
          type="text"
          placeholder="请输入用户名"
          clearable
        />
      </MeQueryItem>

      <MeQueryItem label="登录状态" :label-width="70">
        <n-select
          v-model:value="queryItems.success"
          clearable
          :options="[
            { label: '成功', value: 'true' },
            { label: '失败', value: 'false' },
          ]"
        />
      </MeQueryItem>
    </MeCrud>

    <NDrawer v-model:show="drawerVisible" :width="500">
      <NDrawerContent title="登录详情" closable>
        <NSpace vertical :size="16">
          <NCard title="基本信息" size="small" :bordered="true">
            <p><strong>用户名：</strong> {{ currentLog.username || '未知' }}</p>
            <p>
              <strong>登录状态：</strong>
              <NTag :type="currentLog.success ? 'success' : 'error'" :bordered="false" size="small">
                {{ currentLog.success ? '成功' : '失败' }}
              </NTag>
            </p>
            <p><strong>登录信息：</strong> {{ currentLog.message || '无' }}</p>
            <p><strong>IP 地址：</strong> {{ currentLog.ip || '未知' }}</p>
            <p><strong>登录时间：</strong> {{ currentLog.createTime ? formatDateTime(currentLog.createTime) : '无' }}</p>
          </NCard>

          <NCard title="客户端信息" size="small" :bordered="true">
            <p style="word-break: break-all;"><strong>User-Agent：</strong> {{ currentLog.userAgent || '未知' }}</p>
          </NCard>
        </NSpace>
      </NDrawerContent>
    </NDrawer>
  </CommonPage>
</template>

<script setup>
import { NButton, NCard, NDrawer, NDrawerContent, NSpace, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { MeCrud, MeQueryItem } from '@/components'
import { formatDateTime } from '@/utils'
import api from './login-api'

defineOptions({ name: 'LoginLog' })

const $table = ref(null)
const queryItems = ref({})
const drawerVisible = ref(false)
const currentLog = ref({})

function handleViewDetails(row) {
  currentLog.value = row
  drawerVisible.value = true
}

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  { title: 'ID', key: 'id', width: 60, ellipsis: { tooltip: true } },
  { title: '用户名', key: 'username', width: 100, ellipsis: { tooltip: true } },
  {
    title: '登录状态',
    key: 'success',
    width: 90,
    render: ({ success }) =>
      h(
        NTag,
        { type: success ? 'success' : 'error', bordered: false },
        { default: () => (success ? '成功' : '失败') },
      ),
  },
  { title: '登录信息', key: 'message', width: 150, ellipsis: { tooltip: true } },
  { title: 'IP 地址', key: 'ip', width: 130, ellipsis: { tooltip: true } },
  { title: 'User-Agent', key: 'userAgent', width: 250, ellipsis: { tooltip: true } },
  { title: '登录时间', key: 'createTime', width: 180, render: row => h('span', formatDateTime(row.createTime)) },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    fixed: 'right',
    render: row =>
      h(
        NButton,
        { size: 'small', type: 'primary', text: true, onClick: () => handleViewDetails(row) },
        { default: () => '查看详情' },
      ),
  },
]
</script>
