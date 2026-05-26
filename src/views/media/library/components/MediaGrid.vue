<template>
  <div class="h-full flex flex-col">
    <div class="mb-12 flex flex-shrink-0 items-center justify-between gap-12">
      <div class="flex items-center gap-8">
        <NCheckbox
          :checked="allChecked"
          :indeterminate="someChecked && !allChecked"
          :disabled="!items.length"
          @update:checked="toggleAll"
        >
          全选
        </NCheckbox>
        <span v-if="checkedIds.length" class="text-12 opacity-60">
          已选 {{ checkedIds.length }} 项
        </span>
      </div>
      <div class="flex items-center gap-8">
        <slot name="actions" :checked-ids="checkedIds" :clear-checked="clearChecked" />
      </div>
    </div>

    <NScrollbar class="flex-1">
      <NSpin :show="loading">
        <div v-if="!items.length && !loading" class="py-60 text-center opacity-50">
          <i class="i-fe:image mb-8 text-32" />
          <div>暂无文件</div>
        </div>
        <NImageGroup v-else>
          <div class="grid grid-cols-[repeat(auto-fill,minmax(140px,1fr))] gap-12">
            <div
              v-for="item in items"
              :key="item.id"
              class="group relative cursor-pointer border border-light_border rounded-6 bg-white p-8 transition dark:border-dark_border hover:border-primary dark:bg-#1f1f1f"
              :class="{ 'border-primary! ring-2 ring-primary/30': checkedSet.has(item.id) }"
              @click="onItemClick(item, $event)"
            >
              <div class="absolute left-6 top-6 z-1" @click.stop>
                <NCheckbox
                  :checked="checkedSet.has(item.id)"
                  @update:checked="toggle(item.id)"
                />
              </div>
              <div class="mb-6 aspect-square flex items-center justify-center overflow-hidden rounded-4 bg-#f5f5f5 dark:bg-#2a2a2a">
                <NImage
                  v-if="isImage(item)"
                  :src="item.url"
                  object-fit="cover"
                  class="h-full w-full"
                  :img-props="{ style: 'width:100%;height:100%;object-fit:cover' }"
                  @click.stop
                />
                <i
                  v-else
                  :class="iconForMime(item.mimeType)"
                  class="cursor-pointer text-32 opacity-60 hover:opacity-100"
                  @click.stop="openNonImagePreview(item)"
                />
              </div>
              <div class="truncate text-12" :title="item.filename">
                {{ item.filename }}
              </div>
              <div class="mt-2 flex items-center justify-between text-10 opacity-60">
                <span>{{ humanSize(item.size) }}</span>
                <span>{{ formatShortDate(item.createTime) }}</span>
              </div>

              <div class="absolute right-6 top-6 hidden gap-4 group-hover:flex" @click.stop>
                <NButton
                  circle
                  size="tiny"
                  @click="copyUrl(item.url)"
                >
                  <i class="i-fe:copy text-12" />
                </NButton>
                <NButton
                  v-if="!selectMode"
                  v-permission="'DeleteMedia'"
                  circle
                  size="tiny"
                  type="error"
                  @click="onDelete(item)"
                >
                  <i class="i-fe:trash-2 text-12" />
                </NButton>
              </div>
            </div>
          </div>
        </NImageGroup>
      </NSpin>

      <div v-if="total > pageSize" class="mt-16 flex justify-center">
        <NPagination
          v-model:page="page"
          v-model:page-size="pageSize"
          :item-count="total"
          :page-sizes="[24, 48, 96]"
          show-size-picker
          @update:page="fetchList"
          @update:page-size="fetchList"
        />
      </div>
    </NScrollbar>
  </div>
</template>

<script setup>
import { NButton, NCheckbox, NImage, NImageGroup, NPagination, NScrollbar, NSpin } from 'naive-ui'
import { formatDate } from '@/utils'
import { humanSize, iconForMime } from '../../utils'
import api from '../api'

const props = defineProps({
  configId: { type: [Number, null], default: null },
  folderId: { type: [Number, null], default: 0 },
  filters: { type: Object, default: () => ({}) }, // { filename, mimePrefix, storageType }
  selectMode: { type: Boolean, default: false },
  multiple: { type: Boolean, default: true },
  max: { type: Number, default: 0 },
  acceptPrefix: { type: String, default: '' }, // 'image/' etc.
})

const emit = defineEmits(['update:checked', 'pick', 'delete'])

const items = ref([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(24)
const checkedIds = ref([])
const checkedSet = computed(() => new Set(checkedIds.value))

const allChecked = computed(() =>
  items.value.length > 0 && items.value.every(i => checkedSet.value.has(i.id)),
)
const someChecked = computed(() => checkedIds.value.length > 0)

async function fetchList() {
  loading.value = true
  try {
    const params = {
      pageNo: page.value,
      pageSize: pageSize.value,
      ...props.filters,
    }
    if (props.configId != null)
      params.configId = props.configId
    if (props.folderId != null)
      params.folderId = props.folderId
    if (props.acceptPrefix)
      params.mimePrefix = props.acceptPrefix
    const { data } = await api.page(params)
    items.value = data?.pageData || []
    total.value = data?.total || 0
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
}

watch(
  () => [props.configId, props.folderId, props.filters, props.acceptPrefix],
  () => {
    if (!props.configId)
      return
    page.value = 1
    fetchList()
  },
  { deep: true, immediate: true },
)

function toggle(id) {
  const set = new Set(checkedIds.value)
  if (set.has(id)) {
    set.delete(id)
  }
  else {
    if (!props.multiple) {
      checkedIds.value = [id]
      emit('update:checked', checkedIds.value)
      return
    }
    if (props.max > 0 && set.size >= props.max) {
      window.$message?.warning(`最多选择 ${props.max} 项`)
      return
    }
    set.add(id)
  }
  checkedIds.value = [...set]
  emit('update:checked', checkedIds.value)
}

function toggleAll(checked) {
  if (checked) {
    const ids = items.value.map(i => i.id)
    if (props.max > 0)
      checkedIds.value = ids.slice(0, props.max)
    else
      checkedIds.value = ids
  }
  else {
    checkedIds.value = []
  }
  emit('update:checked', checkedIds.value)
}

function clearChecked() {
  checkedIds.value = []
  emit('update:checked', [])
}

function onItemClick(item) {
  if (props.selectMode) {
    toggle(item.id)
    if (!props.multiple)
      emit('pick', item)
  }
}

function isImage(item) {
  return item.mimeType?.startsWith('image/')
}

function formatShortDate(s) {
  return formatDate(s)
}

async function copyUrl(url) {
  try {
    await navigator.clipboard.writeText(url)
    window.$message.success('已复制 URL')
  }
  catch {
    window.$message.error('复制失败')
  }
}

function openNonImagePreview(item) {
  window.open(item.url, '_blank')
}

function onDelete(item) {
  window.$dialog.warning({
    title: '提示',
    content: `确定删除「${item.filename}」？`,
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.delete(item.id)
        window.$message.success('删除成功')
        emit('delete', item)
        await fetchList()
      }
      catch (err) {
        console.error(err)
      }
    },
  })
}

defineExpose({
  refresh: fetchList,
  clearChecked,
  getCheckedItems: () => items.value.filter(i => checkedSet.value.has(i.id)),
})
</script>
