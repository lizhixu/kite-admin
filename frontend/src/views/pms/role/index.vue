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
      <NButton v-permission="'AddRole'" type="primary" @click="handleAdd()">
        <i class="i-material-symbols:add mr-4 text-18" />
        新增角色
      </NButton>
    </template>

    <MeCrud ref="$table" v-model:query-items="queryItems" :scroll-x="1200" :columns="columns" :get-data="api.read">
      <MeQueryItem label="角色名" :label-width="50">
        <n-input v-model:value="queryItems.name" type="text" placeholder="请输入角色名" clearable />
      </MeQueryItem>
      <MeQueryItem label="状态" :label-width="50">
        <n-select v-model:value="queryItems.enable" clearable :options="[
          { label: '启用', value: 1 },
          { label: '停用', value: 0 },
        ]" />
      </MeQueryItem>
    </MeCrud>
    <MeModal ref="modalRef" width="600px">
      <n-form ref="modalFormRef" label-placement="left" label-align="left" :label-width="80" :model="modalForm">
        <n-form-item label="角色名" path="name" :rule="{
          required: true,
          message: '请输入角色名',
          trigger: ['input', 'blur'],
        }">
          <n-input v-model:value="modalForm.name" />
        </n-form-item>
        <n-form-item label="角色编码" path="code" :rule="{
          required: true,
          message: '请输入角色编码',
          trigger: ['input', 'blur'],
        }">
          <n-input v-model:value="modalForm.code" :disabled="modalAction !== 'add'" />
        </n-form-item>
        <n-form-item label="权限" path="permissionIds">
          <div class="w-full">
            <!-- 搜索 + 操作按钮 -->
            <div class="flex items-center gap-8 mb-8">
              <n-input v-model:value="permSearch" size="small" placeholder="搜索权限名称…" clearable style="flex: 1">
                <template #prefix>
                  <i class="i-fe:search text-14" />
                </template>
              </n-input>
              <n-button size="small" @click="selectAllPermissions">全选</n-button>
              <n-button size="small" @click="clearAllPermissions">清空</n-button>
            </div>
            <!-- 已选数量 -->
            <div style="font-size:12px;color:#999;margin-bottom:6px">
              已选 <b style="color:var(--primary-color)">{{ modalForm.permissionIds?.length || 0 }}</b> / {{
                allPermissionIds.length }} 项
            </div>
            <!-- 权限树 -->
            <n-tree key-field="id" label-field="name" :selectable="false" :data="permissionTree"
              :checked-keys="modalForm.permissionIds"
              :on-update:checked-keys="(keys) => (modalForm.permissionIds = getKeysWithAncestors(keys))"
              :pattern="permSearch" checkable check-on-click default-expand-all class="cus-scroll w-full"
              style="max-height:380px" />
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
    initForm: { enable: true },
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
              onClick: () => handleEdit(row),
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
api.getAllPermissionTree().then(({ data = [] }) => (permissionTree.value = data))

// 权限搜索
const permSearch = ref('')

// 收集所有权限 ID（用于全选）
function collectIds(nodes) {
  const ids = []
  const walk = (list) => {
    list.forEach((n) => {
      ids.push(n.id)
      if (n.children?.length) walk(n.children)
    })
  }
  walk(nodes)
  return ids
}
const allPermissionIds = computed(() => collectIds(permissionTree.value))

function selectAllPermissions() {
  modalForm.value.permissionIds = [...allPermissionIds.value]
}
function clearAllPermissions() {
  modalForm.value.permissionIds = []
}

// 选中子节点时自动将所有祖先一并勾选
const parentMap = computed(() => {
  const map = {}
  const walk = (nodes, pid = null) => {
    nodes.forEach((n) => {
      map[n.id] = pid
      if (n.children?.length) walk(n.children, n.id)
    })
  }
  walk(permissionTree.value)
  return map
})

function getKeysWithAncestors(keys) {
  const result = new Set(keys)
  keys.forEach((id) => {
    let pid = parentMap.value[id]
    while (pid != null) {
      result.add(pid)
      pid = parentMap.value[pid]
    }
  })
  return [...result]
}
</script>
