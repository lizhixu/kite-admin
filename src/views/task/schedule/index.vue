<template>
  <CommonPage>
    <template #action>
      <NSpace>
        <NButton v-permission="'EditTask'" :disabled="!checkedIds.length" type="primary" secondary @click="handleBulkToggle(true)">
          <i class="i-fe:play mr-4" />批量启用
        </NButton>
        <NButton v-permission="'EditTask'" :disabled="!checkedIds.length" secondary @click="handleBulkToggle(false)">
          <i class="i-fe:pause mr-4" />批量停用
        </NButton>
        <NButton v-permission="'DeleteTask'" :disabled="!checkedIds.length" type="error" secondary @click="handleBulkDelete">
          <i class="i-material-symbols:delete-outline mr-4" />批量删除
        </NButton>
        <NButton v-permission="'AddTask'" type="primary" @click="handleAdd()">
          <i class="i-material-symbols:add mr-4 text-18" />新增任务
        </NButton>
      </NSpace>
    </template>

    <!-- 顶部统计 -->
    <NGrid :cols="6" :x-gap="12" :y-gap="12" responsive="screen" class="mb-16">
      <NGridItem :span="1">
        <NCard size="small" hoverable><NStatistic label="任务总数" :value="stats.total || 0" /></NCard>
      </NGridItem>
      <NGridItem :span="1">
        <NCard size="small" hoverable>
          <NStatistic label="启用中">
            <span class="text-success">{{ stats.enabled || 0 }}</span>
          </NStatistic>
        </NCard>
      </NGridItem>
      <NGridItem :span="1">
        <NCard size="small" hoverable>
          <NStatistic label="已停用">
            <span style="color: #999">{{ stats.disabled || 0 }}</span>
          </NStatistic>
        </NCard>
      </NGridItem>
      <NGridItem :span="1">
        <NCard size="small" hoverable><NStatistic label="今日执行" :value="stats.totalToday || 0" /></NCard>
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
          <NStatistic label="今日失败/超时">
            <span class="text-error">{{ (stats.failedToday || 0) + (stats.timeoutToday || 0) }}</span>
          </NStatistic>
        </NCard>
      </NGridItem>
    </NGrid>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1700"
      :columns="columns"
      :get-data="api.read"
      @on-checked="onChecked"
      @on-data-change="loadStats"
    >
      <MeQueryItem label="任务名" :label-width="50">
        <n-input v-model:value="queryItems.name" type="text" placeholder="请输入任务名" clearable />
      </MeQueryItem>
      <MeQueryItem label="类型" :label-width="50">
        <n-select v-model:value="queryItems.type" clearable :options="typeOptions" style="width: 140px" />
      </MeQueryItem>
    </MeCrud>

    <!-- 新增/编辑 -->
    <MeModal ref="modalRef" width="720px">
      <div class="task-form-scroll">
        <n-form
          ref="modalFormRef"
          label-placement="left"
          label-align="left"
          :label-width="90"
          :model="modalForm"
        >
          <n-form-item
            label="任务名"
            path="name"
            :rule="{ required: true, message: '请输入任务名', trigger: ['input', 'blur'] }"
          >
            <n-input v-model:value="modalForm.name" placeholder="便于识别的中文名" />
          </n-form-item>

          <n-form-item
            label="Cron 表达式"
            path="spec"
            :rule="{ required: true, message: '请输入 cron 表达式', trigger: ['input', 'blur'] }"
          >
            <NSpace vertical :size="6" style="width: 100%">
              <n-input v-model:value="modalForm.spec" placeholder="如 */5 * * * * 每 5 分钟" @update:value="onSpecChange" />
              <NSpace :size="6" :wrap="true">
                <NTag
                  v-for="p in cronPresets"
                  :key="p.spec"
                  checkable
                  size="small"
                  @click="applyPreset(p.spec)"
                >
                  {{ p.label }}
                </NTag>
              </NSpace>
              <div v-if="nextRuns.length" class="text-12" style="color: var(--n-text-color-3, #888)">
                下次执行：
                <span v-for="(t, i) in nextRuns" :key="i" style="margin-right: 12px">
                  {{ formatDateTime(t) }}
                </span>
              </div>
              <div v-else-if="specError" class="text-12 text-error">
                {{ specError }}
              </div>
            </NSpace>
          </n-form-item>

          <NGrid :cols="2" :x-gap="12">
            <NGridItem>
              <n-form-item
                label="任务类型"
                path="type"
                :rule="{ required: true, message: '请选择任务类型', trigger: 'change' }"
              >
                <n-select v-model:value="modalForm.type" :options="typeOptions" />
              </n-form-item>
            </NGridItem>
            <NGridItem>
              <n-form-item label="超时(秒)" path="timeout">
                <n-input-number v-model:value="modalForm.timeout" :min="1" :max="3600" style="width: 100%" />
              </n-form-item>
            </NGridItem>
          </NGrid>

          <template v-if="modalForm.type === 'HTTP'">
            <n-form-item
              label="URL"
              path="command"
              :rule="{ required: true, message: '请输入请求地址', trigger: ['input', 'blur'] }"
            >
              <n-input v-model:value="modalForm.command" placeholder="https://example.com/api" />
            </n-form-item>

            <NGrid :cols="2" :x-gap="12">
              <NGridItem>
                <n-form-item label="请求方法" path="httpMethod">
                  <n-select v-model:value="modalForm.httpMethod" :options="httpMethodOptions" />
                </n-form-item>
              </NGridItem>
              <NGridItem>
                <n-form-item label="状态" path="enabled">
                  <NSwitch v-model:value="modalForm.enabled">
                    <template #checked>启用</template>
                    <template #unchecked>停用</template>
                  </NSwitch>
                </n-form-item>
              </NGridItem>
            </NGrid>

            <n-form-item label="请求头" path="httpHeaders">
              <n-input
                v-model:value="modalForm.httpHeaders"
                type="textarea"
                placeholder='JSON 字符串，例如 {"Authorization":"Bearer xxx"}'
                :autosize="{ minRows: 2, maxRows: 4 }"
              />
            </n-form-item>

            <n-form-item label="请求体" path="httpBody">
              <n-input
                v-model:value="modalForm.httpBody"
                type="textarea"
                placeholder="POST/PUT 时的请求体"
                :autosize="{ minRows: 2, maxRows: 4 }"
              />
            </n-form-item>
          </template>

          <n-form-item
            v-if="modalForm.type === 'SHELL'"
            label="Shell 命令"
            path="command"
            :rule="{ required: true, message: '请输入命令', trigger: ['input', 'blur'] }"
          >
            <n-input
              v-model:value="modalForm.command"
              type="textarea"
              placeholder="例如 ls /tmp"
              :autosize="{ minRows: 2, maxRows: 4 }"
            />
          </n-form-item>

          <n-form-item
            v-if="modalForm.type === 'FUNC'"
            label="内置函数"
            path="command"
            :rule="{ required: true, message: '请选择内置函数', trigger: 'change' }"
          >
            <n-select
              v-model:value="modalForm.command"
              :options="funcOptions"
              placeholder="选择已注册的内置函数"
            />
          </n-form-item>

          <n-form-item label="描述" path="description">
            <n-input v-model:value="modalForm.description" type="textarea" :autosize="{ minRows: 1, maxRows: 3 }" />
          </n-form-item>

          <n-form-item v-if="modalForm.type !== 'HTTP'" label="状态" path="enabled">
            <NSwitch v-model:value="modalForm.enabled">
              <template #checked>启用</template>
              <template #unchecked>停用</template>
            </NSwitch>
          </n-form-item>
        </n-form>
      </div>
    </MeModal>

    <!-- 任务详情 -->
    <NDrawer v-model:show="detailVisible" :width="540">
      <NDrawerContent :title="`任务详情 - ${currentTask.name || ''}`" closable>
        <NDescriptions label-placement="left" :column="1" bordered size="small">
          <NDescriptionsItem label="ID">{{ currentTask.id }}</NDescriptionsItem>
          <NDescriptionsItem label="任务名">{{ currentTask.name }}</NDescriptionsItem>
          <NDescriptionsItem label="类型">
            <NTag :type="typeTagMap[currentTask.type]" size="small" :bordered="false">
              {{ currentTask.type }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem label="Cron">
            <code>{{ currentTask.spec }}</code>
          </NDescriptionsItem>
          <NDescriptionsItem label="执行内容">
            <NCode :code="currentTask.command || ''" word-wrap />
          </NDescriptionsItem>
          <NDescriptionsItem v-if="currentTask.type === 'HTTP'" label="HTTP 方法">
            {{ currentTask.httpMethod }}
          </NDescriptionsItem>
          <NDescriptionsItem v-if="currentTask.httpHeaders" label="请求头">
            <JsonViewer :raw="currentTask.httpHeaders" />
          </NDescriptionsItem>
          <NDescriptionsItem v-if="currentTask.httpBody" label="请求体">
            <JsonViewer :raw="currentTask.httpBody" />
          </NDescriptionsItem>
          <NDescriptionsItem label="超时">{{ currentTask.timeout }} 秒</NDescriptionsItem>
          <NDescriptionsItem label="状态">
            <NTag :type="currentTask.enabled ? 'success' : 'default'" size="small" :bordered="false">
              {{ currentTask.enabled ? '启用' : '停用' }}
            </NTag>
          </NDescriptionsItem>
          <NDescriptionsItem label="上次执行">
            {{ currentTask.lastRunAt ? formatDateTime(currentTask.lastRunAt) : '-' }}
          </NDescriptionsItem>
          <NDescriptionsItem label="下次执行">
            {{ currentTask.nextRunAt ? formatDateTime(currentTask.nextRunAt) : '-' }}
          </NDescriptionsItem>
          <NDescriptionsItem label="描述">{{ currentTask.description || '-' }}</NDescriptionsItem>
        </NDescriptions>
      </NDrawerContent>
    </NDrawer>

    <!-- 日志抽屉 -->
    <NDrawer v-model:show="logDrawerVisible" :width="900">
      <NDrawerContent :title="`执行日志 - ${currentTask.name || ''}`" closable>
        <NSpace vertical :size="12">
          <!-- 汇总 chip -->
          <NSpace :size="8" align="center">
            <NTag size="small" :bordered="false">共 {{ logPagination.itemCount }} 条</NTag>
            <NTag size="small" type="success" :bordered="false">成功 {{ logSummary.success }}</NTag>
            <NTag size="small" type="error" :bordered="false">失败 {{ logSummary.failed }}</NTag>
            <NTag size="small" type="warning" :bordered="false">超时 {{ logSummary.timeout }}</NTag>
            <span v-if="logSummary.avg" class="text-12" style="color: #888">
              平均耗时 {{ formatDuration(logSummary.avg) }}
            </span>
          </NSpace>

          <!-- 筛选与工具栏 -->
          <NSpace align="center" :size="8" :wrap="true">
            <NTabs
              v-model:value="logStatus"
              type="segment"
              size="small"
              style="width: 320px"
              @update:value="loadLogs(1)"
            >
              <NTabPane name="" tab="全部" />
              <NTabPane name="SUCCESS" tab="成功" />
              <NTabPane name="FAILED" tab="失败" />
              <NTabPane name="TIMEOUT" tab="超时" />
            </NTabs>
            <n-select
              v-model:value="logTrigger"
              clearable
              placeholder="触发方式"
              size="small"
              style="width: 140px"
              :options="triggerOptions"
              @update:value="loadLogs(1)"
            />
            <NButton size="small" @click="loadLogs(logPagination.page)">
              <i class="i-fe:refresh-cw mr-4" />刷新
            </NButton>
            <NButton size="small" type="primary" secondary @click="handleRun(currentTask)">
              <i class="i-fe:play mr-4" />再执行一次
            </NButton>
            <NSpace align="center" :size="6">
              <span class="text-12" style="color: #888">自动刷新</span>
              <NSwitch v-model:value="autoRefresh" size="small" @update:value="onAutoRefreshChange" />
            </NSpace>
          </NSpace>

          <NDataTable
            :columns="logColumns"
            :data="logRows"
            :loading="logLoading"
            :pagination="logPagination"
            :scroll-x="1100"
            size="small"
            remote
            striped
          />
        </NSpace>
      </NDrawerContent>
    </NDrawer>

    <!-- 输出详情 -->
    <NModal v-model:show="outputVisible" preset="card" title="输出详情" style="width: 880px; max-width: 95vw">
      <NDescriptions label-placement="left" :column="2" bordered size="small" class="mb-12">
        <NDescriptionsItem label="任务名">{{ currentLog.taskName }}</NDescriptionsItem>
        <NDescriptionsItem label="触发方式">
          {{ currentLog.trigger === 'MANUAL' ? '手动' : '定时' }}
        </NDescriptionsItem>
        <NDescriptionsItem label="状态">
          <NTag :type="statusTagMap[currentLog.status]" size="small" :bordered="false">
            {{ currentLog.status }}
          </NTag>
        </NDescriptionsItem>
        <NDescriptionsItem label="耗时">{{ formatDuration(currentLog.duration) }}</NDescriptionsItem>
        <NDescriptionsItem label="开始">
          {{ currentLog.startTime ? formatDateTime(currentLog.startTime) : '-' }}
        </NDescriptionsItem>
        <NDescriptionsItem label="结束">
          {{ currentLog.endTime ? formatDateTime(currentLog.endTime) : '-' }}
        </NDescriptionsItem>
      </NDescriptions>
      <NTabs type="line" animated default-value="steps">
        <NTabPane name="steps" tab="执行过程">
          <template v-if="parsedSteps.length">
            <!-- 顶部统计 + 级别筛选 -->
            <div class="step-toolbar">
              <NSpace :size="6" align="center">
                <NTag size="small" :bordered="false">共 {{ parsedSteps.length }} 步</NTag>
                <NTag size="small" type="info" :bordered="false">INFO {{ stepCounts.INFO }}</NTag>
                <NTag size="small" type="warning" :bordered="false">WARN {{ stepCounts.WARN }}</NTag>
                <NTag size="small" type="error" :bordered="false">ERROR {{ stepCounts.ERROR }}</NTag>
                <span v-if="stepElapsed" class="text-12 text-muted">总耗时 {{ stepElapsed }}</span>
              </NSpace>
              <NSpace :size="4">
                <NTag
                  v-for="lv in stepLevelFilters"
                  :key="lv.value"
                  checkable
                  size="small"
                  :checked="stepFilter === lv.value"
                  @update:checked="stepFilter = lv.value"
                >
                  {{ lv.label }}
                </NTag>
              </NSpace>
            </div>

            <!-- 时间轴 -->
            <div class="step-scroll">
              <NTimeline v-if="filteredSteps.length" size="small">
                <NTimelineItem
                  v-for="s in filteredSteps"
                  :key="s.idx"
                  :type="stepTypeMap[s.level] || 'default'"
                  :line-type="s.level === 'ERROR' ? 'dashed' : 'default'"
                >
                  <template #icon>
                    <i :class="['step-icon', stepIconMap[s.level] || stepIconMap.INFO]" />
                  </template>
                  <template #header>
                    <div class="step-header">
                      <span class="step-idx">#{{ s.idx }}</span>
                      <NTag size="tiny" :type="stepTypeMap[s.level] || 'default'" :bordered="false">
                        {{ s.level }}
                      </NTag>
                      <span class="step-time">{{ formatStepTime(s.time) }}</span>
                      <span class="step-delta">+{{ s.delta }} ms</span>
                    </div>
                  </template>
                  <div class="step-msg">{{ s.message }}</div>
                </NTimelineItem>
              </NTimeline>
              <NEmpty v-else description="当前级别下无记录" size="small" />
            </div>
          </template>
          <NEmpty v-else description="无过程记录" />
        </NTabPane>
        <NTabPane name="output" tab="输出">
          <div class="tab-scroll">
            <JsonViewer v-if="currentLog.output" :raw="currentLog.output" />
            <NEmpty v-else description="无输出" />
          </div>
        </NTabPane>
        <NTabPane name="error" tab="错误" :disabled="!currentLog.error">
          <div class="tab-scroll">
            <JsonViewer v-if="currentLog.error" :raw="currentLog.error" />
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
  NCode,
  NDataTable,
  NDescriptions,
  NDescriptionsItem,
  NDrawer,
  NDrawerContent,
  NEmpty,
  NGrid,
  NGridItem,
  NInputNumber,
  NModal,
  NSpace,
  NStatistic,
  NSwitch,
  NTabPane,
  NTabs,
  NTag,
  NTimeline,
  NTimelineItem,
} from 'naive-ui'
import { MeCrud, MeModal, MeQueryItem, JsonViewer } from '@/components'
import { useCrud } from '@/composables'
import { withPermission } from '@/directives'
import { formatDateTime } from '@/utils'
import api from './api'

defineOptions({ name: 'ScheduleTask' })

const $table = ref(null)
const queryItems = ref({})
const checkedIds = ref([])

const typeOptions = [
  { label: 'HTTP 请求', value: 'HTTP' },
  { label: 'Shell 命令', value: 'SHELL' },
  { label: '内置函数', value: 'FUNC' },
]

const httpMethodOptions = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'].map(v => ({ label: v, value: v }))

const statusOptions = [
  { label: '成功', value: 'SUCCESS' },
  { label: '失败', value: 'FAILED' },
  { label: '超时', value: 'TIMEOUT' },
]

const triggerOptions = [
  { label: '手动', value: 'MANUAL' },
  { label: '定时', value: 'CRON' },
]

function formatDuration(ms) {
  if (ms == null)
    return '-'
  if (ms < 1000)
    return `${ms} ms`
  if (ms < 60_000)
    return `${(ms / 1000).toFixed(2)} s`
  const m = Math.floor(ms / 60_000)
  const s = ((ms % 60_000) / 1000).toFixed(1)
  return `${m}m ${s}s`
}

function parseStepCount(steps) {
  if (!steps)
    return 0
  try {
    const arr = typeof steps === 'string' ? JSON.parse(steps) : steps
    return Array.isArray(arr) ? arr.length : 0
  }
  catch {
    return 0
  }
}

const cronPresets = [
  { label: '每分钟', spec: '*/1 * * * *' },
  { label: '每 5 分', spec: '*/5 * * * *' },
  { label: '每 10 分', spec: '*/10 * * * *' },
  { label: '每小时', spec: '0 * * * *' },
  { label: '每天 0 点', spec: '0 0 * * *' },
  { label: '每天 9 点', spec: '0 9 * * *' },
  { label: '每周一 9 点', spec: '0 9 * * 1' },
  { label: '每月 1 号', spec: '0 0 1 * *' },
]

const funcOptions = ref([])
api.getFuncs().then(({ data = [] }) => {
  funcOptions.value = data.map(name => ({ label: name, value: name }))
})

const stats = ref({})
async function loadStats() {
  try {
    const { data } = await api.stats()
    stats.value = data || {}
  }
  catch (e) {
    console.error(e)
  }
}
onMounted(() => {
  loadStats()
})

const { modalRef, modalFormRef, modalAction, modalForm, handleAdd, handleDelete, handleEdit }
  = useCrud({
    name: '任务',
    initForm: {
      enabled: true,
      type: 'HTTP',
      httpMethod: 'GET',
      timeout: 60,
    },
    doCreate: api.create,
    doDelete: api.delete,
    doUpdate: api.update,
    refresh: (_, keepCurrentPage) => {
      $table.value?.handleSearch(keepCurrentPage)
      loadStats()
    },
  })

// 当 modalForm.spec 变化时刷新下一次执行时间预览
const nextRuns = ref([])
const specError = ref('')
let previewTimer = null
function onSpecChange(spec) {
  clearTimeout(previewTimer)
  if (!spec) {
    nextRuns.value = []
    specError.value = ''
    return
  }
  previewTimer = setTimeout(async () => {
    try {
      const { data } = await api.previewNext(spec, 3)
      nextRuns.value = data || []
      specError.value = ''
    }
    catch (e) {
      nextRuns.value = []
      specError.value = e?.response?.data?.message || 'cron 表达式不合法'
    }
  }, 350)
}
function applyPreset(spec) {
  modalForm.value.spec = spec
  onSpecChange(spec)
}
watch(() => modalForm.value.spec, (v) => {
  if (v)
    onSpecChange(v)
})

onMounted(() => $table.value?.handleSearch())

const typeTagMap = { HTTP: 'info', SHELL: 'warning', FUNC: 'success' }
const statusTagMap = { SUCCESS: 'success', FAILED: 'error', TIMEOUT: 'warning' }

const columns = [
  { type: 'selection', width: 40, fixed: 'left' },
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '任务名',
    key: 'name',
    width: 180,
    ellipsis: { tooltip: true },
    render: row =>
      h('a', { class: 'text-primary cursor-pointer', onClick: () => openDetail(row) }, row.name),
  },
  {
    title: '类型',
    key: 'type',
    width: 90,
    render: ({ type }) =>
      h(NTag, { type: typeTagMap[type] || 'default', bordered: false, size: 'small' }, { default: () => type }),
  },
  { title: 'Cron', key: 'spec', width: 130, ellipsis: { tooltip: true } },
  { title: '执行内容', key: 'command', width: 220, ellipsis: { tooltip: true } },
  {
    title: '上次执行',
    key: 'lastRunAt',
    width: 160,
    render: row => (row.lastRunAt ? formatDateTime(row.lastRunAt) : '-'),
  },
  {
    title: '下次执行',
    key: 'nextRunAt',
    width: 160,
    render: row => (row.nextRunAt ? formatDateTime(row.nextRunAt) : '-'),
  },
  {
    title: '状态',
    key: 'enabled',
    width: 110,
    fixed: 'right',
    render: row =>
      h(
        NSwitch,
        {
          size: 'small',
          rubberBand: false,
          value: row.enabled,
          loading: !!row.enableLoading,
          onUpdateValue: () => handleToggle(row),
        },
        { checked: () => '启用', unchecked: () => '停用' },
      ),
  },
  {
    title: '操作',
    key: 'actions',
    width: 340,
    align: 'right',
    fixed: 'right',
    render(row) {
      return [
        withPermission(
          h(
            NButton,
            { size: 'small', type: 'primary', secondary: true, onClick: () => handleRun(row) },
            { default: () => '执行', icon: () => h('i', { class: 'i-fe:play text-14' }) },
          ),
          'RunTask',
        ),
        h(
          NButton,
          { size: 'small', type: 'primary', secondary: true, style: 'margin-left: 8px;', onClick: () => openLogs(row) },
          { default: () => '日志', icon: () => h('i', { class: 'i-fe:file-text text-14' }) },
        ),
        withPermission(
          h(
            NButton,
            { size: 'small', type: 'primary', style: 'margin-left: 8px;', onClick: () => handleEdit(row) },
            { default: () => '编辑', icon: () => h('i', { class: 'i-material-symbols:edit-outline text-14' }) },
          ),
          'EditTask',
        ),
        withPermission(
          h(
            NButton,
            { size: 'small', type: 'error', style: 'margin-left: 8px;', onClick: () => handleDelete(row.id) },
            { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }) },
          ),
          'DeleteTask',
        ),
      ]
    },
  },
]

function onChecked(ids) {
  checkedIds.value = ids
}

async function handleBulkDelete() {
  if (!checkedIds.value.length)
    return
  $dialog.warning({
    title: '批量删除',
    content: `确定删除选中的 ${checkedIds.value.length} 个任务？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await api.bulkDelete(checkedIds.value)
      $message.success('删除成功')
      checkedIds.value = []
      $table.value?.handleSearch(true)
      loadStats()
    },
  })
}

async function handleBulkToggle(enabled) {
  if (!checkedIds.value.length)
    return
  await api.bulkToggle(checkedIds.value, enabled)
  $message.success(enabled ? '已批量启用' : '已批量停用')
  checkedIds.value = []
  $table.value?.handleSearch(true)
  loadStats()
}

async function handleToggle(row) {
  row.enableLoading = true
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
    row.enableLoading = false
  }
}

async function handleRun(row) {
  try {
    await api.run(row.id)
    $message.success('已触发执行，可在日志中查看结果')
    setTimeout(loadStats, 800)
  }
  catch (err) {
    console.error(err)
  }
}

// ====== 详情抽屉 ======
const detailVisible = ref(false)
function openDetail(row) {
  currentTask.value = row
  detailVisible.value = true
}

// ====== 日志抽屉 ======
const logDrawerVisible = ref(false)
const currentTask = ref({})
const logRows = ref([])
const logLoading = ref(false)
const logStatus = ref('')
const logTrigger = ref(null)
const autoRefresh = ref(false)
let autoRefreshTimer = null
const logSummary = ref({ success: 0, failed: 0, timeout: 0, avg: 0 })
const logPagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: page => loadLogs(page),
  onUpdatePageSize: (size) => {
    logPagination.pageSize = size
    loadLogs(1)
  },
})

function openLogs(row) {
  currentTask.value = row
  logStatus.value = ''
  logTrigger.value = null
  logDrawerVisible.value = true
  loadLogs(1)
}

watch(logDrawerVisible, (v) => {
  if (!v) {
    stopAutoRefresh()
    autoRefresh.value = false
  }
})

function onAutoRefreshChange(v) {
  if (v) {
    autoRefreshTimer = setInterval(() => loadLogs(logPagination.page), 3000)
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

async function loadLogs(page) {
  logPagination.page = page
  logLoading.value = true
  try {
    const { data } = await api.getLogs({
      pageNo: page,
      pageSize: logPagination.pageSize,
      taskId: currentTask.value.id,
      status: logStatus.value || undefined,
      trigger: logTrigger.value || undefined,
    })
    logRows.value = data?.pageData || []
    logPagination.itemCount = data?.total || 0
    // 当前页汇总
    let success = 0; let failed = 0; let timeout = 0; let sum = 0; let n = 0
    for (const r of logRows.value) {
      if (r.status === 'SUCCESS') success++
      else if (r.status === 'FAILED') failed++
      else if (r.status === 'TIMEOUT') timeout++
      if (typeof r.duration === 'number') { sum += r.duration; n++ }
    }
    logSummary.value = { success, failed, timeout, avg: n ? Math.round(sum / n) : 0 }
  }
  finally {
    logLoading.value = false
  }
}

// ====== 输出详情 ======
const outputVisible = ref(false)
const currentLog = ref({})
const stepTypeMap = { INFO: 'info', WARN: 'warning', ERROR: 'error' }
const stepIconMap = {
  INFO: 'i-fe:info',
  WARN: 'i-fe:alert-triangle',
  ERROR: 'i-fe:x-circle',
}
const stepLevelFilters = [
  { label: '全部', value: '' },
  { label: 'INFO', value: 'INFO' },
  { label: 'WARN', value: 'WARN' },
  { label: 'ERROR', value: 'ERROR' },
]
const stepFilter = ref('')

watch(outputVisible, (v) => {
  if (v) stepFilter.value = ''
})

const parsedSteps = computed(() => {
  const raw = currentLog.value.steps
  if (!raw)
    return []
  let arr
  try {
    arr = typeof raw === 'string' ? JSON.parse(raw) : raw
  }
  catch {
    return []
  }
  if (!Array.isArray(arr) || !arr.length)
    return []
  const base = new Date(arr[0].time).getTime()
  return arr.map((s, i) => {
    const t = new Date(s.time).getTime()
    return {
      ...s,
      idx: i + 1,
      delta: Number.isFinite(t - base) ? t - base : 0,
    }
  })
})

const filteredSteps = computed(() =>
  stepFilter.value
    ? parsedSteps.value.filter(s => s.level === stepFilter.value)
    : parsedSteps.value,
)

const stepCounts = computed(() => {
  const c = { INFO: 0, WARN: 0, ERROR: 0 }
  for (const s of parsedSteps.value)
    if (c[s.level] != null) c[s.level]++
  return c
})

const stepElapsed = computed(() => {
  const arr = parsedSteps.value
  if (arr.length < 2)
    return ''
  return formatDuration(arr[arr.length - 1].delta)
})

function formatStepTime(t) {
  const d = new Date(t)
  if (Number.isNaN(d.getTime()))
    return '-'
  const pad = (n, w = 2) => String(n).padStart(w, '0')
  return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}.${pad(d.getMilliseconds(), 3)}`
}

function showOutput(row) {
  currentLog.value = row
  outputVisible.value = true
}

const logColumns = [
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: ({ status }) =>
      h(NTag, { type: statusTagMap[status] || 'default', bordered: false, size: 'small' }, { default: () => status }),
  },
  {
    title: '触发',
    key: 'trigger',
    width: 70,
    render: ({ trigger }) =>
      h(
        NTag,
        { type: trigger === 'MANUAL' ? 'info' : 'default', bordered: false, size: 'small' },
        { default: () => (trigger === 'MANUAL' ? '手动' : '定时') },
      ),
  },
  { title: '开始时间', key: 'startTime', width: 160, render: row => formatDateTime(row.startTime) },
  {
    title: '耗时',
    key: 'duration',
    width: 90,
    render: (row) => {
      const ms = row.duration
      const color = ms == null ? '#999' : ms > 10_000 ? '#d03050' : ms > 3000 ? '#f0a020' : '#18a058'
      return h('span', { style: { color } }, formatDuration(ms))
    },
  },
  {
    title: '步骤',
    key: 'steps',
    width: 60,
    render: row => parseStepCount(row.steps),
  },
  {
    title: '输出预览',
    key: 'output',
    minWidth: 200,
    ellipsis: { tooltip: true },
    render: row => (row.output || row.error || '').split('\n')[0] || '-',
  },
  {
    title: '操作',
    key: 'actions',
    width: 90,
    align: 'right',
    fixed: 'right',
    render: row =>
      h(NButton, { size: 'small', type: 'primary', text: true, onClick: () => showOutput(row) }, { default: () => '详情' }),
  },
]
</script>

<style scoped>
.text-success { color: #18a058; }
.text-error { color: #d03050; }
.text-primary { color: #2080f0; }
.text-12 { font-size: 12px; }
.text-muted { color: #9ca3af; }

.task-form-scroll {
  max-height: calc(100vh - 240px);
  overflow-y: auto;
  padding-right: 6px;
  margin-right: -6px;
}
.task-form-scroll::-webkit-scrollbar {
  width: 6px;
}
.task-form-scroll::-webkit-scrollbar-thumb {
  background: #d4d4d8;
  border-radius: 3px;
}

.step-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  margin-bottom: 12px;
  background: var(--n-action-color, #fafafc);
  border: 1px solid var(--n-border-color, #efeff5);
  border-radius: 6px;
}

.step-scroll {
  max-height: min(460px, calc(100vh - 380px));
  min-height: 220px;
  overflow-y: auto;
  padding: 4px 8px 4px 4px;
}

.step-scroll::-webkit-scrollbar {
  width: 6px;
}
.step-scroll::-webkit-scrollbar-thumb {
  background: #d4d4d8;
  border-radius: 3px;
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

.step-icon {
  display: inline-block;
  width: 14px;
  height: 14px;
}

.step-header {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  line-height: 1.4;
}
.step-idx {
  font-family: ui-monospace, Menlo, Consolas, monospace;
  color: #9ca3af;
  font-size: 12px;
  min-width: 28px;
}
.step-time {
  font-family: ui-monospace, Menlo, Consolas, monospace;
  font-size: 12px;
  color: #6b7280;
}
.step-delta {
  font-size: 12px;
  color: #2080f0;
  background: rgba(32, 128, 240, 0.08);
  padding: 0 6px;
  border-radius: 8px;
}
.step-msg {
  margin-top: 4px;
  padding: 6px 10px;
  font-family: ui-monospace, Menlo, Consolas, monospace;
  font-size: 12px;
  line-height: 1.55;
  color: #374151;
  background: #f8fafc;
  border-left: 3px solid #e5e7eb;
  border-radius: 0 4px 4px 0;
  word-break: break-all;
  white-space: pre-wrap;
}
</style>
