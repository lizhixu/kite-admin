<template>
  <CommonPage>
    <template #action>
      <NSpace>
        <NButton v-permission="'SaveSystemConfig'" type="primary" :loading="saving" @click="handleSave">
          保存配置
        </NButton>
      </NSpace>
    </template>

    <n-spin :show="loading">
      <n-form
        ref="formRef"
        label-placement="left"
        label-align="left"
        :label-width="120"
        :model="form"
        style="max-width: 600px"
      >
        <n-form-item label="后台名称" path="siteName">
          <n-input v-model:value="form.siteName" placeholder="例如: Kite Admin" />
        </n-form-item>
        <n-form-item label="Logo" path="logo">
          <div class="flex flex-col gap-8">
            <div v-if="form.logo" class="flex items-center gap-8">
              <img :src="form.logo" alt="Logo" class="h-40 rounded-4 border border-light_border dark:border-dark_border" @error="handleImageError" />
              <NButton size="small" @click="form.logo = ''">
                清除
              </NButton>
            </div>
            <n-input v-model:value="form.logo" placeholder="Logo 图片 URL" />
          </div>
        </n-form-item>
        <n-form-item label="Favicon" path="favicon">
          <div class="flex flex-col gap-8">
            <div v-if="form.favicon" class="flex items-center gap-8">
              <img :src="form.favicon" alt="Favicon" class="h-24 rounded-4 border border-light_border dark:border-dark_border" @error="handleImageError" />
              <NButton size="small" @click="form.favicon = ''">
                清除
              </NButton>
            </div>
            <n-input v-model:value="form.favicon" placeholder="Favicon 图片 URL" />
          </div>
        </n-form-item>
        <n-form-item label="底部版权" path="copyright">
          <n-input v-model:value="form.copyright" type="textarea" placeholder="例如: © 2024 Kite Admin. All rights reserved." :rows="3" />
        </n-form-item>
      </n-form>
    </n-spin>
  </CommonPage>
</template>

<script setup>
import { NButton, NSpace, NSpin } from 'naive-ui'
import { useSystemConfigStore } from '@/store'
import api from './api'

defineOptions({ name: 'SystemConfig' })

const systemConfigStore = useSystemConfigStore()
const loading = ref(false)
const saving = ref(false)
const formRef = ref(null)

const form = ref({
  siteName: '',
  logo: '',
  favicon: '',
  copyright: '',
})

async function loadConfig() {
  loading.value = true
  try {
    const { data } = await api.getSystemConfig()
    if (data) {
      Object.assign(form.value, data)
    }
  } catch { /* ignore */ }
  loading.value = false
}

async function handleSave() {
  saving.value = true
  try {
    await api.saveSystemConfig(form.value)
    // Refresh global system config store
    await systemConfigStore.fetchConfig()
    $message.success('保存成功')
  } catch { /* ignore */ }
  saving.value = false
}

function handleImageError(e) {
  e.target.style.display = 'none'
}

onMounted(loadConfig)
</script>
