// 全局 $picker 服务
//
// 使用方式：
//   const items = await $picker.open({
//     accept: 'image/',        // 可选 mime 前缀，如 'image/'、'video/'
//     multiple: false,         // 是否多选
//     max: 0,                  // multiple=true 时上限，0 表示不限
//     configId: undefined,     // 限定存储；不传则使用默认存储
//     folderPath: 'avatars',   // 指定目标目录路径（自动创建），优先于 folderId
//     folderId: 0,             // 默认进入的文件夹 ID（folderPath 未设置时生效）
//     uploadable: true,        // 是否允许在弹窗内上传新文件
//     title: '选择媒体',
//   })
//   // items: Media[]，取消返回 []
//
// 该组件挂载在 App.vue 内部以继承主题/Pinia 上下文，挂载完成后
// 自动把 open 方法注册到 window.$picker。

export { default as GlobalMediaPicker } from './MediaPickerModal.vue'
