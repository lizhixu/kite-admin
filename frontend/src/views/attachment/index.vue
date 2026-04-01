<template>
    <CommonPage>
        <div class="flex h-full gap-16">
            <!-- 左侧分组栏 -->
            <div
                class="w-220 flex-shrink-0 flex flex-col bg-white dark:bg-[#18181c] rounded-8 p-12 border border-gray-200 dark:border-gray-800">
                <div class="flex items-center justify-between mb-12 px-8">
                    <h3 class="font-bold text-16 m-0">附件分组</h3>
                    <n-button size="tiny" text type="primary" @click="handleEditGroup()">
                        <i class="i-fe:plus text-16" />
                    </n-button>
                </div>

                <div class="flex-1 overflow-y-auto cus-scroll pr-4">
                    <!-- 默认分组 -->
                    <div class="group-item" :class="{ active: activeGroupId === '' }" @click="activeGroupId = ''">
                        <i class="i-fe:layers text-16 mr-8" />
                        <span class="flex-1">全部附件</span>
                    </div>
                    <div class="group-item" :class="{ active: activeGroupId === 0 }" @click="activeGroupId = 0">
                        <i class="i-fe:folder text-16 mr-8" />
                        <span class="flex-1">未分组</span>
                    </div>

                    <n-divider class="my-8" />

                    <!-- 自定义分组 -->
                    <div v-for="g in groups" :key="g.id" class="group-item" :class="{ active: activeGroupId === g.id }"
                        @click="activeGroupId = g.id">
                        <i class="i-fe:folder text-16 mr-8" />
                        <span class="flex-1 truncate" :title="g.name">{{ g.name }}</span>
                        <n-dropdown trigger="hover" :options="[
                            { label: '编辑名称', key: 'edit' },
                            { label: '删除分组', key: 'delete' },
                        ]" @select="(key) => handleGroupAction(key, g)">
                            <i class="i-fe:more-vertical text-16 opacity-50 hover:opacity-100" @click.stop />
                        </n-dropdown>
                    </div>
                </div>
            </div>

            <!-- 右侧附件区 -->
            <div
                class="w-0 flex-1 flex flex-col bg-white dark:bg-[#18181c] rounded-8 p-16 border border-gray-200 dark:border-gray-800">
                <!-- 顶部工具栏 -->
                <div class="mb-16 flex items-center justify-between gap-12">
                    <div class="flex items-center gap-8">
                        <n-input v-model:value="keyword" placeholder="搜索文件名…" clearable style="width: 220px"
                            @keyup.enter="loadList">
                            <template #prefix>
                                <i class="i-fe:search text-14" />
                            </template>
                        </n-input>
                        <n-button @click="loadList">搜索</n-button>
                    </div>
                    <div class="flex gap-8">
                        <n-button type="primary" @click="triggerUpload">
                            <i class="i-fe:upload mr-4" />上传文件
                        </n-button>
                        <n-button v-permission="'ManageAttachmentConfig'" @click="$router.push('/attachment/settings')">
                            <i class="i-fe:settings mr-4" />存储设置
                        </n-button>
                    </div>
                </div>

                <!-- 隐藏的 input -->
                <input ref="fileInputRef" type="file" multiple style="display: none" @change="handleFileChange" />

                <!-- 文件网格 -->
                <div class="flex-1 overflow-y-auto pr-4 cus-scroll relative">
                    <n-spin :show="loading" class="min-h-[200px]">
                        <div v-if="list.length" class="attachment-grid">
                            <div v-for="item in list" :key="item.id" class="attachment-card">
                                <!-- 预览区域 -->
                                <div class="attachment-preview" @click="handlePreview(item)">
                                    <img v-if="isImage(item.mimeType)" :src="item.url" :alt="item.originalName" />
                                    <div v-else class="file-icon">
                                        <i :class="getFileIcon(item.mimeType)" />
                                    </div>
                                </div>
                                <!-- 文件信息 -->
                                <div class="attachment-info">
                                    <div class="file-name" :title="item.originalName">{{ item.originalName }}</div>
                                    <div class="file-meta">
                                        <span>{{ formatSize(item.fileSize) }}</span>
                                        <span>{{ formatDate(item.createTime) }}</span>
                                    </div>
                                    <div class="file-uploader mt-4 text-11 text-gray-400 flex items-center">
                                        <i class="i-fe:user mr-4" />
                                        <span>{{ item.uploader?.username || '未知' }}</span>
                                    </div>
                                </div>
                                <!-- 操作 -->
                                <div class="attachment-actions">
                                    <n-button size="tiny" text @click="copyUrl(item.url)" title="复制链接">
                                        <i class="i-fe:link text-14" />
                                    </n-button>
                                    <n-button size="tiny" text @click="downloadFile(item)" title="下载">
                                        <i class="i-fe:download text-14" />
                                    </n-button>
                                    <n-button v-permission="'DeleteAttachment'" size="tiny" text type="error"
                                        @click="handleDelete(item)" title="删除">
                                        <i class="i-material-symbols:delete-outline text-14" />
                                    </n-button>
                                </div>
                            </div>
                        </div>
                        <n-empty v-else-if="!loading" description="该分组暂无附件" class="py-60" />
                    </n-spin>
                </div>

                <!-- 分页 -->
                <div v-if="total > 0" class="mt-16 flex justify-end flex-shrink-0">
                    <n-pagination v-model:page="page" :page-size="pageSize" :item-count="total"
                        :page-sizes="[20, 40, 60]" show-size-picker @update:page="loadList"
                        @update:page-size="(s) => { pageSize = s; page = 1; loadList() }" />
                </div>
            </div>
        </div>

        <!-- 图片预览 -->
        <n-image-group>
            <n-image ref="previewRef" v-show="false" :src="previewUrl" />
        </n-image-group>

        <!-- 分组弹窗 -->
        <n-modal v-model:show="groupModalVis" preset="card" :title="groupForm.id ? '编辑分组' : '新建分组'"
            style="width: 400px">
            <n-form ref="groupFormRef" :model="groupForm" label-placement="left" label-width="80">
                <n-form-item label="分组名称" path="name" :rule="{ required: true, message: '请输入分组名称' }">
                    <n-input v-model:value="groupForm.name" placeholder="分组名称" />
                </n-form-item>
                <n-form-item label="排序" path="order">
                    <n-input-number v-model:value="groupForm.order" />
                </n-form-item>
            </n-form>
            <div class="flex justify-end mt-16">
                <n-button class="mr-12" @click="groupModalVis = false">取消</n-button>
                <n-button type="primary" :loading="groupSaving" @click="submitGroup">保存</n-button>
            </div>
        </n-modal>
    </CommonPage>
</template>

<script setup>
import { ref, watch, onMounted, nextTick } from 'vue'
import api from './api'

defineOptions({ name: 'AttachmentMgt' })

const fileInputRef = ref(null)

// ---------------- 分组逻辑 ----------------
const groups = ref([])
const activeGroupId = ref('') // '' 为全部， 0 为未分组，其余为真实 ID

const groupModalVis = ref(false)
const groupFormRef = ref(null)
const groupSaving = ref(false)
const groupForm = ref({ name: '', order: 0 })

async function loadGroups() {
    try {
        const res = await api.getGroups()
        groups.value = res.data || []
    }
    catch (e) {
        console.error(e)
    }
}

function handleEditGroup(g = null) {
    if (g) {
        groupForm.value = { ...g }
    }
    else {
        groupForm.value = { name: '', order: 0 }
    }
    groupModalVis.value = true
}

function submitGroup() {
    groupFormRef.value?.validate(async (errors) => {
        if (!errors) {
            groupSaving.value = true
            try {
                if (groupForm.value.id) {
                    await api.updateGroup(groupForm.value.id, groupForm.value)
                }
                else {
                    await api.createGroup(groupForm.value)
                }
                $message.success('保存成功')
                groupModalVis.value = false
                loadGroups()
            }
            finally {
                groupSaving.value = false
            }
        }
    })
}

function handleGroupAction(key, g) {
    if (key === 'edit') {
        handleEditGroup(g)
    }
    else if (key === 'delete') {
        $dialog.warning({
            title: '删除分组',
            content: `确定要删除分组 [${g.name}] 吗？组内的附件不会被删除，会被移动到未分组。`,
            positiveText: '删除',
            onPositiveClick: async () => {
                await api.deleteGroup(g.id)
                $message.success('已删除')
                if (activeGroupId.value === g.id) {
                    activeGroupId.value = ''
                }
                loadGroups()
            },
        })
    }
}

watch(activeGroupId, () => {
    page.value = 1
    loadList()
})

// ---------------- 附件列表逻辑 ----------------
const list = ref([])
const total = ref(0)
const loading = ref(false)
const keyword = ref('')
const page = ref(1)
const pageSize = ref(20)

async function loadList() {
    loading.value = true
    try {
        const params = { keyword: keyword.value, page: page.value, pageSize: pageSize.value }
        if (activeGroupId.value !== '') {
            params.groupId = activeGroupId.value
        }
        const res = await api.getList(params)
        list.value = res.data?.pageData || []
        total.value = res.data?.total || 0
    }
    finally {
        loading.value = false
    }
}

onMounted(() => {
    loadGroups()
    loadList()
})

// 上传
function triggerUpload() { fileInputRef.value?.click() }
async function handleFileChange(e) {
    const files = [...e.target.files]
    if (!files.length) return
    e.target.value = ''

    $loadingBar.start()
    try {
        for (const f of files) {
            const formData = new FormData()
            formData.append('file', f)
            // 带上当前分组 ID（如果是全部即''，转为设未分组 0）
            formData.append('groupId', activeGroupId.value === '' ? 0 : activeGroupId.value)
            await api.upload(formData)
        }
        $message.success('上传成功')
        loadList()
    }
    catch {
        $message.error('部分文件上传失败')
    }
    finally {
        $loadingBar.finish()
    }
}

// 删除
function handleDelete(item) {
    $dialog.warning({
        title: '提示',
        content: `确定删除 ${item.originalName}？`,
        positiveText: '删除',
        negativeText: '取消',
        onPositiveClick: async () => {
            await api.delete(item.id)
            $message.success('已删除')
            loadList()
        },
    })
}

// 预览
const previewUrl = ref('')
const previewRef = ref(null)
function handlePreview(item) {
    if (isImage(item.mimeType)) {
        previewUrl.value = item.url
        nextTick(() => previewRef.value?.click())
    }
    else {
        window.open(item.url, '_blank')
    }
}

// 工具函数
function isImage(mime) { return mime?.startsWith('image/') || /\.(png|jpg|jpeg|gif|webp|svg)$/i.test(previewUrl.value) }
function getFileIcon(mime) {
    if (!mime) return 'i-fe:file text-36'
    if (mime.startsWith('video/')) return 'i-fe:video text-36'
    if (mime.startsWith('audio/')) return 'i-fe:music text-36'
    if (mime.includes('pdf')) return 'i-fe:file-text text-36'
    if (mime.includes('zip') || mime.includes('rar')) return 'i-fe:archive text-36'
    return 'i-fe:file text-36'
}
function formatSize(bytes) {
    if (!bytes) return '0'
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}
function formatDate(t) {
    if (!t) return ''
    return new Date(t).toLocaleDateString('zh-CN')
}
function copyUrl(url) {
    navigator.clipboard.writeText(url).then(() => $message.success('链接已复制'))
}
function downloadFile(item) {
    const a = document.createElement('a')
    a.href = item.url
    a.download = item.originalName
    a.click()
}
</script>

<style scoped>
.group-item {
    display: flex;
    align-items: center;
    padding: 10px 12px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s, color 0.2s;
    margin-bottom: 4px;
}

.group-item:hover {
    background: rgba(var(--primary-color), 0.05);
}

.group-item.active {
    background: rgba(var(--primary-color), 0.1);
    color: rgb(var(--primary-color));
    font-weight: 500;
}

.attachment-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(170px, 1fr));
    gap: 14px;
}

.attachment-card {
    border: 1px solid var(--n-border-color, #e0e0e6);
    border-radius: 10px;
    overflow: hidden;
    transition: box-shadow 0.2s, border-color 0.2s;
    display: flex;
    flex-direction: column;
}

.attachment-card:hover {
    border-color: rgba(var(--primary-color), 0.5);
    box-shadow: 0 4px 14px rgba(var(--primary-color), 0.1);
}

.attachment-preview {
    height: 120px;
    background: #f5f5fa;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    overflow: hidden;
}

.dark .attachment-preview {
    background: #232328;
}

.attachment-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.2s;
}

.attachment-card:hover .attachment-preview img {
    transform: scale(1.04);
}

.file-icon {
    font-size: 40px;
    color: #aaa;
}

.attachment-info {
    padding: 10px 12px 6px;
    flex: 1;
}

.file-name {
    font-size: 13px;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.file-meta {
    display: flex;
    justify-content: space-between;
    font-size: 11px;
    color: #aaa;
    margin-top: 4px;
}

.attachment-actions {
    display: flex;
    justify-content: flex-end;
    gap: 2px;
    padding: 4px 8px 8px;
    opacity: 0;
    transition: opacity 0.15s;
}

.attachment-card:hover .attachment-actions {
    opacity: 1;
}
</style>
