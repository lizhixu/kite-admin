<template>
  <CommonPage>
    <template #action>
      <NSpace>
        <NSelect
          v-model:value="currentConfigId"
          :options="configOptions"
          placeholder="上传到..."
          style="width: 200px"
          :consistent-menu-width="false"
        />
        <NUpload
          v-permission="'UploadMedia'"
          :show-file-list="false"
          :custom-request="onUpload"
          :multiple="true"
        >
          <NButton type="primary">
            <i class="i-fe:upload mr-4 text-16" />
            上传文件
          </NButton>
        </NUpload>
        <NButton
          v-permission="'DeleteMedia'"
          type="error"
          :disabled="!checkedIds.length"
          @click="handleBulkDelete"
        >
          <i class="i-material-symbols:delete-outline mr-4 text-16" />
          批量删除{{ checkedIds.length ? `(${checkedIds.length})` : '' }}
        </NButton>
        <NButton v-permission="'ManageStorage'" @click="openConfigDrawer">
          <i class="i-fe:settings mr-4 text-16" />
          存储设置
        </NButton>
      </NSpace>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1300"
      :columns="columns"
      :get-data="api.read"
      @on-checked="ids => (checkedIds = ids)"
    >
      <MeQueryItem label="文件名" :label-width="60">
        <NInput v-model:value="queryItems.filename" placeholder="请输入文件名" clearable />
      </MeQueryItem>
      <MeQueryItem label="类型" :label-width="60">
        <NSelect
          v-model:value="queryItems.mimePrefix"
          clearable
          :options="[
            { label: '图片', value: 'image/' },
            { label: '视频', value: 'video/' },
            { label: '音频', value: 'audio/' },
            { label: '文档', value: 'application/' },
            { label: '文本', value: 'text/' },
          ]"
        />
      </MeQueryItem>
      <MeQueryItem label="存储后端" :label-width="70">
        <NSelect
          v-model:value="queryItems.storageType"
          clearable
          :options="[
            { label: '本地', value: 'LOCAL' },
            { label: 'S3', value: 'S3' },
          ]"
        />
      </MeQueryItem>
    </MeCrud>

    <StorageConfigDrawer ref="configDrawerRef" @configs-changed="loadConfigs" />
  </CommonPage>
</template>

<script setup>
import { NButton, NImage, NInput, NSelect, NSpace, NTag, NUpload } from 'naive-ui'
import { MeCrud, MeQueryItem } from '@/components'
import { withPermission } from '@/directives'
import { formatDateTime } from '@/utils'
import api from './api'
import StorageConfigDrawer from './StorageConfigDrawer.vue'

defineOptions({ name: 'MediaMgt' })

const $table = ref(null)
const queryItems = ref({})
const checkedIds = ref([])
const configs = ref([])
const currentConfigId = ref(null)
const configDrawerRef = ref(null)

const configOptions = computed(() =>
  configs.value
    .filter(c => c.enabled)
    .map(c => ({
      label: `${c.name}${c.isDefault ? ' (默认)' : ''}`,
      value: c.id,
    })),
)

const configMap = computed(() => {
  const m = {}
  for (const c of configs.value) m[c.id] = c
  return m
})

async function loadConfigs() {
  try {
    const { data = [] } = await api.listConfigs()
    configs.value = data || []
    if (!currentConfigId.value) {
      const def = configs.value.find(c => c.isDefault)
      currentConfigId.value = def?.id ?? configs.value[0]?.id ?? null
    }
  }
  catch (error) {
    console.error(error)
  }
}

onMounted(async () => {
  await loadConfigs()
  $table.value?.handleSearch()
})

function openConfigDrawer() {
  configDrawerRef.value?.open()
}

async function onUpload({ file, onFinish, onError, onProgress }) {
  try {
    await api.upload(file.file, currentConfigId.value, (pct) => {
      onProgress?.({ percent: pct })
    })
    $message.success(`${file.name} 上传成功`)
    onFinish?.()
    $table.value?.handleSearch(true)
  }
  catch (err) {
    console.error(err)
    onError?.()
  }
}

function handleBulkDelete() {
  if (!checkedIds.value.length)
    return
  $dialog.warning({
    title: '提示',
    content: `确定删除选中的 ${checkedIds.value.length} 个文件？`,
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.bulkDelete(checkedIds.value)
        $message.success('删除成功')
        checkedIds.value = []
        $table.value?.handleSearch(true)
      }
      catch (err) {
        console.error(err)
      }
    },
  })
}

function handleDelete(id) {
  $dialog.warning({
    title: '提示',
    content: '确定删除该文件？',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.delete(id)
        $message.success('删除成功')
        $table.value?.handleSearch(true)
      }
      catch (err) {
        console.error(err)
      }
    },
  })
}

async function copyUrl(url) {
  try {
    await navigator.clipboard.writeText(url)
    $message.success('已复制 URL')
  }
  catch {
    $message.error('复制失败')
  }
}

function humanSize(bytes) {
  if (bytes == null)
    return '-'
  if (bytes < 1024)
    return `${bytes} B`
  if (bytes < 1024 * 1024)
    return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024)
    return `${(bytes / 1024 / 1024).toFixed(1)} MB`
  return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`
}

function iconForMime(mime = '') {
  if (mime.startsWith('image/'))
    return 'i-fe:image'
  if (mime.startsWith('video/'))
    return 'i-fe:film'
  if (mime.startsWith('audio/'))
    return 'i-fe:music'
  if (mime.includes('pdf'))
    return 'i-fe:file-text'
  if (mime.includes('zip') || mime.includes('compressed'))
    return 'i-fe:archive'
  return 'i-fe:file'
}

const columns = [
  { type: 'selection', width: 40, fixed: 'left' },
  {
    title: '预览',
    key: 'preview',
    width: 80,
    render(row) {
      if (row.mimeType?.startsWith('image/')) {
        return h(NImage, {
          src: row.url,
          width: 48,
          height: 48,
          objectFit: 'cover',
          style: 'border-radius: 4px;',
        })
      }
      return h('div', {
        class: 'f-c-c',
        style: 'width:48px;height:48px;border-radius:4px;background:var(--n-action-color,#f5f5f5);',
      }, [h('i', { class: `${iconForMime(row.mimeType)} text-22 op-60` })])
    },
  },
  { title: '文件名', key: 'filename', width: 220, ellipsis: { tooltip: true } },
  { title: '类型', key: 'mimeType', width: 140, ellipsis: { tooltip: true } },
  {
    title: '大小',
    key: 'size',
    width: 100,
    render: row => h('span', humanSize(row.size)),
  },
  {
    title: '存储',
    key: 'storageType',
    width: 120,
    render(row) {
      const cfg = configMap.value[row.configId]
      const label = cfg ? `${cfg.name}` : row.storageType
      return h(NTag, {
        type: row.storageType === 'S3' ? 'info' : 'default',
        bordered: false,
        size: 'small',
      }, { default: () => label })
    },
  },
  { title: '上传人', key: 'uploaderName', width: 100, ellipsis: { tooltip: true } },
  { title: '上传时间', key: 'createTime', width: 170, render: row => h('span', formatDateTime(row.createTime)) },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    align: 'right',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          { size: 'small', type: 'primary', secondary: true, onClick: () => copyUrl(row.url) },
          {
            default: () => '复制URL',
            icon: () => h('i', { class: 'i-fe:copy text-14' }),
          },
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
            {
              default: () => '删除',
              icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
            },
          ),
          'DeleteMedia',
        ),
      ]
    },
  },
]
</script>
