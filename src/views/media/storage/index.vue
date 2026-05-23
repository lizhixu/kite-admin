<template>
  <CommonPage title="存储配置">
    <template #action>
      <NSpace>
        <NButton v-permission="'ManageStorage'" type="primary" @click="handleAdd">
          <i class="i-material-symbols:add mr-4 text-16" />
          新增存储
        </NButton>
        <NButton @click="refresh">
          <i class="i-fe:refresh-cw mr-4 text-16" />
          刷新
        </NButton>
      </NSpace>
    </template>

    <NDataTable
      :columns="columns"
      :data="configs"
      :loading="loading"
      :row-key="row => row.id"
      :bordered="true"
      size="small"
    />

    <MeModal ref="modalRef" width="600px">
      <div class="storage-form-scroll">
        <NForm
          ref="modalFormRef"
          label-placement="left"
          label-align="left"
          :label-width="90"
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

          <template v-if="modalForm.type === 'LOCAL'">
            <NFormItem label="存储目录" path="localDir">
              <NInput v-model:value="modalForm.localDir" placeholder="如：./uploads" />
            </NFormItem>
            <NFormItem label="URL 路径" path="publicPrefix">
              <NInput v-model:value="modalForm.publicPrefix" placeholder="如：/uploads(由后端 Gin 直出的路径)" />
            </NFormItem>
            <NFormItem label="自定义域名" path="customDomain">
              <NInput
                v-model:value="modalForm.customDomain"
                placeholder="可选,如:https://cdn.example.com,设置后图片 URL 走 CDN"
              />
            </NFormItem>
          </template>

          <template v-else>
            <NFormItem
              label="Endpoint"
              path="endpoint"
              :rule="{ required: true, message: '请输入 endpoint', trigger: ['input', 'blur'] }"
            >
              <NInput v-model:value="modalForm.endpoint" placeholder="如：oss-cn-hangzhou.aliyuncs.com" />
            </NFormItem>
            <NGrid :cols="2" :x-gap="12">
              <NGridItem>
                <NFormItem label="Region">
                  <NInput v-model:value="modalForm.region" placeholder="如：cn-hangzhou" />
                </NFormItem>
              </NGridItem>
              <NGridItem>
                <NFormItem
                  label="Bucket"
                  path="bucket"
                  :rule="{ required: true, message: '请输入 bucket', trigger: ['input', 'blur'] }"
                >
                  <NInput v-model:value="modalForm.bucket" />
                </NFormItem>
              </NGridItem>
            </NGrid>
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
            <NFormItem label="自定义域名">
              <NInput v-model:value="modalForm.customDomain" placeholder="可选，如：https://cdn.example.com" />
            </NFormItem>
          </template>

          <NFormItem label="允许扩展名">
            <NInput v-model:value="modalForm.allowExtensions" placeholder="逗号分隔，如 jpg,png,pdf；留空不限" />
          </NFormItem>

          <NGrid :cols="2" :x-gap="12">
            <NGridItem>
              <NFormItem label="最大大小(MB)">
                <NInputNumber v-model:value="modalForm.maxSizeMB" :min="1" :max="10240" style="width: 100%" />
              </NFormItem>
            </NGridItem>
            <NGridItem v-if="modalForm.type === 'S3'">
              <NFormItem label="使用 HTTPS">
                <NSwitch v-model:value="modalForm.useSSL" />
              </NFormItem>
            </NGridItem>
          </NGrid>

          <NGrid :cols="2" :x-gap="12">
            <NGridItem>
              <NFormItem label="启用">
                <NSwitch v-model:value="modalForm.enabled" />
              </NFormItem>
            </NGridItem>
            <NGridItem>
              <NFormItem label="设为默认">
                <NSwitch v-model:value="modalForm.isDefault" />
              </NFormItem>
            </NGridItem>
          </NGrid>
        </NForm>
      </div>
    </MeModal>
  </CommonPage>
</template>

<script setup>
import {
  NButton,
  NDataTable,
  NForm,
  NFormItem,
  NGrid,
  NGridItem,
  NInput,
  NInputNumber,
  NRadio,
  NRadioGroup,
  NSpace,
  NSwitch,
  NTag,
} from 'naive-ui'
import { CommonPage, MeModal } from '@/components'
import { useCrud } from '@/composables'
import api from './api'

defineOptions({ name: 'StorageConfig' })

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
  doCreate: data => api.create(data),
  doUpdate: data => api.update(data),
  doDelete: id => api.delete(id),
  refresh,
})

function handleAdd() {
  crudHandleAdd()
}

function handleEdit(row) {
  crudHandleEdit({ ...row, secretKey: '' })
}

async function refresh() {
  try {
    loading.value = true
    const { data = [] } = await api.list()
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
    window.$message.success('已设为默认')
    await refresh()
  }
  catch (err) {
    console.error(err)
  }
}

async function testConnection(row) {
  try {
    const { data } = await api.test(row.id)
    window.$message.success(`连接成功，耗时 ${data?.elapsedMs ?? 0} ms`)
  }
  catch (err) {
    window.$message.error(`连接失败：${err?.message || '未知错误'}`)
  }
}

onMounted(refresh)

const columns = [
  { title: '名称', key: 'name', width: 160, ellipsis: { tooltip: true } },
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
  { title: 'Bucket / 目录', key: 'target', width: 200, ellipsis: { tooltip: true }, render: row => row.type === 'S3' ? row.bucket : row.localDir },
  { title: '限制', key: 'limit', width: 140, render: row => `≤ ${row.maxSizeMB || '∞'} MB` },
  {
    title: '操作',
    key: 'actions',
    width: 320,
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

<style scoped>
.storage-form-scroll {
  max-height: calc(100vh - 240px);
  overflow-y: auto;
  padding-right: 6px;
  margin-right: -6px;
}
.storage-form-scroll::-webkit-scrollbar {
  width: 6px;
}
.storage-form-scroll::-webkit-scrollbar-thumb {
  background: #d4d4d8;
  border-radius: 3px;
}
</style>
