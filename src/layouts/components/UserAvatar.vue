<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:50:42
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <n-dropdown :options="options" @select="handleSelect">
    <div id="user-dropdown" class="flex cursor-pointer items-center">
      <n-avatar round :size="36" :src="userStore.avatar" />
      <div v-if="userStore.userInfo" class="ml-12 flex-col flex-shrink-0 items-center">
        <span class="text-14">{{ userStore.nickName ?? userStore.username }}</span>
        <span class="text-12 opacity-50">[{{ userStore.currentRole?.name }}]</span>
      </div>
    </div>
  </n-dropdown>

  <RoleSelect ref="roleSelectRef" />
</template>

<script setup>
import api from '@/api'
import { RoleSelect } from '@/layouts/components'
import { useAuthStore, useUserStore } from '@/store'

const router = useRouter()
const userStore = useUserStore()
const authStore = useAuthStore()

const options = reactive([
  {
    label: '个人资料',
    key: 'profile',
    icon: () => h('i', { class: 'i-fe:user?mask text-14' }),
  },
  {
    label: '切换角色',
    key: 'toggleRole',
    icon: () => h('i', { class: 'i-fe:repeat?mask text-14' }),
    show: computed(() => userStore.roles.length > 1),
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h('i', { class: 'i-fe:log-out?mask text-14' }),
  },
])

const roleSelectRef = ref(null)
function handleSelect(key) {
  switch (key) {
    case 'profile':
      router.push('/profile')
      break
    case 'toggleRole':
      roleSelectRef.value?.open({
        onOk() {
          window.location.href = '/'
        },
      })
      break
    case 'logout':
      $dialog.confirm({
        title: '提示',
        type: 'info',
        content: '确认退出？',
        async confirm() {
          try {
            await api.logout()
          }
          catch (error) {
            console.error(error)
          }
          authStore.logout()
          $message.success('已退出登录')
        },
      })
      break
  }
}
</script>
