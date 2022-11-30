package dtos

import (
	"time"

	"final-project-backend/internal/models"
)

type CreatePostRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary" binding:"required"`
	CategoryID int64  `json:"category_id" binding:"required"`
	TypeID     int64  `json:"type_id" binding:"required"`
	ImgUrl     string `json:"img_url" binding:"required"`
	AuthorName string `json:"author_name" binding:"required"`
}

type CreatePostResponse struct {
	ID int64 `json:"id"`
}

type PostResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Content      string    `json:"content,omitempty"`
	Summary      string    `json:"summary"`
	CategoryID   int64     `json:"category_id,omitempty"`
	Category     string    `json:"category"`
	TypeID       int64     `json:"type_id,omitempty"`
	Type         string    `json:"type"`
	ImgUrl       string    `json:"img_url,omitempty"`
	ImgThumbnail string    `json:"img_thumbnail"`
	AuthorName   string    `json:"author_name"`
	ShareCount   int       `json:"share_count"`
	LikeCount    int       `json:"like_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type PostResponseCompact struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Summary      string    `json:"summary"`
	Category     string    `json:"category"`
	Type         string    `json:"type"`
	ImgThumbnail string    `json:"img_thumbnail"`
	AuthorName   string    `json:"author_name"`
	ShareCount   int       `json:"share_count"`
	LikeCount    int       `json:"like_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type GetPostByIDWithHistoryResponse struct {
	PostResponse
	HistoryResponseDTO
}

type GetAllPostResponse struct {
	Data []*PostResponseCompact `json:"data"`
	PaginationResponse
}

type GetAllTypesResponse = []*models.PostType

type GetAllCategoriesResponse = []*models.Category

type LikePostRequest struct {
	IsLike *bool `json:"is_like" binding:"required"`
}

type PostsRequestQuery struct {
	Search     string
	CategoryID int64
	TypeID     int64
	Sort       string
	Limit      int
	Page       int
}

func FormatPost(post *models.Post) *PostResponse {
	return &PostResponse{
		ID:           post.ID,
		Title:        post.Title,
		Slug:         post.Slug,
		Content:      post.Content,
		Summary:      post.Summary,
		CategoryID:   post.CategoryID,
		Category:     post.Category.Name,
		TypeID:       post.TypeID,
		Type:         post.Type.Name,
		ImgUrl:       post.ImgUrl,
		ImgThumbnail: post.ImgThumbnail,
		AuthorName:   post.AuthorName,
		CreatedAt:    post.Model.CreatedAt,
		LikeCount:    post.LikeCount,
		ShareCount:   post.ShareCount,
	}
}

func FormatPostCompact(post *models.Post) *PostResponseCompact {
	return &PostResponseCompact{
		ID:           post.ID,
		Title:        post.Title,
		Slug:         post.Slug,
		Summary:      post.Summary,
		Category:     post.Category.Name,
		Type:         post.Type.Name,
		ImgThumbnail: post.ImgThumbnail,
		AuthorName:   post.AuthorName,
		CreatedAt:    post.Model.CreatedAt,
		LikeCount:    post.LikeCount,
		ShareCount:   post.ShareCount,
	}
}

func FormatPostsCompact(posts []*models.Post) []*PostResponseCompact {
	formattedPosts := []*PostResponseCompact{}
	for _, post := range posts {
		formattedPost := FormatPostCompact(post)
		formattedPosts = append(formattedPosts, formattedPost)
	}
	return formattedPosts
}
