package routes

import (
	"backend/controllers"
	"backend/middleware"
	"backend/plugins/core"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())

	authCtrl := &controllers.AuthController{}
	userCtrl := &controllers.UserController{}
	roleCtrl := &controllers.RoleController{}
	permCtrl := &controllers.PermissionController{}
	pluginCtrl := &controllers.PluginController{}
	attachCtrl := &controllers.AttachmentController{}

	// 静态文件：本地上传附件
	r.Static("/uploads", "./uploads")

	// 认证相关路由（无需认证）
	auth := r.Group("/auth")
	{
		auth.POST("/login", authCtrl.Login)
		auth.GET("/captcha", authCtrl.GetCaptcha)
		auth.POST("/logout", authCtrl.Logout)
	}

	// 需要认证的路由
	api := r.Group("")
	api.Use(middleware.AuthMiddleware())

	// 初始化插件路由注册器
	core.InitRouter(api)
	{
		// 认证相关
		api.POST("/auth/current-role/switch/:roleCode", authCtrl.SwitchRole)

		// 用户相关
		api.GET("/user/detail", userCtrl.GetDetail)
		api.GET("/user", userCtrl.GetList)
		api.POST("/user", middleware.RequirePermission("AddUser"), userCtrl.Create)
		api.DELETE("/user/:id", middleware.RequirePermission("DeleteUser"), userCtrl.Delete)
		api.PATCH("/user/:id", middleware.RequirePermission("EditUser"), userCtrl.Update)
		api.PATCH("/user/password/reset/:id", middleware.RequirePermission("ResetPassword"), userCtrl.ResetPassword)

		// 角色相关
		api.GET("/role/page", roleCtrl.GetPage)
		api.GET("/role", roleCtrl.GetAll)
		api.POST("/role", middleware.RequirePermission("AddRole"), roleCtrl.Create)
		api.PATCH("/role/:id", middleware.RequirePermission("EditRole"), roleCtrl.Update)
		api.DELETE("/role/:id", middleware.RequirePermission("DeleteRole"), roleCtrl.Delete)
		api.PATCH("/role/users/add/:id", middleware.RequirePermission("AssignPermission"), roleCtrl.AddUsers)
		api.PATCH("/role/users/remove/:id", middleware.RequirePermission("AssignPermission"), roleCtrl.RemoveUsers)

		// 权限相关
		api.GET("/role/permissions/tree", permCtrl.GetRolePermissionsTree)
		api.GET("/permission/menu/tree", permCtrl.GetMenuTree)
		api.GET("/permission/menu/validate", permCtrl.ValidateMenuPath)
		api.GET("/permission/tree", permCtrl.GetTree)
		api.GET("/permission/button/:parentId", permCtrl.GetButtonsByParentID)
		api.POST("/permission", middleware.RequirePermission("AddResource"), permCtrl.Create)
		api.PATCH("/permission/:id", middleware.RequirePermission("EditResource"), permCtrl.Update)
		api.DELETE("/permission/:id", middleware.RequirePermission("DeleteResource"), permCtrl.Delete)

		// 插件管理
		api.GET("/plugin/list", pluginCtrl.GetList)
		api.POST("/plugin/install/:code", middleware.RequirePermission("InstallPlugin"), pluginCtrl.Install)
		api.POST("/plugin/uninstall/:code", middleware.RequirePermission("UninstallPlugin"), pluginCtrl.Uninstall)
		api.PATCH("/plugin/enable/:code", middleware.RequirePermission("EnablePlugin"), pluginCtrl.ToggleEnable)

		// 附件管理
		api.GET("/attachment", attachCtrl.GetList)
		api.POST("/attachment/upload", attachCtrl.Upload)
		api.DELETE("/attachment/:id", middleware.RequirePermission("DeleteAttachment"), attachCtrl.Delete)
		api.GET("/attachment/config", attachCtrl.GetConfig)
		api.POST("/attachment/config", middleware.RequirePermission("ManageAttachmentConfig"), attachCtrl.SaveConfig)

		// 附件分组
		api.GET("/attachment/group", attachCtrl.GetGroups)
		api.POST("/attachment/group", attachCtrl.CreateGroup)
		api.PATCH("/attachment/group/:id", attachCtrl.UpdateGroup)
		api.DELETE("/attachment/group/:id", attachCtrl.DeleteGroup)
	}
}
