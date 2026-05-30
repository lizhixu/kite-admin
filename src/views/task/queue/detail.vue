<template>
  <CommonPage back>
    <template #title-suffix>
      <NTag :bordered="false" type="info" size="small">
        {{ queue.name }}
      </NTag>
      <NTag v-if="isRegistered" type="success" :bordered="false" size="small" class="ml-8">
        已注册
      </NTag>
      <NTag v-else type="error" :bordered="false" size="small" class="ml-8">
        未注册
      </NTag>
    </template>
    <template #action>
      <NSpace :size="8">
        <NButton v-permission="'AddQueueJob'" type="primary" size="small" @click="openAddJob">
          <i class="i-fe:plus mr-4" />投递任务
        </NButton>
        <NButton v-permission="'KickQueueJob'" type="warning" size="small" secondary @click="handleKickAll">
          <i class="i-fe:rotate-ccw mr-4" />Kick 全部失败
        </NButton>
        <NButton v-permission="'DeleteQueue'" type="error" size="small" secondary @click="handleClearJobs">
          <i class="i-fe:trash-2 mr-4" />清空{{ jobStatus ? statusLabelMap[jobStatus] : '全部' }}
        </NButton>
        <NButton v-permission="'DeleteQueue'" type="error" size="small" secondary @click="handleCleanupOld">
          <i class="i-fe:archive mr-4" />清理 7 天前
        </NButton>
        <NButton v-permission="'EditQueue'" type="primary" size="small" secondary @click="openEdit">
          <i class="i-fe:edit mr-4" />编辑
        </NButton>
        <NButton v-permission="'DeleteQueue'" type="error" size="small" secondary @click="handleDelete">
          <i class="i-fe:trash-2 mr-4" />删除
        </NButton>
      </NSpace>
    </template>

    <!-- Queue 信息条 -->
    <NCard size="small" class="mb-16">
      <NSpace align="center" :size="16" :wrap="true">
        <span v-if="queue.description" class="text-12" style="color: #666">
          {{ queue.description }}
        </span>
        <NTag size="small" :bordered="false">
          并发 {{ queue.concurrency }}
        </NTag>
        <NTag size="small" :bordered="false">
          超时 {{ queue.timeout }}s
        </NTag>
        <NTag size="small" :bordered="false">
          重试 {{ queue.maxRetries }}
        </NTag>
        <NTag size="small" :bordered="false">
          总 {{ queue.totalJobs || 0 }}
        </NTag>
        <NTag size="small" type="success" :bordered="false">
          成 {{ queue.completedJobs || 0 }}
        </NTag>
        <NTag size="small" type="error" :bordered="false">
          败 {{ queue.failedJobs || 0 }}
        </NTag>
        <NSwitch
          size="small"
          :value="queue.status === 'RUNNING'"
          :loading="toggleLoading"
          @update:value="handleToggle"
        >
          <template #checked>运行</template>
          <template #unchecked>暂停</template>
        </NSwitch>
      </NSpace>
    </NCard>

    <!-- Job 过滤栏 -->
    <NSpace align="center" :size="12" :wrap="true" class="mb-12">
      <NTabs
        v-model:value="jobStatus"
        type="segment"
        size="small"
        style="width: 460px"
        @update:value="loadJobs(1)"
      >
        <NTabPane name="" tab="全部" />
        <NTabPane name="PENDING" tab="待执行" />
        <NTabPane name="RUNNING" tab="运行中" />
        <NTabPane name="SUCCESS" tab="成功" />
        <NTabPane name="FAILED" tab="失败" />
      </NTabs>
      <NDatePicker
        v-model:value="jobDateRange"
        type="daterange"
        size="small"
        clearable
        format="yyyy-MM-dd"
        style="width: 240px"
        @update:value="loadJobs(1)"
      />
      <NSpace align="center" :size="6">
        <span class="text-12" style="color: #888">自动刷新</span>
        <NSwitch v-model:value="autoRefresh" size="small" @update:value="onAutoRefreshChange" />
      </NSpace>
      <NButton size="small" @click="loadJobs(jobsPagination.page)">
        <i class="i-fe:refresh-cw mr-4" />刷新
      </NButton>
      <NTag size="small" :bordered="false">
        共 {{ jobsPagination.itemCount }} 条
      </NTag>
      <NTag size="small" type="warning" :bordered="false">
        待执行 {{ jobSummary.pending }}
      </NTag>
      <NTag size="small" type="info" :bordered="false">
        运行中 {{ jobSummary.running }}
      </NTag>
      <NTag size="small" type="success" :bordered="false">
        成功 {{ jobSummary.success }}
      </NTag>
      <NTag size="small" type="error" :bordered="false">
        失败 {{ jobSummary.failed }}
      </NTag>
    </NSpace>

    <!-- Job 表格 -->
    <NDataTable
      :columns="jobColumns"
      :data="jobRows"
      :loading="jobsLoading"
      :pagination="jobsPagination"
      :scroll-x="1500"
      size="small"
      remote
      striped
    />

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
          <n-input v-model:value="modalForm.description" type="textarea" :autosize="{ minRows: 1, maxRows: 3 }" />
        </n-form-item>
        <NGrid :cols="2" :x-gap="12">
          <NGridItem>
            <n-form-item label="并发数" path="concurrency">
              <n-input-number v-model:value="modalForm.concurrency" :min="1" :max="100" style="width: 100%" />
            </n-form-item>
          </NGridItem>
          <NGridItem>
            <n-form-item label="超时(秒)" path="timeout">
              <n-input-number v-model:value="modalForm.timeout" :min="1" :max="3600" style="width: 100%" />
            </n-form-item>
          </NGridItem>
        </NGrid>
        <n-form-item label="默认重试次数" path="maxRetries">
          <n-input-number v-model:value="modalForm.maxRetries" :min="0" :max="10" style="width: 200px" />
        </n-form-item>
        <n-form-item label="状态" path="status">
          <NSwitch
            :value="modalForm.status === 'RUNNING'"
            @update:value="v => modalForm.status = v ? 'RUNNING' : 'PAUSED'"
          >
            <template #checked>运行</template>
            <template #unchecked>暂停</template>
          </NSwitch>
        </n-form-item>
      </n-form>
    </MeModal>

    <!-- 投递任务弹窗 -->
    <NModal v-model:show="addJobVisible" preset="card" title="投递任务" style="width: 640px; max-width: 95vw">
      <n-form label-placement="left" :label-width="100">
        <n-form-item label="Payload (JSON)">
          <n-input
            v-model:value="newJobPayload"
            type="textarea"
            placeholder='例如 {"orderId": 123, "action": "ship"}'
            :autosize="{ minRows: 4, maxRows: 10 }"
          />
        </n-form-item>
        <NGrid :cols="2" :x-gap="12">
          <NGridItem>
            <n-form-item label="最大重试次数">
              <n-input-number v-model:value="newJobRetries" :min="0" :max="10" style="width: 100%" />
            </n-form-item>
          </NGridItem>
          <NGridItem>
            <n-form-item label="优先级">
              <n-input-number v-model:value="newJobPriority" :min="0" :max="9999" style="width: 100%" />
            </n-form-item>
          </NGridItem>
        </NGrid>
        <n-form-item label="延迟到">
          <NDatePicker
            v-model:value="newJobDelayUntil"
            type="datetime"
            clearable
            format="yyyy-MM-dd HH:mm:ss"
            style="width: 100%"
          />
        </n-form-item>
        <n-form-item label="唯一键">
          <n-input v-model:value="newJobUniqueKey" placeholder="同队列内未完成任务唯一，例如 order:123" clearable />
        </n-form-item>
      </n-form>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="addJobVisible = false">取消</NButton>
          <NButton type="primary" :loading="addingJob" @click="submitNewJob">投递</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- Job 详情弹窗 -->
    <NModal v-model:show="jobDetailVisible" preset="card" title="任务详情" style="width: 720px; max-width: 95vw">
      <NDescriptions label-placement="left" :column="2" bordered size="small" class="mb-12">
        <NDescriptionsItem label="ID">
          {{ currentJob.id }}
        </NDescriptionsItem>
        <NDescriptionsItem label="状态">
          <NTag :type="jobStatusTag[currentJob.status]" size="small" :bordered="false">
            {{ currentJob.status }}
          </NTag>
        </NDescriptionsItem>
        <NDescriptionsItem label="重试">
          {{ currentJob.retryCount || 0 }} / {{ currentJob.maxRetries || 0 }}
        </NDescriptionsItem>
        <NDescriptionsItem label="优先级">
          {{ currentJob.priority || 0 }}
        </NDescriptionsItem>
        <NDescriptionsItem label="延迟到">
          {{ currentJob.delayUntil ? formatDateTime(currentJob.delayUntil) : '-' }}
        </NDescriptionsItem>
        <NDescriptionsItem label="唯一键">
          {{ currentJob.uniqueKey || '-' }}
        </NDescriptionsItem>
        <NDescriptionsItem label="耗时">
          {{ formatDuration(currentJob.duration) }}
        </NDescriptionsItem>
        <NDescriptionsItem label="创建时间">
          {{ currentJob.createdAt ? formatDateTime(currentJob.createdAt) : '-' }}
        </NDescriptionsItem>
        <NDescriptionsItem label="开始时间">
          {{ currentJob.startedAt ? formatDateTime(currentJob.startedAt) : '-' }}
        </NDescriptionsItem>
        <NDescriptionsItem label="完成时间" :span="2">
          {{ currentJob.completedAt ? formatDateTime(currentJob.completedAt) : '-' }}
        </NDescriptionsItem>
      </NDescriptions>
      <NTabs type="line" animated default-value="payload">
        <NTabPane name="payload" tab="Payload">
          <div class="tab-scroll">
            <JsonViewer v-if="currentJob.payload" :raw="currentJob.payload" />
            <NEmpty v-else description="无 payload" />
          </div>
        </NTabPane>
        <NTabPane name="result" tab="结果">
          <div class="tab-scroll">
            <JsonViewer v-if="currentJob.result" :raw="currentJob.result" />
            <NEmpty v-else description="无结果" />
          </div>
        </NTabPane>
        <NTabPane name="error" tab="错误" :disabled="!currentJob.error">
          <div class="tab-scroll">
            <JsonViewer v-if="currentJob.error" :raw="currentJob.error" />
            <NEmpty v-else description="无错误" />
          </div>
        </NTabPane>
      </NTabs>
    </NModal>
  </CommonPage>
</template>

<script setup>
import {
  NButton,
  NCard,
  NDataTable,
  NDatePicker,
  NDescriptions,
  NDescriptionsItem,
  NEmpty,
  NGrid,
  NGridItem,
  NInputNumber,
  NModal,
  NSpace,
  NSwitch,
  NTabPane,
  NTabs,
  NTag,
} from 'naive-ui'
import { MeModal, JsonViewer } from '@/components'
import { useCrud } from '@/composables'
import { formatDateTime } from '@/utils'
import api from './api'

defineOptions({ name: 'QueueDetail' })

const route = useRoute()
const router = useRouter()
const queueId = route.params.id

const statusLabelMap = {
  RUNNING: '运行中',
  PAUSED: '已暂停',
  PENDING: '待执行',
  SUCCESS: '成功',
  FAILED: '失败',
}
const jobStatusTag = {
  PENDING: 'warning',
  RUNNING: 'info',
  SUCCESS: 'success',
  FAILED: 'error',
}

function formatDuration(ms) {
  if (ms == null) return '-'
  if (ms < 1000) return `${ms} ms`
  if (ms < 60_000) return `${(ms / 1000).toFixed(2)} s`
  const m = Math.floor(ms / 60_000)
  const s = ((ms % 60_000) / 1000).toFixed(1)
  return `${m}m ${s}s`
}

function isDelayedPending(job) {
  return job.status === 'PENDING' && job.delayUntil && new Date(job.delayUntil).getTime() > Date.now()
}

// ====== Queue 数据 ======
const queue = ref({})
const toggleLoading = ref(false)
const handlerNames = ref(new Set())

const isRegistered = computed(() => handlerNames.value.has(queue.value.name))

async function loadQueue() {
  try {
    const { data } = await api.getOne(queueId)
    queue.value = data || {}
  }
  catch (e) {
    console.error(e)
    $message.error('队列不存在')
    router.replace('/task/queue')
  }
}

async function loadHandlers() {
  try {
    const { data = [] } = await api.getHandlers()
    handlerNames.value = new Set(data || [])
  }
  catch (e) {
    console.error(e)
  }
}

async function handleToggle() {
  toggleLoading.value = true
  try {
    await api.toggle(queueId)
    $message.success('操作成功')
    loadQueue()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    toggleLoading.value = false
  }
}

function handleDelete() {
  $dialog.warning({
    title: '删除队列',
    content: `确定删除队列「${queue.value.name}」？所有任务记录也将被删除。`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await api.delete(queueId)
      $message.success('删除成功')
      router.replace('/task/queue')
    },
  })
}

async function handleKickAll() {
  $dialog.warning({
    title: 'Kick 全部失败任务',
    content: `确定把队列「${queue.value.name}」中所有 FAILED 任务复活为 PENDING？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const { data } = await api.kickAll(queueId)
      $message.success(`已复活 ${data?.affected || 0} 条`)
      loadJobs(jobsPagination.page)
      loadQueue()
    },
  })
}

// ====== 编辑 ======
const { modalRef, modalFormRef, modalForm, handleEdit }
  = useCrud({
    name: '队列',
    initForm: {
      status: 'RUNNING',
      concurrency: 3,
      timeout: 60,
      maxRetries: 0,
    },
    doUpdate: api.update,
    refresh: () => {
      loadQueue()
    },
  })

function openEdit() {
  handleEdit({ ...queue.value })
}

// ====== 投递任务 ======
const addJobVisible = ref(false)
const newJobPayload = ref('')
const newJobRetries = ref(0)
const newJobPriority = ref(0)
const newJobDelayUntil = ref(null)
const newJobUniqueKey = ref('')
const addingJob = ref(false)

function openAddJob() {
  newJobPayload.value = ''
  newJobRetries.value = queue.value.maxRetries || 0
  newJobPriority.value = 0
  newJobDelayUntil.value = null
  newJobUniqueKey.value = ''
  addJobVisible.value = true
}

async function submitNewJob() {
  if (newJobPayload.value) {
    try {
      JSON.parse(newJobPayload.value)
    }
    catch (e) {
      $message.error('Payload 不是合法 JSON')
      return
    }
  }
  addingJob.value = true
  try {
    const { data } = await api.addJob(queueId, {
      payload: newJobPayload.value,
      maxRetries: newJobRetries.value,
      priority: newJobPriority.value || 0,
      delayUntil: newJobDelayUntil.value ? new Date(newJobDelayUntil.value).toISOString() : null,
      uniqueKey: newJobUniqueKey.value || '',
    })
    if (newJobUniqueKey.value && data?.uniqueKey === newJobUniqueKey.value)
      $message.success(`投递成功，任务 #${data.id}`)
    else
      $message.success('投递成功')
    addJobVisible.value = false
    loadJobs(1)
    loadQueue()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    addingJob.value = false
  }
}

// ====== Jobs ======
const jobRows = ref([])
const jobsLoading = ref(false)
const jobStatus = ref('')
const jobDateRange = ref(null)
const autoRefresh = ref(false)
let autoRefreshTimer = null
const jobSummary = ref({ pending: 0, running: 0, success: 0, failed: 0 })
const jobsPagination = reactive({
  page: 1,
  pageSize: 20,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: page => loadJobs(page),
  onUpdatePageSize: (size) => {
    jobsPagination.pageSize = size
    loadJobs(1)
  },
})

function onAutoRefreshChange(v) {
  if (v) {
    autoRefreshTimer = setInterval(() => loadJobs(jobsPagination.page), 3000)
  }
  else {
    stopAutoRefresh()
  }
}

function stopAutoRefresh() {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer)
    autoRefreshTimer = null
  }
}

onBeforeUnmount(stopAutoRefresh)

function buildJobParams() {
  const params = {
    pageNo: jobsPagination.page,
    pageSize: jobsPagination.pageSize,
  }
  if (jobStatus.value) params.status = jobStatus.value
  if (jobDateRange.value && jobDateRange.value.length === 2) {
    params.from = new Date(jobDateRange.value[0]).toISOString()
    const endDay = new Date(jobDateRange.value[1])
    endDay.setHours(23, 59, 59, 999)
    params.to = endDay.toISOString()
  }
  return params
}

async function loadJobs(page) {
  jobsPagination.page = page
  jobsLoading.value = true
  try {
    const { data } = await api.getJobs(queueId, buildJobParams())
    jobRows.value = data?.pageData || []
    jobsPagination.itemCount = data?.total || 0
    let pending = 0
    let running = 0
    let success = 0
    let failed = 0
    for (const r of jobRows.value) {
      if (r.status === 'PENDING') pending++
      else if (r.status === 'RUNNING') running++
      else if (r.status === 'SUCCESS') success++
      else if (r.status === 'FAILED') failed++
    }
    jobSummary.value = { pending, running, success, failed }
  }
  finally {
    jobsLoading.value = false
  }
}

async function handleClearJobs() {
  const label = jobStatus.value ? statusLabelMap[jobStatus.value] : '全部'
  $dialog.warning({
    title: '清空任务',
    content: `确定清空"${label}"任务？此操作不可恢复。`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const params = {}
      if (jobStatus.value) params.status = jobStatus.value
      await api.clearJobs(queueId, params)
      $message.success('清空成功')
      loadJobs(1)
      loadQueue()
    },
  })
}

async function handleCleanupOld() {
  $dialog.warning({
    title: '清理 7 天前完成记录',
    content: '删除 7 天前已完成（completed_at < 7 天前）的任务记录，PENDING / RUNNING 不受影响。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      const before = new Date(Date.now() - 7 * 24 * 3600 * 1000).toISOString()
      await api.clearJobs(queueId, { before })
      $message.success('清理完成')
      loadJobs(1)
      loadQueue()
    },
  })
}

async function handleKickJob(row) {
  await api.kickJob(row.id)
  $message.success('已复活为 PENDING')
  loadJobs(jobsPagination.page)
}

async function handleDeleteJob(row) {
  $dialog.warning({
    title: '删除任务',
    content: `确定删除任务 #${row.id}？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await api.deleteJob(row.id)
      $message.success('删除成功')
      loadJobs(jobsPagination.page)
    },
  })
}

// ====== Job 详情 ======
const jobDetailVisible = ref(false)
const currentJob = ref({})
function showJobDetail(row) {
  currentJob.value = row
  jobDetailVisible.value = true
}

const jobColumns = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: row =>
      h(
        NTag,
        { type: isDelayedPending(row) ? 'default' : jobStatusTag[row.status] || 'default', bordered: false, size: 'small' },
        { default: () => isDelayedPending(row) ? '延迟等待' : row.status },
      ),
  },
  {
    title: '优先级',
    key: 'priority',
    width: 80,
    render: row => row.priority || 0,
  },
  {
    title: '延迟到',
    key: 'delayUntil',
    width: 160,
    render: row => (row.delayUntil ? formatDateTime(row.delayUntil) : '-'),
  },
  {
    title: '唯一键',
    key: 'uniqueKey',
    minWidth: 160,
    ellipsis: { tooltip: true },
    render: row => row.uniqueKey || '-',
  },
  {
    title: 'Payload',
    key: 'payload',
    minWidth: 220,
    ellipsis: { tooltip: true },
    render: row => row.payload || '-',
  },
  {
    title: '结果',
    key: 'result',
    minWidth: 180,
    ellipsis: { tooltip: true },
    render: row => row.result || '-',
  },
  {
    title: '重试',
    key: 'retry',
    width: 70,
    render: row => `${row.retryCount || 0}/${row.maxRetries || 0}`,
  },
  {
    title: '耗时',
    key: 'duration',
    width: 90,
    render: row => formatDuration(row.duration),
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 160,
    render: row => (row.createdAt ? formatDateTime(row.createdAt) : '-'),
  },
  {
    title: '完成时间',
    key: 'completedAt',
    width: 160,
    render: row => (row.completedAt ? formatDateTime(row.completedAt) : '-'),
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    align: 'right',
    fixed: 'right',
    render(row) {
      const buttons = [
        h(
          NButton,
          { size: 'small', type: 'primary', text: true, onClick: () => showJobDetail(row) },
          { default: () => '详情' },
        ),
      ]
      if (row.status === 'FAILED') {
        buttons.push(
          h(
            NButton,
            { size: 'small', type: 'warning', text: true, style: 'margin-left: 12px;', onClick: () => handleKickJob(row) },
            { default: () => 'Kick' },
          ),
        )
      }
      buttons.push(
        h(
          NButton,
          { size: 'small', type: 'error', text: true, style: 'margin-left: 12px;', onClick: () => handleDeleteJob(row) },
          { default: () => '删除' },
        ),
      )
      return buttons
    },
  },
]

// ====== 初始化 ======
onMounted(() => {
  loadQueue()
  loadHandlers()
  loadJobs(1)
})
</script>

<style scoped>
.ml-8 {
  margin-left: 8px;
}
.mb-12 {
  margin-bottom: 12px;
}
.mb-16 {
  margin-bottom: 16px;
}
.text-12 {
  font-size: 12px;
}

.tab-scroll {
  max-height: min(500px, calc(100vh - 340px));
  min-height: 200px;
  overflow-y: auto;
  padding-right: 4px;
}
.tab-scroll::-webkit-scrollbar {
  width: 6px;
}
.tab-scroll::-webkit-scrollbar-thumb {
  background: #d4d4d8;
  border-radius: 3px;
}
</style>
