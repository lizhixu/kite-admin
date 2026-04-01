package announcement

import (
	"backend/plugins/core"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	core.Register(&AnnouncementPlugin{})
}

type AnnouncementPlugin struct{}

func (p *AnnouncementPlugin) Info() core.PluginInfo {
	return core.PluginInfo{
		Code:        "announcement",
		Name:        "系统公告",
		Description: "系统全站公告管理插件",
		Version:     "1.0.0",
		Author:      "kite",
		// 声明同一个"业务应用"分组（若已被其他插件创建，则自动跳过）
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
				Name:       "Announcement",
				Path:       "/plugins/announcement",
				Component:  "/src/views/plugins/announcement/index.vue",
				Icon:       "i-fe:bell",
				Title:      "系统公告",
				Order:      2,
				ParentCode: "BusinessApps",
			},
		},
	}
}

func (p *AnnouncementPlugin) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Announcement{})
}

func (p *AnnouncementPlugin) RegisterRoutes(group *gin.RouterGroup) {
	ctrl := &Controller{}
	api := group.Group("/plugins/announcement")
	{
		api.GET("", ctrl.GetList)
		api.POST("", ctrl.Create)
		api.DELETE("/:id", ctrl.Delete)
	}
}

func (p *AnnouncementPlugin) Uninstall(db *gorm.DB) error {
	return db.Migrator().DropTable(&Announcement{})
}
