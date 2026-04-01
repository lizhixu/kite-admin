package task_manager

import "time"

// Task 插件演示数据模型: 任务
type Task struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title"`
	Content    string    `json:"content"`
	Status     string    `gorm:"default:'TODO'" json:"status"` // TODO, DOING, DONE
	Assignee   string    `json:"assignee"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
}
