package task_manager

import (
	"backend/plugins/core"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	// 在系统启动时自动注册该插件
	core.Register(&TaskManagerPlugin{})
}

type TaskManagerPlugin struct{}

func (p *TaskManagerPlugin) Info() core.PluginInfo {
	return core.PluginInfo{
		Code:        "task_manager",
		Name:        "任务管理",
		Description: "任务看板、进度跟踪的演示插件",
		Version:     "1.0.0",
		Author:      "kite",
		// 声明独立的一级菜单分组"业务应用"
		MenuGroups: []core.MenuGroup{
			{
				Code:  "BusinessApps",
				Name:  "业务应用",
				Icon:  "i-fe:grid",
				Order: 10,
			},
		},
		FrontendRoutes: []core.FrontendRoute{
			{
				Name:       "TaskManager",
				Path:       "/plugins/task-manager",
				Component:  "/src/views/plugins/task_manager/index.vue",
				Icon:       "i-fe:check-square",
				Title:      "任务管理",
				Order:      1,
				ParentCode: "BusinessApps", // 挂在"业务应用"分组下
			},
		},
	}
}

func (p *TaskManagerPlugin) Migrate(db *gorm.DB) error {
	// 插件安装时自动建表
	return db.AutoMigrate(&Task{})
}

func (p *TaskManagerPlugin) RegisterRoutes(group *gin.RouterGroup) {
	// 插件注册自己的API路由
	ctrl := &Controller{}
	tasks := group.Group("/plugins/tasks")
	{
		tasks.GET("", ctrl.GetList)
		tasks.POST("", ctrl.Create)
		tasks.PATCH("/:id/status", ctrl.UpdateStatus)
		tasks.DELETE("/:id", ctrl.Delete)
	}
}

func (p *TaskManagerPlugin) Uninstall(db *gorm.DB) error {
	// 插件卸载时删表
	return db.Migrator().DropTable(&Task{})
}
