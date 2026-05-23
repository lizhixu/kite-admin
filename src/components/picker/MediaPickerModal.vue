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
          :disabled="!!props.opts.configId"
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

      <div v-if="activeTab === 'exist'" class="min-h-0 flex flex-1 gap-12">
        <div class="w-200 flex-shrink-0 border-r border-light_border pr-12 dark:border-dark_border">
          <FolderTree
            v-model:selected-folder-id="currentFolderId"
            :config-id="currentConfigId"
            readonly
          />
        </div>
        <div class="min-w-0 flex flex-col flex-1">
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
        </div>
      </div>

      <div v-else class="picker-upload min-h-0 flex flex-1 flex-col">
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
                上传后将自动加入"已选"
              </div>
            </div>
          </NUploadDragger>
        </NUpload>
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
const uploadedItems = ref([]) // 累积本次会话上传的 media
const gridRef = ref(null)

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
  // 给 <input type=file> 用的 accept,如 "image/*"
  const a = acceptPrefix.value
  if (!a)
    return undefined
  if (a.endsWith('/'))
    return `${a}*`
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

let _resolveFn = null

async function open(userOpts = {}) {
  // 重置状态
  checkedIds.value = []
  uploadedItems.value = []
  activeTab.value = 'exist'
  filters.filename = ''
  filters.mimePrefix = null

  // 加载配置
  await loadConfigs()

  // 决定默认 configId/folderId
  if (userOpts.configId != null) {
    currentConfigId.value = userOpts.configId
  }
  else {
    const def = configs.value.find(c => c.isDefault && c.enabled)
    currentConfigId.value = def?.id ?? configs.value.find(c => c.enabled)?.id ?? null
  }
  currentFolderId.value = userOpts.folderId || 0

  visible.value = true
  return new Promise((resolve) => {
    _resolveFn = resolve
  })
}

watch(currentConfigId, () => {
  currentFolderId.value = 0
  checkedIds.value = []
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
      // 自动加入已选
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
    // 刷新已有列表,并切回选择 tab 让用户看到结果
    gridRef.value?.refresh()
  }
  catch (err) {
    console.error(err)
    onError?.()
  }
}

function confirm() {
  if (!checkedIds.value.length)
    return
  // 优先用 grid 的已加载项,补充上传新增项
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
  // Modal 关闭未确认时返回空数组
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
