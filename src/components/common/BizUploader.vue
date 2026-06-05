<template>
  <div class="biz-uploader" :class="`mode-${mode}`">
    <!-- avatar -->
    <template v-if="mode === 'avatar'">
      <div
        class="avatar-wrap"
        :style="{ width: `${size}px`, height: `${size}px` }"
        @click="trigger"
      >
        <NAvatar round :size="size" :src="firstUrl" />
        <div class="avatar-mask">
          <NProgress
            v-if="firstUploading"
            type="circle"
            :percentage="firstUploading.percent"
            :stroke-width="10"
            style="width: 60%"
          />
          <template v-else>
            <i class="i-fe:camera text-20" />
            <span class="mt-2 text-10">更换</span>
          </template>
        </div>
      </div>
    </template>

    <!-- image-card -->
    <template v-else-if="mode === 'image-card'">
      <NImageGroup>
        <div class="card-grid">
          <div v-for="(url, i) in urls" :key="`v-${url}`" class="card">
            <NImage :src="url" object-fit="cover" class="card-img" />
            <div class="card-mask">
              <i
                class="i-fe:trash-2 cursor-pointer text-16 hover:text-red-400"
                @click.stop="removeAt(i)"
              />
            </div>
          </div>
          <div v-for="up in uploading" :key="`up-${up.id}`" class="card uploading-card">
            <NProgress
              type="circle"
              :percentage="up.percent"
              :stroke-width="10"
              style="width: 60%"
            />
            <div v-if="up.error" class="error-tip">
              {{ up.error }}
            </div>
          </div>
          <div v-if="canAdd" class="card add-card" @click="trigger">
            <i class="i-fe:plus text-24" />
            <div class="mt-4 text-12 opacity-60">
              {{ addText }}
            </div>
          </div>
        </div>
      </NImageGroup>
    </template>

    <!-- button -->
    <template v-else>
      <NButton :loading="!!uploading.length" :disabled="!canAdd" @click="trigger">
        <template #icon>
          <i class="i-fe:upload" />
        </template>
        {{ buttonText }}
      </NButton>
      <div v-if="urls.length || uploading.length" class="mt-8 flex flex-col gap-4">
        <div
          v-for="(url, i) in urls"
          :key="`v-${url}`"
          class="flex items-center gap-8 rounded-4 bg-#f5f5f5 px-8 py-4 text-12 dark:bg-#2a2a2a"
        >
          <i class="i-fe:file flex-shrink-0 opacity-60" />
          <span class="flex-1 truncate">{{ filenameOf(url) }}</span>
          <i
            class="i-fe:x flex-shrink-0 cursor-pointer opacity-60 hover:opacity-100"
            @click="removeAt(i)"
          />
        </div>
        <div
          v-for="up in uploading"
          :key="`up-${up.id}`"
          class="rounded-4 bg-#f5f5f5 px-8 py-4 text-12 dark:bg-#2a2a2a"
        >
          <div class="flex items-center gap-8">
            <i class="i-fe:upload-cloud flex-shrink-0 opacity-60" />
            <span class="flex-1 truncate">{{ up.name }}</span>
            <span class="flex-shrink-0 opacity-60">{{ up.percent }}%</span>
          </div>
          <NProgress
            :percentage="up.percent"
            :show-indicator="false"
            :height="3"
            class="mt-4"
            :status="up.error ? 'error' : 'default'"
          />
          <div v-if="up.error" class="mt-2 text-red-400">
            {{ up.error }}
          </div>
        </div>
      </div>
    </template>

    <input
      ref="inputRef"
      type="file"
      hidden
      :multiple="multiple"
      :accept="acceptAttr"
      @change="onFileChange"
    >
  </div>
</template>

<script setup>
import { NAvatar, NButton, NImage, NImageGroup, NProgress } from 'naive-ui'
import { computed, onMounted, reactive, ref } from 'vue'
import { getBizUploadSpecs, getCachedBizUploadSpec } from '@/composables/useBizUploadSpecs'
import { request } from '@/utils'
import { mediaAccessUrl } from '@/views/media/library/url'

const props = defineProps({
  // 业务类型；后端 bizUploadSpecs 必须包含该 key
  bizType: { type: String, required: true },
  // 'button' | 'image-card' | 'avatar'
  mode: { type: String, default: 'button' },
  // 多选；avatar 模式强制 false
  multiple: { type: Boolean, default: false },
  // 多选上限，0 = 不限
  maxCount: { type: Number, default: 0 },
  // accept 属性；不传则按 bizType 推断
  accept: { type: String, default: '' },
  // 自定义最大体积（MB），覆盖后端 spec；0 = 用 spec
  maxSize: { type: Number, default: 0 },
  // button 模式按钮文案
  buttonText: { type: String, default: '上传文件' },
  // image-card 添加按钮文案
  addText: { type: String, default: '上传' },
  // avatar 模式头像尺寸
  size: { type: Number, default: 100 },
})

const emit = defineEmits(['update:value', 'success', 'error', 'remove'])

const value = defineModel('value')

const inputRef = ref(null)
const uploading = ref([]) // { id, name, percent, error }
let uploadSeq = 0

const isAvatarMode = computed(() => props.mode === 'avatar')
const effectiveMultiple = computed(() => props.multiple && !isAvatarMode.value)

const urls = computed(() => {
  if (!value.value) return []
  return Array.isArray(value.value) ? value.value.filter(Boolean) : [value.value]
})

const firstUrl = computed(() => urls.value[0] || '')
const firstUploading = computed(() => uploading.value[0])

const canAdd = computed(() => {
  if (!effectiveMultiple.value) return uploading.value.length === 0
  if (props.maxCount > 0 && urls.value.length + uploading.value.length >= props.maxCount)
    return false
  return true
})

const spec = ref(null)
onMounted(async () => {
  spec.value = getCachedBizUploadSpec(props.bizType)
  if (!spec.value) {
    try {
      const all = await getBizUploadSpecs()
      spec.value = all?.[props.bizType] || null
    }
    catch (err) {
      console.warn('[BizUploader] 加载白名单失败，跳过客户端校验', err)
    }
  }
})

const acceptAttr = computed(() => {
  if (props.accept) return props.accept
  const s = spec.value
  if (!s) return ''
  const parts = []
  if (s.allowMimePrefix?.length) {
    for (const p of s.allowMimePrefix) parts.push(p.endsWith('/') ? `${p}*` : p)
  }
  if (s.allowExtensions?.length) {
    for (const ext of s.allowExtensions) parts.push(`.${ext}`)
  }
  return parts.join(',')
})

const effectiveMaxSize = computed(() => props.maxSize || spec.value?.maxSizeMB || 0)

function trigger() {
  if (!canAdd.value) return
  inputRef.value?.click()
}

function validateFile(file) {
  const s = spec.value
  const maxMB = effectiveMaxSize.value
  if (maxMB > 0 && file.size > maxMB * 1024 * 1024) {
    return `文件 ${file.name} 超过 ${maxMB}MB 上限`
  }
  if (s?.allowExtensions?.length) {
    const ext = (file.name.split('.').pop() || '').toLowerCase()
    if (!s.allowExtensions.includes(ext))
      return `文件 ${file.name} 扩展名不在允许范围（${s.allowExtensions.join('/')}）`
  }
  if (s?.allowMimePrefix?.length) {
    const ok = s.allowMimePrefix.some(p => (file.type || '').startsWith(p))
    if (!ok) return `文件 ${file.name} 类型不允许（${file.type || '未知'}）`
  }
  return ''
}

async function onFileChange(e) {
  const files = Array.from(e.target.files || [])
  e.target.value = ''
  if (!files.length) return

  let picked = files
  if (!effectiveMultiple.value) picked = [files[0]]
  if (effectiveMultiple.value && props.maxCount > 0) {
    const room = Math.max(0, props.maxCount - urls.value.length - uploading.value.length)
    picked = picked.slice(0, room)
    if (room === 0) {
      window.$message?.warning(`最多上传 ${props.maxCount} 个文件`)
      return
    }
  }

  for (const file of picked) {
    const err = validateFile(file)
    if (err) {
      window.$message?.error(err)
      continue
    }
    await doUpload(file)
  }
}

async function doUpload(file) {
  const id = ++uploadSeq
  const entry = reactive({ id, name: file.name, percent: 0, error: '' })
  uploading.value.push(entry)

  const fd = new FormData()
  fd.append('bizType', props.bizType)
  fd.append('file', file)

  try {
    const res = await request.post('/upload', fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (ev) => {
        if (ev.total) entry.percent = Math.round((ev.loaded / ev.total) * 100)
      },
    })
    const media = res?.data
    const accessUrl = mediaAccessUrl(media)
    if (!accessUrl) throw new Error('上传响应缺少 accessUrl')
    if (effectiveMultiple.value) {
      const next = Array.isArray(value.value) ? value.value.slice() : []
      next.push(accessUrl)
      value.value = next
    }
    else {
      value.value = accessUrl
    }
    emit('success', media, file)
    removeUploading(id)
  }
  catch (err) {
    console.error(err)
    entry.error = err?.response?.data?.message || err?.message || '上传失败'
    emit('error', err, file)
    setTimeout(() => removeUploading(id), 3000)
  }
}

function removeUploading(id) {
  const idx = uploading.value.findIndex(u => u.id === id)
  if (idx >= 0) uploading.value.splice(idx, 1)
}

function removeAt(i) {
  const removed = urls.value[i]
  if (effectiveMultiple.value) {
    const next = (Array.isArray(value.value) ? value.value : []).slice()
    next.splice(i, 1)
    value.value = next
  }
  else {
    value.value = ''
  }
  emit('remove', removed, i)
}

function filenameOf(url) {
  try {
    return decodeURIComponent(url.split('/').pop() || url)
  }
  catch {
    return url
  }
}

</script>

<style scoped>
.biz-uploader {
  display: inline-block;
}

/* avatar mode */
.avatar-wrap {
  position: relative;
  display: inline-block;
  cursor: pointer;
  border-radius: 50%;
  overflow: hidden;
}
.avatar-mask {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: rgba(0, 0, 0, 0.45);
  opacity: 0;
  transition: opacity 0.2s;
  pointer-events: none;
}
.avatar-wrap:hover .avatar-mask {
  opacity: 1;
}

/* image-card mode */
.card-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.card {
  position: relative;
  width: 96px;
  height: 96px;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid var(--light_border, #e5e7eb);
}
.card-img {
  width: 100%;
  height: 100%;
}
.card-img :deep(img) {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.card-mask {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.45);
  opacity: 0;
  transition: opacity 0.2s;
  color: #fff;
}
.card:hover .card-mask {
  opacity: 1;
}
.uploading-card {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}
:global(.dark) .uploading-card {
  background: #2a2a2a;
}
.uploading-card .error-tip {
  position: absolute;
  bottom: 4px;
  left: 4px;
  right: 4px;
  font-size: 10px;
  color: #f87171;
  text-align: center;
  line-height: 1.2;
}
.add-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #9ca3af;
  border-style: dashed;
}
.add-card:hover {
  color: var(--primary-color, #2080f0);
  border-color: currentColor;
}
</style>
