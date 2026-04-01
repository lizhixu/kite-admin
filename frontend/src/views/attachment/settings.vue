<template>
    <CommonPage back>
        <template #title-suffix>
            <n-tag class="ml-12" type="info" @click="$router.push('/attachment')" style="cursor:pointer">
                附件管理
            </n-tag>
        </template>
        <template #action>
            <n-button type="primary" :loading="saving" @click="handleSave">
                <i class="i-fe:save mr-4" />保存配置
            </n-button>
        </template>

        <n-spin :show="loading">
            <n-form ref="formRef" :model="form" label-placement="left" label-align="right" :label-width="130"
                :style="{ maxWidth: '640px' }">
                <!-- 存储类型 -->
                <n-divider title-placement="left">存储位置</n-divider>
                <n-form-item label="存储类型" path="storageType">
                    <n-radio-group v-model:value="form.storageType">
                        <n-radio value="local">
                            <div class="radio-card" :class="{ active: form.storageType === 'local' }">
                                <i class="i-fe:hard-drive text-20 mb-4" />
                                <div class="font-medium">本地存储</div>
                                <div class="text-12 text-gray-400">保存到服务器本地磁盘</div>
                            </div>
                        </n-radio>
                        <n-radio value="s3">
                            <div class="radio-card" :class="{ active: form.storageType === 's3' }">
                                <i class="i-fe:cloud text-20 mb-4" />
                                <div class="font-medium">S3 对象存储</div>
                                <div class="text-12 text-gray-400">AWS S3 / 兼容接口</div>
                            </div>
                        </n-radio>
                    </n-radio-group>
                </n-form-item>

                <!-- 本地存储配置 -->
                <template v-if="form.storageType === 'local'">
                    <n-form-item label="存储目录" path="localPath">
                        <n-input v-model:value="form.localPath" placeholder="uploads" />
                        <template #feedback>
                            相对于服务器工作目录的路径，默认 uploads
                        </template>
                    </n-form-item>
                </template>

                <!-- S3 配置 -->
                <template v-if="form.storageType === 's3'">
                    <n-form-item label="Endpoint" path="s3Endpoint" :rule="{ required: true, message: '请填写 Endpoint' }">
                        <n-input v-model:value="form.s3Endpoint" placeholder="https://s3.amazonaws.com" />
                    </n-form-item>
                    <n-form-item label="Bucket" path="s3Bucket" :rule="{ required: true, message: '请填写 Bucket' }">
                        <n-input v-model:value="form.s3Bucket" placeholder="my-bucket" />
                    </n-form-item>
                    <n-form-item label="Region" path="s3Region">
                        <n-input v-model:value="form.s3Region" placeholder="us-east-1" />
                    </n-form-item>
                    <n-form-item label="Access Key" path="s3AccessKey">
                        <n-input v-model:value="form.s3AccessKey" />
                    </n-form-item>
                    <n-form-item label="Secret Key" path="s3SecretKey">
                        <n-input v-model:value="form.s3SecretKey" type="password" show-password-on="click" />
                    </n-form-item>
                    <n-form-item label="公开访问 URL" path="s3PublicUrl">
                        <n-input v-model:value="form.s3PublicUrl"
                            placeholder="https://cdn.example.com（不填则使用 Endpoint/Bucket）" />
                    </n-form-item>
                </template>

                <!-- 存储规则 -->
                <n-divider title-placement="left">存储规则</n-divider>
                <n-form-item label="最大文件大小" path="maxSizeMB">
                    <n-input-number v-model:value="form.maxSizeMB" :min="0" :max="2048" style="width: 160px" />
                    <span class="ml-8 text-gray-400">MB（0 = 不限制）</span>
                </n-form-item>
                <n-form-item label="允许的文件类型" path="allowedTypes">
                    <div class="w-full">
                        <div class="flex flex-wrap gap-8 mb-8">
                            <n-tag v-for="t in quickTypes" :key="t.value" checkable
                                :checked="selectedQuickTypes.includes(t.value)"
                                @update:checked="(v) => toggleQuickType(t.value, v)">
                                {{ t.label }}
                            </n-tag>
                        </div>
                        <n-input v-model:value="form.allowedTypes"
                            placeholder="或手动填写 MIME 前缀，逗号分隔，如: image/,application/pdf（空=不限）" clearable />
                    </div>
                </n-form-item>
            </n-form>
        </n-spin>
    </CommonPage>
</template>

<script setup>
import api from './api'

defineOptions({ name: 'AttachmentSettings' })

const formRef = ref(null)
const loading = ref(false)
const saving = ref(false)
const form = ref({
    storageType: 'local',
    localPath: 'uploads',
    maxSizeMB: 10,
    allowedTypes: '',
    s3Endpoint: '',
    s3Bucket: '',
    s3Region: '',
    s3AccessKey: '',
    s3SecretKey: '',
    s3PublicUrl: '',
})

const quickTypes = [
    { label: '图片', value: 'image/' },
    { label: '视频', value: 'video/' },
    { label: '音频', value: 'audio/' },
    { label: 'PDF', value: 'application/pdf' },
    { label: 'Office', value: 'application/vnd' },
    { label: '压缩包', value: 'application/zip' },
]

const selectedQuickTypes = computed(() => {
    const types = form.value.allowedTypes?.split(',').map(s => s.trim()).filter(Boolean) || []
    return quickTypes.filter(q => types.includes(q.value)).map(q => q.value)
})

function toggleQuickType(value, checked) {
    const current = form.value.allowedTypes?.split(',').map(s => s.trim()).filter(Boolean) || []
    const next = checked ? [...new Set([...current, value])] : current.filter(t => t !== value)
    form.value.allowedTypes = next.join(',')
}

async function loadConfig() {
    loading.value = true
    try {
        const res = await api.getConfig()
        if (res.data) Object.assign(form.value, res.data)
    }
    finally {
        loading.value = false
    }
}

async function handleSave() {
    saving.value = true
    try {
        await api.saveConfig(form.value)
        $message.success('配置已保存')
    }
    finally {
        saving.value = false
    }
}

onMounted(loadConfig)
</script>

<style scoped>
.n-radio-group {
    display: flex;
    gap: 12px;
}

.radio-card {
    border: 1.5px solid var(--n-border-color, #e0e0e6);
    border-radius: 10px;
    padding: 14px 24px;
    cursor: pointer;
    text-align: center;
    min-width: 140px;
    transition: border-color 0.2s, background 0.2s;
}

.radio-card.active {
    border-color: rgb(var(--primary-color));
    background: rgba(var(--primary-color), 0.06);
}
</style>
