<template>
  <CommonPage>
    <template #action>
      <NSpace>
        <NButton type="primary" :loading="saving" @click="handleSave">
          保存配置
        </NButton>
        <NButton :loading="testing" :disabled="!form.enabled" @click="handleTest">
          发送测试邮件
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
        <n-form-item label="SMTP 服务器" path="host">
          <n-input v-model:value="form.host" placeholder="例如: smtp.qq.com" />
        </n-form-item>
        <n-form-item label="端口" path="port">
          <n-input-number v-model:value="form.port" :min="1" :max="65535" style="width: 200px" />
        </n-form-item>
        <n-form-item label="用户名" path="username">
          <n-input v-model:value="form.username" placeholder="邮箱账号" />
        </n-form-item>
        <n-form-item label="密码/授权码" path="password">
          <n-input
            v-model:value="form.password"
            type="password"
            show-password-on="click"
            placeholder="邮箱密码或授权码"
          />
        </n-form-item>
        <n-form-item label="发件人名称" path="fromName">
          <n-input v-model:value="form.fromName" placeholder="例如: Kite Admin" />
        </n-form-item>
        <n-form-item label="发件人邮箱" path="fromEmail">
          <n-input v-model:value="form.fromEmail" placeholder="例如: noreply@example.com" />
        </n-form-item>
        <n-form-item label="启用">
          <n-switch v-model:value="form.enabled" />
        </n-form-item>
      </n-form>
    </n-spin>
  </CommonPage>
</template>

<script setup>
import { NButton, NSpace, NSpin } from 'naive-ui'
import api from './api'

defineOptions({ name: 'EmailConfig' })

const loading = ref(false)
const saving = ref(false)
const testing = ref(false)
const formRef = ref(null)

const form = ref({
  host: '',
  port: 587,
  username: '',
  password: '',
  fromName: '',
  fromEmail: '',
  enabled: false,
})

async function loadConfig() {
  loading.value = true
  try {
    const { data } = await api.getEmailConfig()
    if (data) {
      Object.assign(form.value, data)
    }
  } catch { /* ignore */ }
  loading.value = false
}

async function handleSave() {
  saving.value = true
  try {
    await api.saveEmailConfig(form.value)
    $message.success('保存成功')
  } catch { /* ignore */ }
  saving.value = false
}

async function handleTest() {
  testing.value = true
  try {
    await api.testEmailConfig()
    $message.success('测试邮件已发送，请检查收件箱')
  } catch { /* ignore */ }
  testing.value = false
}

onMounted(loadConfig)
</script>
