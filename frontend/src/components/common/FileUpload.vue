<!--
  FileUpload 可复用上传组件
  Props:
    - accept: string       文件 accept, 默认 '*'
    - maxSizeMB: number    最大文件大小（MB），0=不限，默认 10
    - multiple: boolean    是否多选，默认 true
    - listType: 'text' | 'picture-card'  展示样式，默认 'picture-card'
    - modelValue: Array    已上传文件列表 { url, name, id }（v-model）
  Events:
    - update:modelValue    文件列表变化
    - success(attachment)  单个文件上传成功
    - error(err)           上传失败
-->
<template>
    <div class="file-upload-wrap">
        <!-- picture-card 模式 -->
        <div v-if="listType === 'picture-card'" class="picture-list">
            <!-- 已上传预览 -->
            <div v-for="(file, i) in modelValue" :key="i" class="picture-item">
                <img v-if="isImage(file)" :src="file.url" :alt="file.name" />
                <div v-else class="picture-fallback">
                    <i class="i-fe:file text-28" />
                </div>
                <div class="picture-mask">
                    <i class="i-fe:eye text-18 cursor-pointer" @click="handlePreview(file)" />
                    <i class="i-material-symbols:delete-outline text-18 cursor-pointer" @click="handleRemove(i)" />
                </div>
            </div>

            <!-- 上传触发按钮 -->
            <div v-if="!maxCount || modelValue.length < maxCount" class="picture-add" @click="triggerInput">
                <n-spin v-if="uploading" :size="18" />
                <template v-else>
                    <i class="i-fe:plus text-24" />
                    <span class="text-12 mt-4">上传</span>
                </template>
            </div>
        </div>

        <!-- text 模式 -->
        <div v-else class="text-list">
            <div v-for="(file, i) in modelValue" :key="i" class="text-item">
                <i class="i-fe:paperclip text-14 mr-6" />
                <a :href="file.url" target="_blank" class="file-link">{{ file.name }}</a>
                <i class="i-material-symbols:delete-outline text-14 ml-auto cursor-pointer text-red"
                    @click="handleRemove(i)" />
            </div>
            <n-button size="small" :loading="uploading" @click="triggerInput">
                <i class="i-fe:upload mr-4" />选择文件
            </n-button>
        </div>

        <!-- 进度条 -->
        <n-progress v-if="uploading && progress > 0" :percentage="progress" class="mt-8" />

        <input ref="inputRef" type="file" :multiple="multiple" :accept="accept" style="display: none"
            @change="handleChange" />
    </div>
</template>

<script setup>
import { request } from '@/utils'

const props = defineProps({
    modelValue: { type: Array, default: () => [] },
    accept: { type: String, default: '*' },
    maxSizeMB: { type: Number, default: 10 },
    multiple: { type: Boolean, default: true },
    listType: { type: String, default: 'picture-card' }, // 'picture-card' | 'text'
    maxCount: { type: Number, default: 0 }, // 0=不限
})

const emit = defineEmits(['update:modelValue', 'success', 'error'])

const inputRef = ref(null)
const uploading = ref(false)
const progress = ref(0)

function triggerInput() { inputRef.value?.click() }

async function handleChange(e) {
    const files = [...e.target.files]
    e.target.value = ''
    if (!files.length) return

    // 客户端检验
    for (const f of files) {
        if (props.maxSizeMB > 0 && f.size > props.maxSizeMB * 1024 * 1024) {
            $message.error(`文件 ${f.name} 超过最大限制 ${props.maxSizeMB}MB`)
            return
        }
    }

    uploading.value = true
    progress.value = 0
    const newList = [...props.modelValue]

    try {
        for (const f of files) {
            const fd = new FormData()
            fd.append('file', f)
            fd.append('groupId', 0) // 默认传未分组
            const res = await request.post('/attachment/upload', fd, {
                headers: { 'Content-Type': 'multipart/form-data' },
                onUploadProgress: (e) => {
                    progress.value = Math.round((e.loaded / e.total) * 100)
                },
            })
            const att = res.data
            newList.push({ id: att.id, url: att.url, name: att.originalName, mimeType: att.mimeType })
            emit('success', att)
        }
        emit('update:modelValue', newList)
        $message.success('上传成功')
    }
    catch (err) {
        emit('error', err)
        $message.error('上传失败')
    }
    finally {
        uploading.value = false
        progress.value = 0
    }
}

function handleRemove(index) {
    const list = [...props.modelValue]
    list.splice(index, 1)
    emit('update:modelValue', list)
}

function handlePreview(file) {
    if (isImage(file)) window.open(file.url, '_blank')
}

function isImage(file) {
    return file.mimeType?.startsWith('image/') || /\.(png|jpg|jpeg|gif|webp|svg)$/i.test(file.url)
}
</script>

<style scoped>
.file-upload-wrap {
    width: 100%;
}

/* picture-card */
.picture-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
}

.picture-item,
.picture-add {
    width: 88px;
    height: 88px;
    border-radius: 8px;
    border: 1.5px dashed var(--n-border-color, #d9d9d9);
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    overflow: hidden;
}

.picture-item img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.picture-fallback {
    color: #aaa;
}

.picture-mask {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.45);
    color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    opacity: 0;
    transition: opacity 0.2s;
}

.picture-item:hover .picture-mask {
    opacity: 1;
}

.picture-add {
    cursor: pointer;
    flex-direction: column;
    color: #aaa;
    transition: border-color 0.2s, color 0.2s;
}

.picture-add:hover {
    border-color: rgb(var(--primary-color));
    color: rgb(var(--primary-color));
}

/* text */
.text-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.text-item {
    display: flex;
    align-items: center;
    padding: 4px 8px;
    background: rgba(0, 0, 0, 0.03);
    border-radius: 6px;
    font-size: 13px;
}

.file-link {
    color: rgb(var(--primary-color));
    text-decoration: none;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 280px;
}
</style>
