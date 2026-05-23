<template>
  <CommonPage title="媒体库">
    <template #action>
      <NSpace>
        <NSelect
          v-model:value="currentConfigId"
          :options="configOptions"
          placeholder="选择存储..."
          style="width: 200px"
          :consistent-menu-width="false"
        />
        <NButton
          v-permission="'UploadMedia'"
          type="primary"
          :disabled="!currentConfigId"
          @click="openUpload"
        >
          <i class="i-fe:upload mr-4 text-16" />
          上传文件
        </NButton>
      </NSpace>
    </template>

    <div class="h-full flex gap-12 overflow-hidden">
      <div class="w-240 flex-shrink-0 border-r border-light_border pr-12 dark:border-dark_border">
        <FolderTree
          ref="folderTreeRef"
          v-model:selected-folder-id="currentFolderId"
          :config-id="currentConfigId"
          @change="onFolderChange"
        />
      </div>

      <div class="min-w-0 flex flex-col flex-1">
        <div class="mb-12 flex flex-shrink-0 items-center gap-12">
          <NRadioGroup v-if="canViewAll" v-model:value="scope" size="small">
            <NRadioButton value="mine">
              我的
            </NRadioButton>
            <NRadioButton value="all">
              全部
            </NRadioButton>
          </NRadioGroup>
          <NInput
            v-model:value="filters.filename"
            placeholder="搜索文件名"
            clearable
            style="width: 240px"
          >
            <template #prefix>
              <i class="i-fe:search opacity-60" />
            </template>
          </NInput>
          <NSelect
            v-model:value="filters.mimePrefix"
            placeholder="文件类型"
            clearable
            style="width: 140px"
            :options="mimeOptions"
          />
          <div v-if="currentFolderPath" class="ml-auto text-12 opacity-60">
            当前路径：<span class="font-mono">/{{ currentFolderPath }}</span>
          </div>
        </div>

        <div class="min-h-0 flex-1">
          <MediaGrid
            ref="gridRef"
            :config-id="currentConfigId"
            :folder-id="currentFolderId"
            :filters="effectiveFilters"
          >
            <template #actions="{ checkedIds, clearChecked }">
              <NButton
                v-if="checkedIds.length"
                size="small"
                @click="openMoveDialog(checkedIds, clearChecked)"
              >
                <i class="i-fe:move mr-4 text-14" />
                移动
              </NButton>
              <NButton
                v-permission="'DeleteMedia'"
                size="small"
                type="error"
                :disabled="!checkedIds.length"
                @click="handleBulkDelete(checkedIds, clearChecked)"
              >
                <i class="i-material-symbols:delete-outline mr-4 text-14" />
                删除{{ checkedIds.length ? `(${checkedIds.length})` : '' }}
              </NButton>
            </template>
          </MediaGrid>
        </div>
      </div>
    </div>

    <NModal v-model:show="moveDialog.show" preset="dialog" title="移动到文件夹" :show-icon="false" style="width: 420px">
      <div class="mb-8 text-12 opacity-60">
        将 {{ moveDialog.ids.length }} 个文件移动到：
      </div>
      <NTreeSelect
        v-model:value="moveDialog.targetId"
        :options="moveTreeOptions"
        default-expand-all
        placeholder="选择目标文件夹"
        key-field="key"
        label-field="label"
      />
      <template #action>
        <NButton size="small" @click="moveDialog.show = false">
          取消
        </NButton>
        <NButton size="small" type="primary" :loading="moveDialog.loading" @click="confirmMove">
          确定
        </NButton>
      </template>
    </NModal>

    <UploadDialog
      ref="uploadDialogRef"
      :configs="configs"
      :config-id="currentConfigId"
      :folder-id="currentFolderId"
      :folder-path="currentFolderPath"
      @uploaded="onBatchUploaded"
    />
  </CommonPage>
</template>

<script setup>
import { NButton, NInput, NModal, NRadioButton, NRadioGroup, NSelect, NSpace, NTreeSelect } from 'naive-ui'
import { useRoute } from 'vue-router'
import { CommonPage } from '@/components'
import { buildFolderTree } from '../utils'
import api from './api'
import FolderTree from './components/FolderTree.vue'
import MediaGrid from './components/MediaGrid.vue'
import UploadDialog from './components/UploadDialog.vue'

defineOptions({ name: 'MediaLibrary' })

const route = useRoute()
const canViewAll = computed(
  () => !!route.meta?.btns?.some(b => b.code === 'ViewAllMedia'),
)
const scope = ref('mine')

const configs = ref([])
const currentConfigId = ref(null)
const currentFolderId = ref(0)
const currentFolderPath = ref('')
const folderTreeRef = ref(null)
const gridRef = ref(null)

const filters = reactive({
  filename: '',
  mimePrefix: null,
})

const effectiveFilters = computed(() => ({
  ...filters,
  scope: canViewAll.value ? scope.value : 'mine',
}))

const configOptions = computed(() =>
  configs.value
    .filter(c => c.enabled)
    .map(c => ({
      label: `${c.name}${c.isDefault ? ' (默认)' : ''}`,
      value: c.id,
    })),
)

const mimeOptions = [
  { label: '图片', value: 'image/' },
  { label: '视频', value: 'video/' },
  { label: '音频', value: 'audio/' },
  { label: '文档', value: 'application/' },
  { label: '文本', value: 'text/' },
]

async function loadConfigs() {
  try {
    const { data = [] } = await api.listConfigs()
    configs.value = data || []
    if (!currentConfigId.value) {
      const def = configs.value.find(c => c.isDefault && c.enabled)
      currentConfigId.value = def?.id ?? configs.value.find(c => c.enabled)?.id ?? null
    }
  }
  catch (err) {
    console.error(err)
  }
}

watch(currentConfigId, () => {
  currentFolderId.value = 0
  currentFolderPath.value = ''
})

function onFolderChange(node) {
  currentFolderPath.value = node?.path || ''
}

onMounted(loadConfigs)

const uploadDialogRef = ref(null)
function openUpload() {
  uploadDialogRef.value?.open()
}
function onBatchUploaded(count) {
  if (count > 0) {
    window.$message.success(`已上传 ${count} 个文件`)
    gridRef.value?.refresh()
  }
}

function handleBulkDelete(ids, clearChecked) {
  if (!ids.length)
    return
  window.$dialog.warning({
    title: '提示',
    content: `确定删除选中的 ${ids.length} 个文件？`,
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.bulkDelete(ids)
        window.$message.success('删除成功')
        clearChecked?.()
        gridRef.value?.refresh()
      }
      catch (err) {
        console.error(err)
      }
    },
  })
}

const moveDialog = reactive({ show: false, ids: [], targetId: 0, loading: false, clear: null })
const moveTreeData = ref([])

const moveTreeOptions = computed(() => mapTreeForSelect(moveTreeData.value))

function mapTreeForSelect(nodes) {
  return (nodes || []).map(n => ({
    key: n.id ?? n.key ?? 0,
    value: n.id ?? n.key ?? 0,
    label: n.label ?? n.name,
    children: n.children?.length ? mapTreeForSelect(n.children) : undefined,
  }))
}

async function openMoveDialog(ids, clear) {
  moveDialog.ids = [...ids]
  moveDialog.clear = clear
  moveDialog.targetId = 0
  try {
    const { data = [] } = await api.listFolders(currentConfigId.value)
    moveTreeData.value = buildFolderTree(data, '根目录')
  }
  catch (err) {
    console.error(err)
  }
  moveDialog.show = true
}

async function confirmMove() {
  moveDialog.loading = true
  try {
    await api.move(moveDialog.ids, moveDialog.targetId || 0)
    window.$message.success('已移动')
    moveDialog.show = false
    moveDialog.clear?.()
    gridRef.value?.refresh()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    moveDialog.loading = false
  }
}
</script>
