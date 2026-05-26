<template>
  <CommonPage>
    <template #action>
      <NSpace v-if="selected">
        <NButton v-permission="'SaveEmailTemplate'" type="primary" :loading="saving" @click="handleSave">
          <i class="i-material-symbols:save-outline mr-4" />保存
        </NButton>
        <NButton @click="showPreview = true">
          <i class="i-fe:eye mr-4" />预览
        </NButton>
      </NSpace>
    </template>

    <n-spin :show="loading">
      <div class="template-layout">
        <!-- Left: Template List -->
        <n-card size="small" class="template-sidebar" title="模板列表" :bordered="true">
          <n-list clickable hoverable :show-divider="false">
            <n-list-item
              v-for="tmpl in templates"
              :key="tmpl.id"
              class="template-item"
              :class="{ active: selectedId === tmpl.id }"
              @click="selectTemplate(tmpl)"
            >
              <n-thing>
                <template #header>
                  <div class="flex items-center gap-6">
                    <i class="i-fe:mail text-14" style="color: var(--primary-color);" />
                    <span class="text-13">{{ tmpl.name }}</span>
                  </div>
                </template>
                <template #header-extra>
                  <n-tag v-if="tmpl.isBuiltin" size="tiny" type="success" :bordered="false">内置</n-tag>
                </template>
                <template #description>
                  <span class="text-12" style="opacity: 0.5">{{ tmpl.scene }}</span>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>

        <!-- Right: Editor -->
        <div class="template-editor cus-scroll">
          <template v-if="selected">
            <n-card size="small" :bordered="true">
              <template #header>
                <div class="flex items-center gap-8">
                  <span>{{ selected.name }}</span>
                  <n-tag size="small" :bordered="false" type="info">{{ selected.scene }}</n-tag>
                </div>
              </template>

              <n-form
                ref="formRef"
                label-placement="left"
                label-align="left"
                :label-width="80"
                :model="form"
              >
                <n-form-item label="模板名称" path="name">
                  <n-input v-model:value="form.name" placeholder="模板名称" />
                </n-form-item>
                <n-form-item label="邮件主题" path="subject">
                  <n-input v-model:value="form.subject" placeholder="支持 {{变量}} 占位符" />
                </n-form-item>
                <n-form-item label="模板内容" path="content">
                  <n-input
                    v-model:value="form.content"
                    type="textarea"
                    :autosize="{ minRows: 16, maxRows: 30 }"
                    placeholder="HTML 模板内容，支持 {{变量}} 占位符"
                    class="code-editor"
                  />
                </n-form-item>
              </n-form>
            </n-card>
          </template>

          <n-empty v-else description="请从左侧选择一个模板" class="py-48" />
        </div>

        <!-- Far Right: Variable Reference -->
        <n-card v-if="selected" size="small" class="template-vars" title="可用变量" :bordered="true">
          <n-space vertical :size="8">
            <div v-for="v in variables" :key="v.key" class="var-item">
              <n-code :code="`{{${v.key}}}`" language="text" />
              <span class="var-desc">{{ v.desc }}</span>
            </div>
          </n-space>
        </n-card>
      </div>
    </n-spin>

    <!-- Preview Modal -->
    <n-modal
      v-model:show="showPreview"
      preset="card"
      title="邮件预览"
      style="width: 700px; max-width: 90vw;"
    >
      <n-spin :show="previewing">
        <template v-if="previewData">
          <n-descriptions :column="1" bordered size="small" class="mb-12">
            <n-descriptions-item label="主题">{{ previewData.subject }}</n-descriptions-item>
          </n-descriptions>
          <n-divider style="margin: 0" />
          <div class="preview-body" v-html="previewData.htmlBody" />
        </template>
        <n-empty v-else description="暂无预览数据" />
      </n-spin>
    </n-modal>
  </CommonPage>
</template>

<script setup>
import { NCard, NCode, NDivider, NList, NListItem, NSpace, NSpin, NTag, NThing } from 'naive-ui'
import api from './api'

defineOptions({ name: 'EmailTemplate' })

const loading = ref(false)
const saving = ref(false)
const previewing = ref(false)
const showPreview = ref(false)

const templates = ref([])
const selected = ref(null)
const selectedId = ref(null)
const form = ref({ name: '', subject: '', content: '' })
const previewData = ref(null)

const variables = [
  { key: 'title', desc: '消息/邮件标题' },
  { key: 'content', desc: '内容（已渲染为 HTML）' },
  { key: 'username', desc: '收件人昵称' },
  { key: 'currentTime', desc: '发送时间' },
  { key: 'siteURL', desc: '系统地址' },
]

async function loadTemplates() {
  loading.value = true
  try {
    const { data } = await api.getEmailTemplates()
    templates.value = data || []
    if (templates.value.length > 0 && !selectedId.value) {
      selectTemplate(templates.value[0])
    }
  } catch { /* ignore */ }
  loading.value = false
}

function selectTemplate(tmpl) {
  selected.value = tmpl
  selectedId.value = tmpl.id
  form.value = {
    name: tmpl.name,
    subject: tmpl.subject,
    content: tmpl.content,
  }
}

async function handleSave() {
  if (!selected.value) return
  saving.value = true
  try {
    const { data } = await api.saveEmailTemplate(selected.value.id, form.value)
    if (data) {
      Object.assign(selected.value, data)
    }
    $message.success('保存成功')
  } catch { /* ignore */ }
  saving.value = false
}

watch(showPreview, async (val) => {
  if (val && selected.value) {
    previewing.value = true
    try {
      const { data } = await api.previewEmailTemplate(selected.value.id, {})
      previewData.value = data
    } catch {
      previewData.value = null
    }
    previewing.value = false
  }
})

onMounted(loadTemplates)
</script>

<style scoped>
.template-layout {
  display: flex;
  gap: 16px;
  height: calc(100vh - 200px);
}

.template-sidebar {
  width: 240px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.template-sidebar :deep(.n-card__content) {
  padding: 0;
  flex: 1;
  overflow-y: auto;
}

.template-editor {
  flex: 1;
  min-width: 0;
}

.template-vars {
  width: 220px;
  flex-shrink: 0;
  overflow-y: auto;
}

.template-item {
  padding: 8px 12px;
  transition: background-color 0.2s;
}
.template-item.active {
  background-color: rgba(24, 160, 88, 0.06);
}
.template-item.active :deep(.n-thing-header__title) {
  color: #18a058;
}

.var-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.var-desc {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.4);
  padding-left: 2px;
}
:root.dark .var-desc {
  color: rgba(255, 255, 255, 0.4);
}

.code-editor :deep(textarea) {
  font-family: 'Fira Code', 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  line-height: 1.6;
  tab-size: 2;
}

.preview-body {
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 6px;
  padding: 20px;
  background: rgba(0, 0, 0, 0.02);
  max-height: 500px;
  overflow-y: auto;
  margin-top: 12px;
}
:root.dark .preview-body {
  border-color: rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.04);
}
</style>