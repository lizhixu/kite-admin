<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2024/04/01 15:52:40
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <CommonPage>
    <template #action>
      <NButton v-permission="'AddRole'" type="primary" @click="handleOpenAddRole()">
        <i class="i-material-symbols:add mr-4 text-18" />
        新增角色
      </NButton>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1200"
      :columns="columns"
      :get-data="api.read"
    >
      <MeQueryItem label="角色名" :label-width="50">
        <n-input v-model:value="queryItems.name" type="text" placeholder="请输入角色名" clearable />
      </MeQueryItem>
      <MeQueryItem label="状态" :label-width="50">
        <n-select
          v-model:value="queryItems.enable"
          clearable
          :options="[
            { label: '启用', value: 1 },
            { label: '停用', value: 0 },
          ]"
        />
      </MeQueryItem>
    </MeCrud>
    <MeModal
      ref="modalRef"
      width="min(1080px, calc(100vw - 32px))"
      :content-style="{ maxHeight: 'calc(100vh - 8px)' }"
    >
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
        class="role-form-scroll"
        :class="{ 'is-permission-pinned': isPermissionPinned }"
      >
        <n-form-item
          label="角色名"
          path="name"
          :rule="{
            required: true,
            message: '请输入角色名',
            trigger: ['input', 'blur'],
          }"
        >
          <n-input v-model:value="modalForm.name" />
        </n-form-item>
        <n-form-item
          label="角色编码"
          path="code"
          :rule="{
            required: true,
            message: '请输入角色编码',
            trigger: ['input', 'blur'],
          }"
        >
          <n-input v-model:value="modalForm.code" :disabled="modalAction !== 'add'" />
        </n-form-item>
        <n-form-item ref="permissionFormItemRef" label="权限" path="permissionIds" class="permission-form-item">
          <div class="permission-panel w-full">
            <div class="permission-panel__header">
              <div>
                <div class="permission-panel__title">
                  权限配置
                </div>
                <div class="permission-panel__desc">
                  勾选菜单控制左侧访问入口，勾选按钮控制页面内操作权限
                </div>
              </div>
              <div class="permission-summary">
                <span>已选 <b>{{ selectedPermissionCount }}</b></span>
                <span>菜单 <b>{{ selectedMenuCount }}</b></span>
                <span>按钮 <b>{{ selectedButtonCount }}</b></span>
              </div>
              <NButton
                size="tiny"
                secondary
                type="primary"
                class="permission-pin-btn"
                :title="isPermissionPinned ? '取消固定' : '固定'"
                @click="handleTogglePermissionPinned"
              >
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  aria-hidden="true"
                >
                  <path v-if="isPermissionPinned" d="M5 5l14 14" />
                  <template v-if="isPermissionPinned">
                    <path d="M12 17v5" />
                    <path d="M9 22h6" />
                    <path d="M14 4l4 4-2.5 2.5" />
                    <path d="M10.5 13.5 8 16l-4-4 2.5-2.5" />
                  </template>
                  <template v-else>
                    <path d="M12 17v5" />
                    <path d="M9 22h6" />
                    <path d="M17 4l3 3-5 5v3l-2 2-6-6 2-2h3z" />
                  </template>
                </svg>
              </NButton>
            </div>

            <div class="permission-toolbar">
              <n-input
                v-model:value="permissionSearch"
                placeholder="搜索权限名称、编码或路径"
                clearable
                class="permission-search"
              >
                <template #prefix>
                  <i class="i-fe:search text-14 text-#999" />
                </template>
              </n-input>

              <div class="permission-copy">
                <n-select
                  v-model:value="copyRoleId"
                  :options="availableCopyRoleOptions"
                  :loading="copyRolesLoading"
                  clearable
                  filterable
                  placeholder="从已有角色复制权限"
                  class="permission-copy__select"
                />
                <NButton
                  secondary
                  type="primary"
                  :disabled="!copyRoleId"
                  @click="handleCopyRolePermissions"
                >
                  复制
                </NButton>
              </div>
            </div>

            <div class="permission-actions">
              <n-space :size="8">
                <NButton size="small" secondary @click="handleCheckAllPermissions">
                  全选
                </NButton>
                <NButton size="small" secondary @click="handleClearPermissions">
                  清空
                </NButton>
                <NButton size="small" secondary @click="handleExpandAllPermissions">
                  展开全部
                </NButton>
                <NButton size="small" secondary @click="handleCollapseAllPermissions">
                  收起全部
                </NButton>
                <NButton
                  size="small"
                  secondary
                  :disabled="!matchedPermissionIds.length"
                  @click="handleCheckMatchedPermissions"
                >
                  勾选搜索结果
                </NButton>
              </n-space>
              <span v-if="permissionSearch" class="permission-match-tip">
                匹配 {{ matchedPermissionIds.length }} 项
              </span>
            </div>

            <n-tree
              key-field="id"
              label-field="name"
              :selectable="false"
              :data="permissionTree"
              :pattern="permissionSearch"
              :filter="filterPermissionNode"
              :show-irrelevant-nodes="false"
              :checked-keys="modalForm.permissionIds"
              :on-update:checked-keys="handlePermissionCheckedKeysChange"
              :expanded-keys="permissionExpandedKeys"
              :on-update:expanded-keys="keys => (permissionExpandedKeys = keys)"
              :render-label="renderPermissionLabel"
              checkable check-on-click block-line
              class="cus-scroll role-permission-tree w-full"
            />
          </div>
        </n-form-item>
        <n-form-item label="状态" path="enable">
          <NSwitch v-model:value="modalForm.enable">
            <template #checked>
              启用
            </template>
            <template #unchecked>
              停用
            </template>
          </NSwitch>
        </n-form-item>
      </n-form>
    </MeModal>
  </CommonPage>
</template>

<script setup>
import { NButton, NSwitch } from 'naive-ui'
import { MeCrud, MeModal, MeQueryItem } from '@/components'
import { useCrud } from '@/composables'
import { withPermission } from '@/directives'
import api from './api'

defineOptions({ name: 'RoleMgt' })

const router = useRouter()

const $table = ref(null)
/** QueryBar筛选参数（可选） */
const queryItems = ref({})

onMounted(() => {
  $table.value?.handleSearch()
})

const { modalRef, modalFormRef, modalAction, modalForm, handleAdd, handleDelete, handleEdit }
  = useCrud({
    name: '角色',
    doCreate: api.create,
    doDelete: api.delete,
    doUpdate: api.update,
    initForm: { enable: true, permissionIds: [] },
    refresh: (_, keepCurrentPage) => $table.value?.handleSearch(keepCurrentPage),
  })

const columns = [
  { title: '角色名', key: 'name' },
  { title: '角色编码', key: 'code' },
  {
    title: '状态',
    key: 'enable',
    render: row =>
      h(
        NSwitch,
        {
          size: 'small',
          rubberBand: false,
          value: row.enable,
          loading: !!row.enableLoading,
          disabled: row.code === 'SUPER_ADMIN',
          onUpdateValue: () => handleEnable(row),
        },
        {
          checked: () => '启用',
          unchecked: () => '停用',
        },
      ),
  },
  {
    title: '操作',
    key: 'actions',
    width: 320,
    align: 'right',
    fixed: 'right',
    render(row) {
      return [
        withPermission(
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              secondary: true,
              onClick: () =>
                router.push({ path: `/pms/role/user/${row.id}`, query: { roleName: row.name } }),
            },
            {
              default: () => '分配用户',
              icon: () => h('i', { class: 'i-fe:user-plus text-14' }),
            },
          ),
          'AssignPermission',
        ),
        withPermission(
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              style: 'margin-left: 12px;',
              disabled: row.code === 'SUPER_ADMIN',
              onClick: () => handleOpenEditRole(row),
            },
            {
              default: () => '编辑',
              icon: () => h('i', { class: 'i-material-symbols:edit-outline text-14' }),
            },
          ),
          'EditRole',
        ),
        withPermission(
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              style: 'margin-left: 12px;',
              disabled: row.code === 'SUPER_ADMIN',
              onClick: () => handleDelete(row.id),
            },
            {
              default: () => '删除',
              icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
            },
          ),
          'DeleteRole',
        ),
      ]
    },
  },
]

async function handleEnable(row) {
  row.enableLoading = true
  try {
    await api.update({ id: row.id, enable: !row.enable })
    row.enableLoading = false
    $message.success('操作成功')
    $table.value?.handleSearch()
  }
  catch (error) {
    console.error(error)
    row.enableLoading = false
  }
}

const permissionTree = ref([])
const permissionSearch = ref('')
const permissionExpandedKeys = ref([])
const copyRoleId = ref(null)
const copyRoleOptions = ref([])
const copyRolesLoading = ref(false)
const isPermissionPinned = ref(false)
const permissionFormItemRef = ref(null)

const flatPermissions = computed(() => flattenPermissions(permissionTree.value))
const allPermissionIds = computed(() => flatPermissions.value.map(item => item.id))
const selectedPermissionIds = computed(() => new Set(modalForm.value.permissionIds || []))
const selectedPermissionCount = computed(() => selectedPermissionIds.value.size)
const selectedMenuCount = computed(() => flatPermissions.value.filter(item => selectedPermissionIds.value.has(item.id) && item.type === 'MENU').length)
const selectedButtonCount = computed(() => flatPermissions.value.filter(item => selectedPermissionIds.value.has(item.id) && item.type === 'BUTTON').length)
const matchedPermissionIds = computed(() => {
  const keyword = permissionSearch.value?.trim().toLowerCase()
  if (!keyword)
    return []
  return flatPermissions.value
    .filter(item => isPermissionMatched(item, keyword))
    .map(item => item.id)
})
const availableCopyRoleOptions = computed(() => {
  return copyRoleOptions.value.filter(role => role.value !== modalForm.value.id)
})

api.getAllPermissionTree().then(({ data = [] }) => {
  permissionTree.value = data
  permissionExpandedKeys.value = flattenPermissions(data).map(item => item.id)
})

watch(permissionSearch, (value) => {
  if (value)
    handleExpandAllPermissions()
})

function flattenPermissions(tree = []) {
  const result = []
  const walk = (nodes = []) => {
    nodes.forEach((node) => {
      result.push(node)
      if (node.children?.length)
        walk(node.children)
    })
  }
  walk(tree)
  return result
}

function isPermissionMatched(item, keyword) {
  return [item.name, item.code, item.path, item.type]
    .filter(Boolean)
    .some(value => String(value).toLowerCase().includes(keyword))
}

function filterPermissionNode(pattern, node) {
  const keyword = pattern?.trim().toLowerCase()
  if (!keyword)
    return true
  return isPermissionMatched(node, keyword)
}

function renderPermissionLabel({ option }) {
  const isButton = option.type === 'BUTTON'
  return h('div', { class: 'permission-node-label' }, [
    h('span', { class: 'permission-node-label__name' }, option.name),
    h('span', { class: isButton ? 'permission-node-label__tag is-button' : 'permission-node-label__tag is-menu' }, isButton ? '按钮' : '菜单'),
    option.code ? h('span', { class: 'permission-node-label__code' }, option.code) : null,
    option.path ? h('span', { class: 'permission-node-label__path' }, option.path) : null,
  ])
}

function normalizePageRows(data) {
  if (Array.isArray(data))
    return data
  return data?.pageData || data?.records || data?.list || data?.rows || []
}

function resetPermissionAssignState() {
  permissionSearch.value = ''
  copyRoleId.value = null
  isPermissionPinned.value = false
  permissionExpandedKeys.value = [...allPermissionIds.value]
}

function getElement(target) {
  return target?.$el || target
}

async function handleTogglePermissionPinned() {
  if (isPermissionPinned.value) {
    isPermissionPinned.value = false
    return
  }

  const formEl = getElement(modalFormRef.value)
  const permissionItemEl = getElement(permissionFormItemRef.value)
  if (formEl && permissionItemEl)
    formEl.scrollTop = Math.max(permissionItemEl.offsetTop - 85, 0)
  await nextTick()
  isPermissionPinned.value = true
}

async function ensureCopyRoleOptionsLoaded() {
  if (copyRoleOptions.value.length)
    return
  copyRolesLoading.value = true
  try {
    const { data } = await api.read({ pageNo: 1, pageSize: 100 })
    copyRoleOptions.value = normalizePageRows(data).map(role => ({
      ...role,
      label: role.name,
      value: role.id,
    }))
  }
  finally {
    copyRolesLoading.value = false
  }
}

function handleOpenAddRole() {
  resetPermissionAssignState()
  ensureCopyRoleOptionsLoaded()
  handleAdd()
}

function handleOpenEditRole(row) {
  resetPermissionAssignState()
  ensureCopyRoleOptionsLoaded()
  handleEdit({
    ...row,
    permissionIds: [...(row.permissionIds || [])],
  })
}

function getPermissionById(id) {
  return flatPermissions.value.find(item => item.id === id)
}

function getDescendantPermissionIds(node) {
  return flattenPermissions(node?.children || []).map(item => item.id)
}

function handlePermissionCheckedKeysChange(keys) {
  const oldKeys = modalForm.value.permissionIds || []
  const oldSet = new Set(oldKeys)
  const newSet = new Set(keys)
  const addedKey = keys.find(key => !oldSet.has(key))
  const removedKey = oldKeys.find(key => !newSet.has(key))

  if (addedKey !== undefined) {
    const node = getPermissionById(addedKey)
    const descendantIds = getDescendantPermissionIds(node)
    modalForm.value.permissionIds = Array.from(new Set([...keys, ...descendantIds]))
    return
  }

  if (removedKey !== undefined) {
    const node = getPermissionById(removedKey)
    const descendantIds = getDescendantPermissionIds(node)
    if (descendantIds.length) {
      const removedSet = new Set([removedKey, ...descendantIds])
      modalForm.value.permissionIds = oldKeys.filter(key => !removedSet.has(key))
      return
    }
  }

  modalForm.value.permissionIds = keys
}

function handleCheckAllPermissions() {
  modalForm.value.permissionIds = [...allPermissionIds.value]
}

function handleClearPermissions() {
  modalForm.value.permissionIds = []
}

function handleExpandAllPermissions() {
  permissionExpandedKeys.value = [...allPermissionIds.value]
}

function handleCollapseAllPermissions() {
  permissionExpandedKeys.value = []
}

function handleCheckMatchedPermissions() {
  modalForm.value.permissionIds = Array.from(new Set([
    ...(modalForm.value.permissionIds || []),
    ...matchedPermissionIds.value,
  ]))
}

function handleCopyRolePermissions() {
  const sourceRole = copyRoleOptions.value.find(role => role.value === copyRoleId.value)
  if (!sourceRole)
    return

  const applyCopy = () => {
    modalForm.value.permissionIds = [...(sourceRole.permissionIds || [])]
    $message.success(`已复制【${sourceRole.name}】的权限`)
  }

  if (modalForm.value.permissionIds?.length) {
    $dialog.warning({
      title: '确认覆盖权限？',
      content: '复制后会覆盖当前已选择的权限。',
      positiveText: '确认',
      negativeText: '取消',
      onPositiveClick: applyCopy,
    })
    return
  }

  applyCopy()
}
</script>

<style scoped>
.role-form-scroll {
  --role-modal-body-height: min(700px, calc(100vh - 180px));
  --permission-pin-top-gap: 96px;

  height: var(--role-modal-body-height);
  max-height: var(--role-modal-body-height);
  overflow-y: auto;
  padding-right: 6px;
  margin-right: -6px;
}

.role-form-scroll.is-permission-pinned {
  overflow-y: hidden;
}

.permission-form-item {
  height: var(--role-modal-body-height);
}

.permission-form-item :deep(.n-form-item-blank) {
  height: 100%;
  min-height: 0;
}

.permission-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #fff;
  overflow: hidden;
}

.permission-panel__header {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 10px 48px 10px 14px;
  background: linear-gradient(180deg, #fafafa 0%, #fff 100%);
  border-bottom: 1px solid #f0f0f0;
}

.permission-panel__title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.permission-panel__desc {
  margin-top: 2px;
  font-size: 12px;
  line-height: 16px;
  color: #888;
}

.permission-summary {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
  font-size: 12px;
  color: #666;
}

.permission-pin-btn {
  position: absolute;
  top: 10px;
  right: 14px;
  z-index: 2;
  width: 26px;
  min-width: 26px;
  padding: 0;
}

.permission-pin-btn svg {
  width: 14px;
  height: 14px;
}

.permission-summary span {
  padding: 4px 8px;
  border-radius: 999px;
  background: #f5f7fb;
}

.permission-summary b {
  color: #18a058;
}

.permission-toolbar {
  display: flex;
  gap: 12px;
  padding: 12px 16px 8px;
}

.permission-search {
  flex: 1;
  min-width: 240px;
}

.permission-copy {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.permission-copy__select {
  width: 240px;
}

.permission-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 0 16px 12px;
}

.permission-match-tip {
  flex-shrink: 0;
  font-size: 12px;
  color: #888;
}

.role-permission-tree {
  flex: 1;
  min-height: 0;
  padding: 8px 12px 12px;
  overflow: auto;
  border-top: 1px solid #f0f0f0;
  background: #fcfcfd;
}

:deep(.permission-node-label) {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
  line-height: 22px;
}

:deep(.permission-node-label__name) {
  color: #333;
}

:deep(.permission-node-label__tag) {
  display: inline-flex;
  align-items: center;
  height: 18px;
  padding: 0 6px;
  border-radius: 4px;
  font-size: 12px;
  line-height: 18px;
}

:deep(.permission-node-label__tag.is-menu) {
  color: #2080f0;
  background: rgba(32, 128, 240, 0.1);
}

:deep(.permission-node-label__tag.is-button) {
  color: #18a058;
  background: rgba(24, 160, 88, 0.1);
}

:deep(.permission-node-label__code),
:deep(.permission-node-label__path) {
  font-size: 12px;
  color: #999;
}

@media (max-width: 900px) {
  .permission-panel__header,
  .permission-toolbar,
  .permission-actions {
    align-items: stretch;
    flex-direction: column;
  }

  .permission-copy {
    width: 100%;
  }

  .permission-copy__select {
    flex: 1;
    width: auto;
  }
}
</style>
