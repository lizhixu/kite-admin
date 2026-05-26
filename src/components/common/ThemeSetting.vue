<template>
  <n-tooltip trigger="hover" placement="bottom">
    <template #trigger>
      <div id="theme-setting" class="mr-16 h-32 w-32 f-c-c">
        <n-color-picker
          :value="appStore.primaryColor"
          :swatches="primaryColors"
          :on-update:value="(v) => appStore.setPrimaryColor(v)"
          :render-label="() => ''"
        />
      </div>
    </template>
    设置主题色
  </n-tooltip>
</template>

<script setup>
import { getPresetColors } from '@arco-design/color'
import { useAppStore } from '@/store'

const appStore = useAppStore()

const primaryColors = Object.entries(getPresetColors()).map(([, value]) => value.primary)
</script>
