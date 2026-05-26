<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:50:28
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <router-link class="h-60 f-c-c" to="/">
    <img v-if="systemConfigStore.logo" :src="systemConfigStore.logo" class="h-36 w-36 rounded-4 object-contain" @error="handleLogoError">
    <TheLogo v-else />
    <h2
      v-show="!appStore.collapsed"
      class="ml-10 max-w-140 flex-shrink-0 text-16 color-primary font-bold"
    >
      {{ systemConfigStore.siteName }}
    </h2>
  </router-link>
</template>

<script setup>
import { useAppStore, useSystemConfigStore } from '@/store'

const appStore = useAppStore()
const systemConfigStore = useSystemConfigStore()

function handleLogoError(e) {
  e.target.style.display = 'none'
  e.target.nextElementSibling?.style?.removeProperty('display')
}
</script>
