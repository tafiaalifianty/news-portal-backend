package models

import (
	"time"

	"gorm.io/gorm"
)

type History struct {
	UserID       int64     `json:"user_id" gorm:"primaryKey"`
	PostID       int64     `json:"post_id" gorm:"primaryKey"`
	Post         Post      `json:"post" gorm:"foreignKey:post_id"`
	LastAccessed time.Time `json:"last_accessed" gorm:"autoUpdateTime:true"`
	IsLiked      bool      `json:"is_liked"`
	IsShared     bool      `json:"is_shared"`
	DeletedAt    gorm.DeletedAt
	CreatedAt    time.Time
}

func (History) BeforeCreate(db *gorm.DB) error {
	return nil
}
