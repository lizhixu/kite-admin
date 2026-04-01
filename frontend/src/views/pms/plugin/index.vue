<template>
    <CommonPage>
        <div class="mb-16">
            <h3 class="font-bold">插件管理</h3>
            <p class="text-gray-400 text-12">管理系统的功能扩展插件，安装后刷新页面即可生效侧边栏菜单。</p>
        </div>

        <n-spin :show="loading">
            <n-grid :x-gap="16" :y-gap="16" cols="1 s:2 m:3 l:4 xl:5 2xl:6" responsive="screen">
                <n-gi v-for="item in plugins" :key="item.code">
                    <n-card hoverable class="h-full plugin-card" size="small">
                        <template #header>
                            <div class="flex items-center justify-between">
                                <span class="font-bold text-16 truncate" :title="item.name">{{ item.name }}</span>
                                <n-tag :type="item.isInstalled ? 'success' : 'default'" size="small">
                                    {{ item.isInstalled ? '已安装' : '未安装' }}
                                </n-tag>
                            </div>
                        </template>

                        <div class="text-12 text-gray-400 mb-12 flex space-x-12">
                            <span>版本: {{ item.version }}</span>
                            <span>作者: {{ item.author }}</span>
                        </div>

                        <p class="text-14 mb-16 h-40 overflow-hidden text-ellipsis line-clamp-2"
                            :title="item.description">
                            {{ item.description || '暂无描述' }}
                        </p>

                        <template #footer>
                            <div class="flex items-center justify-between">
                                <div v-if="item.isInstalled" class="flex items-center">
                                    <span class="text-12 mr-8">状态：</span>
                                    <n-switch v-model:value="item.enable" size="small" :loading="item.enableLoading"
                                        @update:value="handleToggleEnable(item)" />
                                </div>
                                <div v-else></div>

                                <div class="flex space-x-8">
                                    <n-button v-if="!item.isInstalled" type="primary" size="small"
                                        :loading="item.actionLoading" @click="handleInstall(item)">
                                        安装
                                    </n-button>
                                    <n-button v-if="item.isInstalled" type="error" secondary size="small"
                                        :loading="item.actionLoading" @click="handleUninstall(item)">
                                        卸载
                                    </n-button>
                                </div>
                            </div>
                        </template>
                    </n-card>
                </n-gi>
            </n-grid>

            <n-empty v-if="!loading && plugins.length === 0" description="暂无可用插件" class="my-32" />
        </n-spin>
    </CommonPage>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from './api'



const loading = ref(false)
const plugins = ref([])

async function loadData() {
    try {
        loading.value = true
        const res = await api.getList()
        plugins.value = res.data || []
    } catch (error) {
        console.error(error)
    } finally {
        loading.value = false
    }
}

async function handleInstall(item) {
    $dialog.info({
        title: '确认安装',
        content: `确定要安装插件「${item.name}」吗？安装后会初始化数据结构。`,
        positiveText: '确认',
        negativeText: '取消',
        onPositiveClick: async () => {
            try {
                item.actionLoading = true
                await api.install(item.code)
                $message.success('安装成功，需要刷新页面重新加载菜单')
                setTimeout(() => window.location.reload(), 1500)
            } catch (error) {
                console.error(error)
            } finally {
                item.actionLoading = false
                loadData()
            }
        }
    })
}

async function handleUninstall(item) {
    $dialog.warning({
        title: '确认卸载',
        content: `警告：卸载插件「${item.name}」将会清除所有相关业务数据且无法恢复。确定要继续吗？`,
        positiveText: '确认卸载',
        negativeText: '取消',
        onPositiveClick: async () => {
            try {
                item.actionLoading = true
                await api.uninstall(item.code)
                $message.success('卸载成功，需要刷新页面重新加载菜单')
                setTimeout(() => window.location.reload(), 1500)
            } catch (error) {
                console.error(error)
            } finally {
                item.actionLoading = false
                loadData()
            }
        }
    })
}

async function handleToggleEnable(item) {
    try {
        item.enableLoading = true
        const targetStatus = item.enable
        // reverse status briefly to await api call
        item.enable = !targetStatus
        await api.toggleEnable(item.code, targetStatus)
        item.enable = targetStatus
        $message.success(`${targetStatus ? '启用' : '禁用'}成功，刷新页面后生效菜单变更`)
    } catch (error) {
        console.error(error)
    } finally {
        item.enableLoading = false
    }
}

onMounted(() => {
    loadData()
})
</script>

<style scoped>
.plugin-card {
    transition: all 0.3s ease;
}

.plugin-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>
