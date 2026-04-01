package models

import "time"

// StorageConfig 存储配置（全局单条记录，ID 固定为 1）
type StorageConfig struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	StorageType  string `gorm:"default:'local'" json:"storageType"` // local | s3
	LocalPath    string `gorm:"default:'uploads'" json:"localPath"`
	MaxSizeMB    int    `gorm:"default:10" json:"maxSizeMB"`    // 单文件最大 MB，0=不限
	AllowedTypes string `gorm:"default:''" json:"allowedTypes"` // 逗号分隔的 MIME 前缀，空=不限
	// S3 配置
	S3Endpoint  string `gorm:"default:''" json:"s3Endpoint"`
	S3Bucket    string `gorm:"default:''" json:"s3Bucket"`
	S3Region    string `gorm:"default:''" json:"s3Region"`
	S3AccessKey string `gorm:"default:''" json:"s3AccessKey"`
	S3SecretKey string `gorm:"default:''" json:"s3SecretKey"`
	S3PublicURL string `gorm:"default:''" json:"s3PublicUrl"` // 对外访问的基础 URL
}

// AttachmentGroup 附件分组
type AttachmentGroup struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Order      int       `gorm:"default:0" json:"order"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
}

// Attachment 附件记录
type Attachment struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	GroupID      uint      `gorm:"default:0;index" json:"groupId"` // 0=未分组
	FileName     string    `gorm:"not null" json:"fileName"`       // 存储后的文件名（含唯一前缀）
	OriginalName string    `gorm:"not null" json:"originalName"`   // 原始文件名
	FileSize     int64     `json:"fileSize"`                       // 字节数
	MimeType     string    `json:"mimeType"`
	StorageType  string    `json:"storageType"` // local | s3
	StoragePath  string    `json:"storagePath"` // 存储路径或 S3 key
	Url          string    `json:"url"`         // 可访问 URL
	UploaderID   uint      `json:"uploaderId"`
	Uploader     *User     `gorm:"foreignKey:UploaderID;constraint:false" json:"uploader"` // 关联的用户信息（不建物理外键，避免历史脏数据报错）
	CreateTime   time.Time `gorm:"autoCreateTime" json:"createTime"`
}
