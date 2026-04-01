<template>
    <CommonPage>
        <div class="mb-16 flex justify-between items-center">
            <h3 class="font-bold">任务看板 (插件演示)</h3>
            <n-button type="primary" size="small" @click="handleCreate">
                <template #icon><i class="i-fe:plus" /></template>
                新建任务
            </n-button>
        </div>

        <!-- 简单的看板式布局 -->
        <div class="flex gap-x-16 h-[calc(100vh-200px)]">
            <!-- TODO 列 -->
            <div class="flex-1 bg-gray-50 dark:bg-dark-300 rounded-lg p-16 flex flex-col">
                <h4 class="mb-12 font-bold text-gray-500 flex items-center">
                    <div class="w-8 h-8 rounded-full bg-blue-500 mr-8"></div> 待处理 ({{ todos.length }})
                </h4>
                <div class="flex-1 overflow-y-auto space-y-12 pr-4 cus-scroll">
                    <n-card v-for="task in todos" :key="task.id" size="small"
                        class="shadow-sm cursor-pointer hover:shadow-md transition">
                        <template #header>
                            <span class="text-14 font-bold">{{ task.title }}</span>
                        </template>
                        <p class="text-12 text-gray-400 mb-8">{{ task.content || '暂无描述' }}</p>
                        <div class="flex justify-between items-center mt-8">
                            <span class="text-12 text-gray-400">{{ formatTime(task.createTime) }}</span>
                            <n-button-group size="tiny">
                                <n-button type="success" ghost
                                    @click="api.updateStatus(task.id, 'DOING').then(loadData)">进行中</n-button>
                                <n-button type="error" ghost @click="handleDelete(task)">删除</n-button>
                            </n-button-group>
                        </div>
                    </n-card>
                </div>
            </div>

            <!-- DOING 列 -->
            <div class="flex-1 bg-gray-50 dark:bg-dark-300 rounded-lg p-16 flex flex-col">
                <h4 class="mb-12 font-bold text-gray-500 flex items-center">
                    <div class="w-8 h-8 rounded-full bg-orange-500 mr-8"></div> 进行中 ({{ doings.length }})
                </h4>
                <div class="flex-1 overflow-y-auto space-y-12 pr-4 cus-scroll">
                    <n-card v-for="task in doings" :key="task.id" size="small" class="shadow-sm">
                        <template #header>
                            <span class="text-14 font-bold">{{ task.title }}</span>
                        </template>
                        <div class="flex justify-between items-center mt-8">
                            <span class="text-12 text-gray-400">{{ formatTime(task.createTime) }}</span>
                            <n-button-group size="tiny">
                                <n-button type="info" ghost
                                    @click="api.updateStatus(task.id, 'TODO').then(loadData)">回退</n-button>
                                <n-button type="success" ghost
                                    @click="api.updateStatus(task.id, 'DONE').then(loadData)">完成</n-button>
                            </n-button-group>
                        </div>
                    </n-card>
                </div>
            </div>

            <!-- DONE 列 -->
            <div class="flex-1 bg-gray-50 dark:bg-dark-300 rounded-lg p-16 flex flex-col">
                <h4 class="mb-12 font-bold text-gray-500 flex items-center">
                    <div class="w-8 h-8 rounded-full bg-green-500 mr-8"></div> 已完成 ({{ dones.length }})
                </h4>
                <div class="flex-1 overflow-y-auto space-y-12 pr-4 cus-scroll">
                    <n-card v-for="task in dones" :key="task.id" size="small" class="shadow-sm opacity-60">
                        <template #header>
                            <span class="text-14 font-bold line-through">{{ task.title }}</span>
                        </template>
                        <div class="flex justify-between items-center mt-8">
                            <span class="text-12 text-gray-400">{{ formatTime(task.createTime) }}</span>
                            <n-button size="tiny" type="error" ghost @click="handleDelete(task)">删除</n-button>
                        </div>
                    </n-card>
                </div>
            </div>
        </div>

        <!-- 创建弹窗 -->
        <n-modal v-model:show="modalVisible" preset="card" title="新建任务" style="width: 500px">
            <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80">
                <n-form-item label="任务标题" path="title">
                    <n-input v-model:value="form.title" placeholder="请输入任务标题" />
                </n-form-item>
                <n-form-item label="任务描述" path="content">
                    <n-input v-model:value="form.content" type="textarea" placeholder="请输入任务详情" />
                </n-form-item>
                <div class="flex justify-end mt-16">
                    <n-button class="mr-12" @click="modalVisible = false">取消</n-button>
                    <n-button type="primary" :loading="submitLoading" @click="submitCreate">提交</n-button>
                </div>
            </n-form>
        </n-modal>
    </CommonPage>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from './api'

const tasks = ref([])
const todos = computed(() => tasks.value.filter(t => t.status === 'TODO'))
const doings = computed(() => tasks.value.filter(t => t.status === 'DOING'))
const dones = computed(() => tasks.value.filter(t => t.status === 'DONE'))

async function loadData() {
    try {
        const res = await api.getList()
        tasks.value = res.data || []
    } catch (error) {
        console.error(error)
    }
}

// 表单相关
const modalVisible = ref(false)
const formRef = ref()
const submitLoading = ref(false)
const form = ref({ title: '', content: '' })
const rules = { title: { required: true, message: '请输入标题' } }

function handleCreate() {
    form.value = { title: '', content: '' }
    modalVisible.value = true
}

async function submitCreate() {
    formRef.value?.validate(async (errors) => {
        if (!errors) {
            try {
                submitLoading.value = true
                await api.create(form.value)
                $message.success('创建成功')
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

function handleDelete(task) {
    $dialog.warning({
        title: '提示',
        content: `确定删除任务 ${task.title} 吗？`,
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: async () => {
            await api.delete(task.id)
            $message.success('已删除')
            loadData()
        }
    })
}

function formatTime(t) {
    if (!t) return ''
    const date = new Date(t)
    return `${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
}

onMounted(() => {
    loadData()
})
</script>
