package controllers

import (
	"backend/config"
	"backend/models"
	"backend/storage"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AttachmentController struct{}

// getStorageConfig 从数据库获取配置（不存在则返回默认值）
func getStorageConfig() models.StorageConfig {
	var cfg models.StorageConfig
	// 用 Find+Limit 而非 First，避免记录不存在时 GORM 打印 record not found 日志
	config.DB.Where("id = ?", 1).Limit(1).Find(&cfg)
	if cfg.ID == 0 {
		cfg = models.StorageConfig{
			ID:          1,
			StorageType: "local",
			LocalPath:   "uploads",
			MaxSizeMB:   10,
		}
	}
	return cfg
}

// buildDriver 根据配置构造存储驱动
func buildDriver(cfg models.StorageConfig, baseURL string) storage.Driver {
	if cfg.StorageType == "s3" {
		return &storage.S3Driver{
			Config: storage.S3Config{
				Endpoint:  cfg.S3Endpoint,
				Bucket:    cfg.S3Bucket,
				Region:    cfg.S3Region,
				AccessKey: cfg.S3AccessKey,
				SecretKey: cfg.S3SecretKey,
				PublicURL: cfg.S3PublicURL,
			},
		}
	}
	localPath := cfg.LocalPath
	if localPath == "" {
		localPath = "uploads"
	}
	return &storage.LocalDriver{
		BasePath: localPath,
		BaseURL:  baseURL,
	}
}

// GetConfig GET /attachment/config
func (ac *AttachmentController) GetConfig(c *gin.Context) {
	cfg := getStorageConfig()
	c.JSON(http.StatusOK, models.Response{Code: 0, Message: "OK", Data: cfg, OriginUrl: c.Request.URL.Path})
}

// SaveConfig POST /attachment/config
func (ac *AttachmentController) SaveConfig(c *gin.Context) {
	var req models.StorageConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: err.Error(), OriginUrl: c.Request.URL.Path})
		return
	}
	req.ID = 1 // 单条记录
	result := config.DB.Save(&req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "保存失败", OriginUrl: c.Request.URL.Path})
		return
	}
	c.JSON(http.StatusOK, models.Response{Code: 0, Message: "OK", Data: req, OriginUrl: c.Request.URL.Path})
}

// Upload POST /attachment/upload
func (ac *AttachmentController) Upload(c *gin.Context) {
	cfg := getStorageConfig()

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: "未收到文件", OriginUrl: c.Request.URL.Path})
		return
	}
	defer file.Close()

	// 检查文件大小
	if cfg.MaxSizeMB > 0 && header.Size > int64(cfg.MaxSizeMB)*1024*1024 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:      400,
			Message:   fmt.Sprintf("文件超过最大限制 %dMB", cfg.MaxSizeMB),
			OriginUrl: c.Request.URL.Path,
		})
		return
	}

	// 检查文件类型
	mimeType := header.Header.Get("Content-Type")
	if cfg.AllowedTypes != "" {
		allowed := strings.Split(cfg.AllowedTypes, ",")
		ok := false
		for _, a := range allowed {
			if strings.HasPrefix(mimeType, strings.TrimSpace(a)) {
				ok = true
				break
			}
		}
		if !ok {
			c.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: "文件类型不允许", OriginUrl: c.Request.URL.Path})
			return
		}
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	uniqueName := fmt.Sprintf("%d_%s%s", time.Now().UnixMilli(), sanitizeFileName(strings.TrimSuffix(header.Filename, ext)), ext)

	// 构建访问基础 URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + c.Request.Host

	driver := buildDriver(cfg, baseURL)
	storagePath, url, err := driver.Upload(file, uniqueName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "上传失败: " + err.Error(), OriginUrl: c.Request.URL.Path})
		return
	}

	// 获取上传者 ID
	var uploaderID uint
	if id, exists := c.Get("userId"); exists {
		if uid, ok := id.(uint); ok {
			uploaderID = uid
		}
	}

	// 获取 groupId
	var groupID uint
	if gid := c.PostForm("groupId"); gid != "" {
		fmt.Sscanf(gid, "%d", &groupID)
	}

	attach := models.Attachment{
		GroupID:      groupID,
		FileName:     uniqueName,
		OriginalName: header.Filename,
		FileSize:     header.Size,
		MimeType:     mimeType,
		StorageType:  cfg.StorageType,
		StoragePath:  storagePath,
		Url:          url,
		UploaderID:   uploaderID,
	}
	if err := config.DB.Create(&attach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "记录写入失败", OriginUrl: c.Request.URL.Path})
		return
	}
	c.JSON(http.StatusOK, models.Response{Code: 0, Message: "OK", Data: attach, OriginUrl: c.Request.URL.Path})
}

// GetList GET /attachment
func (ac *AttachmentController) GetList(c *gin.Context) {
	keyword := c.Query("keyword")
	groupIdStr := c.Query("groupId")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	var page, pageSize int
	fmt.Sscanf(pageStr, "%d", &page)
	fmt.Sscanf(pageSizeStr, "%d", &pageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := config.DB.Model(&models.Attachment{})
	if keyword != "" {
		db = db.Where("original_name LIKE ?", "%"+keyword+"%")
	}
	if groupIdStr != "" {
		var gid int
		if _, err := fmt.Sscanf(groupIdStr, "%d", &gid); err == nil {
			db = db.Where("group_id = ?", gid)
		}
	}

	var total int64
	db.Count(&total)

	var list []models.Attachment
	db.Preload("Uploader").Order("create_time DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "OK",
		Data: models.PageData{
			PageData: list,
			Total:    total,
		},
		OriginUrl: c.Request.URL.Path,
	})
}

// Delete DELETE /attachment/:id
func (ac *AttachmentController) Delete(c *gin.Context) {
	id := c.Param("id")
	var attach models.Attachment
	if err := config.DB.First(&attach, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{Code: 404, Message: "附件不存在", OriginUrl: c.Request.URL.Path})
		return
	}

	// 先删物理文件
	cfg := getStorageConfig()
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + c.Request.Host
	driver := buildDriver(cfg, baseURL)
	_ = driver.Delete(attach.StoragePath) // 物理删除失败不阻塞 DB 删除

	if err := config.DB.Delete(&attach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "删除失败", OriginUrl: c.Request.URL.Path})
		return
	}
	c.JSON(http.StatusOK, models.Response{Code: 0, Message: "OK", Data: true, OriginUrl: c.Request.URL.Path})
}

// ---------------- 分组管理 ----------------

// GetGroups GET /attachment/group
func (ac *AttachmentController) GetGroups(c *gin.Context) {
	var list []models.AttachmentGroup
	config.DB.Order("`order` ASC, create_time ASC").Find(&list)
	c.JSON(http.StatusOK, models.Response{Code: 0, Message: "OK", Data: list, OriginUrl: c.Request.URL.Path})
}

// CreateGroup POST /attachment/group
func (ac *AttachmentController) CreateGroup(c *gin.Context) {
	var req models.AttachmentGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: err.Error(), OriginUrl: c.Request.URL.Path})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: "分组名不能为空", OriginUrl: c.Request.URL.Path})
		return
	}
	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "创建失败", OriginUrl: c.Request.URL.Path})
		return
	}
	c.JSON(http.StatusOK, models.Response{Code: 0, Message: "OK", Data: req, OriginUrl: c.Request.URL.Path})
}

// UpdateGroup PATCH /attachment/group/:id
func (ac *AttachmentController) UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var req models.AttachmentGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Code: 400, Message: err.Error(), OriginUrl: c.Request.URL.Path})
		return
	}
	var group models.AttachmentGroup
	if err := config.DB.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{Code: 404, Message: "分组不存在", OriginUrl: c.Request.URL.Path})
		return
	}
	if req.Name != "" {
		group.Name = req.Name
	}
	group.Order = req.Order
	config.DB.Save(&group)
	c.JSON(http.StatusOK, models.Response{Code: 0, Message: "OK", Data: group, OriginUrl: c.Request.URL.Path})
}

// DeleteGroup DELETE /attachment/group/:id
func (ac *AttachmentController) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	// 将该分组内的附件移动到"未分组" (GroupID=0)
	config.DB.Model(&models.Attachment{}).Where("group_id = ?", id).Update("group_id", 0)

	if err := config.DB.Delete(&models.AttachmentGroup{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Message: "删除失败", OriginUrl: c.Request.URL.Path})
		return
	}
	c.JSON(http.StatusOK, models.Response{Code: 0, Message: "OK", Data: true, OriginUrl: c.Request.URL.Path})
}

// sanitizeFileName 移除文件名中的危险字符
func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "..", "_")
	if len(name) > 50 {
		name = name[:50]
	}
	return name
}
