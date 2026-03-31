package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (uc *UserController) GetDetail(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.User
	if err := config.DB.Preload("Profile").Preload("Roles").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:      404,
			Message:   "User not found",
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	roleCode := c.GetString("roleCode")
	var currentRole *models.Role
	for _, role := range user.Roles {
		if role.Code == roleCode {
			currentRole = &role
			break
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "OK",
		Data: gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"enable":      user.Enable,
			"createTime":  user.CreateTime,
			"updateTime":  user.UpdateTime,
			"profile":     user.Profile,
			"roles":       user.Roles,
			"currentRole": currentRole,
		},
		OriginUrl: c.Request.URL.Path,
	})
}

func (uc *UserController) GetList(c *gin.Context) {
	username := c.Query("username")
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))

	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	query := config.DB.Model(&models.User{}).Preload("Roles").Preload("Profile")
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	var total int64
	query.Count(&total)

	var users []models.User
	offset := (pageNo - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Find(&users)

	var pageData []gin.H
	for _, user := range users {
		userData := gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"enable":     user.Enable,
			"createTime": user.CreateTime,
			"updateTime": user.UpdateTime,
			"roles":      user.Roles,
		}
		if user.Profile != nil {
			userData["gender"] = user.Profile.Gender
			userData["avatar"] = user.Profile.Avatar
			userData["address"] = user.Profile.Address
			userData["email"] = user.Profile.Email
		}
		pageData = append(pageData, userData)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "OK",
		Data: models.PageData{
			PageData: pageData,
			Total:    total,
		},
		OriginUrl: c.Request.URL.String(),
	})
}

func (uc *UserController) Create(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:      400,
			Message:   err.Error(),
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:      500,
			Message:   "Failed to hash password",
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Enable:   true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:      500,
			Message:   "Failed to create user",
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:      0,
		Message:   "OK",
		Data:      user,
		OriginUrl: c.Request.URL.Path,
	})
}

func (uc *UserController) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Enable *bool `json:"enable"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:      400,
			Message:   err.Error(),
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	if req.Enable != nil {
		if err := config.DB.Model(&models.User{}).Where("id = ?", id).Update("enable", *req.Enable).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:      500,
				Message:   "Failed to update user",
				OriginUrl: c.Request.URL.Path,
			})
			return
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:      0,
		Message:   "OK",
		Data:      true,
		OriginUrl: c.Request.URL.Path,
	})
}

func (uc *UserController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:      500,
			Message:   "Failed to delete user",
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:      0,
		Message:   "OK",
		Data:      true,
		OriginUrl: c.Request.URL.Path,
	})
}

func (uc *UserController) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:      400,
			Message:   err.Error(),
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:      500,
			Message:   "Failed to hash password",
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:      500,
			Message:   "Failed to reset password",
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:      0,
		Message:   "OK",
		Data:      true,
		OriginUrl: c.Request.URL.Path,
	})
}
