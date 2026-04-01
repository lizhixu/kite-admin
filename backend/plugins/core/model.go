package core

import "time"

// InstalledPlugin 记录已安装的插件
type InstalledPlugin struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"unique;not null" json:"code"`
	Enable      bool      `gorm:"default:true" json:"enable"`
	InstalledAt time.Time `gorm:"autoCreateTime" json:"installedAt"`
}
