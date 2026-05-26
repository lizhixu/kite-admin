<template>
  <NModal
    v-model:show="visible"
    preset="card"
    :title="title"
    style="width: 960px; max-width: 95vw"
    :bordered="false"
    :mask-closable="false"
    :on-after-leave="handleAfterLeave"
  >
    <div class="picker-body flex flex-col">
      <!-- 顶部条 -->
      <div class="mb-12 flex flex-shrink-0 items-center gap-12">
        <NSelect
          v-model:value="currentConfigId"
          :options="configOptions"
          :consistent-menu-width="false"
          placeholder="选择存储"
          style="width: 180px"
          size="small"
          :disabled="!!opts.configId"
        />
        <NRadioGroup v-model:value="activeTab" size="small">
          <NRadioButton value="exist">
            <i class="i-fe:image mr-4 text-12" /> 选择已有
          </NRadioButton>
          <NRadioButton v-if="opts.uploadable !== false" value="upload">
            <i class="i-fe:upload mr-4 text-12" /> 上传新增
          </NRadioButton>
        </NRadioGroup>
        <div class="ml-auto text-12 opacity-60">
          {{ multiple ? `已选 ${checkedIds.length}${max > 0 ? ` / ${max}` : ''}` : (checkedIds.length ? '已选 1 项' : '请选择') }}
        </div>
      </div>

      <!-- 目标目录提示 -->
      <div
        v-if="folderPathLabel"
        class="mb-8 flex flex-shrink-0 items-center gap-6 rounded-4 bg-#f5f5f5 px-10 py-6 text-12 dark:bg-#2a2a2a"
      >
        <i class="i-fe:folder text-14 opacity-60" />
        <span class="opacity-60">目标目录：</span>
        <span class="font-mono">{{ folderPathLabel }}</span>
        <span v-if="opts.folderPath" class="ml-4 text-10 opacity-40">（由调用方指定）</span>
      </div>

      <div class="min-h-0 flex flex-1 gap-12">
        <!-- 左侧文件夹树（两个 tab 共享） -->
        <div class="w-200 flex-shrink-0 border-r border-light_border pr-12 dark:border-dark_border">
          <FolderTree
            v-model:selected-folder-id="currentFolderId"
            :config-id="currentConfigId"
            :readonly="activeTab === 'exist' && !!opts.folderPath"
          />
        </div>

        <!-- 右侧内容区 -->
        <div class="min-w-0 flex flex-col flex-1">
          <template v-if="activeTab === 'exist'">
            <NInput
              v-model:value="filters.filename"
              size="small"
              placeholder="搜索文件名"
              clearable
              class="mb-8"
            >
              <template #prefix>
                <i class="i-fe:search opacity-60" />
              </template>
            </NInput>
            <div class="min-h-0 flex-1">
              <MediaGrid
                ref="gridRef"
                :config-id="currentConfigId"
                :folder-id="currentFolderId"
                :filters="filters"
                select-mode
                :multiple="multiple"
                :max="max"
                :accept-prefix="acceptPrefix"
                @update:checked="ids => (checkedIds = ids)"
                @pick="onPickSingle"
              />
            </div>
          </template>

          <template v-else>
            <div class="picker-upload min-h-0 flex flex-1 flex-col">
              <NUpload
                :show-file-list="true"
                :custom-request="onUpload"
                :multiple="multiple"
                :accept="acceptAttr"
                :disabled="!currentConfigId"
                class="picker-upload-inner"
              >
                <NUploadDragger>
                  <div class="py-24 text-center">
                    <i class="i-fe:upload-cloud text-40 opacity-50" />
                    <div class="mt-8">
                      点击或拖拽文件到此处上传
                    </div>
                    <div class="mt-4 text-12 opacity-50">
                      文件将上传到左侧选中的目录
                    </div>
                  </div>
                </NUploadDragger>
              </NUpload>
            </div>
          </template>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-8">
        <NButton @click="cancel">
          取消
        </NButton>
        <NButton type="primary" :disabled="!canConfirm" @click="confirm">
          确定{{ checkedIds.length ? ` (${checkedIds.length})` : '' }}
        </NButton>
      </div>
    </template>
  </NModal>
</template>

<script setup>
import {
  NButton,
  NInput,
  NModal,
  NRadioButton,
  NRadioGroup,
  NSelect,
  NUpload,
  NUploadDragger,
} from 'naive-ui'
import api from '@/views/media/library/api'
import FolderTree from '@/views/media/library/components/FolderTree.vue'
import MediaGrid from '@/views/media/library/components/MediaGrid.vue'

const props = defineProps({
  opts: { type: Object, default: () => ({}) },
})

const visible = ref(false)
const activeTab = ref('exist')
const configs = ref([])
const currentConfigId = ref(null)
const currentFolderId = ref(0)
const checkedIds = ref([])
const uploadedItems = ref([])
const gridRef = ref(null)
const folderPathLabel = ref('')

const filters = reactive({
  filename: '',
  mimePrefix: null,
})

const opts = computed(() => props.opts || {})
const multiple = computed(() => !!opts.value.multiple)
const max = computed(() => opts.value.max || 0)
const acceptPrefix = computed(() => opts.value.accept || '')
const title = computed(() => opts.value.title || '选择媒体')

const acceptAttr = computed(() => {
  const a = acceptPrefix.value
  if (!a) return undefined
  if (a.endsWith('/')) return `${a}*`
  return a
})

const configOptions = computed(() =>
  configs.value
    .filter(c => c.enabled)
    .map(c => ({
      label: `${c.name}${c.isDefault ? ' (默认)' : ''}`,
      value: c.id,
    })),
)

const canConfirm = computed(() => checkedIds.value.length > 0)

async function loadConfigs() {
  try {
    const { data = [] } = await api.listConfigs()
    configs.value = data || []
  }
  catch (err) {
    console.error(err)
  }
}

// 根据 folderPath 解析（或自动创建）文件夹，返回 folderId
async function resolveFolderPath(configId, folderPath) {
  if (!folderPath || !configId) return 0
  try {
    const { data } = await api.resolveFolder(configId, folderPath, true)
    return data?.id || 0
  }
  catch (err) {
    console.error('resolve folder path failed:', err)
    return 0
  }
}

// 获取文件夹的完整路径标签
function getFolderPathLabel(folders, folderId) {
  if (!folderId) return '/'
  const map = new Map(folders.map(f => [f.id, f]))
  const parts = []
  let cur = map.get(folderId)
  while (cur) {
    parts.unshift(cur.name)
    cur = cur.parentId ? map.get(cur.parentId) : null
  }
  return '/' + parts.join('/')
}

let _resolveFn = null

async function open(userOpts = {}) {
  // 重置状态
  checkedIds.value = []
  uploadedItems.value = []
  activeTab.value = 'exist'
  filters.filename = ''
  filters.mimePrefix = null
  folderPathLabel.value = ''

  // 加载配置
  await loadConfigs()

  // 决定默认 configId
  if (userOpts.configId != null) {
    currentConfigId.value = userOpts.configId
  }
  else {
    const def = configs.value.find(c => c.isDefault && c.enabled)
    currentConfigId.value = def?.id ?? configs.value.find(c => c.enabled)?.id ?? null
  }

  // 决定 folderId：优先 folderPath，其次 folderId
  if (userOpts.folderPath) {
    currentFolderId.value = await resolveFolderPath(currentConfigId.value, userOpts.folderPath)
    // 加载文件夹树以获取路径标签
    try {
      const { data = [] } = await api.listFolders(currentConfigId.value)
      folderPathLabel.value = getFolderPathLabel(data || [], currentFolderId.value)
    }
    catch {
      folderPathLabel.value = userOpts.folderPath
    }
  }
  else {
    currentFolderId.value = userOpts.folderId || 0
    if (currentFolderId.value) {
      try {
        const { data = [] } = await api.listFolders(currentConfigId.value)
        folderPathLabel.value = getFolderPathLabel(data || [], currentFolderId.value)
      }
      catch {
        folderPathLabel.value = ''
      }
    }
  }

  visible.value = true
  return new Promise((resolve) => {
    _resolveFn = resolve
  })
}

watch(currentConfigId, () => {
  currentFolderId.value = 0
  checkedIds.value = []
})

watch(currentFolderId, async (newId) => {
  if (!newId) {
    folderPathLabel.value = opts.value.folderPath ? folderPathLabel.value : ''
    return
  }
  // 如果已有 folderPath 且 folderId 没变，不重新计算
  if (opts.value.folderPath && newId === currentFolderId.value) return
  // 实时更新路径标签
  try {
    const { data = [] } = await api.listFolders(currentConfigId.value)
    folderPathLabel.value = getFolderPathLabel(data || [], newId)
  }
  catch {}
})

function onPickSingle(item) {
  if (!multiple.value) {
    emitResolve([item])
    visible.value = false
  }
}

async function onUpload({ file, onFinish, onError }) {
  try {
    const res = await api.upload(file.file, {
      configId: currentConfigId.value,
      folderId: currentFolderId.value || undefined,
    })
    const media = res?.data
    if (media) {
      uploadedItems.value.push(media)
      if (multiple.value) {
        if (max.value === 0 || checkedIds.value.length < max.value)
          checkedIds.value = [...checkedIds.value, media.id]
      }
      else {
        checkedIds.value = [media.id]
      }
    }
    window.$message.success(`${file.name} 上传成功`)
    onFinish?.()
    gridRef.value?.refresh()
  }
  catch (err) {
    console.error(err)
    onError?.()
  }
}

function confirm() {
  if (!checkedIds.value.length) return
  const fromGrid = gridRef.value?.getCheckedItems?.() || []
  const map = new Map()
  for (const it of fromGrid)
    map.set(it.id, it)
  for (const it of uploadedItems.value)
    map.set(it.id, it)
  const result = checkedIds.value.map(id => map.get(id)).filter(Boolean)
  emitResolve(result)
  visible.value = false
}

function cancel() {
  visible.value = false
}

function emitResolve(items) {
  const fn = _resolveFn
  _resolveFn = null
  fn?.(items || [])
}

function handleAfterLeave() {
  if (_resolveFn) {
    const fn = _resolveFn
    _resolveFn = null
    fn([])
  }
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.$picker = { open }
  }
})

defineExpose({ open })
</script>

<style scoped>
.picker-body {
  height: min(640px, calc(100vh - 220px));
  min-height: 360px;
}

.picker-upload {
  overflow-y: auto;
  padding-right: 4px;
}
.picker-upload::-webkit-scrollbar {
  width: 6px;
}
.picker-upload::-webkit-scrollbar-thumb {
  background: #d4d4d8;
  border-radius: 3px;
}

.picker-upload-inner {
  width: 100%;
}
</style>
