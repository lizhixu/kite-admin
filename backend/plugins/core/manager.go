package core

import (
	"backend/config"
	"backend/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	// 定义全局插件注册表
	pluginRegistry = make(map[string]Plugin)
	// 保存传入的 RouterGroup，以便动态注册路由
	routerGroup *gin.RouterGroup
)

// Register 注册插件 (由各个插件在 init() 中调用)
func Register(p Plugin) {
	info := p.Info()
	pluginRegistry[info.Code] = p
}

// InitRouter 保存 RouterGroup 供插件注册路由使用
func InitRouter(r *gin.RouterGroup) {
	routerGroup = r
}

// GetAll 返回所有已注册的插件列表 (附带安装状态)
func GetAll() []map[string]interface{} {
	var result []map[string]interface{}

	var installed []InstalledPlugin
	config.DB.Find(&installed)
	installedMap := make(map[string]InstalledPlugin)
	for _, inst := range installed {
		installedMap[inst.Code] = inst
	}

	for code, plugin := range pluginRegistry {
		info := plugin.Info()
		inst, isInstalled := installedMap[code]

		item := map[string]interface{}{
			"code":        info.Code,
			"name":        info.Name,
			"description": info.Description,
			"version":     info.Version,
			"author":      info.Author,
			"isInstalled": isInstalled,
		}

		if isInstalled {
			item["enable"] = inst.Enable
			item["installedAt"] = inst.InstalledAt
		}

		result = append(result, item)
	}

	return result
}

func ptrStr(s string) *string { return &s }
func ptrUint(i uint) *uint    { return &i }
func ptrBool(b bool) *bool    { return &b }

// Install 安装插件
func Install(code string) error {
	plugin, exists := pluginRegistry[code]
	if !exists {
		return errors.New("plugin not found")
	}

	var count int64
	config.DB.Model(&InstalledPlugin{}).Where("code = ?", code).Count(&count)
	if count > 0 {
		return errors.New("plugin already installed")
	}

	// 1. 执行插件自身的建表逻辑
	if err := plugin.Migrate(config.DB); err != nil {
		return fmt.Errorf("plugin migration failed: %v", err)
	}

	// 2. 先创建插件声明的菜单分组 (MenuGroups)
	info := plugin.Info()
	for _, group := range info.MenuGroups {
		// 如果已存在同名 code 的菜单，则跳过创建
		var existCount int64
		config.DB.Model(&models.Permission{}).Where("code = ?", group.Code).Count(&existCount)
		if existCount > 0 {
			continue
		}
		groupPerm := models.Permission{
			Name:   group.Name,
			Code:   group.Code,
			Type:   "MENU",
			Icon:   ptrStr(group.Icon),
			Order:  group.Order,
			Show:   ptrBool(true),
			Enable: ptrBool(true),
		}
		if group.ParentCode != "" {
			var parentMenu models.Permission
			if err := config.DB.Where("code = ?", group.ParentCode).First(&parentMenu).Error; err == nil {
				groupPerm.ParentID = ptrUint(parentMenu.ID)
			}
		}
		config.DB.Create(&groupPerm)
	}

	// 3. 注入前端路由到权限表 (按照每个 route 声明的 ParentCode 挂载)
	for _, route := range info.FrontendRoutes {
		perm := models.Permission{
			Name:      route.Title,
			Code:      route.Name,
			Type:      "MENU",
			Path:      ptrStr(route.Path),
			Component: ptrStr(route.Component),
			Icon:      ptrStr(route.Icon),
			Order:     route.Order,
			Show:      ptrBool(true),
			Enable:    ptrBool(true),
		}
		// 根据 ParentCode 查找并设置父菜单
		if route.ParentCode != "" {
			var parentMenu models.Permission
			if err := config.DB.Where("code = ?", route.ParentCode).First(&parentMenu).Error; err == nil {
				perm.ParentID = ptrUint(parentMenu.ID)
			}
		}
		// ParentCode 为空 => 不设 ParentID，自动成为顶级菜单
		config.DB.Create(&perm)
	}

	// 3. 注册 API 路由
	if routerGroup != nil {
		plugin.RegisterRoutes(routerGroup)
	}

	// 4. 记录安装状态
	inst := InstalledPlugin{
		Code:   code,
		Enable: true,
	}
	return config.DB.Create(&inst).Error
}

// Uninstall 卸载插件
func Uninstall(code string) error {
	plugin, exists := pluginRegistry[code]
	if !exists {
		return errors.New("plugin not found")
	}

	var inst InstalledPlugin
	if err := config.DB.Where("code = ?", code).First(&inst).Error; err != nil {
		return errors.New("plugin not installed")
	}

	// 1. 卸载数据库结构
	if err := plugin.Uninstall(config.DB); err != nil {
		return fmt.Errorf("plugin uninstall failed: %v", err)
	}

	// 2. 移除相关菜单权限 (先删子菜单，再删分组)
	info := plugin.Info()
	for _, route := range info.FrontendRoutes {
		config.DB.Where("code = ?", route.Name).Delete(&models.Permission{})
	}
	for _, group := range info.MenuGroups {
		// 只删除该插件自己创建的分组，防止删除共享分组
		// （检查分组下是否还有其他菜单，如没有才删除）
		var childCount int64
		var groupPerm models.Permission
		if err := config.DB.Where("code = ?", group.Code).First(&groupPerm).Error; err == nil {
			config.DB.Model(&models.Permission{}).Where("parent_id = ?", groupPerm.ID).Count(&childCount)
			if childCount == 0 {
				config.DB.Delete(&groupPerm)
			}
		}
	}

	// 3. 删除安装记录
	return config.DB.Delete(&inst).Error
}

// ToggleEnable 启用/禁用插件
func ToggleEnable(code string, enable bool) error {
	var inst InstalledPlugin
	if err := config.DB.Where("code = ?", code).First(&inst).Error; err != nil {
		return errors.New("plugin not installed")
	}

	inst.Enable = enable
	return config.DB.Save(&inst).Error
}

// InitInstalled 初始化已安装的插件路由 (系统启动时调用)
func InitInstalled() {
	var installed []InstalledPlugin
	config.DB.Where("enable = ?", true).Find(&installed)

	for _, inst := range installed {
		if plugin, exists := pluginRegistry[inst.Code]; exists {
			if routerGroup != nil {
				plugin.RegisterRoutes(routerGroup)
			}
		}
	}
}
