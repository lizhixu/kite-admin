<template>
  <div class="h-full flex flex-col">
    <div class="mb-12 flex items-center justify-between">
      <span class="text-14 font-medium opacity-70">文件夹</span>
      <NButton
        v-if="!readonly"
        v-permission="'ManageFolder'"
        size="tiny"
        type="primary"
        secondary
        :disabled="!configId"
        @click="onCreate(null)"
      >
        <i class="i-fe:folder-plus mr-2 text-14" />
        新建
      </NButton>
    </div>

    <NScrollbar class="flex-1">
      <NTree
        block-line
        :data="treeData"
        :selected-keys="selectedKeys"
        :default-expanded-keys="defaultExpandedKeys"
        :node-props="nodeProps"
        :render-suffix="renderSuffix"
        @update:selected-keys="onSelect"
      />
    </NScrollbar>

    <NDropdown
      placement="bottom-start"
      trigger="manual"
      :show="showCtxMenu"
      :options="ctxOptions"
      :x="ctxX"
      :y="ctxY"
      @clickoutside="showCtxMenu = false"
      @select="onCtxSelect"
    />

    <NModal v-model:show="renameShow" preset="dialog" title="重命名文件夹" :show-icon="false">
      <NInput
        v-model:value="renameInput"
        placeholder="请输入新名称"
        @keydown.enter="confirmRename"
      />
      <template #action>
        <NButton size="small" @click="renameShow = false">
          取消
        </NButton>
        <NButton size="small" type="primary" :loading="renameLoading" @click="confirmRename">
          确定
        </NButton>
      </template>
    </NModal>

    <NModal v-model:show="createShow" preset="dialog" title="新建文件夹" :show-icon="false">
      <div class="mb-8 text-12 opacity-60">
        将创建在：<span class="font-mono">{{ createUnderLabel }}</span>
      </div>
      <NInput
        v-model:value="createInput"
        placeholder="请输入文件夹名称"
        @keydown.enter="confirmCreate"
      />
      <template #action>
        <NButton size="small" @click="createShow = false">
          取消
        </NButton>
        <NButton size="small" type="primary" :loading="createLoading" @click="confirmCreate">
          确定
        </NButton>
      </template>
    </NModal>
  </div>
</template>

<script setup>
import { NButton, NDropdown, NInput, NModal, NScrollbar, NTree } from 'naive-ui'
import { buildFolderTree } from '../../utils'
import api from '../api'

const props = defineProps({
  configId: { type: [Number, null], default: null },
  selectedFolderId: { type: [Number, null], default: 0 },
  readonly: { type: Boolean, default: false },
})

const emit = defineEmits(['update:selectedFolderId', 'change'])

const folders = ref([])
const treeData = computed(() => buildFolderTree(folders.value))
const selectedKeys = computed(() => [props.selectedFolderId ?? 0])
const defaultExpandedKeys = ref([0])

async function refresh() {
  if (!props.configId) {
    folders.value = []
    return
  }
  try {
    const { data = [] } = await api.listFolders(props.configId)
    folders.value = data || []
  }
  catch (err) {
    console.error(err)
  }
}

watch(() => props.configId, refresh, { immediate: true })

function onSelect(keys, _options, meta) {
  const k = keys[0]
  if (k == null)
    return
  emit('update:selectedFolderId', k)
  emit('change', meta?.node)
}

// --- 右键菜单 ---
const showCtxMenu = ref(false)
const ctxX = ref(0)
const ctxY = ref(0)
const ctxNode = ref(null)
const ctxOptions = computed(() => {
  if (props.readonly)
    return []
  const isRoot = !ctxNode.value || ctxNode.value.id === 0
  const opts = [
    { label: '在此新建子文件夹', key: 'create', icon: () => h('i', { class: 'i-fe:folder-plus' }) },
  ]
  if (!isRoot) {
    opts.push(
      { label: '重命名', key: 'rename', icon: () => h('i', { class: 'i-fe:edit-2' }) },
      { type: 'divider' },
      { label: '删除', key: 'delete', icon: () => h('i', { class: 'i-fe:trash-2' }) },
    )
  }
  return opts
})

function nodeProps({ option }) {
  return {
    onContextmenu(e) {
      if (props.readonly)
        return
      e.preventDefault()
      ctxNode.value = option
      ctxX.value = e.clientX
      ctxY.value = e.clientY
      showCtxMenu.value = false
      nextTick(() => {
        showCtxMenu.value = true
      })
    },
  }
}

function renderSuffix({ option }) {
  if (props.readonly || option.id === 0)
    return null
  // 悬浮显示 path 提示，便于诊断
  return h('span', { class: 'text-10 opacity-40' }, option.path)
}

function onCtxSelect(key) {
  showCtxMenu.value = false
  const node = ctxNode.value
  if (key === 'create')
    onCreate(node)
  else if (key === 'rename')
    openRename(node)
  else if (key === 'delete')
    confirmDelete(node)
}

// --- 新建 ---
const createShow = ref(false)
const createInput = ref('')
const createLoading = ref(false)
const createUnder = ref(null) // null = 根

const createUnderLabel = computed(() => {
  if (!createUnder.value || createUnder.value.id === 0)
    return '/'
  return `/${createUnder.value.path}`
})

function onCreate(node) {
  if (!props.configId) {
    window.$message.warning('请先选择存储')
    return
  }
  createUnder.value = node
  createInput.value = ''
  createShow.value = true
}

async function confirmCreate() {
  const name = createInput.value.trim()
  if (!name)
    return
  createLoading.value = true
  try {
    const parentId = createUnder.value?.id || null
    await api.createFolder({ name, parentId: parentId || null, configId: props.configId })
    window.$message.success('已创建')
    createShow.value = false
    await refresh()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    createLoading.value = false
  }
}

// --- 重命名 ---
const renameShow = ref(false)
const renameInput = ref('')
const renameLoading = ref(false)
const renameTarget = ref(null)

function openRename(node) {
  renameTarget.value = node
  renameInput.value = node.label
  renameShow.value = true
}

async function confirmRename() {
  const name = renameInput.value.trim()
  if (!name || !renameTarget.value)
    return
  renameLoading.value = true
  try {
    await api.renameFolder(renameTarget.value.id, name)
    window.$message.success('已重命名')
    renameShow.value = false
    await refresh()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    renameLoading.value = false
  }
}

// --- 删除 ---
function confirmDelete(node) {
  window.$dialog.warning({
    title: '删除文件夹',
    content: `确定删除文件夹 “${node.label}” 吗？非空文件夹将提示是否级联删除。`,
    positiveText: '删除',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.deleteFolder(node.id, false)
        window.$message.success('已删除')
        if (props.selectedFolderId === node.id)
          emit('update:selectedFolderId', 0)
        await refresh()
      }
      catch (err) {
        // 后端返回 400 表示非空，尝试级联
        const msg = err?.message || ''
        if (msg.includes('not empty')) {
          window.$dialog.warning({
            title: '文件夹非空',
            content: '该文件夹下存在子文件夹或文件，是否一并删除（不可恢复）？',
            positiveText: '级联删除',
            negativeText: '取消',
            async onPositiveClick() {
              try {
                await api.deleteFolder(node.id, true)
                window.$message.success('已删除')
                if (props.selectedFolderId === node.id)
                  emit('update:selectedFolderId', 0)
                await refresh()
              }
              catch (err2) {
                console.error(err2)
              }
            },
          })
        }
      }
    },
  })
}

defineExpose({ refresh })
</script>
