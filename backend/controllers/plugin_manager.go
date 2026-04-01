package controllers

import (
	"backend/models"
	"backend/plugins/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PluginController struct{}

// GetList 获取所有插件列表
func (pc *PluginController) GetList(c *gin.Context) {
	plugins := core.GetAll()
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "OK",
		Data:    plugins,
	})
}

// Install 安装插件
func (pc *PluginController) Install(c *gin.Context) {
	code := c.Param("code")
	if err := core.Install(code); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "安装成功",
	})
}

// Uninstall 卸载插件
func (pc *PluginController) Uninstall(c *gin.Context) {
	code := c.Param("code")
	if err := core.Uninstall(code); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "卸载成功",
	})
}

// ToggleEnable 启用/禁用插件
func (pc *PluginController) ToggleEnable(c *gin.Context) {
	code := c.Param("code")
	var req struct {
		Enable bool `json:"enable"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	if err := core.ToggleEnable(code, req.Enable); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "操作成功",
	})
}
