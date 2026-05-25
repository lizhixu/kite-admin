<template>
  <CommonPage>
    <!-- 顶部统计 -->
    <NGrid :cols="6" :x-gap="12" :y-gap="12" responsive="screen" class="mb-16">
      <NGridItem :span="1">
        <NCard size="small" hoverable>
          <NStatistic label="队列总数" :value="stats.total || 0" />
        </NCard>
      </NGridItem>
      <NGridItem :span="1">
        <NCard size="small" hoverable>
          <NStatistic label="运行中">
            <span class="text-success">{{ stats.running || 0 }}</span>
          </NStatistic>
        </NCard>
      </NGridItem>
      <NGridItem :span="1">
        <NCard size="small" hoverable>
          <NStatistic label="已暂停">
            <span style="color: #999">{{ stats.paused || 0 }}</span>
          </NStatistic>
        </NCard>
      </NGridItem>
      <NGridItem :span="1">
        <NCard size="small" hoverable>
          <NStatistic label="待执行">
            <span class="text-warning">{{ stats.jobPending || 0 }}</span>
          </NStatistic>
        </NCard>
      </NGridItem>
      <NGridItem :span="1">
        <NCard size="small" hoverable>
          <NStatistic label="今日成功">
            <span class="text-success">{{ stats.successToday || 0 }}</span>
          </NStatistic>
        </NCard>
      </NGridItem>
      <NGridItem :span="1">
        <NCard size="small" hoverable>
          <NStatistic label="今日失败">
            <span class="text-error">{{ stats.failedToday || 0 }}</span>
          </NStatistic>
        </NCard>
      </NGridItem>
    </NGrid>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1200"
      :columns="columns"
      :get-data="api.read"
      @on-data-change="loadStats"
    >
      <MeQueryItem label="名称" :label-width="50">
        <n-input
          v-model:value="queryItems.name"
          type="text"
          placeholder="请输入队列名"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="状态" :label-width="50">
        <n-select
          v-model:value="queryItems.status"
          clearable
          :options="statusOptions"
          style="width: 140px"
        />
      </MeQueryItem>
    </MeCrud>

    <!-- 编辑弹窗 -->
    <MeModal ref="modalRef" width="560px">
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="100"
        :model="modalForm"
      >
        <n-form-item label="队列名">
          <n-input :value="modalForm.name" disabled />
        </n-form-item>
        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="modalForm.description"
            type="textarea"
            :autosize="{ minRows: 1, maxRows: 3 }"
          />
        </n-form-item>
        <NGrid :cols="2" :x-gap="12">
          <NGridItem>
            <n-form-item label="并发数" path="concurrency">
              <NInputNumber
                v-model:value="modalForm.concurrency"
                :min="1"
                :max="100"
                style="width: 100%"
              />
            </n-form-item>
          </NGridItem>
          <NGridItem>
            <n-form-item label="超时(秒)" path="timeout">
              <NInputNumber
                v-model:value="modalForm.timeout"
                :min="1"
                :max="3600"
                style="width: 100%"
              />
            </n-form-item>
          </NGridItem>
        </NGrid>
        <n-form-item label="默认重试次数" path="maxRetries">
          <NInputNumber
            v-model:value="modalForm.maxRetries"
            :min="0"
            :max="10"
            style="width: 200px"
          />
        </n-form-item>
        <n-form-item label="状态" path="status">
          <NSwitch
            :value="modalForm.status === 'RUNNING'"
            @update:value="v => (modalForm.status = v ? 'RUNNING' : 'PAUSED')"
          >
            <template #checked>运行</template>
            <template #unchecked>暂停</template>
          </NSwitch>
        </n-form-item>
      </n-form>
    </MeModal>
  </CommonPage>
</template>

<script setup>
import {
  NButton,
  NCard,
  NGrid,
  NGridItem,
  NInputNumber,
  NSpace,
  NStatistic,
  NSwitch,
  NTag,
} from 'naive-ui'
import { MeCrud, MeModal, MeQueryItem } from '@/components'
import { useCrud } from '@/composables'
import { withPermission } from '@/directives'
import api from './api'

defineOptions({ name: 'QueueMgt' })

const router = useRouter()
const $table = ref(null)
const queryItems = ref({})

const statusOptions = [
  { label: '运行中', value: 'RUNNING' },
  { label: '已暂停', value: 'PAUSED' },
]

// ====== 已注册 handler ======
const handlerNames = ref(new Set())
async function loadHandlers() {
  try {
    const { data = [] } = await api.getHandlers()
    handlerNames.value = new Set(data || [])
  }
  catch {
    // ignore
  }
}
function isRegistered(name) {
  return handlerNames.value.has(name)
}

// ====== Stats ======
const stats = ref({})
async function loadStats() {
  try {
    const { data } = await api.stats()
    stats.value = data || {}
  }
  catch {
    // ignore
  }
}

onMounted(() => {
  loadHandlers()
  loadStats()
})

// ====== CRUD ======
const {
  modalRef,
  modalFormRef,
  modalForm,
  handleDelete,
  handleEdit,
} = useCrud({
  name: '队列',
  initForm: {
    status: 'RUNNING',
    concurrency: 3,
    timeout: 60,
    maxRetries: 0,
  },
  doDelete: api.delete,
  doUpdate: api.update,
  refresh: (_, keepCurrentPage) => {
    $table.value?.handleSearch(keepCurrentPage)
    loadStats()
  },
})

onMounted(() => $table.value?.handleSearch())

// ====== 表格列 ======
const columns = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '队列名',
    key: 'name',
    width: 200,
    ellipsis: { tooltip: true },
    render: row =>
      h(
        'a',
        {
          class: 'text-primary cursor-pointer',
          onClick: () => router.push(`/task/queue/${row.id}`),
        },
        row.name,
      ),
  },
  {
    title: '已注册',
    key: 'registered',
    width: 80,
    render: row =>
      isRegistered(row.name)
        ? h(
            NTag,
            { type: 'success', bordered: false, size: 'small' },
            { default: () => '是' },
          )
        : h(
            NTag,
            { type: 'error', bordered: false, size: 'small' },
            { default: () => '否' },
          ),
  },
  {
    title: '描述',
    key: 'description',
    minWidth: 180,
    ellipsis: { tooltip: true },
  },
  { title: '并发', key: 'concurrency', width: 60 },
  {
    title: '超时',
    key: 'timeout',
    width: 70,
    render: row => `${row.timeout}s`,
  },
  {
    title: '统计',
    key: 'counts',
    width: 200,
    render(row) {
      return h(NSpace, { size: 4 }, () => [
        h(
          NTag,
          { size: 'small', bordered: false },
          { default: () => `总 ${row.totalJobs || 0}` },
        ),
        h(
          NTag,
          { size: 'small', type: 'success', bordered: false },
          { default: () => `成 ${row.completedJobs || 0}` },
        ),
        h(
          NTag,
          { size: 'small', type: 'error', bordered: false },
          { default: () => `败 ${row.failedJobs || 0}` },
        ),
      ])
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    fixed: 'right',
    render: row =>
      h(
        NSwitch,
        {
          size: 'small',
          rubberBand: false,
          value: row.status === 'RUNNING',
          loading: !!row.toggleLoading,
          onUpdateValue: () => handleToggle(row),
        },
        { checked: () => '运行', unchecked: () => '暂停' },
      ),
  },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    align: 'right',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            onClick: () => router.push(`/task/queue/${row.id}`),
          },
          { default: () => '详情' },
        ),
        withPermission(
          h(
            NButton,
            {
              size: 'small',
              type: 'warning',
              secondary: true,
              style: 'margin-left: 8px;',
              onClick: () => handleKickAll(row),
            },
            { default: () => 'Kick' },
          ),
          'KickQueueJob',
        ),
        withPermission(
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              style: 'margin-left: 8px;',
              onClick: () => handleEdit(row),
            },
            { default: () => '编辑' },
          ),
          'EditQueue',
        ),
        withPermission(
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              style: 'margin-left: 8px;',
              onClick: () => handleDelete(row.id),
            },
            { default: () => '删除' },
          ),
          'DeleteQueue',
        ),
      ]
    },
  },
]

// ====== 操作 ======
async function handleToggle(row) {
  row.toggleLoading = true
  try {
    await api.toggle(row.id)
    $message.success('操作成功')
    $table.value?.handleSearch(true)
    loadStats()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    row.toggleLoading = false
  }
}

async function handleKickAll(row) {
  $dialog.warning({
    title: 'Kick 全部失败任务',
    content: `确定把队列「${row.name}」中所有 FAILED 任务复活为 PENDING？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const { data } = await api.kickAll(row.id)
      $message.success(`已复活 ${data?.affected || 0} 条`)
      $table.value?.handleSearch(true)
      loadStats()
    },
  })
}
</script>

<style scoped>
.text-success {
  color: #18a058;
}
.text-error {
  color: #d03050;
}
.text-warning {
  color: #f0a020;
}
.text-primary {
  color: #2080f0;
}
.ml-8 {
  margin-left: 8px;
}
</style>
