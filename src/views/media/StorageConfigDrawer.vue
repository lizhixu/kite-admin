<template>
  <NDrawer v-model:show="show" :width="720" placement="right">
    <NDrawerContent title="存储设置" closable>
      <NSpace vertical :size="16">
        <NSpace>
          <NButton type="primary" @click="handleAdd">
            <i class="i-material-symbols:add mr-4 text-16" />
            新增存储配置
          </NButton>
          <NButton @click="refresh">
            <i class="i-fe:refresh-cw mr-4 text-16" />
            刷新
          </NButton>
        </NSpace>

        <NDataTable
          :columns="columns"
          :data="configs"
          :loading="loading"
          :row-key="row => row.id"
          :bordered="true"
          size="small"
        />
      </NSpace>

      <MeModal ref="modalRef" width="560px">
        <NForm
          ref="modalFormRef"
          label-placement="left"
          label-align="left"
          :label-width="100"
          :model="modalForm"
        >
          <NFormItem
            label="名称"
            path="name"
            :rule="{ required: true, message: '请输入名称', trigger: ['input', 'blur'] }"
          >
            <NInput v-model:value="modalForm.name" placeholder="如：阿里 OSS 主存储" />
          </NFormItem>

          <NFormItem label="类型" path="type">
            <NRadioGroup v-model:value="modalForm.type" :disabled="modalAction === 'edit'">
              <NRadio value="LOCAL">
                本地磁盘
              </NRadio>
              <NRadio value="S3">
                S3 协议
              </NRadio>
            </NRadioGroup>
          </NFormItem>

          <!-- LOCAL 字段 -->
          <template v-if="modalForm.type === 'LOCAL'">
            <NFormItem label="存储目录" path="localDir">
              <NInput v-model:value="modalForm.localDir" placeholder="如：./uploads" />
            </NFormItem>
            <NFormItem label="URL 前缀" path="publicPrefix">
              <NInput v-model:value="modalForm.publicPrefix" placeholder="如：/uploads 或 https://cdn.example.com" />
            </NFormItem>
          </template>

          <!-- S3 字段 -->
          <template v-else>
            <NFormItem
              label="Endpoint"
              path="endpoint"
              :rule="{ required: true, message: '请输入 endpoint', trigger: ['input', 'blur'] }"
            >
              <NInput v-model:value="modalForm.endpoint" placeholder="如：oss-cn-hangzhou.aliyuncs.com" />
            </NFormItem>
            <NFormItem label="Region">
              <NInput v-model:value="modalForm.region" placeholder="如：cn-hangzhou / us-east-1" />
            </NFormItem>
            <NFormItem
              label="Bucket"
              path="bucket"
              :rule="{ required: true, message: '请输入 bucket', trigger: ['input', 'blur'] }"
            >
              <NInput v-model:value="modalForm.bucket" />
            </NFormItem>
            <NFormItem
              label="AccessKey"
              path="accessKey"
              :rule="{ required: true, message: '请输入 AccessKey', trigger: ['input', 'blur'] }"
            >
              <NInput v-model:value="modalForm.accessKey" />
            </NFormItem>
            <NFormItem
              label="SecretKey"
              path="secretKey"
              :rule="{ required: modalAction === 'add', message: '请输入 SecretKey', trigger: ['input', 'blur'] }"
            >
              <NInput
                v-model:value="modalForm.secretKey"
                type="password"
                show-password-on="mousedown"
                :placeholder="modalAction === 'edit' ? '留空表示保持不变' : ''"
              />
            </NFormItem>
            <NFormItem label="使用 HTTPS">
              <NSwitch v-model:value="modalForm.useSSL" />
            </NFormItem>
            <NFormItem label="自定义域名">
              <NInput v-model:value="modalForm.customDomain" placeholder="可选，如：https://cdn.example.com" />
            </NFormItem>
          </template>

          <NFormItem label="允许扩展名">
            <NInput v-model:value="modalForm.allowExtensions" placeholder="逗号分隔，如 jpg,png,pdf；留空不限" />
          </NFormItem>
          <NFormItem label="最大大小 (MB)">
            <NInputNumber v-model:value="modalForm.maxSizeMB" :min="1" :max="10240" />
          </NFormItem>
          <NFormItem label="启用">
            <NSwitch v-model:value="modalForm.enabled" />
          </NFormItem>
          <NFormItem label="设为默认">
            <NSwitch v-model:value="modalForm.isDefault" />
          </NFormItem>
        </NForm>
      </MeModal>
    </NDrawerContent>
  </NDrawer>
</template>

<script setup>
import {
  NButton,
  NDataTable,
  NDrawer,
  NDrawerContent,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NRadio,
  NRadioGroup,
  NSpace,
  NSwitch,
  NTag,
} from 'naive-ui'
import { MeModal } from '@/components'
import { useCrud } from '@/composables'
import api from './api'

const emit = defineEmits(['configsChanged'])

const show = ref(false)
const configs = ref([])
const loading = ref(false)

const {
  modalRef,
  modalFormRef,
  modalForm,
  modalAction,
  handleAdd: crudHandleAdd,
  handleEdit: crudHandleEdit,
  handleDelete,
} = useCrud({
  name: '存储配置',
  initForm: {
    type: 'S3',
    useSSL: true,
    enabled: true,
    isDefault: false,
    maxSizeMB: 50,
  },
  doCreate: data => api.createConfig(data),
  doUpdate: data => api.updateConfig(data),
  doDelete: id => api.deleteConfig(id),
  refresh: async () => {
    await refresh()
    emit('configsChanged')
  },
})

function handleAdd() {
  crudHandleAdd()
}

function handleEdit(row) {
  // 编辑时把 secretKey 清空（后端拒绝回传），用户留空则保留原值
  crudHandleEdit({ ...row, secretKey: '' })
}

async function refresh() {
  try {
    loading.value = true
    const { data = [] } = await api.listConfigs()
    configs.value = data || []
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
}

async function setDefault(row) {
  try {
    await api.setDefault(row.id)
    $message.success('已设为默认')
    await refresh()
    emit('configsChanged')
  }
  catch (err) {
    console.error(err)
  }
}

async function testConnection(row) {
  try {
    const { data } = await api.testConfig(row.id)
    $message.success(`连接成功，耗时 ${data?.elapsedMs ?? 0} ms`)
  }
  catch (err) {
    console.error(err)
    $message.error(`连接失败：${err?.message || '未知错误'}`)
  }
}

function open() {
  show.value = true
  refresh()
}

defineExpose({ open })

const columns = [
  { title: '名称', key: 'name', width: 140, ellipsis: { tooltip: true } },
  {
    title: '类型',
    key: 'type',
    width: 80,
    render: row =>
      h(NTag, { size: 'small', bordered: false, type: row.type === 'S3' ? 'info' : 'default' }, { default: () => row.type }),
  },
  {
    title: '默认',
    key: 'isDefault',
    width: 70,
    render: row => (row.isDefault
      ? h(NTag, { size: 'small', type: 'success', bordered: false }, { default: () => '默认' })
      : h('span', '-')),
  },
  {
    title: '启用',
    key: 'enabled',
    width: 70,
    render: row => h(NTag, {
      size: 'small',
      type: row.enabled ? 'success' : 'warning',
      bordered: false,
    }, { default: () => (row.enabled ? '是' : '否') }),
  },
  { title: 'Bucket / 目录', key: 'target', width: 170, ellipsis: { tooltip: true }, render: row => row.type === 'S3' ? row.bucket : row.localDir },
  {
    title: '操作',
    key: 'actions',
    width: 280,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'tiny', type: 'primary', secondary: true, onClick: () => testConnection(row) }, { default: () => '测试' }),
          !row.isDefault && h(NButton, { size: 'tiny', type: 'primary', secondary: true, onClick: () => setDefault(row) }, { default: () => '设为默认' }),
          h(NButton, { size: 'tiny', type: 'primary', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
          h(NButton, { size: 'tiny', type: 'error', onClick: () => handleDelete(row.id) }, { default: () => '删除' }),
        ].filter(Boolean),
      })
    },
  },
]
</script>
