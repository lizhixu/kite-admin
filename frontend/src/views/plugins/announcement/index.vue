<template>
    <CommonPage>
        <div class="mb-16 flex justify-between items-center">
            <div>
                <h3 class="font-bold">系统公告 (插件演示)</h3>
                <p class="text-12 text-gray-400">向全平台用户发布重要通知</p>
            </div>
            <n-button type="primary" size="small" @click="handleCreate">发布公告</n-button>
        </div>

        <n-alert v-if="topNotice" type="warning" title="置顶公告" class="mb-16">
            <p class="font-bold text-16 mb-8">{{ topNotice.title }}</p>
            <p>{{ topNotice.content }}</p>
            <template #icon>
                <i class="i-fe:bell" />
            </template>
        </n-alert>

        <div class="space-y-16">
            <n-card v-for="item in commonNotices" :key="item.id" hoverable>
                <div class="flex justify-between items-start">
                    <div>
                        <div class="flex items-center gap-x-8 mb-8">
                            <n-tag :type="getTypeColor(item.type)" size="small">{{ item.type }}</n-tag>
                            <h4 class="font-bold text-16 m-0">{{ item.title }}</h4>
                        </div>
                        <p class="text-gray-500 whitespace-pre-wrap">{{ item.content }}</p>
                        <p class="text-12 text-gray-400 mt-12">{{ formatTime(item.createTime) }}</p>
                    </div>
                    <n-button size="small" type="error" text @click="handleDelete(item)">删除</n-button>
                </div>
            </n-card>
        </div>

        <n-empty v-if="commonNotices.length === 0 && !topNotice" description="暂无历史公告" class="my-32" />

        <!-- 发布弹窗 -->
        <n-modal v-model:show="modalVisible" preset="card" title="发布新公告" style="width: 500px">
            <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80">
                <n-form-item label="公告标题" path="title">
                    <n-input v-model:value="form.title" placeholder="请输入标题" />
                </n-form-item>
                <n-form-item label="通知类型" path="type">
                    <n-select v-model:value="form.type" :options="[
                        { label: '普通通知', value: 'INFO' },
                        { label: '重要警告', value: 'WARNING' },
                        { label: '紧急更新', value: 'URGENT' }
                    ]" />
                </n-form-item>
                <n-form-item label="设为置顶" path="isTop">
                    <n-switch v-model:value="form.isTop" />
                </n-form-item>
                <n-form-item label="公告详情" path="content">
                    <n-input v-model:value="form.content" type="textarea" :rows="4" placeholder="请输入详细内容" />
                </n-form-item>
                <div class="flex justify-end mt-16">
                    <n-button class="mr-12" @click="modalVisible = false">取消</n-button>
                    <n-button type="primary" :loading="submitLoading" @click="submitCreate">发布</n-button>
                </div>
            </n-form>
        </n-modal>
    </CommonPage>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from './api'

const notices = ref([])
const topNotice = computed(() => notices.value.find(n => n.isTop))
const commonNotices = computed(() => notices.value.filter(n => !n.isTop))

async function loadData() {
    try {
        const res = await api.getList()
        notices.value = res.data || []
    } catch (error) {
        console.error(error)
    }
}

// 表单
const modalVisible = ref(false)
const formRef = ref()
const submitLoading = ref(false)
const form = ref({ title: '', content: '', type: 'INFO', isTop: false })
const rules = { title: { required: true, message: '请输入标题' } }

function handleCreate() {
    form.value = { title: '', content: '', type: 'INFO', isTop: false }
    modalVisible.value = true
}

function submitCreate() {
    formRef.value?.validate(async (errors) => {
        if (!errors) {
            try {
                submitLoading.value = true
                await api.create(form.value)
                $message.success('发布成功')
                modalVisible.value = false
                loadData()
            } catch (error) {
                console.error(error)
            } finally {
                submitLoading.value = false
            }
        }
    })
}

function handleDelete(item) {
    $dialog.warning({
        title: '提示',
        content: `确定删除公告 ${item.title} 吗？`,
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: async () => {
            await api.delete(item.id)
            $message.success('已删除')
            loadData()
        }
    })
}

function getTypeColor(type) {
    const map = {
        'INFO': 'info',
        'WARNING': 'warning',
        'URGENT': 'error'
    }
    return map[type] || 'default'
}

function formatTime(t) {
    if (!t) return ''
    const date = new Date(t)
    return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
}

onMounted(() => {
    loadData()
})
</script>
