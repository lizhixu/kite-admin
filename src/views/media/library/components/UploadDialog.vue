<template>
  <NModal
    v-model:show="visible"
    preset="card"
    title="上传文件"
    style="width: 720px; max-width: 95vw"
    :bordered="false"
    :mask-closable="!hasActive"
    :close-on-esc="!hasActive"
    :on-after-leave="onAfterLeave"
  >
    <div class="flex flex-col gap-12">
      <!-- 目标位置 -->
      <div class="flex items-center gap-12 rounded-6 bg-#f5f6fb px-12 py-8 text-13 dark:bg-#1f1f1f">
        <i class="i-fe:hard-drive opacity-60" />
        <span class="opacity-70">存储:</span>
        <span class="font-medium">{{ configLabel }}</span>
        <span class="mx-8 opacity-30">|</span>
        <i class="i-fe:folder opacity-60" />
        <span class="opacity-70">目录:</span>
        <span class="font-mono">/{{ folderPath }}</span>
      </div>

      <!-- 拖拽区 -->
      <NUpload
        ref="uploadRef"
        :custom-request="customRequest"
        :default-upload="true"
        :multiple="true"
        :accept="accept || undefined"
        :file-list="fileList"
        :show-cancel-button="true"
        :show-retry-button="true"
        :show-remove-button="true"
        :disabled="!configId"
        @update:file-list="onFileListChange"
        @before-upload="onBeforeUpload"
      >
        <NUploadDragger>
          <div class="py-24 text-center">
            <i class="i-fe:upload-cloud text-44 opacity-50" />
            <div class="mt-12 text-14">
              点击选择 或 拖拽文件到此区域
            </div>
            <div class="mt-4 text-12 opacity-50">
              支持一次选择多个文件;<span v-if="acceptHint">仅 {{ acceptHint }} 类型;</span>单文件最大 {{ maxSizeMB }}MB
            </div>
          </div>
        </NUploadDragger>
      </NUpload>

      <!-- 汇总 -->
      <div v-if="fileList.length" class="flex items-center justify-between text-12 opacity-70">
        <span>
          总 {{ fileList.length }} · 完成 {{ stats.finished }} · 失败 {{ stats.error }}
          <span v-if="stats.uploading" class="ml-8">上传中 {{ stats.uploading }}</span>
        </span>
        <NButton size="tiny" text type="primary" :disabled="!fileList.length || hasActive" @click="clearList">
          清空列表
        </NButton>
      </div>
    </div>

    <template #footer>
      <div class="flex items-center justify-between">
        <span class="text-12 opacity-60">
          <template v-if="hasActive">
            上传未完成,关闭将中断剩余任务
          </template>
          <template v-else-if="stats.finished">
            已完成 {{ stats.finished }} 个文件
          </template>
        </span>
        <NSpace>
          <NButton @click="close">
            {{ hasActive ? '关闭(中断)' : '关闭' }}
          </NButton>
          <NButton type="primary" :disabled="!stats.finished" @click="confirmClose">
            完成
          </NButton>
        </NSpace>
      </div>
    </template>
  </NModal>
</template>

<script setup>
import { NButton, NModal, NSpace, NUpload, NUploadDragger } from 'naive-ui'
import api from '../api'

const props = defineProps({
  configs: { type: Array, default: () => [] },
  configId: { type: [Number, null], default: null },
  folderId: { type: [Number, null], default: 0 },
  folderPath: { type: String, default: '' },
  accept: { type: String, default: '' }, // e.g. 'image/*'
})

const emit = defineEmits(['uploaded', 'closed'])

const visible = ref(false)
const uploadRef = ref(null)
const fileList = ref([])

const currentConfig = computed(() => props.configs.find(c => c.id === props.configId))
const configLabel = computed(() => currentConfig.value
  ? `${currentConfig.value.name} (${currentConfig.value.type})`
  : '未选择',
)
const maxSizeMB = computed(() => currentConfig.value?.maxSizeMB || 50)
const allowExtensions = computed(() => {
  const raw = currentConfig.value?.allowExtensions
  if (!raw)
    return null
  return raw.split(',').map(s => s.trim().toLowerCase()).filter(Boolean)
})
const acceptHint = computed(() => allowExtensions.value?.join('/') || '')

const stats = computed(() => {
  const s = { uploading: 0, finished: 0, error: 0, pending: 0 }
  for (const f of fileList.value) {
    if (f.status === 'uploading')
      s.uploading++
    else if (f.status === 'finished')
      s.finished++
    else if (f.status === 'error')
      s.error++
    else
      s.pending++
  }
  return s
})

const hasActive = computed(() => stats.value.uploading > 0)

function open() {
  fileList.value = []
  visible.value = true
}

function close() {
  visible.value = false
}

function confirmClose() {
  visible.value = false
}

function clearList() {
  fileList.value = fileList.value.filter(f => f.status === 'uploading')
}

function onAfterLeave() {
  // 关闭后,如果有完成的文件,通知父级刷新
  if (stats.value.finished > 0)
    emit('uploaded', stats.value.finished)
  fileList.value = []
  emit('closed')
}

function onFileListChange(list) {
  fileList.value = list
}

function onBeforeUpload({ file }) {
  // 客户端校验扩展名 + 大小
  const ext = (file.name?.split('.').pop() || '').toLowerCase()
  if (allowExtensions.value && !allowExtensions.value.includes(ext)) {
    window.$message.error(`${file.name} 扩展名不允许 (${allowExtensions.value.join(',')})`)
    return false
  }
  const size = file.file?.size ?? file.size ?? 0
  if (maxSizeMB.value > 0 && size > maxSizeMB.value * 1024 * 1024) {
    window.$message.error(`${file.name} 超过 ${maxSizeMB.value}MB`)
    return false
  }
  return true
}

async function customRequest({ file, onFinish, onError, onProgress }) {
  try {
    await api.upload(
      file.file,
      {
        configId: props.configId,
        folderId: props.folderId || undefined,
      },
      pct => onProgress?.({ percent: pct }),
    )
    onFinish?.()
  }
  catch (err) {
    console.error(err)
    onError?.()
  }
}

defineExpose({ open, close })
</script>
