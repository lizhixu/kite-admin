package core

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Plugin 插件必须实现的接口
type Plugin interface {
	// Info 返回插件的基本信息和前端路由注册信息
	Info() PluginInfo

	// Migrate 插件安装时的数据库建表操作
	Migrate(db *gorm.DB) error

	// RegisterRoutes 插件后台 API 路由注册
	RegisterRoutes(group *gin.RouterGroup)

	// Uninstall 插件卸载时的数据库清理操作
	Uninstall(db *gorm.DB) error
}

// PluginInfo 插件元信息
type PluginInfo struct {
	Code           string          `json:"code"`        // 唯一标识，如 task_manager
	Name           string          `json:"name"`        // 显示名称，如 任务管理
	Description    string          `json:"description"` // 描述
	Version        string          `json:"version"`     // 版本
	Author         string          `json:"author"`      // 作者
	MenuGroups     []MenuGroup     `json:"menuGroups"`  // 插件自定义的父级菜单分组（可选）
	FrontendRoutes []FrontendRoute `json:"frontendRoutes"`
}

// MenuGroup 插件声明的菜单分组，安装时自动创建
type MenuGroup struct {
	Code       string `json:"code"`       // 分组权限Code（唯一）
	Name       string `json:"name"`       // 分组显示名称
	Icon       string `json:"icon"`       // 分组图标
	Order      int    `json:"order"`      // 排序
	ParentCode string `json:"parentCode"` // 父分组Code，为空则为顶级
}

// FrontendRoute 前端路由及菜单项信息
type FrontendRoute struct {
	Name       string `json:"name"`       // 组件名称 / 权限代码（唯一），如 TaskManager
	Path       string `json:"path"`       // 前端路由路径，如 /plugins/task-manager
	Component  string `json:"component"`  // 组件相对路径，如 /src/views/plugins/task_manager/index.vue
	Icon       string `json:"icon"`       // 菜单图标，如 i-fe:list
	Title      string `json:"title"`      // 菜单显示标题
	Order      int    `json:"order"`      // 排序
	ParentCode string `json:"parentCode"` // 父菜单权限Code，为空则挂为顶级菜单
}
