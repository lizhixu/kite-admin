package task_manager

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

// GetList 获取任务列表
func (c *Controller) GetList(ctx *gin.Context) {
	var tasks []Task
	config.DB.Order("id desc").Find(&tasks)

	ctx.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "OK",
		Data:    tasks,
	})
}

// Create 创建任务
func (c *Controller) Create(ctx *gin.Context) {
	var req Task
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: "参数错误"})
		return
	}

	req.Status = "TODO"
	if err := config.DB.Create(&req).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Code: 0, Message: "创建成功", Data: req})
}

// Update 更新任务
func (c *Controller) UpdateStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: "参数错误"})
		return
	}

	if err := config.DB.Model(&Task{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Code: 0, Message: "更新成功"})
}

// Delete 删除任务
func (c *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := config.DB.Delete(&Task{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Code: 0, Message: "删除成功"})
}
