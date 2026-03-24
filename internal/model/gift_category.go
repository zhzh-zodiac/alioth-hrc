package model

import "time"

// GiftCategory 系统分类 UserID 为 nil；用户自定义分类 UserID 指向所属用户。
type GiftCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	IsSystem  bool      `gorm:"not null;default:false" json:"is_system"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
