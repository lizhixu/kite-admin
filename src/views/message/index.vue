<template>
  <CommonPage>
    <template #action>
      <NSpace>
        <NButton v-permission="'DeleteMessage'" :disabled="!checkedIds.length" type="error" secondary @click="handleBulkDelete">
          <i class="i-material-symbols:delete-outline mr-4" />批量删除
        </NButton>
        <NButton v-permission="'SendMessage'" type="primary" @click="handleSend()">
          <i class="i-material-symbols:add mr-4 text-18" />发送消息
        </NButton>
      </NSpace>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1000"
      :columns="columns"
      :get-data="api.read"
      @on-checked="onChecked"
    >
      <MeQueryItem label="标题" :label-width="50">
        <n-input
          v-model:value="queryItems.title"
          type="text"
          placeholder="请输入标题"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="类型" :label-width="50">
        <n-select
          v-model:value="queryItems.type"
          clearable
          :options="typeOptions"
          style="width: 140px"
        />
      </MeQueryItem>
    </MeCrud>

    <!-- 发送消息弹窗 -->
    <MeModal ref="modalRef" width="720px" title="发送消息">
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
        :rules="formRules"
      >
        <n-form-item label="标题" path="title">
          <n-input v-model:value="modalForm.title" placeholder="请输入消息标题" />
        </n-form-item>
        <n-form-item label="类型" path="type">
          <n-select v-model:value="modalForm.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item label="目标" path="targetType">
          <n-radio-group v-model:value="modalForm.targetType">
            <n-radio value="ALL">全员广播</n-radio>
            <n-radio value="ROLE">按角色</n-radio>
            <n-radio value="USER">指定用户</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item v-if="modalForm.targetType === 'ROLE'" label="选择角色" path="roleIds">
          <n-select
            v-model:value="modalForm.roleIds"
            multiple
            filterable
            :options="roleOptions"
            placeholder="选择角色分组"
          />
        </n-form-item>
        <n-form-item v-if="modalForm.targetType === 'USER'" label="选择用户" path="userIds">
          <n-select
            v-model:value="modalForm.userIds"
            multiple
            filterable
            remote
            :loading="userLoading"
            :options="userOptions"
            placeholder="搜索并选择用户"
            :remote-method="searchUsers"
            @focus="searchUsers('')"
          />
        </n-form-item>
        <n-form-item label="内容" path="content">
          <n-tabs type="segment" animated>
            <n-tab-pane name="edit" tab="编辑">
              <n-input
                v-model:value="modalForm.content"
                type="textarea"
                placeholder="支持 Markdown 格式"
                :autosize="{ minRows: 10, maxRows: 20 }"
              />
            </n-tab-pane>
            <n-tab-pane name="preview" tab="预览">
              <div class="modal-preview" v-html="previewHTML" />
            </n-tab-pane>
          </n-tabs>
        </n-form-item>
      </n-form>
    </MeModal>

    <!-- 内容详情抽屉 -->
    <n-drawer v-model:show="detailVisible" width="640">
      <n-drawer-content>
        <template #header>
          <span class="text-16 font-bold">{{ detailRow?.title || '消息详情' }}</span>
        </template>
        <MessageDetail :message="detailRow" show-target />
      </n-drawer-content>
    </n-drawer>
  </CommonPage>
</template>

<script setup>
import {
  NButton,
  NDrawer,
  NDrawerContent,
  NRadio,
  NRadioGroup,
  NSpace,
  NTabPane,
  NTabs,
  NTag,
} from 'naive-ui'
import 'highlight.js/styles/github.css'
import { MeCrud, MeModal, MeQueryItem, MessageDetail } from '@/components'
import { useCrud } from '@/composables'
import { withPermission } from '@/directives'
import { renderMarkdown } from '@/utils/markdown'
import api from './api'
import userApi from '@/views/pms/user/api'

defineOptions({ name: 'MessageList' })

const $table = ref(null)
const queryItems = ref({})
const checkedIds = ref([])

const typeOptions = [
  { label: '系统消息', value: 'SYSTEM' },
  { label: '通知公告', value: 'NOTICE' },
  { label: '公告', value: 'ANNOUNCEMENT' },
]

const formRules = {
  title: [{ required: true, message: '请输入消息标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入消息内容', trigger: 'blur' }],
}

// ====== 角色选择 ======
const roleOptions = ref([])
async function fetchRoles() {
  try {
    const { data } = await userApi.getAllRoles()
    roleOptions.value = (data || []).map(r => ({ label: r.name, value: r.id }))
  } catch { /* ignore */ }
}
fetchRoles()

// ====== 用户选择 ======
const userLoading = ref(false)
const userOptions = ref([])
let searchTimer = null

async function searchUsers(query) {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(async () => {
    userLoading.value = true
    try {
      const { data } = await userApi.read({ username: query, pageNo: 1, pageSize: 50 })
      const list = Array.isArray(data) ? data : data?.pageData || []
      userOptions.value = list.map(u => ({
        label: u.profile?.nickName || u.username,
        value: u.id,
      }))
    } catch { /* ignore */ }
    userLoading.value = false
  }, 300)
}

// ====== Markdown preview ======
const previewHTML = computed(() => renderMarkdown(modalForm.value.content || ''))

// ====== CRUD ======
const {
  modalRef,
  modalFormRef,
  modalForm,
  handleAdd: _handleAdd,
} = useCrud({
  name: '消息',
  initForm: {
    title: '',
    content: '',
    type: 'SYSTEM',
    targetType: 'ALL',
    userIds: [],
    roleIds: [],
  },
  doCreate: api.create,
  doDelete: api.delete,
  refresh: (_, keepCurrentPage) => {
    $table.value?.handleSearch(keepCurrentPage)
  },
})

function handleSend() {
  _handleAdd()
}

function handleDelete(rowId) {
  $dialog.warning({
    title: '确认删除',
    content: '删除后不可恢复，确定要删除该消息吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await api.delete(rowId)
        $message.success('删除成功')
        $table.value?.handleSearch(true)
      } catch { /* ignore */ }
    },
  })
}

function onChecked(keys) {
  checkedIds.value = keys
}

function handleBulkDelete() {
  if (!checkedIds.value.length) return
  $dialog.warning({
    title: '批量删除',
    content: `确定要删除选中的 ${checkedIds.value.length} 条消息吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await api.bulkDelete(checkedIds.value)
        $message.success('删除成功')
        checkedIds.value = []
        $table.value?.handleSearch(true)
      } catch { /* ignore */ }
    },
  })
}

onMounted(() => $table.value?.handleSearch())

// ====== 详情抽屉 ======
const detailVisible = ref(false)
const detailRow = ref(null)

function showDetail(row) {
  detailRow.value = row
  detailVisible.value = true
}

// ====== 表格列 ======
const columns = [
  { type: 'selection', width: 40 },
  { title: 'ID', key: 'id', width: 60 },
  {
    title: '标题',
    key: 'title',
    minWidth: 200,
    ellipsis: { tooltip: true },
  },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render: row =>
      h(
        NTag,
        { type: typeTagType(row.type), bordered: false, size: 'small' },
        { default: () => typeLabel(row.type) },
      ),
  },
  {
    title: '目标',
    key: 'targetType',
    width: 100,
    render(row) {
      const map = { ALL: { label: '全员', type: 'default' }, ROLE: { label: '按角色', type: 'warning' }, USER: { label: '指定用户', type: 'info' } }
      const cfg = map[row.targetType] || { label: row.targetType, type: 'default' }
      return h(NTag, { type: cfg.type, bordered: false, size: 'small' }, { default: () => cfg.label })
    },
  },
  { title: '发送人', key: 'senderName', width: 100 },
  {
    title: '发送时间',
    key: 'createTime',
    width: 170,
    render: row => formatTime(row.createTime),
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
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
            onClick: () => showDetail(row),
          },
          { default: () => '查看' },
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
          'DeleteMessage',
        ),
      ]
    },
  },
]

// ====== Helpers ======
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
.message-detail :deep(.n-descriptions) {
  margin-bottom: 12px;
}

.modal-preview {
  min-height: 240px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 4px;
  padding: 16px;
  background: rgba(0, 0, 0, 0.02);
  overflow-y: auto;
  max-height: 460px;
}
:root.dark .modal-preview {
  border-color: rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.04);
}
.modal-preview :deep(h1),
.modal-preview :deep(h2),
.modal-preview :deep(h3),
.modal-preview :deep(h4) {
  margin: 0.8em 0 0.4em;
  font-weight: 600;
  line-height: 1.4;
}
.modal-preview :deep(h1) {
  font-size: 1.4em;
}
.modal-preview :deep(h2) {
  font-size: 1.2em;
}
.modal-preview :deep(h3) {
  font-size: 1.1em;
}
.modal-preview :deep(p) {
  margin: 0.5em 0;
  line-height: 1.7;
}
.modal-preview :deep(ul),
.modal-preview :deep(ol) {
  padding-left: 1.5em;
  margin: 0.5em 0;
}
.modal-preview :deep(blockquote) {
  margin: 0.6em 0;
  padding: 0.4em 0.8em;
  border-left: 3px solid #18a058;
  background: rgba(24, 160, 88, 0.06);
  color: rgba(0, 0, 0, 0.65);
}
:root.dark .modal-preview :deep(blockquote) {
  background: rgba(24, 160, 88, 0.1);
  color: rgba(255, 255, 255, 0.65);
}
.modal-preview :deep(code) {
  background: rgba(0, 0, 0, 0.06);
  padding: 2px 5px;
  border-radius: 3px;
  font-size: 0.9em;
  font-family: 'Consolas', monospace;
}
:root.dark .modal-preview :deep(code) {
  background: rgba(255, 255, 255, 0.08);
}
.modal-preview :deep(pre) {
  background: rgba(0, 0, 0, 0.04);
  border-radius: 4px;
  padding: 12px;
  overflow-x: auto;
  margin: 0.6em 0;
}
:root.dark .modal-preview :deep(pre) {
  background: rgba(255, 255, 255, 0.06);
}
.modal-preview :deep(pre code) {
  background: none;
  padding: 0;
}
.modal-preview :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 0.6em 0;
}
.modal-preview :deep(th),
.modal-preview :deep(td) {
  border: 1px solid rgba(0, 0, 0, 0.08);
  padding: 6px 10px;
  text-align: left;
}
:root.dark .modal-preview :deep(th),
:root.dark .modal-preview :deep(td) {
  border-color: rgba(255, 255, 255, 0.12);
}
.modal-preview :deep(th) {
  background: rgba(0, 0, 0, 0.03);
  font-weight: 600;
}
:root.dark .modal-preview :deep(th) {
  background: rgba(255, 255, 255, 0.05);
}
</style>