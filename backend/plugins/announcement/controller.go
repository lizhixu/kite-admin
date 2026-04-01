package announcement

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

// GetList 获取公告列表
func (c *Controller) GetList(ctx *gin.Context) {
	var list []Announcement
	config.DB.Order("is_top desc, id desc").Find(&list)

	ctx.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "OK",
		Data:    list,
	})
}

// Create 发布公告
func (c *Controller) Create(ctx *gin.Context) {
	var req Announcement
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: "参数错误"})
		return
	}

	if err := config.DB.Create(&req).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "发布失败"})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Code: 0, Message: "发布成功", Data: req})
}

// Delete 删除公告
func (c *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := config.DB.Delete(&Announcement{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Code: 0, Message: "删除成功"})
}
