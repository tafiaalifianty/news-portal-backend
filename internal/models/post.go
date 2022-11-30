package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID           int64    `gorm:"primary_key"`
	Slug         string   `json:"slug"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Summary      string   `json:"summary"`
	CategoryID   int64    `json:"category_id"`
	Category     Category `json:"category" gorm:"foreignKey:category_id"`
	TypeID       int64    `json:"type_id"`
	Type         PostType `json:"type" gorm:"foreignKey:type_id"`
	ImgThumbnail string   `json:"img_thumbnail"`
	ImgUrl       string   `json:"img_url"`
	AuthorName   string   `json:"author_name"`
	ShareCount   int      `json:"share_count"`
	LikeCount    int      `json:"like_count"`
}

type TrendingPost struct {
	Post
	ReaderCount int   `json:"reader_count"`
	PostID      int64 `json:"post_id"`
}
