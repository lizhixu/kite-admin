package announcement

import "time"

// Announcement 插件演示数据模型: 系统公告
type Announcement struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	Type       string    `gorm:"default:'INFO'" json:"type"` // INFO, WARNING, URGENT
	IsTop      bool      `gorm:"default:false" json:"isTop"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
}
