package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())

	authCtrl := &controllers.AuthController{}
	userCtrl := &controllers.UserController{}
	roleCtrl := &controllers.RoleController{}
	permCtrl := &controllers.PermissionController{}

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
	{
		// 认证相关
		api.POST("/auth/current-role/switch/:roleCode", authCtrl.SwitchRole)

		// 用户相关
		api.GET("/user/detail", userCtrl.GetDetail)
		api.GET("/user", userCtrl.GetList)
		api.POST("/user", userCtrl.Create)
		api.DELETE("/user/:id", userCtrl.Delete)
		api.PATCH("/user/:id", userCtrl.Update)
		api.PATCH("/user/password/reset/:id", userCtrl.ResetPassword)

		// 角色相关
		api.GET("/role/page", roleCtrl.GetPage)
		api.GET("/role", roleCtrl.GetAll)
		api.POST("/role", roleCtrl.Create)
		api.PATCH("/role/:id", roleCtrl.Update)
		api.DELETE("/role/:id", roleCtrl.Delete)
		api.PATCH("/role/users/add/:id", roleCtrl.AddUsers)
		api.PATCH("/role/users/remove/:id", roleCtrl.RemoveUsers)

		// 权限相关
		api.GET("/role/permissions/tree", permCtrl.GetRolePermissionsTree)
		api.GET("/permission/menu/tree", permCtrl.GetMenuTree)
		api.GET("/permission/tree", permCtrl.GetTree)
		api.GET("/permission/button/:parentId", permCtrl.GetButtonsByParentID)
		api.POST("/permission", permCtrl.Create)
		api.PATCH("/permission/:id", permCtrl.Update)
		api.DELETE("/permission/:id", permCtrl.Delete)
	}
}
